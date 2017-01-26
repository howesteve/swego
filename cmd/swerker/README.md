# Swiss Ephemeris worker
The Swiss Ephemeris is not safe for use from multiple threads. To use multiple
processor cores, the Swiss Ephemeris is linked into a simple worker that can be
parallelized by running multiple copies simultaneously. The worker act as a RPC
server that can handle a single request at a time.

## Implemented library functions
Most of the Swiss Ephemeris functions are implemented and exposed as functions
that can be called as RPC.

The following functions that are unsupported via RPC:
- `swe_get_library_path`
- `swe_set_astro_models`
- `swe_get_astro_models`
- `swe_set_interpolate_nut`
- `swe_csnorm`
- `swe_difcsn`
- `swe_difcs2n`
- `swe_cs2timestr`
- `swe_cs2lonlatstr`
- `swe_cs2degstr`

Some Swiss Ephemeris functions depend on global state within the library. This
state is modified via functions like `swe_set_jpl_file`. These functions are
exposed as RPC functions and could be called as a _context call_. Each RPC call
has a context with zero or more context calls. These context calls are executed
before the actual library function is called. Clients may call these functions
directly, however they are discouraged to do so as they may unable to
predictably target a specific worker or a set of workers that will handle the
call.

The following functions modify the internal state of the library:
- `swe_close`
- `swe_set_ephe_path`
- `swe_set_jpl_file`
- `swe_set_topo`
- `swe_set_sid_mode`
- `swe_set_lapse_rate`
- `swe_set_tid_acc`
- `swe_set_delta_t_userdef`

## Wire format
Requests and responses are serialized in the [MessagePack][msgpack] format.

[msgpack]: https://msgpack.org/

## Stateful functions in Swiss Ephemeris
Some Swiss Ephemeris functions depend on global state within the library. This
state is modified via functions like `swe_set_jpl_file`. These library
functions are exposed as RPC functions. To bridge the gap between the Swiss
Ephemeris and world of RPCs each RPC call has a context that allows to do
_context calls_ that will modify the internal state. All context calls will be
executed before the requested function that is called via RPC is executed.

## Worker implementation
The RPC function handlers are internally defined as an array for quick
dispatch. This means RPC functions are identified by the index in this array of
handlers. The order of the handlers follows `swephexp.h` except for the
heliacal functions, as they are defined after `swe_split_deg` followed by
`swe_difdegn`. This means the index identifying a function is depended on the
version of the library and is unstable.

That said, there is an exception: the first entry in the handlers array (index
0). The first entry must always a RPC function called `rpc_funcs`. It returns
an array that contains the name of each function in the internal array of
handlers. This allows to create a mapping between function name and index in
the handler array.

## RPC protocol
### Request
A request is an array that contains three values:
- Context: an `array` of zero or more context calls, may be `nil`.
- Function: the function index encoded as an `uint8_t`.
- Arguments: an `array` that matches the types of the library function, may be
`nil` if there are no arguments.

A context call is the same as request array except the `array` contains only a
function and arguments, no context.

A request that calls `swe_calc` with a couple of context calls looks like this:
```
00000000  36 30 3c 93 93 92 0d 93  cb 40 14 77 77 8d d6 16  |60<......@.ww...|
00000010  f8 cb 40 4a 0a aa a7 de  d6 bb 00 92 0e 93 01 00  |..@J............|
00000020  00 92 0b 91 a9 64 65 34  33 31 2e 65 70 68 04 93  |.....de431.eph..|
00000030  cb 41 42 72 8e d4 80 f1  2c 00 ce 00 01 81 01 3e  |.ABr....,......>|
```

It corresponds to the following JSON:
```json
[
  [
    [13, [5.116667, 52.083333, 0]],
    [14, [1, 0, 0]],
    [11, ["de431.eph"]]
  ],
  4,
  [2417949.660185, 0, 98561]
]
```

And corresponds to the following C code:
```c
swe_set_topo(5.116667, 52.083333, 0);
swe_set_sid_mode(SE_SIDM_LAHIRI, 0, 0);
swe_set_jpl_file("de431.eph");

int fl = SEFLG_JPLEPH | SEFLG_SPEED | SEFLG_TOPOCTR | SEFLG_SIDEREAL;
double xx[6];
char err[AS_MAXCH];
int rc = swe_calc(2417949.660185, SE_SUN, fl, xx, err);
```

### Response
A response can either be an array with return values or an error string.

The response for the example request looks like this:
```
00000000  36 32 3c 93 ce 00 01 81  41 96 cb 40 70 8e e8 b7  |62<.....A..@p...|
00000010  e1 56 05 cb bf 5e 72 64  be 9d 1b 96 cb 3f ef 77  |.V...^rd.....?.w|
00000020  c2 d6 24 d2 ed cb 3f f0  61 6f 8d 18 b0 00 cb bf  |..$...?.ao......|
00000030  6e d9 f9 96 82 f0 cc cb  bf 1c bc 80 24 1e 00 00  |n...........$...|
00000040  a0 3e                                             |.>|
```

It corresponds to the following JSON:
```json
[
  98625,
  [
    264.931816,
    -0.001858,
    0.983369,
    1.023788,
    -0.003766,
    -0.000110
  ],
  ""
]
```

And corresponds to the following values in C:
```c
rc = SEFLG_JPLEPH | SEFLG_NONUT | SEFLG_SPEED | SEFLG_TOPOCTR | SEFLG_SIDEREAL;
xx = [264.931816, -0.001858, 0.983369, 1.023788, -0.003766, -0.000110];
err = "";
```

## Implementation specifics
### swedenw-stdio
This program is designed to run as subprocess of the client. It will read
requests from stdin and write responses to stdout. The process will exit when a
newline (`\n`) is read from stdin. Each request and response is framed in a
[Lich][lich] data element so the logic in the worker is fairly simple. Framing
enables reading the input into a buffer and parse the request incrementally.

The client is able to call functions that change the Swiss Ephemeris library
state. With this ability comes the responsibility for the client to initialize
the worker properly by calling `swe_set_ephe_path` on start up and `swe_close`
before quitting the process.

[lich]: https://github.com/rentzsch/lich
