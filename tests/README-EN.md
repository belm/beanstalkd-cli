# Beanstalkd CLI Test Suite

[🇨🇳 Chinese](README.md) | [🇺🇸 English](README-EN.md)

A comprehensive collection of automated tests that validates every feature of the Beanstalkd CLI and Web UI integration.

## Test Files Overview

### 1. `connection_test.go`
- Basic connectivity
- Dial timeouts
- Failure scenarios
- Multiple concurrent connections
- Connection benchmarking

### 2. `job_operations_test.go`
- Put jobs (standard, high priority, delayed, UTF‑8, JSON)
- Reserve jobs with timeout
- Delete jobs
- Release jobs with priority/delay updates
- Bury jobs
- Touch jobs
- Kick jobs (bulk + single)
- Benchmark for `put`

### 3. `peek_operations_test.go`
- Peek by ID
- Peek ready/delayed/buried jobs
- Graceful handling of missing jobs

### 4. `tube_operations_test.go`
- List tubes
- Manage multiple tubes
- Tube statistics
- Isolation guarantees
- Concurrent tube access

### 5. `stats_test.go`
- Server metrics
- Job statistics

### 6. `integration_test.go`
- Producer ➜ consumer workflow
- Retry and release flow
- Priority queue ordering

### Helper & Scripts
- `test_helper.go` – shared helpers (cleanup, assertions, builders)
- `run_tests.sh` – guided test runner with coverage output

## Running the Tests

### Option 1 – Makefile (recommended)
```bash
make test            # run all tests
make test-verbose    # verbose mode
make bench           # run benchmarks
make coverage        # generate coverage report (HTML)
```

### Option 2 – Test script
```bash
cd tests
./run_tests.sh
```

### Option 3 – Native go toolchain
```bash
cd tests
go test -v                  # run everything
go test -v -run TestPutJob  # run a specific test
go test -bench=. -benchmem  # benchmarks
```

## Coverage Summary

✅ Connection handling  
✅ Job lifecycle (put/reserve/delete/release/bury/kick/touch)  
✅ Peek operations  
✅ Tube management  
✅ Statistics endpoints  
✅ Priority & delay mechanics  
✅ Error handling  
✅ Concurrency checks  
✅ End-to-end integration flows  
✅ Performance benchmarks

## Best Practices

1. **Isolation** – dedicate distinct tube names inside each test case.  
2. **Cleanup** – rely on helper functions with `defer` to remove temp jobs.  
3. **Skipping** – tests auto-skip when Beanstalkd is unreachable.  
4. **Logging** – rich `t.Log` output makes CI debugging easier.  
5. **Concurrency** – goroutines probe isolation and race conditions.

## Troubleshooting

- **Failures** – verify Beanstalkd is running on `127.0.0.1:11300`; inspect logs.  
- **Skipped tests** – indicates the service could not be reached.  
- **Timeouts** – increase reserve timeouts or run in blocking mode.

## Continuous Integration

Example GitHub Actions snippet:
```yaml
steps:
  - name: Start beanstalkd
    run: beanstalkd -l 127.0.0.1 -p 11300 &
  - name: Run test suite
    run: cd tests && go test -v
  - name: Coverage report
    run: cd tests && go test -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html
```

---

The suite is production-ready and can be dropped into any CI/CD pipeline to guard against regressions.
