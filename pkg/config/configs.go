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
	Minio struct {
		EndPoint        string `env-required:"true" yaml:"end_point" env:"MINIO_ENDPOINT"`
		AccessKeyID     string `env-required:"true" yaml:"access_key_id" env:"MINIO_ACCESS_KEY_ID"`
		SecretAccessKey string `env-required:"true" yaml:"secret_access_key" env:"MINIO_SECRET_ACCESS_KEY"`
		Location        string `env-required:"true" yaml:"location" env:"MINIO_DEFAULT_REGION"`
		BucketName      string `env-required:"true" yaml:"bucket_name" env:"MINIO_BUCKET"`
		RootFolder      string `env-required:"true" yaml:"root_folder" env:"MINIO_ROOT_FOLDER"`
		UseSSL          bool   `env-required:"true" yaml:"use_ssl" env:"MINIO_USE_SSL"`
	}
)
