package discovery

//import (
//	"encoding/json"
//	"fmt"
//	"github.com/pkg/errors"
//	"google.golang.org/grpc/resolver"
//	"strings"
//)
//
///*
//	@Author:Malyue
//	@Description:定义服务实例的结构和属性
//	@CreatedAt:2023/7/6
//*/
//
//type Server struct {
//	Name    string `json:"name"`
//	Addr    string `json:"addr"`
//	Version string `json:"version"`
//	Weight  int64  `json:"weight"`
//}
//
//// BuildPrefix 如果Server版本为空，返回/<服务名称>/的格式，否则返回'/<服务名称>/<版本>/'的格式
//func BuildPrefix(server Server) string {
//	if server.Version == "" {
//		return fmt.Sprintf("/%s/", server.Name)
//	}
//	return fmt.Sprintf("/%s/%s/", server.Name, server.Version)
//}
//
//// BuildRegisterPath 根据给定的Server结构体构建完整的服务注册中心路径，首先调用`BuildPrefix`函数获取前缀路径，然后与`Server`地址拼接起来
//func BuildRegisterPath(server Server) string {
//	return fmt.Sprintf("%s%s", BuildPrefix(server), server.Addr)
//}
//
//// ParseValue 解析从服务注册中心获取的值（即服务的元数据），接受一个字节数组作为输入，并尝试将其解析为`Server`结构体，如果解析成功，将返回解析后的Server对象
//func ParseValue(value []byte) (Server, error) {
//	server := Server{}
//	if err := json.Unmarshal(value, &server); err != nil {
//		return server, err
//	}
//	return server, nil
//}
//
//// SplitPath 拆分服务注册中心路径，并从中提取出Server的信息，接受一个路径字符串作为输入，并将其按照`/`进行划分，然后将从拆分后的路径中提取出地址，将值赋给Server结构体的Addr字段
//func SplitPath(path string) (Server, error) {
//	server := Server{}
//	strs := strings.Split(path, "/")
//	if len(strs) == 0 {
//		return server, errors.New("invalid path")
//	}
//
//	server.Addr = strs[len(strs)-1]
//
//	return server, nil
//}
//
//// Exist helper function
//func Exist(l []resolver.Address, addr resolver.Address) bool {
//	for i := range l {
//		if l[i].Addr == addr.Addr {
//			return true
//		}
//	}
//
//	return false
//}
//
//// Remove helper function
//func Remove(s []resolver.Address, addr resolver.Address) ([]resolver.Address, bool) {
//	for i := range s {
//		if s[i].Addr == addr.Addr {
//			s[i] = s[len(s)-1]
//			return s[:len(s)-1], true
//		}
//	}
//	return nil, false
//}
//
//func BuildResolverUrl(app string) string {
//	return schema + ":///" + app
//}
