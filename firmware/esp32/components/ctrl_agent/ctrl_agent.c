#include "ctrl_agent.h"
#include "esp_log.h"

static const char *TAG = "ctrl_agent";

static ctrl_agent_state_t s_state = CTRL_AGENT_STATE_IDLE;
static ctrl_msg_cb_t s_msg_cb = NULL;

esp_err_t ctrl_agent_init(void)
{
    s_state = CTRL_AGENT_STATE_IDLE;
    ESP_LOGI(TAG, "Control agent initialized");
    return ESP_OK;
}

esp_err_t ctrl_agent_deinit(void)
{
    s_state = CTRL_AGENT_STATE_IDLE;
    s_msg_cb = NULL;
    ESP_LOGI(TAG, "Control agent deinitialized");
    return ESP_OK;
}

esp_err_t ctrl_agent_start(void)
{
    if (s_state != CTRL_AGENT_STATE_IDLE) {
        return ESP_OK;
    }

    s_state = CTRL_AGENT_STATE_CONNECTED;
    ESP_LOGI(TAG, "Control agent started");
    return ESP_OK;
}

esp_err_t ctrl_agent_stop(void)
{
    s_state = CTRL_AGENT_STATE_IDLE;
    ESP_LOGI(TAG, "Control agent stopped");
    return ESP_OK;
}

ctrl_agent_state_t ctrl_agent_get_state(void)
{
    return s_state;
}

esp_err_t ctrl_agent_send_response(const char *method, const uint8_t *data, size_t len)
{
    if (s_state == CTRL_AGENT_STATE_IDLE) {
        return ESP_ERR_INVALID_STATE;
    }

    return ESP_OK;
}

void ctrl_agent_register_message_callback(ctrl_msg_cb_t cb)
{
    s_msg_cb = cb;
}
