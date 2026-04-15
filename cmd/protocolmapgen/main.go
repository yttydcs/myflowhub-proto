package main

// 本文件提供 `protocolmapgen` 命令入口，用于从 Proto 真相源生成协议映射文档。

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yttydcs/myflowhub-proto/internal/protocolmap"
)

// main 解析命令行参数，并执行协议映射文档的写入或一致性校验。
func main() {
	var (
		protocolRoot = flag.String("protocol-root", "protocol", "protocol root directory (contains subproto folders)")
		outPath      = flag.String("out", filepath.FromSlash("docs/protocol_map.md"), "output markdown path")
		write        = flag.Bool("write", false, "write updated doc to -out")
		check        = flag.Bool("check", false, "check doc is up-to-date (exit non-zero if not)")
	)
	flag.Parse()

	if (*write && *check) || (!*write && !*check) {
		_, _ = fmt.Fprintln(os.Stderr, "usage: protocolmapgen -write|-check [-protocol-root protocol] [-out docs/protocol_map.md]")
		os.Exit(2)
	}

	res, err := protocolmap.UpdateDocFile(*protocolRoot, *outPath)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	if *check {
		if res.Changed {
			_, _ = fmt.Fprintln(os.Stderr, "docs/protocol_map.md is out of date.")
			_, _ = fmt.Fprintln(os.Stderr, "Run: go run ./cmd/protocolmapgen -write -out docs/protocol_map.md")
			os.Exit(1)
		}
		return
	}

	changed, err := protocolmap.WriteFileIfChanged(*outPath, res.Doc)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "write doc:", err)
		os.Exit(1)
	}
	if changed {
		_, _ = fmt.Fprintln(os.Stdout, "updated:", *outPath)
	} else {
		_, _ = fmt.Fprintln(os.Stdout, "no change:", *outPath)
	}
}
