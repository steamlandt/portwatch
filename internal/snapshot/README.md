# snapshot

The `snapshot` package provides point-in-time capture and comparison of open port lists.

## Overview

A `Store` persists port snapshots to a JSON file on disk. On each monitor cycle the current port list can be saved, and the previous snapshot loaded for comparison.

`Diff` computes which ports were **added** or **removed** between two snapshots, making it easy to drive alerts.

## Usage

```go
store := snapshot.New("/var/lib/portwatch/snapshot.json")

// Save current ports
store.Save(currentPorts)

// Load previous snapshot
prev, _ := store.Load()

// Compare
added, removed := snapshot.Diff(prev, snapshot.Snapshot{Ports: currentPorts})
```

## File format

```json
{
  "ports": [
    {"number": 80, "protocol": "tcp"},
    {"number": 443, "protocol": "tcp"}
  ],
  "captured_at": "2024-01-15T10:30:00Z"
}
```
