package conf

import (
	"k8s.io/apimachinery/pkg/util/yaml"
	"os"
	"sync"
)

var once sync.Once

type Conf struct {
	AirCenter struct {
		Host string
	}
	Sim800C struct {
		Port string
		Baud int
	}
}

var conf *Conf

// Get 读取配置
// 默认从配置文件取，如果配置文件中的db节点内容为空，则从环境变量取
// 如果配置文件不存在，则db从环境变量取，其他值使用默认值
func Get() *Conf {
	once.Do(func() {
		if conf == nil {
			filePath := getConfFilePath()
			conf = getConfiguration(filePath)
		}
	})
	return conf
}

// getConfiguration 读取配置
// 优先从配置文件读取，如果数据库相关配置为空，则从环境变量读取
func getConfiguration(filePath *string) *Conf {
	if file, err := os.ReadFile(*filePath); err != nil {
		return getConfFromEnv()
	} else {
		c := Conf{}
		err := yaml.Unmarshal(file, &c)
		if err != nil {
			return getConfFromEnv()
		}
		return &c
	}
}

func getConfFromEnv() *Conf {
	return &Conf{AirCenter: struct{ Host string }{Host: os.Getenv("aircenter__host")}}
}
