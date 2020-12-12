package config

type Config struct {
	StartTriggerMessage     string
	HibernateTriggerMessage string
}

func GetConfig() (Config, error) {
	config := Config{
		StartTriggerMessage:     "start",
		HibernateTriggerMessage: "sleep",
	}
	return config, nil
}
