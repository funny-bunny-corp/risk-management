package domain

type ScoringResult struct {
	Score       ScoreCard           `json:"score"`
	Transaction TransactionAnalysis `json:"transaction"`
}

type ScoreCard struct {
	ValueScore        ValueScoreCard        `json:"valueScore"`
	SellerScore       SellerScoreCard       `json:"sellerScore"`
	AverageValueScore AverageValueScoreCard `json:"averageValueScore"`
	CurrencyScore     CurrencyScoreCard     `json:"currencyScore"`
}

type Transaction struct {
	Id string `json:"id"`
}

type ValueScoreCard struct {
	Score int `json:"score"`
}

func (vsc *ValueScoreCard) Value() int {
	return vsc.Score * 10
}

type SellerScoreCard struct {
	Score int `json:"score"`
}

func (ssc *SellerScoreCard) Value() int {
	return ssc.Score * 5
}

type AverageValueScoreCard struct {
	Score int `json:"score"`
}

func (avsc *AverageValueScoreCard) Value() int {
	return avsc.Score * 3
}

type CurrencyScoreCard struct {
	Score int `json:"score"`
}

func (csc *CurrencyScoreCard) Value() int {
	return csc.Score * 1
}
