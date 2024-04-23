# How to add REST RPC API for OpenIM Chat
</br>

OpenIM Chat server works with the OpenIM Server to offer the powerfull IM service. It has many REST APIs which are called by the OpenIM Server, works as an Application Server to fullfill your system demands.

You can add your business extensions to the OpenIM Chat server by adding REST APIs and let the OpenIM Server call these extend APIs to fulfill your system demands. The following will show you how to add a REST API for it.

## Protobuf source file modifications

First, your should add the REST API request and response message declarations, and RPC method declarations to the '.proto' protobuf source code, then call 'protoc' command to generate '.pb.go' go protobuf implementation codes.

For example, to add a REST API to generate login token for video meeting, you may want to add a request message named <font color="#FF8000">GetTokenForVideoMeetingReq</font> and a response message named <font color="#FF8000">GetTokenForVideoMeetingResp</font>. This can be done by adding this two messages to '/pkg/proto/chat/chat.proto'. You also should add a RPC method named <font color="#FF8000">GetTokenForVideoMeeting()</font> fors the <font color="#FF8000">chat</font> service.

### chat.proto:

```proto
syntax = "proto3";
package OpenIMChat.chat;
import "pub/wrapperspb.proto";
import "pub/sdkws.proto";
import "common/common.proto";
option go_package = "github.com/openimsdk/chat/pkg/protocol/chat";

...

message GetTokenForVideoMeetingReq {
  string room = 1;
  string identity = 2;
}

message GetTokenForVideoMeetingResp {
  string serverUrl = 1;
  string token = 2;
}

service chat {
   ...
    // Audio/video call and video meeting
    rpc GetTokenForVideoMeeting(GetTokenForVideoMeetingReq) returns (GetTokenForVideoMeetingResp);
}
```

## Generate protobuf implementations
</br>

Then, we start a terminal to run 'pkg/proto/gen.sh' shell script, which will call protoc command to regenerate chat.pb.go with protobuf implementation codes for your add messages.


## Add Check() method for the request message
</br>

To check the parameters in the request message, we should add a Check() method for newly added Request message. Take chat.proto as an example, we add the a Check() member method for GetTokenForVideoMeetingReq class in 'pkg/proto/chat/chat.go'.

```go
func (x *GetTokenForVideoMeetingReq) Check() error {
    if x.Room == "" {
        errs.ErrArgs.WrapMsg("Room is empty")
    }
    if x.Identity == "" {
        errs.ErrArgs.WrapMsg("User Identity is empty")
    }
    return nil
}
```


## Add the URL for REST RPC API
</br>

Now, we should add a URL for REST RPC API.

For example, to add a REST API to generate login token for video meeting, you may want to add a RPC API named 'get_token', which is a member of the 'user' group, the REST URL can be '/user/rtc/get_token'.

This URL should be add to NewChatRoute() method in 'internal/api/router.go'.

```go
func NewChatRoute(router gin.IRouter, discov discoveryregistry.SvcDiscoveryRegistry) {
    ...

	user := router.Group("/user", mw.CheckToken)
	user.POST("/update", chat.UpdateUserInfo)                 // Edit personal information
	user.POST("/find/public", chat.FindUserPublicInfo)        // Get user's public information
	user.POST("/find/full", chat.FindUserFullInfo)            // Get all information of the user
	user.POST("/search/full", chat.SearchUserFullInfo)        // Search user's public information
	user.POST("/search/public", chat.SearchUserPublicInfo)    // Search all information of the user

    // your added code begins here
	user.POST("/rtc/get_token", chat.GetTokenForVideoMeeting) // Get token for video meeting

    ...
}
```

## Implement the REST service

### Implement the REST Service Api

We add REST RPC API implementation to the corresponding service implementation go file located in '/internal/api/'. For the chat service, we add codes to '/internal/api/chat.go'.

```go
...
func (o *ChatApi) GetTokenForVideoMeeting(c *gin.Context)
{
	a2r.Call(chat.ChatClient.GetTokenForVideoMeeting, o.chatClient, c)
}
...
```

### Implement the REST Service logic

We implement the REST Service logic in go files of the path '/internal/rpc/service_name/group_name.go'. For the user group of <font color="#FF8000">chat</font> service, the implementation logic should be added to '/internal/rpc/chat/user.go'.

Here we call the GetLiveKitServerUrl() and GetLiveKitToken() func in the rtc package to allocate a new GetTokenForVideoMeetingResp message and return it to the RPC caller.

```go
package chat

import (
    ...
    "github.com/openimsdk/chat/pkg/common/rtc"
    ...
)

...

func (o *chatSvr) GetTokenForVideoMeeting(ctx context.Context, req *chat.GetTokenForVideoMeetingReq) (*chat.GetTokenForVideoMeetingResp, error)
{
    if _, _, err := mctx.Check(ctx); err != nil {
        return nil, err
    }
    serverUrl := rtc.GetLiveKitServerUrl()
    token, err := rtc.GetLiveKitToken(req.Room, req.Identity)
    if err != nil {
        return nil, err
    }
    return &chat.GetTokenForVideoMeetingResp{
        ServerUrl: serverUrl,
        Token:     token,
    }, err
}
```

## Done

Now we have finished all the works to add an REST RPC API for OpenIM Chat. We may call

```shell
make build
```
to re-compile the chat project and call

```shell
make start
```
to start the OpenIM Chat server and test the <font color="#FF8000">chat</font> service's new REST RPC API.

### 