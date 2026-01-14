package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/tarm/serial"      // 需要安装此库
	"go.bug.st/serial/enumerator" // 扫描端口功能
)

// 定义与 Pico 输出匹配的结构体
type PicoData struct {
	RPM  int `json:"rpm"`
	Duty int `json:"duty"`
}
type FanStatus struct {
	ID   string `json:"id"`
	RPM  int    `json:"rpm"`
	Duty int    `json:"duty"`
}

func OpenPico() (*serial.Port, error) {
	portName := FindPicoPortV2()
	if portName == "" {
		return nil, fmt.Errorf("未发现 Pico 设备")
	}

	// 配置串口参数
	config := &serial.Config{
		Name: portName,
		Baud: 115200,
		// tarm/serial 默认是 8N1 模式，这与 MicroPython 匹配
	}

	// 真正打开串口，返回 *serial.Port
	s, err := serial.OpenPort(config)
	if err != nil {
		return nil, fmt.Errorf("打开串口失败: %v", err)
	}

	return s, nil
}

func FindPicoPortV2() string {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return ""
	}
	for _, port := range ports {
		if port.IsUSB {
			// Raspberry Pi Pico 的 VID 是 2E8A
			// 某些 MicroPython 固件可能显示为大写或小写，所以用 EqualFold
			if strings.EqualFold(port.VID, "2E8A") {
				log.Printf("检测到 Pico 设备: %s (PID: %s)", port.Name, port.PID)
				return port.Name
			}
		}
	}
	return ""
}
