package protocolmap

// 本文件承载协议映射生成流程中与 `doc` 相关的逻辑。

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProtocolMapDocUpToDate(t *testing.T) {
	repoRoot := findRepoRoot(t)
	protocolRoot := filepath.Join(repoRoot, "protocol")
	docPath := filepath.Join(repoRoot, "docs", "protocol_map.md")

	res, err := UpdateDocFile(protocolRoot, docPath)
	if err != nil {
		t.Fatalf("update doc: %v", err)
	}
	if !res.Changed {
		return
	}

	oldDoc, _ := os.ReadFile(docPath)
	t.Fatalf("docs/protocol_map.md is out of date.\nRun:\n  go run ./cmd/protocolmapgen -write -out docs/protocol_map.md\n\n%s", DiffHint(oldDoc, res.Doc))
}

func findRepoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("cannot find repo root (go.mod) from %s", wd)
		}
		dir = parent
	}
}
