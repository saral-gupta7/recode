import amqplib from "amqplib";
import path from "path";

import ffmpeg from "fluent-ffmpeg";

const QUEUE_NAME = process.env.QUEUE_NAME || "video_processing";
const RABBITMQ_URL = process.env.RABBITMQ_URL || "amqp://localhost:5672";
const RETRY_DELAY_MS = 5000;

async function sleep(ms: number) {
  await new Promise((resolve) => setTimeout(resolve, ms));
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

      connection.on("close", async () => {
        console.error("[-] RabbitMQ connection closed. Reconnecting...");
        await sleep(RETRY_DELAY_MS);
        void startWorker();
      });

      connection.on("error", (error: unknown) => {
        console.error("[-] RabbitMQ connection error:", error);
      });

      channel.consume(
        QUEUE_NAME,
        async (msg: any) => {
          if (!msg) return;

          try {
            const jobData = JSON.parse(msg.content.toString());
            if (!jobData?.filePath || !jobData?.fileName) {
              console.error("[X] Invalid job payload:", jobData);
              channel.reject(msg, false);
              return;
            }

            console.log(`\n[↓] Pulled job from queue: ${jobData.fileName}`);

            const inputPath = jobData.filePath;
            const outputPath = path.resolve(
              inputPath,
              "..",
              `processed-${jobData.fileName}.gif`,
            );

            console.log("[⚙️] Starting FFmpeg processing...");

            ffmpeg(inputPath)
              .output(outputPath)
              .setDuration(3)
              .on("end", () => {
                console.log(`[✓] Success! Saved to: ${outputPath}`);
                channel!.ack(msg);
              })
              .on("error", (err) => {
                console.error("[X] FFmpeg failed:", err.message);
                channel!.reject(msg, false);
              })
              .run();
          } catch (error) {
            console.error("[X] Failed to process queue message:", error);
            channel.reject(msg, false);
          }
        },
        {
          noAck: false,
        },
      );

      return;
    } catch (error) {
      console.error("[-] Failed to start worker:", error);
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

startWorker();
