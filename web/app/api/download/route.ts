import { NextRequest, NextResponse } from "next/server";
import fs from "fs/promises";
import path from "path";
import { getUploadDir, sanitizeFileName } from "@/lib/storage";

function contentTypeForFile(name: string) {
  const ext = path.extname(name).toLowerCase();
  if (ext === ".gif") return "image/gif";
  if (ext === ".mp4") return "video/mp4";
  if (ext === ".mp3") return "audio/mpeg";
  if (ext === ".png") return "image/png";
  return "application/octet-stream";
}

export async function GET(request: NextRequest) {
  const searchParams = request.nextUrl.searchParams;
  const fileName = searchParams.get("file");

  if (!fileName)
    return new NextResponse("Missing file parameter", { status: 400 });

  const sanitizedFileName = sanitizeFileName(fileName);
  if (sanitizedFileName !== fileName) {
    return new NextResponse("Invalid file parameter", { status: 400 });
  }

  const filePath = path.join(getUploadDir(), sanitizedFileName);

  try {
    const fileBuffer = await fs.readFile(filePath);

    return new NextResponse(fileBuffer, {
      headers: {
        "Content-Disposition": `attachment; filename="${fileName}"`,
        "Content-Type": contentTypeForFile(fileName),
      },
    });
  } catch {
    return new NextResponse("File not found", { status: 404 });
  }
}
