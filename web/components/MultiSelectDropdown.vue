<template>
  <div class="relative" ref="dropdownRef">
    <button
      type="button"
      @click="toggleDropdown"
      :disabled="disabled"
      :class="[
        'relative w-full bg-white border border-slate-300 rounded-lg shadow-sm pl-3 pr-10 py-2.5 text-left cursor-default focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 sm:text-sm min-h-[42px]',
        disabled ? 'bg-slate-50 text-slate-400 cursor-not-allowed' : 'hover:border-slate-400',
        error ? 'border-red-300 focus:ring-red-500 focus:border-red-500' : ''
      ]"
    >
      <div class="flex flex-wrap gap-1">
        <span
          v-for="selectedOption in selectedOptions"
          :key="selectedOption[valueKey]"
          class="inline-flex items-center px-2 py-1 rounded-md text-xs font-medium bg-blue-100 text-blue-800"
        >
          {{ selectedOption[labelKey] }}
          <button
            type="button"
            @click.stop="removeOption(selectedOption)"
            class="ml-1 inline-flex items-center justify-center w-4 h-4 rounded-full hover:bg-blue-200"
          >
            <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
            </svg>
          </button>
        </span>
        <span v-if="selectedOptions.length === 0" class="text-slate-500">
          {{ placeholder }}
        </span>
      </div>
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
        
        <div v-if="showSelectAll" class="px-3 py-2 border-b border-slate-200">
          <button
            type="button"
            @click="toggleSelectAll"
            :class="[
              'w-full text-left px-2 py-1 text-sm rounded hover:bg-slate-100',
              allSelected ? 'text-blue-600 font-medium' : 'text-slate-700'
            ]"
          >
            {{ allSelected ? 'Deselect All' : 'Select All' }}
          </button>
        </div>
        
        <div v-if="filteredOptions.length === 0" class="px-3 py-2 text-sm text-slate-500">
          {{ searchQuery ? 'No options found' : 'No options available' }}
        </div>
        
        <div v-else>
          <button
            v-for="option in filteredOptions"
            :key="option[valueKey]"
            type="button"
            @click="toggleOption(option)"
            :class="[
              'relative w-full px-3 py-2 text-left cursor-default select-none hover:bg-blue-50 hover:text-blue-900',
              isSelected(option) ? 'bg-blue-100 text-blue-900' : 'text-slate-900'
            ]"
          >
            <div class="flex items-center">
              <div class="flex items-center h-5">
                <input
                  :checked="isSelected(option)"
                  type="checkbox"
                  class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-slate-300 rounded"
                  @click.stop
                />
              </div>
              <div class="ml-3">
                <span class="block truncate font-normal">
                  {{ option[labelKey] }}
                </span>
                <span v-if="option.description" class="block text-xs text-slate-500">
                  {{ option.description }}
                </span>
              </div>
            </div>
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
    type: Array,
    default: () => []
  },
  options: {
    type: Array,
    required: true
  },
  placeholder: {
    type: String,
    default: 'Select options'
  },
  disabled: {
    type: Boolean,
    default: false
  },
  searchable: {
    type: Boolean,
    default: true
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
  },
  showSelectAll: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

const dropdownRef = ref(null)
const searchInput = ref(null)
const isOpen = ref(false)
const searchQuery = ref('')

const selectedOptions = computed(() => {
  return props.options.filter(option => 
    props.modelValue.includes(option[props.valueKey])
  )
})

const filteredOptions = computed(() => {
  if (!props.searchable || !searchQuery.value) {
    return props.options
  }
  
  const query = searchQuery.value.toLowerCase()
  return props.options.filter(option =>
    option[props.labelKey].toLowerCase().includes(query) ||
    (option.description && option.description.toLowerCase().includes(query))
  )
})

const allSelected = computed(() => {
  return props.options.length > 0 && props.options.every(option => 
    props.modelValue.includes(option[props.valueKey])
  )
})

const isSelected = (option) => {
  return props.modelValue.includes(option[props.valueKey])
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

const toggleOption = (option) => {
  const value = option[props.valueKey]
  const newValue = [...props.modelValue]
  
  if (isSelected(option)) {
    const index = newValue.indexOf(value)
    if (index > -1) {
      newValue.splice(index, 1)
    }
  } else {
    newValue.push(value)
  }
  
  emit('update:modelValue', newValue)
  emit('change', newValue)
}

const removeOption = (option) => {
  const value = option[props.valueKey]
  const newValue = props.modelValue.filter(v => v !== value)
  emit('update:modelValue', newValue)
  emit('change', newValue)
}

const toggleSelectAll = () => {
  if (allSelected.value) {
    emit('update:modelValue', [])
  } else {
    const allValues = props.options.map(option => option[props.valueKey])
    emit('update:modelValue', allValues)
  }
  emit('change', props.modelValue)
}

const handleClickOutside = (event) => {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target)) {
    isOpen.value = false
    searchQuery.value = ''
  }
}

watch(() => props.modelValue, () => {
  emit('change', props.modelValue)
})

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>
