package flowcontract

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	protocolflow "github.com/yttydcs/myflowhub-proto/protocol/flow"
)

var rawMessageType = reflect.TypeOf(json.RawMessage{})

type Contract struct {
	Enums     []EnumDoc
	KindSpecs []KindSpec
	Types     []TypeDoc
}

type EnumDoc struct {
	Name      string
	TSName    string
	ConstName string
	Values    []string
}

type KindSpec struct {
	Kind   string
	GoName string
	TSName string
	Type   reflect.Type
}

type TypeDoc struct {
	GoName string
	TSName string
	Note   string
	Fields []FieldDoc
}

type FieldDoc struct {
	GoName   string
	JSONName string
	GoType   string
	TSType   string
	Optional bool
	Note     string
}

func BuildContract() (*Contract, error) {
	contract := &Contract{
		Enums: []EnumDoc{
			{
				Name:      "NodeKind",
				TSName:    "FlowNodeKind",
				ConstName: "FLOW_NODE_KINDS",
				Values:    toStrings(protocolflow.SupportedNodeKinds),
			},
			{
				Name:      "BindingSourceKind",
				TSName:    "FlowBindingSourceKind",
				ConstName: "FLOW_BINDING_SOURCE_KINDS",
				Values:    toStrings(protocolflow.SupportedBindingSourceKinds),
			},
			{
				Name:      "BranchMatchOp",
				TSName:    "FlowBranchMatchOp",
				ConstName: "FLOW_BRANCH_MATCH_OPS",
				Values:    toStrings(protocolflow.SupportedBranchMatchOps),
			},
		},
		KindSpecs: []KindSpec{
			{Kind: string(protocolflow.NodeKindCall), GoName: "CallSpec", TSName: "FlowCallSpec", Type: reflect.TypeOf(protocolflow.CallSpec{})},
			{Kind: string(protocolflow.NodeKindCompose), GoName: "ComposeSpec", TSName: "FlowComposeSpec", Type: reflect.TypeOf(protocolflow.ComposeSpec{})},
			{Kind: string(protocolflow.NodeKindTransform), GoName: "TransformSpec", TSName: "FlowTransformSpec", Type: reflect.TypeOf(protocolflow.TransformSpec{})},
			{Kind: string(protocolflow.NodeKindSetVar), GoName: "SetVarSpec", TSName: "FlowSetVarSpec", Type: reflect.TypeOf(protocolflow.SetVarSpec{})},
			{Kind: string(protocolflow.NodeKindBranch), GoName: "BranchSpec", TSName: "FlowBranchSpec", Type: reflect.TypeOf(protocolflow.BranchSpec{})},
			{Kind: string(protocolflow.NodeKindForeach), GoName: "ForeachSpec", TSName: "FlowForeachSpec", Type: reflect.TypeOf(protocolflow.ForeachSpec{})},
			{Kind: string(protocolflow.NodeKindSubflow), GoName: "SubflowSpec", TSName: "FlowSubflowSpec", Type: reflect.TypeOf(protocolflow.SubflowSpec{})},
		},
	}

	typeOrder := []reflect.Type{
		reflect.TypeOf(protocolflow.Trigger{}),
		reflect.TypeOf(protocolflow.NodeUILayout{}),
		reflect.TypeOf(protocolflow.Edge{}),
		reflect.TypeOf(protocolflow.BindingSource{}),
		reflect.TypeOf(protocolflow.InputBinding{}),
		reflect.TypeOf(protocolflow.ComposeSpec{}),
		reflect.TypeOf(protocolflow.SetVarSpec{}),
		reflect.TypeOf(protocolflow.TransformExpr{}),
		reflect.TypeOf(protocolflow.TransformSpec{}),
		reflect.TypeOf(protocolflow.BranchMatch{}),
		reflect.TypeOf(protocolflow.BranchCase{}),
		reflect.TypeOf(protocolflow.BranchSpec{}),
		reflect.TypeOf(protocolflow.ForeachSpec{}),
		reflect.TypeOf(protocolflow.SubflowSpec{}),
		reflect.TypeOf(protocolflow.CallSpec{}),
		reflect.TypeOf(protocolflow.Graph{}),
	}

	for _, typ := range typeOrder {
		doc, err := buildTypeDoc(typ)
		if err != nil {
			return nil, err
		}
		contract.Types = append(contract.Types, doc)
	}

	return contract, nil
}

func buildTypeDoc(typ reflect.Type) (TypeDoc, error) {
	if typ.Kind() != reflect.Struct {
		return TypeDoc{}, fmt.Errorf("type %s is not a struct", typ.String())
	}

	doc := TypeDoc{
		GoName: typ.Name(),
		TSName: tsTypeName(typ),
		Note:   typeNotes[typ.Name()],
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}
		jsonName, optional, skip := parseJSONTag(field)
		if skip {
			continue
		}
		doc.Fields = append(doc.Fields, FieldDoc{
			GoName:   field.Name,
			JSONName: jsonName,
			GoType:   goTypeString(field.Type),
			TSType:   tsFieldType(typ.Name(), field),
			Optional: optional,
			Note:     fieldNotes[fieldKey(typ.Name(), field.Name)],
		})
	}

	return doc, nil
}

func parseJSONTag(field reflect.StructField) (name string, optional bool, skip bool) {
	tag := field.Tag.Get("json")
	if tag == "-" {
		return "", false, true
	}
	if tag == "" {
		return field.Name, false, false
	}
	parts := strings.Split(tag, ",")
	name = parts[0]
	if name == "" {
		name = field.Name
	}
	for _, part := range parts[1:] {
		if part == "omitempty" {
			optional = true
		}
	}
	if field.Type.Kind() == reflect.Pointer {
		optional = true
	}
	return name, optional, false
}

func tsFieldType(typeName string, field reflect.StructField) string {
	if override, ok := fieldTypeOverrides[fieldKey(typeName, field.Name)]; ok {
		return override
	}
	return tsTypeString(field.Type)
}

func tsTypeString(typ reflect.Type) string {
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	switch {
	case typ == rawMessageType:
		return "JSONValue"
	case typ == reflect.TypeOf(protocolflow.NodeKind("")):
		return "FlowNodeKind"
	case typ == reflect.TypeOf(protocolflow.BindingSourceKind("")):
		return "FlowBindingSourceKind"
	case typ == reflect.TypeOf(protocolflow.BranchMatchOp("")):
		return "FlowBranchMatchOp"
	}

	switch typ.Kind() {
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Slice:
		return tsTypeString(typ.Elem()) + "[]"
	case reflect.Map:
		return fmt.Sprintf("Record<string, %s>", tsTypeString(typ.Elem()))
	case reflect.Struct:
		return tsTypeName(typ)
	default:
		return "unknown"
	}
}

func tsTypeName(typ reflect.Type) string {
	switch typ {
	case reflect.TypeOf(protocolflow.Trigger{}):
		return "FlowTrigger"
	case reflect.TypeOf(protocolflow.NodeUILayout{}):
		return "FlowNodeUILayout"
	case reflect.TypeOf(protocolflow.Edge{}):
		return "FlowEdge"
	case reflect.TypeOf(protocolflow.BindingSource{}):
		return "FlowBindingSource"
	case reflect.TypeOf(protocolflow.InputBinding{}):
		return "FlowInputBinding"
	case reflect.TypeOf(protocolflow.CallSpec{}):
		return "FlowCallSpec"
	case reflect.TypeOf(protocolflow.ComposeSpec{}):
		return "FlowComposeSpec"
	case reflect.TypeOf(protocolflow.SetVarSpec{}):
		return "FlowSetVarSpec"
	case reflect.TypeOf(protocolflow.TransformExpr{}):
		return "FlowTransformExpr"
	case reflect.TypeOf(protocolflow.TransformSpec{}):
		return "FlowTransformSpec"
	case reflect.TypeOf(protocolflow.BranchMatch{}):
		return "FlowBranchMatch"
	case reflect.TypeOf(protocolflow.BranchCase{}):
		return "FlowBranchCase"
	case reflect.TypeOf(protocolflow.BranchSpec{}):
		return "FlowBranchSpec"
	case reflect.TypeOf(protocolflow.ForeachSpec{}):
		return "FlowForeachSpec"
	case reflect.TypeOf(protocolflow.SubflowSpec{}):
		return "FlowSubflowSpec"
	case reflect.TypeOf(protocolflow.Graph{}):
		return "FlowGraph"
	case reflect.TypeOf(protocolflow.Node{}):
		return "FlowNode"
	default:
		if typ.Name() == "" {
			return "unknown"
		}
		return typ.Name()
	}
}

func goTypeString(typ reflect.Type) string {
	switch {
	case typ == rawMessageType:
		return "json.RawMessage"
	case typ == reflect.TypeOf(protocolflow.NodeKind("")):
		return "NodeKind"
	case typ == reflect.TypeOf(protocolflow.BindingSourceKind("")):
		return "BindingSourceKind"
	case typ == reflect.TypeOf(protocolflow.BranchMatchOp("")):
		return "BranchMatchOp"
	}

	switch typ.Kind() {
	case reflect.Pointer:
		return "*" + goTypeString(typ.Elem())
	case reflect.Slice:
		return "[]" + goTypeString(typ.Elem())
	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", goTypeString(typ.Key()), goTypeString(typ.Elem()))
	default:
		if typ.Name() != "" {
			return typ.Name()
		}
		return typ.String()
	}
}

func fieldKey(typeName, fieldName string) string {
	return typeName + "." + fieldName
}

func toStrings[T ~string](values []T) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		out = append(out, string(value))
	}
	return out
}

var typeNotes = map[string]string{
	"Graph":       "The top-level graph uses kind-tagged nodes and branch-aware edges.",
	"CallSpec":    "New-write call nodes use args_template plus inputs. Legacy args are intentionally omitted.",
	"ForeachSpec": "The body graph is isolated per item and the result is collected from result_node_id.",
	"SubflowSpec": "Subflow nodes synchronously execute another flow on the same executor.",
}

var fieldNotes = map[string]string{
	"Graph.Nodes":               "Each node uses kind-tagged specs.",
	"Node.Spec":                 "Decode with the spec type mapped from kind.",
	"Edge.Case":                 "Only used on outgoing branch edges.",
	"CallSpec.ArgsTemplate":     "Must be a JSON object in the stable contract.",
	"ComposeSpec.Template":      "Must be a JSON object in the stable contract.",
	"SetVarSpec.Template":       "May be any JSON value.",
	"TransformExpr.Required":    "Only applies when source is selected.",
	"BranchMatch.Value":         "Omit only when op is exists.",
	"ForeachSpec.Body":          "Must be a valid child graph.",
	"ForeachSpec.ResultNodeID":  "Must exist inside the foreach body graph.",
	"SubflowSpec.InputTemplate": "Must be a JSON object in the stable contract.",
}

var fieldTypeOverrides = map[string]string{
	"Node.Spec":                 "FlowNodeSpec",
	"CallSpec.ArgsTemplate":     "JSONObject",
	"ComposeSpec.Template":      "JSONObject",
	"SetVarSpec.Template":       "JSONValue",
	"BranchMatch.Value":         "JSONValue",
	"SubflowSpec.InputTemplate": "JSONObject",
}
