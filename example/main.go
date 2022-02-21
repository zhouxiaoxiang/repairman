package main

import (
	"github.com/pieterclaerhout/go-log"
	"github.com/zhouxiaoxiang/repairman/v5"
)

const (
	APP     = "repairman"
	VERSION = "V5.1.0"
)

func main() {
	log.PrintColors = true
	log.Info("######################################")
	log.Infof("%15s:%s", APP, VERSION)
	log.Info("######################################")
	man := repairman.NewRepairman()
	man.RepairWeb()
	man.RepairJar()
}
