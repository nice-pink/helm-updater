package models

type Config struct {
	Apps       []App
	BaseFolder string
	Notify     Notication
	Helm       Helm
}

type Notication struct {
	Enable    bool
	Webhook   string
	ChannelId string
	Token     string
}

type Helm struct {
	CachePath    string
	RepoFilePath string
	CleanUp      bool
}
