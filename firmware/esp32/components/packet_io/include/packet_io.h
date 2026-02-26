#ifndef XSHARE_PACKET_IO_H
#define XSHARE_PACKET_IO_H

#include "esp_err.h"

typedef struct {
  void *reserved;
} packet_io_config_t;

esp_err_t packet_io_init(const packet_io_config_t *config);
esp_err_t packet_io_start(const packet_io_config_t *config);

#endif  // XSHARE_PACKET_IO_H
