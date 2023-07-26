package admin

import "github.com/OpenIMSDK/tools/utils"

func (x *GetClientConfigResp) ApiFormat() {
	utils.InitMap(&x.Config)
}
