<template>
  <div class="message-publish flex-space-between" >
    <div class="message-hgeader">
      <div class="message-topic flex-align">
        <div class="message-topic-input flex_1">
          <a-input v-model="msgData.topic" ref="topicRef" :style="{width:'90%'}" :placeholder="$t('connections.enter-tig')" />
        </div>
        <div class="message-topic-select">
           <a-dropdown v-model:popup-visible="dropdownVisible" @select="handleTopicSelect" position="tr">
              <div> 主题记录 <icon-up v-if="dropdownVisible"/><icon-down v-else/></div>
              <template #content>
                <a-doption  v-for="item in topicHistory" :value="item">
                  <span style="float: left; max-width: 160px;min-width: 80px; overflow: hidden; text-overflow: ellipsis" >{{item.topic}}</span>
                  <span style="color: #8492a6; font-size: 12px; margin-left:15px">QoS:{{ item.qos }}</span>
                  <span style="float: right; color: #8492a6; font-size: 13px; margin-left: 4px">retain:{{ item.retain ? '1' : '0' }}</span>
                </a-doption>
              </template>
            </a-dropdown>
        </div>
      </div>
    </div>
    <div class="message-editor flex_1" :style="{height:`${inputHeight-36}px`}">
      <div class="publish-editor" >
        <CodeEditor
          v-model:value="msgData.payload"
          :language="languageType"
          :theme="theme=='dark'?'vs-dark':'vs'"
          :options="editorOptions"
          @editorDidMount="handleEditorDidMount"
        />
        </div>
        <div class="message-option flex">
          <div class="message-option-input flex_1">
            <a-space>
             <a-select v-model="payloadType" :style="{width:'95px',borderRadius:'50px'}" size="mini" :trigger-props="{ autoFitPopupMinWidth: true }">
              <a-optgroup label="使用以下格式编码要发布的 Payload">
                <a-option v-for="item in payloadList" :value="item">{{ item }}</a-option>
              </a-optgroup>
             </a-select>
             <a-select v-model="msgData.qos" :style="{width:'95px',borderRadius:'50px'}" size="mini" :trigger-props="{ autoFitPopupMinWidth: true }">
              <a-option v-for="qos in [0, 1, 2]" :key="qos" :label="`QoS ${qos}`" :value="qos">
                 {{ qos }}
                  <span style=" color: var(--color-neutral-6); margin-left: 12px">{{ $t(`connections.qos${qos}`) }}</span>
                </a-option>
             </a-select>
             <a-checkbox class="border-checkbox" v-model="msgData.retain" size="mini">Retain</a-checkbox>
            <a-tooltip content="仅在 MQTT 5.0 中启用" position="top">
             <a-button size="mini" shape="round" v-if="mqtt_version!='5.0'" :disabled="true">Meta</a-button>
            </a-tooltip>
            <a-popconfirm okText="保存" position="top" @ok="handleMeta">
              <template #content>
                <div style="width: 500px;">
                  <a-form ref="topicFormRef" :model="metaData" auto-label-width>
                    <a-form-item field="response_topic" label="内容类型">
                      <a-input v-model="metaData.response_topic" :placeholder="$t('data.form.enter')" />
                    </a-form-item>
                    <a-form-item field="payload_format" label="有效载荷指示器">
                      <a-switch v-model="metaData.payload_format" type="round">
                        <template #checked>ON</template>
                        <template #unchecked>OFF</template>
                      </a-switch>
                    </a-form-item>
                    <a-form-item field="message_expiry" label="消息过期时间">
                      <a-input-number v-model="metaData.message_expiry" :placeholder="$t('data.form.enter')" allow-clear hide-button>
                        <template #suffix>秒</template>
                      </a-input-number>
                    </a-form-item>
                    <a-form-item field="topic_alias" label="主题别名">
                      <a-input-number v-model="metaData.topic_alias" :placeholder="$t('data.form.enter')" />
                    </a-form-item>
                    <a-form-item field="response_topic" label="响应主题">
                      <a-input v-model="metaData.response_topic" :placeholder="$t('data.form.enter')" />
                    </a-form-item>
                    <a-form-item field="correlation_data" label="对比数据">
                      <a-input v-model="metaData.correlation_data" :placeholder="$t('data.form.enter')" />
                    </a-form-item>
                    <a-form-item field="subscription_identifier" label="订阅标识符" style="margin-bottom: 8px;">
                      <a-input-number v-model="metaData.subscription_identifier" :placeholder="$t('data.form.enter')" />
                    </a-form-item>
                    <!--用户属性-->
                    <KeyValueEditor :bordered="false" v-model="metaData.user_properties" :cardMaxHeight="140"/>
                  </a-form>
                </div>
              </template>
              <a-button size="mini" shape="round" v-if="mqtt_version=='5.0'">Meta</a-button>
            </a-popconfirm>
             <a-dropdown v-model:popup-visible="messageVisible" @select="handleMessagesSelect" position="top">
              <div class="message-option-round flex-align">发送记录 <icon-up v-if="messageVisible"/><icon-down v-else/></div>
              <template #content>
                <a-doption v-for="item in messageHistory" :value="item">
                  <span style="float: left; max-width: 200px;min-width: 120px; overflow: hidden; text-overflow: ellipsis" >{{item.payload}}</span>
                </a-doption>
              </template>
              </a-dropdown>
            </a-space>
          </div>
          <div class="message-editor-send">
            <a-button
              class="message-send-btn"
                type="primary"
                shape="circle"
                :status="loading?'warning':'normal'"
                @click="loading ? handleStop() : handleSubmit()"
                :disabled="!loading && !msgData.payload?.trim()"
              >
                <icon-font v-if="!loading" name="icon-send-m" size="15"/>
                <icon-record-stop v-else size="20"/>
              </a-button>
          </div>
        </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref,computed,watch,nextTick,onMounted } from 'vue'
  import { CodeEditor, type EditorOptions } from 'monaco-editor-vue3';
  import convertPayload from '/@/utils/convertPayload';
  import { Message } from '@arco-design/web-vue';
  import { useAppStore } from '@/store';
  import { MessageModel,PushPropertiesModel,MessageHistoryModel,TopicHistoryModel } from '@/types/global'
  import { MqttMessageRecordService,MqttTopicRecordService } from "/#/gofly/internal/service";
  import KeyValueEditor from './KeyValueEditor.vue'
  const appStore = useAppStore();
  const theme = computed(() => {
    return appStore.theme;
  });
  const props = defineProps<{
    loading: boolean,
    inputHeight: number,
    id: number,
    mqtt_version: string,
   }>()
  const emit = defineEmits<{
    submit: [text: MessageModel]
    stop: []
  }>()

  // 手动切换语法高亮
  let editorInstance: any = null
  const handleEditorDidMount = (editor: any) => {
    editorInstance = editor
  }
  const syncEditorLanguage = () => {
    nextTick(() => {
      if (!editorInstance) return
      const model = editorInstance.getModel()
      if (!model) return
      model.setLanguage(languageType.value)
    })
  }
  const payloadType = ref("JSON")
  const msgData = ref<MessageModel>({connection_id:props.id,qos:0,retain: false,topic: '',msg_type: 'publish',payload: JSON.stringify({ msg: 'hello' }, null, 2),create_at:""})
  //mqtt 5.0
  const metaData = ref<PushPropertiesModel>({
      user_properties:[],
      response_topic: "",
      content_type: "",
      correlation_data:  "",
      message_expiry: undefined,
      topic_alias: undefined,
      payload_format:false,
      subscription_identifier: undefined,
  })
  const payloadList = ref(['Plaintext', 'JSON', 'Base64', 'Hex', 'CBOR', 'MsgPack'])
  //主题配置记录和内容记录
  const topicHistory=ref<TopicHistoryModel[]>([])
  const messageHistory=ref<MessageHistoryModel[]>([])
  const getMessageHistory=async()=>{
    const res = await MqttMessageRecordService.GetList()
    if(res?.code==0){
      messageHistory.value=res?.data
    }else{
      messageHistory.value=[]
    }
  }
  const getTopicHistory=async()=>{
    const res = await MqttTopicRecordService.GetList()
    if(res?.code==0){
      topicHistory.value = res?.data.map(item => {
        return {
          ...item, // 复制原有字段
          retain:item.retain==1?true:false // 修改需要的字段
        };
      });
    }else{
      topicHistory.value =[]
    }
  }
  onMounted(async()=>{
    await getMessageHistory()
    await getTopicHistory()
    const mqtt_meta = localStorage.getItem("mqtt_meta")
    if(mqtt_meta){
      metaData.value=JSON.parse(mqtt_meta)
    }
  })

  //监听类型切换
  watch(payloadType, (newVal, oldVal) => {
    convertPayload(msgData.value.payload, newVal, oldVal)
      .then((res) => {
        msgData.value.payload = res
        syncEditorLanguage()
      })
      .catch((error: Error) => {
        const errorMsg = error.toString()
        Message.error(errorMsg)
        payloadType.value = oldVal
      })
  })
  //编辑器语言切换
  const languageType = computed(() => {
     if (['CBOR', 'JSON', 'MsgPack'].includes(payloadType.value)) {
      return 'json'
    } else {
      return 'plaintext'
    }
  });
  
  const dropdownVisible=ref(false)
  const messageVisible=ref(false)
  const topicRef=ref()
  //发送
  const handleSubmit =async () => {
  const text = msgData.value.payload.trim()
    if (!text || props.loading) return
    if(msgData.value.topic==""){
      Message.warning({content:"请输入主题名称",id:"topic"})
      topicRef.value.focus()
      return
    }
    msgData.value.connection_id=props.id
    // if(languageType.value=="json"){
    //   msgData.value.payload=JSON.parse(text)
    // }
    await saveHistory()
    emit('submit', msgData.value)
  }

  // 停止生成
  const handleStop = () => {
    emit('stop')
  }
  // 保存记录
  const saveHistory = async() => {
    await MqttMessageRecordService.Save(JSON.stringify({payload:msgData.value.payload,payload_type:payloadType.value}))
    await MqttTopicRecordService.Save(JSON.stringify({qos:msgData.value.qos,retain:msgData.value.retain?1:0,topic:msgData.value.topic}))
    RefreshRecord(true)
  }
  //重新获取记录数据
  const RefreshRecord=async(getData:boolean)=>{
    if(getData){
      await getMessageHistory()
      getTopicHistory()
    }else{
      topicHistory.value =[]
      messageHistory.value=[]
    }
  }
  //历史主题配置
  const handleTopicSelect=(item:any)=>{
   if (item) {
      const { retain, topic, qos } = item
      Object.assign(msgData.value, { retain, topic, qos })
    }
  }
  //发送记录
  const handleMessagesSelect=(item:any)=>{
   if (item) {
      payloadType.value = item.payloadType
      Object.assign(msgData.value, { payload: item.payload })
      syncEditorLanguage()
    }
  }
  //保存MQTT5.0参数
  const handleMeta=()=>{
    localStorage.setItem("mqtt_meta", JSON.stringify(metaData.value))
  }
  // 代码编辑器参数
  const editorOptions: EditorOptions = {
    fontSize: 13,
    wordWrap: 'off',
    lineNumbers: 'off',
    lineNumbersMinChars: 1,
    renderLineHighlight: 'none',
    scrollbar: {
      horizontal: 'auto',
      horizontalScrollbarSize: 8,
      vertical: 'auto',
      verticalScrollbarSize: 8,
      useShadows: false,
      alwaysConsumeMouseWheel: false,
    },
    minimap: { enabled: false },
    smoothScrolling: true,
    scrollBeyondLastLine: false,
    matchBrackets: 'near',
    folding: false,
    lightbulb: { enabled: false },
    selectionHighlight: false,
    occurrencesHighlight: false,
    overviewRulerBorder: false,
    automaticLayout: true,
    renderWhitespace: 'none',
    renderControlCharacters: false,
  };
  defineExpose({
    RefreshRecord
  });
</script>

<style lang="less" scoped>
.message-publish {
  .message-hgeader{
    border-bottom: var(--color-neutral-2) 1px solid;
    .message-topic{
      padding: 2px 16px 2px 0px;
    }
  }
  .publish-editor{
    height: calc(100% - 35px);
    width: 100%;
    padding: 10px 1px 1px 1px;
  }
  .message-option{
    height: 35px;
    padding: 0px 12px;
    .message-option-round{
      background-color: var(--color-fill-2);
      border-radius: 50px;
      min-height: 24px;
      padding: 0px 12px;
      &:hover{
        background-color: var(--color-fill-3);
      }
    }
    .message-editor-send{
      padding-right: 5px;
     .message-send-btn{
      margin-top: -10px;
     }
    }
  }
}
.message-topic-input{
  :deep(.arco-input-wrapper){
    background-color:transparent;
    border: 0px !important;
  }
}

.border-checkbox{
  background-color: var(--color-fill-2);
  border-radius: 50px;
  min-height: 24px;
  padding: 0px 12px;
}
</style>
