package tradingview

import (
	"fmt"
	"testing"
)

func Test_Search(t *testing.T) {
	//res := Search("ETHUSDT")
	//fmt.Printf("%+v\n", res.SourceId)

	symbolId := GetSymbolId("ethUSDT")
	fmt.Printf("%+v\n", symbolId)
}
