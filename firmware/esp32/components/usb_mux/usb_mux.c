#include "usb_mux.h"
#include "esp_log.h"
#include <string.h>

static const char *TAG = "usb_mux";

static bool s_initialized = false;
static bool s_connected = false;
static usb_mux_frame_cb_t s_frame_cb = NULL;

esp_err_t usb_mux_init(void)
{
    if (s_initialized) {
        return ESP_OK;
    }

    s_initialized = true;
    ESP_LOGI(TAG, "USB MUX initialized");
    return ESP_OK;
}

esp_err_t usb_mux_deinit(void)
{
    if (!s_initialized) {
        return ESP_OK;
    }

    s_initialized = false;
    s_connected = false;
    ESP_LOGI(TAG, "USB MUX deinitialized");
    return ESP_OK;
}

esp_err_t usb_mux_send(usb_mux_channel_t channel, const uint8_t *data, size_t len)
{
    if (!s_initialized) {
        return ESP_ERR_INVALID_STATE;
    }

    if (data == NULL || len == 0) {
        return ESP_ERR_INVALID_ARG;
    }

    return ESP_OK;
}

esp_err_t usb_mux_receive(usb_mux_channel_t channel, uint8_t *data, size_t *len, uint32_t timeout_ms)
{
    if (!s_initialized) {
        return ESP_ERR_INVALID_STATE;
    }

    if (data == NULL || len == NULL) {
        return ESP_ERR_INVALID_ARG;
    }

    return ESP_ERR_TIMEOUT;
}

bool usb_mux_is_connected(void)
{
    return s_connected;
}

void usb_mux_register_frame_callback(usb_mux_frame_cb_t cb)
{
    s_frame_cb = cb;
}
