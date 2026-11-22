package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_NFR06_ExitCode(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantCode   int
		wantStdout string
		wantStderr string
	}{
		{name: "version", args: []string{"--version"}, wantCode: ExitOK, wantStdout: "musubi version 1.2.3\n"},
		{name: "no args shows help", args: nil, wantCode: ExitOK, wantStdout: "Usage:"},
		{name: "common options accepted", args: []string{"--config", "c.toml", "--root", "r", "--no-color", "--quiet", "--verbose"}, wantCode: ExitOK},
		{name: "unknown flag is usage error", args: []string{"--unknown"}, wantCode: ExitUsage, wantStderr: "unknown flag: --unknown"},
		{name: "unknown command is usage error", args: []string{"nope"}, wantCode: ExitUsage, wantStderr: "musubi:"},
		{name: "-v is not version", args: []string{"-v"}, wantCode: ExitUsage, wantStderr: "unknown shorthand flag"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			got := Run("1.2.3", tt.args, &stdout, &stderr)
			if got != tt.wantCode {
				t.Errorf("exit code = %d, want %d (stderr: %q)", got, tt.wantCode, stderr.String())
			}
			if tt.wantStdout != "" && !strings.Contains(stdout.String(), tt.wantStdout) {
				t.Errorf("stdout = %q, want to contain %q", stdout.String(), tt.wantStdout)
			}
			if tt.wantStderr != "" && !strings.Contains(stderr.String(), tt.wantStderr) {
				t.Errorf("stderr = %q, want to contain %q", stderr.String(), tt.wantStderr)
			}
		})
	}
}

func TestNewRootCommand_CommonOptions(t *testing.T) {
	var opts GlobalOptions
	cmd := NewRootCommand("dev", &opts)
	cmd.SetArgs([]string{"--config", "x/config.toml", "--root", "testdata/a/.claude/skills", "--no-color", "--quiet", "--verbose"})
	cmd.SetOut(&bytes.Buffer{})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	want := GlobalOptions{Config: "x/config.toml", Root: "testdata/a/.claude/skills", NoColor: true, Quiet: true, Verbose: true}
	if opts != want {
		t.Errorf("opts = %+v, want %+v", opts, want)
	}
}
