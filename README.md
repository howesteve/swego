# Swiss Ephemeris library for Go [![GoDoc](https://godoc.org/github.com/howesteve/swego?status.svg)](https://godoc.org/github.com/howesteve/swego)

This is a fork from the original work at [https://github.com/astrotools/swego](https://github.com/astrotools/swego)

I have updated swiss ephemeris files, added module support, and made it compile on go >= 1.18.

This is not a very quality repository. The original projects (Swiss Ephemeris and swego) have 
a lot of problems and I did my best to fix them and come up with something usable, but that's 
about it.

This repository contains multiple ways to interface with the Swiss Ephemeris.
- `swecgo` interfaces with the C library via cgo.
- `swerker` interfaces with the C library via a separate worker or workers.
  - `swerker-stdio` is a worker that runs as a subprocess.

## Pronunciation

The name of this package is pronounced _swie-go_, like cgo: cee-go.

## Docs for sweph files download (used by the lib)

- File needed for astrological purposes is **de441.eph** (de440.eph is the same 
  but cover more years)
- ftp://ssd.jpl.nasa.gov/pub/eph/planets/bsp
- https://rhodesmill.org/skyfield/planets.html
- https://ssd.jpl.nasa.gov/planets/eph_export.html
- <https://www.astro.com/ftp/swisseph/doc/swisseph.htm#_Toc58931065>
- <https://ipnpr.jpl.nasa.gov/progress_report/42-196/196C.pdf>
