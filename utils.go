package tradingview

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"math/rand"
	"strings"
	"time"
)

// GetRandomString ...
func GetRandomString(length int) string {
	var src = rand.NewSource(time.Now().UnixNano())
	var characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var letterIdxBits int = 6
	var letterIdxMask int64 = 1<<letterIdxBits - 1
	var letterIdxMax = 63 / letterIdxBits

	requestID := make([]byte, length)
	for i, cache, remain := length-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(characters) {
			requestID[i] = characters[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(requestID)
}

// GetStringRepresentation ...
func GetStringRepresentation(data interface{}) string {
	str, _ := json.Marshal(data)
	return string(str)
}

//------------------------------------------------------

type Instrument struct {
	Symbol             string `json:"symbol,omitempty"`
	Description        string `json:"description,omitempty"`
	Type               string `json:"type,omitempty"`
	CurrencyCode       string `json:"currency_code,omitempty"`
	CurrencyLogoId     string `json:"currency-logoid,omitempty"`
	BaseCurrencyLogoId string `json:"base-currency-logoid,omitempty"`
	/*
		ProviderId         string `json:"provider_id,omitempty"`
		Source2            struct {
			Id          string `json:"id,omitempty"`
			Name        string `json:"name,omitempty"`
			Description string `json:"description,omitempty"`
		} `json:"source2"`

	*/
	Exchange string `json:"exchange,omitempty"`
	SourceId string `json:"source_id,omitempty"`
}

// Search for a symbol based on query and category
func Search(symbol string) *Instrument {
	url := "https://symbol-search.tradingview.com/symbol_search/"

	var instrumentList []Instrument
	//查一下是否存在
	request := gorequest.New()
	_, _, errs := request.Get(url).Param("text", symbol).EndStruct(&instrumentList) //搜索下这个instrument
	if len(errs) > 0 {
		//找不到
		return nil
	}

	for _, instrument := range instrumentList {
		if instrument.Symbol == strings.ToUpper(symbol) {
			//找到了
			return &instrument
		}
	}
	return nil
}

// 从 btcusdt -> binace:btcusdt
func GetSymbolId(symbol string) string {
	ins := Search(symbol)
	if ins == nil {
		return ""
	}
	symbolId := strings.ToUpper(fmt.Sprintf("%s:%s", ins.SourceId, symbol))
	return symbolId
}
