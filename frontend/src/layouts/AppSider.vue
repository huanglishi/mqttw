<template>
  <a-layout id="app-layout-sider" >
    <a-layout-sider
      v-model="collapsed"
      theme="light"
      class="layout-sider"
      style="width: auto;min-width: 70px;"
    >
    <div class="sider-menu">
      <div class="user">
        <a-avatar style="padding: 1px;cursor: pointer;" @click="handleOpenMQTTW">
          <img src="@/assets/logo.png">
        </a-avatar>
      </div>
      <a-menu 
        class="menu-item" 
        theme="light" 
        mode="inline"
        :selectedKeys ="[current]"
        @menu-item-click="menuHandle"
      >
        <template v-for="menuInfo in menulist" >
          <a-menu-item :key="menuInfo.name" v-if="!menuInfo.meta.hideInMenu">
            <div class="menu-btn">
              <div class="menu-icon"><icon-font :name="current==menuInfo.name?menuInfo.meta.selecteicon:menuInfo.meta.icon" :size="19"/></div>
              <div class="menu-text">{{ menuInfo.meta.locale?$t(menuInfo.meta.locale):menuInfo.meta.title }}</div>
              <a-tag v-if="!collapsed && menuInfo.name==='gateway'" size="small" color="orange" class="menu-tag">实验性</a-tag>
            </div>
          </a-menu-item>
        </template>
      </a-menu>
      <div class="footer" :class="{factive:current=='setting'}" @click="handleSetting">
        <div class="menu-btn">
          <div class="menu-icon"><icon-font name="icon-xitong" class="icon" :size="19"/></div>
          <div class="menu-text">{{ $t("menu.setting") }}</div>
        </div>
      </div>
    </div>
    </a-layout-sider>
    <a-layout>
      <a-layout-content class="layout-content">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>
<script setup lang="ts">
  import { ref,computed } from 'vue';
  import { useRouter,useRoute } from 'vue-router';
  import routerMap from '@/router/routerMap';
  import { Browser } from '@wailsio/runtime';
  const menulist = computed(() => routerMap.find((item)=>item.path=="/")?.children);
  const router = useRouter();
  const route = useRoute();
  const collapsed=ref(true)
  const current=ref(route.name)
  //切换菜单
  const menuHandle=(key)=> {
    current.value = key ? key: current.value;
    const linkInfo =menulist.value?.find((item)=>item.name==key)
    if(linkInfo)
    router.push({ name: linkInfo.name, params: linkInfo.params})
  }
  //点击设置
  const handleSetting=()=>{
    current.value ="setting"
    router.push({ name: "setting"})
  }
  //打开MQTTW网站介绍
  const handleOpenMQTTW=()=>{
    Browser.OpenURL("https://goflys.cn/mqttw")
  }
</script>
<style lang="less" scoped>
// 嵌套
#app-layout-sider {
   height: 100%;
  .sider-menu{//菜单
    padding: 5px 0px;
    height: 100%;
    display: flex;
    flex-flow: column;
    .user {//头像
      text-align: center;
      padding: 8px;
      margin-bottom:3px;
      // border-bottom: 1px solid var(--color-neutral-3);
      .name{
        padding-top: 5px;
        font-size: 12;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
    .menu-item{
      flex: 1;
    }
    .footer{
      text-align: center;
      padding: 8px 0px;
      display: flex;
      align-items:center;
      justify-content:center;
      cursor: pointer;
      .menu-btn{
        .menu-icon{
         padding-bottom: 3px;
        }
      }
      &:hover{
        background-color: var(--color-fill-2);
      }
      &.factive{
        background-color: var(--color-fill-2);
      }
    }
  }
  .layout-sider {
    // border-top: 1px solid var(--color-neutral-3);
    border-right: 1px solid var(--color-neutral-3);
  }
  .menu-item {
    .arco-menu-item  {
      padding: 8px 8px !important;
      line-height: unset;
      overflow: hidden;
      display: flex;
      align-items:center;
      justify-content:center;
      .menu-btn{
        .menu-icon{
         padding-bottom: 3px;
         text-align: center;
         .arco-icon{
          margin-right: 0px;
         }
        }
        .menu-text{
           white-space: normal;
        }
        .menu-tag{
          margin-left: 4px;
          flex-shrink: 0;
          line-height: 1;
        }
      }
    }
    :deep(.arco-menu-inner){
      padding: 0;
    }
  }
  .layout-content {
    background: var(--color-menu-light-bg);
  }
}
</style>
