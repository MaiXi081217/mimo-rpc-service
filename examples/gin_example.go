package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mimo/mimo-rpc-service/service"
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

	// 连接 NVMe 控制器
	r.POST("/api/nvme/attach", func(c *gin.Context) {
		var req struct {
			Name   string `json:"name" binding:"required"`
			Trtype string `json:"trtype"`
			Traddr string `json:"traddr" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		trtype := req.Trtype
		if trtype == "" {
			trtype = "PCIe"
		}

		result, err := bdevSvc.AttachNvmeController(req.Name, trtype, req.Traddr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	})

	r.Run(":8080")
}

