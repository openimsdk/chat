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

## 🛫 Quick start 

> **Note**: You can get started quickly with OpenIM Chat.

### 📦 Installation

```bash
git clone https://github.com/openimsdk/chat openim-chat && export openim-chat=$(pwd)/openim-chat && cd $openim-chat && make
```

### Developing chat

You can deploy OpenIM Chat in two ways, either from source (which requires openIM-server to be installed) or with [docker compose](https://github.com/openimsdk/openim-docker)

**Here's how to deploy from source code:**

If you wish to deploy chat, then you should first install and deploy OpenIM, this [open-im-server repository](https://github.com/openimsdk/open-im-server)

First, install openim-server in a new directory or location repository

```bash
git clone -b release-v3.5 https://github.com/OpenIMSDK/Open-IM-Server.git openim && export openim=$(pwd)/openim && cd $openim
sudo docker compose up -d
```

**Setting configuration items:**

```bash
make init
```

> Then modify the configuration file `config/config.yaml` according to your needs
> Note: If you want to use the mysql database, you need to modify the mysql configuration item in the configuration file. If you want to use the mongo database, you need to modify the mongo configuration item in the configuration file


Then go back to the chat directory, Installing Chat

**Starting MySQL:**

> Recent OpenIM versions have discontinued the MySQL component, necessitating a separate MySQL installation for chat deployment.

```bash
docker run -d \
  --name mysql \
  -p 13306:3306 \
  -p 3306:33060 \
  -v "$(pwd)/components/mysql/data:/var/lib/mysql" \
  -v "/etc/localtime:/etc/localtime" \
  -e MYSQL_ROOT_PASSWORD="openIM123" \
  --restart always \
  mysql:5.7
```

**MySQL Compatibility Note:**

The MySQL version 5.7 primarily supports the linux/amd64 architecture. Users operating on Mac or other architectures, or those with specific version requirements for MySQL, can opt for the latest version of the open-source MariaDB. MariaDB is compatible with MySQL and can be a suitable alternative.

```bash
docker run -d \
  --name mysql \
  -p 13306:3306 \
  -p 3306:33060 \
  -v "$(pwd)/components/mysql/data:/var/lib/mysql" \
  -v "/etc/localtime:/etc/localtime" \
  -e MYSQL_ROOT_PASSWORD="openIM123" \
  --restart always \
  mariadb:10.6
```

If you want to use a local or remote database, edit the `config/config.yaml` section after the initial configuration


**Install Chat:**

```bash
$ make build
$ make start
$ make check
```

## 🛫 Quick start 

> **Note**: You can get started quickly with chat.

### 🚀 Run

> **Note**: 
> We need to run the backend server first

```bash
$ make build

# OR build Specifying binary
$ make build BINS=admin-api

# OR build multiarch
$ make build-multiarch
$ make build-multiarch BINS="admin-api"

# OR use scripts build source code
$ ./scripts/build_all.sh
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
$ make all
```

### Chat Start

```bash
$ make start_all
# OR use scripts start
$ ./scripts/start_all.sh
```

### Chat Detection

```bash
$ make check
# OR use scripts check
$ ./scripts/check_all.sh --print-screen
```

### Chat Stop

```bash
$ make stop
# OR use scripts stop
$ ./scripts/stop_all.sh
```

## Contributing

Contributions to this project are welcome! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for details.

## Community Meetings
We want anyone to get involved in our community, we offer gifts and rewards, and we welcome you to join us every Thursday night.

We take notes of each [biweekly meeting](https://github.com/openimsdk/open-im-server/issues/381) in [GitHub discussions](https://github.com/openimsdk/open-im-server/discussions/categories/meeting), and our minutes are written in [Google Docs](https://docs.google.com/document/d/1nx8MDpuG74NASx081JcCpxPgDITNTpIIos0DS6Vr9GU/edit?usp=sharing).


## Who are using open-im-server
The [user case studies](https://github.com/openimsdk/community/blob/main/ADOPTERS.md) page includes the user list of the project. You can leave a [📝comment](https://github.com/openimsdk/open-im-server/issues/379) to let us know your use case.

![avatar](https://github.com/openimsdk/OpenIM-Docs/blob/main/docs/images/WechatIMG20.jpeg)

## 🚨 License

chat is licensed under the  Apache 2.0 license. See [LICENSE](https://github.com/openimsdk/chat/tree/main/LICENSE) for the full license text.
