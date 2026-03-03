import { NextRequest, NextResponse } from "next/server";
import fs from "fs/promises";
import path from "path";

import amqplib from "amqplib";
import { getUploadDir, sanitizeFileName } from "@/lib/storage";

const RABBITMQ_URL = process.env.RABBITMQ_URL || "amqp://localhost:5672";
const QUEUE_NAME = process.env.QUEUE_NAME || "video_processing";

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

    const sanitizedOriginalName = sanitizeFileName(file.name);
    if (!sanitizedOriginalName) {
      return NextResponse.json(
        { error: "Uploaded filename is invalid!" },
        { status: 400 },
      );
    }

    const arrayBuffer = await file.arrayBuffer();
    const buffer = new Uint8Array(arrayBuffer);

    const uploadDir = getUploadDir();
    await fs.mkdir(uploadDir, { recursive: true });
    const safeName = `${Date.now()}-${sanitizedOriginalName}`;
    const filePath = path.join(uploadDir, safeName);
    await fs.writeFile(filePath, buffer);

    const connection = await amqplib.connect(RABBITMQ_URL);
    const channel = await connection.createChannel();
    await channel.assertQueue(QUEUE_NAME, { durable: true });

    const jobPayload = {
      filePath: filePath,
      fileName: safeName,
      action: "convertToGif", // We can make this dynamic later
    };

    channel.sendToQueue(QUEUE_NAME, Buffer.from(JSON.stringify(jobPayload)), {
      persistent: true,
    });

    await channel.close();
    await connection.close();

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
