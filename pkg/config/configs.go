package configs

type (
	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Host string `env-required:"true" yaml:"host" env:"HTTP_HOST"`
		Port int    `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	HTTP2 struct {
		Host string `env-required:"true" yaml:"host" env:"HTTP_HOST_2"`
		Port int    `env-required:"true" yaml:"port" env:"HTTP_PORT_2"`
	}

	Log struct {
		Level             string `env-required:"true" yaml:"log_level" env:"LOG_LEVEL"`
		Dev               bool   `env-required:"true" yaml:"log_dev" env:"LOG_DEV"`
		DisableCaller     bool   `env-required:"true" yaml:"log_disable_caller" env:"LOG_DISABLE_CALLER"`
		DisableStacktrace bool   `env-required:"true" yaml:"log_disable_stacktrace" env:"LOG_DISABLE_STACKTRACE"`
		Encoding          string `env-required:"true" yaml:"log_encoding" env:"LOG_ENCODING"`
	}
	Metrics struct {
		HostMetric string `env-required:"true" yaml:"host" env:"HTTP_HOST_METRIC"`
		PortMetric int    `env-required:"true" yaml:"port" env:"HTTP_PORT_METRIC"`
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
	Jeager struct {
		HostJeager        string `env-required:"true" yaml:"host_jaeger" env:"HOST_JAEGER"`
		ServiceNameJeager string `env-required:"true" yaml:"service_name_jaeger" env:"SERVICE_NAME_JAEGER"`
		LogSpans          bool   `env-required:"true" yaml:"log_spans" env:"LOG_SPANS"`
	}
)
