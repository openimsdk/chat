# OpenIM Chat Configuration Files and Common Configuration Item Modifications Guide

## Configuration Files Explanation

| Configuration File     | Description                                                                                                                                |
| ---------------------- | ------------------------------------------------------------------------------------------------------------------------------------------ |
| **redis.yml**          | Configurations for Redis password, address, etc.                                                                                           |
| **mongodb.yml**        | Configurations for MongoDB username, password, address, etc.                                                                               |
| **share.yml**          | Common configurations needed by various OpenIM services, such as secret.                                                                   |
| **discovery.yml**      | Service discovery account, password, service name and other configurations.                                                                |
| **chat-api-chat.yml**  | Configurations for listening IP, port, etc., in chat-api-chat service.                                                                     |
| **chat-api-admin.yml** | Configurations for listening IP, port, etc., in chat-api-admin service.                                                                    |
| **chat-rpc-chat.yml**  | Configurations for listening IP, port, login registration verification code, registration allowance, and livekit in chat-rpc-chat service. |
| **chat-rpc-admin.yml** | Configurations for listening IP, port, chat backend token expiration policy, and chat backend Secret in chat-rpc-admin service.            |

## Common Configuration Item Modifications

| Configuration Item to Modify              | Configuration File   |
| ----------------------------------------- | -------------------- |
| Modify OpenIM secret                      | `share.yml`          |
| Production environment logs               | `log.yml`            |
| Modify chat Admin username and password   | `share.yml`          |
| Modify verification code                  | `chat-rpc-chat.yml`  |
| Allow user registration                   | `chat-rpc-chat.yml`  |
| Modify chat Admin token expiration policy | `chat-rpc-admin.yml` |
| Modify chat Admin Secret                  | `chat-rpc-admin.yml` |

## Starting Multiple Instances of an OpenIM Service

To launch multiple instances of an OpenIM service, you just need to increase the corresponding number of ports and modify the `start-config.yml` file in the project's root directory, then restart the service for the changes to take effect. For example, the configuration for launching 2 instances of `chat-rpc` is as follows:

```yaml
rpc:
  registerIP: ""
  listenIP: 0.0.0.0
  ports: [30300, 30301]
```

Modify `start-config.yml`:

```yaml
serviceBinaries:
  chat-rpc: 2
```
