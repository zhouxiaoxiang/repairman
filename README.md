# 快速修复器

  快速修复微服务容器中的CLASS、JAR、SQL
  
## 获取发布版

  [下载](https://media.githubusercontent.com/media/zhouxiaoxiang/tools/main/repairman)

  如下载失败，可以直接从 `https://github.com/zhouxiaoxiang/tools/`中下载repairman

## 用法

### class

  按照包名放入`BOOT-INF/classes/`

### sql

  sql脚本放入`BOOT-INF/classes/db/migration/`

### jar

  将jar文件加到`BOOT-INF/lib/`

### 前端代码

  `client-plat-1.1.10.tar.gz` 或者 `client-plat-1.1.10`
  
  错误格式会忽略哦

### 开始修复

  `./repairman`

  如果存在JAVA目录，程序会提示用户选择微服务容器

## 以下为编程方式，普通用户可以忽略

### 导入

```go
import "github.com/zhouxiaoxiang/repairman/v5"
```

### 示例

```go
man := repairman.NewRepairman()
man.RepairWeb()
man.RepairJar()
```
