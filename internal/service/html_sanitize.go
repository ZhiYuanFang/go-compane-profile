package service

import (
	"regexp"
	"strings"
)

var (
	reScriptStyle = regexp.MustCompile(`(?is)<(script|style)[^>]*>.*?</(script|style)>`)
	reTags        = regexp.MustCompile(`(?is)</?([a-z0-9]+)([^>]*)>`)
	reOnAttr      = regexp.MustCompile(`(?i)\s+on[a-z]+\s*=\s*("[^"]*"|'[^']*'|[^\s>]+)`)
	reHrefJS      = regexp.MustCompile(`(?i)\s+(href|src)\s*=\s*("\s*javascript:[^"]*"|'\s*javascript:[^']*')`)
	reStyleAttr   = regexp.MustCompile(`(?i)\s+style\s*=\s*("([^"]*)"|'([^']*)')`)
)

var allowedActivityTags = map[string]bool{
	"p": true, "br": true, "span": true, "strong": true, "b": true, "div": true,
}

// SanitizeActivityHTML keeps a small HTML subset for font size, bold, and color.
func SanitizeActivityHTML(html string) string {
	s := strings.TrimSpace(html)
	if s == "" {
		return ""
	}
	s = reScriptStyle.ReplaceAllString(s, "")
	s = reOnAttr.ReplaceAllString(s, "")
	s = reHrefJS.ReplaceAllString(s, "")
	s = reTags.ReplaceAllStringFunc(s, func(tag string) string {
		m := reTags.FindStringSubmatch(tag)
		if len(m) < 2 {
			return ""
		}
		name := strings.ToLower(m[1])
		if !allowedActivityTags[name] {
			return ""
		}
		closing := strings.HasPrefix(strings.TrimSpace(tag), "</")
		if closing {
			return "</" + name + ">"
		}
		if name == "br" {
			return "<br/>"
		}
		attrs := m[2]
		style := ""
		if sm := reStyleAttr.FindStringSubmatch(attrs); len(sm) > 0 {
			raw := sm[2]
			if raw == "" {
				raw = sm[3]
			}
			style = filterAllowedStyles(raw)
		}
		if style != "" {
			return "<" + name + ` style="` + style + `">`
		}
		return "<" + name + ">"
	})
	return s
}

func filterAllowedStyles(raw string) string {
	parts := strings.Split(raw, ";")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		kv := strings.SplitN(p, ":", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(kv[0]))
		val := strings.TrimSpace(kv[1])
		if val == "" || strings.ContainsAny(val, "<>\"'") {
			continue
		}
		switch key {
		case "color", "font-size", "font-weight":
			out = append(out, key+": "+val)
		}
	}
	return strings.Join(out, "; ")
}
