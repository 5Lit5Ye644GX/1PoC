package main

type Transaction struct {
	Error  interface{} `json:"error"`
	ID     string      `json:"id"`
	Result []struct {
		Addresses []string `json:"addresses"`
		Balance   struct {
			Amount int `json:"amount"`
			Assets []struct {
				Assetref string  `json:"assetref"`
				Name     string  `json:"name"`
				Qty      float64 `json:"qty"`
			} `json:"assets"`
		} `json:"balance"`
		Blockhash     string        `json:"blockhash"`
		Blockindex    int           `json:"blockindex"`
		Blocktime     int           `json:"blocktime"`
		Confirmations int           `json:"confirmations"`
		Data          []interface{} `json:"data"`
		Items         []interface{} `json:"items"`
		Myaddresses   []string      `json:"myaddresses"`
		Permissions   []interface{} `json:"permissions"`
		Time          int           `json:"time"`
		Timereceived  int           `json:"timereceived"`
		Txid          string        `json:"txid"`
		Valid         bool          `json:"valid"`
	} `json:"result"`
}

type Balance struct {
	Addresses   string
	Balances    float64
	CurrentTime string
}
