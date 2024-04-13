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
	"base-currency-logoid",
	"currency-logoid",
	"currency_code", //USDT
	"currency_id",   //XTVCUSDT
	"base_currency_id",
	"current_session",
	"description", //具体pair的描述  TRON / TetherUS
	"exchange",    //报价来源交易所  BINANCE
	"format",
	"fractional",  //false
	"is_tradable", //true
	"language",
	"local_description",
	"listed_exchange",
	"logoid",
	"minmov",
	"minmove2",
	"original_name", //BINANCE:TRXUSDT
	"pricescale",    //100000
	"pro_name",      //BINANCE:TRXUSDT
	"short_name",    //TRXUSDT
	"type",          //spot | forex
	"typespecs",     //crypto | cfd
	"update_mode",   //streaming
	"variable_tick_size",
	"value_unit_id",
}