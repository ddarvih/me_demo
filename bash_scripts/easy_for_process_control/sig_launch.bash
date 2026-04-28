#!/usr/bin/env bash
set -euo pipefail
# ВНИМАНИЕ: строки выше удалять нельзя!

working_dir="$(dirname "$0")"

bash "$working_dir/6_sig_receiver.bash" &
pid_receiver=$!

bash "$working_dir/6_sig_sender.bash" $pid_receiver

wait $pid_receiver 2>/dev/null || true
