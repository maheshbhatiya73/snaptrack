<script setup>
import { ref, onMounted } from 'vue'
import { fetchRecentActivity } from '@/lib/api'

const activities = ref([])
const loading = ref(true)

onMounted(async () => {
  loading.value = true
  activities.value = await fetchRecentActivity()
  loading.value = false
})

function formatDate(dateStr) {
  const d = new Date(dateStr)
  return d.toLocaleString()
}

const getLevelColor = (level) => {
  switch(level) {
    case 'info': return 'text-blue-500 bg-blue-100'
    case 'warning': return 'text-yellow-600 bg-yellow-100'
    case 'error': return 'text-red-500 bg-red-100'
    default: return 'text-gray-500 bg-gray-100'
  }
}

const getIcon = (level) => {
  switch(level) {
    case 'info': return 'ℹ️'
    case 'warning': return '⚠️'
    case 'error': return '❌'
    default: return '🔔'
  }
}
</script>

<template>
  <div class="px-4 py-6">
    <h2 class="text-xl font-semibold text-gray-800 mb-4">Recent Activity</h2>

    <div v-if="loading" class="text-gray-500 text-center py-6">
      Loading recent activity...
    </div>

    <div v-else-if="activities.length > 0" class="space-y-4">
      <div
        v-for="log in activities"
        :key="log.id"
        class="flex items-start space-x-4 p-4 bg-white rounded-xl shadow hover:shadow-lg transition"
      >
        <!-- Icon -->
        <div class="flex-shrink-0 w-10 h-10 flex items-center justify-center rounded-full"
             :class="getLevelColor(log.level)">
          <span class="text-lg">{{ getIcon(log.level) }}</span>
        </div>

        <!-- Content -->
        <div class="flex-1">
          <div class="flex justify-between items-center mb-1">
            <span class="font-medium capitalize" :class="getLevelColor(log.level)">
              {{ log.level }}
            </span>
            <span class="text-xs text-gray-400">{{ formatDate(log.created_at) }}</span>
          </div>

          <p class="text-gray-800 font-medium mb-1">{{ log.message }}</p>

          <div v-if="log.entity_type || log.metadata" class="text-xs text-gray-500 space-y-1">
            <div v-if="log.entity_type">
              <span class="font-semibold">Entity:</span> {{ log.entity_type }} 
              <span v-if="log.entity_id">(#{{ log.entity_id }})</span>
            </div>
            <div v-if="log.metadata">
              <span class="font-semibold">Details:</span>
              <pre class="whitespace-pre-wrap break-words bg-gray-50 p-2 rounded">{{ JSON.stringify(log.metadata, null, 2) }}</pre>
            </div>
          </div>
        </div>
      </div>
    </div>

    <p v-else class="text-gray-400 text-center py-6">
      No recent activity found.
    </p>
  </div>
</template>

<style scoped>
/* Optional: Add a subtle timeline line */
</style>
