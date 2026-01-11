package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/tarm/serial" // 需要安装此库
)

var currentPicoRPM int // 全局变量存储 Pico 传回的转速

// 协程：监听 Pico 串口
func ReadPicoSerial(portName string) {
	c := &serial.Config{Name: portName, Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Printf("无法打开串口 %s: %v", portName, err)
		return
	}

	reader := bufio.NewReader(s)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("读取串口错误: %v", err)
			break
		}
		// 假设 Pico 发送格式为 "RPM:1234"
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "RPM:") {
			fmt.Sscanf(line, "RPM:%d", &currentPicoRPM)
		}
	}
}

func FindPicoPort() string {
	// 方法 A: 扫描 by-id 目录
	idPath := "/dev/serial/by-id/"
	files, err := ioutil.ReadDir(idPath)
	if err == nil {
		for _, f := range files {
			if strings.Contains(f.Name(), "Pico") || strings.Contains(f.Name(), "Raspberry_Pi") {
				print("找到串口" + f.Name())
				return idPath + f.Name()
			}
		}
	}

	// 方法 B: 如果 by-id 不存在，尝试默认值
	return ""
}
