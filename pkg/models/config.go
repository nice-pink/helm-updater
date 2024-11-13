package models

type Config struct {
	Apps       []App
	BaseFolder string
	Notify     Notication
}

type Notication struct {
	Enable  bool
	Webhook string
}
