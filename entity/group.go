package entity

type Group struct {
	ID         uint64 `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Permission string `json:"permission,omitempty"`
}

type GroupSender struct {
	Sender
	MemberName         string `json:"memberName,omitempty"`
	SpecialTitle       string `json:"specialTitle,omitempty"`
	Permission         string `json:"permission,omitempty"`
	JoinTimestamp      uint64 `json:"joinTimestamp,omitempty"`
	LastSpeakTimestamp uint64 `json:"lastSpeakTimestamp,omitempty"`
	MuteTimeRemaining  uint64 `json:"muteTimeRemaining,omitempty"`
	Group              Group  `json:"group,omitempty"`
}
