<template>
  <button 
    :class="buttonClasses" 
    :disabled="disabled"
    @click="$emit('click')"
  >
    <slot name="icon" />
    <span v-if="$slots.default" :class="{ 'ml-2': $slots.icon }">
      <slot />
    </span>
  </button>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  variant: {
    type: String,
    default: 'primary',
    validator: (value) => ['primary', 'secondary', 'ghost', 'danger'].includes(value)
  },
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['sm', 'md', 'lg'].includes(value)
  },
  disabled: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  },
  fullWidth: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['click'])

const buttonClasses = computed(() => {
  const baseClasses = 'inline-flex items-center justify-center font-medium rounded-lg border transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed'
  
  const variantClasses = {
    primary: 'bg-accent text-background border-accent hover:bg-accent-hover focus:ring-accent shadow-soft hover:shadow-medium',
    secondary: 'bg-surface text-primary-text border-borders hover:bg-hover focus:ring-borders shadow-soft hover:shadow-medium',
    ghost: 'bg-transparent text-primary-text border-transparent hover:bg-hover focus:ring-borders',
    danger: 'bg-error text-background border-error hover:bg-accent-hover focus:ring-error shadow-soft hover:shadow-medium'
  }
  
  const sizeClasses = {
    sm: 'px-3 py-1.5 text-sm',
    md: 'px-4 py-2 text-sm',
    lg: 'px-6 py-3 text-base'
  }
  
  const widthClasses = props.fullWidth ? 'w-full' : ''
  
  return [
    baseClasses,
    variantClasses[props.variant],
    sizeClasses[props.size],
    widthClasses
  ].filter(Boolean).join(' ')
})
</script>
