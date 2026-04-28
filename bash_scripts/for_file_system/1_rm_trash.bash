#!/usr/bin/env bash
set -euo pipefail
# ВНИМАНИЕ: строки выше удалять нельзя!

trash_dir=$HOME/.trash
trash_log="$HOME/.trash.log"

mkdir -p "$trash_dir" 2>/dev/null || {
  echo "problems with ~/.trash" >&2
  exit 1
}

touch "$trash_log" 2>/dev/null || {
  echo "problems with ~/.trash.log" >&2
  exit 1
}

v_mode=0
p_mode=0

act_v() {
  if [[ $v_mode -eq 1 ]]; then
    echo "$@"
  fi
}

act_p() {
  local cur_file=$1
  if [[ $p_mode -eq 1 ]]; then
    read -p "Want to remove '$cur_file'? [y/n] " decision
    if [[ "$decision" != "y" && "$decision" != "Y" ]]; then
      return 1
    fi
  fi
  return 0
}

while [[ $# -gt 0 ]]; do
  arg="$1"
  if [[ "$arg" == -- ]]; then
    shift
    break
  fi

  if [[ "$arg" == -* ]]; then
    opts="${arg#-}"
    for ((i = 0; i < ${#opts}; i++)); do
      cur="${opts:i:1}"
      case "$cur" in
        v) v_mode=1 ;;
        p) p_mode=1 ;;
        *)
          echo "Unknown option: -$cur" >&2
          exit 1
          ;;
      esac
    done
    shift
    continue
  fi
  break
done

if [[ $# -eq 0 ]]; then
  echo "problem - no file names given" >&2
  exit 1
fi

act_rm() {
  local cur_file="$1"
  if [[ ! -e "$cur_file" ]]; then
    echo "'$cur_file' does not exist" >&2
    return 0
  fi

  if [[ ! -f "$cur_file" ]]; then
    echo "'$cur_file' is not a regular file" >&2
    return 0
  fi

  if ! act_p "$cur_file"; then
    act_v "skipped '$cur_file'"
    return 0
  fi

  act_v "removing $cur_file"
  local original_path=$(realpath -- "$cur_file")

  act_trash_op "$original_path"
}

act_trash_op() {
  local original_path="$1"
  local inode=$(stat -c %i -- "$original_path")
  local timestamp=$(date '+%Y-%m-%d_%H-%M-%S')
  local sz=$(stat -c %s -- "$original_path")

  local link_name="$(basename -- "$original_path")@${timestamp}@${inode}"
  local link_path="$trash_dir/$link_name"

  if ! ln -- "$original_path" "$link_path" 2>/dev/null; then
    echo "Failed to create hard link for '$original_path'" >&2
    return 1
  fi
  if ! rm -- "$original_path" 2>/dev/null; then
    echo "Failed to remove '$original_path'" >&2
    rm -- "$link_path" 2>/dev/null
    return 1
  fi

  echo "$original_path | $link_name | $inode | $sz | $timestamp" >>"$trash_log"

  act_v "Logged: $link_name"
  return 0
}

for cur_removing_file in "$@"; do
  act_rm "$cur_removing_file" || true
done

act_v "Done"
