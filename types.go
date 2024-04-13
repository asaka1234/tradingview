package tradingview

// SocketInterface ...
type SocketInterface interface {
	AddSymbol(symbol string) error
	RemoveSymbol(symbol string) error
	Init() error
	Close() error
}

// SocketMessage ...
type SocketMessage struct {
	Message string      `json:"m"`
	Payload interface{} `json:"p"`
}

// QuoteMessage ...
type QuoteMessage struct {
	Symbol string     `mapstructure:"n"`
	Status string     `mapstructure:"s"`
	Data   *QuoteData `mapstructure:"v"`
}

// QuoteData ...
type QuoteData struct {
	Price  *float64 `mapstructure:"lp" json:"lp"` //价格
	Volume *float64 `mapstructure:"volume" json:"volume"`
	Bid    *float64 `mapstructure:"bid" json:"bid"` //TODO 这个要看一下是否还有
	Ask    *float64 `mapstructure:"ask" json:"ask"` //TODO 这个要看一下是否还有

	PriceTime *float64 `mapstructure:"lp_time" json:"lp_time"` //报价时间, 1712979679， 十位时间戳
	Chp       *float64 `mapstructure:"chp" json:"chp"`         //价格变化百分比，比如-1.2 则指的跌了-1.2%
	Ch        *float64 `mapstructure:"ch" json:"ch"`           //价格变化值，比如跌了20，则是-20

	HighPrice      *float64 `mapstructure:"high_price" json:"high_price"`
	LowPrice       *float64 `mapstructure:"low_price" json:"low_price"`
	OpenPrice      *float64 `mapstructure:"open_price" json:"open_price"`
	PrevClosePrice *float64 `mapstructure:"prev_close_price" json:"prev_close_price"`
}

// Flags ...
type Flags struct {
	Flags []string `json:"flags"`
}

// OnReceiveDataCallback ...
type OnReceiveDataCallback func(symbol string, data *QuoteData)

// OnErrorCallback ...
type OnErrorCallback func(err error, context string)
