cmake_minimum_required(VERSION 3.4.1)

project("xshare_core")

add_library(xshare_core SHARED
    native_bridge.cpp
)

find_library(log-lib log)

target_link_libraries(xshare_core ${log-lib})
