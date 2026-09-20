<div align="center">
  <a href="https://goflys.cn/mqttw" target="_blank">
    <img alt="MQTTW Logo" width="200" src="https://doc.goflys.cn/image/mqttwicon.png"/>
  </a>
</div>

<div align="center">
  <h1>MQTTW</h1>
</div>

<div align="center">

 All-in-one MQTT debugging & stress-testing desktop client

</div>

 <div align="center">

English | [简体中文](./README.md)

</div>

**MQTTW** is a professional MQTT desktop debugging tool built on **Wails3 + Go + Vue3**, with full support for the **MQTT 3.1 / 3.1.1 / 5.0** protocols. It combines connection debugging, message publishing/subscribing, multi-connection management, and large-scale client stress testing in one tool — ideal for IoT development, Broker integration verification, and performance testing scenarios.

MQTTW also serves as a **classic Wails3 desktop application example** — with a clean project structure and complete frontend/backend interaction. It works both as a reference project for developers learning Wails3, and as a solid base for custom MQTT tooling or secondary development.

![MQTT Client Tool](https://docapi.goflys.cn/common/uploadfile/get_image?url=resource/uploads/20260919/ca9f5c4964706e1353fb424904b1dbca.png)

## Features

### 1. Multi-Protocol Connections

- One UI compatible with three generations of MQTT protocols, with automatic version detection
- **MQTT 3.1 / 3.1.1**: `github.com/eclipse/paho.mqtt.golang` (legacy Paho engine)
- **MQTT 5.0**: `github.com/eclipse/paho.golang/autopaho` (new Paho V5, with automatic subscription recovery on reconnect)
- TLS/SSL connections with CA certificate, client certificate, private key, and ALPN configuration
- Clean Session / Clean Start session control
- Auto-reconnect with configurable reconnect period
- Configurable connection timeout and Keep Alive interval

### 2. Full MQTT 5.0 Feature Support

- **User Properties**: custom key-value pairs
- **Subscription Identifier**: range 1 ~ 268435455
- **Topic Alias**: short packets instead of long topics, reducing bandwidth usage
- **Session Expiry Interval**: custom session retention after disconnect
- **Request / Response**: response topic + correlation data for request-reply patterns
- **Problem Information** and **Receive Maximum / Maximum Packet Size** flow-control parameters
- **Will Message**: full configuration for will topic / QoS / Retain / Payload / Will Delay Interval / Content Type

### 3. Connection Management

- Group tree management (group / connection hierarchy) with drag-and-drop sorting
- Multiple simultaneous online connections with independent message flows
- Per-connection unread counters + group-level summary badges
- Connection configurations persisted to SQLite for one-click reuse

### 4. Subscriptions & Messages

- Persistent subscription management with duplicate subscription interception
- Real-time message push with event-driven UI updates
- Automatic JSON detection and syntax highlighting
- Asynchronous message persistence to SQLite (with retry) for later retrieval
- MQTT 5 subscription options: subscription identifier, user properties, etc.

### 5. Message Publishing

- QoS 0 / 1 / 2 + Retain support
- Topic Alias and Payload Format Indicator (UTF-8)
- User properties, response topic, Content Type, message expiry, correlation data
- Complete MQTT 5 publish property panel

### 6. Dedicated Stress-Testing Module

- Fixed **MQTT 3.1** (`github.com/gonzalop/mq` lightweight library) with 3.1.1 / 5.0 options
- Batch creation of thousands to tens of thousands of clients, with controllable connection rate
- Subscription / publish stress testing with controllable rate and total volume
- Real-time metrics: connection success rate, publish/receive rate, publish/receive totals
- **Echo latency percentiles**: P50 / P95 / P99
- Real-time charts (ECharts) + result summary + CSV export
- Live run status and progress; configuration persistence for one-click reuse on next launch

### 7. Settings Panel

- **Startup behavior**: launch at login, auto-restore last online connections on startup, minimize to system tray
- **Message notifications**: system notification + sound toggle on received subscribed messages
- **Global shortcuts**: show/hide window (Ctrl+Shift+H), disconnect all connections (Ctrl+Shift+D), quick publish (Ctrl+Shift+P)
- **Runtime directory**: one-click open the data directory
- **Internationalization**: Simplified Chinese / Traditional Chinese / English
- **About**: version info, update check, GoFly community / GMQT links

### 8. Companion Products

- **GMQT-Broker**: an MQTT server produced by the GoFly community. Can be deployed standalone or installed into the GoFlyGen framework — a ready-to-use companion Broker for MQTTW (<https://goflys.cn/gmqt>)
- **GoFly Full-Stack Development Community**: Go backend, Vue frontend, Wails desktop apps, AI Agent development, AI programming frameworks, IoT development experience, and MQTT hands-on practice (<https://goflys.cn/>)

---

## Getting Started

### Requirements

| Dependency | Version |
| --- | --- |
| Go | 1.25+ |
| Node.js | 18+ (with npm) |
| Wails3 CLI | v3.0.0-beta.9+ |

### Installation & Run

```bash
# 1. Clone the project
git clone https://gitee.com/huang_li_shi_admin/mqttw.git
cd mqttw

# 2. Install dependencies
go mod tidy
cd frontend && npm install && cd ..

# 3. Run in development mode (hot reload for frontend & backend)
wails3 dev

# 4. Production build
wails3 build
```

> The dev-mode frontend runs on port `9245` by default; if the port is occupied, clean up residual processes and re-run `wails3 dev`.

### Basic Usage Flow

1. **Create a connection**: fill in the Broker address, port, and Client ID, choose the MQTT version (3.1 / 3.1.1 / 5.0), and configure Clean Session / auto-reconnect / connection properties
2. **Subscribe to topics**: add subscriptions in the operation panel (QoS, subscription identifier, etc.) and view messages in real time
3. **Publish messages**: fill in the topic and payload, configure QoS / Retain / MQTT 5 publish properties, then send
4. **Stress test**: open the stress-test module, configure client count, connection rate, subscription/publish parameters and duration, then watch the live metric charts
5. **Data storage**: connections, subscriptions, and message records are persisted to `resource/db/data.db` (SQLite)

---

## Technology Stack

### Backend (Go)

| Component | Description |
| --- | --- |
| [Wails3](https://wails.io/) | v3.0.0-beta.9, desktop application framework (window / tray / global shortcuts / events) |
| [paho.mqtt.golang](https://github.com/eclipse/paho.mqtt.golang) | v1.5.1, MQTT 3.1 / 3.1.1 client |
| [paho.golang](https://github.com/eclipse/paho.golang) | v0.23.0 (autopaho), MQTT 5.0 client |
| [gonzalop/mq](https://github.com/gonzalop/mq) | v0.9.10, lightweight client for the stress-test module (MQTT 3.1) |
| [GORM](https://gorm.io/) + GORM Gen | ORM and DAO code generation |
| [modernc.org/sqlite](https://modernc.org/sqlite) | Pure-Go SQLite driver, no CGO |
| [mcp-go](https://github.com/mark3labs/mcp-go) | MCP service, gofly-remote-mcp manages the database |

### Frontend (Vue3)

| Component | Description |
| --- | --- |
| Vue 3 + TypeScript | Composition API |
| Vite | Build tool |
| Arco Design Vue | Component library (incl. icon-font icons) |
| Pinia | State management (connection / message / stress stores) |
| Vue Router | Routing (connect / stress / setting / help / gmqt) |
| ECharts + vue-echarts | Real-time stress-test charts and metric visualization |
| vue-i18n | Internationalization (Simplified Chinese / Traditional Chinese / English) |
| Monaco Editor + PrismJS | Message / JSON syntax highlighting |
| SortableJS | Drag-and-drop sorting for connection groups |

### Project Structure

```text
MQTTW/
├── main.go                     # App entry: window / tray / close interception / service registration
├── go.mod                      # Go module (gofly)
├── Taskfile.yml                # Wails task configuration
├── build/                      # Build configuration (config.yml, etc.)
├── resource/db/data.db         # SQLite database
├── internal/
│   ├── dao/                    # GORM Gen data access layer (auto-generated)
│   ├── service/                # Business services (connection / subscription / message / stress / settings, etc.)
│   ├── mqttclient/             # MQTT client wrappers (legacy Paho / autopaho unified interface)
│   ├── stress/                 # Stress-test engine
│   ├── logic/logic_mqtt/       # Async message persistence, etc.
│   └── utils/                  # Utility functions
├── mcp/                        # MCP service (database management)
└── frontend/
    └── src/
        ├── views/              # Pages (connect / stress / setting / help / gmqt)
        ├── store/              # Pinia state
        ├── components/         # Shared components
        ├── router/             # Routing
        └── locale/             # Internationalization
```

### Development

```
wails3 dev
```

### Build & Package

```
wails3 package
```

#### Windows Packaging

Uses the NSIS installer:

```
wails3 package GOOS=windows INSTALL_SCOPE=user
```

MSIX package, for Microsoft Store distribution or modern Windows deployment:

```
wails3 package GOOS=windows FORMAT=msix
```

#### macOS Packaging

```
wails3 package GOOS=darwin
```

#### Linux Packaging

```
wails3 build GOOS=linux GOARCH=arm64
wails3 package GOOS=linux GOARCH=arm64
```

---

## Learning Resources

### MQTT Protocol

- [MQTT v5.0 Specification (OASIS)](https://docs.oasis-open.org/mqtt/mqtt/v5.0/mqtt-v5.0.html)
- [MQTT v3.1.1 Specification (OASIS)](https://docs.oasis-open.org/mqtt/mqtt/v3.1.1/os/mqtt-v3.1.1-os.html)

### Go & Desktop Apps

- [Go Documentation](https://go.dev/doc/)
- [Go by Example](https://gobyexample.com/)
- [Wails Documentation](https://wails.io/docs/introduction)
- [Wails3 Documentation](https://v3.wails.io/)

### Frontend

- [Vue 3 Documentation](https://vuejs.org/)
- [Arco Design Vue](https://arco.design/vue/)
- [Vite Documentation](https://vitejs.dev/)
- [Pinia](https://pinia.vuejs.org/)
- [ECharts Documentation](https://echarts.apache.org/zh/index.html)

### MQTT Client Libraries

- [eclipse/paho.mqtt.golang (MQTT 3.1/3.1.1)](https://github.com/eclipse/paho.mqtt.golang)
- [eclipse/paho.golang (MQTT 5.0)](https://github.com/eclipse/paho.golang)
- [gonzalop/mq (lightweight stress-test client)](https://github.com/gonzalop/mq)

### Community & Companion

- [GoFly Full-Stack Development Community](https://goflys.cn/)
- [GMQT - MQTT Broker (GoFlyGen framework integration)](https://goflys.cn/gmqt)

---

## Copyright

© 2026 Kunming Lishi Technology Co., Ltd. · All rights reserved
