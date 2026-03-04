import { NextRequest, NextResponse } from "next/server";
import fs from "fs/promises";
import path from "path";
import { getUploadDir, sanitizeFileName } from "@/lib/storage";

export async function GET(request: NextRequest) {
  const searchParams = request.nextUrl.searchParams;
  const fileName = searchParams.get("fileName");

  const extensions = [".mp3", ".gif", ".png", ".mp4"];

  if (!fileName) {
    return NextResponse.json({ error: "Missing fileName" }, { status: 400 });
  }

  const sanitizedFileName = sanitizeFileName(fileName);
  if (sanitizedFileName !== fileName) {
    return NextResponse.json({ error: "Invalid fileName" }, { status: 400 });
  }
  const baseName = path.parse(sanitizedFileName).name;

  for (const ext of extensions) {
    const processedFileName = `processed-${baseName}${ext}`;
    const filePath = path.join(getUploadDir(), processedFileName);

    try {
      await fs.access(filePath);

      return NextResponse.json({
        status: "done",
        downloadUrl: `/api/download?file=${processedFileName}`,
      });
    } catch {
      continue;
    }
  }

  return NextResponse.json({ status: "processing" });
}
