cmake_minimum_required(VERSION 3.22.1)

project(corebridge LANGUAGES C CXX)

add_library(corebridge SHARED
    src/main/cpp/native_bridge.cpp
)

find_library(log-lib log)

target_link_libraries(corebridge
    ${log-lib}
)
