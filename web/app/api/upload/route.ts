import { NextRequest, NextResponse } from "next/server";
import fs from "fs/promises";
import path from "path";

import amqplib from "amqplib";

export async function POST(request: NextRequest) {
  try {
    const formData = await request.formData();

    const title = formData.get("title");
    const description = formData.get("description");

    const file = formData.get("video");

    if (!file || typeof file == "string") {
      return NextResponse.json(
        { error: "Video file uploaded is invalid!" },
        { status: 400 },
      );
    }

    const arrayBuffer = await file.arrayBuffer();
    const buffer = new Uint8Array(arrayBuffer);

    const uploadDir = path.resolve(process.cwd(), "../uploads");
    await fs.mkdir(uploadDir, { recursive: true });
    const safeName = `${Date.now()}-${file.name}`;
    const filePath = path.join(uploadDir, safeName);
    await fs.writeFile(filePath, buffer);

    const connection = await amqplib.connect("amqp://localhost:5672");
    const channel = await connection.createChannel();
    const queue = "video_processing";

    await channel.assertQueue(queue, { durable: true });

    const jobPayload = {
      filePath: filePath,
      fileName: safeName,
      action: "convertToGif", // We can make this dynamic later
    };

    channel.sendToQueue(queue, Buffer.from(JSON.stringify(jobPayload)), {
      persistent: true,
    });

    setTimeout(() => {
      connection.close();
    }, 500);

    return NextResponse.json({
      title: title,
      description: description,
      fileName: safeName,
    });
  } catch (err) {
    console.error(err);
    return NextResponse.json(
      { error: "Failed to save file!" },
      { status: 500 },
    );
  }
}
