package airwallex

import (
	"context"
)

func (t *Client) GetAccount(ctx context.Context) (*Account, error) {

	var r Account
	err := t.get(ctx, "/api/v1/account", nil, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

/*

{
  "account_details" : {
    "attachments" : {
      "additional_files" : [ ]
    },
    "business_details" : {
      "account_usage" : {
        "estimated_monthly_revenue" : {
          "amount" : "50000.0",
          "currency" : "USD"
        },
        "product_reference" : [ "RECEIVE_TRANSFERS", "CONVERT_FUNDS", "ACCEPT_PAYMENTS", "GET_PAID", "CREATE_CARDS", "USE_AWX_API" ]
      },
      "as_trustee" : null,
      "attachments" : {
        "business_documents" : [ {
          "description" : "1. Tutorduck .pdf",
          "file_id" : "10bdbc02-835f-4816-8443-2ac9e5bf51a3",
          "tag" : "CERTIFICATE_OF_INCORPORATION"
        }, {
          "description" : "TutorDuck Limited -2024BR.pdf",
          "file_id" : "82f7bf52-fb65-4099-9603-2d5cf8263033",
          "tag" : "BUSINESS_LICENSE"
        }, {
          "file_id" : "20487397-444b-4db4-a923-8ed44fa756c2",
          "tag" : "OTHER"
        } ]
      },
      "business_address" : null,
      "business_identifiers" : [ {
        "country_code" : "HK",
        "number" : "67535185",
        "type" : "BRN"
      } ],
      "business_name" : "挑得有限公司",
      "business_name_english" : "Tutorduck Limited",
      "business_name_trading" : null,
      "business_start_date" : "2017-03-29",
      "business_structure" : "COMPANY",
      "contact_number" : null,
      "description_of_goods_or_services" : "我們公司在online 經營補習教育業務，分別提供上門一對一補習以及每月訂閱AI答題補習",
      "explanation_for_high_risk_countries_exposure" : null,
      "has_nominee_shareholders" : null,
      "industry_category_code" : "ICCV3_0103XX",
      "no_shareholders_with_over_25percent" : null,
      "operating_country" : [ "HK" ],
      "registration_address" : {
        "address_line1" : "UNIT 1602 16/F LUCKY CENTRE NO.165-171 WAN CHAI ROAD WAN CHAI HK",
        "country_code" : "HK",
        "postcode" : "999077",
        "state" : "Hong Kong",
        "suburb" : "Hong Kong"
      },
      "registration_address_english" : {
        "address_line1" : "UNIT 1602 16/F LUCKY CENTRE NO.165-171 WAN CHAI ROAD WAN CHAI HK",
        "country_code" : "HK",
        "postcode" : "999077",
        "state" : "Hong Kong",
        "suburb" : "Hong Kong"
      },
      "state_of_incorporation" : null,
      "url" : "https://www.tutorduck.com/",
      "urls" : [ "https://www.tutorduck.com/", "http://m.tutorduck.com/api/home/redirect" ]
    },
    "business_person_details" : [ {
      "attachments" : {
        "business_person_documents" : [ ]
      },
      "date_of_birth" : "1995-06-11",
      "first_name" : "嘉寶",
      "first_name_english" : "Ka Po",
      "identifications" : {
        "primary" : {
          "identification_type" : "PASSPORT",
          "issuing_country_code" : "HK",
          "passport" : {
            "effective_at" : "2015-12-17",
            "expire_at" : "2025-12-17",
            "front_file_id" : "f64b905b-ff95-4b46-9834-3ca80a5e7cda",
            "mrz_line1" : "P<CHNCHONG<<KA<PO<<<<<<<<<<<<<<<<<<<<<<<<<<<",
            "mrz_line2" : "KJ04803147CHN9506114M2512174<R7418441<<<<<86",
            "number" : "KJ0480314"
          }
        }
      },
      "last_name" : "莊",
      "last_name_english" : "Chong",
      "nationality" : "HK",
      "person_id" : "09ace539-26d5-4f92-9c49-74512da378e0",
      "roles" : [ "BENEFICIAL_OWNER", "DIRECTOR", "AUTHORISED_PERSON" ]
    } ],
    "store_details" : {
      "cross_border_transaction_percent" : "100.0",
      "dispute_percent" : "0",
      "employee_size" : 0,
      "estimated_transaction_volume" : {
        "average_amount_per_transaction" : "60.0",
        "currency" : "USD",
        "max_amount_per_transaction" : "500.0",
        "monthly_transaction_amount" : "20000.0"
      },
      "financial_statements" : [ ],
      "fulfillment_days" : 16,
      "industry_code" : "ICCV3_0006XX",
      "mcc" : "7372",
      "operating_models" : [ "ONLINE_CHECKOUT" ],
      "payment_distribution" : [ ],
      "selling_to_country_codes" : [ ],
      "shipping_from_country_codes" : [ ],
      "store_description" : "The client main business is an AI short video generation engine, which is mainly suitable for short video creators at home and abroad.",
      "store_name" : "Tutorduck Limited KYC Status",
      "store_photos" : [ ],
      "store_websites" : [ {
        "url" : "https://www.veogo.ai"
      } ]
    },
    "trustee_details" : null,
    "individual_details" : null,
    "legal_entity_id" : "le_qv8IH-4cPsG8XCh095pAnQ",
    "legal_entity_identifier" : null,
    "legal_entity_type" : "BUSINESS"
  },
  "created_at" : "2025-04-11T14:04:18+0000",
  "customer_agreements" : {
    "agreed_to_data_usage" : false,
    "agreed_to_terms_and_conditions" : true,
    "opt_in_for_marketing" : false,
    "terms_and_conditions" : {
      "device_data" : { },
      "service_agreement_type" : "FULL"
    }
  },
  "id" : "acct_Q6NUqVssMGu--2O8trlsVw",
  "nickname" : "Veogo",
  "primary_contact" : {
    "attachments" : {
      "identity_files" : [ ]
    },
    "email" : "support@tutorduck.com",
    "mobile" : "852-55435024"
  },
  "status" : "ACTIVE",
  "view_type" : "COMPLETE"
}
*/

type Account struct {
	AccountDetails struct {
		Attachments struct {
			AdditionalFiles []interface{} `json:"additional_files"`
		} `json:"attachments"`
		BusinessDetails struct {
			AccountUsage struct {
				EstimatedMonthlyRevenue struct {
					Amount   string `json:"amount"`
					Currency string `json:"currency"`
				} `json:"estimated_monthly_revenue"`
				ProductReference []string `json:"product_reference"`
			} `json:"account_usage"`
			AsTrustee   interface{} `json:"as_trustee"`
			Attachments struct {
				BusinessDocuments []struct {
					Description string `json:"description,omitempty"`
					FileId      string `json:"file_id"`
					Tag         string `json:"tag"`
				} `json:"business_documents"`
			} `json:"attachments"`
			BusinessAddress     interface{} `json:"business_address"`
			BusinessIdentifiers []struct {
				CountryCode string `json:"country_code"`
				Number      string `json:"number"`
				Type        string `json:"type"`
			} `json:"business_identifiers"`
			BusinessName                            string      `json:"business_name"`
			BusinessNameEnglish                     string      `json:"business_name_english"`
			BusinessNameTrading                     interface{} `json:"business_name_trading"`
			BusinessStartDate                       string      `json:"business_start_date"`
			BusinessStructure                       string      `json:"business_structure"`
			ContactNumber                           interface{} `json:"contact_number"`
			DescriptionOfGoodsOrServices            string      `json:"description_of_goods_or_services"`
			ExplanationForHighRiskCountriesExposure interface{} `json:"explanation_for_high_risk_countries_exposure"`
			HasNomineeShareholders                  interface{} `json:"has_nominee_shareholders"`
			IndustryCategoryCode                    string      `json:"industry_category_code"`
			NoShareholdersWithOver25Percent         interface{} `json:"no_shareholders_with_over_25percent"`
			OperatingCountry                        []string    `json:"operating_country"`
			RegistrationAddress                     struct {
				AddressLine1 string `json:"address_line1"`
				CountryCode  string `json:"country_code"`
				Postcode     string `json:"postcode"`
				State        string `json:"state"`
				Suburb       string `json:"suburb"`
			} `json:"registration_address"`
			RegistrationAddressEnglish struct {
				AddressLine1 string `json:"address_line1"`
				CountryCode  string `json:"country_code"`
				Postcode     string `json:"postcode"`
				State        string `json:"state"`
				Suburb       string `json:"suburb"`
			} `json:"registration_address_english"`
			StateOfIncorporation interface{} `json:"state_of_incorporation"`
			Url                  string      `json:"url"`
			Urls                 []string    `json:"urls"`
		} `json:"business_details"`
		BusinessPersonDetails []struct {
			Attachments struct {
				BusinessPersonDocuments []interface{} `json:"business_person_documents"`
			} `json:"attachments"`
			DateOfBirth      string `json:"date_of_birth"`
			FirstName        string `json:"first_name"`
			FirstNameEnglish string `json:"first_name_english"`
			Identifications  struct {
				Primary struct {
					IdentificationType string `json:"identification_type"`
					IssuingCountryCode string `json:"issuing_country_code"`
					Passport           struct {
						EffectiveAt string `json:"effective_at"`
						ExpireAt    string `json:"expire_at"`
						FrontFileId string `json:"front_file_id"`
						MrzLine1    string `json:"mrz_line1"`
						MrzLine2    string `json:"mrz_line2"`
						Number      string `json:"number"`
					} `json:"passport"`
				} `json:"primary"`
			} `json:"identifications"`
			LastName        string   `json:"last_name"`
			LastNameEnglish string   `json:"last_name_english"`
			Nationality     string   `json:"nationality"`
			PersonId        string   `json:"person_id"`
			Roles           []string `json:"roles"`
		} `json:"business_person_details"`
		StoreDetails struct {
			CrossBorderTransactionPercent string `json:"cross_border_transaction_percent"`
			DisputePercent                string `json:"dispute_percent"`
			EmployeeSize                  int    `json:"employee_size"`
			EstimatedTransactionVolume    struct {
				AverageAmountPerTransaction string `json:"average_amount_per_transaction"`
				Currency                    string `json:"currency"`
				MaxAmountPerTransaction     string `json:"max_amount_per_transaction"`
				MonthlyTransactionAmount    string `json:"monthly_transaction_amount"`
			} `json:"estimated_transaction_volume"`
			FinancialStatements      []interface{} `json:"financial_statements"`
			FulfillmentDays          int           `json:"fulfillment_days"`
			IndustryCode             string        `json:"industry_code"`
			Mcc                      string        `json:"mcc"`
			OperatingModels          []string      `json:"operating_models"`
			PaymentDistribution      []interface{} `json:"payment_distribution"`
			SellingToCountryCodes    []interface{} `json:"selling_to_country_codes"`
			ShippingFromCountryCodes []interface{} `json:"shipping_from_country_codes"`
			StoreDescription         string        `json:"store_description"`
			StoreName                string        `json:"store_name"`
			StorePhotos              []interface{} `json:"store_photos"`
			StoreWebsites            []struct {
				Url string `json:"url"`
			} `json:"store_websites"`
		} `json:"store_details"`
		TrusteeDetails        interface{} `json:"trustee_details"`
		IndividualDetails     interface{} `json:"individual_details"`
		LegalEntityId         string      `json:"legal_entity_id"`
		LegalEntityIdentifier interface{} `json:"legal_entity_identifier"`
		LegalEntityType       string      `json:"legal_entity_type"`
	} `json:"account_details"`
	CreatedAt          string `json:"created_at"`
	CustomerAgreements struct {
		AgreedToDataUsage          bool `json:"agreed_to_data_usage"`
		AgreedToTermsAndConditions bool `json:"agreed_to_terms_and_conditions"`
		OptInForMarketing          bool `json:"opt_in_for_marketing"`
		TermsAndConditions         struct {
			DeviceData struct {
			} `json:"device_data"`
			ServiceAgreementType string `json:"service_agreement_type"`
		} `json:"terms_and_conditions"`
	} `json:"customer_agreements"`
	Id             string `json:"id"`
	Nickname       string `json:"nickname"`
	PrimaryContact struct {
		Attachments struct {
			IdentityFiles []interface{} `json:"identity_files"`
		} `json:"attachments"`
		Email  string `json:"email"`
		Mobile string `json:"mobile"`
	} `json:"primary_contact"`
	Status   string `json:"status"`
	ViewType string `json:"view_type"`
}
