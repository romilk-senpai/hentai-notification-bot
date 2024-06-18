package tgcontroller

type UserInfo struct {
	Uuid           string  `json:"uuid"`
	Username       string  `json:"username"`
	ChatID         int     `json:"chat_id"`
	SubscribedTags *TagMap `json:"subscribed_tags"`
}

func (u UserInfo) GetUuid() string {
	return u.Uuid
}
