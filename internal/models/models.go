package models

type Currencies struct {
	Code string `json:"code,omitempty"`
	Data Data   `json:"data,omitempty"`
}
type Data struct {
	Time             int    `json:"time,omitempty"`
	Symbol           string `json:"symbol,omitempty"`
	Buy              string `json:"buy,omitempty"`
	Sell             string `json:"sell,omitempty"`
	ChangeRate       string `json:"changeRate,omitempty"`
	ChangePrice      string `json:"changePrice,omitempty"`
	High             string `json:"high,omitempty"`
	Low              string `json:"low,omitempty"`
	Vol              string `json:"vol,omitempty"`
	VolValue         string `json:"volValue,omitempty"`
	Last             string `json:"last,omitempty"`
	AveragePrice     string `json:"averagePrice,omitempty"`
	TakerFeeRate     string `json:"takerFeeRate,omitempty"`
	MakerFeeRate     string `json:"makerFeeRate,omitempty"`
	TakerCoefficient string `json:"takerCoefficient,omitempty"`
	MakerCoefficient string `json:"makerCoefficient,omitempty"`
}

type ResponseData struct {
	Time   string `json:"time,omitempty"`
	Symbol string `json:"symbol,omitempty"`
	Buy    string `json:"buy,omitempty"`
}

type Stat struct {
	FirstQuery        string `json:"firstQuery,omitempty"`
	CountOfAllQueries int    `json:"countOfAllQueries,omitempty"`
}
