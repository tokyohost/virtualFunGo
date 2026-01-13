# 进入你的 go 项目目录
go build -o pico-fan-bridge .
# 将生成的文件移动到系统程序目录
sudo mv pico-fan-bridge /usr/local/bin/
# 确保它有执行权限
sudo chmod +x /usr/local/bin/pico-fan-bridge