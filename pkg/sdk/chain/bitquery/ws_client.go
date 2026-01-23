package bitquery

import (
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
	"time"
)

type WSClient struct {
	endpoint string
}

func NewWSClient() *WSClient {
	return &WSClient{
		endpoint: "wss://streaming.bitquery.io/eap?token=ory_at_RlCavpqd7pLFaz6SXTa9ZpPVYCziBRwIw0TQOYaKZ7Q.A76LPN4ovOupXC2rSeki2Hiz7csyjRglrwcca5Ac0p0",
		//endpoint: "wss://streaming.bitquery.io/eap?token=ory_at_b0KqyZPNpjOVAhCFKr8UrZcqCVFY0Ef_U_s9uoTSEt8.iGHNVhwrgwXNgfTkrg1SRLDOPTzUzVC128jfoLt4nFk",
	}
}

type Connection struct {
	*websocket.Conn
}

func (t *WSClient) Connect() (*Connection, error) {

	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
	}

	conn, _, err := dialer.Dial(t.endpoint, http.Header{
		"Sec-WebSocket-Protocol": []string{"graphql-ws"},
	})
	if err != nil {
		return nil, err
	}

	c := &Connection{conn}

	err = c.init()
	if err != nil {
		return nil, err
	}

	return c, nil

}

func (t *Connection) init() error {
	payload := map[string]interface{}{
		"type": "connection_init",
	}
	return t.WriteJSON(payload)
}

func (t *Connection) ReadMessageLoop(handle func(messageType int, p []byte, err error)) {
	for {
		messageType, message, err := t.ReadMessage()

		if err != nil {
			continue
		}

		if len(message) == 26 && strings.Contains(string(message), "connection_ack") {
			continue
		}
		//{
		//	"type" : "ka"
		//}
		//if len(message) == 14 && strings.ContainsAny(string(message), "type") {
		//	continue
		//}

		handle(messageType, message, err)
	}
}

func (t *Connection) SubscribeSolanaDexTradesSimple() error {

	params := map[string]interface{}{
		"type": "subscribe",
		"payload": map[string]interface{}{
			"query": `
subscription {
  Solana {
    DEXTrades {
      Block {
        Slot
        Time
      }
      ChainId
      Trade {
        Buy {
          Amount
          AmountInUSD
          Currency {
            Decimals
            MintAddress
            Native
            Symbol
            Uri
          }
          Price
          PriceInUSD
        }
        Dex {
          ProgramAddress
          ProtocolFamily
          ProtocolName
        }
        Index
        Sell {
          Amount
          Account {
            Address
            Owner
            Token {
              Owner
            }
          }
          AmountInUSD
          Price
          PriceInUSD
          Currency {
            Decimals
            MintAddress
            Name
            Native
            Symbol
            Uri
          }
        }
      }
      Transaction {
        Index
        Result {
          ErrorMessage
          Success
        }
        Signature
        Signer
      }
    }
  }
}
`,

			//"variables": map[string]interface{}{},
		},
	}

	return t.WriteJSON(params)
}

func (t *Connection) SubscribeSolanaDexTrades() error {

	params := map[string]interface{}{
		"type": "subscribe",
		"payload": map[string]interface{}{
			"query": `
subscription {
  Solana {
    DEXTrades {
      Block {
        Date
        Hash
        Height
        ParentHash
        ParentSlot
        Slot
        Time
      }
      ChainId
      Trade {
        Buy {
          Amount
          AmountInUSD
          Currency {
            CollectionAddress
            Decimals
            EditionNonce
            Fungible
            IsMutable
            Key
            MetadataAddress
            MintAddress
            Name
            Native
            PrimarySaleHappened
            ProgramAddress
            SellerFeeBasisPoints
            Symbol
            TokenCreator {
              Address
              Share
              Verified
            }
            TokenStandard
            UpdateAuthority
            Uri
            VerifiedCollection
            Wrapped
          }
          Order {
            Account
            BuySide
            LimitAmount
            LimitPrice
            Mint
            OrderId
            Owner
            Payer
          }
          Price
          PriceInUSD
          Account {
            Address
            Owner
            Token {
              Owner
            }
          }
        }
        Dex {
          ProgramAddress
          ProtocolFamily
          ProtocolName
        }
        Index
        PriceAsymmetry
        Sell {
          Amount
          Account {
            Address
            Owner
            Token {
              Owner
            }
          }
          AmountInUSD
          Price
          PriceInUSD
          Currency {
            CollectionAddress
            Decimals
            EditionNonce
            Fungible
            IsMutable
            Key
            MetadataAddress
            MintAddress
            Name
            Native
            PrimarySaleHappened
            ProgramAddress
            SellerFeeBasisPoints
            Symbol
            TokenCreator {
              Address
              Share
              Verified
            }
            TokenStandard
            UpdateAuthority
            Uri
            VerifiedCollection
            Wrapped
          }
          Order {
            Account
            BuySide
            LimitAmount
            LimitPrice
            Mint
            OrderId
            Owner
            Payer
          }
        }
        Market {
          MarketAddress
        }
      }
      Transaction {
        Fee
        FeeInUSD
        FeePayer
        Index
        Result {
          ErrorMessage
          Success
        }
        Signature
        Signer
      }
    }
  }
}
`,

			//"variables": map[string]interface{}{},
		},
	}

	return t.WriteJSON(params)
}

type SubscribeSolanaDexPoolsParams struct {
}

func (t *Connection) SubscribeSolanaDexPools() error {

	payload := map[string]interface{}{
		"type": "subscribe",
		"payload": map[string]interface{}{
			"query": `
subscription {
  Solana {
    DEXPools(where: {Transaction: {Result: {Success: true}}, Pool: {Dex: {ProtocolName: {is: "pump"}}}}) {
      Block {
        Time
		Slot
      }
	  Transaction {
		Signer
		Signature
      }
      Pool {
        Base {
          ChangeAmount
          PostAmount
          Price
          PriceInUSD
        }
        Quote {
          ChangeAmount
          PostAmount
          Price
          PriceInUSD
        }
        Dex {
          ProgramAddress
          ProtocolFamily
        }
        Market {
          BaseCurrency {
            MintAddress
            Name
            Symbol
			Uri
          }
          QuoteCurrency {
            MintAddress
            Name
            Symbol
			Uri
          }
          MarketAddress
        }
      }
    }
  }
}
`,
			//"variables": map[string]interface{}{},
		},
	}

	return t.WriteJSON(payload)
}

func (t *Connection) SubscribeSolanaTokenSupplyUpdates() error {
	payload := map[string]interface{}{
		"type": "subscribe",
		"payload": map[string]interface{}{
			"query": `
subscription {
  Solana{
    TokenSupplyUpdates {
      TokenSupplyUpdate {

 		Currency {
			CollectionAddress
			Decimals
			EditionNonce
			Fungible
			IsMutable
			Key
			MetadataAddress
			MintAddress
			Name
			Native
			PrimarySaleHappened
			ProgramAddress
			SellerFeeBasisPoints
			Symbol
			TokenCreator {
			  Address
			  Share
			  Verified
			}
			TokenStandard
			UpdateAuthority
			Uri
			VerifiedCollection
			Wrapped
	  	}

		Amount
        AmountInUSD
        PostBalance
        PostBalanceInUSD
        PreBalance
        PreBalanceInUSD
      }
      Transaction {
        Fee
        FeeInUSD
        FeePayer
        Index
        Result {
          ErrorMessage
          Success
        }
        Signature
        Signer
      }
    }
  }
}
`,
			//"variables": map[string]interface{}{},
		},
	}
	return t.WriteJSON(payload)
}
