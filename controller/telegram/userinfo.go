package tgcontroller

type UserInfo struct {
	Uuid           string                    `json:"uuid"`
	Username       string                    `json:"username"`
	ChatID         int                       `json:"chat_id"`
	SubscribedTags map[string]map[string]int `json:"subscribed_tags"`
}

func (m UserInfo) GetUuid() string {
	return m.Uuid
}
