package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jinzhu/gorm"
	"quote/utils"
)

type Currency struct {
	CommonModel
	Code        string `json:"code" sql:"-"`
	Key         string `json:"key"`                                                                 // 币的唯一标示
	Symbol      string `json:"symbol"`                                                              // 币的简称
	Logo        string `json:"logo"`                                                                // 币的图标
	Visible     bool   `json:"visible"`                                                             // 是否可用
	Erc20       bool   `json:"erc20"`                                                               // 是否erc20
	Erc23       bool   `json:"erc23"`                                                               // 是否而出3
	Fixed       int    `json:"fixed" gorm:"default:6"`                                              // 小数位精度
	OptionsJson string `json:"-" gorm:"column:options;type:varchar(64);default:\"[10,20,50,100]\""` // 转账快捷选项JSON
	Options     []int  `sql:"-" json:"options"`                                                     // 转账快捷选项
}

func InitAllCurrencies(db *utils.GormDB) {
	// InitCurrenciesFromRfinex()
	db.Where("visible = ?", true).Find(&AllCurrencies)
}
func (currency *Currency) IsEthereum() (result bool) {
	if currency.Symbol == "ethereum" || currency.Erc20 || currency.Erc23 {
		result = true
	}
	return
}

func InitCurrenciesFromRfinex() {
	resp, err := http.Get("https://rfinex.vip/api/v4/currencies")
	if err != nil {
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var respo struct {
		Head map[string]string
		Body []Currency
	}
	json.Unmarshal(body, &respo)
	db := utils.DbBegin()
	defer db.DbRollback()
	for _, cu := range respo.Body {
		cu.Symbol = cu.Code
		db.Save(&cu)
	}
	db.DbCommit()
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
	json.Unmarshal([]byte(c.OptionsJson), &c.Options)
}
