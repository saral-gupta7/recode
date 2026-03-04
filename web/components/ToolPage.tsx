"use client";

import { useCallback, useEffect, useState } from "react";
import { useDropzone } from "react-dropzone";
import { AnimatePresence, motion } from "framer-motion";
import { CheckCircle2, Download, UploadCloud } from "lucide-react";
import Image from "next/image";

import type { Operation } from "@/lib/operations";

type ToolPageProps = {
  operation: Operation;
  title: string;
  subtitle: string;
  doneTitle: string;
  downloadLabel: string;
  idleHint: string;
};

type StatusResponse = {
  status: "processing" | "done";
  downloadUrl?: string;
  mimeType?: string;
  framesDirName?: string;
  frameCount?: number;
};

function Preview({ url, mimeType }: { url: string; mimeType: string }) {
  if (mimeType.startsWith("image/")) {
    return (
      <Image
        src={url}
        alt="Processed output preview"
        width={1200}
        height={675}
        unoptimized
        className="max-h-64 w-full rounded-xl border border-zinc-700 object-contain"
      />
    );
  }

  if (mimeType.startsWith("video/")) {
    return (
      <video
        src={url}
        controls
        className="max-h-64 w-full rounded-xl border border-zinc-700"
      />
    );
  }

  if (mimeType.startsWith("audio/")) {
    return <audio src={url} controls className="w-full" />;
  }

  return null;
}

export default function ToolPage({
  operation,
  title,
  subtitle,
  doneTitle,
  downloadLabel,
  idleHint,
}: ToolPageProps) {
  const [status, setStatus] = useState<"idle" | "uploading" | "processing" | "done">("idle");
  const [jobId, setJobId] = useState<string | null>(null);
  const [downloadUrl, setDownloadUrl] = useState<string | null>(null);
  const [mimeType, setMimeType] = useState<string>("application/octet-stream");
  const [framesDirName, setFramesDirName] = useState<string | null>(null);
  const [frameCount, setFrameCount] = useState<number | null>(null);

  const onDrop = useCallback(
    async (acceptedFiles: File[]) => {
      const file = acceptedFiles[0];
      if (!file) return;

      setStatus("uploading");
      const formData = new FormData();
      formData.append("video", file);
      formData.append("operation", operation);

      try {
        const response = await fetch("/api/upload", {
          method: "POST",
          body: formData,
        });

        if (!response.ok) throw new Error("Upload failed");

        const data = await response.json();
        setJobId(data.fileName);
        setStatus("processing");
      } catch (error) {
        console.error("Upload failed", error);
        setStatus("idle");
      }
    },
    [operation],
  );

  useEffect(() => {
    if (!jobId || status !== "processing") return;

    const interval = setInterval(async () => {
      try {
        const res = await fetch(
          `/api/status?fileName=${encodeURIComponent(jobId)}&operation=${operation}`,
        );
        const data = (await res.json()) as StatusResponse;

        if (data.status === "done") {
          setStatus("done");
          setDownloadUrl(data.downloadUrl || null);
          if (data.mimeType) setMimeType(data.mimeType);
          setFramesDirName(data.framesDirName || null);
          setFrameCount(typeof data.frameCount === "number" ? data.frameCount : null);
          clearInterval(interval);
        }
      } catch (error) {
        console.error("Polling error:", error);
      }
    }, 2000);

    return () => clearInterval(interval);
  }, [jobId, operation, status]);

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: { "video/*": [] },
    multiple: false,
  });

  return (
    <section className="mx-auto w-full max-w-3xl">
      <div className="mb-6 text-center sm:mb-8">
        <h1 className="mb-2 text-3xl font-bold tracking-tight text-white sm:text-4xl">
          {title}
        </h1>
        <p className="text-sm text-zinc-400 sm:text-base">{subtitle}</p>
      </div>

      <div {...getRootProps()} className="relative">
        <input {...getInputProps()} />

        <motion.div
          className={`cursor-pointer overflow-hidden rounded-2xl border-2 border-dashed p-6 text-center transition-colors sm:p-10 ${
            isDragActive
              ? "border-blue-500 bg-blue-500/10"
              : "border-zinc-800 bg-zinc-900/50 hover:border-zinc-700"
          }`}
          animate={{ scale: isDragActive ? 1.02 : 1 }}
          transition={{ type: "spring", stiffness: 300, damping: 20 }}
        >
          <AnimatePresence mode="wait">
            {status === "uploading" && (
              <motion.div
                key="uploading"
                initial={{ opacity: 0, y: 8 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -8 }}
                className="flex flex-col items-center gap-4"
              >
                <div className="h-12 w-12 animate-spin rounded-full border-4 border-zinc-800 border-t-blue-500" />
                <p className="text-base font-medium sm:text-lg">Uploading file...</p>
              </motion.div>
            )}

            {status === "processing" && (
              <motion.div
                key="processing"
                initial={{ opacity: 0, scale: 0.92 }}
                animate={{ opacity: 1, scale: 1 }}
                exit={{ opacity: 0, scale: 0.92 }}
                className="flex flex-col items-center gap-4"
              >
                <div className="h-12 w-12 animate-spin rounded-full border-4 border-zinc-800 border-t-cyan-400" />
                <p className="text-base font-medium text-white sm:text-lg">
                  Processing with ffmpeg...
                </p>
                <p className="text-sm text-zinc-500">Please keep this tab open.</p>
              </motion.div>
            )}

            {status === "done" && (
              <motion.div
                key="done"
                initial={{ opacity: 0, y: 8 }}
                animate={{ opacity: 1, y: 0 }}
                className="mx-auto flex w-full max-w-xl flex-col items-center gap-4"
              >
                <CheckCircle2 size={42} className="text-emerald-400" />
                <p className="text-lg font-semibold text-white sm:text-xl">{doneTitle}</p>

                {downloadUrl && <Preview url={downloadUrl} mimeType={mimeType} />}

                {!downloadUrl && framesDirName && (
                  <div className="w-full rounded-xl border border-zinc-700 bg-zinc-900/70 p-4 text-left text-sm text-zinc-300">
                    <p className="font-medium text-white">Frames extracted successfully.</p>
                    <p className="mt-1 text-zinc-400">Folder: {framesDirName}</p>
                    <p className="mt-1 text-zinc-400">Total frames: {frameCount ?? 0}</p>
                  </div>
                )}

                {downloadUrl && (
                  <a
                    href={downloadUrl}
                    download
                    onClick={(e) => e.stopPropagation()}
                    className="mt-2 inline-flex items-center gap-2 rounded-xl bg-emerald-500 px-5 py-3 text-sm font-bold text-zinc-950 transition hover:bg-emerald-400"
                  >
                    <Download size={18} />
                    {downloadLabel}
                  </a>
                )}

                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    setStatus("idle");
                    setJobId(null);
                    setDownloadUrl(null);
                    setFramesDirName(null);
                    setFrameCount(null);
                  }}
                  className="text-sm text-zinc-500 underline hover:text-zinc-300"
                >
                  Process another video
                </button>
              </motion.div>
            )}

            {status === "idle" && (
              <motion.div
                key="idle"
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                exit={{ opacity: 0 }}
                className="flex flex-col items-center gap-4"
              >
                <motion.div
                  animate={{ y: [0, -8, 0] }}
                  transition={{ repeat: Infinity, duration: 3.5, ease: "easeInOut" }}
                >
                  <UploadCloud
                    size={58}
                    className={isDragActive ? "text-blue-500" : "text-zinc-500"}
                  />
                </motion.div>
                <div>
                  <p className="text-lg font-medium sm:text-xl">
                    {isDragActive ? "Drop the video here" : "Drag and drop a video"}
                  </p>
                  <p className="mt-2 text-sm text-zinc-500">{idleHint}</p>
                </div>
              </motion.div>
            )}
          </AnimatePresence>
        </motion.div>
      </div>
    </section>
  );
}
