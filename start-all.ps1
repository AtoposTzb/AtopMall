# 一键启动所有 AtopMall 微服务
# 启动前清理残留进程，启动间添加延时避免端口冲突

$root = $PSScriptRoot
$delaySeconds = 2

Write-Host "===== AtopMall 微服务启动脚本 =====" -ForegroundColor Cyan

# ===== 启动前清理：关闭可能残留的旧进程 =====
Write-Host "检查并清理残留进程..." -ForegroundColor Yellow

# 递归获取子进程
function Get-ChildProcesses {
    param([int]$ParentId)
    $children = Get-CimInstance Win32_Process | Where-Object { $_.ParentProcessId -eq $ParentId }
    foreach ($child in $children) {
        $child
        Get-ChildProcesses -ParentId $child.ProcessId
    }
}

# 找到所有由 start-all.ps1 启动的 powershell 窗口并清理
$shellProcs = Get-CimInstance Win32_Process | Where-Object {
    $_.Name -eq "powershell.exe" -and (
        $_.CommandLine -like "*atopmall_srvs*" -or
        $_.CommandLine -like "*atopmall_web*"
    )
}

foreach ($p in $shellProcs) {
    $children = Get-ChildProcesses -ParentId $p.ProcessId
    foreach ($child in $children) {
        Stop-Process -Id $child.ProcessId -Force -ErrorAction SilentlyContinue
    }
    Stop-Process -Id $p.ProcessId -Force -ErrorAction SilentlyContinue
    Write-Host "  清理残留窗口: PID $($p.ProcessId)" -ForegroundColor DarkGray
}

# 等待端口释放
Start-Sleep -Seconds 3

# ===== Python 微服务 =====
Write-Host "`n启动 Python 微服务..." -ForegroundColor Green
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$root\atopmall_srvs\goods_srv'; python server.py"
Write-Host "  ✓ goods_srv" -ForegroundColor DarkGray
Start-Sleep -Seconds $delaySeconds

Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$root\atopmall_srvs\inventory_srv'; python server.py"
Write-Host "  ✓ inventory_srv" -ForegroundColor DarkGray
Start-Sleep -Seconds $delaySeconds

Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$root\atopmall_srvs\order_srv'; python server.py"
Write-Host "  ✓ order_srv" -ForegroundColor DarkGray
Start-Sleep -Seconds $delaySeconds

Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$root\atopmall_srvs\user_srv'; python server.py"
Write-Host "  ✓ user_srv" -ForegroundColor DarkGray
Start-Sleep -Seconds $delaySeconds

Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$root\atopmall_srvs\userop_srv'; python server.py"
Write-Host "  ✓ userop_srv" -ForegroundColor DarkGray
Start-Sleep -Seconds $delaySeconds

# ===== Go Web 微服务（使用 air 热重载） =====
Write-Host "`n启动 Go Web 微服务..." -ForegroundColor Green
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$root\atopmall_web\goods_web'; air"
Write-Host "  ✓ goods_web" -ForegroundColor DarkGray
Start-Sleep -Seconds $delaySeconds

Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$root\atopmall_web\order_web'; air"
Write-Host "  ✓ order_web" -ForegroundColor DarkGray
Start-Sleep -Seconds $delaySeconds

Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$root\atopmall_web\oss_web'; air"
Write-Host "  ✓ oss_web" -ForegroundColor DarkGray
Start-Sleep -Seconds $delaySeconds

Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$root\atopmall_web\user_web'; air"
Write-Host "  ✓ user_web" -ForegroundColor DarkGray
Start-Sleep -Seconds $delaySeconds

Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$root\atopmall_web\userop_web'; air"
Write-Host "  ✓ userop_web" -ForegroundColor DarkGray

Write-Host "`n===== 所有微服务窗口已启动 =====" -ForegroundColor Cyan
Write-Host "提示：请检查各窗口是否有报错信息" -ForegroundColor Yellow
