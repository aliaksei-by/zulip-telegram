package main

type Zulip struct {
	Key   string `yaml:"key"`
	Email string `yaml:"email"`
	Site  string `yaml:"site"`
}

type User struct {
	Name              string   `yaml:"name"`
	ID                int64    `yaml:"tg-id"`
	ZulipEmail        string   `yaml:"zulip-email"`
	ZulipKey          string   `yaml:"zulip-key"`
	Words             []string `yaml:"words"`
	IgnorePrivateFrom []string `yaml:"ignore-private-from"`
}

type Config struct {
	TGKey     string `yaml:"tg-key"`
	BotUserID int64  `yaml:"bot-user-id"`
	Zulip     Zulip  `yaml:"zulip"`
	Users     []User `yaml:"users"`
}
