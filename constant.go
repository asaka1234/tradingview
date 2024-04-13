package tradingview

const (
	WidgetDataWssAddress = "wss://widgetdata.tradingview.com/socket.io/websocket"
	DataWssAddress       = "wss://data.tradingview.com/socket.io/websocket"
)

type AuthTokenType string

const (
	AuthTokenTypeWidgetUserToken       AuthTokenType = "widget_user_token"
	AuthTokenTypeUnauthorizedUserToken AuthTokenType = "unauthorized_user_token"
)

var FieldList = []string{
	"short_name",    //TRXUSDT
	"pro_name",      //BINANCE:TRXUSDT
	"original_name", //BINANCE:TRXUSDT
	"description",   //APPLE INC / US DOLLAR

	"volume",  //当前volume
	"lp",      //当前报价
	"lp_time", //当前报价时间, 1712979679， 十位时间戳
	"ch",      //价格变化值，比如跌了20，则是-20
	"chp",     //change percent, 价格变化百分比，比如-1.2 则指的跌了-1.2%
	"bid",     //bid报价
	"ask",     //ask报价
	"high_price",
	"low_price",
	"open_price",
	"prev_close_price",

	"type",                 //spot | forex  | index  |  stock  | bond | fund | futures
	"logoid",               //indices/nasdaq-100  是指数的logo
	"base_currency_id",     //base EUR  (EUR/USDT)
	"base-currency-logoid", //base 图片地址(country/EU)
	"currency_id",          //quote USD
	"currency-logoid",      //country/US 这个是图片的地址
	"currency_code",        //USD

	"current_session",
	"exchange", //报价来源交易所  BINANCE
	"format",
	"fractional",  //false
	"is_tradable", //true
	"language",
	"local_description",
	"listed_exchange",
	"minmov",
	"minmove2",
	"pricescale",  //100000
	"typespecs",   //crypto | cfd
	"update_mode", //streaming
	"variable_tick_size",
	"value_unit_id",
}
