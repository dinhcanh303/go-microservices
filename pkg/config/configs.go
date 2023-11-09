package configs

type (
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Host string `env-required:"true" yaml:"host" env:"HTTP_HOST"`
		Port int    `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	HTTPEcho struct {
		HostEcho string `env-required:"true" yaml:"host" env:"HTTP_HOST_ECHO"`
		PortEcho int    `env-required:"true" yaml:"port" env:"HTTP_PORT_ECHO"`
	}
	Request struct {
		RequestPerSecond int `env-required:"true" yaml:"request_per_second" env:"REQUEST_PER_SECOND"`
		RequestBurst     int `env-required:"true" yaml:"request_burst" env:"REQUEST_BURST"`
		RequestMax       int `env-required:"true" yaml:"request_max" env:"REQUEST_MAX"`
		DurationsSecond  int `env-required:"true" yaml:"durations_second" env:"DURATIONS_SECOND"`
	}
)
