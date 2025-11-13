package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mimo/mimo-rpc-service/service"
)

func main() {
	// 创建服务实例（使用默认 socket 地址）
	bdevSvc := service.NewBdevService("")

	// 示例 1: 获取所有 bdev
	fmt.Println("=== 获取所有 bdev ===")
	result, err := bdevSvc.GetBdevs("", 0)
	if err != nil {
		log.Fatalf("获取 bdev 失败: %v", err)
	}
	data, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(data))

	// 示例 2: 创建 RAID bdev
	fmt.Println("\n=== 创建 RAID bdev ===")
	req := service.CreateRaidBdevRequest{
		Name:        "raid1",
		RaidLevel:   "raid10",
		BaseBdevs:   []string{"bdev1", "bdev2", "bdev3", "bdev4"},
		StripSizeKB: 64,
		Superblock:  false,
	}
	result, err = bdevSvc.CreateRaidBdev(req)
	if err != nil {
		log.Printf("创建 RAID bdev 失败: %v", err)
	} else {
		data, _ = json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
	}

	// 示例 3: 连接 NVMe 控制器
	fmt.Println("\n=== 连接 NVMe 控制器 ===")
	result, err = bdevSvc.AttachNvmeController("Nvme0", "PCIe", "0000:01:00.0")
	if err != nil {
		log.Printf("连接 NVMe 控制器失败: %v", err)
	} else {
		data, _ = json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
	}
}

