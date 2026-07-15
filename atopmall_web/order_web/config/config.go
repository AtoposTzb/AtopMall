package config

type GoodsOrderSrvConfig struct {
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

type AlipayConfig struct {
	AppID        string `mapstructure:"app_id" json:"app_id"`
	PrivateKey   string `mapstructure:"private_key" json:"private_key"`
	AliPublicKey string `mapstructure:"ali_public_key" json:"ali_public_key"`
	NotifyURL    string `mapstructure:"notify_url" json:"notify_url"`
	ReturnURL    string `mapstructure:"return_url" json:"return_url"`
	IsProduction bool   `mapstructure:"is_production" json:"is_production"`
	ProductCode  string `mapstructure:"product_code" json:"product_code"`
}

type ServerConfig struct {
	Name             string              `mapstructure:"name" json:"name"`
	Host             string              `mapstructure:"host" json:"host"`
	Port             int                 `mapstructure:"port" json:"port"`
	Tags             []string            `mapstructure:"tags" json:"tags"`
	GoodsSrvInfo     GoodsOrderSrvConfig `mapstructure:"goods_srv" json:"goods_srv"`
	OrderSrvInfo     GoodsOrderSrvConfig `mapstructure:"order_srv" json:"order_srv"`
	InventorySrvInfo GoodsOrderSrvConfig `mapstructure:"inventory_srv" json:"inventory_srv"`
	JWTInfo          JWTConfig           `mapstructure:"jwt" json:"jwt"`
	ConsulInfo       ConsulConfig        `mapstructure:"consul" json:"consul"`
	AlipayInfo       AlipayConfig        `mapstructure:"alipay" json:"alipay"`
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
