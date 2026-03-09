import amqplib from "amqplib";
import path from "path";
import { spawn } from "child_process";
import fs from "fs/promises";

const QUEUE_NAME = process.env.QUEUE_NAME || "video_processing";
const RABBITMQ_URL = process.env.RABBITMQ_URL || "amqp://rabbitmq:5672";
const RETRY_DELAY_MS = 5000;

type JobOperation =
  | "BLACK_AND_WHITE"
  | "REMOVE_AUDIO"
  | "EXTRACT_AUDIO"
  | "FRAME_EXTRACT"
  | "TRIM_30S"
  | "COMPRESS_VIDEO";

type JobPayload = {
  filePath: string;
  fileName: string;
  operation?: string;
};

async function sleep(ms: number) {
  await new Promise((resolve) => setTimeout(resolve, ms));
}

function normalizeOperation(operation?: string): JobOperation {
  if (operation === "REMOVE_AUDIO") return operation;
  if (operation === "EXTRACT_AUDIO") return operation;
  if (operation === "FRAME_EXTRACT") return operation;
  if (operation === "TRIM_30S") return operation;
  if (operation === "COMPRESS_VIDEO") return operation;
  return "BLACK_AND_WHITE";
}

function outputExtensionForOperation(operation: JobOperation) {
  if (operation === "EXTRACT_AUDIO") return "mp3";
  if (operation === "FRAME_EXTRACT") return "png";
  return "mp4";
}

function buildFfmpegArgs(operation: JobOperation, inputPath: string, outputPath: string) {
  switch (operation) {
    case "EXTRACT_AUDIO":
      return [
        "-hide_banner",
        "-i",
        inputPath,
        "-vn",
        "-c:a",
        "libmp3lame",
        "-b:a",
        "192k",
        "-ar",
        "44100",
        "-y",
        outputPath,
      ];

    case "REMOVE_AUDIO":
      return [
        "-hide_banner",
        "-i",
        inputPath,
        "-map",
        "0:v:0",
        "-c:v",
        "libx264",
        "-preset",
        "medium",
        "-crf",
        "20",
        "-pix_fmt",
        "yuv420p",
        "-tag:v",
        "avc1",
        "-an",
        "-movflags",
        "+faststart",
        "-y",
        outputPath,
      ];

    case "FRAME_EXTRACT":
      return [
        "-hide_banner",
        "-i",
        inputPath,
        "-vf",
        "fps=1",
        "-vsync",
        "vfr",
        "-c:v",
        "png",
        "-y",
        outputPath,
      ];

    case "TRIM_30S":
      return [
        "-hide_banner",
        "-i",
        inputPath,
        "-t",
        "30",
        "-map",
        "0:v:0",
        "-map",
        "0:a:0?",
        "-c:v",
        "libx264",
        "-preset",
        "medium",
        "-crf",
        "21",
        "-pix_fmt",
        "yuv420p",
        "-tag:v",
        "avc1",
        "-c:a",
        "aac",
        "-b:a",
        "128k",
        "-movflags",
        "+faststart",
        "-y",
        outputPath,
      ];

    case "COMPRESS_VIDEO":
      return [
        "-hide_banner",
        "-i",
        inputPath,
        "-vf",
        "scale=1280:-2:force_original_aspect_ratio=decrease",
        "-map",
        "0:v:0",
        "-map",
        "0:a:0?",
        "-c:v",
        "libx264",
        "-preset",
        "slow",
        "-crf",
        "28",
        "-pix_fmt",
        "yuv420p",
        "-tag:v",
        "avc1",
        "-c:a",
        "aac",
        "-b:a",
        "96k",
        "-movflags",
        "+faststart",
        "-y",
        outputPath,
      ];

    case "BLACK_AND_WHITE":
    default:
      return [
        "-hide_banner",
        "-i",
        inputPath,
        "-map",
        "0:v:0",
        "-map",
        "0:a:0?",
        "-vf",
        "hue=s=0",
        "-c:v",
        "libx264",
        "-preset",
        "medium",
        "-crf",
        "20",
        "-pix_fmt",
        "yuv420p",
        "-tag:v",
        "avc1",
        "-c:a",
        "aac",
        "-b:a",
        "128k",
        "-movflags",
        "+faststart",
        "-y",
        outputPath,
      ];
  }
}

function runFfmpeg(args: string[]) {
  return new Promise<void>((resolve, reject) => {
    const proc = spawn("ffmpeg", args);

    proc.stdout.on("data", (data) => {
      process.stdout.write(data.toString());
    });

    proc.stderr.on("data", (data) => {
      process.stdout.write(data.toString());
    });

    proc.on("error", (error) => {
      reject(new Error(`Failed to spawn ffmpeg: ${error.message}`));
    });

    proc.on("close", (code) => {
      if (code === 0) return resolve();
      reject(new Error(`ffmpeg exited with code ${code}`));
    });
  });
}

function parseJob(raw: string): JobPayload {
  const parsed = JSON.parse(raw) as JobPayload;

  if (!parsed.filePath || !parsed.fileName) {
    throw new Error("Invalid payload: missing file path or file name");
  }

  return parsed;
}

async function startWorker() {
  while (true) {
    let connection: any;
    let channel: any;

    try {
      connection = await amqplib.connect(RABBITMQ_URL);
      channel = await connection.createChannel();

      await channel.assertQueue(QUEUE_NAME, { durable: true });
      channel.prefetch(1);

      console.log(`[*] Worker running and listening on '${QUEUE_NAME}'`);

      channel.consume(
        QUEUE_NAME,
        async (msg: any) => {
          if (!msg) return;

          let tempFilePath: string | null = null;
          let tempFramesDirPath: string | null = null;

          try {
            const job = parseJob(msg.content.toString());
            const operation = normalizeOperation(job.operation);
            const baseName = path.parse(job.fileName).name;
            const outputDir = path.dirname(job.filePath);
            const ext = outputExtensionForOperation(operation);
            const finalOutputPath = path.resolve(
              outputDir,
              `processed-${baseName}.${ext}`,
            );
            const tempOutputPath = path.resolve(
              outputDir,
              `processed-${baseName}.tmp.${ext}`,
            );
            let outputPath = tempOutputPath;
            let finalFramesDirPath: string | null = null;
            tempFilePath = outputPath;

            if (operation === "FRAME_EXTRACT") {
              finalFramesDirPath = path.resolve(
                outputDir,
                `processed-${baseName}-frames`,
              );
              tempFramesDirPath = `${finalFramesDirPath}.tmp`;
              await fs.rm(tempFramesDirPath, { recursive: true, force: true });
              await fs.mkdir(tempFramesDirPath, { recursive: true });
              outputPath = path.join(tempFramesDirPath, "frame_%06d.png");
              tempFilePath = null;
            }

            const args = buildFfmpegArgs(operation, job.filePath, outputPath);

            console.log(`[>] Processing ${job.fileName} as ${operation}`);
            await runFfmpeg(args);

            if (operation === "FRAME_EXTRACT") {
              await fs.rm(finalFramesDirPath!, { recursive: true, force: true });
              await fs.rename(tempFramesDirPath!, finalFramesDirPath!);
              console.log(
                `[✓] Created frames folder: processed-${baseName}-frames`,
              );
            } else {
              await fs.rm(finalOutputPath, { force: true });
              await fs.rename(outputPath, finalOutputPath);
              console.log(`[✓] Created ${finalOutputPath}`);
            }

            channel.ack(msg);
          } catch (error: any) {
            if (tempFilePath) {
              await fs.rm(tempFilePath, { force: true });
            }
            if (tempFramesDirPath) {
              await fs.rm(tempFramesDirPath, { recursive: true, force: true });
            }
            console.error("[X] Job failed:", error.message);
            channel.reject(msg, false);
          }
        },
        { noAck: false },
      );

      await new Promise((_, reject) => {
        connection.on("close", () => reject(new Error("RabbitMQ closed")));
      });
    } catch (error) {
      console.error("[-] Worker connection error:", error);
      try {
        await channel?.close();
      } catch {}
      try {
        await connection?.close();
      } catch {}
      await sleep(RETRY_DELAY_MS);
    }
  }
}

void startWorker();
