package main

import (
	"fmt"
	"testing"
)

func TestURLSearchParamsFromKLineRequest(t *testing.T) {
	klrs := []*KLineRequest{
		{
			PairSymbol: "ETHUSDT",
			Interval:   "1d",
			Limit:      3,
		},
	}
	for _, klr := range klrs {
		caredResult := urlSearchParamsFromKLineRequest(klr)
		simpleResult := urlSearchParamsFromKLineRequestSimple(klr)
		if caredResult != simpleResult {
			t.Errorf("should be equal for %#v, cared one is %s, simple one is %s\n", *klr, caredResult, simpleResult)
		}
	}
}

func urlSearchParamsFromKLineRequestSimple(klr *KLineRequest) string {
	params := fmt.Sprintf("?symbol=%s&interval=%s&limit=%d", klr.PairSymbol, klr.Interval, klr.Limit)
	if klr.StartTime != 0 {
		params += fmt.Sprintf("&startTime=%d", klr.StartTime)
	}
	if klr.EndTime != 0 {
		params += fmt.Sprintf("&endTime=%d", klr.EndTime)
	}
	return params
}
