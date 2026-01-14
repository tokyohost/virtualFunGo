package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/tarm/serial"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const homePath = "/sys/class/hwmon/"
const isSerialRunning = false

func main() {
	fmt.Println("=== Pico 虚拟风扇已启动 ===")

	for {
		// 1. 尝试初始化：查找串口和 HWMON 路径
		s, hwmonPath, err := initializeHardware()
		if err != nil {
			log.Printf("等待硬件就绪: %v", err)
			time.Sleep(3 * time.Second) // 探测频率
			continue
		}

		fmt.Printf("成功连接！串口已打开，驱动路径: %s\n", hwmonPath)

		// 2. 创建上下文，用于管理这一轮连接的协程生命周期
		ctx, cancel := context.WithCancel(context.Background())

		// 3. 启动桥接业务逻辑
		// 我们传一个 done channel，用来知道业务协程什么时候因为错误退出了
		done := make(chan struct{})
		go func() {
			startBridge(ctx, s, hwmonPath)
			close(done)
		}()

		// 4. 等待信号：要么是业务报错退出，要么是手动关掉
		<-done
		fmt.Println("硬件连接断开，尝试重新恢复...")

		// 清理资源
		cancel()
		s.Close()
		time.Sleep(2 * time.Second)
	}
}

// 初始化硬件：同时找到串口和驱动路径才算成功
func initializeHardware() (*serial.Port, string, error) {
	// 找到串口路径 (by-id)
	s, err := OpenPico() // 使用你之前的 OpenPico
	if err != nil {
		return nil, "", err
	}

	// 找到 C 驱动路径
	if runtime.GOOS == "windows" {
		return s, "", nil
	}
	path, err := findHwmonPath() // 使用你之前的 findHwmonPath
	if err != nil {
		s.Close()
		return nil, "", err
	}

	return s, path, nil
}

// fileExists 检查文件是否存在
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
func startBridge(ctx context.Context, s *serial.Port, hwmonPath string) {
	// 协程 1: 监听驱动 PWM (1, 2, 3) -> 发给 Pico
	go func() {
		// 使用 map 记录每个风扇的上一次 PWM 值
		lastPwms := make(map[string]int)
		// 这里的 ID 对应你 Pico 字典里的 Key
		fanIds := []string{"fan1", "fan2", "fan3"}

		for {
			select {
			case <-ctx.Done():
				return
			default:
				for i, id := range fanIds {
					// 驱动里的文件名通常从 1 开始: pwm1, pwm2...
					fileName := fmt.Sprintf("pwm%d", i+1)
					pwmFile := filepath.Join(hwmonPath, fileName)

					val := readIntFromFile(pwmFile)
					fmt.Println("读取驱动文件: ", pwmFile, " 得到 PWM: ", val)
					if val != lastPwms[id] {
						if err := SetFanSpeed(s, id, val); err != nil {
							log.Printf("写入串口失败: %v", err)
							return
						}
						lastPwms[id] = val
					}
				}
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()

	// 协程 2: 监听 Pico 串口 -> 分发回驱动对应的 fanX_input
	reader := bufio.NewReader(s)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			line, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("读取串口失败: %v", err)
				return
			}

			line = strings.TrimSpace(line)
			fmt.Println("Pico 反馈: ", line)
			if strings.HasPrefix(line, "[") {
				var fans []FanStatus
				if err := json.Unmarshal([]byte(line), &fans); err == nil {
					for _, fan := range fans {
						// 建立 ID 和 驱动文件的映射关系
						// 假设 fan1 -> fan1_input, fan2 -> fan2_input
						// 我们把 "fan1" 里的数字提取出来
						numStr := strings.TrimPrefix(fan.ID, "fan")
						rpmFile := filepath.Join(hwmonPath, fmt.Sprintf("fan%s_input", numStr))
						fmt.Println("更新驱动文件: ", rpmFile, " 为 RPM: ", fan.RPM)

						// 写入对应的驱动文件
						os.WriteFile(rpmFile, []byte(strconv.Itoa(fan.RPM)), 0644)
					}
				}
			}
		}
	}
}

func findHwmonPath() (string, error) {
	var foundPath string
	hwmonDir := homePath

	err := filepath.Walk(hwmonDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 我们只对 hwmonX 目录感兴趣
		if strings.HasPrefix(info.Name(), "hwmon") {
			markerFile := filepath.Join(path, "device", "marker")

			if fileExists(markerFile) {
				content, err := os.ReadFile(markerFile)
				if err == nil && strings.Contains(string(content), "vFanByTk") {
					foundPath = path        // 捕获路径
					return filepath.SkipDir // 找到了，停止继续扫描该目录
				}
			}
		}
		return nil
	})

	if foundPath == "" && err == nil {
		return "", fmt.Errorf("未找到带有 marker 的硬件驱动")
	}
	return foundPath, err
}

// 示例：向 Pico 发送设置转速的 JSON 指令
func SetFanSpeed(s *serial.Port, fanID string, pwmValue int) error {
	percent := int((float64(pwmValue) / 255.0) * 100)
	// 构造带 ID 的指令: {"fan": "fan1", "set_duty": 50}
	cmdObj := map[string]interface{}{
		"fan":      fanID,
		"set_duty": percent,
	}
	fmt.Printf("send to Pico: %+v\n", cmdObj)
	cmdBuf, _ := json.Marshal(cmdObj)
	_, err := s.Write(append(cmdBuf, '\n'))
	return err
}

func readIntFromFile(filePath string) int {
	// 1. 读取文件全部内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		// 如果读取失败（比如驱动还没加载），返回 0 或 -1 返回255，风扇最大
		return 255
	}

	// 2. 去除字符串两端的空白字符（如 \n \t 或空格）
	content := strings.TrimSpace(string(data))

	// 3. 将字符串转换为整数
	val, err := strconv.Atoi(content)
	if err != nil {
		return 255
	}

	return val
}
