#include <stddef.h>
#include <stdbool.h>

bool handlers_test_functions_enabled;

// Buffer resp is NULL if called as context call.
typedef char *(*handler_callback_t)(char *resp, const char **req);

typedef struct handler handler_t;
struct handler {
  char *name;
  size_t argc;
  bool ccall;
  handler_callback_t callback;
};

void handlers_init();
size_t handlers_count();
handler_t *handlers_get(size_t idx);
