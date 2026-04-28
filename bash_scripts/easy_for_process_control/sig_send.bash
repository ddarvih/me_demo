#!/usr/bin/env bash
set -euo pipefail
# ВНИМАНИЕ: строки выше удалять нельзя!

pid_receiver=$1

while read LINE; do
  case $LINE in
    "+")
      kill -USR1 $pid_receiver
      ;;
    "*")
      kill -USR2 $pid_receiver
      ;;
    "TERM")
      kill -TERM $pid_receiver
      ;;
  esac
  sleep 1
done
