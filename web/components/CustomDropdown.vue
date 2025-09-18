<template>
  <div class="relative" ref="dropdownRef">
    <button
      type="button"
      @click="toggleDropdown"
      :disabled="disabled"
      :class="[
        'relative w-full bg-white border border-slate-300 rounded-lg shadow-sm pl-3 pr-10 py-2.5 text-left cursor-default focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 sm:text-sm',
        disabled ? 'bg-slate-50 text-slate-400 cursor-not-allowed' : 'hover:border-slate-400',
        error ? 'border-red-300 focus:ring-red-500 focus:border-red-500' : ''
      ]"
    >
      <span class="block truncate" :class="{ 'text-slate-500': !selectedText }">
        {{ selectedText || placeholder }}
      </span>
      <span class="absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none">
        <svg
          class="w-5 h-5 text-slate-400 transition-transform duration-200"
          :class="{ 'rotate-180': isOpen }"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
        </svg>
      </span>
    </button>

    <Transition
      enter-active-class="transition duration-100 ease-out"
      enter-from-class="transform scale-95 opacity-0"
      enter-to-class="transform scale-100 opacity-100"
      leave-active-class="transition duration-75 ease-in"
      leave-from-class="transform scale-100 opacity-100"
      leave-to-class="transform scale-95 opacity-0"
    >
      <div
        v-if="isOpen"
        class="absolute z-50 mt-1 w-full bg-white shadow-lg max-h-60 rounded-md py-1 text-base ring-1 ring-black ring-opacity-5 overflow-auto focus:outline-none sm:text-sm"
      >
        <div v-if="searchable" class="px-3 py-2 border-b border-slate-200">
          <input
            ref="searchInput"
            v-model="searchQuery"
            type="text"
            placeholder="Search..."
            class="w-full px-2 py-1 text-sm border border-slate-300 rounded focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            @click.stop
          />
        </div>
        
        <div v-if="filteredOptions.length === 0" class="px-3 py-2 text-sm text-slate-500">
          {{ searchQuery ? 'No options found' : 'No options available' }}
        </div>
        
        <div v-else>
          <button
            v-for="option in filteredOptions"
            :key="option.value"
            type="button"
            @click="selectOption(option)"
            :class="[
              'relative w-full px-3 py-2 text-left cursor-default select-none hover:bg-blue-50 hover:text-blue-900',
              isSelected(option) ? 'bg-blue-100 text-blue-900' : 'text-slate-900'
            ]"
          >
            <span class="block truncate font-normal">
              {{ option.label }}
            </span>
            <span
              v-if="isSelected(option)"
              class="absolute inset-y-0 right-0 flex items-center pr-4 text-blue-600"
            >
              <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
              </svg>
            </span>
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  modelValue: {
    type: [String, Number, Object],
    default: null
  },
  options: {
    type: Array,
    required: true
  },
  placeholder: {
    type: String,
    default: 'Select an option'
  },
  disabled: {
    type: Boolean,
    default: false
  },
  searchable: {
    type: Boolean,
    default: false
  },
  error: {
    type: Boolean,
    default: false
  },
  labelKey: {
    type: String,
    default: 'label'
  },
  valueKey: {
    type: String,
    default: 'value'
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

const dropdownRef = ref(null)
const searchInput = ref(null)
const isOpen = ref(false)
const searchQuery = ref('')

const selectedText = computed(() => {
  if (!props.modelValue) return ''
  
  const selectedOption = props.options.find(option => 
    option[props.valueKey] === props.modelValue
  )
  
  return selectedOption ? selectedOption[props.labelKey] : ''
})

const filteredOptions = computed(() => {
  if (!props.searchable || !searchQuery.value) {
    return props.options
  }
  
  const query = searchQuery.value.toLowerCase()
  return props.options.filter(option =>
    option[props.labelKey].toLowerCase().includes(query)
  )
})

const isSelected = (option) => {
  return option[props.valueKey] === props.modelValue
}

const toggleDropdown = () => {
  if (props.disabled) return
  
  isOpen.value = !isOpen.value
  
  if (isOpen.value && props.searchable) {
    nextTick(() => {
      searchInput.value?.focus()
    })
  }
}

const selectOption = (option) => {
  emit('update:modelValue', option[props.valueKey])
  emit('change', option)
  isOpen.value = false
  searchQuery.value = ''
}

const handleClickOutside = (event) => {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target)) {
    isOpen.value = false
    searchQuery.value = ''
  }
}

watch(() => props.modelValue, () => {
  emit('change', props.options.find(option => 
    option[props.valueKey] === props.modelValue
  ))
})

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>
