# ERP 启动基础服务脚本
# 启动 Redis, Etcd, DTM, MySQL 等基础服务的 exe 文件

# 基础服务 exe 路径配置（请根据实际安装路径修改）
$services = @{
    "Redis" = @{
        "Path" = "F:\Redis服务器-x64-5.0.10\redis-server.exe"
        "Args" = ""
        "WorkDir" = "F:\Redis服务器-x64-5.0.10"
    }
    "Etcd" = @{
        "Path" = "F:\etcd-v3.5.14-windows-amd64\etcd.exe"
        "Args" = ""
        "WorkDir" = "F:\etcd-v3.5.14-windows-amd64"
    }
    "DTM" = @{
        "Path" = "G:\项目\go\dtm_1.19.0_windows_amd64\dtm.exe"
        "Args" = " -c conf.yaml"
        "WorkDir" = "G:\项目\go\dtm_1.19.0_windows_amd64"
    }
     "Nginx" = @{
        "Path" = "F:\nginx-1.28.0\nginx.exe"
        "Args" = " -s reload"
        "WorkDir" = "F:\nginx-1.28.0"
    }
}

$missingServices = @()
$startedServices = @()

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "开始启动 ERP 基础服务..." -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

foreach ($serviceName in $services.Keys) {
    $service = $services[$serviceName]
    $exePath = $service.Path
    $args = $service.Args
    $workDir = $service.WorkDir
    
    Write-Host "[$serviceName] 检查服务..." -ForegroundColor Yellow
    
    if (Test-Path $exePath) {
        try {
            Write-Host "[启动] $exePath" -ForegroundColor Green
            
            # Nginx 需要特殊处理，不使用 WindowStyle
            if ($serviceName -eq "Nginx") {
                if ($args) {
                    Start-Process -FilePath $exePath -ArgumentList $args -WorkingDirectory $workDir -NoNewWindow
                } else {
                    Start-Process -FilePath $exePath -WorkingDirectory $workDir -NoNewWindow
                }
            } else {
                if ($args) {
                    Start-Process -FilePath $exePath -ArgumentList $args -WorkingDirectory $workDir -WindowStyle Minimized
                } else {
                    Start-Process -FilePath $exePath -WorkingDirectory $workDir -WindowStyle Minimized
                }
            }
            
            $startedServices += $serviceName
            Write-Host "[成功] $serviceName 启动成功" -ForegroundColor Green
        } catch {
            Write-Host "[错误] $serviceName 启动失败: $_" -ForegroundColor Red
        }
    } else {
        Write-Host "[缺失] $exePath" -ForegroundColor Red
        $missingServices += @{
            "Name" = $serviceName
            "Path" = $exePath
        }
    }
    Write-Host ""
}

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "启动完成!" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "已启动服务数量: $($startedServices.Count)" -ForegroundColor Green
Write-Host "缺失服务数量: $($missingServices.Count)" -ForegroundColor Yellow
Write-Host ""

if ($startedServices.Count -gt 0) {
    Write-Host "已启动的服务:" -ForegroundColor Green
    foreach ($service in $startedServices) {
        Write-Host "  ✓ $service" -ForegroundColor Green
    }
    Write-Host ""
}

if ($missingServices.Count -gt 0) {
    Write-Host "缺失的服务 EXE 文件:" -ForegroundColor Red
    foreach ($service in $missingServices) {
        Write-Host "  ✗ $($service.Name): $($service.Path)" -ForegroundColor Red
    }
    Write-Host ""
    Write-Host "提示: 请修改脚本中的路径配置，或安装相应的服务" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "服务访问信息:" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Redis:  localhost:6379" -ForegroundColor White
Write-Host "Etcd:   localhost:2379" -ForegroundColor White
Write-Host "DTM:    localhost:36790" -ForegroundColor White
Write-Host "Nginx:  localhost:8443" -ForegroundColor White
Write-Host ""

Write-Host "按任意键退出..." -ForegroundColor Gray
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
