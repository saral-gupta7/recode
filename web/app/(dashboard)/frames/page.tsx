import ToolPage from "@/components/ToolPage";

export default function FrameExtractPage() {
  return (
    <ToolPage
      operation="FRAME_EXTRACT"
      title="Frame Extractor"
      subtitle="Extract every frame from your video into a folder."
      doneTitle="All frames are ready"
      downloadLabel="Download Frames"
      idleHint="Outputs sequential PNG frames in a processed folder."
    />
  );
}
