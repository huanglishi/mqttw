<template>
  <a-card title="用户属性" :bordered="bordered" class="gf-mx gf-bottom" :body-style="{padding:dataList&&dataList.length>0?'16px':0}">
    <template #extra>
      <a-link @click="handleAddItem"><icon-plus /></a-link>
    </template>
    <div class="key-value-row" :style="{maxHeight:cardMaxHeight+'px'}">
      <div class="key-value-item flex-align" v-for="(item, index) in dataList" :key="index">
        <div class="item-key flex_1">
          <a-input v-model="item.key" placeholder="Key" allow-clear @input="upModelDate"/>
        </div>
        <div class="item-value flex_1">
          <a-input v-model="item.value" placeholder="Value" allow-clear @input="upModelDate"/>
        </div>
        <div class="item-del ">
          <icon-delete :size="20" v-if="!item.key&&!item.value" @click="handleDelItem(item)"/>
          <a-popconfirm v-else content="您确定删除吗?" @ok="handleDelItem(item)">
            <icon-delete :size="20"/>
          </a-popconfirm>
        </div>
      </div>
    </div>
  </a-card>
</template>

<script setup lang="ts">
  import { ref } from 'vue'
  import { UserProperty } from '@/types/global'
  // v-model 接收 value，触发 update:modelValue
  // 允许 undefined/null：调用方字段可能为可选（如 PushPropertiesModel.user_properties）
  const props = defineProps<{
    modelValue?: UserProperty[] | null
    bordered?: boolean
    cardMaxHeight?: number
  }>()
  const emit = defineEmits<{
    'update:modelValue': [val: UserProperty[]]
  }>()
  const dataList = ref<UserProperty[]>(props.modelValue && props.modelValue.length > 0 ? props.modelValue : [{ key: '', value: '' }])
  //添加key-value
  const handleAddItem=()=>{
    dataList.value.push({ key: '', value: '' })
    upModelDate()
  }
  const handleDelItem=(item:UserProperty)=>{
    if(dataList.value){
      dataList.value = dataList.value.filter(data => data !== item)
      upModelDate()
    }
  }
  //更新model值
  const upModelDate=()=>{
    emit('update:modelValue',dataList.value)
  }
</script>

<style lang="less" scoped>
 .key-value-row{
    overflow-y: auto;
    .key-value-item{
      margin-bottom: 10px;
      &:last-child{
        margin-bottom: 0px;
      }
      .item-value{
        margin-left: 10px;
      }  
      .item-del{
        padding: 0px 8px;
        color: rgb(var(--arcoblue-5));
        cursor: pointer;
        &:hover{
           color: rgb(var(--arcoblue-4));
        }
      }
    }
 }
</style>
