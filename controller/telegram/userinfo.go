package tgcontroller

type UserInfo struct {
	Uuid             string  `json:"uuid"`
	Username         string  `json:"username"`
	ChatID           int     `json:"chat_id"`
	SubscribedTags   *TagMap `json:"subscribed_tags"`
	AddingTags       bool    `json:"addingTags"`
	ManagerMessageID int     `json:"managerMessage"`
}

func NewUserInfo(uuid string, username string, chatID int, subscribedTags *TagMap) *UserInfo {
	return &UserInfo{
		Uuid:             uuid,
		Username:         username,
		ChatID:           chatID,
		SubscribedTags:   subscribedTags,
		AddingTags:       false,
		ManagerMessageID: -1,
	}
}

func (u *UserInfo) GetUuid() string {
	return u.Uuid
}
