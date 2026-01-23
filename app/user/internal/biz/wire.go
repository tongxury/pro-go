package biz

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewUserBiz,
	NewAuthBiz,
	NewAdminBiz,
	NewPromotionBiz,
	NewUserCreditBiz,
)
