<template>
    <div v-if="store.filter?.length == 0" class="flex flex-col pt-6 w-1/2 m-auto text-center">
        <div>
            <div class="text-xl">Connecting...</div>
            <div class="text-base">If this screen does not disappear check your environment variable settings and network
                log.</div>
        </div>
    </div>
</template>
  
<script setup lang="ts">
import { useStore } from '@/stores/store'
const store = useStore()
</script>