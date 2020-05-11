package models

import (
	// "encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"

	"quote/utils"
)

type Currency struct {
	CommonModel
	Key     string `json:"key"`     // 币的唯一标示
	Symbol  string `json:"symbol"`  // 币的简称
	Source  string `json:"source"`  // 币的来源，交易所
	Visible bool   `json:"visible"` // 是否可用
	// Erc20   bool   `json:"erc20"`   // 是否erc20
	// Erc23   bool   `json:"erc23"`   // 是否而出3
	// Logo    string `json:"logo"`    // 币的图标

	//   Type             string `json:"-"`
	//   WebsiteSlug      string `json:"website_slug"`
	//   Website          string `json:"website"`
	//   CoinmarketcapUrl string `json:"coinmarketcap_url"`
	//   CoinmarketcapId  string `json:"coinmarketcap_id"`
	//   ContactAddress   string `json:"contact_address"`
	//   Twitter          string `json:"-"`
	//   Facebook         string `json:"-"`
	//   Github           string `json:"-"`
	//   LogoUrl          string `json:"logo_url"`
	//   Platform         bool   `json:"platform"`
	//   RfinexKey        string `json:"rfinex_key"`
}

func InitAllCurrencies(db *utils.GormDB) {
	db.Where("visible = ?", true).Find(&AllCurrencies)
}

func FindCurrencyById(id int) (Currency, error) {
	for _, currency := range AllCurrencies {
		if currency.Id == id {
			return currency, nil
		}
	}
	var currency Currency
	return currency, fmt.Errorf("No currency can be found.")
}

func FindCurrencyBySymbol(symbol string) (Currency, error) {
	for _, currency := range AllCurrencies {
		if currency.Symbol == symbol {
			return currency, nil
		}
	}
	var currency Currency
	return currency, fmt.Errorf("No currency can be found.")
}

func (c *Currency) AfterFind(db *gorm.DB) {
	// json.Unmarshal([]byte(c.OptionsJson), &c.Options)
}
