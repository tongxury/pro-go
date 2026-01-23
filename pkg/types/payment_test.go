package types

import (
	"fmt"
	"strings"
	"testing"
)

func TestClientReference(t *testing.T) {

	//assert.Equal(t, "1751861118692156303", ParseClientReference("1751861118692156303").UserID)
	//assert.Equal(t, "1751861118692156303", ParseClientReference(conv.S2J(ClientReference{
	//	UserID:        "1751861118692156303",
	//	PromotionCode: "PromotionCode",
	//})).UserID)
	//assert.Equal(t, "PromotionCode", ParseClientReference(conv.S2J(ClientReference{
	//	UserID:        "1751861118692156303",
	//	PromotionCode: "PromotionCode",
	//})).PromotionCode)

	fmt.Println(strings.TrimSpace("                                  \n                       \n") == "")

}
