package file

import (
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

// NewConfig 读取指定的配置文件来获取
func NewConfig[a any](str string) (*a, error) {
	file, e := os.Open(str)
	if e != nil {
		return nil, e
	}
	defer file.Close()

	yamlFile, e := io.ReadAll(file)
	if e != nil {
		return nil, e
	}

	var a1 a

	config := &a1
	e = yaml.Unmarshal(yamlFile, config)
	if e != nil {
		return nil, e
	}

	return config, nil
}
