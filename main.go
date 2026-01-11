package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const pwmDataStart = "pwm1"
const homePath = "/sys/class/hwmon/"

func readPwmValue(pwmFile string) (int, error) {
	data, err := ioutil.ReadFile(pwmFile)
	if err != nil {
		return 0, err
	}

	var pwm int
	_, err = fmt.Sscanf(string(data), "%d", &pwm)
	if err != nil {
		return 0, err
	}
	return pwm, nil
}

func main() {
	for {
		scanHwmonDirectories()

		time.Sleep(1 * time.Second)
	}
}

// checkMarkerFile 读取 marker 文件并检查其内容是否以 vFanByTk 开头
// 如果满足条件，返回 hwmon 目录的路径
func checkMarkerFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 读取文件的第一行并检查
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		content := scanner.Text()
		if strings.HasPrefix(content, "vFanByTk") {
			// 如果内容以 vFanByTk 开头，返回 hwmon 目录路径
			return filepath.Dir(filepath.Dir(filename)), nil
		} else {
			return "", nil // 内容不符合条件，返回空字符串
		}
	} else if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil
}

// fileExists 检查文件是否存在
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func scanHwmonDirectories() error {
	// 获取 /sys/class/hwmon 目录下的所有子目录
	hwmonDir := homePath
	err := filepath.Walk(hwmonDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		//fmt.Printf("check device path: %s\n", path)
		//fmt.Printf("check info: %s\n", info.Name())
		// 只处理 hwmon 目录下的子目录，并忽略其他类型的文件
		if strings.HasPrefix(info.Name(), "hwmon") {
			// 构建 marker 文件的路径
			markerFile := filepath.Join(path, "device", "marker")
			//log.Printf("check marker file: %s\n", markerFile)
			if fileExists(markerFile) {
				// 如果文件存在，检查其内容
				if hwmonPath, err := checkMarkerFile(markerFile); err == nil {
					// 如果返回的路径不为空，打印路径
					fmt.Printf("Found marker in device path: %s\n", hwmonPath)

					pwmValue, err := readPwmValue(filepath.Join(hwmonPath, pwmDataStart))
					if err != nil {
						log.Printf("Error reading PWM value: %v", err)
					} else {
						fmt.Printf("Current PWM value: %d\n", pwmValue)
					}
				} else {
					// 如果有错误打印错误信息
					fmt.Printf("Error reading marker file %s: %v\n", markerFile, err)
				}
			} else {
				//fmt.Printf("No device path found in %s\n", markerFile)
			}
		}
		return nil
	})

	return err
}
