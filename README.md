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
    result, err := bdevSvc.GetBdevs("", 0)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    data, _ := json.MarshalIndent(result, "", "  ")
    fmt.Println(string(data))
    
    // 创建 RAID bdev
    req := service.CreateRaidBdevRequest{
        Name:        "raid1",
        RaidLevel:   "raid10",
        BaseBdevs:   []string{"bdev1", "bdev2", "bdev3", "bdev4"},
        StripSizeKB: 64,
    }
    result, err = bdevSvc.CreateRaidBdev(req)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    // 连接 NVMe 控制器
    result, err = bdevSvc.AttachNvmeController("Nvme0", "PCIe", "0000:01:00.0")
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
- `GetBdevs(name, timeoutMs)` - 获取 bdev 列表（name 和 timeoutMs 为可选参数）

**NVMe 控制器：**
- `AttachNvmeController(name, trtype, traddr)` - 连接 NVMe 控制器（所有参数必需）
- `DetachNvmeController(name, trtype, traddr)` - 断开 NVMe 控制器（name 必需，trtype 和 traddr 可选）

**Malloc Bdev：**
- `CreateMallocBdev(name, uuid, totalSizeMB, blockSize)` - 创建 malloc bdev（totalSizeMB 和 blockSize 必需）
- `DeleteMallocBdev(name)` - 删除 malloc bdev（name 必需）

**RAID Bdev：**
- `CreateRaidBdev(req)` - 创建 RAID bdev（name, raid_level, base_bdevs 必需）
- `DeleteRaidBdev(name)` - 删除 RAID bdev（name 必需）
- `AddRaidBaseBdev(raidBdev, baseBdev)` - 向 RAID 添加基础 bdev（所有参数必需）
- `RemoveRaidBaseBdev(name)` - 从 RAID 移除基础 bdev（name 必需）

**其他：**
- `WipeSuperblock(name, size)` - 清除 bdev 的 superblock（name 必需，size 可选）

更多示例请查看 `examples/` 目录。
