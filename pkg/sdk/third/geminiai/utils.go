package geminiai

import (
	"fmt"
	"google.golang.org/genai"
	"strings"
)

func ResponseToString(resp *genai.GenerateContentResponse) string {
	var b strings.Builder
	for i, x := range resp.Candidates {
		if len(resp.Candidates) > 1 {
			_, _ = fmt.Fprintf(&b, "%d:", i+1)
		}
		b.WriteString(ContentToString(x.Content))
	}
	return b.String()
}

func ContentToString(c *genai.Content) string {
	var b strings.Builder
	if c == nil || c.Parts == nil {
		return ""
	}
	for i, part := range c.Parts {
		if i > 0 {
			_, _ = fmt.Fprintf(&b, ";")
		}
		_, _ = fmt.Fprintf(&b, "%v", part)
	}
	return b.String()
}
