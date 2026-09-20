<div align="center">
  <a href="https://goflys.cn/mqttw" target="_blank">
    <img alt="MQTTW Logo" width="150" src="https://doc.goflys.cn/image/mqttwicon.png"/>
  </a>
</div>

<div align="center">
  <h1>MQTTW</h1>
</div>

<div align="center">

MQTT 调试 · 压测一体化桌面客户端

</div>

<div align="center">

简体中文 | [English](./README_EN.md)


</div>

 

**MQTTW** 是一款基于 **Wails3 + Go + Vue3** 构建的专业 MQTT 桌面调试工具，全面支持 **MQTT 3.1 / 3.1.1 / 5.0** 协议，集连接调试、消息收发、多连接管理与大规模客户端压测于一体，适用于物联网（IoT）开发、Broker 联调验证与性能压测场景。

同时，MQTTW 也是一套**经典的 Wails3 桌面应用开发示例**—— 工程结构清晰、前后端交互完整，既可作为开发者学习 Wails3 的参考项目，也适合在此基础上定制 MQTT 工具或进行二次开发。

![MQTT客户端工具](https://docapi.goflys.cn/common/uploadfile/get_image?url=resource/uploads/20260919/ca9f5c4964706e1353fb424904b1dbca.png)
## 功能说明

### 一、多协议连接

- 一套界面兼容三代 MQTT 协议，自动识别协议版本
- **MQTT 3.1 / 3.1.1**：`github.com/eclipse/paho.mqtt.golang`（老 Paho 引擎）
- **MQTT 5.0**：`github.com/eclipse/paho.golang/autopaho`（新 Paho V5，断线自动恢复订阅）
- TLS/SSL 连接，支持 CA 证书、客户端证书、私钥、ALPN 配置
- Clean Session / Clean Start 会话控制
- 自动重连 + 可配置重连周期
- 连接超时、Keep Alive 心跳可调

### 二、MQTT 5.0 完整特性

- **用户属性（User Properties）**：自定义键值对
- **订阅标识符（Subscription Identifier）**：范围 1 ~ 268435455
- **主题别名（Topic Alias）**：降低带宽占用
- **会话过期（Session Expiry Interval）**
- **请求 / 响应（Request / Response）**：响应主题 + 关联数据
- **问题信息（Problem Information）**、**接收最大值 / 最大数据包**流控参数
- **遗嘱消息（Will）**：遗嘱主题 / QoS / Retain / Payload / Will Delay Interval / Content Type 等完整配置

### 三、连接管理

- 分组树形管理（分组 / 连接层级），支持拖拽排序
- 多连接同时在线，消息流相互独立
- 连接级未读计数 + 分组汇总角标
- 连接配置持久化到 SQLite，一键复用

### 四、订阅与消息

- 订阅持久化管理，重复订阅自动拦截
- 消息实时推送，事件驱动更新界面
- JSON 消息自动检测并高亮展示
- 消息异步入库 SQLite（带重试），可回溯查询
- MQTT 5 订阅选项：订阅标识符、用户属性等

### 五、消息发布

- QoS 0 / 1 / 2 + Retain 保留消息
- 主题别名、Payload 格式指示（UTF-8）
- 用户属性、响应主题、Content Type、消息过期、关联数据
- 完整 MQTT 5 发布属性面板

### 六、独立压测模块（stress）

- 固定 **MQTT 3.1**（`github.com/gonzalop/mq` 轻量库）与 3.1.1 / 5.0 选项
- 批量创建数千 ~ 数万客户端，连接速率可控
- 订阅 / 发布压测，速率与总量可控
- 实时指标：连接成功率、发布 / 接收速率、发布 / 接收总数
- **回显延迟分位**：P50 / P95 / P99
- 实时曲线（ECharts）+ 结果汇总 + CSV 导出
- 运行状态、进度实时可见；配置持久化，下次启动一键复用

### 七、设置面板

- **启动行为**：开机自启、启动时自动恢复上次在线连接、最小化到系统托盘
- **消息通知**：收到订阅消息时系统通知 + 声音提示开关
- **全局快捷键**：显示 / 隐藏窗口（Ctrl+Shift+H）、断开全部连接（Ctrl+Shift+D）、快速发布（Ctrl+Shift+P）
- **运行目录**：一键打开数据目录
- **国际语言**：简体中文 / 繁体中文 / 英语
- **关于**：版本信息、更新检查、GoFly 社区 / GMQT 链接

### 八、配套推荐

- **GMQT-Broker**：GoFly 社区出品的 MQTT 服务器端，可独立部署或安装到 GoFlyGen 框架，作为 MQTTW 的配套 Broker 即装即用（<https://goflys.cn/gmqt>）
- **GoFly 全栈开发社区**：Go 后端、Vue 前端、Wails 桌面应用、AI Agent 智能体开发、AI 编程框架、物联网开发经验与 MQTT 实战（<https://goflys.cn/>）

---

## 使用说明

### 环境要求

| 依赖 | 版本要求 |
| --- | --- |
| Go | 1.25+ |
| Node.js | 18+（含 npm） |
| Wails3 CLI | v3.0.0-beta.9+ |

### 安装与启动

```bash
# 1. 克隆项目
git clone https://gitee.com/huang_li_shi_admin/mqttw.git
cd mqttw

# 2. 安装依赖
go mod tidy
cd frontend && npm install && cd ..

# 3. 开发模式运行（热重载前后端）
wails3 dev

# 4. 生产构建
wails3 build
```

> 开发模式默认前端端口 `9245`；若端口被占用，清理残留进程后重新 `wails3 dev` 即可。

### 基本使用流程

1. **新建连接**：填写 Broker 地址、端口、Client ID，选择 MQTT 版本（3.1 / 3.1.1 / 5.0），配置 Clean Session / 自动重连 / 连接属性
2. **订阅主题**：在操作面板添加订阅（QoS、订阅标识符等），接收消息实时展示
3. **发布消息**：填写主题与 Payload，配置 QoS / Retain / MQTT 5 发布属性后发送
4. **压测**：进入压测模块，配置客户端数量、连接速率、订阅 / 发布参数与时长，启动后实时查看指标曲线
5. **数据存储**：连接、订阅、消息记录均持久化于 `resource/db/data.db`（SQLite）

---

## 开发技术栈

### 后端（Go）

| 组件 | 说明 |
| --- | --- |
| [Wails3](https://wails.io/) | v3.0.0-beta.9，桌面应用框架（窗口 / 托盘 / 全局快捷键 / 事件） |
| [paho.mqtt.golang](https://github.com/eclipse/paho.mqtt.golang) | v1.5.1，MQTT 3.1 / 3.1.1 客户端 |
| [paho.golang](https://github.com/eclipse/paho.golang) | v0.23.0（autopaho），MQTT 5.0 客户端 |
| [gonzalop/mq](https://github.com/gonzalop/mq) | v0.9.10，压测模块轻量客户端（MQTT 3.1） |
| [GORM](https://gorm.io/) + GORM Gen | ORM 与 DAO 代码生成 |
| [modernc.org/sqlite](https://modernc.org/sqlite) | 纯 Go SQLite 驱动，无 CGO |
| [mcp-go](https://github.com/mark3labs/mcp-go) | MCP 服务，gofly-remote-mcp 管理数据库 |

### 前端（Vue3）

| 组件 | 说明 |
| --- | --- |
| Vue 3 + TypeScript | 组合式 API |
| Vite | 构建工具 |
| Arco Design Vue | 组件库（含 icon-font 图标） |
| Pinia | 状态管理（连接 / 消息 / 压测 Store） |
| Vue Router | 路由（连接 / 压测 / 设置 / 帮助 / GMQT） |
| ECharts + vue-echarts | 压测实时曲线、指标可视化 |
| vue-i18n | 国际化（简中 / 繁中 / 英语） |
| Monaco Editor + PrismJS | 消息 / JSON 高亮展示 |
| SortableJS | 连接分组拖拽排序 |

### 项目结构

```text
MQTTW/
├── main.go                     # 应用入口：窗口 / 托盘 / 关闭拦截 / 服务注册
├── go.mod                      # Go module（gofly）
├── Taskfile.yml                # Wails 任务配置
├── build/                      # 构建配置（config.yml 等）
├── resource/db/data.db         # SQLite 数据库
├── internal/
│   ├── dao/                    # GORM Gen 数据访问层（自动生成）
│   ├── service/                # 业务服务（连接 / 订阅 / 消息 / 压测 / 设置等）
│   ├── mqttclient/             # MQTT 客户端封装（老 Paho / autopaho 统一接口）
│   ├── stress/                 # 压测引擎
│   ├── logic/logic_mqtt/       # 消息异步入库等
│   └── utils/                  # 工具函数
├── mcp/                        # MCP 服务（数据库管理）
└── frontend/
    └── src/
        ├── views/              # 页面（connect / stress / setting / help / gmqt）
        ├── store/              # Pinia 状态
        ├── components/         # 公共组件
        ├── router/             # 路由
        └── locale/             # 国际化
```
### 开发启动
```
wails3 dev
```
### 构建和包装
```
wails3 package
```
#### Windows 封装 
使用NSIS安装器
```
wails3 package GOOS=windows INSTALL_SCOPE=user  
```
MSIX 软件包，对于 Microsoft Store 发行版或现代 Windows 部署：
```
wails3 package GOOS=windows FORMAT=msix
```
#### macOS 封装
```
wails3 package GOOS=darwin
```
####  Linux 封装
```
wails3 build GOOS=linux GOARCH=arm64
wails3 package GOOS=linux GOARCH=arm64

```
---

## 学习资料

### MQTT 协议

- [MQTT v5.0 规范（OASIS）](https://docs.oasis-open.org/mqtt/mqtt/v5.0/mqtt-v5.0.html)
- [MQTT v3.1.1 规范（OASIS）](https://docs.oasis-open.org/mqtt/mqtt/v3.1.1/os/mqtt-v3.1.1-os.html)

### Go 与桌面应用

- [Go 官方文档](https://go.dev/doc/)
- [Go by Example](https://gobyexample.com/)
- [Wails 官方文档](https://wails.io/docs/introduction)
- [Wails3 文档](https://v3.wails.io/)

### 前端

- [Vue 3 官方文档](https://vuejs.org/)
- [Arco Design Vue 组件库](https://arco.design/vue/)
- [Vite 文档](https://vitejs.dev/)
- [Pinia 状态管理](https://pinia.vuejs.org/)
- [ECharts 文档](https://echarts.apache.org/zh/index.html)

### MQTT 客户端库

- [eclipse/paho.mqtt.golang（MQTT 3.1/3.1.1）](https://github.com/eclipse/paho.mqtt.golang)
- [eclipse/paho.golang（MQTT 5.0）](https://github.com/eclipse/paho.golang)
- [gonzalop/mq（压测轻量客户端）](https://github.com/gonzalop/mq)

### 社区与配套

- [GoFly 全栈开发社区](https://goflys.cn/)
- [GMQT - MQTT Broker（GoFlyGen 框架集成）](https://goflys.cn/gmqt)

---

## 版权

© 2026 昆明立师科技有限公司 · 保留所有权利
