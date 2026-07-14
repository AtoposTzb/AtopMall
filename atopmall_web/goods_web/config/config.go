package config

type GoodsSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type JWTConfig struct {
	SigningKey        string `mapstructure:"Key" json:"Key"`
	LoginExpireHour   int    `mapstructure:"login-expire-hour" json:"login-expire-hour"`
	RefreshExpireHour int    `mapstructure:"refresh-expire-hour" json:"refresh-expire-hour"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ServerConfig struct {
	Name         string         `mapstructure:"name" json:"name"`
	Host         string         `mapstructure:"host" json:"host"`
	Port         int            `mapstructure:"port" json:"port"`
	Tags         []string       `mapstructure:"tags" json:"tags"`
	GoodsSrvInfo GoodsSrvConfig `mapstructure:"goods_srv" json:"goods_srv"`
	JWTInfo      JWTConfig      `mapstructure:"jwt" json:"jwt"`
	ConsulInfo   ConsulConfig   `mapstructure:"consul" json:"consul"`
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
