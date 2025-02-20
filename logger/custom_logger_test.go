package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ibm/opentalaria/utils"
)

func Test_painter(t *testing.T) {
	type args struct {
		colorCode int
		msg       string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Happy flow retriving INFO", args: args{colorCode: 32, msg: "dummy message"}, want: "[32mdummymessage[0m"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := painter(tt.args.colorCode, tt.args.msg)
			got = utils.TrimWhitespaces(got)
			if got != tt.want {
				t.Errorf("painter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_colorLogLevel(t *testing.T) {
	const (
		info  = "INFO"
		debug = "DEBUG"
		err   = "ERROR"
		warn  = "WARN"
	)
	type args struct {
		level string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 int
	}{
		{name: "Happy flow retriving INFO", args: args{level: info}, want: "[32mINFO[0m", want1: 32},
		{name: "Happy flow retriving DEBUG", args: args{level: debug}, want: "[97mDEBUG[0m", want1: 37},
		{name: "Happy flow retriving INFO", args: args{level: err}, want: "[91mERROR[0m", want1: 91},
		{name: "Happy flow retriving INFO", args: args{level: warn}, want: "[33mWARN[0m", want1: 33},
		{name: "Happy flow retriving default case", args: args{level: "default"}, want: "[97mINFO[0m", want1: 37},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := colorLogLevel(tt.args.level)
			if got != tt.want {
				t.Errorf("colorLogLevel() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("colorLogLevel() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNewCustomeHandler(t *testing.T) {
	var buf bytes.Buffer
	type args struct {
		opts *Options
	}
	tests := []struct {
		name    string
		args    args
		want    *CustomHandler
		wantOut string
	}{
		{name: "Happy flow retriving INFO", want: NewCustomHandler(&buf, nil)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			if got := NewCustomHandler(out, tt.args.opts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCustomeHandler() = %v, want %v", got, tt.want)
			}
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("NewCustomeHandler() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func TestCustomHandler_Handle(t *testing.T) {
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
	var b bytes.Buffer

	recWithAttr := slog.NewRecord(time.Now(), slog.LevelInfo, "some message", pcs[0])
	recWithAttr.AddAttrs(slog.Attr{Key: "key1", Value: slog.AnyValue("value1")})

	recWithGroup := slog.NewRecord(time.Now(), slog.LevelInfo, "message with group", pcs[0])
	recWithGroup.AddAttrs(slog.Group("Group-test", slog.String("c-group", "d-group"), "e-group", "f-group"))

	recWithEmptyGroup := slog.NewRecord(time.Now(), slog.LevelInfo, "message empty with group", pcs[0])
	recWithEmptyGroup.AddAttrs(slog.Group("", slog.String("", ""), "", ""))

	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-11-01 15:04:05")
	recWithTime := slog.NewRecord(time.Now(), slog.LevelInfo, "message with time kind", pcs[0])
	recWithTime.AddAttrs(slog.Time("time", parsedTime), slog.Attr{Key: "key1", Value: slog.AnyValue("value1")})

	type fields struct {
		opts           Options
		preformatted   []byte
		unopenedGroups []string
		indentLevel    int
		mu             *sync.Mutex
		out            io.Writer
	}
	type args struct {
		ctx context.Context
		r   slog.Record
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     bool
		logMsgValue string
		checks      []check // checks is a list of checks to run on the result.
	}{
		{name: "Happy flow retriving INFO", fields: fields{mu: &sync.Mutex{}, out: &b}, args: args{r: slog.NewRecord(time.Now(), slog.LevelInfo, "some message", pcs[0])}, wantErr: false, logMsgValue: "some message", checks: []check{
			hasKey("msg="),
			hasKey("time="),
			hasKey("level="),
		}},
		{name: "Happy flow retriving INFO with attributes", fields: fields{mu: &sync.Mutex{}, out: &b}, args: args{r: recWithAttr}, wantErr: false, logMsgValue: "value1", checks: []check{
			hasKey("key1"),
			hasKey("value1"),
		}},
		{name: "Happy flow retriving INFO with group", fields: fields{mu: &sync.Mutex{}, out: &b}, args: args{r: recWithGroup}, wantErr: false, logMsgValue: "value1", checks: []check{
			hasKey("Group-test"),
			hasKey("c-group"),
			hasKey("d-group"),
			hasKey("e-group"),
			hasKey("f-group"),
		}},
		{name: "Happy flow retriving INFO with kind time", fields: fields{mu: &sync.Mutex{}, out: &b}, args: args{r: recWithTime}, wantErr: false, logMsgValue: "message with time kind", checks: []check{
			hasKey("time:"),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ch := &CustomHandler{
				opts:           tt.fields.opts,
				preformatted:   tt.fields.preformatted,
				unopenedGroups: tt.fields.unopenedGroups,
				indentLevel:    tt.fields.indentLevel,
				mu:             tt.fields.mu,
				out:            tt.fields.out,
			}
			if err := ch.Handle(tt.args.ctx, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("CustomHandler.Handle() error = %v, wantErr %v", err, tt.wantErr)
			}

			got := b.String()

			if !strings.Contains(got, tt.logMsgValue) {
				t.Errorf("Custome logger should contain %s", tt.logMsgValue)
			}

			n := map[string]any{"test": got}
			for _, check := range tt.checks {
				if p := check(n); p != "" {
					t.Errorf("%s: %s", p, tt.name)
				}
			}

		})
	}
}

func TestCustomHandler_WithAttrs(t *testing.T) {
	var b bytes.Buffer
	attrs := []slog.Attr{
		slog.String("username", "johndoe"),
		slog.Int("user_id", 42),
		slog.Bool("active", true),
		slog.Float64("score", 98.5),
	}

	type fields struct {
		opts           Options
		preformatted   []byte
		unopenedGroups []string
		indentLevel    int
		mu             *sync.Mutex
		out            io.Writer
	}
	type args struct {
		attrs []slog.Attr
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{name: "logger WithAttrs", fields: fields{mu: &sync.Mutex{}, out: &b}, args: args{attrs}, want: "user_id:42active:truescore:98.5"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &CustomHandler{
				opts:           tt.fields.opts,
				preformatted:   tt.fields.preformatted,
				unopenedGroups: tt.fields.unopenedGroups,
				indentLevel:    tt.fields.indentLevel,
				mu:             tt.fields.mu,
				out:            tt.fields.out,
			}

			got := ch.WithAttrs(tt.args.attrs)
			chCopy, ok := got.(*CustomHandler)
			if !ok {
				t.Errorf("Custom logger should be type of *CustomHandler")
			}
			trimmed := utils.TrimWhitespaces(string(chCopy.preformatted))
			if !strings.Contains(trimmed, tt.want) {
				t.Errorf("Custom logger should contain %s", tt.want)
			}
		})
	}
}

type check func(map[string]any) string

func hasKey(key string) check {
	return func(m map[string]any) string {
		jsonData, err := json.Marshal(m["test"])
		if err != nil {
			return fmt.Sprintf("error converting map to JSON:%q", err)
		}

		got := string(jsonData)

		if !strings.Contains(got, key) {
			return fmt.Sprintf("missing key %q", key)
		}

		return ""
	}
}
