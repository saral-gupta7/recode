import ToolPage from "@/components/ToolPage";

export default function TrimPage() {
  return (
    <ToolPage
      operation="TRIM_30S"
      title="Trim to 30 Seconds"
      subtitle="Create a 30-second cut from the start of your video."
      doneTitle="Trimmed video is ready"
      downloadLabel="Download MP4"
      idleHint="Fast way to make short clips for social posts."
    />
  );
}
