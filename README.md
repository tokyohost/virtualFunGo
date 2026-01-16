
# VirtualFunGo

`VirtualFunGo` is a **bridge service written in Go** that connects the **virtual_fan Linux kernel driver** with **external fan control hardware**, such as a **Raspberry Pi Pico**, via USB.

It acts as the **missing link** between:

- Linux virtual fan values (`hwmon`)
- User-space control logic
- Real, external fan speed controllers

This project is designed to be used together with  
👉 **virtual_fun_kernel**  
https://github.com/tokyohost/virtual_fun_kernel

---

## ✨ Features

- ✅ Written in **Go**, lightweight and efficient
- ✅ Reads fan speed values from Linux `hwmon`
- ✅ Communicates with **Raspberry Pi Pico** over USB (CDC / serial)
- ✅ Translates virtual fan values into real hardware control signals
- ✅ Suitable for:
    - Mini PCs
    - NAS devices
    - Industrial / embedded systems
    - Motherboards without fan control support

---

## 🧩 Architecture Overview

```
+----------------------+
|   Linux Kernel       |
|  virtual_fan module  |
+----------+-----------+
           |
           | hwmon (/sys)
           v
+----------------------+
|     VirtualFunGo     |
|  (Go bridge service) |
+----------+-----------+
           |
           | USB / Serial
           v
+----------------------+
| Raspberry Pi Pico    |
| Fan control firmware |
+----------+-----------+
           |
           v
+----------------------+
| Physical Fan (PWM)   |
+----------------------+
```

---

## 🧠 How It Works

1. `virtual_fun_kernel` exposes a virtual fan device via `hwmon`
2. VirtualFunGo:
    - Monitors fan speed / target values from `/sys/class/hwmon`
    - Applies custom logic (mapping, scaling, limits)
3. Fan speed values are sent over **USB serial**
4. Raspberry Pi Pico:
    - Receives commands
    - Generates PWM signals
    - Controls real fan hardware

---

## 📦 Requirements

### Linux Side

- Linux with `virtual_fun_kernel` installed
- Go ≥ 1.24
- USB access permission (may require udev rules)

### Hardware Side

- Raspberry Pi Pico (2040)
- Pico firmware supporting USB CDC
- External fan (5V / 12V, PWM supported)

---

## 🚀 Installation

```bash
curl -fsSL https://raw.githubusercontent.com/tokyohost/virtualFunGo/master/install.sh

```

## 🧪 Example Use Cases

- Software fan curve on unsupported hardware
- Quiet NAS fan control
- External smart fan controller
- Temperature-based fan automation
---

## 📜 License

MIT License

---

## 🤝 Related Projects

- **virtual_fun_kernel**  
  https://github.com/tokyohost/virtual_fun_kernel
- **Pi Pico Fan microPython code**
  https://github.com/tokyohost/virtual_funPico
---

## 🙌 Contributing

Issues and pull requests are welcome.

If you are using this project in production or embedded systems, feedback is highly appreciated.