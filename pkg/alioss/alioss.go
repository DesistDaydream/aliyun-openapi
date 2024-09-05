package alioss

type AlidnsHandler struct {
}

// 实例化云解析处理器
func NewAliossHandler(ak, sk, domainName, region string) *AlidnsHandler {
	return &AlidnsHandler{}
}
