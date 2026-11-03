package fetch

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func makeTarGz(t *testing.T, files map[string]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for name, content := range files {
		hdr := &tar.Header{Name: name, Mode: 0o644, Size: int64(len(content)), Typeflag: tar.TypeReg}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestDownloadAndExtract(t *testing.T) {
	data := makeTarGz(t, map[string]string{
		"skills/x/SKILL.md":    "hi",
		"adapters/claude.json": "{}",
	})
	sum := sha256.Sum256(data)
	shaHex := hex.EncodeToString(sum[:])

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(data)
	}))
	defer srv.Close()

	dest := t.TempDir()
	if err := DownloadAndExtract(srv.URL, shaHex, dest); err != nil {
		t.Fatalf("extract: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(dest, "skills", "x", "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hi" {
		t.Errorf("content = %q", got)
	}

	if err := DownloadAndExtract(srv.URL, "deadbeef", t.TempDir()); err == nil {
		t.Error("expected a checksum mismatch error")
	}
}

func TestReleaseTarballURL(t *testing.T) {
	got := ReleaseTarballURL("https://github.com/u/repo/", "v1.2.3")
	want := "https://github.com/u/repo/releases/download/v1.2.3/godot-ai-toolkit-content-v1.2.3.tar.gz"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
