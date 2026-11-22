#!/usr/bin/env python3
"""PreToolUse(Bash) hook: main / develop への直接コミット・push と force push を止める。

.claude/rules/git-workflow.md の「main / develop に直接コミット・push しない。force push しない」を機械的に守らせる。
"""
import json
import os
import re
import subprocess
import sys

PROTECTED = {"main", "develop"}


def deny(reason: str) -> None:
    print(json.dumps({
        "hookSpecificOutput": {
            "hookEventName": "PreToolUse",
            "permissionDecision": "deny",
            "permissionDecisionReason": reason,
        }
    }))
    sys.exit(0)


def current_branch(cwd: str) -> str:
    try:
        out = subprocess.run(
            ["git", "branch", "--show-current"],
            cwd=cwd or None, capture_output=True, text=True, timeout=5,
        )
        return out.stdout.strip()
    except Exception:
        return ""


def main() -> None:
    try:
        data = json.load(sys.stdin)
    except Exception:
        return
    cmd = (data.get("tool_input") or {}).get("command") or ""
    cwd = data.get("cwd") or ""

    # コマンドを ; && || | で分けて、git commit / git push を含む部分だけ見る
    for part in re.split(r"&&|\|\||;|\|", cmd):
        tokens = part.strip().split()
        if len(tokens) >= 2 and tokens[0] == "cd":
            cwd = os.path.join(cwd, os.path.expanduser(tokens[1]))
            continue
        if len(tokens) < 2 or tokens[0] != "git":
            continue
        # git -C <dir> などのオプションを読み飛ばす
        i = 1
        while i < len(tokens) and tokens[i].startswith("-"):
            if tokens[i] == "-C" and i + 1 < len(tokens):
                cwd = os.path.join(cwd, os.path.expanduser(tokens[i + 1]))
            i += 2 if tokens[i] in ("-C", "-c") else 1
        if i >= len(tokens):
            continue
        sub, args = tokens[i], tokens[i + 1:]

        if sub == "push":
            if any(a in ("-f", "--force") or a.startswith("--force") or a.startswith("+") for a in args):
                deny("force push は禁止です（.claude/rules/git-workflow.md）")
            targets = [a.split(":")[-1] for a in args if not a.startswith("-")][1:]
            if any(t in PROTECTED or t in {f"refs/heads/{p}" for p in PROTECTED} for t in targets):
                deny("main / develop への直接 push は禁止です。feature/* から PR を作ってください（.claude/rules/git-workflow.md）")
            if not targets and current_branch(cwd) in PROTECTED:
                deny(f"{current_branch(cwd)} ブランチからの push は禁止です（.claude/rules/git-workflow.md）")

        if sub in ("commit", "merge", "cherry-pick", "revert") and current_branch(cwd) in PROTECTED:
            deny(f"{current_branch(cwd)} ブランチへの直接 {sub} は禁止です。develop から feature/* を切ってください（.claude/rules/git-workflow.md）")


if __name__ == "__main__":
    main()
