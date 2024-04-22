// Copyright © 2023 OpenIM open source community. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apistruct

import "github.com/openimsdk/protocol/sdkws"

type UserRegisterResp struct {
	ImToken   string `json:"imToken"`
	ChatToken string `json:"chatToken"`
	UserID    string `json:"userID"`
}

type LoginResp struct {
	ImToken   string `json:"imToken"`
	ChatToken string `json:"chatToken"`
	UserID    string `json:"userID"`
}

type UpdateUserInfoResp struct{}

type CallbackAfterSendSingleMsgReq struct {
	CommonCallbackReq
	RecvID string `json:"recvID"`
}

type CommonCallbackReq struct {
	SendID           string   `json:"sendID"`
	CallbackCommand  string   `json:"callbackCommand"`
	ServerMsgID      string   `json:"serverMsgID"`
	ClientMsgID      string   `json:"clientMsgID"`
	OperationID      string   `json:"operationID"`
	SenderPlatformID int32    `json:"senderPlatformID"`
	SenderNickname   string   `json:"senderNickname"`
	SessionType      int32    `json:"sessionType"`
	MsgFrom          int32    `json:"msgFrom"`
	ContentType      int32    `json:"contentType"`
	Status           int32    `json:"status"`
	CreateTime       int64    `json:"createTime"`
	Content          string   `json:"content"`
	Seq              uint32   `json:"seq"`
	AtUserIDList     []string `json:"atUserList"`
	SenderFaceURL    string   `json:"faceURL"`
	Ex               string   `json:"ex"`
}

type CallbackAfterSendSingleMsgResp struct {
	CommonCallbackResp
}

type CommonCallbackResp struct {
	ActionCode int32  `json:"actionCode"`
	ErrCode    int32  `json:"errCode"`
	ErrMsg     string `json:"errMsg"`
	ErrDlt     string `json:"errDlt"`
	NextCode   int32  `json:"nextCode"`
}

type TextElem struct {
	Content string `json:"content" validate:"required"`
}

type PictureElem struct {
	SourcePath      string          `mapstructure:"sourcePath"`
	SourcePicture   PictureBaseInfo `mapstructure:"sourcePicture"   validate:"required"`
	BigPicture      PictureBaseInfo `mapstructure:"bigPicture"      validate:"required"`
	SnapshotPicture PictureBaseInfo `mapstructure:"snapshotPicture" validate:"required"`
}

type PictureBaseInfo struct {
	UUID   string `mapstructure:"uuid"`
	Type   string `mapstructure:"type"   validate:"required"`
	Size   int64  `mapstructure:"size"`
	Width  int32  `mapstructure:"width"  validate:"required"`
	Height int32  `mapstructure:"height" validate:"required"`
	Url    string `mapstructure:"url"    validate:"required"`
}

type SendMsgReq struct {
	// RecvID uniquely identifies the receiver and is required for one-on-one or notification chat types.
	RecvID string `json:"recvID" binding:"required_if" message:"recvID is required if sessionType is SingleChatType or NotificationChatType"`
	SendMsg
}

// SendMsg defines the structure for sending messages with various metadata.
type SendMsg struct {
	// SendID uniquely identifies the sender.
	SendID string `json:"sendID" binding:"required"`

	// GroupID is the identifier for the group, required if SessionType is 2 or 3.
	GroupID string `json:"groupID" binding:"required_if=SessionType 2|required_if=SessionType 3"`

	// SenderNickname is the nickname of the sender.
	SenderNickname string `json:"senderNickname"`

	// SenderFaceURL is the URL to the sender's avatar.
	SenderFaceURL string `json:"senderFaceURL"`

	// SenderPlatformID is an integer identifier for the sender's platform.
	SenderPlatformID int32 `json:"senderPlatformID"`

	// Content is the actual content of the message, required and excluded from Swagger documentation.
	Content map[string]any `json:"content" binding:"required" swaggerignore:"true"`

	// ContentType is an integer that represents the type of the content.
	ContentType int32 `json:"contentType" binding:"required"`

	// SessionType is an integer that represents the type of session for the message.
	SessionType int32 `json:"sessionType" binding:"required"`

	// IsOnlineOnly specifies if the message is only sent when the receiver is online.
	IsOnlineOnly bool `json:"isOnlineOnly"`

	// NotOfflinePush specifies if the message should not trigger offline push notifications.
	NotOfflinePush bool `json:"notOfflinePush"`

	// SendTime is a timestamp indicating when the message was sent.
	SendTime int64 `json:"sendTime"`

	// OfflinePushInfo contains information for offline push notifications.
	OfflinePushInfo *sdkws.OfflinePushInfo `json:"offlinePushInfo"`
}
