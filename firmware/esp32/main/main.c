#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "esp_check.h"
#include "softap_mgr.h"
#include "packet_io.h"

void app_main(void) {
  ESP_ERROR_CHECK(softap_mgr_start());
  ESP_ERROR_CHECK(packet_io_start());
}
