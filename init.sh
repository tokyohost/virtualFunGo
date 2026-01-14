sudo apt remove golang-go
sudo apt autoremove
# 关键一步：删除可能残留在 /usr/lib 或 /usr/bin 的旧软链接
sudo rm -rf /usr/lib/go-1.19
sudo rm -rf /usr/bin/go


# 下载
wget https://go.dev/dl/go1.22.5.linux-amd64.tar.gz

# 解压到 /usr/local (如果之前有旧的，先删掉)
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz

# 临时在当前窗口生效
export PATH=/usr/local/go/bin:$PATH

# 永久生效（写入 .bashrc）
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 检查版本，必须显示 1.22.5
go version

# 回到项目目录
cd ~/virtual_fun/virtual_fan_go/virtualFunGo

# 重新编译
go build -o pico-fan-bridge .