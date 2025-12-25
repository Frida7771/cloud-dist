package config

// Config defines runtime settings for the Gin application.
type Config struct {
	Name     string       `mapstructure:"Name"`
	Host     string       `mapstructure:"Host"`
	Port     int          `mapstructure:"Port"`
	MaxBytes int64        `mapstructure:"MaxBytes"`
	Mysql    MysqlConfig  `mapstructure:"Mysql"`
	Redis    RedisConfig  `mapstructure:"Redis"`
	Log      LogConfig    `mapstructure:"Log"`
	HTTP     HTTPSettings `mapstructure:"HTTP"`
}

// MysqlConfig carries datasource info for the ORM layer.
type MysqlConfig struct {
	DataSource string `mapstructure:"DataSource"`
}

// RedisConfig carries Redis connection info.
type RedisConfig struct {
	Addr string `mapstructure:"Addr"`
}

// LogConfig mirrors the old go-zero log section to avoid surprises.
type LogConfig struct {
	ServiceName string `mapstructure:"ServiceName"`
	Mode        string `mapstructure:"Mode"`
	Path        string `mapstructure:"Path"`
	Level       string `mapstructure:"Level"`
	Encoding    string `mapstructure:"Encoding"`
	KeepDays    int    `mapstructure:"KeepDays"`
}

// HTTPSettings allows future Gin-specific tuning; optional in config files.
type HTTPSettings struct {
	ReadTimeoutSeconds  int `mapstructure:"ReadTimeoutSeconds"`
	WriteTimeoutSeconds int `mapstructure:"WriteTimeoutSeconds"`
}
