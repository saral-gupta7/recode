import ToolPage from "@/components/ToolPage";

export default function CompressPage() {
  return (
    <ToolPage
      operation="COMPRESS_VIDEO"
      title="Compress Video"
      subtitle="Reduce file size while keeping a balanced quality."
      doneTitle="Compressed video is ready"
      downloadLabel="Download MP4"
      idleHint="Targets 1280p max with efficient H.264 settings."
    />
  );
}
