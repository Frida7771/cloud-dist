package config

// Config defines runtime settings for the Gin application.
type Config struct {
	Name     string         `mapstructure:"Name"`
	Host     string         `mapstructure:"Host"`
	Port     int            `mapstructure:"Port"`
	MaxBytes int64          `mapstructure:"MaxBytes"`
	Mysql    MysqlConfig    `mapstructure:"Mysql"`
	Redis    RedisConfig    `mapstructure:"Redis"`
	Log      LogConfig      `mapstructure:"Log"`
	HTTP     HTTPSettings   `mapstructure:"HTTP"`
	S3       S3Config       `mapstructure:"S3"`
	SendGrid SendGridConfig `mapstructure:"SendGrid"`
	JWT      JWTConfig      `mapstructure:"JWT"`
	Stripe   StripeConfig   `mapstructure:"Stripe"`
}

// JWTConfig carries JWT configuration.
type JWTConfig struct {
	Key string `mapstructure:"Key"`
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

// S3Config carries AWS S3 configuration.
type S3Config struct {
	AccessKeyID     string `mapstructure:"AccessKeyID"`
	SecretAccessKey string `mapstructure:"SecretAccessKey"`
	Bucket          string `mapstructure:"Bucket"`
	Region          string `mapstructure:"Region"`
	Endpoint        string `mapstructure:"Endpoint"` // Optional custom endpoint
}

// SendGridConfig carries SendGrid email configuration.
type SendGridConfig struct {
	APIKey    string `mapstructure:"APIKey"`
	FromEmail string `mapstructure:"FromEmail"` // Optional sender email
}

// StripeConfig carries Stripe payment configuration.
type StripeConfig struct {
	SecretKey     string `mapstructure:"SecretKey"`
	WebhookSecret string `mapstructure:"WebhookSecret"`
}
