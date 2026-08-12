package job

import "testing"

func TestImageOperationsReportImageKind(t *testing.T) {
	for _, operation := range []Operation{OperationImageCrop, OperationImageAdjust, OperationImagePadding} {
		kind, err := operation.MediaKind()
		if err != nil || kind != MediaKindImage {
			t.Fatalf("MediaKind(%q) = %q, %v", operation, kind, err)
		}
	}
}

func TestUnknownOperationHasNoMediaKind(t *testing.T) {
	if _, err := Operation("video_clip").MediaKind(); err == nil {
		t.Fatal("expected error")
	}
}
