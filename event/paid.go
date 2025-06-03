package event

import "time"

type Paid struct {
	UserID    string    `json:"user_id"`
	CId       int64     `json:"cid"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}

type PaidEvent struct {
	Account struct {
		Id              int64     `json:"id"`
		Uid             string    `json:"uid"`
		Cid             int64     `json:"cid"`
		PaymentMethod   string    `json:"payment_method"`
		PaymentMethodId string    `json:"payment_method_id"`
		Date            time.Time `json:"date"`
		Description     string    `json:"description"`
		Currency        string    `json:"currency"`
		AmountIn        string    `json:"amount_in"`
		AmountOut       string    `json:"amount_out"`
		TransId         string    `json:"trans_id"`
		InvoiceId       string    `json:"invoice_id"`
		RefundId        string    `json:"refund_id"`
		Status          string    `json:"status"`
		Valid           string    `json:"valid"`
		PaymentId       string    `json:"payment_id"`
		Country         string    `json:"country"`
		CardBrand       string    `json:"cardBrand"`
		AutoRenewal     int64     `json:"autoRenewal"`
	} `json:"account"`
	GoodItems []struct {
		GoodsId       string `json:"goods_id"`
		GoodsName     string `json:"goods_name"`
		GoodsCategory string `json:"goods_category"`
		PeriodType    string `json:"period_type"`
	} `json:"good_items"`
	User struct {
		Id                     int64     `json:"Id"`
		Uid                    string    `json:"Uid"`
		Cid                    int64     `json:"Cid"`
		FirstName              string    `json:"FirstName"`
		MiddleName             string    `json:"MiddleName"`
		LastName               string    `json:"LastName"`
		FullName               string    `json:"FullName"`
		Email                  string    `json:"Email"`
		Country                string    `json:"Country"`
		DefaultPaymentMethodId string    `json:"DefaultPaymentMethodId"`
		GmtCreate              time.Time `json:"GmtCreate"`
		GmtUpdate              time.Time `json:"GmtUpdate"`
	} `json:"user"`
}
