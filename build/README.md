# Building OpenIM Chat Images

The chat service images share one Dockerfile. The compose file defines each service and passes a command path and binary name into that Dockerfile.

## Requirements

Docker with the compose and buildx plugins is required. `jq` is required when pushing images with `PUSH=true`.

## Local Build

Build all chat service images locally through compose:

```bash
build/build.sh
```

Build release-optimized local images:

```bash
RELEASE=true build/build.sh
```

Build one service locally through compose:

```bash
docker compose -f build/images/openim-chat/docker-compose.build.yml build openim-chat-api
```

## Push Images

CI uses the same build script with `PUSH=true`. The script reads services from `build/images/openim-chat/docker-compose.build.yml`, then builds and pushes each service defined there.

The required variables are:

```bash
PUSH=true \
RELEASE=true \
IMAGE_TAGS="v1.0.0,sha-abcdef" \
IMAGE_REGISTRIES="openim,ghcr.io/openimsdk,registry.cn-hangzhou.aliyuncs.com/openimsdk" \
build/build.sh
```

`IMAGE_TAGS` and `IMAGE_REGISTRIES` can be comma, space, or newline separated. For each compose service, the script reads `build.context`, `build.dockerfile`, `build.args.CMD_PATH`, and `build.args.BINARY_NAME`, then pushes every `registry/BINARY_NAME:tag` combination.
