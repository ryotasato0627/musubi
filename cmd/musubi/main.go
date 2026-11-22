// Command musubi はリポジトリ内の Claude Code SKILL を管理する CLI。
package main

import (
	"os"
	"runtime/debug"

	"github.com/ryotasato0627/musubi/internal/cli"
)

// version はリリース時に -ldflags "-X main.version=..." で埋め込む。
var version = "dev"

func main() {
	os.Exit(cli.Run(resolveVersion(), os.Args[1:], os.Stdout, os.Stderr))
}

// resolveVersion は埋め込まれたバージョンを返す。埋め込みが無ければ、
// go install で導入したときのモジュールのバージョンを使う（NFR-02）。
func resolveVersion() string {
	if version != "dev" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return version
}
