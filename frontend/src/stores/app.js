import { defineStore } from 'pinia';
import { Position } from '@vue-flow/core';

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
        maxZoom: 3,
        arrow: true,
        layout:'vertical',
        snapToGrid: false,
        lineType: 'default',
        snapGrid: { 0: 15, 1: 15 },
      },
    },
    retrieved: true,
  }),
  getters: {
    inputPosition: (state) => state.settings.editor.layout == 'horizontal' ? Position.Left : Position.Top,
    outputPosition: (state) => state.settings.editor.layout == 'horizontal' ? Position.Right : Position.Bottom,
  },
  actions: {
    loadSettings() {
      return this.settings
    },
    async updateSettings(settings = {}) {
      this.settings = settings
    },
  },
});
