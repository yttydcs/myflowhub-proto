package protocolmap

// Context: This file supports the Proto source-of-truth workflow around model.

// ProtocolMap 是从 `protocol/*/types.go` 提取出来的“协议字典视图”（用于生成衍生物）。
//
// 说明：
// - 本仓库的 single source-of-truth 仍是 Go 源码；本模型仅用于文档生成与一致性校验。
// - model 设计保持“最小可用”，避免与 wire/实现耦合。
type ProtocolMap struct {
	Protocols []Protocol
}

type Protocol struct {
	Dir         string // e.g. "auth"
	Package     string // e.g. "auth"
	DisplayName string // e.g. "Auth"

	SubProtoConst string // e.g. "SubProtoAuth"
	SubProtoID    uint8  // e.g. 2

	Actions []Const // Action* constants (string)
	Consts  []Const // other exported consts (string/int)
	Types   []string
}

type ConstKind uint8

const (
	ConstKindUnknown ConstKind = iota
	ConstKindString
	ConstKindInt
)

type Const struct {
	Name string
	Kind ConstKind

	// For Kind==ConstKindString
	Str string

	// For Kind==ConstKindInt
	Int int64

	// Raw literal from Go source, e.g. `"login"` or `0x01`.
	Raw string
}
