package main

// 本文件提供 `flowcontractgen` 命令入口，用于从 Proto 真相源生成 Flow 合同文档。

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yttydcs/myflowhub-proto/internal/flowcontract"
)

// main 解析命令行参数，并执行 Flow 合同的写入或一致性校验。
func main() {
	var (
		mdOut = flag.String("md-out", filepath.FromSlash("docs/flow_contract.md"), "output markdown path")
		tsOut = flag.String("ts-out", filepath.FromSlash("generated/flow_contract.ts"), "output TypeScript path")
		write = flag.Bool("write", false, "write generated outputs")
		check = flag.Bool("check", false, "check generated outputs are up to date")
	)
	flag.Parse()

	if (*write && *check) || (!*write && !*check) {
		_, _ = fmt.Fprintln(os.Stderr, "usage: flowcontractgen -write|-check [-md-out docs/flow_contract.md] [-ts-out generated/flow_contract.ts]")
		os.Exit(2)
	}

	md, ts, err := flowcontract.RenderAll()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	if *check {
		mdChanged, err := flowcontract.FileChanged(*mdOut, md)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "check markdown:", err)
			os.Exit(1)
		}
		tsChanged, err := flowcontract.FileChanged(*tsOut, ts)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "check typescript:", err)
			os.Exit(1)
		}
		if mdChanged || tsChanged {
			_, _ = fmt.Fprintln(os.Stderr, "flow contract outputs are out of date.")
			_, _ = fmt.Fprintln(os.Stderr, "Run: go run ./cmd/flowcontractgen -write -md-out docs/flow_contract.md -ts-out generated/flow_contract.ts")
			os.Exit(1)
		}
		return
	}

	mdChanged, err := flowcontract.WriteFileIfChanged(*mdOut, md)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "write markdown:", err)
		os.Exit(1)
	}
	tsChanged, err := flowcontract.WriteFileIfChanged(*tsOut, ts)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "write typescript:", err)
		os.Exit(1)
	}

	if mdChanged {
		_, _ = fmt.Fprintln(os.Stdout, "updated:", *mdOut)
	} else {
		_, _ = fmt.Fprintln(os.Stdout, "no change:", *mdOut)
	}
	if tsChanged {
		_, _ = fmt.Fprintln(os.Stdout, "updated:", *tsOut)
	} else {
		_, _ = fmt.Fprintln(os.Stdout, "no change:", *tsOut)
	}
}
