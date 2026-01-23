package enums

const (
	OrderSide_Buy  string = "buy"
	OrderSide_Sell string = "sell"
)

const (
	OrderCategory_Market      string = "market"
	OrderCategory_MarketQuick string = "market_quick"
	OrderCategory_Limit       string = "limit"
)

const (
	OrderStatus_Pending     string = "pending" // 未上链
	OrderStatus_Uncompleted string = "uncompleted"
	OrderStatus_Executing   string = "executing"
	OrderStatus_Completed   string = "completed"
	//OrderStatus_FeeSettled  string = "feeSettled"

	OrderStatus_Cancelled   string = "cancelled"
	OrderStatus_ChainFailed string = "failed"
)

var (
	UncompletedOrderStatuses = []string{OrderStatus_Uncompleted, OrderStatus_Executing}
	SuccessOrderStatuses     = []string{OrderStatus_Completed}
	EndedOrderStatuses       = []string{OrderStatus_Completed, OrderStatus_Cancelled, OrderStatus_ChainFailed}
)

var (
	Trigger_AutoSell = "autoSell"
	Trigger_bot      = "bot"
)
