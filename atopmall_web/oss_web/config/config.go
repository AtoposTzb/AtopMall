package config

type JWTConfig struct {
	SigningKey        string `mapstructure:"Key" json:"Key"`
	LoginExpireHour   int    `mapstructure:"login-expire-hour" json:"login-expire-hour"`
	RefreshExpireHour int    `mapstructure:"refresh-expire-hour" json:"refresh-expire-hour"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type MinIOOssConfig struct {
	Endpoint     string `mapstructure:"endpoint" json:"endpoint"`
	AccessKey    string `mapstructure:"accessKey" json:"accessKey"`
	SecretKey    string `mapstructure:"secretKey" json:"secretKey"`
	UseSSL       bool   `mapstructure:"useSSL" json:"useSSL"`
	BucketName   string `mapstructure:"bucketName" json:"bucketName"`
	PublicPrefix string `mapstructure:"publicPrefix" json:"publicPrefix"`
}

type ServerConfig struct {
	Name       string         `mapstructure:"name" json:"name"`
	Host       string         `mapstructure:"host" json:"host"`
	Port       int            `mapstructure:"port" json:"port"`
	Tags       []string       `mapstructure:"tags" json:"tags"`
	JWTInfo    JWTConfig      `mapstructure:"jwt" json:"jwt"`
	ConsulInfo ConsulConfig   `mapstructure:"consul" json:"consul"`
	MinIOInfo  MinIOOssConfig `mapstructure:"miniooss" json:"miniooss"`
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
