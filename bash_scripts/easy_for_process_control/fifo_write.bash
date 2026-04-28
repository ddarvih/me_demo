#!/usr/bin/env bash
set -uo pipefail
# ВНИМАНИЕ: строки выше удалять нельзя!

while read LINE; do
  echo "$LINE" >pipe
done
