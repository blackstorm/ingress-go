package config

type Config struct {
	IsShowBanner bool
}

func InitFromFlags() Config {
	return Config{
		IsShowBanner: true,
	}
}
