#include <stdbool.h>
#include <stdlib.h>

bool swex_supports_tls();
void swex_set_jpl_file(const char *fname);
void swex_set_jpl_file_len(const char *fname, size_t len);
void swex_set_topo(double geolon, double geolat, double geoalt);
void swex_set_sid_mode(int32_t sidm, double t0, double ayan_t0);
