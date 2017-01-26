#include <swephexp.h>
#include <sweph.h>
#include "sweversion.h"

#include "tr.h"
#include "handlers.h"

bool handlers_test_functions_enabled = false;
static handler_t handlers[];

void handlers_init() {
  // 1.1 swe_set_ephe_path()
  // This is the first function that should be called before any other function
  // of the Swiss Ephemeris. Even if you don’t want to set an ephemeris path
  // and use the Moshier ephemeris, it is nevertheless recommended to call
  // swe_set_ephe_path(NULL), because this function makes important
  // initializations. If you don’t do that, the Swiss Ephemeris may work, but
  // the results may be not 100% consistent.
  swe_set_ephe_path(NULL);
}

static int64_t mp_get_int(const char **data) {
  switch (mp_typeof(**data)) {
    default:
    case MP_INT:    return mp_decode_int(data);
    case MP_UINT:   return (int64_t)mp_decode_uint(data);
    case MP_DOUBLE: return (int64_t)mp_decode_double(data);
    case MP_FLOAT:  return (int64_t)mp_decode_float(data);
  }
}

static double mp_get_double(const char **data) {
  switch (mp_typeof(**data)) {
    default:
    case MP_DOUBLE: return mp_decode_double(data);
    case MP_FLOAT:  return (double)mp_decode_float(data);
    case MP_UINT:   return (double)mp_decode_uint(data);
    case MP_INT:    return (double)mp_decode_int(data);
  }
}

static char *mp_put_int(char *data, int64_t num) {
  if (num < 0) {
    return mp_encode_int(data, num);
  }

  return mp_encode_uint(data, (uint64_t)num);
}

static char *mp_put_str(char *data, const char *str) {
  return mp_encode_str(data, str, strlen(str));
}

static char *h_rpc_funcs(char *resp, __unused const char **req) {
  size_t n = handlers_count();
  resp = mp_encode_array(resp, n);

  for (size_t i = 0; i < n; i++) {
    resp = mp_put_str(resp, handlers[i].name);
  }

  return resp;
}

static char *h_test_crash(char *resp, __unused const char **req) {
  if (!handlers_test_functions_enabled) {
    resp = mp_encode_map(resp, 1);
    resp = mp_encode_str(resp, "err", 3);
    resp = mp_encode_str(resp, "function disabled", 17);
    return resp;
  }

  fprintf(stderr, "DEBUG: func=test_crash\n");
  fprintf(stderr, "ERROR: test_crash called\n");
  exit(EXIT_FAILURE);
}

static char *h_test_error(char *resp, __unused const char **req) {
  if (!handlers_test_functions_enabled) {
    resp = mp_encode_map(resp, 1);
    resp = mp_encode_str(resp, "err", 3);
    resp = mp_encode_str(resp, "function disabled", 17);
  } else {
    resp = mp_encode_map(resp, 2);
    resp = mp_encode_str(resp, "err", 3);
    resp = mp_encode_str(resp, "test_error called", 17);
    resp = mp_encode_str(resp, "dbg", 3);
    resp = mp_encode_str(resp, "func=test_error", 15);
  }

  return resp;
}

static char *h_swe_version(char *resp, __unused const char **req) {
  resp = mp_encode_array(resp, 1);
  resp = mp_put_str(resp, SE_VERSION);
  return resp;
}

typedef int32 (* swe_calc_func)(double, int32_t, int32_t, double *, char *);
static char *hf_swe_calc(char *resp, const char **req, swe_calc_func calc) {
  double jd = mp_get_double(req);
  int32_t pl = (int32_t)mp_get_int(req);
  int32_t fl = (int32_t)mp_get_int(req);

  double xx[6];
  char err[AS_MAXCH] = {0};
  int32_t rv = calc(jd, pl, fl, xx, err);

  resp = mp_encode_array(resp, 3);
  resp = mp_put_int(resp, rv);
  resp = mp_encode_array(resp, 6);
  for (size_t i = 0; i < 6; i++) {
    resp = mp_encode_double(resp, xx[i]);
  }
  resp = mp_put_str(resp, err);
  return resp;
}

static char *h_swe_calc(char *resp, const char **req) {
  return hf_swe_calc(resp, req, swe_calc);
}

static char *h_swe_calc_ut(char *resp, const char **req) {
  return hf_swe_calc(resp, req, swe_calc_ut);
}

typedef int32 (* swe_fixstar_func)(char *, double, int32, double *, char *);
static char *hf_swe_fixstar(char *resp, const char **req, swe_fixstar_func calc) {
  char star[41] = {0};
  uint32_t len = 0;
  const char *name = mp_decode_str(req, &len);
  strncpy(star, name, len);

  double jd = mp_get_double(req);
  int32_t fl = (int32_t)mp_get_int(req);

  double xx[6];
  char err[AS_MAXCH] = {0};
  int32_t rv = calc((char *)star, jd, fl, xx, err);

  resp = mp_encode_array(resp, 4);
  resp = mp_put_str(resp, star);
  resp = mp_put_int(resp, rv);
  resp = mp_encode_array(resp, 6);
  for (size_t i = 0; i < 6; i++) {
    resp = mp_encode_double(resp, xx[i]);
  }
  resp = mp_put_str(resp, err);
  return resp;
}

static char *h_swe_fixstar(char *resp, const char **req) {
  return hf_swe_fixstar(resp, req, swe_fixstar);
}

static char *h_swe_fixstar_ut(char *resp, const char **req) {
  return hf_swe_fixstar(resp, req, swe_fixstar_ut);
}

static char *h_swe_fixstar_mag(char *resp, const char **req) {
  char star[41] = {0};
  uint32_t len = 0;
  const char *name = mp_decode_str(req, &len);
  strncpy(star, name, len);

  double mag;
  char err[AS_MAXCH] = {0};
  int32_t rv = swe_fixstar_mag((char *)star, &mag, err);

  resp = mp_encode_array(resp, 4);
  resp = mp_put_str(resp, star);
  resp = mp_put_int(resp, rv);
  resp = mp_encode_double(resp, mag);
  resp = mp_put_str(resp, err);
  return resp;
}

static char *h_swe_close(char *resp, __unused const char **req) {
  swe_close();
  resp = mp_encode_array(resp, 0);
  return resp;
}

static char *h_swe_set_ephe_path(char *resp, const char **req) {
  uint32_t len = 0;
  const char *path = mp_decode_str(req, &len);

  swe_set_ephe_path((char *)path);

  resp = mp_encode_array(resp, 0);
  return resp;
}

static void swex_set_jpl_file(const char *fname, size_t len) {
  if (strncmp(fname, swed.jplfnam, len) != 0) {
    swe_set_jpl_file((char *)fname);
  }
}

static char *h_swe_set_jpl_file(char *resp, const char **req) {
  uint32_t len = 0;
  const char *fname = mp_decode_str(req, &len);

  swex_set_jpl_file(fname, len);

  if (resp == NULL) {
    return NULL;
  }

  resp = mp_encode_array(resp, 0);
  return resp;
}

static char *h_swe_get_planet_name(char *resp, const char **req) {
  int32_t pl = (int32_t)mp_get_int(req);

  char name[AS_MAXCH] = {0};
  swe_get_planet_name(pl, name);

  resp = mp_encode_array(resp, 1);
  resp = mp_put_str(resp, name);
  return resp;
}

static void swex_set_topo(double geolon, double geolat, double geoalt) {
#if SWEX_VERSION_MAJOR == 2 && SWEX_VERSION_MINOR < 5
  if (swed.geopos_is_set == TRUE
    && swed.topd.geolon == geolon
    && swed.topd.geolat == geolat
    && swed.topd.geoalt == geoalt
  ) {
    return;
  }
#endif

  swe_set_topo(geolon, geolat, geoalt);
}

static char *h_swe_set_topo(char *resp, const char **req) {
  double geolon = mp_get_double(req);
  double geolat = mp_get_double(req);
  double geoalt = mp_get_double(req);

  swex_set_topo(geolon, geolat, geoalt);

  if (resp == NULL) {
    return NULL;
  }

  resp = mp_encode_array(resp, 0);
  return resp;
}

static void swex_set_sid_mode(int32_t sidm, double t0, double ayan_t0) {
  if (swed.ayana_is_set == FALSE
    || swed.sidd.sid_mode != sidm
    || swed.sidd.ayan_t0 != ayan_t0
    || swed.sidd.t0 != t0
  ) {
    swe_set_sid_mode(sidm, t0, ayan_t0);
  }
}

static char *h_swe_set_sid_mode(char *resp, const char **req) {
  int32_t sidm = (int32_t)mp_get_int(req);
  double t0 = mp_get_double(req);
  double ayan_t0 = mp_get_double(req);

  swex_set_sid_mode(sidm, t0, ayan_t0);

  if (resp == NULL) {
    return NULL;
  }

  resp = mp_encode_array(resp, 0);
  return resp;
}

typedef int32 (* swe_get_ayanamsa_ex_func)(double, int32, double *, char *);
static char *hf_swe_get_ayanamsa_ex(char *resp, const char **req, swe_get_ayanamsa_ex_func calc) {
  double jd = mp_get_double(req);
  int32_t fl = (int32_t)mp_get_int(req);

  double aya = 0;
  char err[AS_MAXCH] = {0};
  int32_t rv = calc(jd, fl, &aya, err);

  resp = mp_encode_array(resp, 2);
  resp = mp_put_int(resp, rv);
  resp = mp_put_str(resp, err);
  return resp;
}

static char *h_swe_get_ayanamsa_ex(char *resp, const char **req) {
  return hf_swe_get_ayanamsa_ex(resp, req, swe_get_ayanamsa_ex);
}

static char *h_swe_get_ayanamsa_ex_ut(char *resp, const char **req) {
  return hf_swe_get_ayanamsa_ex(resp, req, swe_get_ayanamsa_ex_ut);
}

typedef double (* swe_get_ayanamsa_func)(double);
static char *hf_swe_get_ayanamsa(char *resp, const char **req, swe_get_ayanamsa_func calc) {
  double jd = mp_get_double(req);

  double aya = calc(jd);

  resp = mp_encode_array(resp, 1);
  resp = mp_encode_double(resp, aya);
  return resp;
}

static char *h_swe_get_ayanamsa(char *resp, const char **req) {
  return hf_swe_get_ayanamsa(resp, req, swe_get_ayanamsa);
}

static char *h_swe_get_ayanamsa_ut(char *resp, const char **req) {
  return hf_swe_get_ayanamsa(resp, req, swe_get_ayanamsa_ut);
}

static char *h_swe_get_ayanamsa_name(char *resp, const char **req) {
  int32_t sidm = (int32_t)mp_get_int(req);

  const char *name = swe_get_ayanamsa_name(sidm);

  resp = mp_encode_array(resp, 1);
  resp = mp_put_str(resp, name);
  return resp;
}

// swe_date_conversion
// swe_julday
// swe_revjul
// swe_utc_to_jd
// swe_jdet_to_utc
// swe_jdut1_to_utc
// swe_utc_time_zone
// swe_houses
// swe_houses_ex
// swe_houses_armc
// swe_house_pos
// swe_house_name
// swe_gauquelin_sector
// swe_sol_eclipse_where
// swe_lun_occult_where
// swe_sol_eclipse_how
// swe_sol_eclipse_when_loc
// swe_lun_occult_when_loc
// swe_sol_eclipse_when_glob
// swe_lun_occult_when_glob
// swe_lun_eclipse_how
// swe_lun_eclipse_when
// swe_lun_eclipse_when_loc
// swe_pheno
// swe_pheno_ut
// swe_refrac
// swe_refrac_extended
// swe_set_lapse_rate /* context */
// swe_azalt
// swe_azalt_rev
// swe_rise_trans_true_hor
// swe_rise_trans
// swe_nod_aps
// swe_nod_aps_ut
// swe_get_orbital_elements
// swe_orbit_max_min_true_distance
// swe_deltat
// swe_deltat_ex
// swe_time_equ
// swe_lmt_to_lat
// swe_lat_to_lmt
// swe_sidtime0
// swe_sidtime
// swe_set_interpolate_nut
// swe_cotrans
// swe_cotrans_sp
// swe_get_tid_acc
// swe_set_tid_acc /* context */
// swe_set_delta_t_userdef /* context */
// swe_degnorm
// swe_radnorm
// swe_rad_midp
// swe_deg_midp
// swe_split_deg
// swe_heliacal_ut
// swe_heliacal_pheno_ut
// swe_vis_limit_mag
// swe_heliacal_angle
// swe_topo_arcus_visionis
// swe_day_of_week

static handler_t handlers[] = {
  {"rpc_funcs",              0, false, h_rpc_funcs}, // keep this always on top!
  {"test_crash",             0, false, h_test_crash},
  {"test_error",             0, false, h_test_error},
  {"swe_version",            0, false, h_swe_version},
  {"swe_calc",               3, false, h_swe_calc},
  {"swe_calc_ut",            3, false, h_swe_calc_ut},
  {"swe_fixstar",            3, false, h_swe_fixstar},
  {"swe_fixstar_ut",         3, false, h_swe_fixstar_ut},
  {"swe_fixstar_mag",        1, false, h_swe_fixstar_mag},
  {"swe_close",              0, true,  h_swe_close},         /* context */
  {"swe_set_ephe_path",      1, true,  h_swe_set_ephe_path}, /* context */
  {"swe_set_jpl_file",       1, true,  h_swe_set_jpl_file},  /* context */
  {"swe_get_planet_name",    1, false, h_swe_get_planet_name},
  {"swe_set_topo",           3, true,  h_swe_set_topo},      /* context */
  {"swe_set_sid_mode",       3, true,  h_swe_set_sid_mode},  /* context */

#if SWEX_VERSION_MAJOR == 2 && SWEX_VERSION_MINOR >= 2
  {"swe_get_ayanamsa_ex",    2, false, h_swe_get_ayanamsa_ex},
  {"swe_get_ayanamsa_ex_ut", 2, false, h_swe_get_ayanamsa_ex_ut},
#endif

  {"swe_get_ayanamsa",       2, false, h_swe_get_ayanamsa},
  {"swe_get_ayanamsa_ut",    2, false, h_swe_get_ayanamsa_ut},
  {"swe_get_ayanamsa_name",  2, false, h_swe_get_ayanamsa_name},
  // swe_date_conversion
  // swe_julday
  // swe_revjul
  // swe_utc_to_jd
  // swe_jdet_to_utc
  // swe_jdut1_to_utc
  // swe_utc_time_zone
  // swe_houses
  // swe_houses_ex
  // swe_houses_armc
  // swe_house_pos
  // swe_house_name
  // swe_gauquelin_sector
  // swe_sol_eclipse_where
  // swe_lun_occult_where
  // swe_sol_eclipse_how
  // swe_sol_eclipse_when_loc
  // swe_lun_occult_when_loc
  // swe_sol_eclipse_when_glob
  // swe_lun_occult_when_glob
  // swe_lun_eclipse_how
  // swe_lun_eclipse_when
  // swe_lun_eclipse_when_loc
  // swe_pheno
  // swe_pheno_ut
  // swe_refrac
  // swe_refrac_extended
  // swe_set_lapse_rate /* context */
  // swe_azalt
  // swe_azalt_rev
  // swe_rise_trans_true_hor
  // swe_rise_trans
  // swe_nod_aps
  // swe_nod_aps_ut

#if SWEX_VERSION_MAJOR == 2 && SWEX_VERSION_MINOR >= 5
  // swe_get_orbital_elements
  // swe_orbit_max_min_true_distance
#endif

  // swe_deltat
  // swe_deltat_ex
  // swe_time_equ
  // swe_lmt_to_lat
  // swe_lat_to_lmt
  // swe_sidtime0
  // swe_sidtime

#if SWEX_VERSION_MAJOR == 2 && SWEX_VERSION_MINOR >= 6
  // swe_set_interpolate_nut /* context */
#endif

  // swe_cotrans
  // swe_cotrans_sp
  // swe_get_tid_acc
  // swe_set_tid_acc /* context */

#if SWEX_VERSION_MAJOR == 2 && SWEX_VERSION_MINOR >= 5
  // swe_set_delta_t_userdef /* context */
#endif

  // swe_degnorm
  // swe_radnorm
  // swe_rad_midp
  // swe_deg_midp
  // swe_split_deg
  // swe_heliacal_ut
  // swe_heliacal_pheno_ut
  // swe_vis_limit_mag
  // swe_heliacal_angle
  // swe_topo_arcus_visionis
  // swe_difdegn
  // swe_difdeg2n
  // swe_difrad2n
  // swe_d2l
  // swe_day_of_week
};

size_t handlers_count() {
  return sizeof(handlers) / sizeof(handler_t);
}

handler_t *handlers_get(size_t idx) {
  if (idx >= handlers_count()) {
    return NULL;
  }

  return &handlers[idx];
}
