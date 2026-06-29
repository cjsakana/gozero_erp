# ERP 停止所有服务脚本
# 停止所有 API 和 RPC 服务进程

$modules = @("auth", "customer", "finance", "hr", "image", "inventory", "product", "production", "purchase", "sale", "supplier", "user")

# 需要停止的进程名称列表
$processNames = @()

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "开始停止 ERP 所有服务..." -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 构建进程名称列表
foreach ($module in $modules) {
    $processNames += $module  # API 进程名
    if ($module -eq "auth") {
        $processNames += "xAuth"  # auth 的 RPC 使用 xAuth
    }
}

$stoppedCount = 0
$notFoundCount = 0

foreach ($processName in $processNames) {
    Write-Host "[检查] $processName..." -ForegroundColor Yellow
    
    $processes = Get-Process -Name $processName -ErrorAction SilentlyContinue
    
    if ($processes) {
        foreach ($proc in $processes) {
            try {
                Write-Host "[停止] $processName (PID: $($proc.Id))" -ForegroundColor Green
                Stop-Process -Id $proc.Id -Force
                $stoppedCount++
            } catch {
                Write-Host "[错误] 无法停止 $processName (PID: $($proc.Id)): $_" -ForegroundColor Red
            }
        }
    } else {
        Write-Host "[未运行] $processName" -ForegroundColor Gray
        $notFoundCount++
    }
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "停止完成!" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "已停止进程数量: $stoppedCount" -ForegroundColor Green
Write-Host "未运行进程数量: $notFoundCount" -ForegroundColor Gray
Write-Host ""

Write-Host "按任意键退出..." -ForegroundColor Gray
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
