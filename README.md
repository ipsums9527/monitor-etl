# Monitor ETL

一个用于适配`晤大浪事 Netdata版 监控副屏`监控网络设备流量数据的 ETL (Extract, Transform, Load) 工具，支持爱快（iKuai）和 RouterOS 设备。

## 功能特性

- 🔌 支持多种路由器 API
  - 爱快（iKuai）路由器
  - MikroTik RouterOS
  - (更多适配中。。。)
- 📊 实时流量监控和数据采集
- 🐳 Docker 容器化部署
- ⚙️ 灵活的配置管理

## 快速开始

### 前置要求
- Go 1.25+
- Docker（推荐）

### 安装和运行
1. 克隆项目

```bash
git clone https://github.com/ipsums9527/monitor-etl.git cd monitor-etl
```

2. 配置文件

```bash
cp config.example.yml config.yml
```
编辑 `config.yml` 配置您的路由器信息：
```yaml
listen: "0.0.0.0"
port: 19999

api:
  type: "ikuai"
  host: "http://192.168.9.1"
  user: "admin"
  password: "123456"

#  type: "routeros"
#  host: "tcp://192.168.9.1:8728"
#  user: "admin"
#  password: "123456"
#  ethers:
#    - name: "ether-cm"
#      isInvert: false
#    - name: "ether-ct"
#      isInvert: false

```

3. 运行应用

**使用 Docker：**

```bash 
make all
``` 
docker-compose.yml
```
version: "3.9"
services:
  Monitor-ETL:
    image: ghcr.io/ipsums9527/monitor-etl:dev
    container_name: monitor-etl
    restart: unless-stopped
    volumes:
      - ./config.yml:/app/config.yml
    ports:
      - "19999:19999"
```

## 配置说明

### 爱快路由器配置
```yaml
listen: "0.0.0.0"
port: 19999

api:
  type: "ikuai"
  host: "http://192.168.9.1"
  user: "admin"
  password: "123456"
```

### RouterOS 配置
```yaml
listen: "0.0.0.0"
port: 19999

api:
  type: "routeros"
  host: "tcp://192.168.9.1:8728"
  user: "admin"
  password: "123456"
  ethers:
    - name: "ether-cm"
      isInvert: false
    - name: "ether-ct"
      isInvert: false
```

项目结构

``` 
monitor-etl/
├── api/            # API 客户端实现
│   ├── ikuai/      # 爱快路由器 API
│   ├── ros/        # RouterOS API
│   └── api.go
├── app/
│   └── server/     # HTTP 服务器
├── config/         # 配置管理
├── control/        # 控制逻辑
├── model/          # 数据模型
├── main.go         # 入口文件
├── config.yml      # 配置文件
└── Dockerfile      # Docker 配置
```

## API 接口
服务启动后，默认监听在 http://0.0.0.0:19999

## 开发
添加新的路由器支持  
在 api/ 目录下创建新的路由器适配器  
实现 `api.SystemDataClient` 接口  
在配置中添加新的类型支持

## 贡献
欢迎提交 Issue 和 Pull Request！
