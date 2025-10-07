package logger

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type ginLikeFormatter struct{}

func (f *ginLikeFormatter) Format(e *logrus.Entry) ([]byte, error) {
	var b bytes.Buffer
	ts := e.Time.Format("2006/01/02 - 15:04:05")
	level := strings.ToUpper(e.Level.String())

	// [APP] 2025/08/14 - 03:37:59 | INFO | message | k1=v1 k2=v2
	b.WriteString("[APP] ")
	b.WriteString(ts)
	b.WriteString(" | ")
	b.WriteString(level)
	b.WriteString(" | ")
	b.WriteString(e.Message)

	if len(e.Data) > 0 {
		keys := make([]string, 0, len(e.Data))
		for k := range e.Data {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		b.WriteString(" | ")
		for i, k := range keys {
			if i > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(k)
			b.WriteByte('=')
			b.WriteString(toString(e.Data[k]))
		}
	}
	b.WriteByte('\n')
	return b.Bytes(), nil
}

func toString(v any) string {
	if t, ok := v.(time.Time); ok {
		return t.Format(time.RFC3339)
	}
	return fmt.Sprint(v)
}
