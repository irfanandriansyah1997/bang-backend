package logger

import (
	"bytes"
	"encoding/json"
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
		return p.w.Write(line)
	}

	out.WriteByte('\n')
	if _, err := p.w.Write(out.Bytes()); err != nil {
		return 0, err
	}

	return len(line), nil
}
