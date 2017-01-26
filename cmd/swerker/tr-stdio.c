#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>

#include "tr.h"
#include "handlers.h"

void tr_init(int argc, char const *argv[]) {
  for (size_t i = 1; i < argc; i++) {
    if (strncmp(argv[i], "-w", 2) == 0) {
      sleep(5);
    }

    if (strncmp(argv[i], "-dangerous_enable_test_functions", 32) == 0) {
      handlers_test_functions_enabled = true;
    }

    if (strncmp(argv[i], "-dangerous_no_funcs_on_init", 27) == 0) {
      exit(EXIT_FAILURE);
    }

    if (strncmp(argv[i], "-dangerous_invalid_funcs_on_init", 32) == 0) {
      puts("invalid func data");
      exit(EXIT_FAILURE);
    }

    if (strncmp(argv[i], "-dangerous_invalid_funcs_types_on_init", 38) == 0) {
      fputs("1<\xc0>", stdout); // msgpack nil value
      exit(EXIT_FAILURE);
    }
  }

  setbuf(stdin, NULL);
  setvbuf(stdout, NULL, _IOFBF, RESPSIZE);

  // Write RPC functions, same as calling rpc_funcs function (index 0).
  char data[RESPSIZE];
  char *buf = data;
  handler_t *h = handlers_get(0);
  buf = h->callback(buf, NULL);

  if (!tr_send(data, buf)) {
    exit(EXIT_FAILURE);
  }
}

char *tr_recv(char *buf) {
  int c = fgetc(stdin);
  if (c == EOF || c == '\n') {
    exit(EXIT_SUCCESS);
  }

  uint64_t len = 0;
  while ('0' <= c && c <= '9') {
    len *= 10;
    len += c - '0';

    c = fgetc(stdin);
    if (c == EOF) {
      char dbg[DBGSIZE];
      size_t dbglen = 0;
#if DEBUG
      dbglen = sprintf(dbg, "len=%llu", len);
#endif
      tr_error("reading unexpected EOF (length)", dbg, dbglen);
      return NULL;
    }
  }

  // We limit input data to REQSIZE bytes to protect against buffer overflows
  // and unbounded buffer allocations.
  if (len > REQSIZE) {
    char dbg[DBGSIZE];
    size_t dbglen = 0;
#if DEBUG
    dbglen = sprintf(dbg, "len=%llu, limit=%d", len, REQSIZE);
#endif
    tr_error("input data is more than request size limit", dbg, dbglen);
    return NULL;
  }

  // char is already received in while loop
  if (c != '<') {
    char dbg[DBGSIZE];
    size_t dbglen = 0;
#if DEBUG
    dbglen = sprintf(dbg, "c='%c' c=%d", c, c);
#endif
    tr_error("reading unexpected open type marker", dbg, dbglen);
    return NULL;
  }

  size_t n = 0;
  while (n < len) {
    c = fgetc(stdin);
    if (c == EOF) {
      tr_error("reading unexpected EOF (body)", NULL, 0);
      return NULL;
    }

    buf[n++] = c;
  }

  c = fgetc(stdin);
  if (c != '>') {
    char dbg[DBGSIZE];
    size_t dbglen = 0;
#if DEBUG
    dbglen = sprintf(dbg, "c='%c' c=%d", c, c);
#endif
    tr_error("reading unexpected close type marker", dbg, dbglen);
    return NULL;
  }

  if (len == 0) {
    tr_error("input data expected", NULL, 0);
    return NULL;
  }

  return buf;
}

bool tr_send(char *data, char *end) {
  size_t len = end - data;

  int n = fprintf(stdout, "%lu<", len);
  fwrite(data, len, sizeof(char), stdout);
  putc('>', stdout);

  if (n < 0 || ferror(stdout)) {
    // Write to stdout is failed, so write error to stderr as last resort.
#if DEBUG
    fprintf(stderr, "DEBUG: len=%zu data=", len);
    fwrite(data, len, sizeof(char), stderr);
    putc('\n', stderr);
#endif
    fprintf(stderr, "ERROR: failed to write reponse\n");
    return false;
  }

  if (fflush(stdout) != 0) {
    // Flush to stdout is failed, so write error to stderr as last resort.
    fprintf(stderr, "ERROR: failed to flush reponse\n");
    return false;
  }

  return true;
}

void tr_error(const char *msg, const char *dbg, size_t dbglen) {
  char data[RESPSIZE];
  char *buf = data;

  int size = 1;
#if DEBUG
  if (dbglen != 0) {
    size = 2;
  }
#endif

  buf = mp_encode_map(buf, size);
  buf = mp_encode_str(buf, "err", 3);
  buf = mp_encode_str(buf, msg, strlen(msg));

#if DEBUG
  if (dbglen != 0) {
    buf = mp_encode_str(buf, "dbg", 3);
    buf = mp_encode_str(buf, dbg, dbglen);
  }
#endif

  tr_send(data, buf);
}
