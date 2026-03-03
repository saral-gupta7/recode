import amqplib from "amqplib";
import path from "path";

import ffmpeg from "fluent-ffmpeg";

const QUEUE_NAME = "video_processing";

async function startWorker() {
  try {
    const connection = await amqplib.connect("amqp://localhost:5672");
    const channel = await connection.createChannel();

    await channel.assertQueue(QUEUE_NAME, { durable: true });

    channel.prefetch(1);

    console.log(
      `[*] Worker booted. Listening for messages in '${QUEUE_NAME}'...`,
    );

    channel.consume(
      QUEUE_NAME,
      async (msg) => {
        if (msg !== null) {
          const jobData = JSON.parse(msg.content.toString());
          console.log(`\n[↓] Pulled job from queue: ${jobData.fileName}`);

          const inputPath = jobData.filePath;
          const outputPath = path.resolve(
            inputPath,
            "..",
            `processed-${jobData.fileName}.gif`,
          );

          console.log(`[⚙️] Starting FFmpeg processing...`);

          ffmpeg(inputPath)
            .output(outputPath)
            .setDuration(3)
            .on("end", () => {
              console.log(`[✓] Success! Saved to: ${outputPath}`);
              channel.ack(msg);
            })
            .on("error", (err) => {
              console.error(`[X] FFmpeg failed:`, err.message);
              channel.reject(msg, false);
            })
            .run();
        }
      },
      {
        noAck: false,
      },
    );
  } catch (error) {
    console.error("[-] Failed to start worker:", error);
  }
}

startWorker();
