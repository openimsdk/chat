# OpenIM Chat 配置文件说明以及常用配置修改说明

## 配置文件说明

| Configuration File     | Description                                                                               |
| ---------------------- | ----------------------------------------------------------------------------------------- |
| **redis.yml**          | Redis 密码、地址等配置                                                                    |
| **mongodb.yml**        | MongoDB 用户名、密码、地址等配置                                                          |
| **log.yml**            | 日志级别及存储目录等配置                                                                  |
| **share.yml**          | OpenIM 各服务所需的公共配置，如 secret、Admin 后台密码等                                  |
| **discovery.yml**      | 服务发现对应的账号密码和服务名等配置                                                      |
| **chat-api-chat.yml**  | chat-api-chat 服务的监听 IP、端口等配置                                                   |
| **chat-api-admin.yml** | chat-api-admin 服务的监听 IP、端口等配置                                                  |
| **chat-rpc-chat.yml**  | chat-rpc-chat.yml 服务的监听 IP、端口以及登录注册验证码、是否允许注册和 livekit 等配置    |
| **chat-rpc-admin.yml** | chat-rpc-admin 服务的监听 IP、端口以及 chat 后台 token 过期策略和 chat 后台 Secret 等配置 |

## 常用配置修改

| 修改配置项                    | 配置文件             |
| ----------------------------- | -------------------- |
| 修改 OpenIM secret            | `share.yml`          |
| 生产环境日志调整              | `log.yml`            |
| 修改 chat Admin               | `share.yml`          |
| 修改验证码相关配置            | `chat-rpc-chat.yml`  |
| 允许用户注册                  | `chat-rpc-chat.yml`  |
| 修改 chat 后台 token 过期策略 | `chat-rpc-admin.yml` |
| 修改 chat 后台 Secret         | `chat-rpc-admin.yml` |

## 启动某个 OpenIM 服务的多个实例

若要启动某个 OpenIM 的多个实例，只需增加对应的端口数，并修改项目根目录下的`start-config.yml`文件，重启服务即可生效。例如，启动 2 个`chat-rpc`实例的配置如下：

```yaml
rpc:
  registerIP: ""
  listenIP: 0.0.0.0
  ports: [30300, 30301]
```

Modify start-config.yml:

```yaml
serviceBinaries:
  chat-rpc: 2
```
