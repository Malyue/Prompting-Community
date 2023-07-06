package server

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"prompting/internal/gateway/config"
	"strings"
)

type HttpRouter struct {
	config *config.HttpConfig
	router *gin.Engine
}

func NewHttpRouter(config *config.HttpConfig) *HttpRouter {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// 设置鉴权
	//checkFunc,err := NewOuth2()

	// 处理跨域
	cors := CORS()
	route1 := getRoute(config.Routes, "apiserver")
	//route2 := getRoute(config.Routes,"")

	routeMap := map[string]http.HandlerFunc{
		route1.Predicate: apiServerHandleFunc(route1),
	}

	//var apiServerResultHandler = func(c *gin.Context) {
	//	if strings.HasPrefix(c.Request.RequestURI, route1.Predicate) {
	//		writeHandler := responseWriter{
	//			c.Writer,
	//			bytes.NewBuffer([]byte{}),
	//		}
	//		c.Writer = writeHandler
	//		c.Next()
	//
	//	}
	//}
	router.Use(cors, func(c *gin.Context) {
		uri := c.Request.RequestURI
		for _, route := range config.Routes {
			if strings.HasPrefix(uri, route.Predicate) {
				uri = strings.TrimPrefix(uri, route.Predicate)
				c.Request.RequestURI = uri
				c.Request.URL, _ = url.Parse(uri)
				c.Status(200)
				routeMap[route.Predicate].ServeHTTP(c.Writer, c.Request)
				break
			}
		}
		c.Abort()
	})
	return &HttpRouter{config: config, router: router}
}

type responseWriter struct {
	gin.ResponseWriter
	b *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	return w.b.Write(b)
}

func getRoute(routes []*config.HttpRoute, id string) *config.HttpRoute {
	for _, route := range routes {
		if id == route.Id {
			return route
		}
	}
	return nil
}

func (h *HttpRouter) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(rw, r)
}
