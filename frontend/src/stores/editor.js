import { defineStore } from "pinia";

export const useEditorStore = defineStore("editor", {
  state: () => ({
    copiedBlocks: [],
  }),
  getters: {},
  actions: {
    addToPasteBlocks(blocks) {
      return (this.copiedBlocks = blocks);
    },
  },
});
