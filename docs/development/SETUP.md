# XShare Development Environment Setup

## Quick Setup (China)

使用阿里云镜像加速下载：

```bash
# 安装 mise
curl https://mise.run | sh

# 配置阿里云镜像
export MISE_GO_DOWNLOAD_MIRROR=https://mirrors.aliyun.com/golang
export MISE_NPM_DOWNLOAD_MIRROR=https://registry.npmmirror.com

# 安装 Go
mise use -g go@1.22

# 安装 Gradle (使用清华镜像)
echo 'export GRADLE_DOWNLOAD_MIRROR=https://mirrors.cloud.tencent.com/gradle' >> ~/.bashrc
```

## Go Core

### 安装 Go

```bash
# 使用 mise (推荐)
mise use -g go@1.22

# 或使用阿里云镜像手动下载
wget https://mirrors.aliyun.com/golang/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### 验证

```bash
cd core/go
go test ./...
go build ./...
```

## Android App

### 安装 Android SDK

```bash
# 使用 Android Studio (推荐)
# 下载：https://developer.android.google.cn/studio

# 或使用命令行工具
wget https://dl.google.com/android/repository/commandlinetools-linux-*.zip
mkdir -p $ANDROID_HOME/cmdline-tools
unzip commandlinetools-*.zip -d $ANDROID_HOME/cmdline-tools
mv $ANDROID_HOME/cmdline-tools/cmdline-tools latest

# 安装 SDK 组件
$ANDROID_HOME/cmdline-tools/latest/bin/sdkmanager "platform-tools" "platforms;android-34" "build-tools;34.0.0"
```

### 安装 Gradle

```bash
# 使用 mise
mise use -g gradle@8.5

# 或使用腾讯镜像
wget https://mirrors.cloud.tencent.com/gradle/gradle-8.5-bin.zip
sudo unzip -d /opt gradle-8.5-bin.zip
export PATH=$PATH:/opt/gradle-8.5/bin
```

### 验证

```bash
cd android
./gradlew :app:testDebugUnitTest
./gradlew assembleDebug
```

## ESP32 Firmware

### 安装 ESP-IDF

```bash
# 使用官方安装脚本 (中国用户建议使用镜像)
git clone -b v5.1 --depth 1 https://github.com/espressif/esp-idf.git
cd esp-idf
./install.sh esp32

# 设置环境变量
. ./export.sh

# 验证安装
idf.py --version
```

### 使用 Docker (推荐)

```bash
# 拉取 ESP-IDF Docker 镜像
docker pull espressif/idf

# 运行构建
docker run --rm -v $PWD:/project -w /project espressif/idf idf.py build
```

### 验证

```bash
cd firmware/esp32
idf.py build
```

## 完整验证

```bash
# 运行验证脚本
bash tools/verify-mvp.sh
```

## 环境变量配置

将以下内容添加到 `~/.bashrc` 或 `~/.zshrc`：

```bash
# Go
export PATH=$PATH:/usr/local/go/bin

# Android SDK
export ANDROID_HOME=$HOME/Android/Sdk
export PATH=$PATH:$ANDROID_HOME/tools:$ANDROID_HOME/platform-tools

# ESP-IDF
export IDF_PATH=$HOME/esp-idf
. $IDF_PATH/export.sh
```

## 国内镜像源

| 资源 | 镜像地址 |
|------|---------|
| Go | https://mirrors.aliyun.com/golang |
| npm | https://registry.npmmirror.com |
| PyPI | https://mirrors.aliyun.com/pypi |
| Gradle | https://mirrors.cloud.tencent.com/gradle |
| Maven | https://maven.aliyun.com/repository/public |
| Docker | https://registry.docker-cn.com |
