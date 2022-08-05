package entity

type IndividualSender struct {
	Sender
	Nickname string `json:"nickname,omitempty"`
	Remark   string `json:"remark,omitempty"`
}
