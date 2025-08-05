#!/usr/bin/env bash
set -euo pipefail

# Cấu hình
TOTAL=10000        # tổng số request
CONCURRENCY=1000  # sau mỗi 10 request chờ một lần
URL="http://localhost:8081/process"
PAYLOAD='{"foo":"bar"}'

for i in $(seq 1 "$TOTAL"); do
  {
    # Gửi request và chỉ lấy HTTP status
    http_code=$(curl -s -o /dev/null \
                     -H "Content-Type: application/json" \
                     -X POST "$URL" \
                     -d "$PAYLOAD" \
                     -w "%{http_code}")

    echo "Req #$i → HTTP $http_code"
  } &    # chạy background

  # Mỗi CONCURRENCY job, chờ tất cả hoàn thành rồi tiếp
  if (( i % CONCURRENCY == 0 )); then
    wait
  fi
done

# Chờ nốt các job còn lại
wait

echo "Done sending $TOTAL requests in batches of $CONCURRENCY."