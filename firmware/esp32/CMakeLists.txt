cmake_minimum_required(VERSION 3.16)

if(NOT DEFINED ENV{IDF_PATH} OR "$ENV{IDF_PATH}" STREQUAL "")
  message(FATAL_ERROR "IDF_PATH is not set. Export IDF_PATH to your ESP-IDF installation directory before configuring.")
endif()

include($ENV{IDF_PATH}/tools/cmake/project.cmake)
project(xshare_esp32)
