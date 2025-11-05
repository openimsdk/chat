package admin

import "github.com/openimsdk/tools/utils/datautil"

func (x *GetClientConfigResp) ApiFormat() {
	datautil.InitMap(&x.Config)
}
