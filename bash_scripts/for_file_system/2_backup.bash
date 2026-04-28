#!/usr/bin/env bash
set -euo pipefail
# ВНИМАНИЕ: строки выше удалять нельзя!

act_copy() {
  local file=$1
  local dest=$2
  cp -p "$file" "$dest" 2>/dev/null || {
    echo "problems with copying" >&2
    return 1
  }
}

act_move() {
  local file=$1
  local dest=$2
  mkdir -p "$(dirname "$dest")" 2>/dev/null || return 1
  mv "$file" "$dest" 2>/dev/null || {
    echo "problems with moving" >&2
    return 1
  }
}

act_compress() {
  local backup_dir=$1
  local ver_dir="$backup_dir/.versions"
  if [[ -d "$ver_dir" ]]; then
    cd "$backup_dir"
    if [[ -f "versions.tar.gz" ]]; then
      tar -uf "versions.tar.gz" ".versions/" 2>/dev/null || return 1
    else
      tar -czf "versions.tar.gz" ".versions/" 2>/dev/null || return 1
    fi
  fi
}

compare_files() { # вернет 0 если разные
  [[ ! -f "$2" ]] && return 0

  local inode1 mtime1 size1 inode2 mtime2 size2
  read inode1 mtime1 size1 < <(stat -c "%i %Y %s" "$1")
  read inode2 mtime2 size2 < <(stat -c "%i %Y %s" "$2")
  [[ "$inode1" == "$inode2" ]] ||
    [[ "$mtime1" == "$mtime2" && "$size1" == "$size2" ]] && return 1
  return 0
}

act_backup() {
  local dir=$1
  local backup_dir=$2
  local new=() updated=() skipped=()

  while IFS= read -r cur_file; do
    local path_tail="${cur_file#$dir/}"
    local file="$backup_dir/$path_tail"

    if [[ ! -f "$file" ]]; then
      mkdir -p "$(dirname "$file")"
      if act_copy "$cur_file" "$file"; then
        new+=("$path_tail")
      fi

    elif compare_files "$file" "$cur_file"; then
      local timestamp
      timestamp=$(date -r "$cur_file" +"%Y-%m-%d_%H-%M-%S")

      local ver_file="$backup_dir/.versions/${path_tail}@${timestamp}"
      mkdir -p "$(dirname "$ver_file")"

      if act_move "$file" "$ver_file" && act_copy "$cur_file" "$file"; then
        updated+=("$path_tail")
      fi

    else
      skipped+=("$path_tail")
    fi

  done < <(find "$dir" -type f)

  local original_path
  original_path=$(realpath --relative-to="$HOME" "$dir" 2>/dev/null || echo "$dir")
  echo "$original_path" >"$backup_dir/.original_path"

  echo "=== Backup | $(date +%Y-%m-%d) ==="
  if [[ ${#new[@]} -gt 0 ]]; then
    echo "New: $(
      IFS=,
      echo "${new[*]}"
    )"
  fi
  if [[ ${#updated[@]} -gt 0 ]]; then
    echo "Updated: $(
      IFS=,
      echo "${updated[*]}"
    )"
  fi
  if [[ ${#skipped[@]} -gt 0 ]]; then
    echo "Skipped: $(
      IFS=,
      echo "${skipped[*]}"
    )"
  fi
  echo "---------------------------"
  echo "Summary: added ${#new[@]}, updated ${#updated[@]}, skipped ${#skipped[@]}"
  echo "==========================="
}

act_check() {
  local dir=$1
  local backup_dir=$2
  local new=() changed=() unchanged=()

  while IFS= read -r cur_file; do
    local path_tail="${cur_file#$dir/}"
    local file="$backup_dir/$path_tail"

    if [[ ! -f "$file" ]]; then
      new+=("$path_tail")
    elif compare_files "$file" "$cur_file"; then
      changed+=("$path_tail")
    else
      unchanged+=("$path_tail")
    fi
  done < <(find "$dir" -type f)

  echo "=== Check | $(date +%Y-%m-%d) ==="
  if [[ ${#changed[@]} -gt 0 ]]; then
    echo "Changed: $(
      IFS=,
      echo "${changed[*]}"
    )"
  fi
  if [[ ${#new[@]} -gt 0 ]]; then
    echo "New: $(
      IFS=,
      echo "${new[*]}"
    )"
  fi
  if [[ ${#unchanged[@]} -gt 0 ]]; then
    echo "Unchanged: $(
      IFS=,
      echo "${unchanged[*]}"
    )"
  fi
}

src_dir=""
compress_mode=0
check_mode=0
while [[ $# -gt 0 ]]; do
  case $1 in
    --compress)
      compress_mode=1
      shift
      ;;
    --check)
      check_mode=1
      shift
      ;;
    *)
      if [[ -z "$src_dir" && -d "$1" ]]; then
        src_dir="$1"
      else
        echo "problems with given $1" >&2
        exit 1
      fi
      shift
      ;;
  esac
done
if [[ -z "$src_dir" ]]; then
  echo "problems - no source dir" >&2
  exit 1
fi

backup_dir="$HOME/Backups/$(basename "$(realpath "$src_dir")")-$(date +%F)"

if [[ check_mode -eq 1 ]]; then
  if [[ ! -d "$backup_dir" ]]; then
    echo "problems with backup dir" >&2
    exit 1
  fi
  act_check "$src_dir" "$backup_dir"
else
  mkdir -p "$backup_dir" || {
    echo "problems making backup dir" >&2
    exit 1
  }

  act_backup "$src_dir" "$backup_dir"

  if [[ compress_mode -eq 1 ]]; then
    act_compress "$backup_dir"
  fi
fi
