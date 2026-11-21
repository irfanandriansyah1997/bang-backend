package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type prettyWriter struct {
	w io.Writer
}

func NewPrettyWriter(w io.Writer) io.Writer {
	return prettyWriter{w: w}
}

func (p prettyWriter) Write(line []byte) (int, error) {
	var out bytes.Buffer
	if err := json.Indent(&out, bytes.TrimSpace(line), "", "  "); err != nil {
		out.Reset()
		out.Write(line) // not JSON: write the line as it is
	} else {
		out.WriteByte('\n')
	}

	if _, err := p.w.Write(out.Bytes()); err != nil {
		return 0, fmt.Errorf("logger: pretty write: %w", err)
	}

	return len(line), nil
}
