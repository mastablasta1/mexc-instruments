package symbols

type SymbolsResponse struct {
	Data      []string `json:"data"`
	Code      int      `json:"code"`
	Msg       string   `json:"msg"`
	Timestamp int64    `json:"timestamp"`
}
