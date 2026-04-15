package flow

// 本文件定义 Proto 中 `flow` 子协议的稳定类型与载荷结构。

import "encoding/json"

const SubProtoFlow uint8 = 6

const (
	ActionSet           = "set"
	ActionSetResp       = "set_resp"
	ActionRun           = "run"
	ActionRunResp       = "run_resp"
	ActionCancelRun     = "cancel_run"
	ActionCancelRunResp = "cancel_run_resp"
	ActionStatus        = "status"
	ActionStatusResp    = "status_resp"
	ActionDetail        = "detail"
	ActionDetailResp    = "detail_resp"
	ActionListRuns      = "list_runs"
	ActionListRunsResp  = "list_runs_resp"
	ActionList          = "list"
	ActionListResp      = "list_resp"
	ActionGet           = "get"
	ActionGetResp       = "get_resp"
	ActionDelete        = "delete"
	ActionDeleteResp    = "delete_resp"
)

const (
	PermFlowSet    = "flow.set"
	PermFlowDelete = "flow.delete"
	PermFlowRun    = "flow.run"
	PermFlowRead   = "flow.read"
)

type Message struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

type Trigger struct {
	Type          string `json:"type"`
	EveryMs       uint64 `json:"every_ms,omitempty"`
	Cron          string `json:"cron,omitempty"`
	DedupWindowMs *int   `json:"dedup_window_ms,omitempty"`
	EventMode     string `json:"event_mode,omitempty"`
	EventName     string `json:"event_name,omitempty"`
	EventTopic    string `json:"event_topic,omitempty"`
	VarOwner      uint32 `json:"var_owner,omitempty"`
	VarName       string `json:"var_name,omitempty"`
}

type NodeKind string

const (
	NodeKindCall      NodeKind = "call"
	NodeKindCompose   NodeKind = "compose"
	NodeKindTransform NodeKind = "transform"
	NodeKindSetVar    NodeKind = "set_var"
	NodeKindBranch    NodeKind = "branch"
	NodeKindForeach   NodeKind = "foreach"
	NodeKindSubflow   NodeKind = "subflow"
)

var SupportedNodeKinds = []NodeKind{
	NodeKindCall,
	NodeKindCompose,
	NodeKindTransform,
	NodeKindSetVar,
	NodeKindBranch,
	NodeKindForeach,
	NodeKindSubflow,
}

type BindingSourceKind string

const (
	BindingSourceNodeResult BindingSourceKind = "node_result"
	BindingSourceTrigger    BindingSourceKind = "trigger"
	BindingSourceFlowMeta   BindingSourceKind = "flow_meta"
	BindingSourceRunMeta    BindingSourceKind = "run_meta"
	BindingSourceLoopItem   BindingSourceKind = "loop_item"
	BindingSourceLoopIndex  BindingSourceKind = "loop_index"
	BindingSourceFlowVar    BindingSourceKind = "flow_var"
)

var SupportedBindingSourceKinds = []BindingSourceKind{
	BindingSourceNodeResult,
	BindingSourceTrigger,
	BindingSourceFlowMeta,
	BindingSourceRunMeta,
	BindingSourceLoopItem,
	BindingSourceLoopIndex,
	BindingSourceFlowVar,
}

type BranchMatchOp string

const (
	BranchMatchEq     BranchMatchOp = "eq"
	BranchMatchNe     BranchMatchOp = "ne"
	BranchMatchGt     BranchMatchOp = "gt"
	BranchMatchGte    BranchMatchOp = "gte"
	BranchMatchLt     BranchMatchOp = "lt"
	BranchMatchLte    BranchMatchOp = "lte"
	BranchMatchExists BranchMatchOp = "exists"
)

var SupportedBranchMatchOps = []BranchMatchOp{
	BranchMatchEq,
	BranchMatchNe,
	BranchMatchGt,
	BranchMatchGte,
	BranchMatchLt,
	BranchMatchLte,
	BranchMatchExists,
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

type Node struct {
	ID             string          `json:"id"`
	Kind           NodeKind        `json:"kind"`
	AllowFail      bool            `json:"allow_fail,omitempty"`
	Retry          *int            `json:"retry,omitempty"`
	RetryBackoffMs *int            `json:"retry_backoff_ms,omitempty"`
	TimeoutMs      *int            `json:"timeout_ms,omitempty"`
	Spec           json.RawMessage `json:"spec"`
}

type Edge struct {
	From string `json:"from"`
	To   string `json:"to"`
	Case string `json:"case,omitempty"`
}

type NodeUILayout struct {
	X int `json:"x,omitempty"`
	Y int `json:"y,omitempty"`
}

type BindingSource struct {
	Kind   BindingSourceKind `json:"kind"`
	NodeID string            `json:"node_id,omitempty"`
	Name   string            `json:"name,omitempty"`
	Path   string            `json:"path,omitempty"`
	Field  string            `json:"field,omitempty"`
}

type InputBinding struct {
	To       string        `json:"to"`
	Source   BindingSource `json:"source"`
	Required bool          `json:"required,omitempty"`
}

type CallSpec struct {
	Target       uint32          `json:"target,omitempty"`
	Method       string          `json:"method"`
	ArgsTemplate json.RawMessage `json:"args_template,omitempty"`
	Inputs       []InputBinding  `json:"inputs,omitempty"`
	UI           *NodeUILayout   `json:"_ui,omitempty"`
}

type ComposeSpec struct {
	Template json.RawMessage `json:"template"`
	Inputs   []InputBinding  `json:"inputs,omitempty"`
	UI       *NodeUILayout   `json:"_ui,omitempty"`
}

type SetVarSpec struct {
	Name     string          `json:"name"`
	Template json.RawMessage `json:"template,omitempty"`
	Inputs   []InputBinding  `json:"inputs,omitempty"`
	UI       *NodeUILayout   `json:"_ui,omitempty"`
}

type TransformExpr struct {
	Literal  json.RawMessage          `json:"literal,omitempty"`
	Source   *BindingSource           `json:"source,omitempty"`
	Required *bool                    `json:"required,omitempty"`
	Op       string                   `json:"op,omitempty"`
	Args     []TransformExpr          `json:"args,omitempty"`
	Object   map[string]TransformExpr `json:"object,omitempty"`
	Array    []TransformExpr          `json:"array,omitempty"`
}

type TransformSpec struct {
	Expr TransformExpr `json:"expr"`
	UI   *NodeUILayout `json:"_ui,omitempty"`
}

type BranchMatch struct {
	Source BindingSource   `json:"source"`
	Op     BranchMatchOp   `json:"op"`
	Value  json.RawMessage `json:"value,omitempty"`
}

type BranchCase struct {
	Name  string      `json:"name"`
	Match BranchMatch `json:"match"`
}

type BranchSpec struct {
	Cases       []BranchCase  `json:"cases"`
	DefaultCase string        `json:"default_case,omitempty"`
	UI          *NodeUILayout `json:"_ui,omitempty"`
}

type ForeachSpec struct {
	Source       BindingSource `json:"source"`
	Required     *bool         `json:"required,omitempty"`
	Body         Graph         `json:"body"`
	ResultNodeID string        `json:"result_node_id"`
	UI           *NodeUILayout `json:"_ui,omitempty"`
}

type SubflowSpec struct {
	FlowID        string          `json:"flow_id"`
	InputTemplate json.RawMessage `json:"input_template,omitempty"`
	Inputs        []InputBinding  `json:"inputs,omitempty"`
	ResultNodeID  string          `json:"result_node_id,omitempty"`
	UI            *NodeUILayout   `json:"_ui,omitempty"`
}

type SetReq struct {
	ReqID         string  `json:"req_id"`
	OriginNode    uint32  `json:"origin_node,omitempty"`
	ExecutorNode  uint32  `json:"executor_node,omitempty"`
	FlowID        string  `json:"flow_id"`
	Name          string  `json:"name,omitempty"`
	MaxActiveRuns *int    `json:"max_active_runs,omitempty"`
	Trigger       Trigger `json:"trigger"`
	Graph         Graph   `json:"graph"`
}

type SetResp struct {
	ReqID  string `json:"req_id"`
	Code   int    `json:"code"`
	Msg    string `json:"msg,omitempty"`
	FlowID string `json:"flow_id,omitempty"`
}

type DeleteReq struct {
	ReqID        string `json:"req_id"`
	OriginNode   uint32 `json:"origin_node,omitempty"`
	ExecutorNode uint32 `json:"executor_node,omitempty"`
	FlowID       string `json:"flow_id"`
}

type DeleteResp struct {
	ReqID  string `json:"req_id"`
	Code   int    `json:"code"`
	Msg    string `json:"msg,omitempty"`
	FlowID string `json:"flow_id,omitempty"`
}

type RunReq struct {
	ReqID        string `json:"req_id"`
	OriginNode   uint32 `json:"origin_node,omitempty"`
	ExecutorNode uint32 `json:"executor_node,omitempty"`
	FlowID       string `json:"flow_id"`
}

type RunResp struct {
	ReqID  string `json:"req_id"`
	Code   int    `json:"code"`
	Msg    string `json:"msg,omitempty"`
	FlowID string `json:"flow_id,omitempty"`
	RunID  string `json:"run_id,omitempty"`
}

type CancelRunReq struct {
	ReqID        string `json:"req_id"`
	OriginNode   uint32 `json:"origin_node,omitempty"`
	ExecutorNode uint32 `json:"executor_node,omitempty"`
	FlowID       string `json:"flow_id"`
	RunID        string `json:"run_id"`
}

type CancelRunResp struct {
	ReqID        string `json:"req_id"`
	Code         int    `json:"code"`
	Msg          string `json:"msg,omitempty"`
	ExecutorNode uint32 `json:"executor_node,omitempty"`
	FlowID       string `json:"flow_id,omitempty"`
	RunID        string `json:"run_id,omitempty"`
	Status       string `json:"status,omitempty"`
}

type StatusReq struct {
	ReqID        string `json:"req_id"`
	OriginNode   uint32 `json:"origin_node,omitempty"`
	ExecutorNode uint32 `json:"executor_node,omitempty"`
	FlowID       string `json:"flow_id"`
	RunID        string `json:"run_id,omitempty"`
}

type NodeStatus struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Code   int    `json:"code,omitempty"`
	Msg    string `json:"msg,omitempty"`
}

type StatusResp struct {
	ReqID        string       `json:"req_id"`
	Code         int          `json:"code"`
	Msg          string       `json:"msg,omitempty"`
	ExecutorNode uint32       `json:"executor_node,omitempty"`
	FlowID       string       `json:"flow_id,omitempty"`
	RunID        string       `json:"run_id,omitempty"`
	Status       string       `json:"status,omitempty"`
	Nodes        []NodeStatus `json:"nodes,omitempty"`
}

type DetailReq struct {
	ReqID        string `json:"req_id"`
	OriginNode   uint32 `json:"origin_node,omitempty"`
	ExecutorNode uint32 `json:"executor_node,omitempty"`
	FlowID       string `json:"flow_id"`
	RunID        string `json:"run_id,omitempty"`
	NodeID       string `json:"node_id"`
	Path         string `json:"path,omitempty"`
}

type DetailResp struct {
	ReqID        string          `json:"req_id"`
	Code         int             `json:"code"`
	Msg          string          `json:"msg,omitempty"`
	ExecutorNode uint32          `json:"executor_node,omitempty"`
	FlowID       string          `json:"flow_id,omitempty"`
	RunID        string          `json:"run_id,omitempty"`
	Path         string          `json:"path,omitempty"`
	Node         *NodeStatus     `json:"node,omitempty"`
	Result       json.RawMessage `json:"result,omitempty"`
}

type ListRunsReq struct {
	ReqID        string `json:"req_id"`
	OriginNode   uint32 `json:"origin_node,omitempty"`
	ExecutorNode uint32 `json:"executor_node,omitempty"`
	FlowID       string `json:"flow_id"`
	Limit        uint32 `json:"limit,omitempty"`
}

type RunSummary struct {
	RunID       string `json:"run_id"`
	Status      string `json:"status"`
	StartedAtMs int64  `json:"started_at_ms,omitempty"`
	EndedAtMs   int64  `json:"ended_at_ms,omitempty"`
	Msg         string `json:"msg,omitempty"`
}

type ListRunsResp struct {
	ReqID        string       `json:"req_id"`
	Code         int          `json:"code"`
	Msg          string       `json:"msg,omitempty"`
	ExecutorNode uint32       `json:"executor_node,omitempty"`
	FlowID       string       `json:"flow_id,omitempty"`
	Runs         []RunSummary `json:"runs,omitempty"`
}

type ListReq struct {
	ReqID        string `json:"req_id"`
	OriginNode   uint32 `json:"origin_node,omitempty"`
	ExecutorNode uint32 `json:"executor_node,omitempty"`
}

type FlowSummary struct {
	FlowID     string `json:"flow_id"`
	Name       string `json:"name,omitempty"`
	EveryMs    uint64 `json:"every_ms,omitempty"`
	LastRunID  string `json:"last_run_id,omitempty"`
	LastStatus string `json:"last_status,omitempty"`
}

type ListResp struct {
	ReqID        string        `json:"req_id"`
	Code         int           `json:"code"`
	Msg          string        `json:"msg,omitempty"`
	ExecutorNode uint32        `json:"executor_node,omitempty"`
	Flows        []FlowSummary `json:"flows,omitempty"`
}

type GetReq struct {
	ReqID        string `json:"req_id"`
	OriginNode   uint32 `json:"origin_node,omitempty"`
	ExecutorNode uint32 `json:"executor_node,omitempty"`
	FlowID       string `json:"flow_id"`
}

type GetResp struct {
	ReqID         string  `json:"req_id"`
	Code          int     `json:"code"`
	Msg           string  `json:"msg,omitempty"`
	ExecutorNode  uint32  `json:"executor_node,omitempty"`
	FlowID        string  `json:"flow_id,omitempty"`
	Name          string  `json:"name,omitempty"`
	MaxActiveRuns *int    `json:"max_active_runs,omitempty"`
	Trigger       Trigger `json:"trigger,omitempty"`
	Graph         Graph   `json:"graph,omitempty"`
}
