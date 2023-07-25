package admin

import "github.com/OpenIMSDK/Open-IM-Server/pkg/utils"

func (x *GetClientConfigResp) ApiFormat() {
	utils.InitMap(&x.Config)
}
