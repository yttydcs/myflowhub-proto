package protocolmap

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	beginMarker = "<!-- BEGIN GENERATED -->"
	endMarker   = "<!-- END GENERATED -->"
)

type UpdateResult struct {
	Changed bool
	Doc     []byte
}

func BuildGeneratedSection(protocolRoot string) (string, error) {
	pm, err := ParseProtocolRoot(protocolRoot)
	if err != nil {
		return "", err
	}
	return RenderGeneratedMarkdown(pm)
}

// UpdateDocFile reads docPath and updates only the generated block.
func UpdateDocFile(protocolRoot string, docPath string) (*UpdateResult, error) {
	docPath = filepath.Clean(docPath)
	src, err := os.ReadFile(docPath)
	if err != nil {
		return nil, fmt.Errorf("read doc: %w", err)
	}

	gen, err := BuildGeneratedSection(protocolRoot)
	if err != nil {
		return nil, err
	}

	newDoc, err := ReplaceGeneratedBlock(string(src), gen)
	if err != nil {
		return nil, err
	}

	changed := normalizeNewlines(string(src)) != normalizeNewlines(newDoc)
	return &UpdateResult{Changed: changed, Doc: []byte(newDoc)}, nil
}

func ReplaceGeneratedBlock(doc string, generated string) (string, error) {
	// Work in LF for reliable indexing, then keep original newline style at write-time by preserving existing markers lines.
	docLF := normalizeNewlines(doc)
	lines := strings.Split(docLF, "\n")

	begin, end := -1, -1
	for i := range lines {
		if strings.TrimSpace(lines[i]) == beginMarker {
			begin = i
			continue
		}
		if strings.TrimSpace(lines[i]) == endMarker && begin >= 0 {
			end = i
			break
		}
	}
	if begin < 0 || end < 0 || end <= begin {
		return "", fmt.Errorf("generated markers not found or invalid: %q ... %q", beginMarker, endMarker)
	}

	genLF := strings.TrimSuffix(normalizeNewlines(generated), "\n")
	genLines := []string{}
	if genLF != "" {
		genLines = strings.Split(genLF, "\n")
	}

	out := make([]string, 0, len(lines)-((end-begin)-1)+len(genLines))
	out = append(out, lines[:begin+1]...)
	out = append(out, genLines...)
	out = append(out, lines[end:]...)

	joined := strings.Join(out, "\n")
	// Preserve original trailing newline if present.
	if strings.HasSuffix(docLF, "\n") && !strings.HasSuffix(joined, "\n") {
		joined += "\n"
	}

	// Preserve CRLF if original doc uses it.
	if hasCRLF(doc) {
		joined = strings.ReplaceAll(joined, "\n", "\r\n")
	}
	return joined, nil
}

func WriteFileIfChanged(path string, content []byte) (bool, error) {
	old, err := os.ReadFile(path)
	if err == nil && normalizeNewlines(string(old)) == normalizeNewlines(string(content)) {
		return false, nil
	}
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		return false, err
	}
	return true, nil
}

func hasCRLF(s string) bool { return strings.Contains(s, "\r\n") }

func normalizeNewlines(s string) string {
	// Convert CRLF to LF; keep lone CR as-is (rare).
	return strings.ReplaceAll(s, "\r\n", "\n")
}

func DiffHint(oldDoc, newDoc []byte) string {
	// Provide a small hint for failure messages without adding dependencies.
	oldN := normalizeNewlines(string(oldDoc))
	newN := normalizeNewlines(string(newDoc))
	if oldN == newN {
		return ""
	}
	// Keep it short: only show first differing line.
	oldLines := strings.Split(oldN, "\n")
	newLines := strings.Split(newN, "\n")
	n := len(oldLines)
	if len(newLines) < n {
		n = len(newLines)
	}
	for i := 0; i < n; i++ {
		if oldLines[i] != newLines[i] {
			return fmt.Sprintf("first diff at line %d:\n- %s\n+ %s\n", i+1, oldLines[i], newLines[i])
		}
	}
	if len(oldLines) != len(newLines) {
		return fmt.Sprintf("doc line count differs: %d vs %d\n", len(oldLines), len(newLines))
	}
	return ""
}

// BytesEqualNormalized is used by tests/CLI to compare docs ignoring CRLF differences.
func BytesEqualNormalized(a, b []byte) bool {
	return bytes.Equal([]byte(normalizeNewlines(string(a))), []byte(normalizeNewlines(string(b))))
}

