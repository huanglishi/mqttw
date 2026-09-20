<template>
  <a-layout-content class="layout-content">
    <div class="connect-nav">
    <div class="connect-nav-left">
      <div class="back" @click="back"><icon-left :size="16" /></div>
      <div class="name-desc"   >
        <div class="name" v-if="formData.title">{{ formData.title }}</div>
        <div class="desc" v-if="formData.client_id">{{ formData.client_id }}</div>
      </div>
    </div>
    <!-- :pid={{ formData.pid }}，id={{ formData.id }}，from_id={{ from_id }} -->
    <div class="connect-nav-mid">{{formData.id>0?"编辑连接":"新建连接"}}</div>
    <div class="connect-nav-right">
      <a-dropdown-button @click="handleSubmit(1)" :loading="loading" :disabled="loading" type="primary">
          连 接
          <template #icon>
              <icon-down />
          </template>
          <template #content>
              <a-doption @click="handleSubmit(0)">
                  <template #icon><icon-save /></template>
                  <template #default>仅保存</template>
              </a-doption>
          </template>
      </a-dropdown-button>
    </div>
  </div>

  <div class="connect-content">
    <a-form ref="formRef" :model="formData" auto-label-width>
      <a-card title="基本信息" class="gf-mx gf-bottom">
          <a-row :gutter="16">
            <a-col :span="24">
              <a-form-item field="title" label="名称" validate-trigger="input" :rules="[{required:true,message:'请填写名称'}]">
                <a-input v-model="formData.title" :placeholder="$t('data.form.enter')" />
              </a-form-item>
            </a-col>
            <a-col :span="24">
              <a-form-item field="host" label="服务器地址" :rules="[{required:true,message:'请填服务器地址'}]">
                  <div class="flex wrap-tow">
                  <a-select :style="{width:'120px'}" v-model="formData.protocol">
                      <a-option value="mqtt">mqtt://</a-option>
                      <a-option value="mqtts">mqtts://</a-option>
                      <a-option value="ws">ws://</a-option>
                      <a-option value="wss">wss://</a-option>
                  </a-select>
                  <div class="flex_1 tow-right">
                      <a-input v-model="formData.host" :placeholder="$t('data.form.enter')" :style="{width:'100%'}"/>
                  </div>
                </div>
              </a-form-item>
            </a-col>
            <a-col :span="24">
              <a-form-item field="port" label="端口" validate-trigger="input" :rules="[{required:true,message:'请填写端口'}]">
                <a-input v-model="formData.port" :placeholder="$t('data.form.enter')" />
              </a-form-item>
            </a-col>
            <a-col :span="24">
              <a-form-item field="client_id" label="Client ID" validate-trigger="input" :rules="[{required:true,message:'请填写Client ID'}]">
                <a-input v-model="formData.client_id" :placeholder="$t('data.form.enter')" />
                <a href="javascript:;" class="icon-btn" @click="handleClientID"><icon-refresh /></a>
              </a-form-item>
            </a-col>
            <a-col :span="24" v-if="formData.protocol === 'ws' || formData.protocol === 'wss'">
              <a-form-item field="path" label="Path" validate-trigger="input" :rules="[{required:true,message:'请填写Path'}]">
                <a-input v-model="formData.path" :placeholder="$t('data.form.enter')" />
              </a-form-item>
            </a-col>
            <a-col :span="24">
              <a-form-item field="username" label="用户名" >
                <a-input v-model="formData.username" :placeholder="$t('data.form.enter')" />
              </a-form-item>
            </a-col>
            <a-col :span="24">
              <a-form-item field="password" label="密码" >
                <a-input-password v-model="formData.password" :placeholder="$t('data.form.enter')" />
              </a-form-item>
            </a-col>
            <a-col :span="24">
              <a-form-item field="ssl" label="SSL/TLS" class="gf-bottom">
                <a-switch v-model="formData.ssl" type="round" :checked-value="1" :unchecked-value="0">
                  <template #checked>ON</template>
                  <template #unchecked>OFF</template>
                </a-switch>
              </a-form-item>
            </a-col>
            <template v-if="formData.ssl">
              <a-form-item field="rejectUnauthorized" label="SSL 安全" tooltip="是否验证服务端证书链和地址名称" class="form-item-mb">
                <a-switch v-model="formData.ssl_security":checked-value="1" :unchecked-value="0">
                  <template #checked>ON</template>
                  <template #unchecked>OFF</template>
                </a-switch>
              </a-form-item>
              <a-col :span="24">
                <a-form-item field="alpn" label="ALPN" validate-trigger="input" class="gf-bottom">
                  <a-input v-model="formData.alpn" :placeholder="$t('data.form.enter')" />
                </a-form-item>
              </a-col>
              <a-form-item field="certType" label="证书类型" class="form-item-mb">
                <a-radio-group v-model="formData.cert_type">
                  <a-radio value="server">{{ $t('form.certType.server') }}</a-radio>
                  <a-radio value="self">{{ $t('form.certType.self') }}</a-radio>
                </a-radio-group>
              </a-form-item>
            </template>
          </a-row>
      </a-card>
      <!--自定义证书-->
      <a-card v-if="formData.cert_type=='self'" title="Certificates" class="gf-mx gf-bottom">
          <a-form-item field="ca" label="CA 文件" validate-trigger="input" class="form-item-mb">
            <a-input v-model="formData.ca" :placeholder="$t('data.form.enter')" allow-clear/>
            <a-button @click="handleSelectFile('ca')" :style="{marginLeft:'10px'}">选择文件</a-button>
          </a-form-item>
          <a-form-item field="cert" label="客户端证书cert" validate-trigger="input" class="form-item-mb">
            <a-input v-model="formData.cert" :placeholder="$t('data.form.enter')" allow-clear/>
            <a-button @click="handleSelectFile('cert')" :style="{marginLeft:'10px'}">选择文件</a-button>
          </a-form-item>
          <a-form-item field="key" label="客户端key文件" validate-trigger="input" class="form-item-mb">
            <a-input v-model="formData.key" :placeholder="$t('data.form.enter')" allow-clear/>
            <a-button @click="handleSelectFile('key')" :style="{marginLeft:'10px'}">选择文件</a-button>
          </a-form-item>
      </a-card>

      <a-card title="高级设置" class="gf-mx gf-bottom" :body-style="{padding:advancedVisible?'16px':0}">
          <template #extra>
            <a-link @click="advancedVisible=!advancedVisible">{{ advancedVisible?"收起":"展开" }}<icon-caret-up v-if="advancedVisible"/><icon-caret-down v-if="!advancedVisible"/></a-link>
          </template>
          <div class="wrap-collapse" v-show="advancedVisible">
              <a-form-item field="mqtt_version" label="MQTT 版本" validate-trigger="input" class="form-item-mb">
                  <a-select v-model="formData.mqtt_version">
                  <a-option value="3.1">3.1</a-option>
                  <a-option value="3.1.1">3.1.1</a-option>
                  <a-option value="5.0">5.0</a-option>
              </a-select>
              </a-form-item>
              <a-form-item field="connect_timeout" label="连接超时时长" validate-trigger="input" class="form-item-mb">
                  <a-input-number v-model="formData.connect_timeout" :placeholder="$t('data.form.enter')" allow-clear hide-button>
                      <template #suffix>秒</template>
                  </a-input-number>
              </a-form-item>
              <a-form-item field="keepalive" label="Keep Alive" validate-trigger="input" class="form-item-mb">
                  <a-input-number v-model="formData.keepalive" :placeholder="$t('data.form.enter')" allow-clear hide-button>
                      <template #suffix>秒</template>
                  </a-input-number>
              </a-form-item>
              <a-form-item field="reconnect" label="自动重连" validate-trigger="input" class="form-item-mb">
                  <a-switch v-model="formData.reconnect" type="round" :checked-value="1" :unchecked-value="0">
                      <template #checked>ON</template>
                      <template #unchecked>OFF</template>
                  </a-switch>
              </a-form-item>
              <a-form-item v-if="formData.reconnect" field="reconnect_period" label="重连周期" validate-trigger="input" class="form-item-mb">
                  <a-input-number v-model="formData.reconnect_period" :placeholder="$t('data.form.enter')" allow-clear/>
              </a-form-item>
              <a-form-item field="clean" :label="formData.mqtt_version === '5.0' ? 'Clean Start' : 'Clean Session'" validate-trigger="input" class="form-item-mb">
                  <a-switch v-model="formData.clean" type="round" :checked-value="1" :unchecked-value="0">
                      <template #checked>ON</template>
                      <template #unchecked>OFF</template>
                  </a-switch>
              </a-form-item>
              <template v-if="formData.mqtt_version === '5.0'">
              <a-form-item field="sessionExpiryInterval" label="会话过期时间" validate-trigger="input" class="form-item-mb">
                <a-input-number v-model="formData.properties.sessionExpiryInterval" :placeholder="$t('data.form.enter')" allow-clear hide-button>
                <template #suffix>秒</template>
                </a-input-number>
              </a-form-item>
              <a-form-item field="receiveMaximum" label="接收最大数值" validate-trigger="input" class="form-item-mb">
                  <a-input-number v-model="formData.properties.receiveMaximum" :min="1" :placeholder="$t('data.form.enter')" allow-clear/>
              </a-form-item>
              <a-form-item field="maximumPacketSize" label="最大数据包大小" validate-trigger="input" class="form-item-mb">
                  <a-input-number v-model="formData.properties.maximumPacketSize" :min="100" :placeholder="$t('data.form.enter')" allow-clear/>
              </a-form-item>
              <a-form-item field="topicAliasMaximum" label="主题别名最大值" validate-trigger="input" class="form-item-mb">
                  <a-input-number v-model="formData.properties.topicAliasMaximum" :min="1" :placeholder="$t('data.form.enter')" allow-clear/>
              </a-form-item>
              <a-form-item field="requestResponseInformation" label="请求响应信息" validate-trigger="input" class="form-item-mb">
                  <a-switch v-model="formData.properties.requestResponseInformation" type="round" :checked-value="1" :unchecked-value="0">
                      <template #checked>ON</template>
                      <template #unchecked>OFF</template>
                  </a-switch>
              </a-form-item>
              <a-form-item field="requestProblemInformation" label="请求失败信息" validate-trigger="input" class="form-item-mb">
                  <a-switch v-model="formData.properties.requestProblemInformation" type="round" :checked-value="1" :unchecked-value="0">
                      <template #checked>ON</template>
                      <template #unchecked>OFF</template>
                  </a-switch>
              </a-form-item>
              </template>
          </div>
      </a-card>
      <!--用户属性-->
      <KeyValueEditor :bordered="true" v-if="formData.mqtt_version === '5.0'" v-model="formData.user_properties" :cardMaxHeight="230"/>

      <a-card title="遗嘱消息" class="gf-mx gf-bottom" :body-style="{padding:willVisible?'16px':0}">
          <template #extra>
              <a-link @click="willVisible=!willVisible">{{ willVisible?"收起":"展开" }}<icon-caret-up v-if="willVisible"/><icon-caret-down v-if="!willVisible"/></a-link>
          </template>
            <div class="wrap-collapse" v-show="willVisible">
              <a-form-item field="willTopic" label="遗嘱消息主题" validate-trigger="input" class="form-item-mb">
                  <a-input v-model="formData.will.willTopic" :placeholder="$t('data.form.enter')" />
              </a-form-item>
              <a-form-item field="willQos" label="遗嘱消息 QoS" validate-trigger="change" class="form-item-mb">
                <a-radio-group v-model="formData.will.willQos">
                  <a-radio :value="0">0</a-radio>
                  <a-radio :value="1">1</a-radio>
                  <a-radio :value="2">2</a-radio>
                </a-radio-group>
              </a-form-item>
              <a-form-item field="willRetain" label="遗嘱消息保留标志" validate-trigger="change" class="form-item-mb">
                  <a-switch v-model="formData.will.willRetain" type="round" :checked-value="1" :unchecked-value="0">
                    <template #checked>ON</template>
                    <template #unchecked>OFF</template>
                  </a-switch>
              </a-form-item>
              <a-form-item field="willPayload" label="遗嘱消息" class="form-item-mb">
                <div class="code-editor">
                  <div class="code-editor-content">
                    <CodeEditor
                      v-model:value="formData.will.willPayload"
                      :language="languageType"
                      :theme="theme=='dark'?'vs-dark':'vs'"
                      :options="editorOptions"
                      @editorDidMount="handleEditorDidMount"
                    />
                  </div>
                  <div class="code-editor-lang-type">
                    <a-radio-group v-model="languageType" @change="syncEditorLanguage">
                      <a-radio value="json">JSON</a-radio>
                      <a-radio value="plaintext">Plaintext</a-radio>
                    </a-radio-group>
                  </div>
                </div>
              </a-form-item>
              <!-- MQTT v5.0 -->
              <template v-if="formData.mqtt_version === '5.0'">
                <a-form-item field="payloadFormatIndicator" label="有效载荷指示器" validate-trigger="change" class="form-item-mb">
                    <a-switch v-model="formData.will.payloadFormatIndicator" type="round" :checked-value="1" :unchecked-value="0">
                      <template #checked>ON</template>
                      <template #unchecked>OFF</template>
                    </a-switch>
                </a-form-item>
                <a-form-item field="willDelayInterval" label="遗嘱消息延迟时间" validate-trigger="input" class="form-item-mb">
                  <a-input-number v-model="formData.will.willDelayInterval" :placeholder="$t('data.form.enter')" allow-clear hide-button>
                    <template #suffix>秒</template>
                  </a-input-number>
                </a-form-item>
                <a-form-item field="messageExpiryInterval" label="消息过期时间" validate-trigger="input" class="form-item-mb">
                  <a-input-number v-model="formData.will.messageExpiryInterval" :placeholder="$t('data.form.enter')" allow-clear hide-button>
                    <template #suffix>秒</template>
                  </a-input-number>
                </a-form-item>
                 <a-form-item field="contentType" label="内容类型" validate-trigger="input" class="form-item-mb">
                  <a-input v-model="formData.will.contentType" :placeholder="$t('data.form.enter')" />
                </a-form-item>
                 <a-form-item field="responseTopic" label="响应主题" validate-trigger="input" class="form-item-mb">
                  <a-input v-model="formData.will.responseTopic" :placeholder="$t('data.form.enter')" />
                </a-form-item>
                 <a-form-item field="correlationData" label="对比数据" validate-trigger="input" class="form-item-mb">
                  <a-input v-model="formData.will.correlationData" :placeholder="$t('data.form.enter')" />
                </a-form-item>
              </template>
          </div>
      </a-card>
      <!-- <a-card title="扩展" class="gf-mx">
          
      </a-card> -->
    </a-form>
   </div>
  </a-layout-content>
</template>

<script lang="ts" setup>
  import { ref,computed,onMounted,watch,nextTick } from 'vue';
  import { FormInstance,Message} from '@arco-design/web-vue';
  import useLoading from '@/hooks/loading';
  import {Dialogs} from "@wailsio/runtime";
  import { CodeEditor, type EditorOptions } from 'monaco-editor-vue3';
  import { useAppStore } from '@/store';
  import { getClientId } from '@/utils/generateRandom'
  import { MqttConnectionForm } from '@/types/connect'
  import KeyValueEditor from '../components/KeyValueEditor.vue'
  //go 数据接口
  import {MqttConnectionService} from "/#/gofly/internal/service";
  const appStore = useAppStore();
  const theme = computed(() => {
    return appStore.theme;
  });
  const props = defineProps<{
    id: number
    pid: number
    from_id: number
  }>()
  const emit = defineEmits<{
    'add-connect': [index: number,from_id:number]
    'success': [id:number,type:number]
  }>()
  // 手动切换语法高亮
  const languageType = ref('plaintext')
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

  //表单
  const willVisible = ref(true)
  const advancedVisible = ref(true)
  const formRef = ref<FormInstance>();

  // 默认表单
  const getDefaultForm = (id: number, pid: number): MqttConnectionForm => ({
    id,
    pid,
    title: "",
    protocol: "mqtt",
    host: "gmqt.goflys.cn",
    port: "1883",
    client_id: getClientId(),
    path: "/mqtt",
    username: "",
    password: "",
    ssl: 0,
    ssl_security: 0,
    alpn: "",
    cert_type: "server",
    ca: "",
    cert: "",
    key: "",
    mqtt_version: "5.0",
    connect_timeout: 10,
    keepalive: 10,
    reconnect: 1,
    reconnect_period: 4000,
    clean: 1,
    properties: {
      sessionExpiryInterval: 0,
      receiveMaximum: undefined,
      maximumPacketSize: undefined,
      topicAliasMaximum: undefined,
      requestResponseInformation: 0,
      requestProblemInformation: 0,
    },
    will: {
      willTopic: "",
      willQos: 0,
      willRetain: 0,
      willPayload: "",
      payloadFormatIndicator: 0,
      willDelayInterval: undefined,
      messageExpiryInterval: undefined,
      contentType: "",
      responseTopic: "",
      correlationData: "",
    },
    user_properties: [],
  })

  const formData = ref<MqttConnectionForm>(getDefaultForm(props.id, props.pid))
  const { loading, setLoading } = useLoading();
  //监听props数据
  watch(props, (newVal, oldVal) => {
    formRef.value?.resetFields()
    formData.value.id=newVal.id
    formData.value.pid=newVal.pid
    if(newVal.id>0){
      getInitData(newVal.id)
    }else{
      formData.value=getDefaultForm(newVal.id, newVal.pid)
    }
  }, { deep: true })

  //返回
  const back = () => {
    if(!props.from_id||props.from_id==0){
      emit('add-connect', 0, 0)
    }else{
      emit('add-connect', 1,props.from_id)
    }
  };
  //获取编辑数据
  const getInitData=async(id:number)=>{
    const res =await MqttConnectionService.GetContent(id);
    if(res.code==0){
      // formRef.value?.resetFields()
      const data=res.data
      const properties=data.properties
      const will=data.will
      delete data["properties"]
      delete data["will"]
      Object.assign(formData.value,data)
      if(res.data.properties)
      Object.assign(formData.value,{properties:JSON.parse(properties)})
      if(res.data.will)
       Object.assign(formData.value,{will:JSON.parse(will)})
      if(res.data.user_properties)
      formData.value.user_properties=JSON.parse(data.user_properties)
      console.log("获取编辑数据",formData.value)
    }
  }
  //初始数据
  onMounted(async ()=>{
    await nextTick()
    if(props.id>0){
      getInitData(props.id)
    }
  })
  //提交数据
  const handleSubmit =async (type:number) => {
    try {
        const validate = await formRef.value?.validate();
        if (!validate) {
          setLoading(true);
          Message.loading({content:"提交中",id:"submit",duration:0})
          const res =await MqttConnectionService.SaveData(JSON.stringify(formData.value));
          if(res.code==0){
            emit('success',res.data,type);
            Message.success({content:res.message,id:"submit",duration:2000})
          }else{
            Message.error({content:res.message,id:"submit",duration:2000})
          }
          setLoading(false);
        }
    } catch (error) {
        setLoading(false);
        Message.loading({content:"提交中",id:"submit",duration:1})
    }
  };
  //选择文件
  const handleSelectFile = async (key: 'ca' | 'cert' | 'key') => {
      const filedata = await Dialogs.OpenFile({Title:"选择证书文件"})
      if (filedata) {
          formData.value[key] = filedata
          console.log('选择附件结果：:', filedata)
      }
  }
  //生成id字符串
  const handleClientID=async()=>{
    formData.value.client_id = getClientId()
  }
  //code编辑器参数
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
</script>

<style scoped lang="less">
  @header-height: 50px;
  //头部导航
  .connect-nav {
    width: 100%;
    height: @header-height;
    box-shadow: 0 1px 6px 0 rgba(0, 0, 0, 0.04);
    font-size: 14px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 16px;
    user-select: none;
    background: var(--color-bg-3);
    border-bottom: 1px var(--color-fill-1) solid;
    .connect-nav-left {
        flex: 1;
        display: flex;
        align-items: center;
        height: 100%;
        .back {
            width: 32px;
            height: 32px;
            display: flex;
            align-items: center;
            justify-content: center;
            cursor: pointer;
            color: var(--color-text-2);
            background-color: var(--color-secondary);
            font-size: 14px;
            border-radius: var(--border-radius-small);
            &:hover {
                background-color: var(--color-secondary-hover);
            }
        }
        .name-desc {
            display: flex;
            flex-direction: column;
            max-width: 200px;
            margin-left: 15px;
            .name {
              margin-bottom: 2px;
              font-size: 16px;
              white-space: nowrap;
              text-overflow: ellipsis;
              overflow: hidden;
            }
            .desc {
              font-size: 12px;
              color: #8f959e;
              white-space: nowrap;
              text-overflow: ellipsis;
              overflow: hidden;
            }
        }
    }
    .connect-nav-mid {
        flex: 2;
        text-align: center;
    }
    .connect-nav-right {
        flex: 1;
        display: flex;
        align-items: center;
        justify-content: flex-end;
    }
   }
   //内容
   .connect-content{
     padding: 10px;
     height:calc(100% - @header-height);
     overflow-y: auto;
   }
   .gf-mx{
     border-radius: 2px;
   }
   .gf-bottom{
    margin-bottom: 10px;
   }
   .wrap-tow{
    width: 100%;
    .tow-right{
        padding-left: 10px;
    }
   }
   .form-item-mb{
     margin-bottom: 10px;
     &:last-child{
         margin-bottom: 0px;
     }
   }
   //代码编辑器
   .code-editor{
       width: 100%;
      .code-editor-content{
        height: 235px;
        width: 100%;
        border: 1px solid var(--color-neutral-4);
        padding: 10px 1px 1px 1px;
        border-top-left-radius: 4px;
        border-top-right-radius: 4px;
        overflow: hidden;
      }
      .code-editor-lang-type{
        width: 100%;
        height: 30px;
        line-height: 30px;
        padding: 0px 12px;
        background: var(--color-neutral-2);
        border: 1px solid var(--color-neutral-4);
        border-top: none;
        text-align: right;
        border-bottom-left-radius: 4px;
        border-bottom-right-radius: 4px;
      }
   }
</style>