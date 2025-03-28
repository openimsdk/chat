package botstruct

import (
	"strings"

	"github.com/openimsdk/chat/pkg/common/constant"
)

func IsAgentUserID(userID string) bool {
	return strings.HasPrefix(userID, constant.AgentUserIDPrefix)
}
