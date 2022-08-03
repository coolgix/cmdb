package utils

//hash函数工具

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"

	"github.com/infraboard/mcube/logger/zap"
)

//应为是对string进行hash
//传递任何对象进来，使用sha1进行hash
func Hash(x any) string {
	hash := sha1.New()
	//通过json的marshal把byte改为string
	b, err := json.Marshal(x)
	if err != nil {
		zap.L().Errorf("hash %v error, %s", x, err)
		return ""
	}
	hash.Write(b)
	//格式化16进制的数据
	return fmt.Sprintf("%x", hash.Sum(nil))
}
