#!/usr/bin/env bash
set -euo pipefail
# ВНИМАНИЕ: строки выше удалять нельзя!
x=1

got_usr1() {
  let x=$x+2
  echo $x
}

got_usr2() {
  let x=$x*2
  echo $x
}

got_term() {
  exit 0
}

trap 'got_usr1' USR1
trap 'got_usr2' USR2
trap 'got_term' TERM

while true; do
  sleep 1
done
