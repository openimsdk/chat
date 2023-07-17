# chat

### Modifying configuration items

Refer to `config/config.yaml` for configuration instructions

## 🧩 Awesome features


## 🛫 Quick start 

> **Note**: You can get started quickly with OpenIM Chat.

### 📦 Installation

```bash
git clone https://github.com/OpenIMSDK/chat openim-chat && export openim-chat=$(pwd)/openim-chat && cd $openim-chat && make
```

### Developing chat
```

```

## 🛫 Quick start 

> **Note**: You can get started quickly with chat.

### 🚀 Run

> **Note**: 
> We need to run the backend server first

```bash
make build
```
```

### 📖 Contributors get up to speed

Be good at using Makefile, it can ensure the quality of your project.

```bash
Usage: make <TARGETS> ...

Targets:
  all                          Build all the necessary targets. 🏗️
  build                        Build binaries by default. 🛠️
  go.build                     Build the binary file of the specified platform. 👨‍💻
  build-multiarch              Build binaries for multiple platforms. 🌍
  tidy                         tidy go.mod 📦
  style                        Code style -> fmt,vet,lint 🎨
  fmt                          Run go fmt against code. ✨
  vet                          Run go vet against code. 🔍
  generate                     Run go generate against code and docs. ✅
  lint                         Run go lint against code. 🔎
  test                         Run unit test ✔️
  cover                        Run unit test with coverage. 🧪
  docker-build                 Build docker image with the manager. 🐳
  docker-push                  Push docker image with the manager. 🔝
  docker-buildx-push           Push docker image with the manager using buildx. 🚢
  copyright-verify             Validate boilerplate headers for assign files. 📄
  copyright-add                Add the boilerplate headers for all files. 📝
  swagger                      Generate swagger document. 📚
  serve-swagger                Serve swagger spec and docs. 🌐
  clean                        Clean all builds. 🧹
  help                         Show this help info. ℹ️
```

> **Note**: 
> It's highly recommended that you run `make all` before committing your code. 🚀

```bash
make all
```

### Chat Start

 ```
 ./scripts/start_all.sh
 ```

### Chat Detection

 ```
cd scripts
./check_all.sh
 ```

### Chat Stop

 ```
cd scripts
./stop_all.sh
 ```

## Contributing

Contributions to this project are welcome! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for details.

## Community Meetings
We want anyone to get involved in our community, we offer gifts and rewards, and we welcome you to join us every Thursday night.

We take notes of each [biweekly meeting](https://github.com/OpenIMSDK/Open-IM-Server/issues/381) in [GitHub discussions](https://github.com/OpenIMSDK/Open-IM-Server/discussions/categories/meeting), and our minutes are written in [Google Docs](https://docs.google.com/document/d/1nx8MDpuG74NASx081JcCpxPgDITNTpIIos0DS6Vr9GU/edit?usp=sharing).


## Who are using Open-IM-Server
The [user case studies](https://github.com/OpenIMSDK/community/blob/main/ADOPTERS.md) page includes the user list of the project. You can leave a [📝comment](https://github.com/OpenIMSDK/Open-IM-Server/issues/379) to let us know your use case.

![avatar](https://github.com/OpenIMSDK/OpenIM-Docs/blob/main/docs/images/WechatIMG20.jpeg)

## 🚨 License

chat is licensed under the  Apache 2.0 license. See [LICENSE](https://github.com/OpenIMSDK/chat/tree/main/LICENSE) for the full license text.
