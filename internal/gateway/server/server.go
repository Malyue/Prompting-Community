package server

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"prompting/internal/gateway/config"
	"sync"
)

type Server struct {
	config     *config.Endpoint
	httpServer *http.Server
	grpcServer *grpc.Server
	mutex      sync.Mutex
	running    bool
}

func NewServer(config *config.Endpoint) *Server {
	httpServer := &http.Server{Handler: NewHttpRouter(config.Http)}
	grpcServer := grpc.NewServer(grpc.MaxRecvMsgSize(10 * 1024 * 1024))
	//RegisterServiceRoute(grpcServer)
	return &Server{
		config:     config,
		httpServer: httpServer,
		grpcServer: grpcServer,
	}
}

func (s *Server) Run() error {
	if !s.running {
		s.mutex.Lock()
		defer s.mutex.Unlock()
		zap.L().Info("http listen the port:" + s.config.Http.Port)
		httpListen, err := net.Listen("tcp", s.config.Http.Port)
		if err != nil {
			return errors.Wrap(err, "http listen fail")
		}
		zap.L().Info("grpc listen the port: " + s.config.Grpc.Addr)
		rpcListen, err := net.Listen("tcp", s.config.Grpc.Addr)
		go func() {
			err := s.httpServer.Serve(httpListen)
			if err != net.ErrClosed {
				return
			}
		}()

		go func() {
			err := s.grpcServer.Serve(rpcListen)
			if err != net.ErrClosed {
				return
			}
		}()
	}

	return nil
}

func (s *Server) Close() error {
	if s.running {
		wg := sync.WaitGroup{}
		wg.Add(2)
		errCh := make(chan error, 1)
		go func() {
			defer wg.Done()
			errCh <- s.httpServer.Close()
		}()

		go func() {
			defer wg.Done()
			s.grpcServer.Stop()
		}()
		wg.Wait()
		err := <-errCh
		return err
	}
	return nil
}
