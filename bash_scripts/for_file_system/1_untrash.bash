#!/usr/bin/env bash
set -euo pipefail

home="$HOME"
trash_dir="$home/.trash"
log_trash="$home/.trash.log"
lost_dir="$home/restore-lost"

if [[ ! -d "$trash_dir" ]]; then
  echo "no trash directory found" >&2
  exit 1
fi

if [[ ! -f "$log_trash" ]]; then
  echo "no log file found" >&2
  exit 1
fi

ignore_mode=0
overwrite_mode=0
unique_mode=0
to_mode=0
target_dir=""

while (($# > 0)); do
  case "$1" in
    --ignore)
      ignore_mode=1
      overwrite_mode=0
      unique_mode=0
      shift
      ;;
    --overwrite)
      ignore_mode=0
      overwrite_mode=1
      unique_mode=0
      shift
      ;;
    --unique)
      ignore_mode=0
      overwrite_mode=0
      unique_mode=1
      shift
      ;;
    --to)
      if [[ -z "${2:-}" ]]; then
        echo "problem - no argument for --to" >&2
        exit 1
      fi
      to_mode=1
      target_dir="$2"
      shift 2
      ;;

    --*)
      echo "unknown option $1  - try without it" >&2
      exit 1
      ;;
    *)
      break
      ;;
  esac
done

if ((ignore_mode + overwrite_mode + unique_mode == 0)); then
  ignore_mode=1
fi

if (($# == 0)); then
  echo "no restore pattern provided" >&2
  exit 1
fi

pattern=$1

parse_inodes() {
  local log_file="$1"
  local pattern="$2"

  declare -gA log_data
  declare -g inodes_sorted=()
  while IFS='|' read -r original_path link_name inode size timestamp; do
    original_path=$(echo "$original_path" | xargs)
    link_name=$(echo "$link_name" | xargs)
    inode=$(echo "$inode" | xargs)
    timestamp=$(echo "$timestamp" | xargs)

    [[ -z "$inode" ]] && continue

    if [[ "$link_name" == *"@"*"@"*"$inode" ]]; then
      file_name="${link_name%%@*}"
      if [[ "$file_name" == $pattern ]]; then
        log_data["$inode"]="$original_path|$link_name|$timestamp"
        if [[ ! " ${inodes_sorted[*]} " =~ " $inode " ]]; then
          inodes_sorted+=("$inode")
        fi
      fi
    fi
  done <"$log_file"
}

act_untrash() {
  local dir="$1"
  local file="$2"
  local link="$3"

  if [[ -f "$dir/$file" ]]; then
    if ((ignore_mode)); then
      echo "unable to restore $file. File already exists in $dir"
      return
    fi

    if ((overwrite_mode)); then
      rm -f "$dir/$file"
      ln "$link" "$dir/$file"
      rm "$link"
      return
    fi

    if ((unique_mode)); then
      local counter=1
      local base="${file%.*}"
      local ext="${file##*.}"

      [[ "$base" == "$file" ]] && ext="" || ext=".$ext"

      while [[ -f "$dir/$base($counter)$ext" ]]; do
        ((counter++))
      done

      ln "$link" "$dir/$base($counter)$ext"
      rm "$link"
      return
    fi
  else
    ln "$link" "$dir/$file"
    rm "$link"
  fi
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

parse_inodes "$log_trash" "$pattern"
for inode in "${inodes_sorted[@]}"; do
  IFS='|' read -r original_path link_name timestamp <<<"${log_data[$inode]}"
  link_path="$trash_dir/$link_name"
  file_name="${link_name%%@*}"

  if [[ ! -f "$link_path" ]]; then
    echo "didnt find $link_name in trash -- skipping" >&2
    continue
  fi

  echo "Found: $original_path (deleted at $timestamp)"
  if ! act_ask "Restore?"; then
    continue
  fi

  if ((to_mode == 0)); then
    original_dir=$(dirname "$original_path")
    if [[ -d "$original_dir" ]]; then
      act_untrash "$original_dir" "$file_name" "$link_path"
    else
      mkdir -p "$lost_dir"
      act_untrash "$lost_dir" "$file_name" "$link_path"
    fi
  else
    mkdir -p "$target_dir"
    act_untrash "$target_dir" "$file_name" "$link_path"
  fi
done
