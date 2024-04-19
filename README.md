# openim-chat

## 📄 License Options for OpenIM Source Code

You may use the OpenIM source code to create compiled versions not originally produced by OpenIM under one of the following two licensing options:

### 1. GNU General Public License v3.0 (GPLv3) 🆓

+ This option is governed by the Free Software Foundation's [GPL v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html).
+ Usage is subject to certain exceptions as outlined in this policy.

### 2. Commercial License 💼

+ Obtain a commercial license by contacting OpenIM.
+ For more details and licensing inquiries, please email 📧 [contact@openim.io](mailto:contact@openim.io).

## 🧩 Awesome features
1. This repository implement a business system, which consists of two parts: User related function and background management function
2. The business system depends on the api of the im system ([open-im-server repository](https://github.com/openimsdk/open-im-server)) and implement various functions by calling the api of the im system
3. User related part includes some regular functions like user login, user register, user info update, etc.
4. Background management provides api for admin to manage the im system containing functions like user management, message mangement,group management,etc.

## :busts_in_silhouette: Community

+ 💬 [Follow our Twitter account](https://twitter.com/founder_im63606)
+ 👫 [Join our Reddit](https://www.reddit.com/r/OpenIMessaging)
+ 🚀 [Join our Slack community](https://join.slack.com/t/openimsdk/shared_invite/zt-22720d66b-o_FvKxMTGXtcnnnHiMqe9Q)
+ :eyes: [Join our wechat (微信群)](https://openim-1253691595.cos.ap-nanjing.myqcloud.com/WechatIMG20.jpeg)
+ 📚 [OpenIM Community](https://github.com/openimsdk/community)
+ 💕 [OpenIM Interest Group](https://github.com/Openim-sigs)

## 🛫 Quick start 

> **Note**: You can get started quickly with OpenIM Chat.

### 📦 Installation

```bash
$ git clone https://github.com/openimsdk/chat openim-chat
$ cd openim-chat
```

### Chat Build

```bash
$ mage
```


### Chat Start

```bash
$ mage start
```

### Chat Check

```bash
$ mage check
```

### Chat Stop

```bash
$ mage stop
```

### 🚀 Boot Sequence
1. start [open-im-server](https://github.com/openimsdk/open-im-server) successfully.
2. modify related component configurations under the `config` folder.
3. build chat `mage`.
4. start chat `mage start`.

## Add REST RPC API

Please refer to "[How to add REST RPC API for OpenIM Chat](./HOW_TO_ADD_REST_RPC_API.md)".

## Setup LiveKit if you want to enable Audio and Video chat

Please refer to "[How to setup LiveKit server](./HOW_TO_SETUP_LIVEKIT_SERVER.md)".

## Contributing

Contributions to this project are welcome! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for details.

## Community Meetings
We want anyone to get involved in our community, we offer gifts and rewards, and we welcome you to join us every Thursday night.

We take notes of each [biweekly meeting](https://github.com/openimsdk/open-im-server/issues/381) in [GitHub discussions](https://github.com/openimsdk/open-im-server/discussions/categories/meeting), and our minutes are written in [Google Docs](https://docs.google.com/document/d/1nx8MDpuG74NASx081JcCpxPgDITNTpIIos0DS6Vr9GU/edit?usp=sharing).


## Who are using open-im-server
The [user case studies](https://github.com/openimsdk/community/blob/main/ADOPTERS.md) page includes the user list of the project. You can leave a [📝comment](https://github.com/openimsdk/open-im-server/issues/379) to let us know your use case.


## 🚨 License

chat is licensed under the  Apache 2.0 license. See [LICENSE](https://github.com/openimsdk/chat/tree/main/LICENSE) for the full license text.
