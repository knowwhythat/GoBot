import { defineStore } from "pinia";
import { Position } from "@vue-flow/core";

export const useAppStore = defineStore("app", {
  state: () => ({
    settings: {
      locale: "zh",
      operateSystem: "win10",
      deleteLogAfter: 30,
      logsLimit: 1000,
      editor: {
        minZoom: 0.3,
        maxZoom: 3,
        arrow: true,
        layout: "vertical",
        snapToGrid: false,
        lineType: "default",
        snapGrid: { 0: 15, 1: 15 },
      },
    },
    retrieved: true,
    loading: false,
  }),
  getters: {
    inputPosition: (state) =>
      state.settings.editor.layout == "horizontal"
        ? Position.Left
        : Position.Top,
    outputPosition: (state) =>
      state.settings.editor.layout == "horizontal"
        ? Position.Right
        : Position.Bottom,
  },
  actions: {
    loadSettings() {},
    async updateSettings(settings = {}) {
      this.settings = settings;
    },
    changeLoadingState(value) {
      this.loading = value;
    },
  },
});
