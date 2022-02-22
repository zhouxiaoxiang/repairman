# 快速修复器

  快速修复微服务容器中的CLASS、JAR、SQL
  
## 获取发布版

`wget https://raw.githubusercontent.com/zhouxiaoxiang/tools/main/repairman`

## 用法

### class

  按照包名放入BOOT-INF/classes/

### sql

  sql脚本放入BOOT-INF/classes/db/migration/

### jar

  将jar文件加到BOOT-INF/lib/

### 前端代码

  client-plat-1.1.10.tar.gz 或者 client-plat-1.1.10

### 开始修复

  `./repairman`

## 编程方式

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
