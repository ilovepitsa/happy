package config

type (
	App struct {
		Name string `yaml="name"`
	}

	Network struct {
		Host string `yaml="host"`
		Port string `yaml="port"`
	}

	Logger struct {
		JSONEnable bool   `yaml="json_enable"`
		Level      string `yaml="level"`
	}

	Config struct {
		App
	}
)
