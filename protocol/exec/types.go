package exec

// 本文件定义 Proto 中 `exec` 子协议的稳定类型与载荷结构。

import "encoding/json"

const SubProtoExec uint8 = 7

const (
	ActionCall     = "call"
	ActionCallResp = "call_resp"

	// 能力注册中心（逐级聚合）动作。
	ActionCapSnapshot  = "cap_snapshot"
	ActionCapUpsert    = "cap_upsert"
	ActionCapWithdraw  = "cap_withdraw"
	ActionCapHeartbeat = "cap_heartbeat"
	ActionCapSyncResp  = "cap_sync_resp"
	ActionCapQuery     = "cap_query"
	ActionCapQueryResp = "cap_query_resp"
)

const (
	PermExecCall     = "exec.call"
	PermExecCapSync  = "exec.cap.sync"
	PermExecCapQuery = "exec.cap.query"
)

type Message struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

type CallReq struct {
	ReqID        string          `json:"req_id"`
	ExecutorNode uint32          `json:"executor_node"`
	TargetNode   uint32          `json:"target_node"`
	Method       string          `json:"method"`
	Args         json.RawMessage `json:"args,omitempty"`
	TimeoutMs    int             `json:"timeout_ms,omitempty"`
}

type CallResp struct {
	ReqID        string          `json:"req_id"`
	Code         int             `json:"code"`
	Msg          string          `json:"msg,omitempty"`
	ExecutorNode uint32          `json:"executor_node,omitempty"`
	TargetNode   uint32          `json:"target_node,omitempty"`
	Method       string          `json:"method,omitempty"`
	Result       json.RawMessage `json:"result,omitempty"`
}

// CapabilityDescriptor 描述一个 provider 节点提供的能力。
// 说明：
// - provider_node 是真正执行 method 的节点。
// - input/output_schema 建议使用 JSON Schema 子集。
type CapabilityDescriptor struct {
	ProviderNode     uint32            `json:"provider_node"`
	Method           string            `json:"method"`
	Version          string            `json:"version,omitempty"`
	InputSchema      json.RawMessage   `json:"input_schema,omitempty"`
	OutputSchema     json.RawMessage   `json:"output_schema,omitempty"`
	DefaultTimeoutMs int               `json:"default_timeout_ms,omitempty"`
	Permissions      []string          `json:"permissions,omitempty"`
	Tags             map[string]string `json:"tags,omitempty"`
}

// CapabilityKey 用于精准撤销某一条能力记录。
type CapabilityKey struct {
	ProviderNode uint32 `json:"provider_node"`
	Method       string `json:"method"`
	Version      string `json:"version,omitempty"`
}

// CapSnapshotReq 用于子节点向父节点上报全量能力快照（带 epoch）。
type CapSnapshotReq struct {
	ReqID    string                 `json:"req_id,omitempty"`
	FromNode uint32                 `json:"from_node"`
	Epoch    uint64                 `json:"epoch"`
	LeaseMs  uint64                 `json:"lease_ms,omitempty"`
	Caps     []CapabilityDescriptor `json:"caps"`
}

// CapUpsertReq 用于增量上报新增/更新能力。
type CapUpsertReq struct {
	ReqID    string                 `json:"req_id,omitempty"`
	FromNode uint32                 `json:"from_node"`
	Epoch    uint64                 `json:"epoch"`
	LeaseMs  uint64                 `json:"lease_ms,omitempty"`
	Caps     []CapabilityDescriptor `json:"caps"`
}

// CapWithdrawReq 用于增量撤销能力。
type CapWithdrawReq struct {
	ReqID    string          `json:"req_id,omitempty"`
	FromNode uint32          `json:"from_node"`
	Epoch    uint64          `json:"epoch"`
	Keys     []CapabilityKey `json:"keys"`
}

// CapHeartbeatReq 用于续租（不携带能力内容）。
type CapHeartbeatReq struct {
	ReqID    string `json:"req_id,omitempty"`
	FromNode uint32 `json:"from_node"`
	Epoch    uint64 `json:"epoch"`
	LeaseMs  uint64 `json:"lease_ms,omitempty"`
}

// CapSyncResp 是 cap_snapshot/cap_upsert/cap_withdraw/cap_heartbeat 的统一响应。
type CapSyncResp struct {
	ReqID     string `json:"req_id,omitempty"`
	Code      int    `json:"code"`
	Msg       string `json:"msg,omitempty"`
	FromNode  uint32 `json:"from_node,omitempty"`
	Epoch     uint64 `json:"epoch,omitempty"`
	Applied   int    `json:"applied,omitempty"`
	Responder uint32 `json:"responder_node,omitempty"`
}

// CapQueryReq 查询聚合能力索引。
type CapQueryReq struct {
	ReqID         string `json:"req_id"`
	RequesterNode uint32 `json:"requester_node,omitempty"`
	Method        string `json:"method,omitempty"`
	Prefix        bool   `json:"prefix,omitempty"`
	ProviderNode  uint32 `json:"provider_node,omitempty"`
	Limit         int    `json:"limit,omitempty"`
	IncludeSchema bool   `json:"include_schema,omitempty"`
}

// CapabilityRoute 表示某能力的候选 provider 及可用路由信息。
type CapabilityRoute struct {
	ProviderNode     uint32            `json:"provider_node"`
	ViaNode          uint32            `json:"via_node,omitempty"`
	Method           string            `json:"method"`
	Version          string            `json:"version,omitempty"`
	DefaultTimeoutMs int               `json:"default_timeout_ms,omitempty"`
	Permissions      []string          `json:"permissions,omitempty"`
	Tags             map[string]string `json:"tags,omitempty"`
	LeaseExpireAt    int64             `json:"lease_expire_at,omitempty"`
	InputSchema      json.RawMessage   `json:"input_schema,omitempty"`
	OutputSchema     json.RawMessage   `json:"output_schema,omitempty"`
}

type CapQueryResp struct {
	ReqID         string            `json:"req_id"`
	Code          int               `json:"code"`
	Msg           string            `json:"msg,omitempty"`
	ResponderNode uint32            `json:"responder_node,omitempty"`
	Total         int               `json:"total,omitempty"`
	Routes        []CapabilityRoute `json:"routes,omitempty"`
}
