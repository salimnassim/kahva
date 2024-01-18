<template>
  <div v-if="store.filter === undefined" class="w-1/2 pt-6 m-auto text-center">
    <div class="text-lg">
      Waiting for server...
    </div>
    <div>
      If this message does not disappear, check the browser console tab for errors.
    </div>
  </div>
  <div class="h-full w-full flex flex-col bg-slate-900 text-slate-400" v-if="store.filter !== undefined">
    <div class="grow overflow-x-auto overflow-y-auto w-full">
      <table class="w-full table-auto">
        <thead>
          <tr class="text-sm text-left border-b border-slate-400">
            <th><input class="w-full bg-slate-900 text-white outline-none" v-model="search" type="text"
                placeholder="Search" /></th>
            <th class="hover:cursor-pointer" @click="sort('name')">Priority</th>
            <th class="hover:cursor-pointer" @click="sort('completed_bytes')">Size</th>
            <th class="hover:cursor-pointer" @click="sort('upload_rate')">Completed</th>
            <th>%</th>
            <th class="hover:cursor-pointer" @click="sort('upload_rate')">Up Rate</th>
            <th class="hover:cursor-pointer" @click="sort('download_rate')">Down Rate</th>
            <th class="hover:cursor-pointer" @click="sort('upload_total')">Up Total</th>
            <th class="hover:cursor-pointer" @click="sort('download_total')">Down Total</th>
            <th class="hover:cursor-pointer" @click="sort('leechers')">Leechers</th>
            <th class="hover:cursor-pointer" @click="sort('seeders')">Seeders</th>
          </tr>
        </thead>
        <tbody>
          <tr class="text-sm even:bg-slate-800 hover:bg-slate-600 hover:text-white" v-for="torrent in store.filter"
            :key="torrent.hash" :data-hash="torrent.hash">
            <td class="flex flex-row justify-between">
              <div>{{ torrent.name }}</div>
              <div class="flex flex-row">
                <div class="text-sm mr-2" v-if="torrent.message" :title="torrent.message">💬</div>
              </div>
            </td>
            <td>{{ priority(torrent.priority) }}</td>
            <td>{{ human(torrent.size_bytes) }}</td>
            <td>{{ human(torrent.completed_bytes) }}</td>
            <td>{{ torrent.size_bytes / torrent.completed_bytes * 100 }}%</td>
            <td>{{ human(torrent.upload_rate) }}</td>
            <td>{{ human(torrent.download_rate) }}</td>
            <td>{{ human(torrent.upload_total) }}</td>
            <td>{{ human(torrent.download_total) }}</td>
            <td>{{ torrent.leechers }}</td>
            <td>{{ torrent.seeders }}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="flex flex-row justify-between bg-black text-white text-sm">
      <div class="flex flex-row">
        <div>Throttle [{{ human(store.system?.throttle_global_up_max_rate ?? 0) }} / {{
          human(store.system?.throttle_global_down_max_rate ?? 0) }}]</div>
        <div class="pl-2">Rate [{{ human(store.system?.throttle_global_up_rate ?? 0) }} / {{
          human(store.system?.throttle_global_down_rate ?? 0) }}]</div>
      </div>
      <div class="flex flex-row">
        <div>rtorrent {{ store.system?.client_version }}/{{ store.system?.library_version }}</div>
        <div class="pl-2">{{ store.system?.pid }}@{{ store.system?.hostname }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useStore } from '@/stores/store'
import { computed } from 'vue';
const store = useStore()

const search = computed({
  get(): string {
    return store.ui.search
  },
  set(query: string): void {
    store.search(query)
  }
})

function priority(priority: number): string {
  switch (priority) {
    case 0: return 'Off';
    case 1: return 'Low';
    case 2: return 'Normal';
    case 3: return 'High';
    default: return '?';
  }
}

// Converts bytes to a human readable format.
function human(bytes: number) {
  if (bytes === undefined || bytes === 0) {
    return '0'
  }
  const index = Math.floor(Math.log(bytes) / Math.log(1024))
  return `${Number((bytes / Math.pow(1024, index)).toFixed(2)) * 1} ${(['B', 'KB', 'MB', 'GB', 'TB'] as const)[index]}`
}

// Changes the UI sorting state.
function sort(key: string) {
  if (store.ui.sorting.key === key) {
    store.sort(key, !store.ui.sorting.direction)
    return
  }
  store.sort(key, false)
}

// Refresh view every n seconds
store.refreshView()
store.refreshSystem()
window.setInterval(function () {
  store.refreshView()
  store.refreshSystem()
}, 1000 * store.ui.updateInterval)
</script>