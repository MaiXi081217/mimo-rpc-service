package client

import (
	"encoding/json"
	"fmt"
	"sync"

	spdk "github.com/spdk/spdk/go/rpc/client"
)

const DefaultSocketAddress = "/var/tmp/mimo.sock"

var (
	clientInstance *spdk.Client
	socketAddr     string
	once           sync.Once
)

// SetSocketAddress 设置 socket 地址
func SetSocketAddress(addr string) {
	socketAddr = addr
}

// GetClient 获取单例 RPC 客户端
func GetClient() (*spdk.Client, error) {
	var err error
	once.Do(func() {
		addr := socketAddr
		if addr == "" {
			addr = DefaultSocketAddress
		}
		c, e := spdk.CreateClientWithJsonCodec(spdk.Unix, addr)
		if e != nil {
			err = fmt.Errorf("failed to connect RPC (%s): %w", addr, e)
			return
		}
		clientInstance = c
	})
	if err != nil {
		return nil, err
	}
	return clientInstance, nil
}

// Call 通用 RPC 调用方法
func Call(method string, params map[string]any) ([]byte, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}
	resp, err := client.Call(method, params)
	if err != nil {
		return nil, fmt.Errorf("RPC call failed (%s): %w", method, err)
	}

	data, err := json.MarshalIndent(resp.Result, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal failed: %w", err)
	}
	return data, nil
}

// BuildParams 构建 RPC 参数（自动过滤空值）
func BuildParams(args map[string]any) map[string]any {
	params := make(map[string]any)
	for k, v := range args {
		if v == nil {
			continue
		}
		switch val := v.(type) {
		case string:
			if val != "" {
				params[k] = val
			}
		default:
			params[k] = val
		}
	}
	return params
}

// Close 关闭 RPC 客户端连接
func Close() {
	if clientInstance != nil {
		clientInstance.Close()
	}
}

