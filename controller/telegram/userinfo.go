package tgcontroller

type UserInfo struct {
	Uuid           string  `json:"uuid"`
	Username       string  `json:"username"`
	ChatID         int     `json:"chat_id"`
	SubscribedTags *TagMap `json:"subscribed_tags"`
	AddingTags     bool    `json:"addingTags"`
	ManagerMessage int     `json:"managerMessage"`
}

func (u *UserInfo) GetUuid() string {
	return u.Uuid
}
