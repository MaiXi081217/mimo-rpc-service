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
// 如果 socketAddr 为空，将使用默认地址 /var/tmp/mimo.sock
func NewBdevService(socketAddr string) *BdevService {
	if socketAddr != "" {
		client.SetSocketAddress(socketAddr)
	}
	return &BdevService{socketAddr: socketAddr}
}

// callRPC 通用 RPC 调用方法，统一处理错误和 JSON 解析
func (s *BdevService) callRPC(method string, params map[string]any) (interface{}, error) {
	result, err := client.Call(method, params)
	if err != nil {
		return nil, err
	}

	var data interface{}
	if err := json.Unmarshal(result, &data); err != nil {
		return nil, fmt.Errorf("unmarshal result failed: %w", err)
	}

	return data, nil
}

// GetBdevs 获取 bdev 列表
// name: bdev 名称，为空则获取所有 bdev
// timeoutMs: 超时时间（毫秒），0 表示不等待
func (s *BdevService) GetBdevs(name string, timeoutMs int) (interface{}, error) {
	params := client.BuildParams(map[string]any{
		"name":       name,
		"timeout_ms": timeoutMs,
	})

	result, err := s.callRPC("bdev_get_bdevs", params)
	if err != nil {
		return nil, fmt.Errorf("get bdev failed: %w", err)
	}
	return result, nil
}

// GetAllBdevs 获取所有 bdev（便利方法）
func (s *BdevService) GetAllBdevs() (interface{}, error) {
	return s.GetBdevs("", 0)
}

// AttachNvmeController 连接 NVMe 控制器
// name: bdev 名称
// trtype: 传输类型，默认为 "PCIe"
// traddr: PCIe 地址（必需）
func (s *BdevService) AttachNvmeController(name, trtype, traddr string) (interface{}, error) {
	if trtype == "" {
		trtype = "PCIe"
	}
	params := client.BuildParams(map[string]any{
		"name":   name,
		"trtype": trtype,
		"traddr": traddr,
	})

	result, err := s.callRPC("bdev_nvme_attach_controller", params)
	if err != nil {
		return nil, fmt.Errorf("attach NVMe controller failed: %w", err)
	}
	return result, nil
}

// AttachNvmeControllerByPCIe 通过 PCIe 地址连接 NVMe 控制器（便利方法）
func (s *BdevService) AttachNvmeControllerByPCIe(name, pcieAddr string) (interface{}, error) {
	return s.AttachNvmeController(name, "PCIe", pcieAddr)
}

// CreateMallocBdev 创建 malloc bdev
// name: bdev 名称（可选）
// uuid: UUID（可选）
// totalSizeMB: 总大小（MB）
// blockSize: 块大小（字节）
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

	result, err := s.callRPC("bdev_malloc_create", params)
	if err != nil {
		return nil, fmt.Errorf("create malloc bdev failed: %w", err)
	}
	return result, nil
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

	result, err := s.callRPC("bdev_raid_create", params)
	if err != nil {
		return nil, fmt.Errorf("create RAID bdev failed: %w", err)
	}
	return result, nil
}

// CreateRaidBdevSimple 简化版创建 RAID bdev（便利方法）
// 使用默认 strip_size_kb=64，superblock=false
func (s *BdevService) CreateRaidBdevSimple(name, raidLevel string, baseBdevs []string) (interface{}, error) {
	req := CreateRaidBdevRequest{
		Name:        name,
		RaidLevel:   raidLevel,
		BaseBdevs:   baseBdevs,
		StripSizeKB: 64,
		Superblock:  false,
	}
	return s.CreateRaidBdev(req)
}

// DetachNvmeController 断开 NVMe 控制器
// name: 控制器名称（必需）
// trtype: 传输类型（可选）
// traddr: 传输地址（可选）
func (s *BdevService) DetachNvmeController(name, trtype, traddr string) (interface{}, error) {
	params := client.BuildParams(map[string]any{
		"name":   name,
		"trtype": trtype,
		"traddr": traddr,
	})

	result, err := s.callRPC("bdev_nvme_detach_controller", params)
	if err != nil {
		return nil, fmt.Errorf("detach NVMe controller failed: %w", err)
	}
	return result, nil
}

// DeleteMallocBdev 删除 malloc bdev
func (s *BdevService) DeleteMallocBdev(name string) (interface{}, error) {
	params := client.BuildParams(map[string]any{
		"name": name,
	})

	result, err := s.callRPC("bdev_malloc_delete", params)
	if err != nil {
		return nil, fmt.Errorf("delete malloc bdev failed: %w", err)
	}
	return result, nil
}

// DeleteRaidBdev 删除 RAID bdev
func (s *BdevService) DeleteRaidBdev(name string) (interface{}, error) {
	params := client.BuildParams(map[string]any{
		"name": name,
	})

	result, err := s.callRPC("bdev_raid_delete", params)
	if err != nil {
		return nil, fmt.Errorf("delete RAID bdev failed: %w", err)
	}
	return result, nil
}

// AddRaidBaseBdev 向 RAID bdev 添加基础 bdev
func (s *BdevService) AddRaidBaseBdev(raidBdev, baseBdev string) (interface{}, error) {
	params := client.BuildParams(map[string]any{
		"raid_bdev": raidBdev,
		"base_bdev": baseBdev,
	})

	result, err := s.callRPC("bdev_raid_add_base_bdev", params)
	if err != nil {
		return nil, fmt.Errorf("add base bdev to RAID failed: %w", err)
	}
	return result, nil
}

// RemoveRaidBaseBdev 从 RAID bdev 移除基础 bdev
func (s *BdevService) RemoveRaidBaseBdev(name string) (interface{}, error) {
	params := client.BuildParams(map[string]any{
		"name": name,
	})

	result, err := s.callRPC("bdev_raid_remove_base_bdev", params)
	if err != nil {
		return nil, fmt.Errorf("remove base bdev from RAID failed: %w", err)
	}
	return result, nil
}

// WipeSuperblock 清除 bdev 的 superblock
// name: bdev 名称（必需）
// size: 清除大小（字节），0 表示使用默认值 1MB
func (s *BdevService) WipeSuperblock(name string, size int) (interface{}, error) {
	params := client.BuildParams(map[string]any{
		"name": name,
		"size": size,
	})

	result, err := s.callRPC("bdev_wipe_superblock", params)
	if err != nil {
		return nil, fmt.Errorf("wipe superblock failed: %w", err)
	}
	return result, nil
}

