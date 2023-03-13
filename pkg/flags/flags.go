// Package flags is a helper module that can be used to provide common functionality for
// flag sets
package flags

import (
	"bytes"
	"flag"
	"fmt"
	"strings"

	"github.com/kr/text"
)

// Usage takes a help string and a flag set and appends the flags
// and their descriptions below the help string and returns it.
// This is to standardize the usage function across the different environments
func Usage(help string, flags *flag.FlagSet) string {
	out := new(bytes.Buffer)
	out.WriteString(strings.TrimSpace(help))
	out.WriteString("\n")
	out.WriteString("\n")

	if flags != nil {
		flags.VisitAll(func(f *flag.Flag) {
			example, _ := flag.UnquoteUsage(f)
			if example != "" {
				fmt.Fprintf(out, "  -%s=<%s>\n", f.Name, example)
			} else {
				fmt.Fprintf(out, "  -%s\n", f.Name)
			}

			indented := wrapAtLength(f.Usage, 5)
			fmt.Fprintf(out, "%s\n\n", indented)
		})
	}
	return strings.TrimRight(out.String(), "\n")
}

// maxLineLength is the maximum width of any line.
const maxLineLength int = 72

// wrapAtLength wraps the given text at the maxLineLength, taking into account
// any provided left padding.
func wrapAtLength(s string, pad int) string {
	wrapped := text.Wrap(s, maxLineLength-pad)
	lines := strings.Split(wrapped, "\n")
	for i, line := range lines {
		lines[i] = strings.Repeat(" ", pad) + line
	}
	return strings.Join(lines, "\n")
}
