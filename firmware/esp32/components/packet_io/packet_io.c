#include "packet_io.h"
#include "esp_log.h"
#include "esp_netif.h"
#include "esp_eth.h"
#include "esp_wifi.h"
#include <string.h>

static const char *TAG = "packet_io";

static bool s_running = false;
static uplink_cb_t s_uplink_cb = NULL;
static downlink_cb_t s_downlink_cb = NULL;

esp_err_t packet_io_start(uplink_cb_t cb)
{
    if (s_running) {
        return ESP_OK;
    }

    s_uplink_cb = cb;
    s_running = true;

    ESP_LOGI(TAG, "Packet IO started");
    return ESP_OK;
}

esp_err_t packet_io_stop(void)
{
    if (!s_running) {
        return ESP_OK;
    }

    s_running = false;
    s_uplink_cb = NULL;

    ESP_LOGI(TAG, "Packet IO stopped");
    return ESP_OK;
}

esp_err_t packet_io_inject_downlink(const uint8_t *pkt, size_t len)
{
    if (!s_running) {
        return ESP_ERR_INVALID_STATE;
    }

    if (pkt == NULL || len == 0) {
        return ESP_ERR_INVALID_ARG;
    }

    if (s_downlink_cb) {
        s_downlink_cb(pkt, len);
    }

    return ESP_OK;
}

bool packet_io_is_running(void)
{
    return s_running;
}

void packet_io_register_uplink_callback(uplink_cb_t cb)
{
    s_uplink_cb = cb;
}

void packet_io_register_downlink_callback(downlink_cb_t cb)
{
    s_downlink_cb = cb;
}
