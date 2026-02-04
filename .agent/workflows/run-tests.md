---
description: Run tests and check code coverage for Watchman
---

# Run Tests and Coverage

This workflow guides running tests and analyzing coverage.

## Prerequisites
- Go 1.25.6+ installed
- Project dependencies installed

## Steps

### 1. Run All Tests with Race Detection
// turbo
```bash
make test
```

### 2. Run Quick Tests (No Race Detector)
// turbo
```bash
make test-short
```

### 3. View Coverage Report
// turbo
```bash
make coverage
```

### 4. Open Coverage in Browser
```bash
make coverage-html
```

### 5. Check Coverage Threshold
Target: **60%+ statement coverage** (project minimum)
Target: **80%+ for internal/** packages

Parse coverage percentage:
// turbo
```powershell
go tool cover -func coverage.out | Select-String "total:"
```

### 6. Run Specific Package Tests
```bash
go test -v -race ./internal/config/...
go test -v -race ./internal/database/...
```

### 7. Run Single Test
```bash
go test -v -race -run TestFunctionName ./internal/package/...
```

### 8. Run Tests with Verbose Output
```bash
go test -v -race -count=1 ./...
```

## Test Patterns (Project Standard)

### Table-Driven Tests
```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:  "valid case",
            input: "test",
            want:  "expected",
        },
        {
            name:    "error case",
            input:   "",
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := FunctionUnderTest(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("got = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Mocking Interfaces
```go
type mockDatabase struct {
    pingErr error
}

func (m *mockDatabase) Ping() error {
    return m.pingErr
}

func TestWithMock(t *testing.T) {
    mock := &mockDatabase{pingErr: nil}
    // Use mock in tests
}
```

## Windows-Specific Testing
Some tests require Windows:
```go
//go:build windows

package service

func TestWindowsService(t *testing.T) {
    // Windows-only test
}
```

## CI Integration
Tests run automatically via GitHub Actions on:
- Every push to main
- Every pull request

## Coverage Reports
Coverage file: `coverage.out`
Format: Go coverage profile

View in different formats:
```bash
# Function-level
go tool cover -func coverage.out

# HTML report
go tool cover -html coverage.out -o coverage.html
```

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Race detector slow | Use `make test-short` |
| Windows-only tests fail on CI | Add `//go:build windows` |
| Flaky tests | Add `-count=1` to disable caching |
| Import cycle | Check package dependencies |

## Checklist
- [ ] All tests pass
- [ ] Race detector finds no issues
- [ ] Coverage meets threshold (60%+)
- [ ] No flaky tests
- [ ] Windows-specific tests properly tagged
