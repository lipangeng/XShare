idf_component_register(
    SRCS "softap_mgr.c"
    INCLUDE_DIRS "include"
    REQUIRES esp_wifi esp_netif lwip
)
