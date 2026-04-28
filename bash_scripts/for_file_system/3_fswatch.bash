#!/usr/bin/env bash
set -euo pipefail
# ВНИМАНИЕ: строки выше удалять нельзя!

watching_dir=$1
log_file="${2:-}"

if [ ! -d "$watching_dir" ]; then
  echo "problem with given dir"
  exit 1
fi

inotifywait -mr -e modify,create,delete,move \
  --format '%T %w%f %e' \
  --timefmt '%Y-%m-%d %H:%M:%S' \
  "$watching_dir" |
  while read -r timestamp filepath event; do
    if [ -n "$log_file" ]; then
      echo "[$timestamp] $filepath $event" >>"$log_file"
    else
      echo "[$timestamp] $filepath $event"
    fi
  done
