# MIMO RPC Service

MIMO 存储系统的 RPC 服务层 SDK，提供简洁的 Go API 用于与 SPDK RPC 服务交互。

## 安装

```bash
go get github.com/mimo/mimo-rpc-service
```

## 快速开始

```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/mimo/mimo-rpc-service/service"
)

func main() {
    // 创建服务实例（使用默认 socket: /var/tmp/mimo.sock）
    bdevSvc := service.NewBdevService("")
    
    // 获取所有 bdev
    result, err := bdevSvc.GetAllBdevs()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    data, _ := json.MarshalIndent(result, "", "  ")
    fmt.Println(string(data))
    
    // 创建 RAID bdev（简化方法）
    result, err = bdevSvc.CreateRaidBdevSimple("raid1", "raid10", 
        []string{"bdev1", "bdev2", "bdev3", "bdev4"})
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    // 连接 NVMe 控制器（便利方法）
    result, err = bdevSvc.AttachNvmeControllerByPCIe("Nvme0", "0000:01:00.0")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
}
```

## 主要 API

### 创建服务

```go
bdevSvc := service.NewBdevService("")  // 使用默认 socket
bdevSvc := service.NewBdevService("/custom/path/sock")  // 自定义 socket
```

### 常用方法

**查询和管理：**
- `GetAllBdevs()` - 获取所有 bdev
- `GetBdevs(name, timeoutMs)` - 获取指定 bdev

**NVMe 控制器：**
- `AttachNvmeControllerByPCIe(name, pcieAddr)` - 连接 NVMe 控制器（便利方法）
- `AttachNvmeController(name, trtype, traddr)` - 连接 NVMe 控制器（完整方法）
- `DetachNvmeController(name, trtype, traddr)` - 断开 NVMe 控制器

**Malloc Bdev：**
- `CreateMallocBdev(name, uuid, totalSizeMB, blockSize)` - 创建 malloc bdev
- `DeleteMallocBdev(name)` - 删除 malloc bdev

**RAID Bdev：**
- `CreateRaidBdevSimple(name, raidLevel, baseBdevs)` - 创建 RAID bdev（简化版，默认 strip_size_kb=64）
- `CreateRaidBdev(req)` - 创建 RAID bdev（完整版）
- `DeleteRaidBdev(name)` - 删除 RAID bdev
- `AddRaidBaseBdev(raidBdev, baseBdev)` - 向 RAID 添加基础 bdev
- `RemoveRaidBaseBdev(name)` - 从 RAID 移除基础 bdev

**其他：**
- `WipeSuperblock(name, size)` - 清除 bdev 的 superblock（默认 1MB）

更多示例请查看 `examples/` 目录。
