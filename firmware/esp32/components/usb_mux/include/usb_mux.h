#ifndef USB_MUX_H
#define USB_MUX_H

#include "esp_err.h"
#include <stdint.h>
#include <stdbool.h>

typedef enum {
    USB_MUX_CHANNEL_CONTROL = 1,
    USB_MUX_CHANNEL_DATA = 2,
    USB_MUX_CHANNEL_OTA = 3
} usb_mux_channel_t;

typedef void (*usb_mux_frame_cb_t)(usb_mux_channel_t channel, const uint8_t *data, size_t len);

esp_err_t usb_mux_init(void);
esp_err_t usb_mux_deinit(void);
esp_err_t usb_mux_send(usb_mux_channel_t channel, const uint8_t *data, size_t len);
esp_err_t usb_mux_receive(usb_mux_channel_t channel, uint8_t *data, size_t *len, uint32_t timeout_ms);
bool usb_mux_is_connected(void);

void usb_mux_register_frame_callback(usb_mux_frame_cb_t cb);

#endif
