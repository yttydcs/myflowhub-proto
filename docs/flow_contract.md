# Flow Contract

Generated from `protocol/flow/types.go`. Do not hand-edit this file; run `go run ./cmd/flowcontractgen -write`.

## Enums

### NodeKind

- `call`
- `compose`
- `transform`
- `set_var`
- `branch`
- `foreach`
- `subflow`

### BindingSourceKind

- `node_result`
- `trigger`
- `flow_meta`
- `run_meta`
- `loop_item`
- `loop_index`
- `flow_var`

### BranchMatchOp

- `eq`
- `ne`
- `gt`
- `gte`
- `lt`
- `lte`
- `exists`

## Node Kind To Spec Mapping

| kind | spec type |
|---|---|
| `call` | `CallSpec` |
| `compose` | `ComposeSpec` |
| `transform` | `TransformSpec` |
| `set_var` | `SetVarSpec` |
| `branch` | `BranchSpec` |
| `foreach` | `ForeachSpec` |
| `subflow` | `SubflowSpec` |

## Types

### Node

| json | Go type | TS type | required | notes |
|---|---|---|---|---|
| `id` | `string` | `string` | yes |  |
| `kind` | `NodeKind` | `FlowNodeKind` | yes | Selects the spec type. |
| `allow_fail` | `bool` | `boolean` | no |  |
| `retry` | `*int` | `number` | no |  |
| `retry_backoff_ms` | `*int` | `number` | no |  |
| `timeout_ms` | `*int` | `number` | no |  |
| `spec` | `json.RawMessage` | `FlowNodeSpec` | yes | Decode with the type mapped from `kind`. |

### Trigger

| json | Go field | Go type | TS type | required | notes |
|---|---|---|---|---|---|
| `type` | `Type` | `string` | `string` | yes |  |
| `every_ms` | `EveryMs` | `uint64` | `number` | no |  |
| `cron` | `Cron` | `string` | `string` | no |  |
| `dedup_window_ms` | `DedupWindowMs` | `*int` | `number` | no |  |
| `event_mode` | `EventMode` | `string` | `string` | no |  |
| `event_name` | `EventName` | `string` | `string` | no |  |
| `event_topic` | `EventTopic` | `string` | `string` | no |  |
| `var_owner` | `VarOwner` | `uint32` | `number` | no |  |
| `var_name` | `VarName` | `string` | `string` | no |  |

### NodeUILayout

| json | Go field | Go type | TS type | required | notes |
|---|---|---|---|---|---|
| `x` | `X` | `int` | `number` | no |  |
| `y` | `Y` | `int` | `number` | no |  |

### Edge

| json | Go field | Go type | TS type | required | notes |
|---|---|---|---|---|---|
| `from` | `From` | `string` | `string` | yes |  |
| `to` | `To` | `string` | `string` | yes |  |
| `case` | `Case` | `string` | `string` | no | Only used on outgoing branch edges. |

### BindingSource

| json | Go field | Go type | TS type | required | notes |
|---|---|---|---|---|---|
| `kind` | `Kind` | `BindingSourceKind` | `FlowBindingSourceKind` | yes |  |
| `node_id` | `NodeID` | `string` | `string` | no |  |
| `name` | `Name` | `string` | `string` | no |  |
| `path` | `Path` | `string` | `string` | no |  |
| `field` | `Field` | `string` | `string` | no |  |

### InputBinding

| json | Go field | Go type | TS type | required | notes |
|---|---|---|---|---|---|
| `to` | `To` | `string` | `string` | yes |  |
| `source` | `Source` | `BindingSource` | `FlowBindingSource` | yes |  |
| `required` | `Required` | `bool` | `boolean` | no |  |

### ComposeSpec

| json | Go field | Go type | TS type | required | notes |
|---|---|---|---|---|---|
| `template` | `Template` | `json.RawMessage` | `JSONObject` | yes | Must be a JSON object in the stable contract. |
| `inputs` | `Inputs` | `[]InputBinding` | `FlowInputBinding[]` | no |  |
| `_ui` | `UI` | `*NodeUILayout` | `FlowNodeUILayout` | no |  |

### SetVarSpec

| json | Go field | Go type | TS type | required | notes |
|---|---|---|---|---|---|
| `name` | `Name` | `string` | `string` | yes |  |
| `template` | `Template` | `json.RawMessage` | `JSONValue` | no | May be any JSON value. |
| `inputs` | `Inputs` | `[]InputBinding` | `FlowInputBinding[]` | no |  |
| `_ui` | `UI` | `*NodeUILayout` | `FlowNodeUILayout` | no |  |

### TransformExpr

| json | Go field | Go type | TS type | required | notes |
|---|---|---|---|---|---|
| `literal` | `Literal` | `json.RawMessage` | `JSONValue` | no |  |
| `source` | `Source` | `*BindingSource` | `FlowBindingSource` | no |  |
| `required` | `Required` | `*bool` | `boolean` | no | Only applies when source is selected. |
| `op` | `Op` | `string` | `string` | no |  |
| `args` | `Args` | `[]TransformExpr` | `FlowTransformExpr[]` | no |  |
| `object` | `Object` | `map[string]TransformExpr` | `Record<string, FlowTransformExpr>` | no |  |
| `array` | `Array` | `[]TransformExpr` | `FlowTransformExpr[]` | no |  |

### TransformSpec

| json | Go field | Go type | TS type | required | notes |
|---|---|---|---|---|---|
| `expr` | `Expr` | `TransformExpr` | `FlowTransformExpr` | yes |  |
| `_ui` | `UI` | `*NodeUILayout` | `FlowNodeUILayout` | no |  |

### BranchMatch

| json | Go field | Go type | TS type | required | notes |
|---|---|---|---|---|---|
| `source` | `Source` | `BindingSource` | `FlowBindingSource` | yes |  |
| `op` | `Op` | `BranchMatchOp` | `FlowBranchMatchOp` | yes |  |
| `value` | `Value` | `json.RawMessage` | `JSONValue` | no | Omit only when op is exists. |

### BranchCase

| json | Go field | Go type | TS type | required | notes |
|---|---|---|---|---|---|
| `name` | `Name` | `string` | `string` | yes |  |
| `match` | `Match` | `BranchMatch` | `FlowBranchMatch` | yes |  |

### BranchSpec

| json | Go field | Go type | TS type | required | notes |
|---|---|---|---|---|---|
| `cases` | `Cases` | `[]BranchCase` | `FlowBranchCase[]` | yes |  |
| `default_case` | `DefaultCase` | `string` | `string` | no |  |
| `_ui` | `UI` | `*NodeUILayout` | `FlowNodeUILayout` | no |  |

### ForeachSpec

The body graph is isolated per item and the result is collected from result_node_id.

| json | Go field | Go type | TS type | required | notes |
|---|---|---|---|---|---|
| `source` | `Source` | `BindingSource` | `FlowBindingSource` | yes |  |
| `required` | `Required` | `*bool` | `boolean` | no |  |
| `body` | `Body` | `Graph` | `FlowGraph` | yes | Must be a valid child graph. |
| `result_node_id` | `ResultNodeID` | `string` | `string` | yes | Must exist inside the foreach body graph. |
| `_ui` | `UI` | `*NodeUILayout` | `FlowNodeUILayout` | no |  |

### SubflowSpec

Subflow nodes synchronously execute another flow on the same executor.

| json | Go field | Go type | TS type | required | notes |
|---|---|---|---|---|---|
| `flow_id` | `FlowID` | `string` | `string` | yes |  |
| `input_template` | `InputTemplate` | `json.RawMessage` | `JSONObject` | no | Must be a JSON object in the stable contract. |
| `inputs` | `Inputs` | `[]InputBinding` | `FlowInputBinding[]` | no |  |
| `result_node_id` | `ResultNodeID` | `string` | `string` | no |  |
| `_ui` | `UI` | `*NodeUILayout` | `FlowNodeUILayout` | no |  |

### CallSpec

New-write call nodes use args_template plus inputs. Legacy args are intentionally omitted.

| json | Go field | Go type | TS type | required | notes |
|---|---|---|---|---|---|
| `target` | `Target` | `uint32` | `number` | no |  |
| `method` | `Method` | `string` | `string` | yes |  |
| `args_template` | `ArgsTemplate` | `json.RawMessage` | `JSONObject` | no | Must be a JSON object in the stable contract. |
| `inputs` | `Inputs` | `[]InputBinding` | `FlowInputBinding[]` | no |  |
| `_ui` | `UI` | `*NodeUILayout` | `FlowNodeUILayout` | no |  |

### Graph

The top-level graph uses kind-tagged nodes and branch-aware edges.

| json | Go field | Go type | TS type | required | notes |
|---|---|---|---|---|---|
| `nodes` | `Nodes` | `[]Node` | `FlowNode[]` | yes | Each node uses kind-tagged specs. |
| `edges` | `Edges` | `[]Edge` | `FlowEdge[]` | yes |  |

