package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// https://github.com/binance/binance-spot-api-docs/blob/master/rest-api.md
// https://github.com/binance/binance-spot-api-docs/blob/master/rest-api_CN.md
type KLine struct {
	OpenTime   uint64 `arrayIndex:"0"`
	OpenPrice  string `arrayIndex:"1"`
	HighPrice  string `arrayIndex:"2"`
	LowPrice   string `arrayIndex:"3"`
	ClosePrice string `arrayIndex:"4"`
	Volume     string `arrayIndex:"5"` // 成交量
	CloseTime  uint64 `arrayIndex:"6"` // 收盘时间

	/*
		ChatGPT:
		Quote asset volume is a metric used in cryptocurrency trading that measures the amount of a particular quote currency that is traded in a specific time period, typically 24 hours. The quote asset is the currency in which the price of a cryptocurrency is denominated, and the volume refers to the total value of trades executed in that currency over a given period of time.

		For example, if the quote asset is US dollars (USD), the quote asset volume would represent the total value of all trades executed in USD in the past 24 hours. This metric is useful for traders and investors to assess the liquidity of a particular cryptocurrency on an exchange, as well as to identify trends and trading opportunities.
	*/
	QuoteAssetVolume       string `arrayIndex:"7"`  // 成交额
	NumberOfTrades         uint64 `arrayIndex:"8"`  // 成交笔数
	BoughtBaseAssetVolume  string `arrayIndex:"9"`  // 主动买入成交量
	BoughtQuoteAssetVolume string `arrayIndex:"10"` // 主动买入成交额

	// Unused         string // ignore this
}

type KLineRequest struct {
	PairSymbol string `http:"symbol"`
	Interval   string `http:"interval"`
	StartTime  uint   `http:"startTime"`
	EndTime    uint   `http:"endTime"`
	Limit      uint   `http:"limit"` // default 500; max 1000

	// test int `http:"test,default=100"`
}

type KLineRequestOption func(klr *KLineRequest)

func WithStartTime(startTime uint) KLineRequestOption {
	return func(klr *KLineRequest) {
		klr.StartTime = startTime
	}
}

func WithEndTime(endTime uint) KLineRequestOption {
	return func(klr *KLineRequest) {
		klr.EndTime = endTime
	}
}

// TODO: struct tag 标识字段格式要求并校验
func NewKLineReq(pairSymbol string, interval KLineInterval, limit uint, options ...KLineRequestOption) *KLineRequest {
	klr := &KLineRequest{
		PairSymbol: pairSymbol,
		Interval:   interval.String(),
		Limit:      limit,
	}
	for _, opt := range options {
		opt(klr)
	}

	return klr
}

func KLinesURL(klr *KLineRequest) string {
	return fmt.Sprintf("%s%s%s", BaseURL, KLINES, urlSearchParamsFromKLineRequest(klr))
}

func urlSearchParamsFromKLineRequest(klr *KLineRequest) string {
	var urlParams []string
	v := reflect.ValueOf(klr).Elem()
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i)
		if fieldInfo.IsExported() {
			var fieldValueString string
			// not working: fmt.Println("zero value:", reflect.DeepEqual(v.Field(i).Interface(), reflect.Zero(v.Field(i).Type())))
			if v.Field(i).Interface() == reflect.Zero(v.Field(i).Type()).Interface() {
				// fmt.Printf("field %s got zero value\n", fieldInfo.Name)
				continue
			}
			switch v.Field(i).Kind() {
			case reflect.String:
				fieldValueString = v.Field(i).String()
			case reflect.Uint:
				fieldValueString = strconv.FormatUint(v.Field(i).Uint(), 10)
			}

			var fieldNameInTag string
			tag := fieldInfo.Tag
			fn, found := tag.Lookup("http")
			if found {
				fieldNameInTag = fn
			} else {
				fieldNameInTag = fieldInfo.Name
			}
			urlParams = append(urlParams, fmt.Sprintf("%s=%s", fieldNameInTag, fieldValueString))
		} else {
			// fmt.Printf("field %s is unexported, got tag: %q\n", fieldInfo.Name, fieldInfo.Tag.Get("http"))
		}
	}

	if len(urlParams) == 0 {
		return ""
	}

	return fmt.Sprintf("?%s", strings.Join(urlParams, "&"))
}

func UnmarshalKLinesJSON(data []byte) ([]*KLine, error) {
	var raws [][]json.RawMessage
	if err := json.Unmarshal(data, &raws); err != nil {
		return nil, err
	}

	for _, raw := range raws {
		unmarshalKLineJSON(raw)
	}

	return nil, nil
}

// TODO maybe cache mapping
func unmarshalKLineJSON(raw []json.RawMessage) {
	// raw is a reflect.Slice
	rv := reflect.ValueOf(raw)
	// rv.Index(i).Kind() is a reflect.Slice
	// fmt.Println(rv.Kind(), rv.Len(), rv.Index(11).Kind())

	var kLine KLine
	v := reflect.ValueOf(&kLine).Elem()
	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag
		idx, err := strconv.Atoi(tag.Get("arrayIndex"))
		if err != nil {
			fmt.Println(err)
		}
		fd := v.Field(idx)
		switch fd.Kind() {
		case reflect.Uint64:
			var u uint64
			if err := json.Unmarshal(rv.Index(idx).Bytes(), &u); err != nil {
				fmt.Println(err)
			}
			// fmt.Println(u)
			fd.Set(reflect.ValueOf(u))
		case reflect.String:
			var s string
			if err := json.Unmarshal(rv.Index(idx).Bytes(), &s); err != nil {
				fmt.Println(err)
			}
			// fmt.Println(s)
			fd.Set(reflect.ValueOf(s))
		default:
			panic("should not happen")
		}

		/*
			var fdv interface{}
			if err := json.Unmarshal(rv.Index(idx).Bytes(), &fdv); err != nil {
				fmt.Println(err)
			}
			fd.Set(reflect.ValueOf(fdv)) // uint64 会被当作 float64 处理，导致无法 Set（value of type float64 is not assignable to type uint64）
		*/

		// json.Unmarshal 拿不到第二个参数，能拿到 uintptr 值；
		// 可以尝试用 unsafe Offset
	}
	fmt.Printf("%#v\nopen time: %d, close time: %d, volume: %d\n", kLine, kLine.OpenTime, kLine.CloseTime, kLine.NumberOfTrades)
}
