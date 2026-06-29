# MySQL数据库导出脚本
# 导出所有erp_*数据库的表结构和数据到deploy/sql/all目录

param(
    [string]$MysqlHost = "localhost",
    [string]$MysqlPort = "3306",
    [string]$MysqlUser = "root",
    [string]$MysqlPassword = "123456",
    [string]$OutputBaseDir = "g:\项目\erp\deploy\sql\all"
)

# 设置控制台编码为UTF-8
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
$OutputEncoding = [System.Text.Encoding]::UTF8

Write-Host "=== MySQL ERP数据库导出工具 ===" -ForegroundColor Green
Write-Host "主机: ${MysqlHost}:${MysqlPort}" -ForegroundColor Cyan
Write-Host "用户: $MysqlUser" -ForegroundColor Cyan
Write-Host "输出目录: $OutputBaseDir" -ForegroundColor Cyan
Write-Host ""

# 检查mysqldump是否可用
try {
    $null = Get-Command mysqldump -ErrorAction Stop
    Write-Host "[√] mysqldump 工具已找到" -ForegroundColor Green
} catch {
    Write-Host "[×] 错误: 未找到mysqldump工具，请确保MySQL已安装并添加到PATH环境变量" -ForegroundColor Red
    exit 1
}

# 检查mysql客户端是否可用
try {
    $null = Get-Command mysql -ErrorAction Stop
    Write-Host "[√] mysql 客户端已找到" -ForegroundColor Green
} catch {
    Write-Host "[×] 错误: 未找到mysql客户端，请确保MySQL已安装并添加到PATH环境变量" -ForegroundColor Red
    exit 1
}

Write-Host ""

# 构建MySQL连接参数
$mysqlArgs = @(
    "-h", $MysqlHost,
    "-P", $MysqlPort,
    "-u", $MysqlUser
)

if ($MysqlPassword -ne "") {
    $mysqlArgs += "-p$MysqlPassword"
}

# 获取所有erp_*数据库
Write-Host "正在查询erp_*数据库..." -ForegroundColor Yellow
$query = "SHOW DATABASES LIKE 'erp_%';"
$databases = & mysql @mysqlArgs -N -e $query 2>&1

if ($LASTEXITCODE -ne 0) {
    Write-Host "[×] 错误: 无法连接到MySQL数据库" -ForegroundColor Red
    Write-Host $databases -ForegroundColor Red
    exit 1
}

# 过滤空行
$databases = $databases | Where-Object { $_ -match '\S' }

if ($databases.Count -eq 0) {
    Write-Host "[!] 未找到任何erp_*数据库" -ForegroundColor Yellow
    exit 0
}

Write-Host "[√] 找到 $($databases.Count) 个数据库:" -ForegroundColor Green
$databases | ForEach-Object { Write-Host "  - $_" -ForegroundColor Cyan }
Write-Host ""

# 创建输出基础目录
if (-not (Test-Path $OutputBaseDir)) {
    New-Item -ItemType Directory -Path $OutputBaseDir -Force | Out-Null
    Write-Host "[√] 创建输出目录: $OutputBaseDir" -ForegroundColor Green
}

# 导出每个数据库
$successCount = 0
$failCount = 0

foreach ($db in $databases) {
    $db = $db.Trim()
    if ([string]::IsNullOrWhiteSpace($db)) {
        continue
    }
    
    Write-Host "----------------------------------------" -ForegroundColor Gray
    Write-Host "正在导出数据库: $db" -ForegroundColor Yellow
    
    # 创建数据库对应的目录
    $dbDir = Join-Path $OutputBaseDir $db
    if (-not (Test-Path $dbDir)) {
        New-Item -ItemType Directory -Path $dbDir -Force | Out-Null
        Write-Host "  [√] 创建目录: $dbDir" -ForegroundColor Green
    }
    
    # 导出文件路径
    $outputFile = Join-Path $dbDir "$db.sql"
    
    # 构建mysqldump参数
    $dumpArgs = @(
        "-h", $MysqlHost,
        "-P", $MysqlPort,
        "-u", $MysqlUser
    )
    
    if ($MysqlPassword -ne "") {
        $dumpArgs += "-p$MysqlPassword"
    }
    
    # 添加导出选项
    $dumpArgs += @(
        "--single-transaction",      # 使用事务保证数据一致性
        "--routines",                 # 导出存储过程和函数
        "--triggers",                 # 导出触发器
        "--events",                   # 导出事件
        "--default-character-set=utf8mb4",  # 使用UTF8MB4字符集
        "--add-drop-database",        # 添加DROP DATABASE语句
        "--databases",                # 指定数据库模式
        $db
    )
    
    # 执行导出
    Write-Host "  正在导出到: $outputFile" -ForegroundColor Cyan
    
    try {
        & mysqldump @dumpArgs | Out-File -FilePath $outputFile -Encoding UTF8
        
        if ($LASTEXITCODE -eq 0) {
            $fileSize = (Get-Item $outputFile).Length
            $fileSizeKB = [math]::Round($fileSize / 1KB, 2)
            Write-Host "  [√] 导出成功! 文件大小: $fileSizeKB KB" -ForegroundColor Green
            $successCount++
        } else {
            Write-Host "  [×] 导出失败!" -ForegroundColor Red
            $failCount++
        }
    } catch {
        Write-Host "  [×] 导出异常: $($_.Exception.Message)" -ForegroundColor Red
        $failCount++
    }
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Gray
Write-Host "导出完成!" -ForegroundColor Green
Write-Host "成功: $successCount 个数据库" -ForegroundColor Green
if ($failCount -gt 0) {
    Write-Host "失败: $failCount 个数据库" -ForegroundColor Red
}
Write-Host "输出目录: $OutputBaseDir" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Gray
