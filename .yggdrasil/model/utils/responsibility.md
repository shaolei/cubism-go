# Utils — Responsibility

## Identity
Utility packages providing shared helper functions with no domain knowledge.

## Responsibility
- **internal/utils**: `ParseVersion(v uint32) string` — converts native version uint32 to "major.minor.patch" string
- **internal/strings**: `GoString(p uintptr) string` — converts a C null-terminated string pointer to Go string without CGO
- **renderer/utils**: `Normalize(x, n, m float32) float32` — maps x from range [n,m] to [-1,1], returns 0 if n==m

## NOT Responsible For
- Any domain logic

## Key Invariants
- `ParseVersion`: format is `major = v >> 24`, `minor = (v >> 16) & 0xFF`, `patch = v & 0xFFFF`
- `GoString`: handles nil pointer (returns empty string), reads until null byte
- `Normalize`: division-by-zero protection when n == m
