package protocolmap

// Context: This file supports the Proto source-of-truth workflow around render.

import (
	"fmt"
	"sort"
	"strings"
)

func RenderGeneratedMarkdown(pm *ProtocolMap) (string, error) {
	if pm == nil {
		return "", fmt.Errorf("nil protocol map")
	}
	var b strings.Builder

	b.WriteString("## SubProto Overview\n")
	b.WriteString("\n")
	b.WriteString("| SubProto | Name | Package |\n")
	b.WriteString("|---:|---|---|\n")
	for _, p := range pm.Protocols {
		b.WriteString(fmt.Sprintf("| %d | %s | `protocol/%s` |\n", p.SubProtoID, p.DisplayName, p.Dir))
	}
	b.WriteString("\n")

	for i := range pm.Protocols {
		p := pm.Protocols[i]
		b.WriteString(fmt.Sprintf("## %s (SubProto=%d)\n", p.DisplayName, p.SubProtoID))
		b.WriteString("\n")

		b.WriteString("**Actions**\n")
		if len(p.Actions) == 0 {
			b.WriteString("- (none)\n")
		} else {
			// already sorted
			for _, a := range p.Actions {
				b.WriteString(fmt.Sprintf("- `%s = %q`\n", a.Name, a.Str))
			}
		}
		b.WriteString("\n")

		b.WriteString("**Payload types**\n")
		if len(p.Types) == 0 {
			b.WriteString("- (none)\n")
		} else {
			for _, t := range p.Types {
				b.WriteString(fmt.Sprintf("- `%s`\n", t))
			}
		}
		b.WriteString("\n")

		other := formatOtherConsts(p.Consts)
		if other != "" {
			b.WriteString("**Other constants**\n")
			b.WriteString(other)
			b.WriteString("\n")
		}
	}

	return b.String(), nil
}

func formatOtherConsts(consts []Const) string {
	if len(consts) == 0 {
		return ""
	}
	// stable: sort by name; caller already did, but keep defensive.
	cc := append([]Const(nil), consts...)
	sort.Slice(cc, func(i, j int) bool { return cc[i].Name < cc[j].Name })

	var b strings.Builder
	for _, c := range cc {
		switch c.Kind {
		case ConstKindString:
			b.WriteString(fmt.Sprintf("- `%s = %q`\n", c.Name, c.Str))
		case ConstKindInt:
			// keep raw literal if available (hex/dec), fallback to int.
			raw := strings.TrimSpace(c.Raw)
			if raw == "" {
				raw = fmt.Sprintf("%d", c.Int)
			}
			b.WriteString(fmt.Sprintf("- `%s = %s`\n", c.Name, raw))
		default:
			// ignore unknown
		}
	}
	return b.String()
}
