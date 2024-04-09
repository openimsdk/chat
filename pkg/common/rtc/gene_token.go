package rtc

import (
	"time"

	"github.com/livekit/protocol/auth"
	"github.com/openimsdk/chat/pkg/common/config"
)

func GetLiveKitToken(room string, identity string) (string, error) {
	apiKey := config.Config.LiveKit.Key
	apiSecret := config.Config.LiveKit.Secret
	return geneLiveKitToken(apiKey, apiSecret, room, identity)
}

func geneLiveKitToken(key string, secret string, room string, identity string) (string, error) {
	at := auth.NewAccessToken(key, secret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     room,
	}
	at.AddGrant(grant).SetIdentity(identity).SetValidFor(time.Hour)
	return at.ToJWT()
}

func GetLiveKitServerUrl() string {
	return config.Config.LiveKit.LiveKitUrl
}
