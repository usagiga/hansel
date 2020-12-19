package config

type Config struct {
	StartTriggerMessage     string
	HibernateTriggerMessage string
	GetStatusMessage        string
}

func GetConfig() (Config, error) {
	config := Config{
		StartTriggerMessage:     "start",
		HibernateTriggerMessage: "sleep",
		GetStatusMessage:        "status",
	}
	return config, nil
}
