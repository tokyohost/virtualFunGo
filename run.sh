# 1. 重新加载 Systemd 配置
sudo systemctl daemon-reload

# 2. 设置开机自启
sudo systemctl enable pico-fan.service

# 3. 立即启动服务
sudo systemctl start pico-fan.service

# 4. 查看运行状态
sudo systemctl status pico-fan.service