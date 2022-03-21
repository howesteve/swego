# Swiss Ephemeris library for Go [![GoDoc](https://godoc.org/github.com/howesteve/swego?status.svg)](https://godoc.org/github.com/howesteve/swego)

This is a fork from the original work at [https://github.com/astrotools/swego](https://github.com/astrotools/swego)

I have updated swiss ephemeris files, added module support, and made it compile on go 1.18.

Not very quality repository - do not put a lot of faith on it.

This repository contains multiple ways to interface with the Swiss Ephemeris.
- `swecgo` interfaces with the C library via cgo.
- `swerker` interfaces with the C library via a separate worker or workers.
  - `swerker-stdio` is a worker that runs as a subprocess.

## Pronunciation
The name of this package is pronounced _swie-go_, like cgo: cee-go.
