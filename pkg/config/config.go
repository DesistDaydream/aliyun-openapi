package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// 认证信息配置
type AuthConfig struct {
	AuthList map[string]Auth `json:"authList" yaml:"authList"`
}
type Auth struct {
	AccessKeyID     string `json:"accessKeyId" yaml:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret" yaml:"accessKeySecret"`
}

func NewAuthInfo(file string) (auth *AuthConfig) {
	// 读取认证信息
	fileByte, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(fileByte, &auth)
	if err != nil {
		panic(err)
	}
	return auth
}

// 判断文件中是否存在域名
func (c *AuthConfig) IsDomainExist(domName string) bool {
	if _, ok := c.AuthList[domName]; ok {
		return true
	} else {
		return false
	}
}
