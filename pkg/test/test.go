package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// 认证信息
type AuthConfig struct {
	AuthList map[string]Auth `json:"authList" yaml:"authList"`
}
type Auth struct {
	AccessKeyID     string `json:"accessKeyId" yaml:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret" yaml:"accessKeySecret"`
}

func main() {
	var auth *AuthConfig

	fileByte, err := os.ReadFile("auth.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(fileByte, &auth)
	if err != nil {
		panic(err)
	}
	fmt.Println(auth.AuthList["desistdaydream.ltd"])
}
