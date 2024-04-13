package tradingview

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"strconv"
	"sync"
)

// Socket ...
type TradingViewWebSocket struct {
	address       string        //websocket地址
	authTokenType AuthTokenType //初始化发送msg时需要指定token鉴权方式

	OnReceiveMarketDataCallback OnReceiveDataCallback
	OnErrorCallback             OnErrorCallback

	mx        sync.Mutex
	conn      *websocket.Conn
	isClosed  bool
	sessionID string
}

// Connect - Connects and returns the trading view socket object
func Connect(
	address string,
	authTokenType AuthTokenType,
	onReceiveMarketDataCallback OnReceiveDataCallback,
	onErrorCallback OnErrorCallback,
) (socket SocketInterface, err error) {
	socket = &TradingViewWebSocket{
		address:                     address,
		authTokenType:               authTokenType,
		OnReceiveMarketDataCallback: onReceiveMarketDataCallback,
		OnErrorCallback:             onErrorCallback,
	}

	err = socket.Init()

	return
}

// Init connects to the tradingview web socket
func (s *TradingViewWebSocket) Init() (err error) {
	s.mx = sync.Mutex{}
	s.isClosed = true
	s.conn, _, err = (&websocket.Dialer{}).Dial(s.address, getHeaders())
	if err != nil {
		s.onError(err, InitErrorContext)
		return
	}

	//链接上服务器后, server会推过来一个初始化确认信息
	err = s.checkFirstReceivedMessage()
	if err != nil {
		return
	}

	//创建一个session_id,这个是这次wss交互的唯一标记
	s.generateSessionID()

	err = s.sendConnectionSetupMessages()
	if err != nil {
		s.onError(err, ConnectionSetupMessagesErrorContext)
		return
	}

	s.isClosed = false
	go s.connectionLoop()

	return
}

// Close ...
func (s *TradingViewWebSocket) Close() (err error) {
	s.isClosed = true
	return s.conn.Close()
}

// AddSymbol ...
func (s *TradingViewWebSocket) AddSymbols(symbolList []interface{}) (err error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	//批量订阅
	symbols := append([]interface{}{s.sessionID}, symbolList...)

	err = s.sendSocketMessage(
		getSocketMessage("quote_add_symbols", symbols),
	)
	//单独加一下这个
	err = s.sendSocketMessage(
		getSocketMessage("quote_fast_symbols", symbols),
	)
	return
}

// RemoveSymbol ...
func (s *TradingViewWebSocket) RemoveSymbols(symbolList []interface{}) (err error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	//批量解除订阅
	symbols := append([]interface{}{s.sessionID}, symbolList...)

	err = s.sendSocketMessage(
		getSocketMessage("quote_remove_symbols", symbols),
	)
	return
}

func (s *TradingViewWebSocket) checkFirstReceivedMessage() (err error) {
	var msg []byte

	_, msg, err = s.conn.ReadMessage()
	if err != nil {
		s.onError(err, ReadFirstMessageErrorContext)
		return
	}

	//这个是过滤掉最前边的~m~2~m~这个前缀，后边就是json了. 这个json字符串就是payload
	payload := msg[getPayloadStartingIndex(msg):]
	var p map[string]interface{}
	//反序列化一下
	err = json.Unmarshal(payload, &p)
	if err != nil {
		s.onError(err, DecodeFirstMessageErrorContext)
		return
	}

	//本质上就是看一下有没有session_id
	if p["session_id"] == nil {
		err = errors.New("cannot recognize the first received message after establishing the connection")
		s.onError(err, FirstMessageWithoutSessionIdErrorContext)
		return
	}

	return
}

// 要发一个
func (s *TradingViewWebSocket) generateSessionID() {
	s.sessionID = "qs_" + GetRandomString(12)
}

func (s *TradingViewWebSocket) sendConnectionSetupMessages() (err error) {

	//指定一下要获取哪些字段
	fields := []string{s.sessionID}
	fields = append(fields, FieldList...)

	messages := []*SocketMessage{
		getSocketMessage("set_auth_token", []string{string(s.authTokenType)}),
		getSocketMessage("quote_create_session", []string{s.sessionID}),

		//~m~34~m~{"m":"set_locale","p":["en","US"]}
		getSocketMessage("set_locale", []string{"en", "US"}),

		getSocketMessage("quote_set_fields", fields), //[]string{s.sessionID, "base-currency-logoid", "ch", "chp", "currency-logoid", "currency_code", "currency_id", "base_currency_id", "current_session", "description", "exchange", "format", "fractional", "is_tradable", "language", "local_description", "listed_exchange", "logoid", "lp", "lp_time", "minmov", "minmove2", "original_name", "pricescale", "pro_name", "short_name", "type", "typespecs", "update_mode", "volume", "variable_tick_size", "value_unit_id"}),
		//getSocketMessage("quote_set_fields", []string{s.sessionID, "lp", "volume", "bid", "ask"}),

		//getSocketMessage("quote_add_symbols", []string{s.sessionID, "FOREXCOM:SPXUSD", "FOREXCOM:NSXUSD", "FX_IDC:EURUSD", "BITSTAMP:BTCUSD", "BITSTAMP:ETHUSD"}),

		//~m~138~m~{"m":"quote_fast_symbols","p":["qs_JLiN50VoHqbu","FOREXCOM:SPXUSD","FOREXCOM:NSXUSD","FX_IDC:EURUSD","BITSTAMP:BTCUSD","BITSTAMP:ETHUSD"]}
		//getSocketMessage("quote_fast_symbols", []string{s.sessionID, "FOREXCOM:SPXUSD", "FOREXCOM:NSXUSD", "FX_IDC:EURUSD", "BITSTAMP:BTCUSD", "BITSTAMP:ETHUSD"}),
	}

	for _, msg := range messages {
		err = s.sendSocketMessage(msg)
		if err != nil {
			return
		}
	}

	return
}

func (s *TradingViewWebSocket) sendSocketMessage(p *SocketMessage) (err error) {
	payload, _ := json.Marshal(p)
	payloadWithHeader := "~m~" + strconv.Itoa(len(payload)) + "~m~" + string(payload)

	err = s.conn.WriteMessage(websocket.TextMessage, []byte(payloadWithHeader))
	if err != nil {
		s.onError(err, SendMessageErrorContext+" - "+payloadWithHeader)
		return
	}
	return
}

func (s *TradingViewWebSocket) connectionLoop() {
	var readMsgError error
	var writeKeepAliveMsgError error

	for readMsgError == nil && writeKeepAliveMsgError == nil {
		if s.isClosed {
			break
		}

		var msgType int
		var msg []byte
		msgType, msg, readMsgError = s.conn.ReadMessage()

		go func(msgType int, msg []byte) {
			if msgType != websocket.TextMessage {
				return
			}

			//这里是服务端发心跳包, 收到后:就原样copy回复一下
			if isKeepAliveMsg(msg) {
				writeKeepAliveMsgError = s.conn.WriteMessage(msgType, msg)
				return
			}

			go s.parsePacket(msg)
		}(msgType, msg)
	}

	if readMsgError != nil {
		s.onError(readMsgError, ReadMessageErrorContext)
	}
	if writeKeepAliveMsgError != nil {
		s.onError(writeKeepAliveMsgError, SendKeepAliveMessageErrorContext)
	}
}

// 负责解析收到的数据
func (s *TradingViewWebSocket) parsePacket(packet []byte) {
	var symbolsArr []string
	var dataArr []*QuoteData

	index := 0
	for index < len(packet) {
		//解析最前边的~m~23~m~, 这个被定义为payload
		//这里的 payloadLength 指的是中间的23， 也就是后边实际json字符串的长度
		payloadLength, err := getPayloadLength(packet[index:])
		if err != nil {
			s.onError(err, GetPayloadLengthErrorContext+" - "+string(packet))
			return
		}

		//这里的headerLength是 ~m~23~m~ 这个payLoad整体的长度
		headerLength := 6 + len(strconv.Itoa(payloadLength))
		//这个是获取具体的返回json字符串
		payload := packet[index+headerLength : index+headerLength+payloadLength]
		index = index + headerLength + len(payload)

		//解析每一个json字符串
		symbol, data, err := s.parseJSON(payload)
		if err != nil {
			break
		}

		dataArr = append(dataArr, data)
		symbolsArr = append(symbolsArr, symbol)
	}

	for i := 0; i < len(dataArr); i++ {
		isDuplicate := false
		for j := i + 1; j < len(dataArr); j++ {
			if GetStringRepresentation(dataArr[i]) == GetStringRepresentation(dataArr[j]) {
				isDuplicate = true
				break
			}
		}
		if !isDuplicate {
			s.OnReceiveMarketDataCallback(symbolsArr[i], dataArr[i])
		}
	}
}

func (s *TradingViewWebSocket) parseJSON(msg []byte) (symbol string, data *QuoteData, err error) {
	var decodedMessage *SocketMessage

	err = json.Unmarshal(msg, &decodedMessage)
	if err != nil {
		s.onError(err, DecodeMessageErrorContext+" - "+string(msg))
		return
	}

	//fmt.Printf("-----result------%+v\n", decodedMessage)

	if decodedMessage.Message == "critical_error" || decodedMessage.Message == "error" {
		err = errors.New("Error -> " + string(msg))
		s.onError(err, DecodedMessageHasErrorPropertyErrorContext)
		return
	}

	//qsd是返回的报价类型
	if decodedMessage.Message != "qsd" {
		err = errors.New("ignored message - Not QSD")
		return
	}

	//payload是具体信息
	if decodedMessage.Payload == nil {
		err = errors.New("Msg does not include 'p' -> " + string(msg))
		s.onError(err, DecodedMessageDoesNotIncludePayloadErrorContext)
		return
	}

	p, isPOk := decodedMessage.Payload.([]interface{})
	//返回的p就两个字段， 一个字段是对应的symbol, 一个是这个symbol对应的详细行情信息(是个json)
	if !isPOk || len(p) != 2 {
		err = errors.New("There is something wrong with the payload - can't be parsed -> " + string(msg))
		s.onError(err, PayloadCantBeParsedErrorContext)
		return
	}

	//把单个symbol的行情json字符串给解析一下
	var decodedQuoteMessage *QuoteMessage
	err = mapstructure.Decode(p[1].(map[string]interface{}), &decodedQuoteMessage)
	if err != nil {
		s.onError(err, FinalPayloadCantBeParsedErrorContext+" - "+string(msg))
		return
	}

	if decodedQuoteMessage.Status != "ok" || decodedQuoteMessage.Symbol == "" || decodedQuoteMessage.Data == nil {
		err = errors.New("There is something wrong with the payload - couldn't be parsed -> " + string(msg))
		s.onError(err, FinalPayloadHasMissingPropertiesErrorContext)
		return
	}
	symbol = decodedQuoteMessage.Symbol
	data = decodedQuoteMessage.Data
	return
}

func (s *TradingViewWebSocket) onError(err error, context string) {
	if s.conn != nil {
		s.conn.Close()
	}
	s.OnErrorCallback(err, context)
}

func getSocketMessage(m string, p interface{}) *SocketMessage {
	return &SocketMessage{
		Message: m,
		Payload: p,
	}
}

func getFlags() *Flags {
	return &Flags{
		Flags: []string{"force_permission"},
	}
}

// ~m~4~m~~h~1
// 内容是 ~h~1  所以只要判断第一个字符是 ~ 就说明问题了
func isKeepAliveMsg(msg []byte) bool {
	return string(msg[getPayloadStartingIndex(msg)]) == "~"
}

// ~m~34~m~{"m":"set_locale","p":["en","US"]}
func getPayloadStartingIndex(msg []byte) int {
	char := ""

	//跳过前边的~m~这3个字符
	index := 3
	for char != "~" {
		//读取中间的34这个字符
		char = string(msg[index])
		index++
	}
	index += 2
	//返回的是json,所以这里的index是json字符串起始的位置
	return index
}

func getPayloadLength(msg []byte) (length int, err error) {
	char := ""
	index := 3
	lengthAsString := ""
	for char != "~" {
		char = string(msg[index])
		if char != "~" {
			lengthAsString += char
		}
		index++
	}
	length, err = strconv.Atoi(lengthAsString)
	return
}

func getHeaders() http.Header {
	headers := http.Header{}

	headers.Set("Accept-Encoding", "gzip, deflate, br")
	headers.Set("Accept-Language", "en-US,en;q=0.9,es;q=0.8")
	headers.Set("Cache-Control", "no-cache")
	headers.Set("Host", "data.tradingview.com")
	headers.Set("Origin", "https://www.tradingview.com")
	headers.Set("Pragma", "no-cache")
	headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.193 Safari/537.36")

	return headers
}
