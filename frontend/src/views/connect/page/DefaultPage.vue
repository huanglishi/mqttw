<template>
  <div class="container">
    <div class="page flex-all-center">
      <!-- Hero：主视觉区 -->
      <div class="hero">
        <div class="hero-icon">
          <svg width="150" viewBox="0 0 240 180" xmlns="http://www.w3.org/2000/svg">
            <path d="M120 60 A60 60 0 0 1 180 120" stroke="#ff992b" stroke-width="14" fill="none" stroke-linecap="round">
              <animate attributeName="stroke-width" values="14;18;14" dur="1.8s" repeatCount="indefinite" />
            </path>
            <path d="M120 60 A60 60 0 0 0 60 120" stroke="#ff992b" stroke-width="14" fill="none" stroke-linecap="round">
              <animate attributeName="stroke-width" values="14;18;14" dur="1.8s" repeatCount="indefinite" />
            </path>
            <path d="M120 80 A40 40 0 0 1 160 120" stroke="#ff992b" stroke-width="14" fill="none" stroke-linecap="round">
              <animate attributeName="stroke-width" values="14;20;14" dur="1.4s" repeatCount="indefinite" />
            </path>
            <path d="M120 80 A40 40 0 0 0 80 120" stroke="#ff992b" stroke-width="14" fill="none" stroke-linecap="round">
              <animate attributeName="stroke-width" values="14;20;14" dur="1.4s" repeatCount="indefinite" />
            </path>
            <circle cx="120" cy="124" r="16" fill="#ff7722"/>
            <path d="M120 124 L142 170 L98 170 Z" fill="#ff7722"/>
          </svg>
        </div>
        <h2 class="hero-title">欢迎使用 <span>MQTTW</span></h2>
        <p class="hero-desc">创建你的第一个 MQTT 连接，或从这里快速获取 Broker 与学习资源</p>
        <div class="hero-actions">
          <a-button type="primary" size="large" @click="handleAddConnect">
            <template #icon><icon-plus /></template>
            <template #default>新建连接</template>
          </a-button>
        </div>
      </div>

      <!-- 引导卡片 -->
      <div class="guide-grid">
        <div class="guide-card">
          <div class="guide-head">
            <span class="guide-num">01</span>
            <span class="guide-icon"><icon-download /></span>
          </div>
          <div class="guide-title">需要私有或在线 MQTT 服务？</div>
          <div class="guide-text">
            免费的 MQTT-Broker 一键私有部署，在 Windows、Mac、Linux 等操作系统 1 分钟快速部署；在线 MQTT 服务提供试用、学习、测试等场景使用，前往
            <a-link @click="openWeb">GMQT-Broker</a-link> 获取。
          </div>
          <div class="guide-action">
            <a-button size="small" @click="openWeb">
              <template #icon><icon-link /></template>
              了解 GMQT-Broker
            </a-button>
          </div>
        </div>

        <div class="guide-card">
          <div class="guide-head">
            <span class="guide-num">02</span>
            <span class="guide-icon"><icon-cloud /></span>
          </div>
          <div class="guide-title">获取在线免费 MQTT 服务器</div>
          <div class="guide-text">
            GoFly 社区免费提供在线 GMQT-Broker 服务，地址：
            <a-link @click="openWebURL('https://gmqt.goflys.cn/')">https://gmqt.goflys.cn/</a-link>
            。默认账号
            <code @click="handleCopyTopic('admin')">admin</code>
            ，密码
            <code @click="handleCopyTopic('admin')">admin</code>
            ，支持 TCP、WS、SSL、WSS 四种服务协议。
          </div>
          <div class="guide-action">
            <a-button size="small" @click="openWebURL('https://gmqt.goflys.cn/')">
              <template #icon><icon-launch /></template>
              在线服务
            </a-button>
            <a-button size="small" @click="handleCopyTopic('admin')">
              <template #icon><icon-copy /></template>
              复制账号
            </a-button>
          </div>
        </div>

        <div class="guide-card">
          <div class="guide-head">
            <span class="guide-num">03</span>
            <span class="guide-icon"><icon-book /></span>
          </div>
          <div class="guide-title">有完整、系统的学习教程吗？</div>
          <div class="guide-text">
            GoFly 社区免费提供
            <a-link @click="openWebURL('https://doc.goflys.cn/docview?id=50&fid=1181')">视频教程和使用文档</a-link>
            ，从协议基础到实战应用，帮助你系统掌握 MQTT。
          </div>
          <div class="guide-action">
            <a-button size="small" @click="openWebURL('https://doc.goflys.cn/docview?id=50&fid=1181')">
              <template #icon><icon-book /></template>
              查看教程文档
            </a-button>
          </div>
        </div>
      </div>

      <div class="foot-tip">提示：首次使用建议从「新建连接」开始，连接成功后即可订阅与发布消息</div>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import {Browser,Clipboard} from "@wailsio/runtime";
  import { Message} from '@arco-design/web-vue';
  const openWeb=()=>{
    Browser.OpenURL("https://goflys.cn/gmqt")
  }
  //在默认浏览器打开网站
  const openWebURL=(url:string)=>{
    Browser.OpenURL(url)
  }
   //复制内容
  const handleCopyTopic=async(text:string)=>{
    await Clipboard.SetText(text)
    Message.success("复制成功")
  }
  const emit = defineEmits<{
    'add-connect': [index: number]
  }>()
  const handleAddConnect=()=>{
    emit('add-connect', 2)
  }
</script>

<style scoped lang="less">
  .page {
    max-width: 1020px;
    margin: 0 auto;
    flex-direction: column;
    height: 100vh;
  }

  .hero {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    padding: 28px 0 20px;
     margin-top: 20px;
    .hero-icon {
      filter: drop-shadow(0 10px 18px rgba(255, 119, 34, 0.22));
      animation: float 3.2s ease-in-out infinite;
    }

    .hero-title {
      margin: 14px 0 0;
      font-size: 26px;
      font-weight: 700;
      color: var(--color-text-1);

      span {
        background: linear-gradient(90deg, #ff7722, #ff992b, #ffc53d);
        -webkit-background-clip: text;
        background-clip: text;
        -webkit-text-fill-color: transparent;
      }
    }

    .hero-desc {
      margin: 8px 0 0;
      font-size: 14px;
      color: var(--color-text-3);
    }

    .hero-actions {
      margin-top: 20px;
    }
  }

  // ===== 引导卡片 =====
  .guide-grid {
    width: 100%;
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 16px;
    margin-top: 8px;
    padding: 0px 10px;
    @media (max-width: 900px) {
      grid-template-columns: 1fr;
    }
  }

  .guide-card {
    display: flex;
    flex-direction: column;
    padding: 20px 18px;
    border-radius: 12px;
    background: var(--color-bg-2);
    border: 1px solid var(--color-border-2);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
    transition: box-shadow 0.2s ease, transform 0.2s ease, border-color 0.2s ease;

    &:hover {
      transform: translateY(-3px);
      border-color: rgb(var(--arcoblue-4));
      box-shadow: 0 12px 24px rgba(0, 0, 0, 0.08);
    }

    .guide-head {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 12px;

      .guide-num {
        font-size: 24px;
        font-weight: 700;
        line-height: 1;
        background: linear-gradient(135deg, #ff7722, #ffb400);
        -webkit-background-clip: text;
        background-clip: text;
        -webkit-text-fill-color: transparent;
        color: #ff7722;
      }

      .guide-icon {
        width: 36px;
        height: 36px;
        display: inline-flex;
        align-items: center;
        justify-content: center;
        border-radius: 10px;
        font-size: 18px;
        color: #ff7722;
        background: rgba(255, 119, 34, 0.12);
      }
    }

    .guide-title {
      font-size: 15px;
      font-weight: 600;
      color: var(--color-text-1);
      margin-bottom: 8px;
    }

    .guide-text {
      flex: 1;
      font-size: 13px;
      line-height: 1.8;
      color: var(--color-text-3);

      code {
        background: var(--color-neutral-2);
        border-radius: 4px;
        font-size: 13px;
        padding: 0 4px 1px 4px;
        cursor: pointer;
        transition: background 0.2s ease;
        user-select: all;

        &:hover {
          background: rgb(var(--arcoblue-1));
          color: rgb(var(--arcoblue-6));
        }
      }
    }

    .guide-action {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
      margin-top: 14px;

      .arco-btn {
        display: inline-flex;
        align-items: center;
      }
    }
  }

  // ===== 底部提示 =====
  .foot-tip {
    margin-top: 22px;
    padding: 8px 16px;
    border-radius: 8px;
    font-size: 12.5px;
    color: var(--color-text-4);
    background: var(--color-fill-1);
  }

  @keyframes float {
    0%, 100% { transform: translateY(0); }
    50% { transform: translateY(-7px); }
  }
</style>
