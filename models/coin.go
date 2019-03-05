package models

type Coin struct {
	CommonModel
	Type             string `json:"-"`
	Name             string `json:"name"`
	Symbol           string `json:"symbol"`
	WebsiteSlug      string `json:"website_slug"`
	Website          string `json:"website"`
	CoinmarketcapUrl string `json:"coinmarketcap_url"`
	CoinmarketcapId  string `json:"coinmarketcap_id"`
	ContactAddress   string `json:"contact_address"`
	Twitter          string `json:"-"`
	Facebook         string `json:"-"`
	Github           string `json:"-"`
	LogoUrl          string `json:"logo_url"`
	Platform         bool   `json:"platform"`
	RfinexKey        string `json:"rfinex_key"`

	Quotes []Quote `json:"quotes"`
}
