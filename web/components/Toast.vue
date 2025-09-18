<template>
  <transition name="toast-fade" appear>
    <div
      v-show="visible"
      :class="[
        'fixed top-5 right-5 z-50 px-6 py-4 rounded-xl shadow-lg flex items-center space-x-3',
        typeClass
      ]"
    >
      <div class="flex-1">
        <p class="font-semibold">{{ message }}</p>
      </div>
      <button @click="close" class="text-white font-bold text-xl leading-none">&times;</button>
    </div>
  </transition>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'

const props = defineProps({
  message: { type: String, required: true },
  type: { type: String, default: 'success' },
  duration: { type: Number, default: 3000 }
})

const visible = ref(false)

const typeClass = computed(() => {
  switch (props.type) {
    case 'success':
      return 'bg-green-500 text-white'
    case 'error':
      return 'bg-red-500 text-white'
    case 'info':
      return 'bg-blue-500 text-white'
    default:
      return 'bg-gray-500 text-white'
  }
})

const close = () => {
  visible.value = false
}

onMounted(() => {
  visible.value = true
  setTimeout(() => {
    visible.value = false
  }, props.duration)
})
</script>

<style>
.toast-fade-enter-active,
.toast-fade-leave-active {
  transition: all 0.3s ease;
}
.toast-fade-enter-from {
  opacity: 0;
  transform: translateY(-10px);
}
.toast-fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
