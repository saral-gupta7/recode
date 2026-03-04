import { NextRequest, NextResponse } from "next/server";
import fs from "fs/promises";
import path from "path";
import amqplib from "amqplib";

import { getUploadDir, sanitizeFileName } from "@/lib/storage";
import { isOperation } from "@/lib/operations";

const RABBITMQ_URL = process.env.RABBITMQ_URL || "amqp://localhost:5672";
const QUEUE_NAME = process.env.QUEUE_NAME || "video_processing";

export async function POST(request: NextRequest) {
  try {
    const formData = await request.formData();
    const file = formData.get("video");
    const rawOperation = formData.get("operation");
    const timestamp = formData.get("timestamp");

    if (!file || typeof file === "string") {
      return NextResponse.json(
        { error: "Video file uploaded is invalid." },
        { status: 400 },
      );
    }

    if (!rawOperation || typeof rawOperation !== "string" || !isOperation(rawOperation)) {
      return NextResponse.json(
        { error: "Invalid or missing operation." },
        { status: 400 },
      );
    }

    const sanitizedOriginalName = sanitizeFileName(file.name);
    if (!sanitizedOriginalName) {
      return NextResponse.json(
        { error: "Uploaded filename is invalid." },
        { status: 400 },
      );
    }

    const uploadDir = getUploadDir();
    await fs.mkdir(uploadDir, { recursive: true });

    // Timestamp prefix keeps jobs unique even for same original filename.
    const safeName = `${Date.now()}-${sanitizedOriginalName}`;
    const filePath = path.join(uploadDir, safeName);

    const arrayBuffer = await file.arrayBuffer();
    await fs.writeFile(filePath, new Uint8Array(arrayBuffer));

    const connection = await amqplib.connect(RABBITMQ_URL);
    const channel = await connection.createChannel();
    await channel.assertQueue(QUEUE_NAME, { durable: true });

    const jobPayload = {
      filePath,
      fileName: safeName,
      operation: rawOperation,
      timestamp: typeof timestamp === "string" ? timestamp : undefined,
    };

    channel.sendToQueue(QUEUE_NAME, Buffer.from(JSON.stringify(jobPayload)), {
      persistent: true,
    });

    await channel.close();
    await connection.close();

    return NextResponse.json({ fileName: safeName, operation: rawOperation });
  } catch (error) {
    console.error(error);
    return NextResponse.json({ error: "Failed to save file." }, { status: 500 });
  }
}
