import { NextRequest, NextResponse } from "next/server";
import path from "path";
import fs from "fs/promises";
import { spawn } from "child_process";

import { getUploadDir, sanitizeFileName } from "@/lib/storage";

export const runtime = "nodejs";

export async function GET(request: NextRequest) {
  const searchParams = request.nextUrl.searchParams;
  const dir = searchParams.get("dir");

  if (!dir) {
    return new NextResponse("Missing dir parameter", { status: 400 });
  }

  const sanitizedDir = sanitizeFileName(dir);
  if (sanitizedDir !== dir || !sanitizedDir.endsWith("-frames")) {
    return new NextResponse("Invalid dir parameter", { status: 400 });
  }

  const uploadDir = getUploadDir();
  const framesDirPath = path.join(uploadDir, sanitizedDir);

  try {
    const stat = await fs.stat(framesDirPath);
    if (!stat.isDirectory()) {
      return new NextResponse("Frames directory not found", { status: 404 });
    }
  } catch {
    return new NextResponse("Frames directory not found", { status: 404 });
  }

  const archiveName = `${sanitizedDir}.tar.gz`;
  const tar = spawn("tar", ["-czf", "-", "-C", uploadDir, sanitizedDir]);

  const stream = new ReadableStream<Uint8Array>({
    start(controller) {
      tar.stdout.on("data", (chunk: Buffer) => {
        controller.enqueue(new Uint8Array(chunk));
      });

      tar.stdout.on("end", () => {
        controller.close();
      });

      tar.on("error", () => {
        controller.error(new Error("Failed to start archive stream"));
      });

      tar.on("close", (code) => {
        if (code !== 0) {
          controller.error(new Error(`tar exited with code ${code}`));
        }
      });
    },
    cancel() {
      tar.kill("SIGTERM");
    },
  });

  return new NextResponse(stream, {
    headers: {
      "Content-Type": "application/gzip",
      "Content-Disposition": `attachment; filename=\"${archiveName}\"`,
    },
  });
}
