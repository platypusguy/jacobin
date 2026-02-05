module jacobin

// go 1.18
// go 1.19   // as of 2022-08 (v. 0.2.1)
// go 1.20   // as of 2023-03-25
// go 1.21   // as of 2023-08-11 (v. 0.4.0) per JACOBIN-330
// go 1.21.4 // as of 2023-11-08
// go 1.24   // as of 2025-02-27 (v. 0.7.0) per JACOBIN-636
// go 1.24.0
go 1.25.x // as of 2026-02-05

require (
	github.com/cespare/xxhash/v2 v2.3.0
	golang.org/x/crypto v0.47.0
	golang.org/x/term v0.39.0
)

require (
	golang.org/x/mod v0.25.0 // indirect
	golang.org/x/sync v0.15.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/telemetry v0.0.0-20240521205824-bda55230c457 // indirect
	golang.org/x/tools v0.34.0 // indirect
)

tool golang.org/x/tools/cmd/deadcode
