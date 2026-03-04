import amqplib from "amqplib";
import path from "path";
import ffmpeg from "fluent-ffmpeg";

const QUEUE_NAME = process.env.QUEUE_NAME || "video_processing";
const RABBITMQ_URL = process.env.RABBITMQ_URL || "amqp://rabbitmq:5672";
const RETRY_DELAY_MS = 5000;

async function sleep(ms: number) {
  await new Promise((resolve) => setTimeout(resolve, ms));
}

type JobOperation =
  | "EXTRACT_AUDIO"
  | "BLACK_AND_WHITE"
  | "FRAME_EXTRACT"
  | "VIDEO_TO_GIF";

function normalizeOperation(
  rawOperation?: string,
  rawAction?: string,
): JobOperation {
  if (rawOperation) return rawOperation as JobOperation;
  if (rawAction === "convertToGif") return "VIDEO_TO_GIF";
  return "VIDEO_TO_GIF";
}

function outputExtensionForOperation(operation: JobOperation) {
  switch (operation) {
    case "EXTRACT_AUDIO":
      return "mp3";
    case "FRAME_EXTRACT":
      return "png";
    case "BLACK_AND_WHITE":
      return "mp4";
    case "VIDEO_TO_GIF":
    default:
      return "gif";
  }
}

async function startWorker() {
  while (true) {
    let connection: any = null;
    let channel: any = null;

    try {
      connection = await amqplib.connect(RABBITMQ_URL);
      channel = await connection.createChannel();

      await channel.assertQueue(QUEUE_NAME, { durable: true });
      channel.prefetch(1);

      console.log(
        `[*] Worker booted. Listening for messages in '${QUEUE_NAME}'...`,
      );

      // Handle connection drops
      connection.on("close", () => {
        console.error("[-] Connection closed. Restarting...");
        return; // The outer while loop handles the restart
      });

      channel.consume(
        QUEUE_NAME,
        async (msg: any) => {
          if (!msg) return;

          try {
            const { filePath, fileName, operation, action, timestamp } =
              JSON.parse(msg.content.toString());

            const resolvedOperation = normalizeOperation(operation, action);
            const outputExtension =
              outputExtensionForOperation(resolvedOperation);

            // Clean up the file name so we don't end up with .mp4.mp4
            const parsedFileName = path.parse(fileName).name;
            const outputPath = path.resolve(
              filePath,
              "..",
              `processed-${parsedFileName}`, // Fixed variable name
            );

            let command = ffmpeg(filePath);

            switch (resolvedOperation) {
              case "EXTRACT_AUDIO":
                command = command
                  .noVideo()
                  .toFormat("mp3")
                  .audioBitrate("192k");
                break;
              case "BLACK_AND_WHITE":
                command = command
                  .videoCodec("libx264")
                  .audioCodec("aac") // 👈 CRITICAL: Forces Apple-friendly audio
                  .videoFilters("hue=s=0")
                  .outputOptions([
                    "-pix_fmt",
                    "yuv420p", // Forces Apple-friendly pixels
                    "-movflags",
                    "+faststart", // Required for web playback
                  ])
                  .toFormat("mp4");
                break;
              case "FRAME_EXTRACT":
                command = command.seekInput(timestamp || "00:00:01").frames(1);
                break;
              case "VIDEO_TO_GIF":
              default:
                command = command
                  .setDuration(5)
                  .fps(10)
                  .size("480x?")
                  .toFormat("gif");
                break;
            }

            command
              .on("end", () => {
                console.log(`[✓] Success! Job ${resolvedOperation} finished.`);
                channel!.ack(msg);
              })
              .on("error", (err) => {
                console.error("[X] FFmpeg failed:", err.message);
                channel!.reject(msg, false);
              })
              .save(outputPath);
          } catch (error) {
            console.error("[X] Parsing failed:", error);
            channel.reject(msg, false);
          }
        },
        { noAck: false },
      );

      return;
    } catch (error) {
      console.error("[-] Worker Error:", error);
      await sleep(RETRY_DELAY_MS);
    }
  }
}

startWorker();
