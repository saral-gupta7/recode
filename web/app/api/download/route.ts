import { NextRequest, NextResponse } from "next/server";
import fs from "fs/promises";
import path from "path";

export async function GET(request: NextRequest) {
  const searchParams = request.nextUrl.searchParams;
  const fileName = searchParams.get("file");

  if (!fileName)
    return new NextResponse("Missing file parameter", { status: 400 });

  const filePath = path.resolve(process.cwd(), "../uploads", fileName);

  try {
    const fileBuffer = await fs.readFile(filePath);

    // Serve the file as an attachment so the browser downloads it
    return new NextResponse(fileBuffer, {
      headers: {
        "Content-Disposition": `attachment; filename="${fileName}"`,
        "Content-Type": "image/gif",
      },
    });
  } catch (error) {
    return new NextResponse("File not found", { status: 404 });
  }
}
