import ToolPage from "@/components/ToolPage";

export default function GrayscalePage() {
  return (
    <ToolPage
      operation="BLACK_AND_WHITE"
      title="Grayscale Video"
      subtitle="Turn your video into a black-and-white MP4."
      doneTitle="Grayscale video is ready"
      downloadLabel="Download MP4"
      idleHint="Output is H.264 MP4 for broad player compatibility."
    />
  );
}
