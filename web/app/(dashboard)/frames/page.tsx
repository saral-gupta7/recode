import ToolPage from "@/components/ToolPage";

export default function FrameExtractPage() {
  return (
    <ToolPage
      operation="FRAME_EXTRACT"
      title="Frame Extractor"
      subtitle="Extract frames at one frame per second into a folder."
      doneTitle="All frames are ready"
      downloadLabel="Download Frames"
      idleHint="Outputs sequential PNG frames at fps=1."
    />
  );
}
