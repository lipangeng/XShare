#include "softap_mgr.h"
#include "esp_wifi.h"
#include "esp_netif.h"
#include "esp_log.h"
#include "lwip/ip_addr.h"
#include <string.h>

static const char *TAG = "softap_mgr";

static softap_mgr_state_t s_state = SOFTAP_MGR_STATE_IDLE;
static softap_mgr_config_t s_config = {
    .ssid = SOFTAP_MGR_SSID_DEFAULT,
    .password = SOFTAP_MGR_PASSWORD_DEFAULT,
    .gateway_ip = SOFTAP_MGR_IP_DEFAULT,
    .subnet_mask = SOFTAP_MGR_SUBNET_DEFAULT,
};

static client_connected_cb_t s_connect_cb = NULL;
static client_disconnected_cb_t s_disconnect_cb = NULL;

esp_err_t softap_mgr_start(void)
{
    if (s_state == SOFTAP_MGR_STATE_STARTED) {
        return ESP_OK;
    }

    ESP_LOGI(TAG, "Starting SoftAP with SSID: %s", s_config.ssid);

    wifi_config_t wifi_config = {
        .ap = {
            .ssid_len = 0,
            .max_connection = 4,
            .authmode = WIFI_AUTH_OPEN,
        },
    };

    strncpy((char *)wifi_config.ap.ssid, s_config.ssid, sizeof(wifi_config.ap.ssid));
    strncpy((char *)wifi_config.ap.password, s_config.password, sizeof(wifi_config.ap.password));

    if (strlen(s_config.password) == 0) {
        wifi_config.ap.authmode = WIFI_AUTH_OPEN;
    } else {
        wifi_config.ap.authmode = WIFI_AUTH_WPA_WPA2_PSK;
    }

    ESP_ERROR_CHECK(esp_wifi_set_mode(WIFI_MODE_AP));
    ESP_ERROR_CHECK(esp_wifi_set_config(WIFI_IF_AP, &wifi_config));
    ESP_ERROR_CHECK(esp_wifi_start());

    s_state = SOFTAP_MGR_STATE_STARTED;
    ESP_LOGI(TAG, "SoftAP started");

    return ESP_OK;
}

esp_err_t softap_mgr_stop(void)
{
    if (s_state != SOFTAP_MGR_STATE_STARTED) {
        return ESP_OK;
    }

    ESP_ERROR_CHECK(esp_wifi_stop());
    s_state = SOFTAP_MGR_STATE_IDLE;
    ESP_LOGI(TAG, "SoftAP stopped");

    return ESP_OK;
}

esp_err_t softap_mgr_set_config(const softap_mgr_config_t *config)
{
    if (config == NULL) {
        return ESP_ERR_INVALID_ARG;
    }

    memcpy(&s_config, config, sizeof(softap_mgr_config_t));
    return ESP_OK;
}

esp_err_t softap_mgr_get_config(softap_mgr_config_t *config)
{
    if (config == NULL) {
        return ESP_ERR_INVALID_ARG;
    }

    memcpy(config, &s_config, sizeof(softap_mgr_config_t));
    return ESP_OK;
}

softap_mgr_state_t softap_mgr_get_state(void)
{
    return s_state;
}

esp_err_t softap_mgr_get_client_count(int *count)
{
    if (count == NULL) {
        return ESP_ERR_INVALID_ARG;
    }

    wifi_sta_list_t sta_list;
    esp_err_t err = esp_wifi_sta_get_list(&sta_list);
    if (err != ESP_OK) {
        return err;
    }

    *count = sta_list.num;
    return ESP_OK;
}

void softap_mgr_register_connect_callback(client_connected_cb_t cb)
{
    s_connect_cb = cb;
}

void softap_mgr_register_disconnect_callback(client_disconnected_cb_t cb)
{
    s_disconnect_cb = cb;
}
