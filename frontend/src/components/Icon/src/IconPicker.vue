<template>
  <div class="icon-container" :class="{ 'is-list': !isGridView }">
    <a-row>
      <section style="flex: 1; margin-right: 8px">
        <a-input v-model="searchValue" placeholder="搜索图标名称" allow-clear size="small" @input="search"
          @clear="search">
          <template #prefix>
            <icon-search />
          </template>
        </a-input>
      </section>

      <a-button size="small" @click="isGridView = !isGridView">
        <template #icon>
          <icon-apps v-if="isGridView" />
          <icon-unordered-list v-else />
        </template>
      </a-button>
    </a-row>

    <section class="icon-list">
      <a-row wrap :gutter="4">
        <a-col v-for="item of currentPageIconList" :key="item" :span="isGridView ? 4 : 8">
          <div class="icon-item" :class="{ active: modelValue === item }" @click="handleSelectedIcon(item)">
            <SvgIcon :size="18" :name="item"  />
            <div class="gi_line_1 icon-name">{{ item }}</div>
          </div>
        </a-col>
      </a-row>
    </section>
    <a-row justify="center" align="center" v-if="total>=42">
      <a-pagination size="mini" :page-size="pageSize" :total="total" :show-size-changer="false"
        @change="onPageChange"></a-pagination>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import {ref,onMounted} from 'vue';
import * as ArcoIcons from '@arco-design/web-vue/es/icon'
import SvgIcon from './SvgIcon.vue';

defineOptions({ name: 'IconSelector' })
interface Props {
  type?: 'arco' 
  modelValue?: string
  enableCopy?: boolean
}
const props = withDefaults(defineProps<Props>(), {
  type: 'arco', // 默认是arco图标类型 
  modelValue: '',
  enableCopy: false
})

const emit = defineEmits(['select', 'update:modelValue'])
const searchValue = ref('') // 搜索词
// 获取 Arco内置图标
let IconList: string[] =[];
 //获取阿里图标
 const getAlilist=async()=> {
  let result: string[] = [];
  const baseUrl = import.meta.env.BASE_URL === '/' ? '' : import.meta.env.BASE_URL;
  const aliIcon = await fetch(`${baseUrl}iconfont/iconfont.json`).then(res => res.json())
  const aliprefix: string = aliIcon?.css_prefix_text ?? '';
  if (aliprefix) {
    result = (aliIcon?.glyphs ?? []).map((item:any) => aliprefix+item.font_class);
  } else if (Array.isArray(aliIcon?.glyphs )) {
    result = (aliIcon?.glyphs ?? []).map((item:any) => item.font_class);
  }
  return result;
}
  const pageSize = 42
  const current = ref(1)
  const total = ref(IconList.length) // 图标总数
  const isGridView = ref(true)// 图标列表grid图视
  onMounted(async() => {
    const alilist = await getAlilist();
     alilist.forEach((icon:string)=>{
      IconList.push(icon)
    })
    total.value=IconList.length
  });
// 当前页的图标列表
const currentPageIconList = ref(IconList.slice(0, pageSize))
// 搜索列表
const searchList = ref<string[]>([])

// 页码改变
const onPageChange = (page: number) => {
  current.value = page
  if (!searchList.value.length) {
    currentPageIconList.value = IconList.slice((page - 1) * pageSize, page * pageSize)
  } else {
    currentPageIconList.value = searchList.value.slice((page - 1) * pageSize, page * pageSize)
  }
}

// 搜索
const search = () => {
  if (searchValue.value) {
    const temp = searchValue.value.toLowerCase()
    searchList.value = IconList.filter((item) => {
      return item.indexOf(temp)>-1
    })
    total.value = searchList.value.length
    currentPageIconList.value = searchList.value.slice(0, pageSize)
  } else {
    searchList.value = []
    total.value = IconList.length
    currentPageIconList.value = IconList.slice((current.value - 1) * pageSize, current.value * pageSize)
  }
}

// 点击选择图标
const handleSelectedIcon = async (icon: string) => {
  emit('select', icon)
  emit('update:modelValue', icon)
}
</script>

<style lang="less" scoped>
.icon-container {
  width: 320px;
  overflow: hidden;
  .icon-list {
    margin-top: 10px;
    margin-bottom: 10px;
    .icon-item {
      height: 30px;
      margin-bottom: 4px;
      overflow: hidden;
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      cursor: pointer;
      border: 1px dashed var(--color-bg-1);

      .icon-name {
        display: none;
      }

      &.active {
        border: 1px dashed rgb(var(--primary-3));
        background-color: rgba(var(--primary-6), 0.05);
      }

      &:not(.active) {
        &:hover {
          border-color: var(--color-border-3);
        }
      }
    }
  }
}

.is-list {
  min-width: 400px;

  .icon-list {
    height: 300px;
    overflow: hidden;
    overflow-y: auto;

    .icon-item {
      display: flex;
      flex-direction: row;
      justify-content: flex-start;
      align-items: center;
      padding-left: 4px;
      box-sizing: border-box;

      .icon-name {
        margin-left: 6px;
        font-size: 12px;
        color: var(--color-text-2);
        display: block;
      }
    }
  }
}
</style>
