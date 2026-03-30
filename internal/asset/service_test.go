package asset

import "testing"

func TestIsSafeSegment(t *testing.T) {
	safe := []string{"images", "article-123", "cover", "my_file"}
	for _, s := range safe {
		if !isSafeSegment(s) {
			t.Errorf("expected safe: %q", s)
		}
	}

	unsafe := []string{"", ".", "..", "../etc", "a/b", "a\\b", "..\\windows"}
	for _, s := range unsafe {
		if isSafeSegment(s) {
			t.Errorf("expected unsafe: %q", s)
		}
	}
}

func TestProcessUpload_PathTraversal(t *testing.T) {
	svc := &Service{} // repo 和 storage 不需要，校验在它们之前

	cases := []struct {
		scope, name string
	}{
		{"../../etc", "passwd"},
		{"images", "../../.env"},
		{"..", "foo"},
		{"images/../..", "bar"},
	}

	for _, c := range cases {
		_, _, err := svc.ProcessUpload("test.jpg", c.scope, c.name)
		if err != ErrUnsafePath {
			t.Errorf("scope=%q name=%q: expected ErrUnsafePath, got %v", c.scope, c.name, err)
		}
	}
}

func TestProcessUpload_Normal(t *testing.T) {
	svc := &Service{}
	record, savedPath, err := svc.ProcessUpload("photo.jpg", "images", "cover")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if record.Path != "/images/cover.jpg" {
		t.Errorf("path = %q, want /images/cover.jpg", record.Path)
	}
	if savedPath != "images/cover.jpg" {
		t.Errorf("savedPath = %q, want images/cover.jpg", savedPath)
	}
}
