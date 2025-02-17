package rpc

import (
	"danmu-http/logger"
	api "danmu-http/rpc/proto"
	"danmu-http/setting"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

var (
	client     api.LiveServiceClient
	conn       *grpc.ClientConn
	clientOnce sync.Once
	mu         sync.RWMutex
)

const (
	dialTimeout      = 5 * time.Second
	keepaliveTime    = 30 * time.Second
	keepaliveTimeout = 10 * time.Second
)

// Init 初始化 RPC 客户端
func Init() error {
	var err error
	clientOnce.Do(func() {
		err = connect()
	})
	return err
}

// connect 建立 gRPC 连接
func connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()

	var err error
	conn, err = grpc.DialContext(ctx, setting.RPCSetting.LiveServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                keepaliveTime,
			Timeout:             keepaliveTimeout,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		logger.Error().Err(err).Msg("failed to connect to gRPC server")
		return err
	}

	client = api.NewLiveServiceClient(conn)
	return nil
}

// GetClient 获取 RPC 客户端，如果连接断开会尝试重连
func GetClient() (api.LiveServiceClient, error) {
	mu.RLock()
	if client != nil && conn != nil && conn.GetState().String() == "READY" {
		defer mu.RUnlock()
		return client, nil
	}
	mu.RUnlock()

	mu.Lock()
	defer mu.Unlock()

	// 双重检查
	if client != nil && conn != nil && conn.GetState().String() == "READY" {
		return client, nil
	}

	// 重新连接
	if err := connect(); err != nil {
		return nil, err
	}

	return client, nil
}

// Close 关闭 RPC 连接
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if conn != nil {
		if err := conn.Close(); err != nil {
			logger.Error().Err(err).Msg("failed to close gRPC connection")
			return err
		}
		conn = nil
		client = nil
	}
	return nil
}

// WithTimeout 使用超时上下文包装 RPC 调用
func WithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
