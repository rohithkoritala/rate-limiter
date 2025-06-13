#!/bin/bash

set -e

REDIS_BASE=~/rate-limiter/redis-cluster
PORTS=(7000 7001 7002 7003 7004 7005)

# Create base directories
mkdir -p "$REDIS_BASE"

# Generate config for each node
for PORT in "${PORTS[@]}"; do
  DIR="$REDIS_BASE/$PORT"
  mkdir -p "$DIR"
  cat > "$DIR/redis.conf" <<EOF
port $PORT
cluster-enabled yes
cluster-config-file nodes-$PORT.conf
cluster-node-timeout 5000
appendonly yes
dir $DIR
EOF
  echo "Generated config for port $PORT"
done

# Start each Redis node
for PORT in "${PORTS[@]}"; do
  redis-server "$REDIS_BASE/$PORT/redis.conf" &
  echo "Started Redis on port $PORT"
done

sleep 2  # Give Redis a moment to start

# Create the cluster (no replicas)
REDIS_CLI=$(which redis-cli)

$REDIS_CLI --cluster create \
  127.0.0.1:7000 \
  127.0.0.1:7001 \
  127.0.0.1:7002 \
  127.0.0.1:7003 \
  127.0.0.1:7004 \
  127.0.0.1:7005 \
  --cluster-replicas 1 <<EOF
yes
EOF

# Show cluster info
$REDIS_CLI -p 7000 cluster nodes
