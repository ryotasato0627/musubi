// Package cli は musubi の引数解析とコマンド定義（cobra）を持つ。
// ロジックは持たず、core / analyzers / rules を呼び出して結果を表示するだけにする（design.md 1・5 章）。
package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// 終了コード（design.md 5.1、NFR-06）。
const (
	// ExitOK は成功を表す。
	ExitOK = 0
	// ExitUsage は使い方の誤り・設定エラーを表す。
	ExitUsage = 2
)

// GlobalOptions は全コマンド共通のオプション（design.md 5.1）。
// この段階では定義だけで、値を使う処理は各コマンドの実装で追加する。
type GlobalOptions struct {
	Config  string // --config: 設定ファイルのパス（既定は .musubi/config.toml）
	Root    string // --root: 探索ルートの上書き（テスト用に 1 つだけ）
	NoColor bool   // --no-color: 色を付けない
	Quiet   bool   // --quiet: 出力を減らす
	Verbose bool   // --verbose: 詳細を出力する
}

// NewRootCommand は musubi のルートコマンドを作る。
// opts には共通オプションの値が入る。
func NewRootCommand(version string, opts *GlobalOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "musubi",
		Short:         "Claude Code の SKILL を管理し、依存関係を検証・可視化する",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	// cobra 既定の --version は -v の短縮形を付けるため、--verbose と紛れないよう短縮形なしで定義する。
	cmd.Flags().Bool("version", false, "バージョンを表示する")

	pf := cmd.PersistentFlags()
	pf.StringVar(&opts.Config, "config", "", "設定ファイルのパス（既定は .musubi/config.toml）")
	pf.StringVar(&opts.Root, "root", "", "探索ルートを上書きする（テスト用。1 つだけ）")
	pf.BoolVar(&opts.NoColor, "no-color", false, "色を付けずに出力する")
	pf.BoolVar(&opts.Quiet, "quiet", false, "出力を減らす")
	pf.BoolVar(&opts.Verbose, "verbose", false, "詳細を出力する")
	return cmd
}

// Run はコマンドを実行し、終了コードを返す。
func Run(version string, args []string, stdout, stderr io.Writer) int {
	var opts GlobalOptions
	cmd := NewRootCommand(version, &opts)
	cmd.SetArgs(args)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	if err := cmd.Execute(); err != nil {
		// この段階で起き得るエラーは引数・オプションの誤りだけなので、使い方の誤りとして扱う。
		_, _ = fmt.Fprintf(stderr, "musubi: %v\n", err)
		return ExitUsage
	}
	return ExitOK
}
