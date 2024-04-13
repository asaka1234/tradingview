package tradingview

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestExx_Signed(t *testing.T) {
	tradingviewsocket, err := Connect(
		DataWssAddress,                     //WidgetDataWssAddress,         //wss地址
		AuthTokenTypeUnauthorizedUserToken, //AuthTokenTypeWidgetUserToken, //鉴权方式
		func(symbol string, data *QuoteData) {
			fmt.Printf("%s\n", symbol)
			fmt.Printf("%+v\n", data)

			if data.Price != nil {
				fmt.Printf("price=%f\n", *data.Price)
			}
			if data.Volume != nil {
				fmt.Printf("volume=%f\n", *data.Volume)
			}
			//如果没有数据,证明没有任何change
			if data.Bid != nil {
				fmt.Printf("bid=%f\n", *data.Bid)
			}
			if data.Ask != nil {
				fmt.Printf("ask=%f\n", *data.Ask)
			}
			//fmt.Printf("%#v\n", *data.Price)
			if data.PriceTime != nil {
				fmt.Printf("PriceTime=%f\n", *data.PriceTime)
			}
			if data.Ch != nil {
				fmt.Printf("ch=%f\n", *data.Ch)
			}
			if data.Chp != nil {
				fmt.Printf("chp=%f\n", *data.Chp)
			}
		},
		func(err error, context string) {
			fmt.Printf("%#v", "error -> "+err.Error())
			fmt.Printf("%#v", "context -> "+context)
		},
	)
	if err != nil {
		panic("Error while initializing the trading view socket -> " + err.Error())
	}
	//tradingviewsocket.AddSymbol("FPMARKETS:USDJPY")
	//tradingviewsocket.AddSymbol("ICMARKETS:EURUSD")
	//tradingviewsocket.AddSymbol("CAPITALCOM:EURUSD")
	//tradingviewsocket.AddSymbol("OANDA:GBPUSD")
	//STOCK
	//tradingviewsocket.AddSymbol("NASDAQ:AAPL")
	//CRYPTO
	tradingviewsocket.AddSymbol("BINANCE:DOGEUSDT")
	//tradingviewsocket.AddSymbol("COINBASE:BTCUSDT")
	//FOREX
	//tradingviewsocket.AddSymbol("VANTAGE:USDJPY")
	//tradingviewsocket.AddSymbol("ICE:USDJPY")

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	<-quit
}
