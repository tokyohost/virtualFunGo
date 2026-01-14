package main

import (
	"bufio"
	"fmt"
	"github.com/tarm/serial" // 需要安装此库
	"log"
	"time"
)

func main1() {
	// 1. 自动寻找端口 (使用你之前的逻辑)
	s, err := OpenPico() // 使用你之前的 OpenPico
	if err != nil {
		log.Fatal("无法打开串口: ", err)
	}
	defer s.Close()
	// 在 main 的 defer s.Close() 之后加入
	go func() {
		scanner := bufio.NewScanner(s)
		for scanner.Scan() {
			fmt.Printf("Pico 状态反馈: %s\n", scanner.Text())
		}
	}()
	fmt.Println("开始呼吸灯测试... 观察面包板上的 LED")

	for {
		// 按照 hwmon 标准从 0 增加到 255
		for d := 0; d <= 255; d += 5 {
			sendDuty(s, d)
			time.Sleep(30 * time.Millisecond) // 稍微加快一点，呼吸更顺滑
		}
		for d := 255; d >= 0; d -= 5 {
			sendDuty(s, d)
			time.Sleep(30 * time.Millisecond)
		}
	}
}

// 示例：向 Pico 发送设置转速的 JSON 指令
func sendDuty(s *serial.Port, pwmValue int) error {
	// 将 0-255 的 hwmon 值转换为 0-100 的百分比
	percent := int((float64(pwmValue) / 255.0) * 100)
	cmd := fmt.Sprintf("{\"set_duty\": %d}\n", percent)
	_, err := s.Write([]byte(cmd))
	if err != nil {
		log.Printf("发送指令失败: %v", err)
		return err
	}
	return nil
}
