import amqplib from "amqplib";
import path from "path";
import ffmpeg from "fluent-ffmpeg";

const QUEUE_NAME = process.env.QUEUE_NAME || "video_processing";
const RABBITMQ_URL = process.env.RABBITMQ_URL || "amqp://rabbitmq:5672";
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
            const { filePath, fileName, operation, timestamp } = JSON.parse(
              msg.content.toString(),
            );

            let extension = "mp4";
            if (operation === "EXTRACT_AUDIO") extension = "mp3";
            if (operation === "VIDEO_TO_GIF") extension = "gif";
            if (operation === "FRAME_EXTRACT") extension = "png";

            const outputPath = path.resolve(
              filePath,
              "..",
              `processed-${fileName}.${extension}`,
            );
            console.log(`\n[↓] Job: ${operation} | File: ${fileName}`);

            let command = ffmpeg(filePath);

            switch (operation) {
              case "EXTRACT_AUDIO":
                command = command
                  .noVideo()
                  .toFormat("mp3")
                  .audioBitrate("192k");
                break;
              case "BLACK_AND_WHITE":
                command = command.videoFilters("format=gray").toFormat("mp4");
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
                console.log(`[✓] Success!`);
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
