# 快速修复器

  快速修复微服务容器中的JAVA类、JAR、SQL和前端代码

## 导入

```go
import "github.com/zhouxiaoxiang/repairman/v5"
```

## 示例

```go
man := repairman.NewRepairman()
man.RepairWeb()
man.RepairJar()
```

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
