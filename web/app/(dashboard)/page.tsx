import ToolPage from "@/components/ToolPage";

export default function VideoToGifPage() {
  return (
    <ToolPage
      operation="VIDEO_TO_GIF"
      title="Video to GIF"
      subtitle="Convert a short clip into a lightweight GIF."
      doneTitle="GIF is ready"
      downloadLabel="Download GIF"
      idleHint="Best results: short clips with clear motion."
    />
  );
}
