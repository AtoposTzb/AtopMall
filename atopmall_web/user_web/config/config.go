package config

type UserSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"Key" json:"Key"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Password string `mapstructure:"password" json:"password"`
	DB       int    `mapstructure:"db" json:"db"`
}

type EmailConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
	Expires  int    `mapstructure:"expires" json:"expires"`
}
type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ServerConfig struct {
	Name        string        `mapstructure:"name" json:"name"`
	Port        int           `mapstructure:"port" json:"port"`
	UserSrvInfo UserSrvConfig `mapstructure:"user_srv" json:"user_srv"`
	JWTInfo     JWTConfig     `mapstructure:"jwt" json:"jwt"`
	RedisInfo   RedisConfig   `mapstructure:"redis" json:"redis"`
	EmailInfo   EmailConfig   `mapstructure:"email" json:"email"`
	ConsulInfo  ConsulConfig  `mapstructure:"consul" json:"consul"`
}

type NacosConfig struct {
	Host                string `mapstructure:"host"`
	Port                uint64 `mapstructure:"port"`
	NamespaceId         string `mapstructure:"namespaceId"`
	TimeoutMs           uint64 `mapstructure:"timeoutMs"`
	NotLoadCacheAtStart bool   `mapstructure:"notLoadCacheAtStart"`
	LogDir              string `mapstructure:"logDir"`
	CacheDir            string `mapstructure:"cacheDir"`
	LogLevel            string `mapstructure:"logLevel"`
	User                string `mapstructure:"user"`
	Password            string `mapstructure:"password"`
	Dataid              string `mapstructure:"dataid"`
	Group               string `mapstructure:"group"`
}
