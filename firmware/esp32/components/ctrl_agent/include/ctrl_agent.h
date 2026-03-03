#ifndef CTRL_AGENT_H
#define CTRL_AGENT_H

#include "esp_err.h"
#include <stdint.h>

typedef enum {
    CTRL_AGENT_STATE_IDLE = 0,
    CTRL_AGENT_STATE_CONNECTED,
    CTRL_AGENT_STATE_FORWARDING
} ctrl_agent_state_t;

typedef void (*ctrl_msg_cb_t)(const char *method, const uint8_t *data, size_t len);

esp_err_t ctrl_agent_init(void);
esp_err_t ctrl_agent_deinit(void);
esp_err_t ctrl_agent_start(void);
esp_err_t ctrl_agent_stop(void);
ctrl_agent_state_t ctrl_agent_get_state(void);

esp_err_t ctrl_agent_send_response(const char *method, const uint8_t *data, size_t len);

void ctrl_agent_register_message_callback(ctrl_msg_cb_t cb);

#endif
