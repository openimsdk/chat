# openim-chat

## 📄 源代码的许可选项

您可以在以下两种许可选项之一下使用 OpenIM 源代码来创建非 OpenIM 原始生产的编译版本：

### 1. 通用公共许可证 v3.0 (GPLv3) 🆓

+ 该选项受自由软件基金会的 [GPL v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html) 管理。
+ 使用受此政策概述的某些例外的约束。

### 2. 商业许可 💼

+ 通过联系 OpenIM 获得商业许可。
+ 有关详细信息和许可查询，请通过电子邮件 📧 [contact@openim.io](mailto:contact@openim.io)。

## 🧩 功能简介

1. 该仓库实现了业务系统，包括两部分：用户系统和后台管理系统。
2. 该系统依赖于 [open-im-server 仓库](https://github.com/openimsdk/open-im-server)，通过调用即时消息系统的 API 实现丰富的业务功能。
3. 用户系统包括一些常规功能，如用户登录、用户注册、用户信息更新等。
4. 后台管理系统包括提供了 API 管理用户、群组和消息等。

## :busts_in_silhouette: 社区

+ 💬 [关注我们的 Twitter 账户](https://twitter.com/founder_im63606)
+ 🚀 [加入我们的 Slack 社区](https://join.slack.com/t/openimsdk/shared_invite/zt-2hljfom5u-9ZuzP3NfEKW~BJKbpLm0Hw)
+ :eyes: [加入我们的微信群](https://openim-1253691595.cos.ap-nanjing.myqcloud.com/WechatIMG20.jpeg)

## 🛫 快速开始

> :warning: **注意**：本项目在 Linux/Windows/Mac 平台以及 ARM 和 AMD 架构下均可正常使用

### 📦 克隆

```bash
git clone https://github.com/openimsdk/chat openim-chat
cd openim-chat
```

### 🛠 初始化

:computer: 第一次编译前，Linux/Mac 平台下执行：

```
sh bootstrap.sh
```

:computer: Windows 执行：

```
bootstrap.bat
```

### 🏗 编译

```bash
mage
```

### 🚀 启动

```bash
mage start
```

### :floppy_disk: 或后台启动 收集日志

```
nohup mage start >> _output/logs/chat.log 2>&1 &
```

### :mag_right: 检测

```bash
mage check
```

### 🛑 停止

```bash
mage stop
```

### 🚀 启动顺序

1. 成功启动 [open-im-server](https://github.com/openimsdk/open-im-server)。
2. 编译 chat `mage`。
3. 启动 chat `mage start`。

## 📞 如果您想启用音视频通话，请配置 LiveKit

:link: 请参考 "[如何设置 LiveKit 服务器](./HOW_TO_SETUP_LIVEKIT_SERVER.md)"。

## :handshake: 贡献

:heart: 欢迎对该项目做出贡献！请查看 [CONTRIBUTING.md](./CONTRIBUTING.md) 了解详情。

## 🚨 许可

:scroll: chat 根据 [GPL-3.0 license](https://github.com/openimsdk/chat#GPL-3.0-1-ov-file) 许可证授权。查看 [LICENSE](https://github.com/openimsdk/chat/tree/main/LICENSE) 获取完整的许可证文本。
