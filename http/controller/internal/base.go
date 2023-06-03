package internal

type Response struct {
	Ok   bool        `json:"ok"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}
