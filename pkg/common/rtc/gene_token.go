package rtc

import (
	"time"

	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/livekit/protocol/auth"
)

func GetLiveKitToken(room string, identity string) (string, error) {
	apiKey := config.Config.LiveKit.Key
	apiSecret := config.Config.LiveKit.Secret

	//fmt.Printf("livekit apiKey=%s\n", apiKey)
	//fmt.Printf("livekit apiSecret=%s\n", apiSecret)
	//
	//fmt.Printf("room=%s\n", room)
	//fmt.Printf("identity=%s\n", identity)

	token, err := geneLiveKitToken(apiKey, apiSecret, room, identity)

	//fmt.Printf("livekit token=%s\n", token)

	return token, err
}

func geneLiveKitToken(key string, secret string, room string, identity string) (string, error) {
	at := auth.NewAccessToken(key, secret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     room,
	}
	at.AddGrant(grant).
		SetIdentity(identity).
		SetValidFor(time.Hour)

	return at.ToJWT()
}

func GetLiveKitServerUrl() string {
	return config.Config.LiveKit.LiveKitUrl
}
