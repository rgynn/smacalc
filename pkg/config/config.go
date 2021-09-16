package config

import (
	"errors"
	"flag"
	"strings"
)

type Data struct {
	Count       int
	APIKey      string
	SymbolNames []string
}

var (
	n           = flag.Int("n", 0, "n count for sma calculations of symbols")
	apiKey      = flag.String("apikey", "", "apikey for finnhub.io")
	symbolsFlag = flag.String("symbols", "", "comma separated list of symbols from finnhub.io, example: BINANCE:BTCUSDT,BINANCE:ETHUSDT,BINANCE:ADAUSDT")
)

func init() {
	flag.Parse()
}

func NewFromFlags() (*Data, error) {
	if *n < 1 {
		return nil, errors.New("flag -n needs to be atleast 1")
	}
	if *apiKey == "" {
		return nil, errors.New("flag -apikey not provided")
	}
	symbolNames := strings.Split(*symbolsFlag, ",")
	if len(symbolNames) < 1 {
		return nil, errors.New("flag -symbols not provided, or incorrectly formatted")
	}
	return &Data{
		Count:       *n,
		APIKey:      *apiKey,
		SymbolNames: symbolNames,
	}, nil
}
