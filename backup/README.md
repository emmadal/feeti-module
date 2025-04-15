# 🌟 Backup Module

This module provides a simple way to backup a PostgreSQL database.

## Features

Light-weight and simple way to backup a PostgreSQL database

### 🚀 Benchmarks

Results of benchmark tests run on Apple M1 (Darwin/ARM64):

#### Long Run (7s)

| Benchmark | Operations | Time/Op | Memory/Op | Allocations/Op |
|---|---:|---:|---:|---:|
| GetBackupTime | 81,262,174 | 97.40 ns | 0 B | 0 |
| GetDBCredentials | 39,666,595 | 211.4 ns | 0 B | 0 |
| FileName | 36,394,063 | 230.1 ns | 48 B | 1 |
| PerformBackup | 266 | 31.44 ms | 853.58 KB | 90 |

#### Short Run (1s)

| Benchmark | Operations | Time/Op | Memory/Op | Allocations/Op |
|---|---:|---:|---:|---:|
| GetBackupTime | 12,332,106 | 97.20 ns | 0 B | 0 |
| GetDBCredentials | 5,691,098 | 210.8 ns | 0 B | 0 |
| FileName | 5,185,474 | 231.3 ns | 48 B | 1 |
| PerformBackup | 34 | 31.77 ms | 853.59 KB | 90 |

> *Command used: `go test -bench=. -benchmem -benchtime=Xs`*
