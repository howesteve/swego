# Swiss Ephemeris library for Go [![GoDoc](https://godoc.org/github.com/dwlnetnl/swego?status.svg)](https://godoc.org/github.com/dwlnetnl/swego)
Package swego allows access to the Swiss Ephemeris from Go.

## Implemented C functions
Currently the following subset of the C API is implemented:
- `swe_version`
- `swe_calc`
- `swe_calc_ut`
- `swe_close`
- `swe_set_ephe_path`
- `swe_set_jpl_file` (via arguments passed to method)
- `swe_get_planet_name`
- `swe_set_topo` (via arguments passed to method)
- `swe_set_sid_mode` (via arguments passed to method)
- `swe_get_ayanamsa` (*you should use `swe_get_ayanamsa_ex`*)
- `swe_get_ayanamsa_ut` (*you should use `swe_get_ayanamsa_ex_ut`*)
- `swe_get_ayanamsa_ex`
- `swe_get_ayanamsa_ex_ut`
- `swe_get_ayanamsa_name`
- `swe_julday`
- `swe_revjul`
- `swe_utc_to_jd`
- `swe_jdet_to_utc`
- `swe_jdut1_to_utc`
- `swe_houses`
- `swe_houses_ex`
- `swe_houses_armc`
- `swe_house_pos`
- `swe_house_name`
- `swe_deltat` (*you should use `swe_deltat_ex`*)
- `swe_deltat_ex`
- `swe_time_equ`
- `swe_lmt_to_lat`
- `swe_lat_to_lmt`
- `swe_sidtime0`
- `swe_sidtime`
- `swe_set_tid_acc` (handled internally by the C library)

### What is the deal with that _via arguments passed to method_?
The reason is to eliminate the number of calls a user has to make. This is in contrast with the C API that requires you to call `swe_set_topo` before `swe_calc`. when you are calculating the topocentric position of Venus. Only calling a single C function is important in the context of Go because you like to minimize the number of calls to C.

Currently the implementation is smart about when to call `swe_set_topo`, but it figures this out in via a C call separate of the calculation. So here's room for improvement but this can be done without changing the public API.
