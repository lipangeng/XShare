---
name: esp32-s3
description: |
  ESP32-S3 资料检索与开发环境安装“主线规则”。当用户提到 ESP32-S3、ESP-IDF、Arduino-ESP32、PlatformIO、烧录/编译/USB/JTAG/PSRAM/Boot 等时，优先加载并遵循本 skill。
---

# ESP32-S3：资料与安装环境遵循规则（OpenCode Skill）

## 你要做的事
当话题涉及 ESP32-S3：
1) 明确“权威资料”从哪里找
2) 明确“安装开发环境”必须遵循哪份文档
3) 遇到分支（ESP-IDF / Arduino / PlatformIO）时按规则选择，并给出对应官方入口

## 资料来源优先级（从高到低）
### A. Espressif 官方（永远第一优先）
- ESP-IDF Programming Guide（目标：ESP32-S3，版本：stable）
    - https://docs.espressif.com/projects/esp-idf/en/stable/esp32s3/
- ESP32-S3 Get Started（安装/配置/编译/烧录主线）
    - https://docs.espressif.com/projects/esp-idf/en/stable/esp32s3/get-started/index.html
- ESP-IDF GitHub（示例、源码、版本信息、issue）
    - https://github.com/espressif/esp-idf

### B. Arduino（仅当用户明确要 Arduino 才走）
- Arduino-ESP32 官方安装文档
    - https://docs.espressif.com/projects/arduino-esp32/en/latest/installing.html
- Arduino-ESP32 入门
    - https://docs.espressif.com/projects/arduino-esp32/en/latest/getting_started.html

### C. PlatformIO（仅当用户明确要 PlatformIO 才走）
- PlatformIO ESP-IDF 框架文档
    - https://docs.platformio.org/en/latest/frameworks/espidf.html
- Espressif 对第三方 PlatformIO 的说明（提示：第三方工具，问题先查 PlatformIO）
    - https://docs.espressif.com/projects/esp-idf/en/stable/esp32/third-party-tools/platformio.html

### D. 板卡厂商资料（只用于“板级差异”）
仅用于：引脚图、USB/串口芯片、Boot/Reset 按键、下载模式、PSRAM 配置提示等。
不要用厂商“快速开始”替代 Espressif 的 ESP-IDF 安装主线。

## 环境安装默认主线（除非用户指定）
默认：ESP-IDF（stable）+ ESP32-S3 Get Started
- 任何“怎么安装/怎么配环境/怎么编译/怎么烧录”的问题：
    - 一律先指向并遵循 ESP32-S3 Get Started（stable）
- 用户若拿来的是旧教程/旧命令：
    - 先对齐到 stable 文档
    - 只有当用户明确要求或项目锁定旧版本时，才切到对应版本的 ESP-IDF 文档页（例如 v5.1/v4.4 等）

## 决策规则（用户没说清楚时怎么选）
1) 量产/稳定/系统能力（Wi-Fi/BLE/USB/低功耗/PSRAM/性能/组件）=> ESP-IDF
2) 快速 demo / 依赖 Arduino 库生态 => Arduino-ESP32
3) 需要跨平台统一工程管理（且用户已在用 PlatformIO）=> PlatformIO（并提示它是第三方维护链路）

## 每次输出必须包含（强制）
- 你选择的主线：ESP-IDF / Arduino / PlatformIO
- 对应“官方入口链接”（至少 1 条）
- 若涉及安装/构建/烧录：明确说明“以该官方文档步骤为准”
- 若出现冲突：以 Espressif 官方文档为准；第三方教程仅作补充

## 常见坑（简短提醒）
- PlatformIO 的 ESP-IDF 链路是第三方，遇到构建/包管理异常优先查 PlatformIO 文档与社区
- Arduino core 与 ESP-IDF 版本有关联：用户要 IDF 新特性时先核对 Arduino-ESP32 对应的 IDF 版本