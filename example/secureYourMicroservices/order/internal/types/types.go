package types

type Order struct {
	ID              string     `json:"id"`
	Items           []LineItem `json:"items"`
	ShippingAddress string     `json:"shipping_address"`
}

type LineItem struct {
	ItemCode string `json:"item_code"`
	Quantity int    `json:"quantity"`
}
