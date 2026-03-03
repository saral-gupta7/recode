"use client";

import { useState, useEffect, useCallback } from "react";
import { useDropzone } from "react-dropzone";
import { motion, AnimatePresence } from "framer-motion";
import { UploadCloud, CheckCircle, Download } from "lucide-react";

export default function Home() {
  const [status, setStatus] = useState<
    "idle" | "uploading" | "processing" | "done"
  >("idle");
  const [jobId, setJobId] = useState<string | null>(null);
  const [downloadUrl, setDownloadUrl] = useState<string | null>(null);

  const onDrop = useCallback(async (acceptedFiles: File[]) => {
    const file = acceptedFiles[0];
    if (!file) return;

    setStatus("uploading");
    const formData = new FormData();
    formData.append("video", file);
    formData.append("title", "Frontend Upload");

    try {
      const response = await fetch("/api/upload", {
        method: "POST",
        body: formData,
      });

      if (!response.ok) throw new Error("Upload failed");

      const data = await response.json();
      console.log("Success:", data);

      setJobId(data.fileName);
      setStatus("processing");
    } catch (error) {
      console.error("Upload failed", error);
      setStatus("idle");
    }
  }, []);

  useEffect(() => {
    // Only poll if we have a job ID and the status is currently processing
    if (!jobId || status !== "processing") return;

    const interval = setInterval(async () => {
      try {
        const res = await fetch(`/api/status?fileName=${jobId}`);
        const data = await res.json();

        if (data.status === "done") {
          setStatus("done");
          setDownloadUrl(data.downloadUrl);
          clearInterval(interval); // Stop polling!
        }
      } catch (error) {
        console.error("Polling error:", error);
      }
    }, 2000); // Check every 2 seconds

    return () => clearInterval(interval); // Cleanup on unmount
  }, [jobId, status]);

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: { "video/*": [] },
    multiple: false,
  });

  return (
    <main className="flex min-h-screen items-center justify-center bg-zinc-950 p-6 text-zinc-100">
      <div className="w-full max-w-xl">
        <div className="mb-8 text-center">
          <h1 className="text-4xl font-bold tracking-tight mb-2 text-white">
            Recode
          </h1>
          <p className="text-zinc-400">
            Lightning fast, login-free media transformations.
          </p>
        </div>

        <motion.div
          {...getRootProps()}
          className={`relative cursor-pointer overflow-hidden rounded-2xl border-2 border-dashed p-12 text-center transition-colors ${
            isDragActive
              ? "border-blue-500 bg-blue-500/10"
              : "border-zinc-800 bg-zinc-900/50 hover:border-zinc-700"
          }`}
          animate={{ scale: isDragActive ? 1.02 : 1 }}
          transition={{ type: "spring", stiffness: 300, damping: 20 }}
        >
          <input {...getInputProps()} />

          <AnimatePresence mode="wait">
            {status === "uploading" && (
              <motion.div
                key="uploading"
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -10 }}
                className="flex flex-col items-center gap-4"
              >
                <div className="h-12 w-12 animate-spin rounded-full border-4 border-zinc-800 border-t-blue-500" />
                <p className="text-lg font-medium">Sending to API...</p>
              </motion.div>
            )}

            {status === "processing" && (
              <motion.div
                key="processing"
                initial={{ opacity: 0, scale: 0.9 }}
                animate={{ opacity: 1, scale: 1 }}
                exit={{ opacity: 0, scale: 0.9 }}
                className="flex flex-col items-center gap-4 text-blue-400"
              >
                <div className="h-12 w-12 animate-spin rounded-full border-4 border-zinc-800 border-t-blue-400" />
                <p className="text-lg font-medium text-white">
                  Worker is crunching the pixels...
                </p>
                <p className="text-sm text-zinc-500 mt-2">
                  Converting via FFmpeg
                </p>
              </motion.div>
            )}

            {status === "done" && downloadUrl && (
              <motion.div
                key="done"
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                className="flex flex-col items-center gap-4 text-emerald-400"
              >
                <CheckCircle size={48} />
                <p className="text-xl font-bold text-white">GIF is Ready!</p>

                {/* e.stopPropagation() prevents the dropzone from opening the file picker again */}
                <a
                  href={downloadUrl}
                  download
                  onClick={(e) => e.stopPropagation()}
                  className="mt-4 flex items-center gap-2 rounded-xl bg-emerald-500 px-6 py-3 text-sm font-bold text-zinc-950 shadow-[0_0_20px_rgba(16,185,129,0.4)] transition-all hover:scale-105 hover:bg-emerald-400 active:scale-95"
                >
                  <Download size={20} />
                  Download Your GIF
                </a>

                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    setStatus("idle");
                    setJobId(null);
                    setDownloadUrl(null);
                  }}
                  className="mt-4 text-sm text-zinc-500 underline transition-colors hover:text-zinc-300"
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
                  animate={{ y: [0, -10, 0] }}
                  transition={{
                    repeat: Infinity,
                    duration: 4,
                    ease: "easeInOut",
                  }}
                >
                  <UploadCloud
                    size={64}
                    className={isDragActive ? "text-blue-500" : "text-zinc-500"}
                  />
                </motion.div>
                <div>
                  <p className="text-xl font-medium">
                    {isDragActive
                      ? "Drop the video here!"
                      : "Drag & drop your video"}
                  </p>
                  <p className="text-sm text-zinc-500 mt-2">
                    or click to browse your files
                  </p>
                </div>
              </motion.div>
            )}
          </AnimatePresence>
        </motion.div>
      </div>
    </main>
  );
}
