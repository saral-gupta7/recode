import { NextRequest, NextResponse } from "next/server";
import fs from "fs/promises";
import path from "path";

import { getUploadDir, sanitizeFileName } from "@/lib/storage";
import { isOperation, OPERATION_OUTPUT } from "@/lib/operations";

export async function GET(request: NextRequest) {
  const searchParams = request.nextUrl.searchParams;
  const fileName = searchParams.get("fileName");
  const operation = searchParams.get("operation");

  if (!fileName) {
    return NextResponse.json({ error: "Missing fileName" }, { status: 400 });
  }

  if (!operation || !isOperation(operation)) {
    return NextResponse.json({ error: "Missing or invalid operation" }, { status: 400 });
  }

  const sanitizedFileName = sanitizeFileName(fileName);
  if (sanitizedFileName !== fileName) {
    return NextResponse.json({ error: "Invalid fileName" }, { status: 400 });
  }

  const baseName = path.parse(sanitizedFileName).name;
  if (operation === "FRAME_EXTRACT") {
    const framesDirName = `processed-${baseName}-frames`;
    const framesDirPath = path.join(getUploadDir(), framesDirName);

    try {
      const files = await fs.readdir(framesDirPath);
      const frames = files.filter((file) => file.toLowerCase().endsWith(".png"));

      if (frames.length > 0) {
        return NextResponse.json({
          status: "done",
          framesDirName,
          frameCount: frames.length,
          mimeType: "application/gzip",
          downloadUrl: `/api/download-frames?dir=${encodeURIComponent(framesDirName)}`,
        });
      }
    } catch {}

    return NextResponse.json({ status: "processing" });
  }

  const output = OPERATION_OUTPUT[operation];
  const processedFileName = `processed-${baseName}.${output.extension}`;
  const filePath = path.join(getUploadDir(), processedFileName);

  try {
    await fs.access(filePath);
    return NextResponse.json({
      status: "done",
      fileName: processedFileName,
      mimeType: output.mimeType,
      downloadUrl: `/api/download?file=${encodeURIComponent(processedFileName)}`,
    });
  } catch {
    return NextResponse.json({ status: "processing" });
  }
}
