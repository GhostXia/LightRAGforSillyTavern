package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/HerSophia/go-lightrag/internal/api"
	"github.com/HerSophia/go-lightrag/internal/config"
	"go.uber.org/zap"
)

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "config.yaml", "配置文件路径")
	host := flag.String("host", "", "服务器主机地址")
	port := flag.Int("port", 0, "服务器端口")
	workingDir := flag.String("working-dir", "", "RAG工作目录")
	inputDir := flag.String("input-dir", "", "输入文档目录")
	model := flag.String("model", "", "LLM模型名称")
	embeddingModel := flag.String("embedding-model", "", "嵌入模型名称")
	logLevel := flag.String("log-level", "", "日志级别")
	apiKey := flag.String("key", "", "API密钥")
	flag.Parse()

	// 初始化配置
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Printf("加载配置文件失败: %v\n", err)
		os.Exit(1)
	}

	// 命令行参数覆盖配置文件
	if *host != "" {
		cfg.Server.Host = *host
	}
	if *port != 0 {
		cfg.Server.Port = *port
	}
	if *workingDir != "" {
		cfg.Rag.WorkingDir = *workingDir
	}
	if *inputDir != "" {
		cfg.Rag.InputDir = *inputDir
	}
	if *model != "" {
		cfg.Llm.Model = *model
	}
	if *embeddingModel != "" {
		cfg.Embedding.Model = *embeddingModel
	}
	if *logLevel != "" {
		cfg.Log.Level = *logLevel
	}
	if *apiKey != "" {
		cfg.Server.ApiKey = *apiKey
	}

	// 初始化日志
	logger, err := initLogger(cfg.Log.Level)
	if err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// 创建并启动服务器
	server := api.NewServer(cfg, logger)
	go func() {
		if err := server.Start(); err != nil {
			logger.Fatal("服务器启动失败", zap.Error(err))
		}
	}()

	logger.Info("服务器已启动",
		zap.String("host", cfg.Server.Host),
		zap.Int("port", cfg.Server.Port))

	// 等待中断信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Info("正在关闭服务器...")
	server.Shutdown()
	logger.Info("服务器已关闭")
}

func initLogger(level string) (*zap.Logger, error) {
	var zapLevel zap.AtomicLevel
	switch level {
	case "debug":
		zapLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		zapLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		zapLevel = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		zapLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		zapLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	cfg := zap.Config{
		Level:            zapLevel,
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return cfg.Build()
}
