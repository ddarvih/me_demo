#!/usr/bin/env bash
set -uo pipefail
# ВНИМАНИЕ: строки выше удалять нельзя!
if [ $# -ne 5 ]; then
  echo "expected 5 args"
  exit 1
fi

target_script=$1
priority_list=($2 $3 $4)
p_s=()
log_file=$5

touch "$log_file" 2>/dev/null || {
  echo "unable to use $log_file to write" >&2
  exit 1
}

for i in {0..2}; do
  nice -n "${priority_list[$i]}" bash "$target_script" &
  p_s[$i]=$!
done
echo "${p_s[@]}"

get_process_info() {
  local pid=$1
  ps -h -o ni -o %cpu -p $pid 2>/dev/null
}

log_to_file() {
  log=$(date +%s)
  for i in {0..2}; do
    local pid=${p_s[$i]}
    info=($(get_process_info "$pid"))
    local cur_priority=${info[0]:-0}
    local cur_cpu=${info[1]:-0}
    log="$log:$pid:$cur_priority:$cur_cpu"
  done
  echo "$log" >>$log_file
}

for i in {1..3}; do
  log_to_file
  sleep 5
done

past_15_s=()

max_res=0
min_res=0
max_cpu=0
min_cpu=0
for i in {0..2}; do
  info=($(get_process_info ${p_s[$i]}))
  cur_cpu=${info[1]//%/}
  cur_cpu=${cur_cpu:-0}
  if [ $i -eq 0 ]; then
    max_cpu=$cur_cpu
    min_cpu=$cur_cpu
  fi

  if (($(echo "$cur_cpu > $max_cpu" | bc -l 2>/dev/null))); then
    max_cpu=$cur_cpu
    max_res=$i
  fi
  if (($(echo "$cur_cpu < $min_cpu" | bc -l 2>/dev/null))); then
    min_cpu=$cur_cpu
    min_res=$i
  fi
done

renice -n +10 ${p_s[$max_res]} 2>/dev/null
renice -n -10 ${p_s[$min_res]} 2>/dev/null || echo "some problems with negative renice (not sudo?)" >&2

for i in {1..3}; do
  log_to_file
  sleep 5
done

kill ${p_s[@]} 2>/dev/null
