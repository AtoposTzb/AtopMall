﻿# 一键关闭所有 AtopMall 微服务
# 关闭 PowerShell 窗口及其所有子进程（python/air）

Write-Host "===== 正在关闭 AtopMall 微服务 =====" -ForegroundColor Cyan

# 递归获取子进程
function Get-ChildProcesses {
    param([int]$ParentId)
    $children = Get-CimInstance Win32_Process | Where-Object { $_.ParentProcessId -eq $ParentId }
    foreach ($child in $children) {
        $child
        Get-ChildProcesses -ParentId $child.ProcessId
    }
}

# 找到所有由 start-all.ps1 启动的 powershell 窗口
$shellProcs = Get-CimInstance Win32_Process | Where-Object {
    $_.Name -eq "powershell.exe" -and (
        $_.CommandLine -like "*atopmall_srvs*" -or
        $_.CommandLine -like "*atopmall_web*"
    )
}

$shellCount = $shellProcs.Count
$childCount = 0

foreach ($p in $shellProcs) {
    # 先杀所有子进程
    $children = Get-ChildProcesses -ParentId $p.ProcessId
    foreach ($child in $children) {
        Stop-Process -Id $child.ProcessId -Force -ErrorAction SilentlyContinue
        $childCount++
    }
    # 再杀 PowerShell 窗口本身
    Stop-Process -Id $p.ProcessId -Force -ErrorAction SilentlyContinue
}

$total = $shellCount + $childCount

if ($total -eq 0) {
    Write-Host "未发现运行中的微服务窗口。" -ForegroundColor Yellow
} else {
    Write-Host "===== 已关闭 $shellCount 个窗口 + $childCount 个子进程，共 $total 个 =====" -ForegroundColor Green
}
