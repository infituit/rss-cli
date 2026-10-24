package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestNewRootCommand(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "help",
			args: []string{"help"},
			want: "Available Commands:",
		},
		{
			name: "version",
			args: []string{"version"},
			want: `"version": "debug"`,
		},
		{
			name: "nested help",
			args: []string{"freshrss", "auth", "--help"},
			want: "save FreshRSS credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := new(bytes.Buffer)
			root := NewRootCommand(context.Background(), stdout, new(bytes.Buffer))
			root.SetArgs(tt.args)

			if err := root.Execute(); err != nil {
				t.Fatalf("execute command: %v", err)
			}
			if !strings.Contains(stdout.String(), tt.want) {
				t.Fatalf("output = %q, want substring %q", stdout.String(), tt.want)
			}
		})
	}
}

func TestContextCommands(t *testing.T) {
	t.Setenv("RSS_CLI_CONFIG_DIR", t.TempDir())

	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "create context",
			args: []string{"context", "create", "home"},
			want: `"status": "created"`,
		},
		{
			name: "use context",
			args: []string{"context", "use", "home"},
			want: `"status": "current"`,
		},
		{
			name: "list contexts",
			args: []string{"--output", "table", "context", "list"},
			want: "home",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := new(bytes.Buffer)
			root := NewRootCommand(context.Background(), stdout, new(bytes.Buffer))
			root.SetArgs(tt.args)

			if err := root.Execute(); err != nil {
				t.Fatalf("execute command: %v", err)
			}
			if !strings.Contains(stdout.String(), tt.want) {
				t.Fatalf("output = %q, want substring %q", stdout.String(), tt.want)
			}
		})
	}
}

func TestTracingConfigCommands(t *testing.T) {
	t.Setenv("RSS_CLI_CONFIG_DIR", t.TempDir())

	for _, args := range [][]string{
		{"context", "create", "home"},
		{"context", "use", "home"},
	} {
		root := NewRootCommand(context.Background(), new(bytes.Buffer), new(bytes.Buffer))
		root.SetArgs(args)
		if err := root.Execute(); err != nil {
			t.Fatalf("execute %v: %v", args, err)
		}
	}

	stdout := new(bytes.Buffer)
	root := NewRootCommand(context.Background(), stdout, new(bytes.Buffer))
	root.SetArgs([]string{"tracing", "config", "init", "--endpoint", "http://otel.example/v1/traces"})
	if err := root.Execute(); err != nil {
		t.Fatalf("configure tracing: %v", err)
	}
	if !strings.Contains(stdout.String(), `"status": "configured"`) {
		t.Fatalf("configure output = %q", stdout.String())
	}

	stdout.Reset()
	root = NewRootCommand(context.Background(), stdout, new(bytes.Buffer))
	root.SetArgs([]string{"tracing", "config", "show"})
	if err := root.Execute(); err != nil {
		t.Fatalf("show tracing: %v", err)
	}
	if !strings.Contains(stdout.String(), `"endpoint": "http://otel.example/v1/traces"`) {
		t.Fatalf("show output = %q", stdout.String())
	}

	stdout.Reset()
	root = NewRootCommand(context.Background(), stdout, new(bytes.Buffer))
	root.SetArgs([]string{"tracing", "config", "disable"})
	if err := root.Execute(); err != nil {
		t.Fatalf("disable tracing: %v", err)
	}
	if !strings.Contains(stdout.String(), `"enabled": false`) {
		t.Fatalf("disable output = %q", stdout.String())
	}
}

func TestTracingConfigRejectsInvalidEndpoint(t *testing.T) {
	t.Setenv("RSS_CLI_CONFIG_DIR", t.TempDir())

	for _, args := range [][]string{
		{"context", "create", "home"},
		{"context", "use", "home"},
	} {
		root := NewRootCommand(context.Background(), new(bytes.Buffer), new(bytes.Buffer))
		root.SetArgs(args)
		if err := root.Execute(); err != nil {
			t.Fatalf("execute %v: %v", args, err)
		}
	}

	root := NewRootCommand(context.Background(), new(bytes.Buffer), new(bytes.Buffer))
	root.SetArgs([]string{"tracing", "config", "init", "--endpoint", "otel.example"})
	if err := root.Execute(); err == nil || !strings.Contains(err.Error(), "invalid tracing endpoint") {
		t.Fatalf("invalid endpoint error = %v", err)
	}
}

func TestContextDeleteRequiresConfirmation(t *testing.T) {
	t.Setenv("RSS_CLI_CONFIG_DIR", t.TempDir())

	stdout := new(bytes.Buffer)
	root := NewRootCommand(context.Background(), stdout, new(bytes.Buffer))
	root.SetArgs([]string{"context", "create", "home"})
	if err := root.Execute(); err != nil {
		t.Fatalf("create context: %v", err)
	}

	root = NewRootCommand(context.Background(), stdout, new(bytes.Buffer))
	root.SetArgs([]string{"context", "delete", "home"})
	if err := root.Execute(); err == nil || !strings.Contains(err.Error(), "--confirm") {
		t.Fatalf("delete without confirmation error = %v", err)
	}
}
