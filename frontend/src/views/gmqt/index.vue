<template>
  <div class="container">
  <div class="gmqt-page">
    <!-- Hero -->
    <div class="hero">
      <div class="hero-content">
        <div class="hero-logo">GMQT</div>
        <h1 class="hero-title">MQTT Broker · 物联网服务器</h1>
        <p class="hero-desc">
          GMQT 是 GoFly 社区打造的 MQTT 服务器端产品，支持独立部署，
          也可一键安装到 GoFlyGen 框架。作为 MQTTW 的配套 Broker，
          让开发者即装即用，快速完成连接、订阅、发布的联调验证与压测。
        </p>
        <div class="hero-tags">
          <a-tag size="large">独立部署</a-tag>
          <a-tag size="large">GoFlyGen 框架集成</a-tag>
          <a-tag size="large">标准 MQTT 协议</a-tag>
          <a-tag size="large">MQTTW 配套</a-tag>
        </div>
        <a-button type="primary" size="large" class="hero-btn" @click="openDoc">
          <icon-global />
          <span>访问 GMQT 社区介绍</span>
        </a-button>
      </div>
    </div>

    <!-- 为什么推荐 -->
    <div class="section">
      <h2 class="section-title">为什么推荐给 MQTTW 用户</h2>
      <p class="section-desc">
        MQTTW 是调试客户端，GMQT 是服务端——两者同属 GoFly 技术团队出品，天然配套。
      </p>
      <div class="reason-grid">
        <a-card v-for="r in reasons" :key="r.title" :bordered="false" class="reason-card">
          <div class="reason-icon" :style="{ background: r.bg, color: r.color }">
            <component :is="iconMap[r.icon]" :size="24" />
          </div>
          <h3 class="reason-title">{{ r.title }}</h3>
          <p class="reason-desc">{{ r.desc }}</p>
        </a-card>
      </div>
    </div>

    <!-- 部署形态 -->
    <div class="section alt-section">
      <h2 class="section-title">两种部署方式，按需选择</h2>
      <div class="deploy-grid">
        <a-card :bordered="false" class="deploy-card">
          <div class="deploy-badge">方式一</div>
          <h3 class="deploy-title">独立部署</h3>
          <p class="deploy-desc">
            GMQT 作为独立 MQTT 服务器运行，不依赖任何框架。
            适合快速搭建本地开发环境、测试环境，或直接作为生产 Broker 使用。
          </p>
          <div class="deploy-tags">
            <a-tag color="arcoblue">即装即用</a-tag>
            <a-tag color="arcoblue">独立进程</a-tag>
            <a-tag color="arcoblue">轻量部署</a-tag>
          </div>
        </a-card>
        <a-card :bordered="false" class="deploy-card">
          <div class="deploy-badge deploy-badge-2">方式二</div>
          <h3 class="deploy-title">GoFlyGen 框架集成</h3>
          <p class="deploy-desc">
            将 GMQT 安装到 GoFlyGen 快速开发框架中，与业务系统共用一套
            运行环境与管理后台，适合已有 GoFlyGen 项目的团队一体化接入。
          </p>
          <div class="deploy-tags">
            <a-tag color="green">框架模块</a-tag>
            <a-tag color="green">统一管理</a-tag>
            <a-tag color="green">业务联动</a-tag>
          </div>
        </a-card>
      </div>
    </div>

    <!-- 核心能力 -->
    <div class="section">
      <h2 class="section-title">核心能力</h2>
      <div class="ability-grid">
        <div v-for="a in abilities" :key="a.title" class="ability-item">
          <icon-check-circle-fill class="ability-icon" />
          <div class="ability-text">
            <div class="ability-title">{{ a.title }}</div>
            <div class="ability-desc">{{ a.desc }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 场景指引 -->
    <div class="section">
      <h2 class="section-title">从 MQTTW 开始，三步用起来</h2>
      <div class="steps">
        <div v-for="(s, i) in steps" :key="s.title" class="step">
          <div class="step-num">{{ i + 1 }}</div>
          <div class="step-title">{{ s.title }}</div>
          <div class="step-desc">{{ s.desc }}</div>
        </div>
      </div>
    </div>

    <!-- CTA -->
    <div class="cta">
      <div class="cta-title">需要一台 MQTT Broker？</div>
      <p class="cta-desc">
        访问 GoFly 社区 GMQT 产品页，查看部署文档、版本说明与社区案例，立即开始使用。
      </p>
      <div class="cta-btns">
        <a-button type="primary" size="large" @click="openDoc">
          <span class="img-middle">
            <icon-font name="icon-MQTT2" :size="13"/>
          </span>
          <span class="img-middle">goflys.cn/gmqt</span>
        </a-button>
        <a-button size="large" @click="openCommunity">
          <span class="img-middle">
            <icon-link :size="13" />
          </span>
          <span class="img-middle">GoFly 社区首页</span>
        </a-button>
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
import { IconUserGroup,IconThunderbolt,IconDownload,IconBook } from '@arco-design/web-vue/es/icon';
const iconMap: any= {
  "icon-user-group": IconUserGroup,
  "icon-thunderbolt": IconThunderbolt,
  "icon-download": IconDownload,
  "icon-book": IconBook,
}
// GMQT 社区产品页
const openDoc = () => {
  Browser.OpenURL('https://goflys.cn/gmqt');
};
// GoFly 社区首页
const openCommunity = () => {
  Browser.OpenURL('https://goflys.cn/');
};

interface ReasonItem {
  title: string;
  desc: string;
  icon: string;
  color: string;
  bg: string;
}

const reasons: ReasonItem[] = [
  {
    title: '同团队配套出品',
    desc: 'GMQT 与 MQTTW 均由 GoFly 技术团队开发，服务端与调试端协议、行为天然对齐，联调少踩坑。',
    icon: 'icon-user-group',
    color: '#165DFF',
    bg: 'rgba(22, 93, 255, 0.1)',
  },
  {
    title: '标准 MQTT 协议',
    desc: '支持标准 MQTT 连接、订阅发布、QoS 与保留消息等机制，与 MQTTW 的全协议调试能力配合使用。',
    icon: 'icon-thunderbolt',
    color: '#00B42A',
    bg: 'rgba(0, 180, 42, 0.1)',
  },
  {
    title: '即装即用',
    desc: '无需从零搭建 Broker，独立部署即可获得完整的 MQTT 服务，几分钟完成本地联调环境。',
    icon: 'icon-download',
    color: '#FF7D00',
    bg: 'rgba(255, 125, 0, 0.1)',
  },
  {
    title: '社区文档与支持',
    desc: 'GoFly 社区提供 GMQT 详细介绍、部署说明与交流渠道，遇到问题有人帮、有文档查。',
    icon: 'icon-book',
    color: '#722ED1',
    bg: 'rgba(114, 46, 209, 0.1)',
  },
];

const abilities = [
  { title: '连接与会话管理', desc: '管理客户端连接、心跳保活与会话生命周期' },
  { title: '主题订阅与发布', desc: '通配符订阅、主题过滤与消息路由' },
  { title: 'QoS 消息质量', desc: '支持 0 / 1 / 2 三级消息投递质量' },
  { title: '保留消息', desc: '新订阅客户端立即获得主题最新保留消息' },
  { title: '遗嘱消息', desc: '客户端异常断开时按预设广播离线状态' },
  { title: '灵活部署', desc: '独立运行或作为 GoFlyGen 框架模块集成' },
];

const steps = [
  { title: '获取 GMQT', desc: '访问社区 GMQT 产品页，按文档下载并部署 Broker' },
  { title: '启动服务', desc: '独立部署或安装到 GoFlyGen 框架，启动 MQTT 服务' },
  { title: 'MQTTW 连接调试', desc: '在 MQTTW 中配置 Broker 地址，开始连接、订阅、发布与压测' },
];
</script>

<style scoped lang="less">
.gmqt-page {
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
  background: linear-gradient(135deg, #0FC6C2 0%, #165DFF 100%);
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
    max-width: 690px;
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

  .hero-btn {
    margin-top: 8px;
    display: inline-flex;
    align-items: center;
    gap: 8px;
    background: #fff;
    color: #165DFF;
    border: none;

    &:hover {
      background: rgba(255, 255, 255, 0.92);
      color: #165DFF;
    }
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

/* ---------- 推荐卡片 ---------- */
.reason-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;

  @media (max-width: 720px) {
    grid-template-columns: 1fr;
  }
}

.reason-card {
  border-radius: 10px;
  background: var(--color-bg-2);
  transition: transform 0.2s ease, box-shadow 0.2s ease;

  &:hover {
    transform: translateY(-3px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  }

  .reason-icon {
    width: 46px;
    height: 46px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .reason-title {
    margin: 14px 0 6px;
    font-size: 16px;
    font-weight: 600;
  }

  .reason-desc {
    font-size: 13px;
    line-height: 1.7;
    color: var(--color-text-3);
    margin: 0;
  }
}

/* ---------- 部署形态 ---------- */
.deploy-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;

  @media (max-width: 720px) {
    grid-template-columns: 1fr;
  }
}

.deploy-card {
  position: relative;
  border-radius: 10px;
  background: var(--color-bg-2);

  .deploy-badge {
    position: absolute;
    top: 0;
    right: 0;
    padding: 4px 14px;
    border-radius: 0 10px 0 10px;
    font-size: 12px;
    font-weight: 600;
    background: rgba(22, 93, 255, 0.12);
    color: #165DFF;
  }

  .deploy-badge-2 {
    background: rgba(0, 180, 42, 0.12);
    color: #00B42A;
  }

  .deploy-title {
    margin: 4px 0 8px;
    font-size: 17px;
    font-weight: 600;
  }

  .deploy-desc {
    font-size: 13px;
    line-height: 1.8;
    color: var(--color-text-3);
    margin-bottom: 12px;
  }
}

/* ---------- 核心能力 ---------- */
.ability-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px 32px;

  @media (max-width: 640px) {
    grid-template-columns: 1fr;
  }
}

.ability-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 10px 0;

  .ability-icon {
    flex: none;
    margin-top: 3px;
    color: #00B42A;
    font-size: 16px;
  }

  .ability-title {
    font-size: 14px;
    font-weight: 600;
  }

  .ability-desc {
    font-size: 12.5px;
    color: var(--color-text-3);
    margin-top: 2px;
  }
}

/* ---------- 步骤 ---------- */
.steps {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;

  @media (max-width: 720px) {
    grid-template-columns: 1fr;
  }
}

.step {
  padding: 20px;
  border-radius: 10px;
  background: var(--color-bg-2);
  border: 1px solid var(--color-border-2);
  text-align: center;

  .step-num {
    width: 36px;
    height: 36px;
    margin: 0 auto 10px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 16px;
    font-weight: 700;
    color: #fff;
    background: linear-gradient(135deg, #165DFF, #0FC6C2);
  }

  .step-title {
    font-size: 15px;
    font-weight: 600;
  }

  .step-desc {
    margin-top: 6px;
    font-size: 12.5px;
    line-height: 1.7;
    color: var(--color-text-3);
  }
}

/* ---------- CTA ---------- */
.cta {
  margin-top: 44px;
  padding: 36px 32px;
  text-align: center;
  border-radius: 12px;
  background: linear-gradient(135deg, #1D2129 0%, #23252E 100%);
  color: #fff;

  .cta-title {
    font-size: 20px;
    font-weight: 600;
  }

  .cta-desc {
    margin: 10px auto 20px;
    max-width: 560px;
    font-size: 13.5px;
    line-height: 1.8;
    color: rgba(255, 255, 255, 0.78);
  }

  .cta-btns {
    display: flex;
    gap: 12px;
    justify-content: center;

    @media (max-width: 480px) {
      flex-direction: column;
    }

    .arco-btn {
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
