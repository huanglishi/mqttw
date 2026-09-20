<template>
  <div id="resize-height" class="resize-height" @mousedown="handleMousedown"></div>
</template>

<script setup lang="ts">
// v-model 接收 value，触发 update:modelValue
const props = defineProps<{
  modelValue: number
}>()
const emit = defineEmits<{
  'update:modelValue': [val: number]
}>()

const handleMousedown = (event: MouseEvent) => {
  let yValue = event.y

  // mousemove
  const handleMouseMove = (moveEvent: MouseEvent) => {
    const yMove = moveEvent.y
    const offset = yMove - yValue
    yValue = moveEvent.y
    // 更新值
    emit('update:modelValue', props.modelValue - offset)
    document.body.classList.add('select-none')
  }

  // mouseup 清理事件
  const handleMouseUp = () => {
    document.removeEventListener('mousemove', handleMouseMove)
    document.removeEventListener('mouseup', handleMouseUp)
    document.body.classList.remove('select-none')
  }

  document.addEventListener('mousemove', handleMouseMove)
  document.addEventListener('mouseup', handleMouseUp)
}
</script>

<style lang="less" scoped>
.resize-height {
  width: 100%;
  height: 3px;
  cursor: row-resize;
  transition: all 0.1s;
  margin-top: -1px;
  &:hover,
  &:active {
    background-color: var(--color-neutral-4);
  }
}
</style>
