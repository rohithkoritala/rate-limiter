#!/bin/bash

KEY="user123"
URL="http://localhost:8080/check?key=$KEY"
REQUESTS=100
CONCURRENCY=10

# Temp files for counters
TMP_DIR=$(mktemp -d)
SUCCESS_FILE="$TMP_DIR/success"
RATE_LIMITED_FILE="$TMP_DIR/ratelimited"
FAILED_FILE="$TMP_DIR/failed"
: > "$SUCCESS_FILE"
: > "$RATE_LIMITED_FILE"
: > "$FAILED_FILE"

echo "Starting stress test: $REQUESTS requests with $CONCURRENCY concurrency"
echo

function worker() {
  for i in $(seq 1 $((REQUESTS / CONCURRENCY))); do
    RESPONSE=$(curl -s -w "%{http_code}" -o /tmp/curl_output.$$ "$URL")
    BODY=$(cat /tmp/curl_output.$$)

    if [[ "$BODY" == *"allowed"* ]]; then
      echo 1 >> "$SUCCESS_FILE"
    elif [[ "$BODY" == *"rate limit exceeded"* ]]; then
      echo 1 >> "$RATE_LIMITED_FILE"
    else
      echo 1 >> "$FAILED_FILE"
    fi
  done
}

start=$(date +%s)

for i in $(seq 1 $CONCURRENCY); do
  worker &
done

wait

end=$(date +%s)
duration=$((end - start))

# Final tallies
SUCCESS=$(wc -l < "$SUCCESS_FILE")
RATE_LIMITED=$(wc -l < "$RATE_LIMITED_FILE")
FAILED=$(wc -l < "$FAILED_FILE")

echo
echo "======================="
echo "Stress Test Complete"
echo "======================="
echo "Total Requests:   $REQUESTS"
echo "Success:          $SUCCESS"
echo "Rate Limited:     $RATE_LIMITED"
echo "Other Failures:   $FAILED"
echo "Total Time:       ${duration}s"
echo "Requests/sec:     $((REQUESTS / duration))"

rm -rf "$TMP_DIR" /tmp/curl_output.$$


