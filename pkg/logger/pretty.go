// Package logger provides a structured logging interface for applications.
package logger

import (
	"bytes"
	"encoding/json"
	"os"
	"sort"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

const (
	ansiReset = "\033[0m"
	ansiBold  = "\033[1m"
	ansiFaint = "\033[2m"

	callerColor   = "\033[38;2;148;163;184m"
	metaKeyColor  = "\033[38;2;94;234;212m"
	metaValColor  = "\033[38;2;226;232;240m"
	textColor     = "\033[38;2;226;232;240m"
	warnKeyColor  = "\033[38;2;251;191;36m"
	warnValColor  = "\033[38;2;253;230;138m"
	errorKeyColor = "\033[38;2;248;113;113m"
	errorValColor = "\033[38;2;254;202;202m"
)

var levelPalette = map[zapcore.Level]string{
	zapcore.DebugLevel:  "\033[38;2;129;140;248m",
	zapcore.InfoLevel:   "\033[38;2;16;185;129m",
	zapcore.WarnLevel:   "\033[38;2;245;158;11m",
	zapcore.ErrorLevel:  "\033[38;2;248;113;113m",
	zapcore.DPanicLevel: "\033[38;2;244;63;94m",
	zapcore.PanicLevel:  "\033[38;2;244;63;94m",
	zapcore.FatalLevel:  "\033[38;2;217;70;239m",
}

var levelEmoji = map[zapcore.Level]string{
	zapcore.DebugLevel:  "🧪",
	zapcore.InfoLevel:   "ℹ️ ", // added space for alignment
	zapcore.WarnLevel:   "⚠️ ", // added space for alignment
	zapcore.ErrorLevel:  "🚨",
	zapcore.DPanicLevel: "🚨",
	zapcore.PanicLevel:  "🚨",
	zapcore.FatalLevel:  "💥",
}

// prettyLogger wraps zap's JSON encoder to produce colorized, indented output suited for terminals.
type prettyLogger struct {
	zapcore.Encoder
}

// Clone ensures derived loggers keep the pretty encoder wrapper.
func (e *prettyLogger) Clone() zapcore.Encoder {
	return &prettyLogger{Encoder: e.Encoder.Clone()}
}

// newPrettyLogger creates a new pretty logger with color support and JSON indentation.
func newPrettyLogger(cfg *zap.Config) *zap.Logger {
	enc := &prettyLogger{Encoder: zapcore.NewJSONEncoder(cfg.EncoderConfig)}
	core := zapcore.NewCore(enc, zapcore.AddSync(os.Stdout), cfg.Level)
	opts := buildPrettyOptions(cfg)
	return zap.New(core, opts...)
}

func buildPrettyOptions(cfg *zap.Config) []zap.Option {
	opts := []zap.Option{zap.ErrorOutput(zapcore.AddSync(os.Stderr))}
	if cfg.Development {
		opts = append(opts, zap.Development())
	}
	if !cfg.DisableCaller {
		opts = append(opts, zap.AddCaller())
	}
	if len(cfg.InitialFields) > 0 {
		keys := make([]string, 0, len(cfg.InitialFields))
		for k := range cfg.InitialFields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		fields := make([]zap.Field, 0, len(keys))
		for _, k := range keys {
			fields = append(fields, zap.Any(k, cfg.InitialFields[k]))
		}
		opts = append(opts, zap.Fields(fields...))
	}
	return opts
}

// EncodeEntry formats a log entry with pretty printing and colorization.
func (e *prettyLogger) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	jsonBuf, err := e.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, err
	}

	raw := append([]byte(nil), jsonBuf.Bytes()...)
	jsonBuf.Reset()

	trimmed := bytes.TrimSpace(raw)
	var payload map[string]any
	if err := json.Unmarshal(trimmed, &payload); err != nil {
		jsonBuf.Write(raw)
		if len(raw) == 0 || raw[len(raw)-1] != '\n' {
			jsonBuf.WriteByte('\n')
		}
		return jsonBuf, nil
	}

	jsonBuf.WriteString(buildHeader(entry, payload))
	meta := filterReserved(payload)
	writeMetadata(jsonBuf, meta, entry.Level)

	return jsonBuf, nil
}

func buildHeader(entry zapcore.Entry, payload map[string]any) string {
	timestamp := headerTimestamp(entry)
	level := headerLevel(entry, payload)
	message := headerMessage(entry, payload)
	callerInfo := resolveCaller(entry, payload)
	emoji := levelEmoji[entry.Level]

	var b strings.Builder
	b.WriteString(styleTime("[" + timestamp + "]"))
	b.WriteByte(' ')
	if emoji != "" {
		b.WriteString(emoji)
		b.WriteByte(' ')
	}
	b.WriteString(styleLevel(level, entry.Level))
	if message != "" {
		b.WriteByte(' ')
		b.WriteString(styleMessage(entry.Level, message))
	}
	if callerInfo != "" {
		b.WriteByte(' ')
		b.WriteString(styleCaller("(" + callerInfo + ")"))
	}
	b.WriteByte('\n')
	return b.String()
}

func headerTimestamp(entry zapcore.Entry) string {
	timestamp := entry.Time
	if timestamp.IsZero() {
		timestamp = time.Now()
	}
	value := timestamp.Format(time.DateTime)
	return value
}

func headerLevel(entry zapcore.Entry, payload map[string]any) string {
	value := strings.ToUpper(entry.Level.String())
	if lvlVal, ok := payload[levelKey]; ok {
		if lvlText, ok := lvlVal.(string); ok && lvlText != "" {
			value = strings.ToUpper(lvlText)
		}
	}
	return value
}

func headerMessage(entry zapcore.Entry, payload map[string]any) string {
	value := entry.Message
	if msgVal, ok := payload[messageKey]; ok {
		if msgText, ok := msgVal.(string); ok {
			value = msgText
		}
	}
	return value
}

func resolveCaller(entry zapcore.Entry, payload map[string]any) string {
	if callerVal, ok := payload[callerKey]; ok {
		if callerText, ok := callerVal.(string); ok && callerText != "" {
			return callerText
		}
	}
	if entry.Caller.Defined {
		return entry.Caller.TrimmedPath()
	}
	return ""
}

func filterReserved(payload map[string]any) map[string]any {
	meta := make(map[string]any)
	for k, v := range payload {
		switch k {
		case timeKey, levelKey, messageKey, callerKey:
			continue
		default:
			meta[k] = v
		}
	}
	if _, ok := meta[nameKey]; !ok {
		if nameVal, has := payload[nameKey]; has {
			meta[nameKey] = nameVal
		}
	}
	return meta
}

func writeMetadata(buf *buffer.Buffer, meta map[string]any, level zapcore.Level) {
	if len(meta) == 0 {
		return
	}

	keyColor, valColor := metaColours(level)
	pretty, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		buf.WriteString(ansiFaint + valColor + metaFallback(meta) + ansiReset)
		buf.WriteByte('\n')
		return
	}

	lines := bytes.Split(pretty, []byte("\n"))
	for i, line := range lines {
		formatted := styleMetaLine(line, keyColor, valColor)
		if formatted == "" {
			continue
		}
		buf.WriteString(formatted)
		if i < len(lines)-1 {
			buf.WriteByte('\n')
		}
	}
	buf.WriteByte('\n')
}

func styleLevel(level string, lvl zapcore.Level) string {
	color := levelPalette[lvl]
	if color == "" {
		color = levelPalette[zapcore.InfoLevel]
	}
	return ansiBold + color + level + ansiReset
}

func styleTime(v string) string {
	return ansiFaint + callerColor + v + ansiReset
}

func styleCaller(v string) string {
	return ansiFaint + callerColor + v + ansiReset
}

func styleMessage(level zapcore.Level, v string) string {
	if v == "" {
		return ""
	}
	colour := messageColour(level)
	return colour + v + ansiReset
}

func messageColour(level zapcore.Level) string {
	switch level {
	case zapcore.WarnLevel:
		return warnValColor
	case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		return errorValColor
	default:
		return textColor
	}
}

func metaFallback(meta map[string]any) string {
	raw, err := json.Marshal(meta)
	if err != nil {
		return "{}"
	}
	return string(raw)
}

func metaColours(level zapcore.Level) (string, string) {
	switch level {
	case zapcore.WarnLevel:
		return warnKeyColor, warnValColor
	case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		return errorKeyColor, errorValColor
	default:
		return metaKeyColor, metaValColor
	}
}

func styleMetaLine(line []byte, keyColor, valColor string) string {
	if len(line) == 0 {
		return ""
	}
	trimmed := bytes.TrimSpace(line)
	if len(trimmed) == 0 {
		return ""
	}
	indentLen := len(line) - len(bytes.TrimLeft(line, " "))
	indent := string(line[:indentLen])
	colonIdx := bytes.IndexByte(trimmed, ':')
	if colonIdx == -1 {
		return indent + ansiFaint + valColor + string(trimmed) + ansiReset
	}
	key := string(trimmed[:colonIdx])
	rest := string(trimmed[colonIdx+1:])
	return indent + keyColor + key + ansiReset + ":" + ansiFaint + valColor + rest + ansiReset
}
