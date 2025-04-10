package imapi

import "github.com/openimsdk/protocol/sdkws"

// SendSingleMsgReq defines the structure for sending a message to multiple recipients.
type SendSingleMsgReq struct {
	// groupMsg should appoint sendID
	SendID          string                 `json:"sendID"`
	Content         string                 `json:"content" binding:"required"`
	OfflinePushInfo *sdkws.OfflinePushInfo `json:"offlinePushInfo"`
	Ex              string                 `json:"ex"`
}
type SendSingleMsgResp struct{}
