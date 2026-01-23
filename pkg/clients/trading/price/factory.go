package pricing

import "store/pkg/sdk/chain/sol/jupiter"

type PricingFactory struct {
	jupiter      IPricing
	jupiterQuote IPricing
	gmgnai       IPricing
}

func NewPricingFactory() *PricingFactory {

	jc := jupiter.NewJupiterClient()
	jw := jupiter.NewSwapClient()

	return &PricingFactory{
		jupiter:      NewJupiter(jc),
		jupiterQuote: NewJupiterQuote(jw),
		gmgnai:       NewGmgnai(),
	}
}

func (t *PricingFactory) Get() IPricing {
	return t.jupiterQuote
}

func (t *PricingFactory) GetGmgnai() IPricing {
	return t.gmgnai
}
