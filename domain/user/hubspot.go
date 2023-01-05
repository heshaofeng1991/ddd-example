package user

type HubspotContactProperties struct {
	Company                   string `json:"company"`
	Email                     string `json:"email"`
	Firstname                 string `json:"firstname"`
	Lastname                  string `json:"lastname"`
	Phone                     string `json:"phone"`
	WhatKindOfServiceYouNeed  string `json:"what_kind_of_service_you_need,omitempty"`
	OtherServiceNeeded        string `json:"other_service_needed,omitempty"`
	ProductCategory           string `json:"product_category,omitempty"`
	OtherProductCategory      string `json:"other_product_category,omitempty"`
	TargetMarket              string `json:"target_market,omitempty"`
	Website                   string `json:"website,omitempty"`
	SellingPlatform           string `json:"selling_platform,omitempty"`
	OtherSellingPlatform      string `json:"other_selling_platform,omitempty"`
	WhatStageIsYourBusinessIn string `json:"what_stage_is_your_business_in_,omitempty"`
	TrafficSource             string `json:"traffic_source,omitempty"`
	StoreIntegration          bool   `json:"store_integration,omitempty"`
	SystemState               string `json:"system_state,omitempty"`
}

type HubspotContact struct {
	HsObjectID string                    `json:"hs_object_id"`
	Properties *HubspotContactProperties `json:"properties"`
}

func NewHubspotContact(user *User) (*HubspotContact, error) {
	contact := &HubspotContact{
		HsObjectID: user.guideInfo.HsObjectID,
		Properties: &HubspotContactProperties{
			Email:            user.Email(),
			Firstname:        user.Username(),
			Phone:            user.Phone(),
			Website:          user.Website(),
			TrafficSource:    user.SourceTag(),
			StoreIntegration: user.GuideInfo().status&4 > 0,
		},
	}

	for _, item := range user.GuideInfo().Questions() {
		switch item.Title {
		case "required_services":
			if values, ok := item.Answer.([]interface{}); ok {
				for _, value := range values {
					if _, ok := RequiredServices[value.(string)]; ok {
						contact.Properties.WhatKindOfServiceYouNeed = contact.Properties.WhatKindOfServiceYouNeed + ";" + value.(string)
					} else {
						contact.Properties.OtherServiceNeeded = value.(string)
					}
				}
			}
		case "product_categories":
			if values, ok := item.Answer.([]interface{}); ok {
				for _, value := range values {
					if _, ok := ProductCategories[value.(string)]; ok {
						contact.Properties.ProductCategory = contact.Properties.ProductCategory + ";" + value.(string)
					} else {
						contact.Properties.OtherProductCategory = value.(string)
					}
				}
			}
		case "target_markets":
			if values, ok := item.Answer.([]interface{}); ok {
				for _, value := range values {
					if _, ok := TargetMarkets[value.(string)]; ok {
						contact.Properties.TargetMarket = contact.Properties.TargetMarket + ";" + value.(string)
					}
				}
			}
		case "business_platforms":
			if values, ok := item.Answer.([]interface{}); ok {
				for _, service := range values {
					if _, ok := SellingPlatforms[service.(string)]; ok {
						contact.Properties.SellingPlatform = contact.Properties.SellingPlatform + ";" + service.(string)
					} else {
						contact.Properties.OtherSellingPlatform = service.(string)
					}
				}
			}
		case "business_stage":
			if value, ok := item.Answer.(string); ok {
				contact.Properties.WhatStageIsYourBusinessIn = value
			}
		case "website":
			if value, ok := item.Answer.(string); ok {
				contact.Properties.Website = value
			}
		}
	}

	if user.guideInfo.status&1 > 0 {
		contact.Properties.SystemState += "registered"
	}
	if user.guideInfo.status&2 > 0 {
		contact.Properties.SystemState += ";answered_questions"
	}
	if user.guideInfo.status&4 > 0 {
		contact.Properties.SystemState += ";store_integrated"
	}
	if user.guideInfo.Finished() && user.guideInfo.status&4 == 0 {
		contact.Properties.SystemState += ";skipped_store_integration"
	}

	return contact, nil
}

type HubspotPayload struct {
	Properties *HubspotContactProperties `json:"properties"`
}

type CRMUpdateProperties struct {
	HsObjectID                string
	WhatKindOfServiceYouNeed  []string
	ProductCategory           []string
	TargetMarket              []string
	Website                   string
	SellingPlatform           []string
	WhatStageIsYourBusinessIn string
	TrafficSource             string
	StoreIntegration          bool
	SystemState               []string
}

var ProductCategories = map[string]bool{
	"Bike & Vehicle":                 true,
	"Books":                          true,
	"Clothing, Shoes & Accessories":  true,
	"Computers/Tablets & Networking": true,
	"Consumer Electronics":           true,
	"Cosmetics & Skincare":           true,
	"Food":                           true,
	"Garden":                         true,
	"Health & Beauty":                true,
	"Home & Decoration":              true,
	"Jewelry & Watch":                true,
	"Kitchenware":                    true,
	"Pet Supplies":                   true,
	"Sports/Fitness & Outdoors":      true,
	"Toys & Games":                   true,
}

var SellingPlatforms = map[string]bool{
	"Shopify":              true,
	"Amazon":               true,
	"WooCommerce":          true,
	"BigCommerce":          true,
	"Wix":                  true,
	"Ebay":                 true,
	"Esty":                 true,
	"Lazada":               true,
	"Magento":              true,
	"Express/LCL only":     true,
	"Crowdfunding":         true,
	"Individual Order":     true,
	"Partnership Requests": true,
	"None/other":           true,
}

var RequiredServices = map[string]bool{
	"Sourcing":                     true,
	"Ecommerce Fulfillment":        true,
	"Dropshipping Fulfillment":     true,
	"Crowdfunding Fulfillment":     true,
	"Warehouse Fulfillment":        true,
	"Subscription Box Fulfillment": true,
	"FBA Fulfillment":              true,
	"Other":                        true,
}

var TargetMarkets = map[string]bool{
	"Andorra":                           true,
	"United Arab Emirates":              true,
	"Afghanistan":                       true,
	"Antigua and Barbuda":               true,
	"Anguilla":                          true,
	"Albania":                           true,
	"Armenia":                           true,
	"Angola":                            true,
	"Antarctica":                        true,
	"Argentina":                         true,
	"American Samoa":                    true,
	"Austria":                           true,
	"Australia":                         true,
	"Aruba":                             true,
	"Åland Islands":                     true,
	"Azerbaijan":                        true,
	"Bosnia and Herzegovina":            true,
	"Barbados":                          true,
	"Bangladesh":                        true,
	"Belgium":                           true,
	"Burkina Faso":                      true,
	"Bulgaria":                          true,
	"Bahrain":                           true,
	"Burundi":                           true,
	"Benin":                             true,
	"Saint Barthélemy":                  true,
	"Bermuda":                           true,
	"Brunei Darussalam":                 true,
	"Bolivia":                           true,
	"Bonaire, Sint Eustatius and Saba":  true,
	"Brazil":                            true,
	"Bahamas":                           true,
	"Bhutan":                            true,
	"Bouvet Island":                     true,
	"Botswana":                          true,
	"Belarus":                           true,
	"Belize":                            true,
	"Canada":                            true,
	"Cocos (Keeling) Islands":           true,
	"Congo, Democratic Republic of the": true,
	"Central African Republic":          true,
	"Congo":                             true,
	"Switzerland":                       true,
	"Côte d'Ivoire":                     true,
	"Cook Islands":                      true,
	"Chile":                             true,
	"Cameroon":                          true,
	"China":                             true,
	"Colombia":                          true,
	"Costa Rica":                        true,
	"Cuba":                              true,
	"Cabo Verde":                        true,
	"Curaçao":                           true,
	"Christmas Island":                  true,
	"Cyprus":                            true,
	"Czechia":                           true,
	"Germany":                           true,
	"Djibouti":                          true,
	"Denmark":                           true,
	"Dominica":                          true,
	"Dominican Republic":                true,
	"Algeria":                           true,
	"Ecuador":                           true,
	"Estonia":                           true,
	"Egypt":                             true,
	"Western Sahara":                    true,
	"Eritrea":                           true,
	"Spain":                             true,
	"Ethiopia":                          true,
	"Finland":                           true,
	"Fiji":                              true,
	"Falkland Islands":                  true,
	"Micronesia":                        true,
	"Faroe Islands":                     true,
	"France":                            true,
	"Gabon":                             true,
	"United Kingdom":                    true,
	"Grenada":                           true,
	"Georgia":                           true,
	"French Guiana":                     true,
	"Guernsey":                          true,
	"Ghana":                             true,
	"Gibraltar":                         true,
	"Greenland":                         true,
	"Gambia":                            true,
	"Guinea":                            true,
	"Guadeloupe":                        true,
	"Equatorial Guinea":                 true,
	"Greece":                            true,
	"South Georgia and the South Sandwich Islands": true,
	"Guatemala":                         true,
	"Guam":                              true,
	"Guinea-Bissau":                     true,
	"Guyana":                            true,
	"Hong Kong":                         true,
	"Heard Island and McDonald Islands": true,
	"Honduras":                          true,
	"Croatia":                           true,
	"Haiti":                             true,
	"Hungary":                           true,
	"Indonesia":                         true,
	"Ireland":                           true,
	"Israel":                            true,
	"Isle of Man":                       true,
	"India":                             true,
	"British Indian Ocean Territory":    true,
	"Iraq":                              true,
	"Iran (Islamic Republic of)":        true,
	"Iceland":                           true,
	"Italy":                             true,
	"Jersey":                            true,
	"Jamaica":                           true,
	"Jordan":                            true,
	"Japan":                             true,
	"Kenya":                             true,
	"Kyrgyzstan":                        true,
	"Cambodia":                          true,
	"Kiribati":                          true,
	"Comoros":                           true,
	"Saint Kitts and Nevis":             true,
	"North Korea":                       true,
	"South Korea":                       true,
	"Kuwait":                            true,
	"Cayman Islands":                    true,
	"Kazakhstan":                        true,
	"Lao People's Democratic Republic":  true,
	"Lebanon":                           true,
	"Saint Lucia":                       true,
	"Liechtenstein":                     true,
	"Sri Lanka":                         true,
	"Liberia":                           true,
	"Lesotho":                           true,
	"Lithuania":                         true,
	"Luxembourg":                        true,
	"Latvia":                            true,
	"Libya":                             true,
	"Morocco":                           true,
	"Monaco":                            true,
	"Moldova, Republic of":              true,
	"Montenegro":                        true,
	"Saint Martin (French part)":        true,
	"Madagascar":                        true,
	"Marshall Islands":                  true,
	"North Macedonia":                   true,
	"Mali":                              true,
	"Myanmar":                           true,
	"Macao":                             true,
	"Northern Mariana Islands":          true,
	"Martinique":                        true,
	"Mauritania":                        true,
	"Montserrat":                        true,
	"Malta":                             true,
	"Mauritius":                         true,
	"Maldives":                          true,
	"Malawi":                            true,
	"Mexico":                            true,
	"Malaysia":                          true,
	"Mozambique":                        true,
	"Namibia":                           true,
	"New Caledonia":                     true,
	"Niger":                             true,
	"Norfolk Island":                    true,
	"Nigeria":                           true,
	"Nicaragua":                         true,
	"Netherlands":                       true,
	"Norway":                            true,
	"Nepal":                             true,
	"Nauru":                             true,
	"Niue":                              true,
	"New Zealand":                       true,
	"Oman":                              true,
	"Panama":                            true,
	"Peru":                              true,
	"French Polynesia":                  true,
	"Papua New Guinea":                  true,
	"Philippines":                       true,
	"Pakistan":                          true,
	"Poland":                            true,
	"Saint Pierre and Miquelon":         true,
	"Pitcairn":                          true,
	"Puerto Rico":                       true,
	"Palestine, State of":               true,
	"Portugal":                          true,
	"Palau":                             true,
	"Paraguay":                          true,
	"Qatar":                             true,
	"Réunion":                           true,
	"Romania":                           true,
	"Serbia":                            true,
	"Russian Federation":                true,
	"Rwanda":                            true,
	"Saudi Arabia":                      true,
	"Solomon Islands":                   true,
	"Seychelles":                        true,
	"Sudan":                             true,
	"Sweden":                            true,
	"Singapore":                         true,
	"Saint Helena, Ascension and Tristan da Cunha": true,
	"Slovenia":                             true,
	"Svalbard and Jan Mayen":               true,
	"Slovakia":                             true,
	"Sierra Leone":                         true,
	"San Marino":                           true,
	"Senegal":                              true,
	"Somalia":                              true,
	"Suriname":                             true,
	"South Sudan":                          true,
	"Sao Tome and Principe":                true,
	"El Salvador":                          true,
	"Sint Maarten (Dutch part)":            true,
	"Syrian Arab Republic":                 true,
	"Eswatini":                             true,
	"Turks and Caicos Islands":             true,
	"Chad":                                 true,
	"French Southern Territories":          true,
	"Togo":                                 true,
	"Thailand":                             true,
	"Tajikistan":                           true,
	"Tokelau":                              true,
	"Timor-Leste":                          true,
	"Turkmenistan":                         true,
	"Tunisia":                              true,
	"Tonga":                                true,
	"The Republic of Turkey":               true,
	"Trinidad and Tobago":                  true,
	"Tuvalu":                               true,
	"Taiwan, Province of China":            true,
	"Tanzania, United Republic of":         true,
	"Ukraine":                              true,
	"Uganda":                               true,
	"United States Minor Outlying Islands": true,
	"United States of America":             true,
	"Uruguay":                              true,
	"Uzbekistan":                           true,
	"Vatican City":                         true,
	"Saint Vincent and the Grenadines":     true,
	"Venezuela (Bolivarian Republic of)":   true,
	"Virgin Islands (British)":             true,
	"Virgin Islands (U.S.)":                true,
	"Viet Nam":                             true,
	"Vanuatu":                              true,
	"Wallis and Futuna":                    true,
	"Samoa":                                true,
	"Yemen":                                true,
	"Mayotte":                              true,
	"South Africa":                         true,
	"Zambia":                               true,
	"Zimbabwe":                             true,
	"Somaliland, Republic of (N. Somalia)": true,
	"Bonaire":                              true,
	"St. Eustatius":                        true,
	"Canary Islands":                       true,
	"Netherlands Antilles":                 true,
	"Kosovo":                               true,
	"New Zealand Territory Islands":        true,
	"Campione/ Lake Lugano (Italy)":        true,
	"St. Lucia":                            true,
	"Northern Ireland (United Kingdom)":    true,
	"St. Barthelemy":                       true,
	"Mongolia":                             true,
	"Channel Islands":                      true,
}
