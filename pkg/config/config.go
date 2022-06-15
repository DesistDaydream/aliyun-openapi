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
