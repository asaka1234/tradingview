package tradingview

// SocketInterface ...
type SocketInterface interface {
	AddSymbols(symbol []interface{}) error
	RemoveSymbols(symbol []interface{}) error
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
	/*名字*/
	ShortName    *string `mapstructure:"short_name" json:"short_name,omitempty"`       //TRXUSDT
	ProName      *string `mapstructure:"pro_name" json:"pro_name,omitempty"`           //BINANCE:TRXUSDT
	OriginalName *string `mapstructure:"original_name" json:"original_name,omitempty"` //BINANCE:TRXUSDT
	Description  *string `mapstructure:"description" json:"description,omitempty"`     //APPLE INC / US DOLLAR
	/*报价*/
	Volume         *float64 `mapstructure:"volume" json:"volume,omitempty"`
	Price          *float64 `mapstructure:"lp" json:"lp,omitempty"`           //价格
	PriceTime      *float64 `mapstructure:"lp_time" json:"lp_time,omitempty"` //报价时间, 1712979679， 十位时间戳
	Chp            *float64 `mapstructure:"chp" json:"chp,omitempty"`         //价格变化百分比，比如-1.2 则指的跌了-1.2%
	Ch             *float64 `mapstructure:"ch" json:"ch,omitempty"`           //价格变化值，比如跌了20，则是-20
	Bid            *float64 `mapstructure:"bid" json:"bid,omitempty"`         //TODO 这个要看一下是否还有
	Ask            *float64 `mapstructure:"ask" json:"ask,omitempty"`         //TODO 这个要看一下是否还有
	HighPrice      *float64 `mapstructure:"high_price" json:"high_price,omitempty"`
	LowPrice       *float64 `mapstructure:"low_price" json:"low_price,omitempty"`
	OpenPrice      *float64 `mapstructure:"open_price" json:"open_price,omitempty"`
	PrevClosePrice *float64 `mapstructure:"prev_close_price" json:"prev_close_price,omitempty"`
	/*logo说明*/
	Type               *string `mapstructure:"type" json:"type,omitempty"`                                 //spot | forex  | index  |  stock  | bond | fund | futures
	LogoId             *string `mapstructure:"logoid" json:"logoid,omitempty"`                             //indices/nasdaq-100  是指数的logo
	BaseCurrencyId     *string `mapstructure:"base_currency_id" json:"base_currency_id,omitempty"`         //base EUR  (EUR/USDT)
	BaseCurrencyLogoId *string `mapstructure:"base-currency-logoid" json:"base-currency-logoid,omitempty"` //base 图片地址(country/EU)
	CurrencyId         *string `mapstructure:"currency_id" json:"currency_id,omitempty"`                   //quote USD
	CurrencyLogoId     *string `mapstructure:"currency-logoid" json:"currency-logoid,omitempty"`           //country/US 这个是图片的地址
	CurrencyCode       *string `mapstructure:"currency_code" json:"currency_code,omitempty"`               //USD
}

// Flags ...
type Flags struct {
	Flags []string `json:"flags"`
}

// OnReceiveDataCallback ...
type OnReceiveDataCallback func(symbol string, data *QuoteData)

// OnErrorCallback ...
type OnErrorCallback func(err error, context string)
