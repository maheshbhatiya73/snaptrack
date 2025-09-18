<template>
  <div></div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { isAuthenticated } from '~/lib/api'

const router = useRouter()

onMounted(() => {
  if (isAuthenticated()) {
    router.replace('/dashboard')
  } else {
    router.replace('/auth/login')
  }
})
</script>
