import { defineStore } from 'pinia';

export const useAppStore = defineStore('app', {
  storageMap: {
    tabs: 'tabs',
    settings: 'settings',
  },
  state: () => ({
    tabs: [],
    copiedEls: {
      edges: [],
      nodes: [],
    },
    settings: {
      locale: 'en',
      operateSystem:'win10',
      deleteLogAfter: 30,
      logsLimit: 1000,
      editor: {
        minZoom: 0.3,
        maxZoom: 1.3,
        arrow: true,
        snapToGrid: false,
        lineType: 'default',
        snapGrid: { 0: 15, 1: 15 },
      },
    },
    retrieved: true,
  }),
  actions: {
    loadSettings() {
      return this.settings
    },
    async updateSettings(settings = {}) {
      this.settings = settings
    },
  },
});
