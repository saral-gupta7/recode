import { NextRequest, NextResponse } from "next/server";
import fs from "fs/promises";
import path from "path";

export async function GET(request: NextRequest) {
  const searchParams = request.nextUrl.searchParams;
  const fileName = searchParams.get("fileName");

  if (!fileName) {
    return NextResponse.json({ error: "Missing fileName" }, { status: 400 });
  }

  // The expected name of the file output by our worker
  const processedFileName = `processed-${fileName}.gif`;
  const filePath = path.resolve(process.cwd(), "../uploads", processedFileName);

  try {
    // If fs.access succeeds, the file exists!
    await fs.access(filePath);
    return NextResponse.json({
      status: "done",
      downloadUrl: `/api/download?file=${processedFileName}`,
    });
  } catch (error) {
    // If fs.access throws an error, the file isn't there yet
    return NextResponse.json({ status: "processing" });
  }
}
