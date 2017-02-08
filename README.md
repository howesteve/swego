# Swiss Ephemeris library for Go [![GoDoc](https://godoc.org/github.com/astrotools/swego?status.svg)](https://godoc.org/github.com/astrotools/swego)

This repository contains multiple ways to interface with the Swiss Ephemeris.
- `swecgo` interfaces with the C library via cgo.
- `swerker` interfaces with the C library via a separate worker or workers.
  - `swerker-stdio` is a worker that runs as a subprocess.

## Pronunciation
The name of this package is pronounced _swie-go_, like cgo: cee-go.
