<template>
  <button class="btn" @click="$emit('click')">
    <slot />
  </button>
</template>

<script setup>
</script>

<style scoped>
.btn {
  padding: 8px 16px;
  background: #4f46e5;
  color: white;
  border-radius: 6px;
  cursor: pointer;
}
.btn:hover {
  background: #4338ca;
}
</style>
