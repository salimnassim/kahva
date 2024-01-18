import { defineStore } from 'pinia'
import { IsViewResponse, type Torrent } from '@/types/torrent'
import { IsSystemResponse, type System } from '@/types/system'

export interface Indexable<T = any> {
  [key: string]: T
}

export interface Sorting {
  key: string
  direction: boolean
}

export interface UserInterface {
  collator: Intl.Collator
  sorting: Sorting
  search: string
  updateInterval: number
  view: string
}

export interface State {
  userInterface: UserInterface
  torrents: Torrent[] | undefined
  system: System | undefined
}

export const useStore = defineStore('store', {
  state: (): State => ({
    userInterface: {
      collator: new Intl.Collator(undefined, { numeric: true, sensitivity: 'base' }),
      sorting: { key: '', direction: false },
      search: '',
      updateInterval: 30,
      view: 'main'
    },
    torrents: undefined,
    system: undefined
  }),
  getters: {
    ui: (state) => {
      return state.userInterface
    },
    // returns a filtered view of torrents
    filter: (state) => {
      if (state.torrents === undefined) {
        return undefined
      }

      let items = [] as Torrent[]

      // if sorting key is not set, reverse torrents so it will display latest torrents first
      if (state.userInterface.sorting.key === '') {
        items = state.torrents.reverse()
      }
      // if sorting key is set, sort by direction where true is ascending
      if (state.userInterface.sorting.key != '' && state.userInterface.sorting.direction) {
        items = state.torrents.sort((a, b) =>
          state.userInterface.collator.compare(
            a[state.userInterface.sorting.key],
            b[state.userInterface.sorting.key]
          )
        )
      }
      // if sorting key is set, sort by direction where false is descending
      if (state.userInterface.sorting.key != '' && !state.userInterface.sorting.direction) {
        items = state.torrents.sort((a, b) =>
          state.userInterface.collator.compare(
            b[state.userInterface.sorting.key],
            a[state.userInterface.sorting.key]
          )
        )
      }

      if (state.userInterface.search != '') {
        items = items.filter((n) => n.name.toLocaleLowerCase().match(state.userInterface.search))
      }

      return items
    }
  },
  actions: {
    async sort(key: string, direction: boolean) {
      this.ui.sorting.key = key
      this.ui.sorting.direction = direction
    },
    async search(query: string) {
      this.ui.search = query
    },
    // refresh view
    async refreshView() {
      try {
        const res = await fetch(
          `${import.meta.env.VITE_BACKEND_BASE_URL}/api/view/${this.userInterface.view}`,
          {
            method: 'GET'
          }
        )
        if (!res.ok) {
          console.error(res.status, res.statusText)
          return
        }
        const json = await res.json()
        if (!IsViewResponse(json)) {
          console.error(`malformed view response`)
          return
        }
        this.torrents = json.torrents
      } catch (err: any) {
        console.error(err)
        return
      }
    },
    // refresh system
    async refreshSystem() {
      try {
        const res = await fetch(`${import.meta.env.VITE_BACKEND_BASE_URL}/api/system`, {
          method: 'GET'
        })
        if (!res.ok) {
          console.error(res.status, res.statusText)
          return
        }
        const json = await res.json()
        if (!IsSystemResponse(json)) {
          console.error(`malformed system response`)
          return
        }
        this.system = json.system
      } catch (err: any) {
        console.error(err)
        return
      }
    }
  }
})
