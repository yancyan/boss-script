package env

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"os"
)

// config keys.
const (
	KEY_TEST = "cookie"
	APP_NAME = "application.name"
)

var Config config = config{}

type config struct {
	Domain       string `json:"domain"`
	Cookie       string `json:"cookie"`
	PartnerChurn partnerChurn
}

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	log.Infof("======  >>> viper config is %+v", viper.AllSettings())

	//	将viper配置映射到struct
	Config.InitCustomConfig()

	log.Infof("======  >>> env.Config is %+v", Config)
}

func (c *config) InitCustomConfig() {
	//f, err := ioutil.ReadFile("config/application.yaml")
	//if err != nil {
	//	panic(err)
	//}
	out, err := yaml.Marshal(viper.AllSettings())
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(out, c)
	if err != nil {
		panic(err)
	}
	c.PartnerChurn = initPartnerChurn()
}
