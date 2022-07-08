package aliclient

type ClientInfo struct {
	AK     string
	SK     string
	Region string
}

// 用于在根命令的 PersistentPreRun 中初始化账号Client
// 以便在子命令中直接引用该变量的值
var Info *ClientInfo
