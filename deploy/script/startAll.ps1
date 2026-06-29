# ERP 启动所有服务脚本
# 使用绝对路径启动所有 API 和 RPC 服务

$baseDir = "g:\项目\erp\app"
$modules = @("auth", "customer", "finance", "hr", "image", "inventory", "product", "production", "purchase", "sale", "supplier", "user")

$missingExes = @()
$startedExes = @()

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "开始启动 ERP 所有服务..." -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

foreach ($module in $modules) {
    # 检查并启动 API 服务
    $apiPath = Join-Path $baseDir "$module\api"
    $apiExes = Get-ChildItem -Path $apiPath -Filter "*.exe" -ErrorAction SilentlyContinue
    
    if ($apiExes) {
        foreach ($exe in $apiExes) {
            $exePath = $exe.FullName
            Write-Host "[启动] $exePath" -ForegroundColor Green
            Start-Process -FilePath $exePath -WorkingDirectory $apiPath
            $startedExes += $exePath
        }
    } else {
        $expectedPath = Join-Path $apiPath "$module.exe"
        Write-Host "[缺失] $expectedPath" -ForegroundColor Yellow
        $missingExes += $expectedPath
    }
    
    # 检查并启动 RPC 服务
    $rpcPath = Join-Path $baseDir "$module\rpc"
    
    # 特殊处理: auth 模块的 RPC 使用 xAuth.exe
    if ($module -eq "auth") {
        $specificExe = Join-Path $rpcPath "xAuth.exe"
        if (Test-Path $specificExe) {
            Write-Host "[启动] $specificExe" -ForegroundColor Green
            Start-Process -FilePath $specificExe -WorkingDirectory $rpcPath
            $startedExes += $specificExe
        } else {
            Write-Host "[缺失] $specificExe" -ForegroundColor Yellow
            $missingExes += $specificExe
        }
    } else {
        $rpcExes = Get-ChildItem -Path $rpcPath -Filter "*.exe" -ErrorAction SilentlyContinue
        
        if ($rpcExes) {
            foreach ($exe in $rpcExes) {
                $exePath = $exe.FullName
                Write-Host "[启动] $exePath" -ForegroundColor Green
                Start-Process -FilePath $exePath -WorkingDirectory $rpcPath
                $startedExes += $exePath
            }
        } else {
            $expectedPath = Join-Path $rpcPath "$module.exe"
            Write-Host "[缺失] $expectedPath" -ForegroundColor Yellow
            $missingExes += $expectedPath
        }
    }
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "启动完成!" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "已启动服务数量: $($startedExes.Count)" -ForegroundColor Green
Write-Host "缺失服务数量: $($missingExes.Count)" -ForegroundColor Yellow
Write-Host ""

if ($missingExes.Count -gt 0) {
    Write-Host "缺失的 EXE 文件列表:" -ForegroundColor Red
    foreach ($missing in $missingExes) {
        Write-Host "  - $missing" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "按任意键退出..." -ForegroundColor Gray
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
