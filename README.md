# Distributed Rate Limiter in Go with Redis Cluster

## 🚀 Overview

A high-performance distributed rate limiter built in Go, with **Redis Clusters** and **atomic Lua scripting** for safe, consistent multi-instance coordination. It provides precise token bucket control and is production-ready for cloud or containerized environments.


## ✨ Features

### 🔗 **Distributed Coordination**

* Uses **Redis Cluster** for centralized, horizontally scalable state.
* Ensures all app instances share and respect the same rate-limiting rules.

### 🧠 **Atomic Token Bucket Algorithm**

* Implemented using **Lua scripts** for atomicity and performance.
* Eliminates race conditions under high concurrency.
* Adjustable rate and burst per key (e.g., user ID, API token).

### ♻️ **TTL-Based Expiry & Cleanup**

* Automatic expiration of rate and token data after inactivity.
* Keeps Redis usage efficient and clean.

### ⚙️ **Key Partitioning for Redis Cluster**

* Keys use hash tags to ensure **CROSSSLOT compatibility**.
* Fully functional in sharded, clustered Redis setups.

### 📊 **Real-Time Statistics**

* Tracks requests that are allowed, rate-limited, or failed.
* Optional metrics endpoint for Prometheus integration.


## How to Run

1. **Spin up Redis Cluster** (6 nodes, no Docker required)
```bash
./setup_cluster.sh
```
2. **Run the Go app:**

```bash
go run main.go
```

3. **Call the rate limiter:**

```bash
curl "http://localhost:8080/check?key=user123"
```

4. **Run the stress test:**

```bash
./test_check.sh
```
Or stress test with Grafana:
```bash
node test.js
```

## 🚧 Roadmap

* [x] Lua-based token bucket
* [x] Redis cluster support with slot-safe keys
* [x] TTL-based cleanup
* [x] Multi-instance coordination
* [x] Stress testing & metrics
* [ ] gRPC API support
* [ ] Prometheus metrics endpoint
* [ ] Kubernetes deployment YAMLs


PRs and ideas welcome! I started this project to learn new skills.

MIT License. Free to use and extend.
