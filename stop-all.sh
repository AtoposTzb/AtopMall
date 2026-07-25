#!/bin/bash
root="/root/AtopMall"
echo -e "\033[36m===== 停止 AtopMall 所有微服务 =====\033[0m"

# 第一步：尝试关闭tmux会话
if tmux has-session -t atop-session 2>/dev/null; then
    tmux kill-session -t atop-session
    echo -e "\033[32m✅ 成功销毁 atop-session tmux会话\033[0m"
else
    echo -e "\033[33m⚠️  无法检测到tmux会话，进入进程兜底清理\033[0m"
fi

# 第二步：兜底强制清理本项目进程（核心，防止tmux检测失效残留服务）
echo -e "\033[36m执行进程兜底清理...\033[0m"
pgrep -f "${root}.*python server.py" | xargs -r kill -9
pgrep -f "${root}.*air" | xargs -r kill -9
echo -e "\033[32m✅ 所有AtopMall相关进程清理完成\033[0m"