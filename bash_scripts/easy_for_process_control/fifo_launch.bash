#!/usr/bin/env bash
set -uo pipefail
# ВНИМАНИЕ: строки выше удалять нельзя!

m_pipe="pipe"
working_dir="$(dirname "$0")"
rm -f "$m_pipe"
mkfifo "$m_pipe" || {
  echo "making pipe problems" >&2
  exit 1
}
trap 'rm -f "$m_pipe"' EXIT

bash "$working_dir/5_fifo_reader.bash" &
pid_reader=$!

bash "$working_dir/5_fifo_writer.bash"

wait $pid_reader 2>/dev/null || true
