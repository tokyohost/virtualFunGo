package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

const pwmFile = "/sys/class/hwmon/hwmon4/pwm1"

func readPwmValue() (int, error) {
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
		pwmValue, err := readPwmValue()
		if err != nil {
			log.Printf("Error reading PWM value: %v", err)
		} else {
			fmt.Printf("Current PWM value: %d\n", pwmValue)
		}
		time.Sleep(1 * time.Second)
	}
}
