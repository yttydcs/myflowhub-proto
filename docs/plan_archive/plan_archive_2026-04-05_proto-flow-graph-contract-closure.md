# Plan - graph-contract-closure

## Workflow Information
- Repo: `D:/project/MyFlowHub3/repo/MyFlowHub-Proto` (primary control repo for this workflow)
- Branch: `refactor/graph-contract-closure`
- Base: `main`
- Worktree: `D:/project/MyFlowHub3/worktrees/proto-graph-contract-closure`
- Current Stage: `4`

## Stage Records

### Initialization
- guide.md: read from `D:/project/MyFlowHub3/guide.md`
- base/worktree confirmation:
  - Primary repo: `MyFlowHub-Proto` -> `D:/project/MyFlowHub3/worktrees/proto-graph-contract-closure`
  - Additional repo: `MyFlowHub-SubProto` -> `D:/project/MyFlowHub3/worktrees/subproto-graph-contract-closure`
  - Additional repo: `MyFlowHub-Server` -> `D:/project/MyFlowHub3/worktrees/server-graph-contract-closure`
  - Additional repo: `MyFlowHub-Win` -> `D:/project/MyFlowHub3/worktrees/win-graph-contract-closure`
  - All branches: `refactor/graph-contract-closure`
  - All bases: `main`
  - Control-plane repo paths remain read-only except for worktree management, merge, release, and convergence.

### Stage 1 - Requirements Analysis
#### Goal
- Close the remaining `flow graph/node contract` gap by making the graph/node wire contract canonical, typed, and auditable across Proto, SubProto, Server, and Win without reintroducing legacy compatibility.

#### Scope
- 必须
  - Canonicalize Flow graph/node wire contract in `MyFlowHub-Proto`.
  - Remove private duplicated graph/node contract structs from `MyFlowHub-SubProto/flow` and consume Proto definitions directly.
  - Make Win consume generated or synced contract artifacts instead of maintaining scattered hard-coded node-kind/source/op truth.
  - Update Server/Win specs so the written contract matches shipped code.
  - Add regression tests that prove advanced graph/node contracts round-trip and validate consistently.
- 可选
  - Add generated human-readable canonical Flow contract documentation from Proto if that is the cleanest way to expose the new truth to Server/Win docs.
  - Add contract sync tooling when it materially reduces future drift.
- 不做
  - No new runtime node semantics beyond the current `call/compose/transform/set_var/branch/foreach/subflow` contract.
  - No legacy `local` / `exec` compatibility.
  - No breaking action renames or SubProto number changes.
  - No new product UX outside the flow editor contract consumption work needed for this closure.

#### Use Cases
- Server `flow.set/get` must accept and return the same advanced graph/node contract for `branch/foreach/subflow/transform/set_var`.
- SubProto runtime must validate and execute advanced graph/node specs without relying on repo-local duplicate contract structs.
- Win editor must author and reload advanced graph/node specs against a canonical contract source rather than handwritten local literals.
- Docs readers must be able to find one canonical place for Flow graph/node contract truth and see downstream mirrored specs updated accordingly.

#### Functional Requirements
- `flow` graph node kinds remain `call | compose | transform | set_var | branch | foreach | subflow`.
- Each node kind must have an explicit typed spec model in Proto, including binding sources and branch ops.
- `graph.edges[].case` rules must remain explicit and branch-only.
- `foreach.body` and `result_node_id` contract must remain explicit and validated.
- `subflow.flow_id / input_template / inputs / result_node_id` contract must remain explicit and validated.
- Win authoring must consume canonical node-kind/source/op/spec-shape truth from generated or synced contract artifacts.
- Server stable specs must describe the same contract actually consumed by code.

#### Non-functional Requirements
- Respect repo boundaries in `repos.md`: Proto holds protocol definitions and generated contract metadata, not runtime handler logic.
- Keep change surface minimal and avoid plan-external behavior changes.
- Preserve existing stable wire field names and semantics.
- Provide deterministic contract generation/sync steps so drift is discoverable in CI or local tests.

#### Inputs / Outputs
- Inputs:
  - `D:/project/MyFlowHub3/repo/MyFlowHub-Server/docs/requirements/flow_data_dag.md`
  - `D:/project/MyFlowHub3/repo/MyFlowHub-Server/docs/specs/flow.md`
  - `D:/project/MyFlowHub3/repo/MyFlowHub-Win/docs/specs/flow-editor-visual-form.md`
  - Current implementations in `MyFlowHub-Proto`, `MyFlowHub-SubProto/flow`, `MyFlowHub-Win/frontend/src/stores/flow.ts`, and `MyFlowHub-Server/tests/integration_flow_round_trip_test.go`
- Outputs:
  - Typed Flow contract definitions and generator/sync outputs in Proto
  - Refactored SubProto consumer code and tests
  - Win editor/store contract consumption updates and tests
  - Server spec + integration test updates

#### Edge Cases
- Unknown node kind or unsupported node spec fields must fail explicitly for new writes.
- `branch` edges with missing or undeclared `edge.case` must still be rejected.
- `foreach` body graphs with missing `result_node_id` target must still be rejected.
- `subflow` self-call or recursion must still be rejected.
- Existing docs still mention old compatibility behavior; those references must be removed or clarified.

#### Acceptance Criteria
- Proto exposes typed Flow graph/node contract definitions for all supported node kinds and reusable enums.
- SubProto no longer owns duplicated wire-level binding/spec structs for Flow graph/node contracts.
- Win imports canonical contract kind/source/op truth from generated or synced artifacts instead of scattered local unions and literal arrays.
- Server `TestIntegrationFlowRoundTrip` covers advanced graph/node round-trip, not only simple compose-node payloads.
- Stable specs and canonical docs no longer claim legacy `local` / `exec` compatibility for new or runtime paths.

#### Risks
- Proto boundary drift: putting too much validation logic into Proto would violate `repos.md`.
- Cross-repo sync drift: generated contract artifacts can go stale if sync rules are unclear.
- Win build/test noise: worktree frontend dependencies or Wails generation may need baseline setup before feature verification.
- Existing control-plane dirtiness in main repo paths means final merge must be surgical and preserve unrelated user changes.

#### Issue List
- None.

### Stage 2 - Architecture Design
#### Overall Solution
- Choose a three-layer closure:
  - Proto becomes the canonical Flow graph/node contract source for wire structs, enums, and generated contract artifacts.
  - SubProto keeps runtime validation/compile/execution logic, but consumes Proto types directly instead of maintaining private duplicates.
  - Win and Server consume canonical generated or mirrored contract outputs and update docs/tests accordingly.
- This closes the contract gap without violating the repo rule that Proto must not absorb runtime handler logic.

#### Alternatives Considered
- Move all graph validation into Proto.
  - Rejected because `repos.md` and `MyFlowHub-Proto/README.md` forbid runtime logic in Proto.
- Keep the current split and just add more tests.
  - Rejected because it leaves the root problem unchanged: contract truth remains duplicated and can drift silently.
- Put the canonical contract only in Server specs.
  - Rejected because code would still not have a typed single source of truth, and Win/SubProto would remain divergent consumers.

#### Module Responsibilities
- `MyFlowHub-Proto`
  - Own `protocol/flow` typed graph/node contract definitions, enums, and generated canonical contract artifacts.
  - Own any generator or sync support needed to expose the canonical contract to humans or TS consumers.
- `MyFlowHub-SubProto`
  - Own Flow runtime decode/validate/compile/execute logic.
  - Consume Proto types directly for graph/node wire contracts.
- `MyFlowHub-Server`
  - Own stable human specs for flow behavior and integration verification.
  - Mirror canonical Proto contract wording where needed, but not invent parallel truth.
- `MyFlowHub-Win`
  - Own editor draft state and UX.
  - Consume canonical contract kinds/enums/spec-shape outputs instead of repo-local hard-coded truth.

#### Data / Call Flow
- Authoring path:
  - Proto typed contract -> generated canonical contract artifact -> synced into Win -> Win draft parse/export -> `flow.set`
- Runtime path:
  - `flow.set` payload -> Proto wire structs -> SubProto decode to Proto node spec types -> SubProto validate/compile -> runtime execute
- Documentation path:
  - Proto canonical contract artifact -> mirrored references in Server `flow.md` and Win editor spec -> workflow archive/change docs

#### Interface Drafts
- Proto additions:
  - typed `NodeKind`, `BindingSourceKind`, `BranchMatchOp`, `InputBinding`, `BindingSource`, and per-node spec structs
  - generated canonical contract output for docs and TS consumption
- SubProto refactor:
  - decode/validate functions accept Proto-defined spec structs
  - graph/node validation remains in SubProto but no longer duplicates wire contract types
- Win consumption:
  - generated/synced `FlowNodeKind` and related enum/type exports replace local handwritten unions/allowed-kind arrays
  - parser/export paths continue to own UI draft adaptation, not wire truth definition

#### Error Handling and Safety
- New-write contract parsing remains strict; unknown or invalid fields must fail early with explicit errors.
- No silent fallback to removed legacy node kinds.
- Sync steps must be explicit and testable so stale generated artifacts are visible.

#### Performance and Testing Strategy
- Avoid repeated JSON schema or contract parsing at runtime where static generated artifacts suffice.
- Proto:
  - `go test ./...`
  - generator check/write verification
- SubProto:
  - `GOWORK=off go test ./flow/... -count=1 -p 1`
- Server:
  - `GOWORK=off go test ./tests -run TestIntegrationFlowRoundTrip -count=1 -p 1`
- Win:
  - targeted frontend flow store/editor tests
  - `GOWORK=off go test ./internal/services/flow ./internal/mcp ./internal/mcpapp -count=1 -p 1`
  - `GOWORK=off wails generate module` if touched bindings or generated imports require it

#### Extensibility Design Points
- Future node kinds should extend Proto first, then regenerate/sync downstream contract artifacts.
- Generated artifacts should capture enums and typed spec shapes, so downstream repos avoid adding new literal unions manually.
- SubProto validation should be organized around Proto types so new node kinds add one decode/validate path rather than duplicating wire structs again.

#### Issue List
- None.

### Stage 3.1 - Planning
#### Project Goal and Current State
- Goal:
  - Finish the last structural Flow completeness gap by unifying graph/node contract truth across Proto, SubProto, Server, and Win.
- Current state:
  - Proto still models `Node.Kind` as `string` and `Node.Spec` as `json.RawMessage`.
  - SubProto already has substantial real contract structs and validation logic, but they live privately in runtime files.
  - Win duplicates node-kind/source/op truth and advanced-node parse/export rules locally.
  - Server stable spec still carries outdated compatibility wording for legacy node kinds.

#### Docs Governance Routing Decision
- 使用 `$m-docs` 校验计划文档路由、requirements/specs 影响和 lessons 查询入口。
- Requirements impact: `none`
  - Existing behavior requirements already live in `repo/MyFlowHub-Server/docs/requirements/flow_data_dag.md`.
  - This workflow changes contract canonicalization and docs truth placement, not product requirements.
- Specs impact: `updated`
  - Update `repo/MyFlowHub-Server/docs/specs/flow.md`
  - Update `repo/MyFlowHub-Win/docs/specs/flow-editor-visual-form.md`
  - Add Proto-side canonical contract doc if generated output is introduced
- Lessons impact: `updated`
  - Revisit `D:/project/MyFlowHub3/docs/lessons/wails-binding-proto-drift.md` if the new generator/sync path changes the recommended quick checks
- Stable truth routing:
  - Wire contract source: `MyFlowHub-Proto/protocol/flow/types.go`
  - Human contract mirror: generated Proto canonical doc plus updated Server/Win specs
  - Workflow execution record: worktree root `plan.md` and later `docs/change/...`

#### Related Requirements / Specs / Lessons
- Related requirements:
  - `D:/project/MyFlowHub3/repo/MyFlowHub-Server/docs/requirements/flow_data_dag.md`
- Related specs:
  - `D:/project/MyFlowHub3/repo/MyFlowHub-Server/docs/specs/flow.md`
  - `D:/project/MyFlowHub3/repo/MyFlowHub-Win/docs/specs/flow-editor-visual-form.md`
  - `D:/project/MyFlowHub3/docs/specs/protocol_map.md`
- Related lessons:
  - `D:/project/MyFlowHub3/docs/lessons/wails-binding-proto-drift.md`

#### Executable Task List
- [x] GC1 Proto canonical Flow graph/node contract and generator artifacts
- [x] GC2 SubProto runtime contract deduplication onto Proto types
- [x] GC3 Win editor/store consumption of canonical contract artifacts
- [x] GC4 Server spec and behavior-level round-trip closure
- [x] GC5 Integration, review, and workflow archive

#### Task Details
##### GC1 - Proto canonical Flow graph/node contract and generator artifacts
- Owner: Main agent
- Worktree: `D:/project/MyFlowHub3/worktrees/proto-graph-contract-closure`
- Plan Path: `D:/project/MyFlowHub3/worktrees/proto-graph-contract-closure/plan.md`
- Goal:
  - Add typed Flow graph/node contract definitions, enums, and canonical generated artifacts in Proto without introducing runtime handler logic.
- Files / Modules:
  - `protocol/flow/types.go`
  - `cmd/flowcontractgen/**` or equivalent generator path
  - `internal/**` generator support as needed
  - `docs/flow_contract.md` or equivalent canonical contract doc
  - `README.md`
- Write Set:
  - Only files under `D:/project/MyFlowHub3/worktrees/proto-graph-contract-closure`
- Acceptance:
  - Proto owns typed spec structs and reusable enums for all supported Flow node kinds.
  - Canonical contract artifact can be regenerated deterministically.
  - No runtime handler logic is added to Proto.
- Test Points:
  - `go test ./...`
  - generator check/write command
- Rollback:
  - Revert new typed contract definitions and generator/docs changes in Proto worktree only.

##### GC2 - SubProto runtime contract deduplication onto Proto types
- Owner: Sub-agent or main agent after GC1 API freeze
- Worktree: `D:/project/MyFlowHub3/worktrees/subproto-graph-contract-closure`
- Plan Path: `D:/project/MyFlowHub3/worktrees/subproto-graph-contract-closure/todo.md`
- Goal:
  - Remove private duplicated Flow wire contract structs and refactor decode/validate paths to use Proto-owned types.
- Files / Modules:
  - `flow/runtime_bindings.go`
  - `flow/handler.go`
  - `flow/types.go`
  - `flow/*_test.go`
  - `flow/go.mod` if local Proto replace is required for the worktree
- Write Set:
  - Only files under `D:/project/MyFlowHub3/worktrees/subproto-graph-contract-closure`
- Acceptance:
  - SubProto no longer owns duplicated Flow node spec/binding wire structs.
  - Existing validation semantics remain intact.
  - New tests cover typed decode and validation paths.
- Test Points:
  - `GOWORK=off go test ./flow/... -count=1 -p 1`
- Rollback:
  - Revert SubProto worktree changes only.

##### GC3 - Win editor/store consumption of canonical contract artifacts
- Owner: Sub-agent or main agent after GC1 generator output is available
- Worktree: `D:/project/MyFlowHub3/worktrees/win-graph-contract-closure`
- Plan Path: `D:/project/MyFlowHub3/worktrees/win-graph-contract-closure/todo.md`
- Goal:
  - Replace scattered local node-kind/source/op truth with canonical generated or synced contract artifacts while preserving existing editor behavior.
- Files / Modules:
  - `frontend/src/stores/flow.ts`
  - `frontend/src/windows/FlowEditorWindow.vue`
  - `frontend/src/components/flow/editor/**`
  - generated contract file path under `frontend/src/**`
  - `frontend/src/stores/flow.test.ts`
  - `docs/specs/flow-editor-visual-form.md`
- Write Set:
  - Only files under `D:/project/MyFlowHub3/worktrees/win-graph-contract-closure`
- Acceptance:
  - Win imports canonical node-kind/source/op truth from generated/synced artifacts.
  - Advanced node parse/export/allowed-kind usage no longer depends on scattered handwritten literal lists.
  - Targeted flow editor/store tests pass.
- Test Points:
  - targeted frontend tests for flow store/editor
  - `GOWORK=off go test ./internal/services/flow ./internal/mcp ./internal/mcpapp -count=1 -p 1`
  - `GOWORK=off wails generate module` if needed
- Rollback:
  - Revert Win worktree changes only.

##### GC4 - Server spec and behavior-level round-trip closure
- Owner: Sub-agent or main agent after GC1 contract/doc output is available
- Worktree: `D:/project/MyFlowHub3/worktrees/server-graph-contract-closure`
- Plan Path: `D:/project/MyFlowHub3/worktrees/server-graph-contract-closure/todo.md`
- Goal:
  - Update Server stable spec to match the new canonical contract and strengthen round-trip integration coverage for advanced graph/node payloads.
- Files / Modules:
  - `docs/specs/flow.md`
  - `tests/integration_flow_round_trip_test.go`
  - mirrored `docs/specs/protocol_map.md` only if Proto generated output requires sync
- Write Set:
  - Only files under `D:/project/MyFlowHub3/worktrees/server-graph-contract-closure`
- Acceptance:
  - No lingering legacy `local` / `exec` compatibility wording remains in stable specs.
  - Integration round-trip test covers advanced graph/node structure and canonical `get` behavior.
- Test Points:
  - `GOWORK=off go test ./tests -run TestIntegrationFlowRoundTrip -count=1 -p 1`
- Rollback:
  - Revert Server worktree changes only.

##### GC5 - Integration, review, and workflow archive
- Owner: Main agent
- Worktree: `D:/project/MyFlowHub3/worktrees/proto-graph-contract-closure`
- Plan Path: `D:/project/MyFlowHub3/worktrees/proto-graph-contract-closure/plan.md`
- Goal:
  - Integrate sub-agent results, resolve conflicts, run cross-repo verification, perform Stage 3.3 review, and prepare Stage 4 archives.
- Files / Modules:
  - worktree-local control docs
  - synced docs/artifacts across the four worktrees
- Write Set:
  - Integration commits inside the four worktrees plus later archive docs during Stage 4
- Acceptance:
  - All task outputs are reconciled.
  - Review checklist passes.
  - Stage 4 archive is ready.
- Test Points:
  - cross-repo targeted test matrix from GC1-GC4
- Rollback:
  - Revert only the specific repo worktree branches that fail acceptance.

#### Dependencies
- GC1 defines the canonical contract surface and unblocks GC2-GC4.
- GC2 depends on the Proto contract API settled by GC1.
- GC3 depends on GC1 generated/synced contract output.
- GC4 depends on GC1 canonical wording/output and may reference GC2 runtime behavior when strengthening integration tests.
- GC5 depends on GC1-GC4 completion.

#### Risks and Notes
- Proto generator scope must stay contract-focused; if it starts encoding runtime validation rules, return to Stage 2 and redesign.
- Win worktree may require `npm install` and `GOWORK=off wails generate module` before trustworthy frontend verification.
- Server and workspace `protocol_map` files are mirrored/generated; do not hand-edit protected generated regions.
- Main repo control paths currently contain unrelated dirty state; integration and closeout must not overwrite them.

#### Parallelism Assessment
- Assessment:
  - GC1 is the critical-path seed task and should stay with the main agent first.
  - After GC1 contract surface and generator output stabilize, GC2/GC3/GC4 have clean repo-level write-set separation and can proceed in parallel.
  - GC5 remains main-agent-only for integration and acceptance.
- Planned delegation in Stage `3.2`:
  - Worker A: GC2 in `subproto-graph-contract-closure`
  - Worker B: GC3 in `win-graph-contract-closure`
  - Worker C: GC4 in `server-graph-contract-closure`
- Required context packets:
  - must include the task ID, repo, branch, base, worktree path, governing plan path, write set, acceptance, tests, rollback, and contract references from Stage 1/2.

#### Issue List
- None.

阻塞：否
进入 3.2

### Stage 3.2 - Implementation
#### Task Results
- GC1
  - `MyFlowHub-Proto` now owns canonical typed Flow graph/node contract definitions in `protocol/flow/types.go`.
  - Added `cmd/flowcontractgen` plus `internal/flowcontract/**` to generate `docs/flow_contract.md` and `generated/flow_contract.ts`.
  - Removed plan-external legacy compatibility by keeping canonical node/source/op truth limited to current supported contract surface.
- GC2
  - `MyFlowHub-SubProto/flow` no longer treats private wire structs as the source of truth for Flow graph/node contracts.
  - Runtime decode/validate/materialize now consumes Proto canonical types directly and rejects legacy `call.args` / `local` / `exec` payloads by design.
  - Added a local Proto `replace` in `flow/go.mod` for this workflow worktree to freeze against the GC1 API.
- GC3
  - `MyFlowHub-Win` now consumes `frontend/src/generated/flow_contract.ts` as the canonical source for `FlowNodeKind`, `FlowBindingSourceKind`, and `FlowBranchMatchOp`.
  - `frontend/src/stores/flow.ts` now owns only Win-local draft projection, root/body source-kind filtering, label mapping, and strict parse/export validation.
  - Inspector and binding-dialog option lists now derive from the canonical generated artifact instead of handwritten literal arrays.
- GC4
  - `MyFlowHub-Server/docs/specs/flow.md` now points back to Proto canonical contract sources and removes lingering legacy compatibility wording.
  - `TestIntegrationFlowRoundTrip` now covers advanced graph/node round-trip for `set_var`, `transform`, `branch`, `foreach`, `subflow`, `edge.case`, and canonical `get` payload equality.
- GC5
  - Main agent revalidated the four-repo closure, confirmed the Win synced TS artifact matches Proto generated output, and prepared Stage 4 archive artifacts.

#### Verification Matrix
- Proto
  - `GOWORK=off go test ./...`
  - Result: passed
- SubProto
  - `GOWORK=off go test ./... -count=1 -p 1`
  - Result: passed
- Server
  - `GOWORK=off go test ./tests -run TestIntegrationFlowRoundTrip -count=1 -p 1`
  - Result: passed
- Win
  - `GOWORK=off wails generate module`
  - Result: passed
  - `GOWORK=off go test ./internal/services/flow ./internal/mcp ./internal/mcpapp -count=1 -p 1`
  - Result: passed
  - `npx vitest run src/stores/flow.test.ts --maxWorkers=1 --no-file-parallelism`
  - Result: passed
  - `npx vitest run src/components/flow/editor/FlowNodeInspector.test.ts --maxWorkers=1 --no-file-parallelism`
  - Result: passed
  - `npx vitest run src/components/flow/editor/FlowBodyNodeInspector.test.ts --maxWorkers=1 --no-file-parallelism`
  - Result: passed
  - `npx vitest run src/components/flow/editor/FlowFieldBindingDialog.test.ts --pool threads --maxWorkers=1 --no-file-parallelism`
  - Result: passed
  - `npx vitest run src/windows/FlowEditorWindow.test.ts --pool threads --maxWorkers=1 --no-file-parallelism`
  - Result: passed
- Cross-repo sync
  - Compared `generated/flow_contract.ts` in Proto and Win synced copy.
  - Result: identical

### Stage 3.3 - Code Review
- 需求覆盖：通过
  - GC1-GC4 now close the last remaining `flow graph/node contract` gap and keep legacy compatibility intentionally removed.
- 架构合理性：通过
  - Proto owns canonical contract types and generated artifacts; SubProto keeps runtime semantics; Win and Server consume or mirror Proto truth instead of redefining it.
- 性能风险（N+1 / 重复计算 / 多余 I/O / 锁竞争）：通过
  - Changes are mostly static contract definitions, decode-path type substitutions, generated artifact imports, and test/spec updates with no new repeated runtime I/O.
- 可读性与一致性：通过
  - Shared enums and option lists now come from one generated artifact; duplicated literal unions/arrays were removed from Win and duplicated wire structs were removed from SubProto.
- 可扩展性与配置化：通过
  - Future node/source/op additions now have one required first step in Proto, then downstream sync/adapter work instead of silent parallel edits.
- 稳定性与安全：通过
  - Unknown or removed legacy payload shapes now fail explicitly instead of being silently accepted.
- 测试覆盖情况：通过
  - Proto/SubProto/Server targeted Go tests pass; Win store/editor targeted vitest cases and Go-side tests pass under stable single-worker settings.
- 子Agent治理与审计（任务映射、上下文完整性、文件所有权、结果复核、冲突处理、记录完整性）：通过
  - GC2, GC3, and GC4 were dispatched with bounded repo-level write sets and governing plan references.
  - Main agent independently revalidated GC2-GC4 outputs and handled final integration acceptance.
  - GC3 sub-agent completion arrived late, but its resulting write set matched the already reviewed local state and was revalidated by the main agent before closure.

阻塞：否
进入 4
