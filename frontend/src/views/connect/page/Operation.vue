<template>
  <a-layout-header class="layout-header" ref="headerRef">
      <div class="option-header flex flex-middle">
        <div class="flex_1 flex flex-middle">
          <div class="option-title">{{ formData.title }}</div>
          <div class="option-up-down" >
            <a-link @click="collapseShow=!collapseShow">
              <icon-double-up v-if="collapseShow" :size="18"/>
              <icon-double-down v-else :size="18"/>
            </a-link>
          </div>
        </div>
        <div class="">
          <a-space>
            <a-link style="color: var(--color-neutral-10);" @click="handleStartAddStop">
              <template v-if="formData.connected==1">
              <icon-refresh v-if="connectedLoading" :size="18" spin style="color: rgb(var(--arcoblue-6));"/>
              <icon-poweroff v-else :size="18" style="color: rgb(var(--red-6));"/>
              </template>
              <icon-play-arrow-fill v-else :size="18"/>
            </a-link>
            <a-link style="color: var(--color-neutral-10);" @click="handleEdit" :disabled="formData.connected==1"><icon-edit :size="18"/></a-link>
            <a-dropdown position="br">
              <a-link style="color: var(--color-neutral-10);"><icon-more :size="18"/></a-link>
              <template #content>
                <a-doption @click="handleClearRecord">
                  <template #icon><icon-font name="icon-clear-data" size="14"/></template>
                  <template #default>清除记录</template>
                </a-doption>
                <a-doption @click="handleDelMessage">
                  <template #icon><icon-font name="icon-shanchuduihua" size="14"/></template>
                  <template #default>删除消息</template>
                </a-doption>
                <a-doption class="doption-delete" @click="handleDelConnect">
                  <template #icon><icon-delete /></template>
                  <template #default>删除连接</template>
                </a-doption>
              </template>
            </a-dropdown>
          </a-space>
        </div>
      </div>
      <collapse-transition>
        <div class="option-info" v-show="collapseShow">
          <a-form ref="formRef" :disabled="formData.connected==1" size="small" :model="formData" layout="vertical" auto-label-width>
            <a-row :gutter="16">
              <a-col :span="8">
                <a-form-item field="title" label="名称" validate-trigger="input" :rules="[{required:true,message:'请填写名称'}]" style="margin-bottom: 10px;">
                  <a-input v-model="formData.title" :placeholder="$t('data.form.enter')" @change="handleConnenct({id:formData.id,title:formData.title})"/>
                </a-form-item>
              </a-col>
              <a-col :span="8">
                <a-form-item field="client_id" label="Client ID" validate-trigger="input" :rules="[{required:true,message:'请填写端口'}]" style="margin-bottom: 10px;">
                  <a-input v-model="formData.client_id" :placeholder="$t('data.form.enter')" @change="handleConnenct({id:formData.id,client_id:formData.client_id})"/>
                  <a href="javascript:;" class="icon-btn" @click="handleClientID"><icon-refresh /></a>
                </a-form-item>
              </a-col>
              <a-col :span="8">
                <a-form-item field="username" label="用户名" style="margin-bottom: 10px;">
                  <a-input v-model="formData.username" :placeholder="$t('data.form.enter')" @change="handleConnenct({id:formData.id,username:formData.username})"/>
                </a-form-item>
              </a-col>
              <a-col :span="8">
                <a-form-item field="password" label="密码">
                  <a-input-password v-model="formData.password" :placeholder="$t('data.form.enter')" @change="handleConnenct({id:formData.id,password:formData.password})"/>
                </a-form-item>
              </a-col>
              <a-col :span="8">
                <a-form-item field="keepalive" label="Keep Alive">
                  <a-input-number v-model="formData.keepalive" :placeholder="$t('data.form.enter')" @change="handleConnenct({id:formData.id,keepalive:formData.keepalive})"/>
                </a-form-item>
              </a-col>
              <a-col :span="8">
                <a-form-item field="clean" :label="formData.mqtt_version === '5.0' ? 'Clean Start' : 'Clean Session'" >
                  <a-checkbox v-model="formData.clean" :value="1" @change="handleConnenct({id:formData.id,keepcleanalive:formData.clean})">{{ formData.clean?true:false }}</a-checkbox>
                </a-form-item>
              </a-col>
            </a-row>
          </a-form>
        </div>
      </collapse-transition>
  </a-layout-header >
  <a-layout>
    <a-layout-sider
      hide-trigger
      :style="{width: '230px',height:`calc(100vh - ${headerHeight}px)`}"
      class="connections-sider"
    >
      <div class="topic-box flex-space-between">
        <div class="topic-content flex_1">
            <template v-for="item in topicList">
              <a-dropdown trigger="contextMenu" alignPoint :popup-max-height="false" :style="{display:'block'}"> 
                <div :class="['topic-item flex-align', { active: item.id === topicActiveIndex, disabled: item.disabled==1}]" 
                @click="selectTopic(item.id||0,item.topic)"
                :style="{backgroundColor:hexToRgba(item.color),color:(item.id === topicActiveIndex?'var(--color-bg-3)':item.color)}">
                  <div class="topic-color-line" :style="{background: item.id === topicActiveIndex?'var(--color-bg-3)':`${item.color}`}"></div> 
                  <div class="topic-item-text ellipsis-l1e flex_1">
                    <a-tooltip :background-color="copySuccess?'#00b42a':''" @popup-visible-change="handleCopyTip" :content="getTopicText(item.alias || item.topic)" position="top">
                      <a @click.stop="handleCopyTopic(item.alias || item.topic)">{{ item.alias || item.topic }}</a>
                    </a-tooltip>
                  </div>
                  <div class="topic-item-qos">QoS {{ item.qos }}</div>
                </div>
                <template #content>
                  <a-doption @click="handleEditTopic(item)">
                    <template #icon><icon-edit /></template>
                    <template #default>编 辑</template>
                  </a-doption>
                  <a-doption @click="handleDisabledTopic(item)">
                    <template #icon><icon-check-circle v-if="item.disabled"/><icon-minus-circle v-else/></template>
                    <template #default>{{item.disabled?"启 用":"禁 用"}}</template>
                  </a-doption>
                  <a-doption @click="handleDalTopic(item)" class="doption-delete">
                    <template #icon><icon-delete /></template>
                    <template #default>删 除</template>
                  </a-doption>
                </template>
              </a-dropdown>
            </template>
            <div v-if="!topicList||topicList.length==0" style="margin-top: 35%;">
              <a-dropdown trigger="contextMenu" alignPoint :style="{display:'block'}">
                <a-empty  description="请添加订阅主题"/>
                <template #content>
                  <a-doption @click="handleAddTopic">
                    <template #icon><icon-plus /></template>
                    <template #default>添加订阅</template>
                  </a-doption>
                </template>
              </a-dropdown>
            </div>
        </div>
        <div class="topic-footer" @click="handleAddTopic">
         <span class="common-icon img-middle"><icon-plus /></span>
         <span class="common-text ellipsis-l1e img-middle">添加订阅</span>
        </div>
      </div>
    </a-layout-sider>
    <a-layout-content :style="{height:`calc(100vh - ${headerHeight}px)`,overflow:'hidden'}">
      <div class="message-box flex-space-between" :style="{height:`calc(100vh - ${headerHeight}px)`}">
        <div class="message-header flex-align">
          <div class="message-header-left flex_1"></div>
          <div class="message-header-right">
            <a-space >
              <template #split>
                <a-divider direction="vertical" :margin="6"/>
              </template>
              <a :class="['message-type-btn', { active: msgType === 'all' }]" @click="changeValue('all')">全部</a>
              <a :class="['message-type-btn', { active: msgType === 'received' }]" @click="changeValue('received')">已接收</a>
              <a :class="['message-type-btn', { active: msgType === 'publish' }]" @click="changeValue('publish')">已发送</a>
            </a-space>
          </div>
        </div>
        <div class="message-content flex_1" ref="messageContentRef">
          <!--信息对话区-->
          <ChatView :list="msglist" @scrollTop="handleScrollTop" @refresh-message="handleRefreshMessage"/>
        </div>
        <div class="message-footer" :style="{ height: `${inputHeight+2}px` }">
          <ResizeHeight v-model="inputHeight" />
          <MessagePublish 
            ref="messagePublishRef"
            :id="id" 
            :loading="sendLoading" 
            :mqtt_version="formData.mqtt_version" 
            :inputHeight="inputHeight" 
            :style="{ height: `${inputHeight}px` }"
            @submit="handleSendMsg"
          />
        </div>
      </div>
    </a-layout-content>
  </a-layout>
  <!--主题表单-->
  <a-modal v-model:visible="topicData.visible" width="700px" draggable :closable="true" :mask-closable="false" title-align="start" @before-ok="handleOkTopic">
    <template #title>
      <span>
        <a-popover title="订阅说明" position="bl"> 
          <icon-exclamation-circle :size="18" class="img-middle topic-exclamation"/>
          <template #content>
            <p style="width: 370px;">添加订阅后，且MQTTW连接服务器后会自动重新订阅本地保存的主题，无论主题是否在 MQTT Broker 持久化。</p>
          </template>
        </a-popover>
         添加订阅</span>
    </template>
    <div class="topic-wrap" :style="{height: `${viewHeight<450?viewHeight:450}px`}">
      <a-form ref="topicFormRef" :model="topicData.form" auto-label-width>
        <a-form-item field="topic" label="Topic" validate-trigger="input" tooltip="可订阅一个或多个主题。订阅多个主题时，请用逗号, 分隔每个主题。例如：test1,test2" :rules="[{required:true,message:'请填写订阅主题'}]">
          <a-textarea v-model="topicData.form.topic" auto-size />
        </a-form-item>
        <a-form-item field="qos" label="QoS" validate-trigger="input" :rules="[{required:true,message:'请选择QoS'}]">
          <a-select v-model="topicData.form.qos" :trigger-props="{ autoFitPopupMinWidth: true }">
            <a-option v-for="qos in [0, 1, 2]" :key="qos" :value="qos">
              {{ qos }}
              <span style="min-width: 100px;max-width: 160px;color: var(--color-neutral-6); margin-left: 12px;float: right;">{{ $t(`connections.qos${qos}`) }}</span>
            </a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="topicColor" label="标记颜色" validate-trigger="input"  >
          <a-input v-model="topicData.form.color" placeholder="选择或输入标签颜色、展示Topic标签时颜色" allow-clear >
            <template #suffix>
              <a-color-picker v-model="topicData.form.color" showPreset :trigger-props="{position: 'right',showArrow:true}"/>
            </template>
          </a-input>
          <a href="javascript:;" class="icon-btn" @click="handleTopicColor"><icon-refresh /></a>
        </a-form-item>
        <a-form-item field="alias" label="别名" validate-trigger="input" tooltip="为多主题设置别名时，也使用逗号分隔符（,）">
          <a-textarea v-model="topicData.form.alias" auto-size />
        </a-form-item>
        <template v-if="formData.mqtt_version=='5.0'">
          <a-form-item field="subscription_identifier" label="订阅标识符" validate-trigger="input" tooltip="取值范围 1 ~ 268435455">
            <a-input-number :min="1" :max="268435455" v-model="topicData.form.subscription_identifier" :placeholder="$t('data.form.enter')" />
          </a-form-item>
          <a-form-item field="nl" label="禁止本地转发" validate-trigger="input"  >
            <a-radio-group v-model="topicData.form.nl">
              <a-radio :value="0">false</a-radio>
              <a-radio :value="1">true</a-radio>
            </a-radio-group>
          </a-form-item>
            <a-form-item field="rap" label="发布时状态保留" validate-trigger="input"  >
            <a-radio-group v-model="topicData.form.rap">
              <a-radio :value="0">false</a-radio>
              <a-radio :value="1">true</a-radio>
            </a-radio-group>
          </a-form-item>
          <a-form-item field="rh" label="保留消息处理" validate-trigger="input"  >
            <a-select v-model="topicData.form.rh"  >
            <a-option v-for="qos in [0, 1, 2]" :key="qos" :label="`${qos}`" :value="qos">
                <div class="flex-align">
                  <div class="flex_1">{{ qos }}</div>
                  <div class="">
                    <span style="color: var(--color-neutral-6); margin-left: 12px;">{{ $t(`connections.topic${qos}`) }}</span>
                    <a-tag style=" margin-left: 12px;">Retain:{{ qos==0?'false':'true' }}</a-tag>
                  </div>
                </div>
              </a-option>
            </a-select>
          </a-form-item>
        </template>
      </a-form>
    </div>
  </a-modal>
</template>

<script lang="ts" setup>
  import { ref,computed, onMounted,onUnmounted, onBeforeUnmount, nextTick,watch } from 'vue';
  import { LayoutHeader,FormInstance,Modal,Message,Notification } from '@arco-design/web-vue'
  import { Clipboard, Events } from "@wailsio/runtime";
  import ResizeHeight from '../components/ResizeHeight.vue'
  import ChatView from '../components/ChatView.vue'
  import MessagePublish from '../components/MessagePublish.vue'
  import { hexToRgba } from '@/utils/index'
  import { getRandomColor,getNextTopicColor } from '@/utils/colors'
  import { getClientId } from '@/utils/generateRandom'
  import { MessageModel } from '@/types/global'
  import dayjs from 'dayjs'
  import { matchTopicMethod } from '@/utils/topicMatch'
  import { useMqttStore } from '@/store'
  const mqttStore = useMqttStore()

  //go 数据接口
  import {
    MqttSubscriptionService,
    MqttConnectionService,
    MqttMessageService,
    MqttClientService,
  } from "/#/gofly/internal/service";

  const props = defineProps<{
    id: number
    start: boolean
  }>()
  const emit = defineEmits<{
    'add-connect': [index: number,from_id:number]
  }>()
  const messagePublishRef=ref<InstanceType<typeof MessagePublish>>()
  const topicActiveIndex=ref<number>()
  const topicActive=ref<string>("all")
  const formRef=ref<FormInstance>();
  const topicFormRef=ref<FormInstance>();
  // 获取窗口高度
  const viewHeight = ref(window.innerHeight)
  let timer: number | null = null
  async function refreshSize() {
    viewHeight.value = window.innerHeight-162;//减去弹框上下高度
  }
  // 防抖包装
  function handleResize() {
    if (timer) clearTimeout(timer)
    timer = window.setTimeout(() => {
      refreshSize()
    }, 80)
  }

  //对话窗口
  const rawInputHeight = ref(180)
  const MIN_HEIGHT = 120
  const MAX_HEIGHT = 450
  const inputHeight = computed({
    get() {
      return rawInputHeight.value
    },
    set(val: number) {
      // 钳位，限制区间
      rawInputHeight.value = Math.max(MIN_HEIGHT, Math.min(MAX_HEIGHT, val))
    }
  })
  const msgType = ref("all")
  const sendLoading = ref(false);//发送消息状态
  // 组件实例Ref
  const headerRef=ref<InstanceType<typeof LayoutHeader> | null>(null)
  const messageContentRef = ref<HTMLDivElement | null>(null)
  // 实时header高度
  const headerHeight = ref(0)
  let resizeObserver: ResizeObserver | null = null

  const formData=ref({
      id:0,
      title:"",
      client_id:"",
      username:"",
      password:"",
      keepalive:10,
      clean:true,
      mqtt_version:"5.0",
      connected:0,
    })
  const collapseShow=ref(false)
  const connectedLoading=ref(false)//连接动画状态
  //获取连接数据
  const getConnectionData=async()=>{
    const res =await MqttConnectionService.GetInfo(props.id);
    if(res.code==0){
      formData.value=res.data
      //判断打开界面是否进行连接
      if(props.start&&formData.value.connected==0){
        handleStartAddStop()
      }
    }else{
      formRef.value?.resetFields()
    }
  }
  //连接、断开
  const handleStartAddStop=async()=>{
    formData.value.connected= formData.value.connected==1?0:1
    if(formData.value.connected==1){//连接
      collapseShow.value=false
      try {
        connectedLoading.value=true
        const res =await MqttClientService.Connect(props.id);
        connectedLoading.value=false
        if(res.code==0){
          Notification.info({
            title: '连接和订阅主题结果',
            content: res.message,
            position: 'bottomRight',
            closable: true,
          })
        }else{
          formData.value.connected=0
          Message.error({content:res.message,id:"connect",duration:2000})
        }
      } catch (error) {
        connectedLoading.value=false
      }
    }else{//断开
      const res =await MqttClientService.Disconnect(props.id);
      if(res.code==0){
          Notification.info({
            title: '连接断开结果',
            content: res.message,
            position: 'bottomRight',
            closable: true,
          })
        }else{
          Message.error({content:res.message,id:"connect",duration:2000})
        }
    }
  }
  //编辑
  const handleEdit=()=>{
    emit('add-connect', 2,props.id)
  }
  //信息对话记录
  const msglist=ref<MessageModel[]>([])
  //获取消息数据
  const getMessageData=async()=>{
    const res =await MqttMessageService.GetList(props.id,topicActive.value,msgType.value);
    if(res.code==0){
      msglist.value=res.data.map(item => {
        return {
          ...item, // 复制原有字段
          color:getTopicColor(item.topic) // 修改需要的字段
        };
      });
    }else{
      msglist.value=[]
    }
    await nextTick()
    handleScrollTop()
  }
  //刷新信息列表
  const handleRefreshMessage=()=> {
    topicActive.value="all"
    getMessageData()
  }
  //获取主题颜色
  const getTopicColor=(topic:string):string=>{
     return topicList.value.find((sub: TopicFormItem)=> matchTopicMethod(sub.topic, topic))?.color || getRandomColor()
  }
  // 接收 MQTT 消息，推送到 msglist（mqtt:all_message 由后端 Connect 回调发出）
  const handleMqttMessage = (ev: any) => {
    if (!ev?.data) return
    const data = JSON.parse(ev.data)
    if (!data || data.connection_id !== props.id) return   // 只接收当前连接的消息
    // 按当前选中的主题过滤（topicActive = "all" 时不限）
    if (topicActive.value !== 'all' && data.topic !== topicActive.value) return
    msglist.value.push({
      connection_id: data.connection_id,
      msg_type: 'received',                      // ChatView 据此显示为左侧接收消息
      qos: data.qos,
      retain: !!data.retain,
      topic: data.topic,
      payload: data.payload,                     // 后端已转字符串，直接可用
      color: getTopicColor(data.topic),
      user_properties: data.user_properties,
      response_topic: data.response_topic,
      content_type: data.content_type,
      correlation_data: data.correlation_data,
      message_expiry: data.message_expiry,
      create_at: dayjs().format('YYYY-MM-DD HH:mm:ss'),
    })
    nextTick(() => handleScrollTop())            // 滚动到底部
  }

  //发送消息
  const handleSendMsg=async(msgData:MessageModel)=>{
      try {
        Message.loading({content:"发送中",id:"save",duration:0})
        sendLoading.value=true
        //1.发送到mqtt服务器
        var metaData ={}
        if(formData.value.mqtt_version=="5.0"){
          const mqtt_meta = localStorage.getItem("mqtt_meta")
          if(mqtt_meta){
            metaData =JSON.parse(mqtt_meta)
          }
        }
         const sendRes =await MqttClientService.Publish(JSON.stringify({
          connection_id: msgData.connection_id,
          topic: msgData.topic,
          payload: msgData.payload,
          qos: msgData.qos,
          retain: msgData.retain,
          ...metaData
        }))
         //2.保存数据到本地
         if(sendRes.code==0){
            let retainVal=0
            if(msgData.retain){
              retainVal==1
            }
            Object.assign(msgData,{retain:retainVal})
            const res =await MqttMessageService.Save(JSON.stringify(msgData));
            sendLoading.value=false
            Message.loading({content:"发送中",id:"save",duration:2})
            if(res.code==0){
              getMessageData()
              // Message.success({content:res.message,id:"save",duration:2000})
            }else{
              Message.error({content:res.message,id:"save",duration:2000})
            }
         }else{
          Message.error({content:sendRes.message,id:"save",duration:2000})
          sendLoading.value=false
         }
      } catch (error) {
        sendLoading.value=false
        Message.loading({content:"发送中",id:"save",duration:2})
      }
  }
  //清除记录
  const handleClearRecord=async()=>{
      try {
        Message.loading({content:"清除记录中",id:"save",duration:0})
        const res =await MqttMessageService.ClearRecord();
        if(res.code==0){
          Message.success({content:res.message,id:"save",duration:2000})
          if (messagePublishRef.value) {
            messagePublishRef.value.RefreshRecord(false)
          }
        }else{
          Message.error({content:res.message,id:"save",duration:2000})
        }
      } catch (error) {
        Message.loading({content:"清除记录中",id:"save",duration:2})
      }
  }
  //添加主题
  interface TopicFormItem {
    id:number | undefined
    connection_id: number
    topic: string
    qos: number
    color: string
    alias: string
    subscription_identifier: number | undefined
    nl: number
    rap: number
    rh: number
    disabled: number
  }
  // 获取表单默认值
  const getDefaultTopicForm = (): TopicFormItem => ({
    id:undefined,
    connection_id:props.id,
    topic: "topic/#",
    qos: 0,
    color: getNextTopicColor(topicList.value.length||0)||"#16C9FF",
    alias: "",
    subscription_identifier: undefined,
    nl: 0,
    rap: 0,
    rh: 0,
    disabled: 0,
  })
  const topicList=ref<TopicFormItem[]>([])
  const topicData=ref({
    visible:false,
    form:getDefaultTopicForm()
  })
  //刷新颜色
  const handleTopicColor=()=>{
    topicData.value.form.color=getRandomColor()
  }
  //选择主题
  const selectTopic=(id:number,topic:string)=>{
    if(topicActiveIndex.value==id){
      topicActiveIndex.value=0
      topicActive.value="all"
      getMessageData()
    }else{
      topicActiveIndex.value=id
      topicActive.value=topic
      getMessageData()
    }
  }
  //监听props数据
  watch(
  () => ({ ...props }), // getter 返回新对象，简单浅拷贝
  async(newVal, oldVal) => {
    await Events.Off(`mqtt:message:${oldVal.id}`)
    await Events.On(`mqtt:message:${newVal.id}` , handleMqttMessage)
    getConnectionData()
    getTopicData()
    topicActiveIndex.value=0
    topicActive.value="all"
    sendLoading.value=false
    getMessageData()
  }, { deep: true } )
  //获取主题数据
  const getTopicData=async()=>{
    const res =await MqttSubscriptionService.GetList(props.id);
    if(res.code==0){
       topicList.value=res.data
    }else{
       topicList.value=[]
    }
  }
  // 显示添加表单弹框
  const handleAddTopic=()=>{
    topicData.value.visible=true
    topicData.value.form=getDefaultTopicForm()
  }
  //提交主题
  const handleOkTopic = async(done) => {
    try {
      const validate = await topicFormRef.value?.validate();
      if (!validate) {
        Message.loading({content:"提交中",id:"save",duration:0})
        const res =await MqttSubscriptionService.Save(JSON.stringify(topicData.value.form));
        if(res.code==0){
          getTopicData()
          done()
          Message.success({content:res.message,id:"save",duration:2000})
          if(res.exdata=="add"){
            MqttClientService.Subscribe(res.data)
          }else{
           await MqttClientService.UnSubscribe(res.data,res.exdata)
            MqttClientService.Subscribe(res.data)
          }
        }else{
          Message.error({content:res.message,id:"save",duration:2000})
          done(false)
        }
      }else{
        done(false)
      }
    } catch (error) {
      Message.loading({content:"提交中",id:"save",duration:2})
    }
  };
  //编辑主题
  const handleEditTopic=(item:TopicFormItem)=>{
    topicData.value.form= Object.assign({},getDefaultTopicForm(),item)
    topicData.value.visible=true
  }
  //禁用主题
  const handleDisabledTopic=async(item:TopicFormItem)=>{
    if(item.id){
      try {
        Message.loading({content:"提交中",id:"save",duration:0})
        const res =await MqttSubscriptionService.Disable(item.id,item.disabled?0:1);
        if(res.code==0){
          getTopicData()
          Message.success({content:res.message,id:"save",duration:2000})
          if(item.disabled==0){
            MqttClientService.UnSubscribe(item.connection_id,item.topic)
          }else{
            MqttClientService.Subscribe(item.id)
          }
        }else{
          Message.error({content:res.message,id:"save",duration:2000})
        }
      } catch (error) {
        Message.loading({content:"提交中",id:"save",duration:2})
      }
    }
  }
  //删除主题
  const handleDalTopic=async(item:TopicFormItem)=>{
      Modal.warning({
        title: `您确认删除 ${item.topic} 主题？`,
        content: '删除后将无法恢复，需要谨慎操作哦！',
        hideCancel:false,
        titleAlign:"start",
        onOk:async()=>{
          try {
              Message.loading({content:"删除中",id:"del",duration:0})
              const res =await MqttSubscriptionService.Del(item.id||0);
              if(res.code==0){
                Message.success({content:res.message,id:"del",duration:2000})
                getTopicData()
                MqttClientService.UnSubscribe(item.connection_id,item.topic)
              }else{
                Message.error({content:res.message,id:"del",duration:2000})
              }
          } catch (error) {
            Message.loading({content:"删除中",id:"del",duration:2})
          }
        }
      });
  }
  //删除信息内容
  const handleDelMessage=()=>{
    Modal.warning({
      title: `您确认删除信息内容吗？`,
      content: '删除后将无法恢复，需要谨慎操作哦！',
      hideCancel:false,
      titleAlign:"start",
      onOk:async()=>{
        try {
            Message.loading({content:"删除内容中",id:"del",duration:0})
            const res =await MqttMessageService.ClearMessage(props.id);
            if(res.code==0){
              msglist.value=[]
              Message.success({content:res.message,id:"del",duration:2000})
            }else{
              Message.error({content:res.message,id:"del",duration:2000})
            }
        } catch (error) {
          Message.loading({content:"删除内容中",id:"del",duration:2})
        }
      }
    });
  }
  //删除连接
  const handleDelConnect=async()=>{
    Modal.warning({
      title: `您确认删除 ${formData.value.title} 连接吗？`,
      content: '删除后将无法恢复，需要谨慎操作哦！',
      hideCancel:false,
      titleAlign:"start",
      onOk:async()=>{
        try {
            Message.loading({content:"删除连接中",id:"del",duration:0})
            const res =await MqttConnectionService.Del(formData.value.id);
            if(res.code==0){
              Message.success({content:res.message,id:"del",duration:2000})
              emit('add-connect', 0,props.id)
            }else{
              Message.error({content:res.message,id:"del",duration:2000})
            }
        } catch (error) {
          Message.loading({content:"删除连接中",id:"del",duration:2})
        }
      }
    });
  }
  //复制
  const copySuccess=ref(false)
  const handleCopyTopic=async(text:string)=>{
     await Clipboard.SetText(text)
      copySuccess.value = true
      setTimeout(() => {
        copySuccess.value = false
      }, 1500)
  }
  const handleCopyTip=(visible:boolean)=>{
    if(!visible){
      copySuccess.value = false
    }
  }
  const getTopicText=(text:string)=>{
    if (copySuccess.value) {
      return "复制成功"
    }else{
      return text
    }
  }
  // 信息滚动到底部
  const handleScrollTop=()=>{
    nextTick(() => {
      if (messageContentRef.value) {
        // 方式1：直接设置滚动高度（最常用）
        messageContentRef.value.scrollTop = messageContentRef.value.scrollHeight
        // 方式2：平滑滚动
        // messageContentRef.value.scrollTo({
        //   top: messageContentRef.value.scrollHeight,
        //   behavior: 'smooth'
        // })
      }
    })
  }
  // 监听header高度变化
  onMounted(async ()=>{
    await getConnectionData()
    await getTopicData()
    topicActiveIndex.value=0
    topicActive.value="all"
    getMessageData()
    await nextTick()
    Events.On(`mqtt:message:${props.id}` , handleMqttMessage)
    mqttStore.setCurrentConn(props.id)
    if(!headerRef.value) return
    const dom = headerRef.value.$el as HTMLElement
    const updateHeight = ()=>{
      headerHeight.value = dom.offsetHeight
    }
    updateHeight()
    resizeObserver = new ResizeObserver(()=>{
      updateHeight()
    })
    resizeObserver.observe(dom)
    //高度监听
    await refreshSize()
    window.addEventListener('resize', handleResize)
  })
  onUnmounted(() => {
    Events.Off(`mqtt:message:${props.id}`)
    mqttStore.setCurrentConn(0)
    window.removeEventListener('resize', handleResize)
    if(timer) clearTimeout(timer)
  })
  onBeforeUnmount(()=>{
    resizeObserver?.disconnect()
  })
  //生成id字符串
  const handleClientID=async()=>{
    formData.value.client_id = getClientId()
    handleConnenct({id:formData.value.id,client_id:formData.value.client_id})
  }
  //更新连接内容
  const handleConnenct=(fields:any)=>{
    MqttConnectionService.Update(JSON.stringify(fields));
  }
  //消息查询
  const changeValue=(val: string)=> {
    msgType.value=val
    getMessageData()
  }
</script>

<style scoped lang="less">
  .layout-header{
    border-bottom: var(--color-neutral-2) 1px solid;
  }
  .option-header {
    height: 50px;
    padding: 0px 16px;
    .option-up-down{
      cursor: pointer;
      padding: 6px;
      color: rgb(var(--arcoblue-5));
    }
  }
  :deep(.arco-form-item-label-col){
    margin-bottom: 0px;
  }
  //表单
  .option-info{
    padding: 0px 16px;
  }
  //淡入淡出
  .fade-collapse-enter-from,
  .fade-collapse-leave-to {
    opacity: 0;
    max-height: 0px;
  }
  .fade-collapse-enter-active,
  .fade-collapse-leave-active {
    transition: all 0.24s ease;
    overflow: hidden;
  }
  .fade-collapse-enter-to,
  .fade-collapse-leave-from {
    opacity: 1;
    max-height: 1000px; /* 设置一个足够大的值，大于内容最大高度 */
  }
  //内容操作
  //主题
  .topic-box{
    height: 100%;
    overflow-y: auto;
    .topic-content{
      overflow-y: auto;
      padding: 5px;
      .topic-list{
       overflow-y: auto;
      }
      .topic-item{
        margin-bottom: 10px;
        padding: 10px;
        border-radius: 10px;
        overflow: hidden;
        cursor: pointer;
        &:last-child{
          margin-bottom: 0px;
        }
        &:hover {
          transform: translateX(2px);
          box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08), 0 1px 3px rgba(0, 0, 0, 0.06);
        }
        &.active{
          background: var(--color-neutral-8) !important;
          box-shadow: none;
        }
        &.disabled{
          background: transparent !important;
          border: 1px solid var(--color-neutral-3);
          box-shadow: none;
          cursor: not-allowed;
          color: var(--color-neutral-3) !important;
          &:hover {
            transform: none;
            box-shadow: none;
          }
          .topic-color-line{
            background: var(--color-neutral-3) !important;
          }
        }
        .topic-color-line{
          width: 4px;
          height: 22px;
          border-radius: 0 4px 4px 0;
        }
        .topic-item-text{
          padding: 0px 5px;
          a:hover{
            opacity: 0.8;
          }
        }
      }
    }
    .topic-footer{
      text-align: center;
      padding: 12px 5px;
      border-top: var(--color-neutral-2) 1px solid;
      cursor: pointer;
      &:hover{
        background-color: var(--color-neutral-1);
      }
    }
  }
  //消息内容
  .message-box{
    height: 100%;
    .message-header{
      height: 30px;
      padding: 0px 16px;
      .message-type-btn{
        color: var(--color-neutral-8);
        font-size: 12px;
        cursor: pointer;
        &.active {
          color: rgb(var(--arcoblue-6));
        }
        &:hover{
          color: var(--color-neutral-10);
        }
      }
    }
    .message-content{
      padding: 0px 16px;
      background-color: var(--color-fill-1);
      border-bottom: var(--color-neutral-2) 1px solid;
      overflow-y: auto;
    }
  }
  //主题
  .topic-exclamation{
    cursor: pointer;
    &:hover{
      color: rgb(var(--arcoblue-5));
    }
  }
  .topic-wrap{
    min-height: 339px;
    overflow-y: auto;
  }
  :deep(.arco-select-option-content){
    width: 100%;
  }
</style>