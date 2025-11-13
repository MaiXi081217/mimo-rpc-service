package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mimo/mimo-rpc-service/client"
)

// BdevService bdev 相关服务
type BdevService struct {
	socketAddr string
}

// NewBdevService 创建 bdev 服务实例
func NewBdevService(socketAddr string) *BdevService {
	if socketAddr != "" {
		client.SetSocketAddress(socketAddr)
	}
	return &BdevService{socketAddr: socketAddr}
}

// GetBdevs 获取 bdev 列表
func (s *BdevService) GetBdevs(name string, timeoutMs int) (interface{}, error) {
	params := client.BuildParams(map[string]any{
		"name":       name,
		"timeout_ms": timeoutMs,
	})

	result, err := client.Call("bdev_get_bdevs", params)
	if err != nil {
		return nil, fmt.Errorf("get bdev failed: %w", err)
	}

	var data interface{}
	if err := json.Unmarshal(result, &data); err != nil {
		return nil, fmt.Errorf("unmarshal result failed: %w", err)
	}

	return data, nil
}

// AttachNvmeController 连接 NVMe 控制器
func (s *BdevService) AttachNvmeController(name, trtype, traddr string) (interface{}, error) {
	params := client.BuildParams(map[string]any{
		"name":   name,
		"trtype": trtype,
		"traddr": traddr,
	})

	result, err := client.Call("bdev_nvme_attach_controller", params)
	if err != nil {
		return nil, fmt.Errorf("attach NVMe controller failed: %w", err)
	}

	var data interface{}
	if err := json.Unmarshal(result, &data); err != nil {
		return nil, fmt.Errorf("unmarshal result failed: %w", err)
	}

	return data, nil
}

// CreateMallocBdev 创建 malloc bdev
func (s *BdevService) CreateMallocBdev(name, uuid string, totalSizeMB float64, blockSize int) (interface{}, error) {
	if totalSizeMB <= 0 || blockSize <= 0 {
		return nil, fmt.Errorf("total_size and block_size must be positive")
	}

	numBlocks := int((totalSizeMB * 1024 * 1024) / float64(blockSize))

	params := client.BuildParams(map[string]any{
		"name":       name,
		"uuid":       uuid,
		"block_size": blockSize,
		"num_blocks": numBlocks,
	})

	result, err := client.Call("bdev_malloc_create", params)
	if err != nil {
		return nil, fmt.Errorf("create malloc bdev failed: %w", err)
	}

	var data interface{}
	if err := json.Unmarshal(result, &data); err != nil {
		return nil, fmt.Errorf("unmarshal result failed: %w", err)
	}

	return data, nil
}

// CreateRaidBdevRequest 创建 RAID bdev 请求
type CreateRaidBdevRequest struct {
	Name        string   `json:"name"`
	RaidLevel   string   `json:"raid_level"`
	BaseBdevs   []string `json:"base_bdevs"`
	StripSizeKB int      `json:"strip_size_kb,omitempty"`
	UUID        string   `json:"uuid,omitempty"`
	Superblock  bool     `json:"superblock,omitempty"`
}

// CreateRaidBdev 创建 RAID bdev
func (s *BdevService) CreateRaidBdev(req CreateRaidBdevRequest) (interface{}, error) {
	if req.Name == "" || req.RaidLevel == "" || len(req.BaseBdevs) == 0 {
		return nil, fmt.Errorf("name, raid_level, and base_bdevs are required")
	}

	// 如果 BaseBdevs 只有一个元素且包含空格，尝试分割
	baseList := req.BaseBdevs
	if len(baseList) == 1 && strings.Contains(baseList[0], " ") {
		baseList = strings.Fields(baseList[0])
	}

	params := client.BuildParams(map[string]any{
		"name":          req.Name,
		"raid_level":    req.RaidLevel,
		"base_bdevs":    baseList,
		"strip_size_kb": req.StripSizeKB,
		"uuid":          req.UUID,
		"superblock":    req.Superblock,
	})

	result, err := client.Call("bdev_raid_create", params)
	if err != nil {
		return nil, fmt.Errorf("create RAID bdev failed: %w", err)
	}

	var data interface{}
	if err := json.Unmarshal(result, &data); err != nil {
		return nil, fmt.Errorf("unmarshal result failed: %w", err)
	}

	return data, nil
}

