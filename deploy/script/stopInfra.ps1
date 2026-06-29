# ERP 停止基础服务脚本
# 停止 Redis, Etcd, DTM, MySQL 等基础服务进程

# 基础服务进程名称（不包含 .exe 后缀）
$serviceProcesses = @(
    "redis-server",
    "etcd",
    "dtm",
    "nginx"
)

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "开始停止 ERP 基础服务..." -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$stoppedCount = 0
$notFoundCount = 0

foreach ($processName in $serviceProcesses) {
    Write-Host "[检查] $processName..." -ForegroundColor Yellow
    
    # Nginx 需要特殊处理，使用 taskkill 停止所有进程
    if ($processName -eq "nginx") {
        $processes = Get-Process -Name $processName -ErrorAction SilentlyContinue
        if ($processes) {
            try {
                Write-Host "[停止] 使用 taskkill 停止所有 nginx 进程" -ForegroundColor Green
                taskkill /F /IM nginx.exe /T 2>$null
                if ($LASTEXITCODE -eq 0) {
                    $stoppedCount += $processes.Count
                    Write-Host "[成功] 已停止 $($processes.Count) 个 nginx 进程" -ForegroundColor Green
                } else {
                    Write-Host "[错误] 停止 nginx 失败" -ForegroundColor Red
                }
                Start-Sleep -Milliseconds 500
            } catch {
                Write-Host "[错误] 无法停止 nginx: $_" -ForegroundColor Red
            }
        } else {
            Write-Host "[未运行] $processName" -ForegroundColor Gray
            $notFoundCount++
        }
    } else {
        # 其他服务使用标准方式停止
        $processes = Get-Process -Name $processName -ErrorAction SilentlyContinue
        
        if ($processes) {
            foreach ($proc in $processes) {
                try {
                    Write-Host "[停止] $processName (PID: $($proc.Id))" -ForegroundColor Green
                    Stop-Process -Id $proc.Id -Force
                    $stoppedCount++
                    Start-Sleep -Milliseconds 500
                } catch {
                    Write-Host "[错误] 无法停止 $processName (PID: $($proc.Id)): $_" -ForegroundColor Red
                }
            }
        } else {
            Write-Host "[未运行] $processName" -ForegroundColor Gray
            $notFoundCount++
        }
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

if ($stoppedCount -gt 0) {
    Write-Host "提示: 基础服务已停止" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "按任意键退出..." -ForegroundColor Gray
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
