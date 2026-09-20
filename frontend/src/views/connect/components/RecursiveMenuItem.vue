<template>
  <template v-for="item in list" :key="item.id">
    <a-dropdown trigger="contextMenu" alignPoint :popup-max-height="false" :style="{display:'block'}" v-if="item.is_group==1" > 
      <a-sub-menu :key="item.id" @contextmenu.stop class="sort-group" :data-group="item.id">
        <template #icon>
          <Icon :name="openKeys.includes(item.id)?'icon-file-open':'icon-file-close'" :size="16" />
        </template>
        <template #title>
          {{ item.title }}
          <a-badge
              v-if="getGroupUnread(item) > 0"
              :count="getGroupUnread(item)"
              class="badge-count"
            />
        </template>
        <div class="sort-group" :data-group="item.id">
          <RecursiveMenuItem 
          v-if="item.children&&item.children.length>0"
          :list="item.children" 
          :openKeys="openKeys" 
          :menu-type="menuType"
          @group-event="handleChildGroupEvent"
          @connect-event="handleChildConnectEvent"
          />
        </div>
      </a-sub-menu>
      <template #content>
        <a-doption @click="hanldeGroup('create',item)">
          <template #icon><icon-plus /></template>
          <template #default>新建分组</template>
        </a-doption>
        <a-doption @click="hanldeGroup('edit',item)">
          <template #icon><icon-edit /></template>
          <template #default>重命名分组</template>
        </a-doption>
        <a-doption class="doption-delete" @click="hanldeGroup('delete',item)">
          <template #icon><icon-delete /></template>
          <template #default>删除分组</template>
        </a-doption>
        <a-doption @click="hanldeGroup('addconnect',item)">
          <template #icon><icon-font name="icon-connect5" size="13"/></template>
          <template #default>新建连接</template>
        </a-doption>
        <!-- <a-doption>
          <template #icon><icon-font name="icon-connect-batch3" size="13"/></template>
          <template #default>批量加连接</template>
        </a-doption>
        <a-doption>
          <template #icon><icon-play-arrow /></template>
          <template #default>批量启动</template>
        </a-doption>
        <a-doption>
          <template #icon><icon-font name="icon-a-ziyuan5" size="13"/></template>
          <template #default>批量断开</template>
        </a-doption> -->
      </template>
    </a-dropdown>
    <a-dropdown trigger="contextMenu" alignPoint :popup-max-height="false" :style="{display:'block'}"  v-else > 
    <a-menu-item
      :key="item.id"
      :data-key="item.id"
      @contextmenu.stop
    >
      <template #icon>
        <Icon name="icon-connect5" :size="15" />
      </template>
      {{ item.title }}
      <a-badge v-if="mqttStore.msgCount[item.id]" :count="mqttStore.msgCount[item.id]" class="badge-count"/>
    </a-menu-item>
    <template #content>
      <a-doption @click="hanldeConnect('copy',item)">
        <template #icon><icon-font name="icon-connect-batch1" size="13"/></template>
        <template #default>复制连接</template>
      </a-doption>
      <a-doption @click="hanldeConnect('edit',item)">
        <template #icon><icon-edit /></template>
        <template #default>编 辑</template>
      </a-doption>
      <a-doption class="doption-delete" @click="hanldeConnect('delete',item)">
        <template #icon><icon-delete /></template>
        <template #default>删 除</template>
      </a-doption>
    </template>
    </a-dropdown>
  </template>
</template>
<script setup lang="ts">
import { Icon } from '@/components/Icon'
import { useMqttStore } from '@/store'
const mqttStore = useMqttStore()
defineOptions({ name: 'RecursiveMenuItem' })

export interface MenuItem {
  id: number
  pid: number
  is_group: number
  title: string
  icon?: string
  weigh?: number
  children?: MenuItem[]
}

defineProps<{
  list: MenuItem[] ,
  openKeys: number[],
  menuType: number[], // 从父组件传入数组，不要内部ref
}>()

const emit = defineEmits<{
  'group-event': [type: string,item:MenuItem]
  'connect-event': [type: string,item:MenuItem]
}>()

//操作分组
const hanldeGroup = (type: string,item:MenuItem) => {
  emit('group-event', type,item)
}
//操作连接
const hanldeConnect = (type: string,item:MenuItem) => {
  emit('connect-event', type,item)
}
// 子层分组操作透传
const handleChildGroupEvent = (type: string, item: MenuItem) => {
  emit('group-event', type, item)
}
// 子层连接操作透传
const handleChildConnectEvent = (type: string, item: MenuItem) => {
  emit('connect-event', type, item)
}
// 递归统计分组下所有连接（含子分组）的未读总数
const getGroupUnread = (item: MenuItem): number => {
  let total = 0
  const walk = (nodes: MenuItem[] | undefined) => {
    if (!nodes) return
    for (const node of nodes) {
      if (node.is_group == 1) {
        walk(node.children)              // 子分组递归
      } else {
        total += mqttStore.msgCount[node.id] || 0   // 连接直接累加
      }
    }
  }
  walk(item.children)
  return total
}

</script>
<style>
.sort-group {
  min-height: 8px;
}
.arco-menu-item,
.arco-menu-inline-header {
  -webkit-user-select: none;
  user-select: none;
}
.arco-menu-icon{
  margin-right: 8px !important;
}
.badge-count{
  position: absolute;
  top: 10px;
  right: 1px;
}
</style>
