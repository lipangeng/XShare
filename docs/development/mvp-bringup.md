# MVP Bring-up and Verification

This guide is the minimum bring-up flow for the XShare MVP workspace.

## Prerequisites

- `go` on `PATH` (for core tests).
- `buf` on `PATH` (for protocol generation).
- Java runtime available (`java` on `PATH`) for Gradle.
- Android Gradle wrapper present and executable at `android/gradlew`.

If any prerequisite is missing, `tools/verify-mvp.sh` exits non-zero with a clear error.

## Fast Verification

From repository root:

```bash
bash tools/verify-mvp.sh
```

This runs, in order:

1. Prerequisite checks.
2. `buf generate` in `protocol/`.
3. `go test ./...` in `core/go/`.
4. `./gradlew :app:testDebugUnitTest` in `android/`.

## Focused Controller Contract Check

To validate controller start/stats contract only:

```bash
cd core/go
go test ./pkg/controller -v
```

The contract test `TestForwardStartThenStatsAvailableContract` ensures stats are readable immediately after starting forwarding.
