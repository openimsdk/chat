package rtc

import (
	"fmt"
	"testing"
)

func TestGetLiveKitToken(t *testing.T) {
	apiKey := "APIDXJxJeCL8haY"
	apiSecret := "ak1qulJ3nfXeflQHWBdmQDc4js4ueMc5OnxoORVJC2xA"

	token, err := geneLiveKitToken(apiKey, apiSecret, "testRoom", "testUser001")
	if err == nil {
		fmt.Printf("Livekit token: %s\n", token)
	} else {
		t.Errorf("geneLiveKitToken failed! err=%s\n", err.Error())
	}
}
