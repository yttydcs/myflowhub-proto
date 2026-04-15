package protocolmap

// 本文件承载协议映射生成流程中与 `parse` 相关的逻辑。

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var displayNameOverrides = map[string]string{
	"management": "Management",
	"auth":       "Auth",
	"varstore":   "VarStore",
	"topicbus":   "TopicBus",
	"file":       "File",
	"flow":       "Flow",
	"exec":       "Exec",
}

// ParseProtocolRoot 扫描 `protocol/` 根目录并汇总每个子协议的映射信息。
func ParseProtocolRoot(protocolRoot string) (*ProtocolMap, error) {
	protocolRoot = filepath.Clean(protocolRoot)
	entries, err := os.ReadDir(protocolRoot)
	if err != nil {
		return nil, fmt.Errorf("read protocol root: %w", err)
	}
	var dirs []string
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e.Name())
		}
	}
	sort.Strings(dirs)

	var out ProtocolMap
	for _, dir := range dirs {
		p, err := parseProtocolDir(filepath.Join(protocolRoot, dir), dir)
		if err != nil {
			return nil, err
		}
		out.Protocols = append(out.Protocols, *p)
	}

	// Validate: unique SubProto id.
	seen := make(map[uint8]string, len(out.Protocols))
	for _, p := range out.Protocols {
		if p.SubProtoConst == "" {
			return nil, fmt.Errorf("protocol %q missing SubProto const", p.Dir)
		}
		if prev, ok := seen[p.SubProtoID]; ok {
			return nil, fmt.Errorf("duplicate SubProto id %d: %q and %q", p.SubProtoID, prev, p.Dir)
		}
		seen[p.SubProtoID] = p.Dir
	}

	sort.Slice(out.Protocols, func(i, j int) bool {
		return out.Protocols[i].SubProtoID < out.Protocols[j].SubProtoID
	})
	return &out, nil
}

// parseProtocolDir 解析单个子协议目录中的包名、常量、动作和值类型列表。
func parseProtocolDir(dirPath string, dirName string) (*Protocol, error) {
	files, err := listGoFiles(dirPath)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("protocol dir %q has no .go files", dirName)
	}

	fset := token.NewFileSet()
	var pkgName string

	declaredTypes := make(map[string]struct{})
	var consts []Const
	constValues := make(map[string]Const) // best-effort for ident resolution within package

	for _, f := range files {
		af, err := parser.ParseFile(fset, f, nil, 0)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", f, err)
		}
		if pkgName == "" {
			pkgName = af.Name.Name
		} else if pkgName != af.Name.Name {
			return nil, fmt.Errorf("protocol dir %q has multiple packages: %q and %q", dirName, pkgName, af.Name.Name)
		}

		for _, decl := range af.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			switch gen.Tok {
			case token.TYPE:
				for _, spec := range gen.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					if !ts.Name.IsExported() {
						continue
					}
					declaredTypes[ts.Name.Name] = struct{}{}
				}
			case token.CONST:
				parsed, err := parseConstDecl(gen, constValues)
				if err != nil {
					return nil, fmt.Errorf("protocol dir %q const parse: %w", dirName, err)
				}
				for _, c := range parsed {
					// keep only exported consts
					if ast.IsExported(c.Name) {
						consts = append(consts, c)
						constValues[c.Name] = c
					}
				}
			}
		}
	}

	p := &Protocol{
		Dir:     dirName,
		Package: pkgName,
	}
	if dn, ok := displayNameOverrides[dirName]; ok {
		p.DisplayName = dn
	} else {
		p.DisplayName = upperFirst(dirName)
	}

	// Extract SubProto const.
	for _, c := range consts {
		if strings.HasPrefix(c.Name, "SubProto") {
			if p.SubProtoConst != "" {
				return nil, fmt.Errorf("protocol %q has multiple SubProto consts: %q and %q", dirName, p.SubProtoConst, c.Name)
			}
			if c.Kind != ConstKindInt {
				return nil, fmt.Errorf("protocol %q SubProto const %q is not int", dirName, c.Name)
			}
			if c.Int < 0 || c.Int > 255 {
				return nil, fmt.Errorf("protocol %q SubProto id out of range: %d", dirName, c.Int)
			}
			p.SubProtoConst = c.Name
			p.SubProtoID = uint8(c.Int)
		}
	}

	// Extract actions and other consts.
	for _, c := range consts {
		switch {
		case strings.HasPrefix(c.Name, "SubProto"):
			// skip (already handled)
		case strings.HasPrefix(c.Name, "Action"):
			if c.Kind != ConstKindString {
				return nil, fmt.Errorf("protocol %q action %q is not string", dirName, c.Name)
			}
			p.Actions = append(p.Actions, c)
		default:
			p.Consts = append(p.Consts, c)
		}
	}

	sort.Slice(p.Actions, func(i, j int) bool {
		if p.Actions[i].Str == p.Actions[j].Str {
			return p.Actions[i].Name < p.Actions[j].Name
		}
		return p.Actions[i].Str < p.Actions[j].Str
	})
	sort.Slice(p.Consts, func(i, j int) bool {
		return p.Consts[i].Name < p.Consts[j].Name
	})

	for tn := range declaredTypes {
		p.Types = append(p.Types, tn)
	}
	sort.Strings(p.Types)

	// basic sanity: action value uniqueness within protocol
	actionSeen := map[string]string{}
	for _, a := range p.Actions {
		if prev, ok := actionSeen[a.Str]; ok {
			return nil, fmt.Errorf("protocol %q duplicate action value %q: %s and %s", dirName, a.Str, prev, a.Name)
		}
		actionSeen[a.Str] = a.Name
	}

	return p, nil
}

// listGoFiles 返回目录下参与协议扫描的非测试 Go 文件。
func listGoFiles(dirPath string) ([]string, error) {
	ents, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("read dir %s: %w", dirPath, err)
	}
	var files []string
	for _, e := range ents {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		files = append(files, filepath.Join(dirPath, name))
	}
	sort.Strings(files)
	return files, nil
}

// parseConstDecl 解析一个 const 声明块，并尽量沿用前一条隐式赋值规则。
func parseConstDecl(gen *ast.GenDecl, known map[string]Const) ([]Const, error) {
	if gen == nil || gen.Tok != token.CONST {
		return nil, nil
	}
	var out []Const

	var prevValues []ast.Expr
	var prevType ast.Expr

	for _, spec := range gen.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		values := vs.Values
		typ := vs.Type
		if len(values) == 0 {
			values = prevValues
			typ = prevType
		} else {
			prevValues = values
			prevType = typ
		}

		if len(values) == 0 {
			// implicit const without previous expression (invalid or unsupported)
			return nil, errors.New("const value missing (implicit values not supported without previous expression)")
		}
		if len(values) != len(vs.Names) {
			return nil, fmt.Errorf("const names/values mismatch: %d names, %d values", len(vs.Names), len(values))
		}

		for i, n := range vs.Names {
			if n == nil {
				continue
			}
			val, err := evalConstExpr(values[i], known)
			if err != nil {
				// If explicit type is present, still treat as const (type info is not required for doc).
				_ = typ
				return nil, fmt.Errorf("const %s: %w", n.Name, err)
			}
			val.Name = n.Name
			out = append(out, val)
			known[n.Name] = val
		}
	}
	return out, nil
}

// evalConstExpr 计算协议映射生成器支持的常量字面量与标识符引用。
func evalConstExpr(expr ast.Expr, known map[string]Const) (Const, error) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		switch e.Kind {
		case token.STRING:
			s, err := strconv.Unquote(e.Value)
			if err != nil {
				return Const{}, err
			}
			return Const{Kind: ConstKindString, Str: s, Raw: e.Value}, nil
		case token.INT:
			i, err := strconv.ParseInt(e.Value, 0, 64)
			if err != nil {
				return Const{}, err
			}
			return Const{Kind: ConstKindInt, Int: i, Raw: e.Value}, nil
		case token.CHAR:
			r, _, tail, err := strconv.UnquoteChar(e.Value, '\'')
			if err != nil {
				return Const{}, err
			}
			if tail != "" {
				return Const{}, fmt.Errorf("invalid char literal: %s", e.Value)
			}
			return Const{Kind: ConstKindInt, Int: int64(r), Raw: e.Value}, nil
		default:
			return Const{}, fmt.Errorf("unsupported literal kind: %v", e.Kind)
		}
	case *ast.UnaryExpr:
		if e.Op != token.ADD && e.Op != token.SUB {
			return Const{}, fmt.Errorf("unsupported unary op: %s", e.Op.String())
		}
		v, err := evalConstExpr(e.X, known)
		if err != nil {
			return Const{}, err
		}
		if v.Kind != ConstKindInt {
			return Const{}, fmt.Errorf("unary op on non-int const")
		}
		if e.Op == token.SUB {
			v.Int = -v.Int
			v.Raw = "-" + v.Raw
		}
		return v, nil
	case *ast.Ident:
		if e.Name == "iota" {
			return Const{}, fmt.Errorf("iota is not supported (use explicit literals)")
		}
		if v, ok := known[e.Name]; ok {
			// preserve the ident reference as raw for readability
			v.Raw = e.Name
			return v, nil
		}
		return Const{}, fmt.Errorf("unknown ident: %s", e.Name)
	default:
		return Const{}, fmt.Errorf("unsupported const expr: %T", expr)
	}
}

// upperFirst 为缺省展示名提供最小的人类可读格式。
func upperFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
