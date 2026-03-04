import ToolPage from "@/components/ToolPage";

export default function RemoveAudioPage() {
  return (
    <ToolPage
      operation="REMOVE_AUDIO"
      title="Remove Audio"
      subtitle="Create a silent MP4 from your input video."
      doneTitle="Silent video is ready"
      downloadLabel="Download MP4"
      idleHint="Keeps video only and strips the audio track."
    />
  );
}
