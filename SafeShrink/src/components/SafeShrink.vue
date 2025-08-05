<template>
  <div class="min-h-screen bg-gray-900 text-white">
    <!-- Header -->
    <header class="bg-gray-800 border-b border-gray-700">
      <div class="max-w-6xl mx-auto px-4 py-6">
        <div class="flex items-center gap-3">
          <ZapIcon class="w-8 h-8 text-blue-400" />
          <h1 class="text-3xl font-bold bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent">
            SafeShrink
          </h1>
          <span 
            :class="wasmReady ? 'bg-green-500/20 text-green-400' : 'bg-yellow-500/20 text-yellow-400'"
            class="px-2 py-1 text-xs rounded-full"
          >
            {{ wasmReady ? 'WASM Ready' : 'Loading...' }}
          </span>
        </div>
        <p class="text-gray-400 mt-2">
          Compress images locally with WebAssembly - Your images never leave your browser
        </p>
      </div>
    </header>

    <div class="max-w-6xl mx-auto px-4 py-8">
      <div class="grid lg:grid-cols-3 gap-8">
        <!-- Upload Section -->
        <div class="lg:col-span-2 space-y-6">
          <!-- Drop Zone -->
          <div
            :class="dropZoneClasses"
            @dragover.prevent="isDragOver = true"
            @dragleave.prevent="isDragOver = false"
            @drop.prevent="handleDrop"
            class="border-2 border-dashed rounded-xl p-8 text-center transition-all duration-300"
          >
            <div v-if="selectedFile" class="space-y-4">
              <FileImageIcon class="w-12 h-12 text-green-400 mx-auto" />
              <div>
                <p class="text-lg font-medium text-green-400">{{ selectedFile.name }}</p>
                <p class="text-gray-400">{{ formatBytes(selectedFile.size) }}</p>
              </div>
              <button
                @click="resetAll"
                class="inline-flex items-center gap-2 px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors"
              >
                <Trash2Icon class="w-4 h-4" />
                Clear
              </button>
            </div>
            <div v-else class="space-y-4">
              <UploadIcon class="w-12 h-12 text-gray-400 mx-auto" />
              <div>
                <p class="text-xl font-medium">Drop your image here</p>
                <p class="text-gray-400">or click to browse</p>
              </div>
              <button
                @click="$refs.fileInput.click()"
                class="inline-flex items-center gap-2 px-6 py-3 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors font-medium"
              >
                <ImageIcon class="w-5 h-5" />
                Choose Image
              </button>
            </div>
            <input
              ref="fileInput"
              type="file"
              accept="image/*"
              class="hidden"
              @change="handleFileInputChange"
            />
          </div>

          <!-- Compression Button -->
          <button
            v-if="selectedFile"
            @click="compressImage"
            :disabled="isProcessing || !wasmReady"
            class="w-full py-4 bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 disabled:from-gray-600 disabled:to-gray-600 rounded-xl font-medium text-lg transition-all duration-300 flex items-center justify-center gap-3"
          >
            <div v-if="isProcessing" class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
            <ZapIcon v-else class="w-5 h-5" />
            {{ isProcessing ? 'Compressing...' : 'Compress Image' }}
          </button>

          <!-- Results -->
          <div v-if="compressionStats" class="bg-gray-800 rounded-xl p-6">
            <h3 class="text-lg font-semibold mb-4 flex items-center gap-2">
              <InfoIcon class="w-5 h-5 text-blue-400" />
              Compression Results
            </h3>
            <div class="grid grid-cols-3 gap-4 mb-4">
              <div class="text-center">
                <p class="text-2xl font-bold text-red-400">
                  {{ formatBytes(compressionStats.originalSize) }}
                </p>
                <p class="text-gray-400 text-sm">Original</p>
              </div>
              <div class="text-center">
                <p class="text-2xl font-bold text-green-400">
                  {{ formatBytes(compressionStats.compressedSize) }}
                </p>
                <p class="text-gray-400 text-sm">Compressed</p>
              </div>
              <div class="text-center">
                <p class="text-2xl font-bold text-blue-400">
                  {{ compressionStats.ratio.toFixed(1) }}%
                </p>
                <p class="text-gray-400 text-sm">Reduction</p>
              </div>
            </div>
            <button
              @click="downloadCompressed"
              class="w-full py-3 bg-green-600 hover:bg-green-700 rounded-lg font-medium transition-colors flex items-center justify-center gap-2"
            >
              <DownloadIcon class="w-5 h-5" />
              Download Compressed Image
            </button>
          </div>
        </div>

        <!-- Settings Panel -->
        <div class="space-y-6">
          <div class="bg-gray-800 rounded-xl p-6">
            <h3 class="text-lg font-semibold mb-4 flex items-center gap-2">
              <SettingsIcon class="w-5 h-5 text-blue-400" />
              Compression Settings
            </h3>
            
            <div class="space-y-4">
              <div>
                <label class="block text-sm font-medium mb-2">
                  Quality: {{ settings.quality }}%
                </label>
                <input
                  v-model.number="settings.quality"
                  type="range"
                  min="1"
                  max="100"
                  class="w-full h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer slider"
                />
                <div class="flex justify-between text-xs text-gray-400 mt-1">
                  <span>Lower size</span>
                  <span>Higher quality</span>
                </div>
              </div>

              <div>
                <label class="block text-sm font-medium mb-2">
                  Max Width (px)
                </label>
                <input
                  v-model.number="settings.maxWidth"
                  type="number"
                  placeholder="No limit"
                  class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg focus:border-blue-400 focus:outline-none"
                />
              </div>

              <div>
                <label class="block text-sm font-medium mb-2">
                  Max Height (px)
                </label>
                <input
                  v-model.number="settings.maxHeight"
                  type="number"
                  placeholder="No limit"
                  class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg focus:border-blue-400 focus:outline-none"
                />
              </div>

              <div class="flex items-center">
                <input
                  v-model="settings.keepAspect"
                  type="checkbox"
                  id="keepAspect"
                  class="w-4 h-4 text-blue-600 bg-gray-700 border-gray-600 rounded focus:ring-blue-500"
                />
                <label for="keepAspect" class="ml-2 text-sm">
                  Keep aspect ratio
                </label>
              </div>
            </div>
          </div>

          <!-- Preview Images -->
          <div v-if="previewUrl" class="bg-gray-800 rounded-xl p-4">
            <h4 class="font-medium mb-3">Original</h4>
            <img
              :src="previewUrl"
              alt="Original"
              class="w-full rounded-lg border border-gray-600"
            />
          </div>

          <div v-if="compressedUrl" class="bg-gray-800 rounded-xl p-4">
            <h4 class="font-medium mb-3">Compressed</h4>
            <img
              :src="compressedUrl"
              alt="Compressed"
              class="w-full rounded-lg border border-gray-600"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

// Icons (you can replace these with actual icon components)
const ZapIcon = { template: '<svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path></svg>' }
const UploadIcon = { template: '<svg class="w-12 h-12" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"></path></svg>' }
const DownloadIcon = { template: '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>' }
const ImageIcon = { template: '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>' }
const FileImageIcon = { template: '<svg class="w-12 h-12" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>' }
const Trash2Icon = { template: '<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>' }
const SettingsIcon = { template: '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path></svg>' }
const InfoIcon = { template: '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>' }

// Reactive state
const isDragOver = ref(false)
const selectedFile = ref(null)
const previewUrl = ref('')
const compressedUrl = ref('')
const isProcessing = ref(false)
const compressionStats = ref(null)
const wasmReady = ref(false)
const fileInput = ref(null)

const settings = ref({
  quality: 85,
  maxWidth: 0,
  maxHeight: 0,
  keepAspect: true
})

// Computed properties
const dropZoneClasses = computed(() => {
  if (isDragOver.value) {
    return 'border-blue-400 bg-blue-400/10'
  } else if (selectedFile.value) {
    return 'border-green-400 bg-green-400/5'
  } else {
    return 'border-gray-600 bg-gray-800/50 hover:border-gray-500'
  }
})

// Initialize WebAssembly
onMounted(async () => {
  try {
    if (!window.Go) {
      // Load wasm_exec.js
      const script = document.createElement('script')
      script.src = '/wasm_exec.js'
      script.onload = async () => {
        const go = new window.Go()
        const result = await WebAssembly.instantiateStreaming(
          fetch('/safeshrink.wasm'),
          go.importObject
        )
        go.run(result.instance)
        wasmReady.value = true
      }
      document.head.appendChild(script)
    } else {
      wasmReady.value = true
    }
  } catch (error) {
    console.error('Failed to initialize WASM:', error)
  }
})

// Event handlers
const handleDrop = (e) => {
  isDragOver.value = false
  const files = Array.from(e.dataTransfer.files)
  if (files.length > 0) {
    handleFileSelect(files[0])
  }
}

const handleFileInputChange = (e) => {
  const file = e.target.files[0]
  if (file) {
    handleFileSelect(file)
  }
}

const handleFileSelect = (file) => {
  if (!file || !file.type.startsWith('image/')) {
    alert('Please select a valid image file')
    return
  }

  selectedFile.value = file
  compressedUrl.value = ''
  compressionStats.value = null

  // Create preview
  const reader = new FileReader()
  reader.onload = (e) => {
    previewUrl.value = e.target.result
  }
  reader.readAsDataURL(file)
}

const compressImage = async () => {
  if (!selectedFile.value || !wasmReady.value || !window.compressImage) {
    alert('WASM module not ready or no file selected')
    return
  }

  isProcessing.value = true

  try {
    // Convert file to base64
    const reader = new FileReader()
    reader.onload = async (e) => {
      const base64Data = e.target.result.split(',')[1]
      
      // Call WASM function
      const result = window.compressImage(base64Data, settings.value)
      
      if (result.success) {
        const compressedBlob = base64ToBlob(result.compressedData, 'image/jpeg')
        const compressedUrl_val = URL.createObjectURL(compressedBlob)
        
        compressedUrl.value = compressedUrl_val
        compressionStats.value = {
          originalSize: result.originalSize,
          compressedSize: result.compressedSize,
          ratio: result.compressionRatio
        }
      } else {
        alert('Compression failed: ' + result.message)
      }
      
      isProcessing.value = false
    }
    
    reader.readAsDataURL(selectedFile.value)
  } catch (error) {
    console.error('Compression error:', error)
    alert('Compression failed: ' + error.message)
    isProcessing.value = false
  }
}

const base64ToBlob = (base64, mimeType) => {
  const byteCharacters = atob(base64)
  const byteNumbers = new Array(byteCharacters.length)
  for (let i = 0; i < byteCharacters.length; i++) {
    byteNumbers[i] = byteCharacters.charCodeAt(i)
  }
  const byteArray = new Uint8Array(byteNumbers)
  return new Blob([byteArray], { type: mimeType })
}

const downloadCompressed = () => {
  if (!compressedUrl.value) return
  
  const a = document.createElement('a')
  a.href = compressedUrl.value
  a.download = `compressed_${selectedFile.value.name.split('.')[0]}.jpg`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
}

const formatBytes = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

const resetAll = () => {
  selectedFile.value = null
  previewUrl.value = ''
  compressedUrl.value = ''
  compressionStats.value = null
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}
</script>

<style scoped>
.slider::-webkit-slider-thumb {
  appearance: none;
  height: 20px;
  width: 20px;
  border-radius: 50%;
  background: #3b82f6;
  cursor: pointer;
  box-shadow: 0 0 2px 0 #555;
  transition: background .15s ease-in-out;
}

.slider::-webkit-slider-thumb:hover {
  background: #2563eb;
}

.slider::-moz-range-thumb {
  height: 20px;
  width: 20px;
  border-radius: 50%;
  background: #3b82f6;
  cursor: pointer;
  border: none;
  box-shadow: 0 0 2px 0 #555;
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>