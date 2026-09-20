<template>
 <div class="chat-wrap">
  <!-- 对话消息列表 -->
  <template v-for="item in list">
    <div :class="['chat-msg-item flex', item.msg_type === 'received' ? 'chat-received' : 'chat-publish']">
      <div class="avatar" :style="{background:item.msg_type === 'received'?hexToRgba(item.color):'rgb(var(--orange-4))'}">
         <icon-font v-if="item.msg_type === 'received'" name="icon-receive_message" size="22" color="rgb(var(--orange-5))"/>
         <icon-font v-else name="icon-chat-publish" color="#ffffff" size="22"/>
      </div>
      <!--接收消息-->
      <div class="msg-body" v-if="item.msg_type=='received'">
        <div class="msg-wrap" :style="{background:hexToRgba(item.color),color:item.color,'--arrow-color': item.color?hexToRgba(item.color):'#ffffff'}" >
          <div class="chat-option">
            <a-space style="width: 100%;">
              <div class="copy-btn" @click.stop="handleCopyTopic(item.topic)">Topic: {{ item.topic }}</div>
              <div class="">QoS: {{ item.qos }}</div>
              <div class="retain" v-if="item.retain">Retained</div>
            </a-space>
          </div>
          <div class="chat-content">
            <div v-if="isJson(item.payload)" >
              <pre style="margin:0; padding: 4px; background:#f7f8fa; border-radius:4px;">
                <code class="language-json" v-html="highlightJson(item.payload)"></code>
              </pre>
            </div>
            <div v-else class="chat-text">
              {{ item.payload }}
            </div>
          </div>
        </div>
        <div class="message-action-bar">
          <a-space style="width: 100%;">
            <a-tooltip content="复制 Payload">
             <a-link @click="handleCopyTopic(item.payload)" style="color: var(--color-neutral-8);"><icon-copy :size="15"/></a-link>
            </a-tooltip>
            <a-tooltip content="复制 Topic">
             <a-link @click="handleCopyTopic(item.topic)"><icon-font name="icon-a-1" color="var(--color-neutral-8)" size="15"/></a-link>
            </a-tooltip>
            <a-tooltip content="删除">
             <a-link @click="handleDelMsg(item.id)" style="color: var(--color-neutral-8);"><icon-delete :size="15"/></a-link>
            </a-tooltip>
            <span class="creat-time-text">{{dayjs(item.created_at).format("YYYY-MM-DD HH:mm:ss.SSS")}}</span>
          </a-space>
        </div>
      </div>
      <!--发送消息-->
      <div class="msg-body" v-else>
        <div class="msg-wrap">
          <div class="chat-option">
            <a-space style="width: 100%;">
              <div class="">Topic: {{ item.topic }}</div>
              <div class="">QoS: {{ item.qos }}</div>
            </a-space>
          </div>
          <div class="chat-content">
            <div v-if="isJson(item.payload)">
              <pre style="margin:0; padding: 4px; background:var(--color-bg-4); border-radius:4px;">
                <code class="language-json" v-html="highlightJson(item.payload)"></code>
              </pre>
            </div>
            <div v-else class="chat-text">
              {{ item.payload }}
            </div>
          </div>
        </div>
        <div class="message-action-bar" style="text-align: right;">
          <a-space >
            <span class="creat-time-text">{{dayjs(item.created_at).format("YYYY-MM-DD HH:mm:ss.SSS")}}</span>
            <a-tooltip content="复制 Payload">
             <a-link @click="handleCopyTopic(item.payload)" style="color: var(--color-neutral-8);"><icon-copy :size="15"/></a-link>
            </a-tooltip>
            <a-tooltip content="复制 Topic">
             <a-link @click="handleCopyTopic(item.topic)"><icon-font name="icon-a-1" color="var(--color-neutral-8)" size="15"/></a-link>
            </a-tooltip>
            <a-tooltip content="删除">
             <a-link @click="handleDelMsg(item.id)" style="color: var(--color-neutral-8);"><icon-delete :size="15"/></a-link>
            </a-tooltip>
          </a-space>
        </div>
      </div>
    </div>
  </template>
 </div>
</template>
<script setup lang="ts">
  import { ref } from 'vue';
  import { Clipboard } from "@wailsio/runtime";
  import { Message,Modal } from '@arco-design/web-vue';
  import { isValidJson } from '@/utils/jsonBigint'
  import { hexToRgba } from '@/utils/index'
  import { MqttMessageService } from "/#/gofly/internal/service";
  import dayjs from 'dayjs';
  import Prism from 'prismjs'
  import 'prismjs/components/prism-json'   // JSON 语言组件（核心包不含 json）
  import 'prismjs/themes/prism.css'        // 主题，想深色可换 prism-tomorrow.css
  const props = defineProps<{
    list: any[]
  }>()
  const emit = defineEmits<{
    'refresh-message': []
  }>()
  const isJson=(text:any) => {
    return isValidJson(text)
  }
  const highlightJson = (text: string) => {
    return Prism.highlight(text, Prism.languages.json, 'json')
  }
  //复制内容
  const handleCopyTopic=async(text:string)=>{
    await Clipboard.SetText(text)
    Message.success("复制成功")
  }
  //删除单条内容
  const handleDelMsg=async(id:number)=>{
    Modal.warning({
      title: `您确认删除该条内容吗？`,
      content: '删除后将无法恢复，需要谨慎操作哦！',
      hideCancel:false,
      titleAlign:"start",
      onOk:async()=>{
        try {
            Message.loading({content:"删除中",id:"del",duration:0})
            const res =await MqttMessageService.Del(id);
            if(res.code==0){
              Message.success({content:res.message,id:"del",duration:2000})
              emit("refresh-message")
            }else{
              Message.error({content:res.message,id:"del",duration:2000})
            }
        } catch (error) {
          Message.loading({content:"删除中",id:"del",duration:2})
        }
      }
    });
  }
</script>
<style scoped lang="less">
.chat-wrap{
  .chat-msg-item{
    gap: 12px;
    padding: 12px 0;
    &:last-child{
      margin-bottom: 10px;
    }
    &.chat-publish {
      flex-direction: row-reverse;
      .msg-body { 
        align-items: flex-end;
        max-width: calc(100% - 90px);
      }
      .msg-wrap { 
        background: rgb(var(--orange-4)); 
        width: 100%;
         &:before{
          right: -11px;
          vertical-align: middle;
          border-right-color: transparent;
          border-left-color: rgb(var(--orange-4));
        }
        .chat-option{
           color: #ffffff;
           border-bottom: var(--color-fill-1) 1px dashed;
        }
        .chat-text{
          color: #ffffff;
        }
      }
    }
    //接收信息
    &.chat-received {
      .msg-body {
          align-items: flex-start;
          max-width: calc(100% - 90px);
      }
      .msg-wrap { 
        background: var(--color-bg-3); 
        width: 100%;
        &:before{
          content: " ";
          position: absolute;
          top: 12px;
          right: 100%;
          border: 6px solid transparent;
          border-right-color: var(--arrow-color);
          // border-right-color: var(--color-bg-3);
        }
        .chat-option{
           border-bottom: var(--color-bg-3) 1.5px dashed;
        }
      }
    }
    .avatar{
      flex-shrink: 0;
      width: 36px;
      height: 36px;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      background: var(--color-bg-3);
    }
    .msg-wrap { 
      position: relative;
      border-radius: 10px;
      // padding: 12px;
      &:before{
          content: " ";
          position: absolute;
          top: 12px;
          right: 100%;
          border: 6px solid transparent;
      }
      .chat-option{
        padding: 12px 12px 5px 12px;
        .copy-btn{
          cursor: pointer;
        }
        .retain{
          padding: 2px 8px;
          border-radius: 12px;
          background-color: rgb(var(--orange-3));
          border: rgb(var(--orange-4)) 1px solid;
          color: #ffffff;
        }
      }
      .chat-content{
         padding: 10px 12px 12px 12px;
      }
    }
    .message-action-bar{
      margin-top: 5px;
      margin-bottom: -10px;
      .creat-time-text{
        color: var(--color-neutral-6);
        font-size: 13px;
      }
    }
  }
}

:deep(pre) {
  margin: 0;
  padding: 0;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  white-space: pre;
  font-size: inherit;
  border: none;
  overflow: auto;
  display: flex;
  justify-content: start;
}

</style>
