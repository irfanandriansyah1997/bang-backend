#!/bin/sh
# Enforce the repo commit format:  <type>(<feature>): <message>
#   e.g.  feat(docs): adding readme.md with setup steps
#         fix(work-order): reject complete without photos
# Usage: commit-msg.sh <path-to-commit-message-file>   (git passes it to the commit-msg hook)

msg_file="$1"
subject=$(head -n1 "$msg_file")

# Let git's own generated messages through
case "$subject" in
  Merge\ *|Revert\ *|fixup!\ *|squash!\ *) exit 0 ;;
esac

types='feat|fix|refactor|perf|test|docs|style|build|ci|chore|revert'
# type: required, lowercase | (feature): required, lowercase kebab-case | optional ! for breaking change
pattern="^($types)\([a-z0-9]+(-[a-z0-9]+)*\)!?: [^ ].*$"

if ! printf '%s' "$subject" | grep -qE "$pattern"; then
  echo "✗ Commit message must be:  <type>(<feature>): <message>"
  echo "  example: feat(docs): adding readme.md with setup steps"
  echo "  type:    feat, fix, refactor, perf, test, docs, style, build, ci, chore, revert"
  echo "  feature: required, lowercase kebab-case (docs, auth, work-order, ...)"
  echo "  got:     $subject"
  exit 1
fi

length=$(printf '%s' "$subject" | wc -m | tr -d ' ')
if [ "$length" -gt 72 ]; then
  echo "✗ Commit subject is $length characters; keep it at 72 or fewer."
  exit 1
fi
