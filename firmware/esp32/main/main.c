#include <stdio.h>
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "softap_mgr.h"
#include "packet_io.h"

void app_main(void)
{
    printf("XShare Gateway starting...\n");

    esp_err_t ret = softap_mgr_start();
    if (ret != ESP_OK) {
        printf("Failed to start SoftAP manager\n");
        return;
    }

    ret = packet_io_start(NULL);
    if (ret != ESP_OK) {
        printf("Failed to start packet IO\n");
        return;
    }

    printf("XShare Gateway running\n");
}
