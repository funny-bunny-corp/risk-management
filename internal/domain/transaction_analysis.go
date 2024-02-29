package domain

import "time"

type TransactionAnalysis struct {
	Participants Participants `json:"participants"`
	Order        Checkout     `json:"order"`
	Payment      Payment      `json:"payment"`
}

type Checkout struct {
	Id          string    `json:"id"`
	PaymentType CardInfo  `json:"payment_type"`
	At          time.Time `json:"at"`
}

type Payment struct {
	Id       string `json:"id"`
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
	Status   string `json:"status"`
}

type Participants struct {
	Buyer  BuyerInfo  `json:"buyer"`
	Seller SellerInfo `json:"seller"`
}

type BuyerInfo struct {
	Document string `json:"document"`
	Name     string `json:"name"`
}

type SellerInfo struct {
	SellerId string `json:"seller_id"`
}

type CardInfo struct {
	CardInfo string `json:"card_info"`
	Token    string `json:"token"`
}
