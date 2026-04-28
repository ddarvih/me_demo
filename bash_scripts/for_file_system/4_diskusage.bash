#!/usr/bin/env bash
set -euo pipefail
# ВНИМАНИЕ: строки выше удалять нельзя!

dir="$1"

if [[ ! -d "$dir" ]]; then
  echo "problem with given dir"
  exit 1
fi

mapfile -t files < <(find "$dir" -type f)
files_counter="${#files[@]}"
mapfile -t biggest10 < <(find "$dir" -type f -printf "%s %p\n" | sort -nr | head -10)

X=$(find "$dir" -type f -name "*.log" -printf "%s\n" 2>/dev/null |
  awk '{s+=$1} END {print s+0}')
Y=$(find "$dir" -type f -name "*.tmp" -printf "%s\n" 2>/dev/null |
  awk '{s+=$1} END {print s+0}')
Z=$(find "$dir" -type f -name ".*" -printf "%s\n" 2>/dev/null |
  awk '{s+=$1} END {print s+0}')
T=$(find "$dir" -type f -printf "%s\n" |
  awk '{s+=$1} END {print s+0}')

echo "=== Disk Report | $(date +"%Y-%m-%d %H:%M:%S") ==="
echo "Files: $files_counter"
echo "Top-10:"

top_ind=1
while [[ $top_ind -le 10 ]]; do
  cur="${biggest10[$((top_ind - 1))]:-}"
  [[ -z "$cur" ]] && break

  read -r size file <<<"$cur"
  echo "  $top_ind. $file (${size} B)"
  ((top_ind++))
done

echo "Logs: $X B"
echo "Temp: $Y B"
echo "Hidden: $Z B"
echo "--------------------------------------"
echo "Total: $T B"
echo "======================================"
