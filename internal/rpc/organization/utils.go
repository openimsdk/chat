package organization

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"math/big"
	"math/rand"
	"strconv"
	"time"
)

func genDepartmentID() string {
	r := utils.Md5(strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.FormatUint(rand.Uint64(), 10))
	bi := big.NewInt(0)
	bi.SetString(r[0:8], 16)
	return bi.String()
}
