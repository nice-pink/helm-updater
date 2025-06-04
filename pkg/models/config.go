package models

type Config struct {
	Apps       []App
	BaseFolder string
	Notify     Notication
}

type Notication struct {
	Enable    bool
	ChannelId string
	Webhook   string
}
