#ifndef PACKET_IO_H
#define PACKET_IO_H

#include "esp_err.h"
#include <stdint.h>
#include <stdbool.h>

typedef void (*uplink_cb_t)(const uint8_t *pkt, size_t len);
typedef void (*downlink_cb_t)(const uint8_t *pkt, size_t len);

esp_err_t packet_io_start(uplink_cb_t cb);
esp_err_t packet_io_stop(void);
esp_err_t packet_io_inject_downlink(const uint8_t *pkt, size_t len);
bool packet_io_is_running(void);

void packet_io_register_uplink_callback(uplink_cb_t cb);
void packet_io_register_downlink_callback(downlink_cb_t cb);

#endif
