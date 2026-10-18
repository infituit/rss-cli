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
