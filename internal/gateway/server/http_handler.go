package server

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
	"net/http"
	"prompting/internal/gateway/config"
)

func apiServerHandleFunc(route *config.HttpRoute) http.HandlerFunc {
	// 创建一个对象，添加一个自定义的元数据处理函数，从HTTP请求的头部中提取相关信息，并将其作为元数据附加到gRPC请求中
	mux := runtime.NewServeMux(runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
		m := make(map[string]string)
		//m["currentUserId"] = request.Header.Get("currentUserId")
		//m["tenantId"] = request.Header.Get("tenantId")
		return metadata.New(m)
	}))

	//opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	// 将grpc注册到http server中，然后就可以转发了
	//err := boss.RegisterBoosServiceHandlerFromEndpoint(context.Background(), mux, route.Endpoint, opts) //注册boss服务

	return mux.ServeHTTP

}
