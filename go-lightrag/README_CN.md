# Go-LightRAG 中文说明

## 项目简介

这是使用Go语言重构的LightRAG项目，旨在提供更高效、更简洁的RAG（检索增强生成）系统实现。原始项目基于Python，本项目保留了核心功能，同时利用Go语言的并发性能和编译优势，提供更好的性能和更简单的部署体验。

## 重构优势

1. **更高的性能**：利用Go语言的并发特性，提供更快的响应速度和更高的吞吐量
2. **更低的资源占用**：Go语言的内存管理更加高效，减少资源消耗
3. **更简单的部署**：编译为单一二进制文件，无需复杂的依赖管理
4. **保持API兼容性**：与原Python版本保持API兼容，方便迁移

## 安装说明

### 从源码编译

1. 确保已安装Go 1.21或更高版本
2. 克隆仓库
   ```bash
   git clone https://github.com/HerSophia/go-lightrag.git
   cd go-lightrag
   ```
3. 编译项目
   ```bash
   go build -o lightrag ./cmd/server
   ```

### 使用预编译二进制文件

1. 从[发布页面](https://github.com/HerSophia/go-lightrag/releases)下载适合您系统的二进制文件
2. 解压缩文件
3. 将二进制文件放置在您选择的目录中

## 配置说明

1. 复制示例配置文件
   ```bash
   cp config.yaml.example config.yaml
   ```
2. 编辑配置文件，设置您的OpenAI API密钥和其他参数
   ```yaml
   # LLM配置
   llm:
     provider: "openai"
     model: "gpt-4o-mini"
     api_key: "您的OpenAI API密钥"
   
   # 嵌入模型配置
   embedding:
     provider: "openai"
     model: "text-embedding-3-large"
     api_key: "您的OpenAI API密钥"
   ```

## 运行说明

```bash
# 使用默认配置文件运行
./lightrag

# 指定配置文件路径
./lightrag --config /path/to/config.yaml

# 指定服务器主机和端口
./lightrag --host 127.0.0.1 --port 8080

# 指定工作目录
./lightrag --working-dir /path/to/rag_storage
```

## API使用说明

### 查询接口

```bash
curl -X POST "http://localhost:9621/api/query" \
  -H "Content-Type: application/json" \
  -d '{"query":"什么是RAG技术?", "mode":"hybrid"}'
```

### 插入文本接口

```bash
curl -X POST "http://localhost:9621/api/insert" \
  -H "Content-Type: application/json" \
  -d '{"text":"RAG（检索增强生成）是一种结合检索系统和生成式AI的技术..."}'
```

### OpenAI兼容接口

```bash
curl -X POST "http://localhost:9621/v1/chat/completions" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4",
    "messages": [
      {"role": "system", "content": "你是一个有用的助手。"},
      {"role": "user", "content": "什么是RAG技术?"}
    ]
  }'
```

## 与原项目的区别

- 使用Go语言的强类型系统提供更好的代码可维护性
- 利用Go的goroutine和channel提供更好的并发性能
- 编译为单一二进制文件，简化部署
- 更高效的内存管理
- 保持API兼容性，方便迁移

## 常见问题

1. **Q: 如何切换到不同的存储后端？**  
   A: 在配置文件中修改storage部分的配置，支持memory、neo4j、postgres等多种后端。

2. **Q: 如何使用非OpenAI的模型？**  
   A: 在配置文件中修改llm和embedding部分的provider参数，目前支持openai、ollama等提供商。

3. **Q: 如何调整RAG的检索参数？**  
   A: 在配置文件的rag部分修改相关参数，如chunk_size、top_k等。

## 贡献指南

欢迎提交Pull Request或Issue来帮助改进项目！