package types

import (
	"fmt"
	"github.com/stripe/stripe-go/v82"
	"store/pkg/sdk/conv"
	"strings"
)

type ClientReference struct {
	UserID        string
	PromotionCode string
}

func (t *ClientReference) ToString() *string {
	return stripe.String(conv.S2J(t))
}

func NewClientReference(userID, promotionCode string) *ClientReference {
	return &ClientReference{
		UserID: userID, PromotionCode: promotionCode,
	}
}

func ParseClientReference(params string) ClientReference {

	userID := conv.Int64(params)
	if userID > 0 {
		return ClientReference{
			UserID: conv.String(userID),
		}
	}

	var p ClientReference
	_ = conv.J2S([]byte(params), &p)

	return p
}

func NewLookupKey(level, cycle string) string {
	return fmt.Sprintf("%s-%s", level, cycle)
}

func ParseLookupKey(lookupKey string) (string, string) {
	parts := strings.Split(lookupKey, "-")
	return parts[0], parts[1]
}
