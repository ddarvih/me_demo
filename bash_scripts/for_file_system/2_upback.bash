#!/usr/bin/env bash
set -euo pipefail
# ВНИМАНИЕ: строки выше удалять нельзя!

to_mode=0
target_dir=""
backup_name=""
while [[ $# -gt 0 ]]; do
  case "$1" in
    --to)
      to_mode=1
      target_dir="$2"
      shift 2
      ;;
    *)
      if [[ -z "$backup_name" ]]; then
        backup_name="$1"
      fi
      shift
      ;;
  esac
done

if [[ -z "$backup_name" ]]; then
  echo "problem - no backup name" >&2
  exit 1
fi

backup_path="$HOME/Backups/$backup_name"
if [[ ! -d "$backup_path" ]]; then
  echo "backup not found: $backup_path" >&2
  exit 1
fi

if [[ $to_mode -eq 0 ]]; then
  if [[ -f "$backup_path/.original_path" ]]; then
    target_dir=$(cat "$backup_path/.original_path")
    [[ "$target_dir" != /* ]] && target_dir="$HOME/$target_dir"
  else
    echo "problem - no original path found in backup" >&2
    exit 1
  fi
fi
target_dir="$(realpath -m -- "$target_dir")"

mkdir -p "$target_dir" || {
  echo "problems making restore dir $target_dir" >&2
  exit 1
}

find_latest_version() {
  local rel="$1"
  local dir_v="$backup_path/.versions"
  [[ ! -d "$dir_v" ]] && return 0

  local escaped_rel_path=$(printf '%s' "$rel" | sed 's/[.[\*^$(){}+?|]/\\&/g' | sed 's/\//\\\//g')
  local rel_versions=$(find "$dir_v" -type f -regex ".*/${escaped_rel_path}@[0-9]\{4\}-[0-9]\{2\}-[0-9]\{2\}_[0-9]\{2\}-[0-9]\{2\}-[0-9]\{2\}" 2>/dev/null)
  [[ -z "$rel_versions" ]] && return 0

  local latest_timestamp=$(
    printf '%s\n' "$rel_versions" |
      sed 's/.*@//' |
      sort -r |
      head -n 1
  )

  find "$dir_v" -type f -name "$(basename "$rel")@${latest_timestamp}" \
    -path "*/$(dirname "$rel")/*" 2>/dev/null | head -n 1
}

act_ask() {
  local q="$1"
  local ans

  while true; do
    echo "$q [y/n]"
    read ans
    case "$ans" in
      [Yy]) return 0 ;;
      [Nn]) return 1 ;;
      *) echo "Answer [y/n]" ;;
    esac
  done
}

parse_rel_paths() {
  local backup_path="$1"
  local files_list
  files_list=$(
    {
      find "$backup_path" \
        -type f \
        ! -path "$backup_path/.versions/*" \
        ! -name ".original_path" \
        ! -name "versions.tar.gz"
      find "$backup_path/.versions" -type f 2>/dev/null || true
    }
  )
  declare -ga rel_files_paths
  mapfile -t rel_files_paths < <(
    printf '%s\n' "$files_list" |
      sed "s|$backup_path/||" |
      sed "s|$backup_path/.versions/||" |
      sed 's/@.*//' |
      sort -u
  )
}

restored=0
skipped=0

parse_rel_paths "$backup_path"

for rel in "${rel_files_paths[@]}"; do
  v_latest=$(find_latest_version "$rel" || true)
  src=""
  if [[ -n "$v_latest" ]]; then
    src="$v_latest"
  elif [[ -f "$backup_path/$rel" ]]; then
    src="$backup_path/$rel"
  else
    continue
  fi

  dest="$target_dir/$rel"
  mkdir -p "$(dirname "$dest")"
  if [[ -e "$dest" ]]; then
    echo "file exists: $dest"
    if ask_yn "Overwrite?"; then
      cp -p -- "$src" "$dest"
      echo "Restored: $rel"
      restored=$((restored + 1))
    else
      echo "Skipped: $rel"
      skipped=$((skipped + 1))
    fi
  else
    cp -p -- "$src" "$dest"
    echo "Restored: $rel"
    restored=$((restored + 1))
  fi
done

if [[ $restored -gt 0 || $skipped -gt 0 ]]; then
  echo "---------------------------"
  echo "Summary: restored $restored, skipped $skipped"
fi
