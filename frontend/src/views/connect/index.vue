<template>
  <a-layout id="app-menu">
    <a-layout-sider
      theme="light"
      class="layout-sider"
      collapsible 
      hide-trigger
      :collapsed="collapsed"
      :width="250"
      :collapsed-width="49"
    >
      <div class="flex flex-middle menu-header">
        <transition name="menu-title-fade" >
          <div class="flex_1 title" v-if="!collapsed">{{ $t('connect.title') }}</div>
        </transition>
        <div class="flex-all-center">
          <a-space>
            <transition name="menu-title-fade" >
              <a-button @click="hanldeClearData" v-if="!collapsed&&(menuData&&menuData.length>0)" :style="{ padding: '0 6px', height: '25px', lineHeight: '25px'}" >
                <icon-eraser />
              </a-button>
            </transition>
            <a-dropdown position="br">
              <transition name="menu-title-fade" >
              <a-button v-if="!collapsed" :style="{ padding: '0 6px', height: '25px', lineHeight: '25px'}" >
                <icon-plus />
              </a-button>
              </transition>
              <template #content>
                <a-doption @click="handleAddConnect(0)">
                  <template #icon><icon-font name="icon-connect5" size="13"/></template>
                  <template #default>新建连接</template>
                </a-doption>
                <a-doption @click="handleAddGroup">
                  <template #icon><icon-folder :size="15"/></template>
                  <template #default>新建分组</template>
                </a-doption>
              </template>
            </a-dropdown>
            <a-button :style="{ padding: '0 6px', height: '25px', lineHeight: '25px'}" @click="()=>collapsed=!collapsed">
              <icon-menu-unfold v-if="collapsed" />
              <icon-menu-fold v-else />
            </a-button>
          </a-space>
        </div>
      </div>
      <a-menu
        ref="menuRef"
        :key="menuVersion"
        :collapsed="collapsed"
        v-model:open-keys="openKeys"
        v-model:selected-keys="selectedKeys"
        :style="{ width: `100%`, height: 'calc(100% - 37px)' }"
        @menu-item-click="handleOpenConnect"
      >
        <div class="sort-group" data-group="0">
          <RecursiveMenuItem 
            :list="menuData" 
            :openKeys="openKeys"
            :menu-type="menuType"
            @group-event="onGroupEvent"
            @connect-event="onConnectEvent"
          />
        </div>
        <div v-if="!menuData||menuData.length==0" style="margin-top: 35%;">
          <a-empty  description="请添加连接/分组"/>
        </div>
      </a-menu>
    </a-layout-sider>
    <a-layout class="layout-wrap">
      <component 
        :is="PageMap[activeKey]" 
        :id="openId" 
        :pid="AddPId" 
        :start="isStart" 
        :from_id="fromId" 
        @add-connect="onChangePage" 
        @success="onAddForm"
      />
    </a-layout>
  </a-layout>
  <!--添加新平台-->
  <a-modal v-model:visible="groupModal.visible" simple hide-title :mask-closable="false" @cancel="groupModal.visible = false;" @before-ok="handleBeforeOk">
    <a-form ref="formAddRef" :model="groupModal.data" layout="vertical" size="large">
      <a-form-item field="title" label="分组名称" :rules="[{required:true,message:'请填写分组名称'}]" style="margin-bottom: 0px;">
        <a-input v-model="groupModal.data.title" :style="{width:'320px'}" placeholder="请输入" allow-clear @press-enter="handleBeforeOk"/>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onBeforeUnmount, nextTick } from 'vue'
import Sortable from 'sortablejs'
import RecursiveMenuItem, { type MenuItem } from './components/RecursiveMenuItem.vue'
import { FormInstance,Modal,Message} from '@arco-design/web-vue';
//go 数据接口
import {MqttConnectionService} from "/#/gofly/internal/service";
//页面
import DefaultPage from './page/DefaultPage.vue'
import Operation from './page/Operation.vue'
import AddForm from './page/AddForm.vue'
import { useMqttStore } from '@/store'
const mqttStore = useMqttStore()
//编辑、操作页面
const PageMap: any= {
  0: DefaultPage,
  1: Operation,
  2: AddForm
}
const activeKey = ref(0)
const openId = ref(0)//连接数据id
const AddPId = ref(0)//父级id
const fromId = ref(0)//从上个连接返回数据id，用于返回
const isStart = ref(false)// 打开操作页面是否启动连接

const menuType = ref<number[]>([])      // 唯一一份状态，父组件维护
const collapsed = ref(false)            //是否折叠菜单
const openKeys = ref<number[]>([])      //展开的子菜单 key 数组
const selectedKeys = ref<number[]>([])  //选中的菜单项 key 数组
//分组
const formAddRef = ref<FormInstance>();
const groupModal=ref({
  visible:false,
  data:{id:0,pid:0,title:""}
})
//添加连接
const handleAddConnect=(pid:number)=>{
  activeKey.value=2
  AddPId.value=pid
  openId.value=0
  fromId.value=0
}
//切换页面
const onChangePage=(index:number,from_id:number)=>{
  isStart.value=false
  if(index==1){
    openId.value=from_id
    fromId.value=0
  }else{
    fromId.value=from_id
  }
  activeKey.value=index
  getData()
}
//打开连接操作
const handleOpenConnect=(key:any)=>{
  activeKey.value=1
  openId.value=key
  isStart.value=false
  mqttStore.setCurrentConn(Number(key))   // 清零未读 + 记录当前连接
}
//添加分组
const handleAddGroup=()=>{
  groupModal.value.visible=true
  groupModal.value.data={id:0,pid:0,title:""}
}
//提交分组
const handleBeforeOk = async(done) => {
  try {
      const validate = await formAddRef.value?.validate();
      if (!validate) {
        Message.loading({content:"提交中",id:"save",duration:0})
        const res =await MqttConnectionService.SaveGroup(JSON.stringify({...groupModal.value.data,is_group:1}));
          if(res.code==0){
            groupModal.value.visible=false
            getData()
            done()
            Message.success({content:res.message,id:"save",duration:2000})
          }else{
            Message.error({content:res.message,id:"save",duration:2000})
          }
      }else{
        done(false)
      }
    } catch (error) {
      Message.loading({content:"提交中",id:"save",duration:2})
    }
};
//提交连接成功返回
const onAddForm=(id:number,type:number)=>{
  getData()
  console.log("提交连接成功返回:",id)
  openId.value=id
  selectedKeys.value=[id]//默认选中
  if(type==1){
    activeKey.value=1
    isStart.value=true
  }else{
    isStart.value=false
  }
}
//拖拽排序
type SortableInstance = InstanceType<typeof Sortable>
let sortableList: SortableInstance[] = []
const menuData = ref<MenuItem[]>([])
const menuVersion = ref(0)   // 强制重建菜单的版本号
/**
 * 递归删除节点，返回被移除节点
 */
function removeChildByKey(tree: MenuItem[], targetKey: number): MenuItem | null {
  targetKey = Number(targetKey)
  // 顶层直接节点（pid=0 的连接/分组）
  const topIdx = tree.findIndex(c => Number(c.id) === targetKey)
  if (topIdx > -1) return tree.splice(topIdx, 1)[0]
  // 递归 children
  for (const node of tree) {
    if (node.children?.length) {
      const idx = node.children.findIndex(c => Number(c.id) === targetKey)
      if (idx > -1) return node.children.splice(idx, 1)[0]
      const found = removeChildByKey(node.children, targetKey)
      if (found) return found
    }
  }
  return null
}
// 在整个树（含嵌套子分组）中按 id 找节点
function findNodeById(tree: MenuItem[], id: number): MenuItem | null {
  for (const node of tree) {
    if (Number(node.id) === Number(id)) return node
    if (node.children?.length) {
      const found = findNodeById(node.children, id)
      if (found) return found
    }
  }
  return null
}
function handleDragEnd(fromGroup: number, toGroup: number, itemKey: number, newIndex: number) {
  const childItem = removeChildByKey(menuData.value, itemKey)
  if (!childItem) {
    console.warn('未找到拖拽节点', itemKey)
    return
  }
  if (toGroup === 0) {
    // 拖到根级
    const dup = menuData.value.findIndex(c => Number(c.id) === Number(childItem.id))
    if (dup > -1) menuData.value.splice(dup, 1)
    menuData.value.splice(newIndex, 0, childItem)
  } else {
    const targetSub = findNodeById(menuData.value, toGroup)
    if (!targetSub) return
    if (!targetSub.children) targetSub.children = []
    const dup = targetSub.children.findIndex(c => Number(c.id) === Number(childItem.id))
    if (dup > -1) targetSub.children.splice(dup, 1)
    targetSub.children.splice(newIndex, 0, childItem)
  }
}

// 拖拽结束后把整棵树扁平化，传给后端更新 pid + weigh
const saveOrder = async () => {
  const list: { id: number; pid: number; weigh: number }[] = []
  const flatten = (tree: MenuItem[], pid: number) => {
    tree.forEach((n, i) => {
      list.push({ id: n.id, pid, weigh: i + 1 })
      if (n.children?.length) flatten(n.children, n.id)
    })
  }
  flatten(menuData.value, 0)
  const res = await MqttConnectionService.UpdateOrder(JSON.stringify(list))
  if (res.code !== 0) {
    Message.error({ content: res.message, id: 'order', duration: 2000 })
    getData()   // 保存失败回滚
  }
}

async function initSortable() {
  sortableList.forEach(s => s.destroy())
  sortableList = []
  await nextTick()
  const groups = document.querySelectorAll('.sort-group') as NodeListOf<HTMLElement>
  groups.forEach((el: HTMLElement) => {
    const sort = Sortable.create(el, {
      group: {
        name: 'menu-item-transfer',
        pull: true,
        put: true,
      },
      animation: 150,
      ghostClass: 'sortable-ghost',
      filter: '.arco-menu-item-disabled',
      draggable: '.arco-menu-item',
      onEnd: (evt) => {
       const itemEl = evt.item as HTMLElement
          const fromGroup = parseInt(evt.from.getAttribute('data-group') || '0', 10)
          const toGroup   = parseInt(evt.to.getAttribute('data-group')   || '0', 10)
          const itemKey   = parseInt(itemEl.getAttribute('data-key') || '0', 10)
          const newIndex  = evt.newIndex ?? 0
          if (!itemKey || !Number.isInteger(fromGroup) || !Number.isInteger(toGroup)) return
          handleDragEnd(fromGroup, toGroup, itemKey, newIndex)
          menuVersion.value++
          nextTick(() => initSortable())
          saveOrder()
      },
    })
    sortableList.push(sort)
  })
}

// sub‑menu展开关闭DOM重建，重新初始化拖拽
watch(openKeys, async () => {
  await nextTick()
  initSortable()
})

onMounted(async() => {
  await getData()
  initSortable()
})
//获取数据列表
const getData=async()=>{
   const res= await MqttConnectionService.GetList()
  if(res.code==0){
    menuData.value=res.data
  }
}
//删除数据
const delData=async(id:number,title:string)=>{
  Modal.warning({
    title: `是否确认删除${title}？`,
    content: '删除后将无法恢复，需要谨慎操作哦！',
    hideCancel:false,
    titleAlign:"start",
    onOk:async()=>{
      try {
          Message.loading({content:"删除中",id:"del",duration:0})
          const res =await MqttConnectionService.Del(id);
          if(res.code==0){
            Message.success({content:res.message,id:"del",duration:2000})
            getData()
            activeKey.value=0
            AddPId.value=0
            openId.value=0
            fromId.value=0
          }else{
            Message.error({content:res.message,id:"del",duration:2000})
          }
      } catch (error) {
        Message.loading({content:"删除中",id:"del",duration:2})
      }
    }
  });
}

//复制数据
const CopyData=async(id:number)=>{
  try {
    Message.loading({content:"复制数据中",id:"copy",duration:0})
    const res =await MqttConnectionService.Copy(id);
    if(res.code==0){
      Message.success({content:res.message,id:"copy",duration:2000})
      getData()
    }else{
      Message.error({content:res.message,id:"copy",duration:2000})
    }
  } catch (error) {
    Message.loading({content:"复制数据中",id:"copy",duration:2})
  }
}
onBeforeUnmount(() => {
  sortableList.forEach(s => s.destroy())
  sortableList = []
})
//分组数据操作
const onGroupEvent=(type: string,item:MenuItem)=>{
  if(type=="create"){
    groupModal.value.visible=true
    groupModal.value.data={id:0,pid:item.id,title:""}
  }else if(type=="edit"){
    groupModal.value.visible=true
    groupModal.value.data={id:item.id,pid:item.pid,title:item.title}
  }else if(type=="delete"){
    delData(item.id,item.title)
  }else if(type=="addconnect"){
    handleAddConnect(item.id)
  }
}
//连接数据操作
const onConnectEvent=(type: string,item:MenuItem)=>{
  if(type=="copy"){
    CopyData(item.id)
  }else if(type=="edit"){
    activeKey.value=2
    openId.value=item.id
    AddPId.value=item.pid
    selectedKeys.value=[item.id]//默认选中
  }else if(type=="delete"){
    delData(item.id,item.title)
  }
}
//清除全部数据
const hanldeClearData=async()=>{
  Modal.warning({
    title: `是否确认清除全部数据？`,
    content: '删除后将无法恢复，需要谨慎操作哦！',
    hideCancel:false,
    titleAlign:"start",
    onOk:async()=>{
      try {
          Message.loading({content:"删除中",id:"del",duration:0})
          const res =await MqttConnectionService.ClearData();
          if(res.code==0){
            Message.success({content:res.message,id:"del",duration:2000})
            getData()
            activeKey.value=0
            AddPId.value=0
            openId.value=0
            fromId.value=0
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

<style lang="less" scoped>
#app-menu {
  height: 100vh;
  overflow: hidden;
}
:deep(.arco-layout-sider-children) {
  overflow: hidden;
}
:deep(.sort-group) {
  min-height: 8px;
}
:deep(.sortable-ghost) {
  background-color: var(--color-neutral-3) !important;
  opacity: 0.75;
}
:deep(.sortable-drag) {
  background-color: var(--color-neutral-4) !important;  
  border-radius: 4px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  transition: none !important;  
  opacity: 1 !important;
  z-index: 9999;
}
:deep(.sortable-chosen) {
  opacity: 0.4;
}
.menu-header{
  padding: 5px 10px;
  width: 100%;
  overflow: hidden;
  .title{
    font-size: 15px;
    font-weight: 600;
  }
}
.layout-wrap{
  height: 100vh;
  overflow: hidden;
  // background-color: var(--color-bg-3);
  padding-left: 1px;
}
.layout-content{
  height: 100%;
  background-color: var(--color-fill-1);
}
</style>
