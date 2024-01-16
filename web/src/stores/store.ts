import { defineStore } from 'pinia'
import { type Torrent } from '@/types/torrent'

export interface Indexable<T = any> {
  [key: string]: T
}

export interface Sorting {
  key: string
  direction: boolean
}

export interface State {
  ok: boolean
  collator: Intl.Collator
  sorting: Sorting
  torrents: Torrent[]
}

export const useStore = defineStore('store', {
  state: (): State => ({
    ok: false,
    collator: new Intl.Collator(undefined, { numeric: true, sensitivity: 'base' }),
    sorting: {} as Sorting,
    torrents: [] as Torrent[]
  }),
  getters: {
    // returns a filtered view of torrents
    filter: (state) => {
      // if sorting key is not set, reverse torrents so it will display latest torrents first
      if (state.sorting.key === '') {
        return state.torrents.reverse()
      }
      // if sorting key is set, sort by direction where true is ascending
      if (state.sorting.key != '' && state.sorting.direction) {
        return state.torrents.sort((a, b) =>
          state.collator.compare(a[state.sorting.key], b[state.sorting.key])
        )
      }
      // if sorting key is set, sort by direction where false is descending
      if (state.sorting.key != '' && !state.sorting.direction) {
        return state.torrents.sort((a, b) =>
          state.collator.compare(b[state.sorting.key], a[state.sorting.key])
        )
      }
    }
  },
  actions: {
    // ping to check if backend is alive
    async ping() {
      const res = await fetch(`${import.meta.env.VITE_BACKEND_BASE_URL}/api/ping`, {
        method: 'GET',
      })
      const json = await res.json()
      console.log(json)
    },
    // refresh torrents
    async refresh() {
      const res = await fetch(`${import.meta.env.VITE_BACKEND_BASE_URL}/api/torrents`, {
        method: 'GET'
      })
    }
  }
})
