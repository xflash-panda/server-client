package main

import (
	"log"
	"time"

	"github.com/xflash-panda/server-client/pkg"
)

func main() {
	// 创建客户端
	apiConfig := &pkg.Config{
		APIHost: "http://127.0.0.1:8080",
		Token:   "your-token-here",
		Timeout: 5 * time.Second,
		Debug:   true,
	}
	client := pkg.New(apiConfig)

	// 1. 获取节点配置 (新接口 - 只需要 nodeID)
	config, err := client.Config(1, pkg.Hysteria2)
	if err != nil {
		log.Fatalf("获取配置失败: %v", err)
	}
	log.Printf("节点配置: %v", config)

	// 2. 注册节点 (修改后的接口 - 只返回 register_id)
	registerId, err := client.Register(1, pkg.Hysteria2, "test-hostname", 8080, "127.0.0.1")
	if err != nil {
		log.Fatalf("注册失败: %v", err)
	}
	log.Printf("注册成功，Register ID: %d", registerId)

	// 3. 使用 registerId 进行后续操作
	users, err := client.Users(registerId, pkg.Hysteria2)
	if err != nil {
		log.Fatalf("获取用户列表失败: %v", err)
	}
	log.Printf("用户数量: %d", len(*users))
}
