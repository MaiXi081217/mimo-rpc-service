# MIMO RPC Service

MIMO 存储系统的 RPC 服务层 SDK，提供简洁的 Go API 用于与 SPDK RPC 服务交互。

## 特性

-  **简洁易用**: 提供清晰的 API 接口
-  **类型安全**: 使用 Go 类型系统，避免参数错误
-  **自动参数处理**: 自动过滤空值和无效参数
-  **灵活配置**: 支持自定义 socket 地址
-  **独立项目**: 不依赖 CLI 或其他系统组件

## 安装

```bash
go get github.com/mimo/mimo-rpc-service
```

## 快速开始

### 基本使用

```go
package main

import (
    "fmt"
    "github.com/mimo/mimo-rpc-service/service"
)

func main() {
    // 创建服务实例（使用默认 socket: /var/tmp/mimo.sock）
    bdevSvc := service.NewBdevService("")
    
    // 或者使用自定义 socket
    // bdevSvc := service.NewBdevService("/custom/path/sock")
    
    // 获取所有 bdev
    result, err := bdevSvc.GetBdevs("", 0)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Bdevs: %+v\n", result)
}
```

### Gin Web 服务示例

```go
package main

import (
    "net/http"
    "github.com/mimo/mimo-rpc-service/service"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // 创建服务实例
    bdevSvc := service.NewBdevService("/var/tmp/mimo.sock")

    // 获取 bdev 列表
    r.GET("/api/bdevs", func(c *gin.Context) {
        name := c.Query("name")
        result, err := bdevSvc.GetBdevs(name, 0)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, result)
    })

    // 创建 RAID bdev
    r.POST("/api/raid/create", func(c *gin.Context) {
        var req service.CreateRaidBdevRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        result, err := bdevSvc.CreateRaidBdev(req)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, result)
    })

    r.Run(":8080")
}
```

## API 文档

### BdevService

#### NewBdevService(socketAddr string) *BdevService

创建 bdev 服务实例。

- `socketAddr`: RPC socket 地址，为空则使用默认地址 `/var/tmp/mimo.sock`

#### GetBdevs(name string, timeoutMs int) (interface{}, error)

获取 bdev 列表。

- `name`: bdev 名称（可选，为空则返回所有）
- `timeoutMs`: 超时时间（毫秒，0 表示不等待）

返回: bdev 列表数据

#### AttachNvmeController(name, trtype, traddr string) (interface{}, error)

连接 NVMe 控制器。

- `name`: bdev 名称
- `trtype`: 传输类型（如 "PCIe"）
- `traddr`: 传输地址（如 PCIe 地址）

#### CreateMallocBdev(name, uuid string, totalSizeMB float64, blockSize int) (interface{}, error)

创建 malloc bdev。

- `name`: bdev 名称（可选）
- `uuid`: UUID（可选）
- `totalSizeMB`: 总大小（MB）
- `blockSize`: 块大小（字节）

#### CreateRaidBdev(req CreateRaidBdevRequest) (interface{}, error)

创建 RAID bdev。

### CreateRaidBdevRequest

```go
type CreateRaidBdevRequest struct {
    Name        string   `json:"name"`                  // RAID bdev 名称（必需）
    RaidLevel   string   `json:"raid_level"`            // RAID 级别：raid0, raid1, raid10, concat（必需）
    BaseBdevs   []string `json:"base_bdevs"`            // 基础 bdev 列表（必需）
    StripSizeKB int      `json:"strip_size_kb,omitempty"` // 条带大小（KB，可选）
    UUID        string   `json:"uuid,omitempty"`         // UUID（可选）
    Superblock  bool     `json:"superblock,omitempty"`   // 是否启用 superblock（可选）
}
```

## 项目结构

```
mimo-rpc-service/
├── client/          # RPC 客户端封装
│   └── client.go
├── service/         # 服务层 API
│   └── bdev.go
├── go.mod
└── README.md
```

## 依赖

- `github.com/spdk/spdk/go/rpc` - SPDK Go RPC 客户端

## 注意事项

 **此包仅包含 RPC 相关功能**，不包含系统初始化、系统更新等功能。

系统管理功能请使用 MIMO CLI 工具。


