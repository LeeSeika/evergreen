package main

import (
	"context"
	"evergreen/dao/mysql"
	"evergreen/dao/redis"
	"evergreen/logger"
	"evergreen/pkg/snowflake"
	"evergreen/router"
	"evergreen/settings"
	"evergreen/util"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	var configFileName string
	flag.StringVar(&configFileName, "conf", "./conf/config.yaml", "程序读取的配置文件全路径")

	err := settings.Init(configFileName)
	if err != nil {
		fmt.Printf("settings init error:%s\n", err)
		return
	}

	err = logger.Init()
	if err != nil {
		fmt.Printf("logger init error:%s\n", err)
		return
	}
	defer zap.L().Sync()

	err = mysql.Init()
	if err != nil {
		fmt.Printf("logger init error:%s\n", err)
		return
	}
	defer mysql.Close()

	err = redis.Init()
	if err != nil {
		fmt.Printf("redis init error:%s\n", err)
		return
	}

	err = snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineId)
	if err != nil {
		fmt.Printf("snowflake init error:%s\n", err)
		return
	}

	err = util.InitTrans("zh")
	if err != nil {
		fmt.Printf("translator init error:%s\n", err)
		return
	}

	engine := router.Setup()

	port := settings.Conf.Port
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: engine,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: %s\n", zap.Error(err))
		}
	}()
	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	log.Println("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

}
