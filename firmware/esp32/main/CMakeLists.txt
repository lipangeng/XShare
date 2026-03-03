idf_component_register(
    SRCS "main.c"
    INCLUDE_DIRS "."
    REQUIRES softap_mgr packet_io usb_mux ctrl_agent
)
