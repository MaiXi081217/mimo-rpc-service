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

	// 示例 1: 获取所有 bdev（使用便利方法）
	fmt.Println("=== 获取所有 bdev ===")
	result, err := bdevSvc.GetAllBdevs()
	if err != nil {
		log.Fatalf("获取 bdev 失败: %v", err)
	}
	data, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(data))

	// 示例 2: 创建 RAID bdev（使用简化方法）
	fmt.Println("\n=== 创建 RAID bdev（简化方法）===")
	result, err = bdevSvc.CreateRaidBdevSimple("raid1", "raid10", []string{"bdev1", "bdev2", "bdev3", "bdev4"})
	if err != nil {
		log.Printf("创建 RAID bdev 失败: %v", err)
	} else {
		data, _ = json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
	}

	// 示例 3: 创建 RAID bdev（完整方法，带所有参数）
	fmt.Println("\n=== 创建 RAID bdev（完整方法）===")
	req := service.CreateRaidBdevRequest{
		Name:        "raid2",
		RaidLevel:   "raid1",
		BaseBdevs:   []string{"bdev5", "bdev6"},
		StripSizeKB: 128,
		Superblock:  true,
	}
	result, err = bdevSvc.CreateRaidBdev(req)
	if err != nil {
		log.Printf("创建 RAID bdev 失败: %v", err)
	} else {
		data, _ = json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
	}

	// 示例 4: 连接 NVMe 控制器（使用便利方法）
	fmt.Println("\n=== 连接 NVMe 控制器（便利方法）===")
	result, err = bdevSvc.AttachNvmeControllerByPCIe("Nvme0", "0000:01:00.0")
	if err != nil {
		log.Printf("连接 NVMe 控制器失败: %v", err)
	} else {
		data, _ = json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
	}
}

