# Setting Up LiveKit Server for OpenIM Chat

OpenIM Chat uses the LiveKit server as the media server to support video calls and video meeting services.

## About LiveKit

[LiveKit](https://github.com/livekit/livekit-server) is an open-source WebRTC SFU written in Go, built on top of the excellent [Pion](https://github.com/pion) project. For more information, visit the [LiveKit website](https://livekit.io/).

## Quick Start

To self-host LiveKit, start the server with the following Docker command:
> replace `your-server-ip` with your server IP address

```bash
docker run -d \
    -p 7880:7880 \
    -p 7881:7881 \
    -p 7882:7882/udp \
    -v $PWD/livekit/livekit.yaml:/livekit.yaml \
    livekit/livekit-server \
    --config /livekit.yaml \
    --bind 0.0.0.0 \
    --node-ip=your-server-ip
```

## Viewing Logs

To check the server logs and ensure everything is running correctly, use the following command:

```bash
docker logs livekit/livekit-server
```

## Configuring the LiveKit Address in OpenIM Chat

- **If Chat is deployed from source code**, update the `config/chat-rpc-chat.yml` file to configure the LiveKit server address:

```yaml
liveKit:
  url: "ws://127.0.0.1:7880"  # ws://your-server-ip:7880 or wss://your-domain/path
```

- **If Chat is deployed via Docker**, add the following environment variable to the `docker-compose.yaml` file:

```yaml
CHATENV_CHAT_RPC_CHAT_LIVEKIT_URL="ws://your-server-ip:7880"  # or wss://your-domain/path
```

Open the following ports: TCP: 7880-7881, UDP: 7882, and UDP: 50000-60000.

By following these steps, you can set up and configure the LiveKit server for use with OpenIM Chat.



## More about Deploying LiveKi

For detailed instructions on deploying LiveKit, refer to the self-hosting [deployment documentation](https://docs.livekit.io/realtime/self-hosting/deployment/).
