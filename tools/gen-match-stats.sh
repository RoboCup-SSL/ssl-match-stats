#!/bin/bash

set -euo pipefail

readonly base_dir="${1}"
readonly parallel=20

function process() {
  local -r dir=$1
  local -r tournament=$2
  local -r division=$3
  local -r match_stats_dir="$dir/match-stats"

  echo "Processing $dir"

  if [[ ! -d "$match_stats_dir" ]]; then
    mkdir "$match_stats_dir"
    echo "Generating match stats for $dir"
    ssl-match-stats \
      -targetDir "${match_stats_dir}" \
      -parallel="$parallel" \
      "$dir/"*.log.gz
  fi

  echo "Importing match stats for $dir into DB"
  ssl-match-stats-db \
    -sqlDbSource="postgres://ssl_match_stats:ssl_match_stats@localhost:5432/ssl_match_stats?sslmode=disable" \
    -parallel="$parallel" \
    -tournament="$tournament" \
    -division="$division" \
    "${match_stats_dir}"/*.bin
}

process "$base_dir/2023/div-a" RoboCup2023 DivA
process "$base_dir/2023/div-b" RoboCup2023 DivB
process "$base_dir/2022/div-a" RoboCup2022 DivA
process "$base_dir/2022/div-b" RoboCup2022 DivB
process "$base_dir/2021" RoboCup2021 none
process "$base_dir/2019/div-a" RoboCup2019 DivA
process "$base_dir/2019/div-b" RoboCup2019 DivB
process "$base_dir/2018/div-a" RoboCup2018 DivA
process "$base_dir/2018/div-b" RoboCup2018 DivB
process "$base_dir/2017" RoboCup2017 none
process "$base_dir/2016" RoboCup2016 none