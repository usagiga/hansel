package config

type Config struct {
	StartTriggerMessage     string
	HibernateTriggerMessage string
	GetStatusTriggerMessage string
}

func GetConfig() (Config, error) {
	config := Config{
		StartTriggerMessage:     "start",
		HibernateTriggerMessage: "sleep",
		GetStatusTriggerMessage: "status",
	}
	return config, nil
}
