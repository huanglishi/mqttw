<template>
  <!-- 必须写 name="collapse-transition"，Vue 才会追加对应的class -->
  <Transition
    name="collapse-transition"
    @before-enter="beforeEnter"
    @enter="enter"
    @after-enter="afterEnter"
    @before-leave="beforeLeave"
    @leave="leave"
    @after-leave="afterLeave"
  >
    <slot />
  </Transition>
</template>

<script setup lang="ts">
type TransitionHook = (el: Element) => void
type TransitionHookWithDone = (el: Element, done: () => void) => void

const emit = defineEmits<{
  afterEnter: []
  afterLeave: []
}>()

const beforeEnter: TransitionHook = (el) => {
  const node = el as HTMLElement
  // 缓存原始padding
  node.dataset.oldPaddingTop = node.style.paddingTop
  node.dataset.oldPaddingBottom = node.style.paddingBottom

  node.style.maxHeight = '0px'
  node.style.paddingTop = '0px'
  node.style.paddingBottom = '0px'
  node.style.opacity = '0'
  node.style.overflow = 'hidden'
}

const enter: TransitionHookWithDone = (el, done) => {
  const node = el as HTMLElement
  // 强制回流，保证动画触发
  node.getBoundingClientRect()
  // 开启过渡
  node.style.transition = 'max-height 240ms ease-in-out, padding 240ms ease-in-out, opacity 240ms ease-in-out'
  node.style.maxHeight = `${node.scrollHeight}px`
  node.style.opacity = '1'

  const onTransitionEnd = () => {
    node.removeEventListener('transitionend', onTransitionEnd)
    done()
  }
  node.addEventListener('transitionend', onTransitionEnd)
}

const afterEnter: TransitionHook = (el) => {
  const node = el as HTMLElement
  // 动画结束清除maxHeight，支持内容自适应
  node.style.maxHeight = ''
  node.style.paddingTop = node.dataset.oldPaddingTop || ''
  node.style.paddingBottom = node.dataset.oldPaddingBottom || ''
  node.style.overflow = ''
  emit('afterEnter')
}

const beforeLeave: TransitionHook = (el) => {
  const node = el as HTMLElement
  node.style.maxHeight = `${node.scrollHeight}px`
  node.style.opacity = '1'
  node.style.overflow = 'hidden'
}

const leave: TransitionHookWithDone = (el, done) => {
  const node = el as HTMLElement
  // 强制回流
  node.getBoundingClientRect()
  node.style.transition = 'max-height 240ms ease-in-out, padding 240ms ease-in-out, opacity 240ms ease-in-out'
  node.style.maxHeight = '0px'
  node.style.opacity = '0'

  const onTransitionEnd = () => {
    node.removeEventListener('transitionend', onTransitionEnd)
    done()
  }
  node.addEventListener('transitionend', onTransitionEnd)
}

const afterLeave: TransitionHook = (el) => {
  const node = el as HTMLElement
  node.style.maxHeight = ''
  node.style.paddingTop = node.dataset.oldPaddingTop || ''
  node.style.paddingBottom = node.dataset.oldPaddingBottom || ''
  emit('afterLeave')
}
</script>

<style scoped>
/* scoped 下必须 :deep() 选中插槽内的根DOM */
:deep(.collapse-transition-enter-active),
:deep(.collapse-transition-leave-active) {
  overflow: hidden;
}
</style>
