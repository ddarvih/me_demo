#!/usr/bin/env bash
set -uo pipefail
# ВНИМАНИЕ: строки выше удалять нельзя!

touch data.log 2>/dev/null || {
  echo "data.log file problems" >&2
  exit 1
}

if [[ ! -p pipe ]]; then
  echo "pipe problems" >&2
  exit 1
fi

LINE=""
mode="+"

calc() {
  local args=($1)
  local res
  case $mode in
    "+")
      res=0
      for a in "${args[@]}"; do
        res=$((res + a))
      done
      ;;
    "*")
      res=1
      for a in "${args[@]}"; do
        res=$((res * a))
      done
      ;;
  esac
  echo "${mode}:${args[*]}:${res}" >>data.log
}

tail -f pipe | while read -r LINE; do
  case $LINE in
    "+")
      mode="+"
      ;;
    "*")
      mode="*"
      ;;
    "STOP")
      exit 0
      ;;
  esac

  if [[ $LINE =~ ^[-[:digit:][:space:]]+$ ]] && [[ -n "$LINE" ]]; then
    calc "$LINE"
  fi
done
