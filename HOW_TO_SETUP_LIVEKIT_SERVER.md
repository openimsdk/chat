# How to setup LiveKit server

OpenIM Chat uses LiveKit server as the media server to support video call and video meeting services.


## Something about LiveKit

[Livekit](https://github.com/livekit/livekit-server) is an open source WebRTC SFU written in go, built on top of the excellent [Pion](https://github.com/pion) project. You can get more information about it on its website [livekit.io](https://livekit.io/).


## Setup LiveKit server on Linux

Please follow the following instructions to setup a LiveKit server to work with OpenIM Chat and OpenIM server on Linux server.

### Docker installation

For self hosting user, we suggest you install LiveKit server by docker pull. You can get a server ready for use in a short time.

```bash
sudo docker pull livekit/livekit-server
```

For cloud deployment, you may follow the [Deploy to a VM](https://docs.livekit.io/realtime/self-hosting/vm/#Deploy-to-a-VM) on [livekit.io](https://docs.livekit.io/).

### Generate configuration

To generate configuration file for LiveKit server, please refer to [Generate configuration](https://docs.livekit.io/realtime/self-hosting/vm/#Generate-configuration).

The keys section of the generated .yaml file is the "apiKey: apiSecret" pair. This key pair should be set to the LiveKit section as the value of <font color="#FF8000">key</font> and <font color="#FF8000">secret</font> arguments.

### Generate access token for video call participant

Please refer to [Generating tokens](https://docs.livekit.io/realtime/server/generating-tokens/) on [livekit.io](https://docs.livekit.io/).

### Run it

For self hosting user, you may start LiveKit with:

```bash
docker run --rm \
    -p 7880:7880 \
    -p 7881:7881 \
    -p 7882:7882/udp \
    -v $PWD/livekit.yaml:/livekit.yaml \
    livekit/livekit-server \
    --config /livekit.yaml \
    --bind 0.0.0.0
```

For cloud deployment user, you may follow the [Deploy to a VM](https://docs.livekit.io/realtime/self-hosting/vm/#Deploy-to-a-VM) on [livekit.io](https://docs.livekit.io/).

## More about Deploying LiveKit

Please refer to the self hosting [Deploying LiveKit](https://docs.livekit.io/realtime/self-hosting/deployment/) documentation.
