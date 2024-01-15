package stakan

type StakanDataStruct struct {
	LastUpdateID int64      `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
	Timestamp    int64      `json:"timestamp"`
}
