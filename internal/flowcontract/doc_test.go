package flowcontract

// 本文件承载 Flow 合同生成流程中与 `doc` 相关的逻辑。

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGeneratedFilesUpToDate(t *testing.T) {
	repoRoot := findRepoRoot(t)
	mdPath := filepath.Join(repoRoot, "docs", "flow_contract.md")
	tsPath := filepath.Join(repoRoot, "generated", "flow_contract.ts")

	md, ts, err := RenderAll()
	if err != nil {
		t.Fatalf("render outputs: %v", err)
	}

	mdChanged, err := FileChanged(mdPath, md)
	if err != nil {
		t.Fatalf("check markdown: %v", err)
	}
	tsChanged, err := FileChanged(tsPath, ts)
	if err != nil {
		t.Fatalf("check typescript: %v", err)
	}
	if !mdChanged && !tsChanged {
		return
	}

	t.Fatalf("flow contract outputs are out of date.\nRun:\n  go run ./cmd/flowcontractgen -write -md-out docs/flow_contract.md -ts-out generated/flow_contract.ts")
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
