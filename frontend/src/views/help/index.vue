<template>
  <div class="container">
  <div class="help-page">
    <!-- Hero -->
    <div class="hero">
      <div class="hero-bg"></div>
      <div class="hero-content">
        <div class="hero-logo">MQTTW</div>
        <h1 class="hero-title">MQTT 调试 · 压测一体化桌面工具</h1>
        <p class="hero-desc">
          基于 Wails3 + Go 打造的专业 MQTT 客户端，全面支持 MQTT 3.1 / 3.1.1 / 5.0
          协议，集连接调试、消息收发、多连接管理与大规模压测于一体。
        </p>
        <div class="hero-tags">
          <a-tag color="arcoblue" size="large">MQTT 3.1</a-tag>
          <a-tag color="arcoblue" size="large">MQTT 3.1.1</a-tag>
          <a-tag color="arcoblue" size="large">MQTT 5.0</a-tag>
          <a-tag color="green" size="large">多连接并行</a-tag>
          <a-tag color="green" size="large">实时指标</a-tag>
          <a-tag color="orange" size="large">独立压测</a-tag>
        </div>
      </div>
    </div>

    <!-- 功能模块 -->
    <div class="section">
      <h2 class="section-title">核心功能</h2>
      <div class="feature-grid">
        <a-card v-for="f in features" :key="f.title" :bordered="false" class="feature-card">
          <div class="feature-icon" :style="{ background: f.bg, color: f.color }">
            <icon-font :name="f.icon" :size="26" />
          </div>
          <h3 class="feature-title">{{ f.title }}</h3>
          <p class="feature-desc">{{ f.desc }}</p>
          <ul class="feature-points">
            <li v-for="p in f.points" :key="p">{{ p }}</li>
          </ul>
        </a-card>
      </div>
    </div>

    <!-- MQTT 5.0 特性 -->
    <div class="section alt-section">
      <h2 class="section-title">MQTT 5.0 完整支持</h2>
      <p class="section-desc">
        不仅兼容 3.1 / 3.1.1，MQTTW 完整实现了 MQTT 5.0 新增特性，满足生产级调试需求。
      </p>
      <div class="v5-grid">
        <div v-for="item in v5Features" :key="item.title" class="v5-item">
          <div class="v5-dot"></div>
          <div class="v5-text">
            <div class="v5-title">{{ item.title }}</div>
            <div class="v5-desc">{{ item.desc }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 压测亮点 -->
    <div class="section">
      <h2 class="section-title">独立压测模块</h2>
      <p class="section-desc">
        脱离调试链路独立运行的大规模客户端压测引擎，指标实时可视，结果一目了然。
      </p>
      <div class="stress-stats">
        <div v-for="s in stressStats" :key="s.label" class="stress-stat">
          <div class="stat-value">{{ s.value }}</div>
          <div class="stat-label">{{ s.label }}</div>
        </div>
      </div>
      <div class="stress-points">
        <div v-for="p in stressPoints" :key="p" class="stress-point">
          <icon-check-circle-fill class="point-icon" />
          {{ p }}
        </div>
      </div>
    </div>

    <!-- 技术架构 -->
    <div class="section alt-section">
      <h2 class="section-title">技术架构</h2>
      <div class="arch-row">
        <div class="arch-box">
          <div class="arch-layer" style="--layer: 1">
            <div class="arch-name">前端</div>
            <div class="arch-tech">Vue 3 · Arco Design · ECharts · Pinia · Vite</div>
          </div>
          <div class="arch-arrow">▲</div>
          <div class="arch-layer" style="--layer: 2">
            <div class="arch-name">通信桥</div>
            <div class="arch-tech">Wails3 Runtime · 事件推送（Event Emit / On）</div>
          </div>
          <div class="arch-arrow">▲</div>
          <div class="arch-layer" style="--layer: 3">
            <div class="arch-name">后端</div>
            <div class="arch-tech">Go · Paho（3.1/3.1.1）· autopaho（5.0）· gonzalop/mq（压测）</div>
          </div>
          <div class="arch-arrow">▲</div>
          <div class="arch-layer" style="--layer: 4">
            <div class="arch-name">数据层</div>
            <div class="arch-tech">SQLite · GORM Gen · 连接/订阅/消息持久化</div>
          </div>
        </div>
      </div>
    </div>

    <!-- GoFly 社区 -->
    <div class="section community-section">
      <h2 class="section-title">GoFly 全栈开发社区</h2>
      <div class="community-card">
        <div class="community-logo">
          <div class="logo-text">GoFly</div>
          <div class="logo-sub">goflys.cn</div>
        </div>
        <div class="community-info">
          <div class="community-name">MQTTW 由 GoFly 技术团队打造</div>
          <p class="community-desc">
            GoFly 全栈开发社区（goflys.cn）汇聚全栈开发者，围绕 Go 后端、
            Vue 前端、Wails 桌面应用、AI Agent 智能体开发、AI 编程框架，
            以及物联网（IoT）开发经验与 MQTT 实战等技术方向，提供开源项目、
            实战教程与技术交流。加入社区，与开发者一起成长。
          </p>
          <div class="community-tags">
            <a-tag color="arcoblue" size="medium">开源项目</a-tag>
            <a-tag color="green" size="medium">实战教程</a-tag>
            <a-tag color="orange" size="medium">技术问答</a-tag>
            <a-tag color="purple" size="medium">全栈交流</a-tag>
            <a-tag color="cyan" size="medium">AI Agent 开发</a-tag>
            <a-tag color="gold" size="medium">AI 编程框架</a-tag>
            <a-tag color="lime" size="medium">物联网开发</a-tag>
          </div>
          <a-button type="primary" size="large" class="community-btn" @click="openCommunity">
            <icon-link />
            <span>访问 GoFly 社区</span>
          </a-button>
        </div>
      </div>
    </div>

    <!-- 版权 -->
    <div class="copyright">
      © 2026 昆明立师科技有限公司 · 保留所有权利
    </div>
  </div>
  </div>
</template>

<script lang="ts" setup>
import { Browser } from '@wailsio/runtime';

// 打开 GoFly 社区官网
const openCommunity = () => {
  Browser.OpenURL('https://goflys.cn/');
};

interface Feature {
  title: string;
  desc: string;
  icon: string;
  color: string;
  bg: string;
  points: string[];
}

const features: Feature[] = [
  {
    title: '多协议连接',
    desc: '一套界面兼容三代 MQTT 协议，自动识别协议版本，无需切换工具。',
    icon: 'icon-duoxieyizhichi',
    color: '#165DFF',
    bg: 'rgba(22, 93, 255, 0.1)',
    points: [
      'MQTT 3.1 / 3.1.1（Paho 老引擎）',
      'MQTT 5.0（Paho V5 / autopaho）',
      '协议自动协商（Auto）',
      'TLS/SSL + CA/证书/密钥/ALPN',
    ],
  },
  {
    title: '连接管理',
    desc: '分组树组织全部连接，多连接并行在线，连接参数一目了然。',
    icon: 'icon-connect5',
    color: '#00B42A',
    bg: 'rgba(0, 180, 42, 0.1)',
    points: [
      '分组 / 连接树形管理，拖拽排序',
      '多连接同时在线，独立消息流',
      '自动重连 + 可配置重连周期',
      'Clean Session / Clean Start 会话控制',
    ],
  },
  {
    title: '订阅与消息',
    desc: '订阅管理持久化，消息实时推送，JSON 智能高亮，未读角标提醒。',
    icon: 'icon-messages',
    color: '#FF7D00',
    bg: 'rgba(255, 125, 0, 0.1)',
    points: [
      '订阅持久化，重复订阅自动拦截',
      '消息实时推送，连接级未读计数',
      'JSON 消息自动检测并高亮展示',
      '消息入库 SQLite，可回溯',
    ],
  },
  {
    title: '消息发布',
    desc: '灵活的消息发布面板，QoS / Retain / MQTT 5 发布属性全覆盖。',
    icon: 'icon-send1',
    color: '#F53F3F',
    bg: 'rgba(245, 63, 63, 0.1)',
    points: [
      'QoS 0 / 1 / 2 + Retain 保留消息',
      '主题别名、Payload 格式指示',
      '订阅标识符、用户属性自定义',
      '响应主题 / Content Type / 消息过期',
    ],
  },
  {
    title: '遗嘱消息',
    desc: '完整的 Will 配置，异常掉线时按预期广播离线状态。',
    icon: 'icon-acewill-leaveback',
    color: '#722ED1',
    bg: 'rgba(114, 46, 209, 0.1)',
    points: [
      '遗嘱主题 / QoS / Retain / Payload',
      'Will Delay Interval 延迟发布',
      'Payload 格式指示 / Content Type',
      '响应主题与关联数据',
    ],
  },
  {
    title: '独立压测',
    desc: '脱离调试链路的压测引擎，数千到数万客户端同时在线实测。',
    icon: 'icon-stress',
    color: '#0FC6C2',
    bg: 'rgba(15, 198, 194, 0.1)',
    points: [
      '批量客户端 + 连接速率控制',
      '订阅 / 发布压测，速率与总量可控',
      '回显延迟 P50 / P95 / P99 分位',
      '实时曲线 + 结果汇总 + CSV 导出',
    ],
  },
];

const v5Features = [
  { title: '用户属性（User Properties）', desc: '自定义键值对，携带业务上下文' },
  { title: '订阅标识符（Subscription Identifier）', desc: '区分同一主题不同订阅来源，范围 1 ~ 268435455' },
  { title: '主题别名（Topic Alias）', desc: '短报文代替长主题，显著降低带宽占用' },
  { title: '会话过期（Session Expiry Interval）', desc: '连接断开后会话保留时长自定义' },
  { title: '请求 / 响应（Request / Response）', desc: '响应主题 + 关联数据，实现请求应答模式' },
  { title: '问题信息（Problem Information）', desc: '服务器返回更详细的错误诊断信息' },
  { title: '接收最大值 / 最大数据包', desc: '流控参数透明可配，贴近服务器能力' },
  { title: '保留消息与消息过期', desc: '消息生命周期精细化控制' },
];

const stressStats = [
  { value: '3.1/3.1.1/5.0', label: '压测协议' },
  { value: '数千~数万', label: '并发客户端' },
  { value: '实时 500ms', label: '指标推送' },
  { value: 'P50/P95/P99', label: '延迟分位' },
];

const stressPoints = [
  '连接进度实时可见，运行状态与时长随时掌握',
  '连接成功率、发布/接收速率、回显延迟多维度监控',
  '失败原因自动归类，Top 结果一目了然',
  '配置持久化本地存储，下次启动一键复用',
];
</script>

<style scoped lang="less">
.help-page {
  max-width: 1080px;
  margin: 0 auto;
  padding: 0 20px 40px;
  color: var(--color-text-1);
}

/* ---------- Hero ---------- */
.hero {
  position: relative;
  border-radius: 12px;
  overflow: hidden;
  margin-top: 20px;
  padding: 56px 40px;
  text-align: center;
  background: linear-gradient(135deg, #165DFF 0%, #0FC6C2 100%);
  color: #fff;

  .hero-content {
    position: relative;
    z-index: 1;
  }

  .hero-logo {
    display: inline-block;
    font-size: 42px;
    font-weight: 800;
    letter-spacing: 4px;
    padding: 4px 24px;
    border: 2px solid rgba(255, 255, 255, 0.7);
    border-radius: 10px;
    background: rgba(255, 255, 255, 0.12);
    backdrop-filter: blur(4px);
  }

  .hero-title {
    margin: 22px 0 12px;
    font-size: 28px;
    font-weight: 600;
  }

  .hero-desc {
    max-width: 640px;
    margin: 0 auto 24px;
    font-size: 15px;
    line-height: 1.8;
    color: rgba(255, 255, 255, 0.92);
  }

  :deep(.arco-tag) {
    background: rgba(255, 255, 255, 0.16);
    border-color: rgba(255, 255, 255, 0.4);
    color: #fff;
    margin: 0 4px 8px;
  }
}

/* ---------- 通用 section ---------- */
.section {
  margin-top: 44px;

  .section-title {
    font-size: 22px;
    font-weight: 600;
    margin-bottom: 8px;
  }

  .section-desc {
    font-size: 14px;
    color: var(--color-text-3);
    margin: 0 0 20px;
    line-height: 1.7;
  }
}

.alt-section {
  margin-top: 44px;
  padding: 28px 32px;
  border-radius: 12px;
  background: var(--color-fill-2);
}

/* ---------- 功能卡片 ---------- */
.feature-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;

  @media (max-width: 900px) {
    grid-template-columns: repeat(2, 1fr);
  }
  @media (max-width: 640px) {
    grid-template-columns: 1fr;
  }
}

.feature-card {
  border-radius: 10px;
  background: var(--color-bg-2);
  transition: transform 0.2s ease, box-shadow 0.2s ease;

  &:hover {
    transform: translateY(-3px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  }

  .feature-icon {
    width: 52px;
    height: 52px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .feature-title {
    margin: 16px 0 8px;
    font-size: 17px;
    font-weight: 600;
  }

  .feature-desc {
    font-size: 13px;
    line-height: 1.7;
    color: var(--color-text-3);
    margin-bottom: 10px;
  }

  .feature-points {
    margin: 0;
    padding-left: 18px;
    font-size: 12.5px;
    line-height: 2;
    color: var(--color-text-2);
  }
}

/* ---------- MQTT5 ---------- */
.v5-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px 32px;

  @media (max-width: 640px) {
    grid-template-columns: 1fr;
  }
}

.v5-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 10px 0;

  .v5-dot {
    flex: none;
    width: 8px;
    height: 8px;
    border-radius: 50%;
    margin-top: 7px;
    background: linear-gradient(135deg, #165DFF, #0FC6C2);
  }

  .v5-title {
    font-size: 14px;
    font-weight: 600;
  }

  .v5-desc {
    font-size: 12.5px;
    color: var(--color-text-3);
    margin-top: 2px;
  }
}

/* ---------- 压测 ---------- */
.stress-stats {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 14px;
  margin-bottom: 24px;

  @media (max-width: 640px) {
    grid-template-columns: repeat(2, 1fr);
  }

  .stress-stat {
    padding: 20px 12px;
    text-align: center;
    border-radius: 10px;
    background: var(--color-bg-2);
    border: 1px solid var(--color-border-2);

    .stat-value {
      font-size: 22px;
      font-weight: 700;
      background: linear-gradient(135deg, #165DFF, #0FC6C2);
      -webkit-background-clip: text;
      background-clip: text;
      color: transparent;
    }

    .stat-label {
      margin-top: 6px;
      font-size: 13px;
      color: var(--color-text-3);
    }
  }
}

.stress-points {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;

  @media (max-width: 640px) {
    grid-template-columns: 1fr;
  }

  .stress-point {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 13.5px;
    color: var(--color-text-2);

    .point-icon {
      color: #00B42A;
      font-size: 16px;
    }
  }
}

/* ---------- 架构 ---------- */
.arch-row {
  display: flex;
  justify-content: center;
}

.arch-box {
  width: 100%;
  max-width: 620px;
}

.arch-layer {
  padding: 16px 20px;
  border-radius: 10px;
  background: var(--color-bg-2);
  border: 1px solid var(--color-border-2);
  text-align: center;

  .arch-name {
    font-size: 15px;
    font-weight: 600;
  }

  .arch-tech {
    margin-top: 4px;
    font-size: 12.5px;
    color: var(--color-text-3);
  }
}

.arch-arrow {
  text-align: center;
  color: var(--color-text-4);
  line-height: 1;
  padding: 8px 0;
  font-size: 12px;
}

/* ---------- GoFly 社区 ---------- */
.community-card {
  display: flex;
  gap: 28px;
  align-items: center;
  padding: 28px 32px;
  border-radius: 12px;
  background: linear-gradient(135deg, #1D2129 0%, #23252E 100%);
  color: #fff;

  @media (max-width: 640px) {
    flex-direction: column;
    text-align: center;
  }

  .community-logo {
    flex: none;
    width: 128px;
    height: 128px;
    border-radius: 20px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 6px;
    background: linear-gradient(135deg, #165DFF, #0FC6C2);

    .logo-text {
      font-size: 30px;
      font-weight: 800;
      letter-spacing: 1px;
    }

    .logo-sub {
      font-size: 13px;
      opacity: 0.85;
    }
  }

  .community-info {
    flex: 1;
    min-width: 0;

    .community-name {
      font-size: 18px;
      font-weight: 600;
    }

    .community-desc {
      margin: 10px 0 14px;
      font-size: 13.5px;
      line-height: 1.8;
      color: rgba(255, 255, 255, 0.78);
    }

    :deep(.arco-tag) {
      background: rgba(255, 255, 255, 0.12);
      border-color: rgba(255, 255, 255, 0.28);
      color: #fff;
      margin-right: 8px;
      margin-bottom: 8px;
    }

    .community-btn {
      margin-top: 6px;
      display: inline-flex;
      align-items: center;
      gap: 8px;
    }
  }
}

/* ---------- 版权 ---------- */
.copyright {
  margin-top: 48px;
  padding-top: 24px;
  border-top: 1px solid var(--color-border-2);
  text-align: center;
  font-size: 13px;
  color: var(--color-text-3);
}
</style>
