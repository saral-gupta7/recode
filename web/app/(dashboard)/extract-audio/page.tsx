import ToolPage from "@/components/ToolPage";

export default function ExtractAudioPage() {
  return (
    <ToolPage
      operation="EXTRACT_AUDIO"
      title="Extract Audio"
      subtitle="Export the audio track from your video as MP3."
      doneTitle="Audio is ready"
      downloadLabel="Download MP3"
      idleHint="Good for voice, podcast, and background music extraction."
    />
  );
}
