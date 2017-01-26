#if !defined(DEBUG) || DEBUG == 0
#define NDEBUG 1
#else
#undef NDEBUG
#endif

#include <stdlib.h>
#include <stdbool.h>
#include <stdio.h>
#include <string.h>

#include "tr.h"
#include "handlers.h"

int main(int argc, char const *argv[]) {
  handlers_init();
  tr_init(argc, argv);

  char req[REQSIZE];
  char resp[RESPSIZE];

  while (true) {
nextreq:
    memset(req, 0, REQSIZE);
    memset(resp, 0, RESPSIZE);
    const char *reqbuf = req;
    char *respbuf = resp;

    reqbuf = tr_recv(req);
    if (reqbuf == NULL) {
      goto nextreq;
    }

    const char *input = (char *)req;
    if (!mp_check(&input, reqbuf)) {
      tr_error("input is not valid msgpack", NULL, 0);
      goto nextreq;
    }

    // Reset pointer to start of request buffer.
    reqbuf = req;

    uint32_t fields = mp_decode_array(&reqbuf);
    if (fields != 3) {
      char dbg[DBGSIZE];
      size_t dbglen = 0;
#if DEBUG
      dbglen = sprintf(dbg, "size=%u", fields);
#endif
      tr_error("array with 3 values expected (envelope)", dbg, dbglen);
      goto nextreq;
    }

    // Execute context calls first.
    // The type of the context value is either array or nil.
    if (mp_typeof(*reqbuf) == MP_NIL) {
      mp_decode_nil(&reqbuf);
    } else {
      uint32_t size = mp_decode_array(&reqbuf);
      for (size_t i = 0; i < size; i++) {
        uint32_t fields = mp_decode_array(&reqbuf);
        if (fields != 2) {
          char dbg[DBGSIZE];
          size_t dbglen = 0;
#if DEBUG
          dbglen = sprintf(dbg, "size=%u", fields);
#endif
          tr_error("array with 2 values expected (ccall envelope)", dbg, dbglen);
          goto nextreq;
        }

        uint8_t idx = mp_load_u8(&reqbuf);
        handler_t *h = handlers_get(idx);
        if (h == NULL) {
          char dbg[DBGSIZE];
          size_t dbglen = 0;
#if DEBUG
          dbglen = sprintf(dbg, "func=%u", idx);
#endif
          tr_error("invalid index (ccall function)", dbg, dbglen);
          goto nextreq;
        }

        // The type of the arguments value is either array or nil.
        if (mp_typeof(*reqbuf) == MP_NIL) {
          mp_decode_nil(&reqbuf);
        } else {
          uint32_t argc = mp_decode_array(&reqbuf);
          if (h->argc != argc) {
            char dbg[DBGSIZE];
            size_t dbglen = 0;
#if DEBUG
            dbglen = sprintf(dbg, "func=%u(%s) argc=%zu/%u", idx, h->name, h->argc, argc);
#endif
            tr_error("invalid number of arguments (ccall function)", dbg, dbglen);
            goto nextreq;
          }
        }

        if (!h->ccall) {
          char dbg[DBGSIZE];
          size_t dbglen = 0;
#if DEBUG
            dbglen = sprintf(dbg, "func=%u(%s)", idx, h->name);
#endif
            tr_error("function is invalid as context call", dbg, dbglen);
        }

        h->callback(NULL, &reqbuf);
      }
    }

    // Execute actual call.
    uint8_t idx = mp_load_u8(&reqbuf);
    handler_t *h = handlers_get(idx);
    if (h == NULL) {
      char dbg[DBGSIZE];
      size_t dbglen = 0;
#if DEBUG
      dbglen = sprintf(dbg, "func=%u", idx);
#endif
      tr_error("invalid index (function)", dbg, dbglen);
      goto nextreq;
    }

    // The type of the arguments value is either array or nil.
    if (mp_typeof(*reqbuf) == MP_NIL) {
      // If the type is nil, invalidate the request buffer.
      reqbuf = NULL;
    } else {
      uint32_t argc = mp_decode_array(&reqbuf);
      if (h->argc != argc) {
        char dbg[DBGSIZE];
        size_t dbglen = 0;
#if DEBUG
        dbglen = sprintf(dbg, "func=%u(%s) argc=%zu/%u", idx, h->name, h->argc, argc);
#endif
        tr_error("invalid number of arguments", dbg, dbglen);
        goto nextreq;
      }
    }

    respbuf = h->callback(respbuf, &reqbuf);
    if (respbuf == NULL) {
      char dbg[DBGSIZE];
      size_t dbglen = 0;
#if DEBUG
      dbglen = sprintf(dbg, "func=%u(%s)", idx, h->name);
#endif
      tr_error("function call failed", dbg, dbglen);
      goto nextreq;
    }

    if (!tr_send(resp, respbuf)) {
      return EXIT_FAILURE;
    }
  }

  return EXIT_SUCCESS;
}
