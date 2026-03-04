import ToolPage from "@/components/ToolPage";

export default function ExtractAudioPage() {
  return (
    <ToolPage
      operation="EXTRACT_AUDIO"
      title="Extract Audio"
      subtitle="Extract the soundtrack and export it as MP3."
      doneTitle="Audio track is ready"
      downloadLabel="Download MP3"
      idleHint="Works best for videos with a clean primary audio track."
    />
  );
}
