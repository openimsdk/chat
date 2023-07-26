package config

var imConfig struct {
	Mysql struct {
		Address       []string `yaml:"address"`
		Username      string   `yaml:"username"`
		Password      string   `yaml:"password"`
		Database      string   `yaml:"database"`
		MaxOpenConn   int      `yaml:"maxOpenConn"`
		MaxIdleConn   int      `yaml:"maxIdleConn"`
		MaxLifeTime   int      `yaml:"maxLifeTime"`
		LogLevel      int      `yaml:"logLevel"`
		SlowThreshold int      `yaml:"slowThreshold"`
	} `yaml:"mysql"`

	Mongo struct {
		Uri         string   `yaml:"uri"`
		Address     []string `yaml:"address"`
		Database    string   `yaml:"database"`
		Username    string   `yaml:"username"`
		Password    string   `yaml:"password"`
		MaxPoolSize int      `yaml:"maxPoolSize"`
	} `yaml:"mongo"`

	Redis struct {
		Address  []string `yaml:"address"`
		Username string   `yaml:"username"`
		Password string   `yaml:"password"`
	} `yaml:"redis"`

	Rpc struct {
		RegisterIP string `yaml:"registerIP"`
		ListenIP   string `yaml:"listenIP"`
	} `yaml:"rpc"`

	Log struct {
		StorageLocation     string `yaml:"storageLocation"`
		RotationTime        uint   `yaml:"rotationTime"`
		RemainRotationCount uint   `yaml:"remainRotationCount"`
		RemainLogLevel      int    `yaml:"remainLogLevel"`
		IsStdout            bool   `yaml:"isStdout"`
		IsJson              bool   `yaml:"isJson"`
		WithStack           bool   `yaml:"withStack"`
	} `yaml:"log"`

	Manager struct {
		UserID   []string `yaml:"userID"`
		Nickname []string `yaml:"nickname"`
	} `yaml:"manager"`

	Secret      string `yaml:"secret"`
	TokenPolicy struct {
		Expire int64 `yaml:"expire"`
	} `yaml:"tokenPolicy"`
}
