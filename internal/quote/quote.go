package quote

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

type ResponseQuote struct {
	ShortName                  string  `json:"shortName"`
	Symbol                     string  `json:"symbol"`
	MarketState                string  `json:"marketState"`
	RegularMarketChange        float64 `json:"regularMarketChange"`
	RegularMarketChangePercent float64 `json:"regularMarketChangePercent"`
	RegularMarketPrice         float64 `json:"regularMarketPrice"`
	PostMarketChange           float64 `json:"postMarketChange"`
	PostMarketChangePercent    float64 `json:"postMarketChangePercent"`
	PostMarketPrice            float64 `json:"postMarketPrice"`
	PreMarketChange            float64 `json:"preMarketChange"`
	PreMarketChangePercent     float64 `json:"preMarketChangePercent"`
	PreMarketPrice             float64 `json:"preMarketPrice"`
}

type Quote struct {
	ResponseQuote
	Price                   float64
	Change                  float64
	ChangePercent           float64
	IsActive                bool
	IsRegularTradingSession bool
}

type Response struct {
	QuoteResponse struct {
		Quotes []ResponseQuote `json:"result"`
		Error  interface{}     `json:"error"`
	} `json:"quoteResponse"`
}

func transformResponseQuote(responseQuote ResponseQuote) Quote {

	if responseQuote.MarketState == "REGULAR" {
		return Quote{
			ResponseQuote:           responseQuote,
			Price:                   responseQuote.RegularMarketPrice,
			Change:                  responseQuote.RegularMarketChange,
			ChangePercent:           responseQuote.RegularMarketChangePercent,
			IsActive:                true,
			IsRegularTradingSession: true,
		}
	}

	if responseQuote.MarketState == "POST" {
		return Quote{
			ResponseQuote:           responseQuote,
			Price:                   responseQuote.PostMarketPrice + responseQuote.RegularMarketPrice,
			Change:                  responseQuote.PostMarketChange + responseQuote.RegularMarketChange,
			ChangePercent:           responseQuote.PostMarketChangePercent + responseQuote.RegularMarketChangePercent,
			IsActive:                true,
			IsRegularTradingSession: false,
		}
	}

	if responseQuote.MarketState == "PRE" {
		return Quote{
			ResponseQuote:           responseQuote,
			Price:                   responseQuote.PreMarketPrice,
			Change:                  responseQuote.PreMarketChange,
			ChangePercent:           responseQuote.PreMarketChangePercent,
			IsActive:                true,
			IsRegularTradingSession: false,
		}
	}

	// temporary for testing
	// if responseQuote.MarketState == "CLOSED" {
	// 	return Quote{
	// 		ResponseQuote:           responseQuote,
	// 		Price:                   responseQuote.RegularMarketPrice,
	// 		Change:                  responseQuote.RegularMarketChange,
	// 		ChangePercent:           responseQuote.RegularMarketChangePercent,
	// 		IsActive:                true,
	// 		IsRegularTradingSession: true,
	// 	}
	// }

	return Quote{
		ResponseQuote:           responseQuote,
		Price:                   responseQuote.RegularMarketPrice,
		Change:                  0.0,
		ChangePercent:           0.0,
		IsActive:                false,
		IsRegularTradingSession: false,
	}

}

func transformResponseQuotes(responseQuotes []ResponseQuote) []Quote {

	quotes := make([]Quote, 0)
	for _, responseQuote := range responseQuotes {
		quotes = append(quotes, transformResponseQuote(responseQuote))
	}
	return quotes

}

func GetQuotes(symbols []string) []Quote {
	symbolsString := strings.Join(symbols, ",")
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v7/finance/quote?lang=en-US&region=US&corsDomain=finance.yahoo.com&symbols=%s", symbolsString)
	response, _ := resty.New().R().
		SetResult(&Response{}).
		Get(url)

	return transformResponseQuotes((response.Result().(*Response)).QuoteResponse.Quotes)
}
