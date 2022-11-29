package record

type SendMessageResponseData struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	MessageId int    `json:"messageId"`
}

func (s *SendMessageResponseData) getType() Type {
	return SendMessageResponse
}

func GetSendMessageResponse() Base {
	return &BaseImp{Data: &SendMessageResponseData{}}
}
