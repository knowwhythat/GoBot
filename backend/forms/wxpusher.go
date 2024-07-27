package forms

type GenerateCode struct {
	AppToken  string `json:"appToken"`
	Extra     string `json:"extra"`
	ValidTime int64  `json:"validTime"`
}

type WxPusherResp[T any] struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Data    T      `json:"data"`
	Success bool   `json:"success"`
}

type CodeData struct {
	Expires  int64  `json:"expires"`
	Code     string `json:"code"`
	ShortURL string `json:"shortUrl"`
	Extra    string `json:"extra"`
	URL      string `json:"url"`
}
