package config

var (
	App AppConfig
	DB  DBConfig
	MQ  MQConfig
)

type AppConfig struct {
	Env  string `mapstructure:"env"`
	Name string `mapstructure:"name"`
	Log  struct {
		Path    string `mapstructure:"path"`
		MaxSize int    `mapstructure:"max_size"`
		MaxAge  int    `mapstructure:"max_age"`
	} `mapstructure:"log"`
	Pagination struct {
		DefaultSize int `mapstructure:"default_size"`
		MaxSize     int `mapstructure:"max_size"`
	} `mapstructure:"pagination"`
	Port                string `mapstructure:"port"`
	UpRate              int    `mapstructure:"up_rate"`
	DownRate            int    `mapstructure:"down_rate"`
	SpeedUpTime         string `mapstructure:"speed_up_cron"`
	EventBridgeUrl      string `mapstructure:"event_bridge_url"`
	AESKeyStr           string `mapstructure:"aes_key"`
	AESKEY              []byte
	SpeedUpInterval     int     `mapstructure:"speed_up_interval"`
	ExpectedSuccessRate float64 `mapstructure:"expected_success_rate"`
	SlowDownInterval    int     `mapstructure:"slow_down_interval"`
}

type DBConfig struct {
	Master DBConfigOptions `mapstructure:"master"`
	Slave  DBConfigOptions `mapstructure:"slave"`
}

type MQConfig struct {
	CAFilePath       string `mapstructure:"ca_file_path"`
	CerFilePath      string `mapstructure:"cer_file_path"`
	KeyFilePath      string `mapstructure:"key_file_path"`
	Url              string `mapstructure:"url"`
	PasswordFilePath string `mapstructure:"password_file_path"`
}

type DBConfigOptions struct {
	Type             string `mapstructure:"type"`
	Dsn              string `mapstructure:"dsn"`
	PasswordFilePath string `mapstructure:"password_file_path"`
	MaxOpen          int    `mapstructure:"max_open"`
	MaxIdle          int    `mapstructure:"max_idle"`
	MaxLifeTime      int    `mapstructure:"max_life_time"`
}
