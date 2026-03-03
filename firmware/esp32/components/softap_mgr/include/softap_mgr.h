#ifndef SOFTAP_MGR_H
#define SOFTAP_MGR_H

#include "esp_err.h"
#include <stdint.h>
#include <stdbool.h>

#define SOFTAP_MGR_SSID_DEFAULT "XShare_Gateway"
#define SOFTAP_MGR_PASSWORD_DEFAULT ""
#define SOFTAP_MGR_IP_DEFAULT "192.168.4.1"
#define SOFTAP_MGR_SUBNET_DEFAULT "255.255.255.0"

typedef enum {
    SOFTAP_MGR_STATE_IDLE = 0,
    SOFTAP_MGR_STATE_STARTED,
    SOFTAP_MGR_STATE_ERROR
} softap_mgr_state_t;

typedef struct {
    char ssid[32];
    char password[64];
    char gateway_ip[16];
    char subnet_mask[16];
} softap_mgr_config_t;

typedef void (*client_connected_cb_t)(const char *mac, const char *ip);
typedef void (*client_disconnected_cb_t)(const char *mac, const char *ip);

esp_err_t softap_mgr_start(void);
esp_err_t softap_mgr_stop(void);
esp_err_t softap_mgr_set_config(const softap_mgr_config_t *config);
esp_err_t softap_mgr_get_config(softap_mgr_config_t *config);
softap_mgr_state_t softap_mgr_get_state(void);
esp_err_t softap_mgr_get_client_count(int *count);

void softap_mgr_register_connect_callback(client_connected_cb_t cb);
void softap_mgr_register_disconnect_callback(client_disconnected_cb_t cb);

#endif
