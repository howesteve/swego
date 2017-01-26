#include <stdbool.h>
#include <stdint.h>
#include "msgpuck.h"

#define REQSIZE 1024
#define RESPSIZE 1024
#define DBGSIZE 512

void tr_init(int argc, char const *argv[]);
char *tr_recv(char *buf);
bool tr_send(char *data, char *end);
void tr_error(const char *msg, const char *dbg, size_t dbglen);
