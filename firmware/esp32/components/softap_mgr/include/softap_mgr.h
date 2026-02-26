#ifndef XSHARE_SOFTAP_MGR_H
#define XSHARE_SOFTAP_MGR_H

#include "esp_err.h"

typedef struct {
  void *reserved;
} softap_mgr_config_t;

esp_err_t softap_mgr_init(const softap_mgr_config_t *config);
esp_err_t softap_mgr_start(const softap_mgr_config_t *config);

#endif  // XSHARE_SOFTAP_MGR_H
