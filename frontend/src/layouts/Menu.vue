<template>
  <a-layout id="app-menu">
    <a-layout-sider
      theme="light"
      class="layout-sider"
      style="width: 120px;"
    >
      <a-menu 
        theme="light" 
        mode="inline" 
        :selectedKeys="[current]"
        @menu-item-click="changeMenu">
        <a-menu-item v-for="item in menulist" :key="item.name">
          <div class="flex-all-center">
             <icon-font :name="item?.meta.icon" /> <span class="menu-tile">{{ item?.meta.title }}</span>
          </div>
        </a-menu-item>
      </a-menu>
    </a-layout-sider>
    <a-layout>
      <a-layout-content>
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>
<script setup lang="ts">
  import { ref,computed } from 'vue';
  import { useRouter,useRoute } from 'vue-router';
  import routerMap from '@/router/routerMap';
  const route = useRoute();
  const router = useRouter();
  const current=ref(route.name)
  const menulist = computed(() => {
    const childrenlist= routerMap.find((item)=>item.path=="/")?.children
    const nowname=route.path.split("/")
    return childrenlist?.find((item)=>item.name==nowname[1])?.children
  });
  const changeMenu=(key)=>{
    current.value = key;
    const linkInfo =menulist.value?.find((item)=>item.name==key)
    if(linkInfo)
    router.push({ name: linkInfo.name, params: linkInfo.params})
  }
</script>
<style lang="less" scoped>
#app-menu {
  height: 100%;
  text-align: center;
  .layout-sider {
    // border-top: 1px solid var(--color-neutral-3);
    border-right: 1px solid var(--color-neutral-3);
    // background-color: #FAFAFA;
    overflow: auto;
    :deep(.arco-menu .arco-menu-item .arco-icon){
      margin-right: 5px;
      display: inline-block;
    }
  }
}
.menu-tile{
  padding-left: 3px;
}
</style>
