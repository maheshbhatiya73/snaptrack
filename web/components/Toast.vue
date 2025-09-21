<template>
  <transition name="toast-fade" appear>
    <div
      v-show="visible"
      class="fixed top-5 right-5 z-50 bg-white border rounded-lg shadow-lg flex items-start space-x-3 p-4 max-w-sm"
    >
      <div :class="iconClass">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path :d="iconPath" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"/>
        </svg>
      </div>
      <div class="flex-1">
        <p :class="textClass">{{ message }}</p>
      </div>
      <button @click="close" class="text-gray-400 hover:text-gray-600 transition-colors">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
        </svg>
      </button>
    </div>
  </transition>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

const props = defineProps({
  message: { type: String, required: true },
  type: { type: String, default: 'success' },
  duration: { type: Number, default: 3000 }
})

const visible = ref(false)

const iconClass = computed(() => {
  switch (props.type) {
    case 'success':
      return 'text-green-600'
    case 'error':
      return 'text-red-600'
    case 'info':
      return 'text-blue-600'
    default:
      return 'text-gray-600'
  }
})

const textClass = computed(() => {
  switch (props.type) {
    case 'success':
      return 'text-green-800 font-medium'
    case 'error':
      return 'text-red-800 font-medium'
    case 'info':
      return 'text-blue-800 font-medium'
    default:
      return 'text-gray-800 font-medium'
  }
})

const iconPath = computed(() => {
  switch (props.type) {
    case 'success':
      return 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z'
    case 'error':
      return 'M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
    case 'info':
      return 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
    default:
      return 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
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
