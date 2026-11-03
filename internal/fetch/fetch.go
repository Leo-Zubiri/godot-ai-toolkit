// Package fetch downloads a content tarball over HTTPS, verifies its SHA-256,
// and extracts it — using only the standard library.
package fetch

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// maxArchiveBytes caps a single extracted file to guard against decompression
// bombs. Content tarballs are small (skills + adapters).
const maxArchiveBytes = 64 << 20 // 64 MiB

// DownloadAndExtract fetches a .tar.gz from url, verifies it against expectedSHA
// (hex; empty skips verification), and extracts it into destDir.
func DownloadAndExtract(url, expectedSHA, destDir string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download %s: HTTP %d", url, resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if expectedSHA != "" {
		sum := sha256.Sum256(data)
		got := hex.EncodeToString(sum[:])
		if !strings.EqualFold(got, expectedSHA) {
			return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedSHA, got)
		}
	}
	return extractTarGz(data, destDir)
}

func extractTarGz(data []byte, destDir string) error {
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	cleanDest := filepath.Clean(destDir)
	if err := os.MkdirAll(cleanDest, 0o755); err != nil {
		return err
	}
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := filepath.Join(cleanDest, filepath.FromSlash(hdr.Name))
		if target != cleanDest && !strings.HasPrefix(target, cleanDest+string(os.PathSeparator)) {
			return fmt.Errorf("unsafe path in archive: %q", hdr.Name)
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(hdr.Mode).Perm())
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, io.LimitReader(tr, maxArchiveBytes)); err != nil {
				f.Close()
				return err
			}
			if err := f.Close(); err != nil {
				return err
			}
		default:
			// Skip symlinks and other entry types for safety.
		}
	}
	return nil
}

// ReleaseTarballURL builds the conventional content-tarball URL for a tag.
func ReleaseTarballURL(repo, tag string) string {
	repo = strings.TrimSuffix(repo, "/")
	return fmt.Sprintf("%s/releases/download/%s/godot-ai-toolkit-content-%s.tar.gz", repo, tag, tag)
}
