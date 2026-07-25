#!/bin/bash
# AtopMall Linux 一键微服务启动脚本
# 请先安装tmux工具
root="/root/AtopMall"
delaySeconds=2

echo -e "\033[36m===== AtopMall 微服务启动脚本 =====\033[0m"
echo "请先安装tmux工具，否则无法正常启动"
echo "安装tmux命令：apt install tmux -y or yum install tmux -y"
# 1. 清理旧进程
echo -e "\033[33m检查并清理残留服务进程...\033[0m"
pkill -f "python server.py" 2>/dev/null
pkill -f "air" 2>/dev/null
sleep 3

# 关键修复：如果没有tmux服务，创建一个常驻主会话 atop-session
if ! tmux has-session -t atop-session 2>/dev/null; then
    tmux new -d -s atop-session
fi

# 2. 启动 Python SRV 微服务
echo -e "\n\033[32m启动 Python 微服务...\033[0m"
tmux new-window -t atop-session -n goods_srv -d "cd ${root}/atopmall_srvs/goods_srv && python server.py; bash"
echo -e "\033[90m  ✓ goods_srv\033[0m"
sleep ${delaySeconds}

tmux new-window -t atop-session -n inventory_srv -d "cd ${root}/atopmall_srvs/inventory_srv && python server.py; bash"
echo -e "\033[90m  ✓ inventory_srv\033[0m"
sleep ${delaySeconds}

tmux new-window -t atop-session -n order_srv -d "cd ${root}/atopmall_srvs/order_srv && python server.py; bash"
echo -e "\033[90m  ✓ order_srv\033[0m"
sleep ${delaySeconds}

tmux new-window -t atop-session -n user_srv -d "cd ${root}/atopmall_srvs/user_srv && python server.py; bash"
echo -e "\033[90m  ✓ user_srv\033[0m"
sleep ${delaySeconds}

tmux new-window -t atop-session -n userop_srv -d "cd ${root}/atopmall_srvs/userop_srv && python server.py; bash"
echo -e "\033[90m  ✓ userop_srv\033[0m"
sleep ${delaySeconds}

# 3. 启动 Go Web 微服务（air热重载）
echo -e "\n\033[32m启动 Go Web 微服务...\033[0m"
tmux new-window -t atop-session -n goods_web -d "cd ${root}/atopmall_web/goods_web && air; bash"
echo -e "\033[90m  ✓ goods_web\033[0m"
sleep ${delaySeconds}

tmux new-window -t atop-session -n order_web -d "cd ${root}/atopmall_web/order_web && air; bash"
echo -e "\033[90m  ✓ order_web\033[0m"
sleep ${delaySeconds}

tmux new-window -t atop-session -n oss_web -d "cd ${root}/atopmall_web/oss_web && air; bash"
echo -e "\033[90m  ✓ oss_web\033[0m"
sleep ${delaySeconds}

tmux new-window -t atop-session -n user_web -d "cd ${root}/atopmall_web/user_web && air; bash"
echo -e "\033[90m  ✓ user_web\033[0m"
sleep ${delaySeconds}

tmux new-window -t atop-session -n userop_web -d "cd ${root}/atopmall_web/userop_web && air; bash"
echo -e "\033[90m  ✓ userop_web\033[0m"

echo -e "\n\033[36m===== 所有微服务后台启动完成 =====\033[0m"
echo -e "\033[33m操作说明：\033[0m"
echo "  tmux ls                    查看会话"
echo "  tmux a -t atop-session     进入总会话查看所有窗口(建议xhsell执行)"
echo "  2.打开窗口可视化列表：Ctrl+b 松开 → w，方向键选择窗口回车进入"
echo "  Ctrl+b 然后 d 脱离窗口，服务后台持续运行"