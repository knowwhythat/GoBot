import { onUnmounted, onMounted } from "vue";
import Mousetrap from "mousetrap";
import { isObject } from "@/utils/helper";

const defaultShortcut = {
  "editor:save": {
    id: "editor:save",
    combo: "mod+s",
  },
  "editor:cut-block": {
    id: "editor:cut-block",
    combo: "mod+x",
  },
  "editor:copy-block": {
    id: "editor:copy-block",
    combo: "mod+c",
  },
  "editor:paste-block": {
    id: "editor:paste-block",
    combo: "mod+v",
  },
  "editor:del-block": {
    id: "editor:del-block",
    combo: "del",
  },
  "editor:execute-flow": {
    id: "editor:execute-flow",
    combo: "shift+f10",
  },
  "editor:debug-flow": {
    id: "editor:debug-flow",
    combo: "shift+f9",
  },
  "editor:terminate-flow": {
    id: "editor:terminate-flow",
    combo: "mod+f2",
  },
};

export function getReadableShortcut(str) {
  const list = {
    option: {
      win: "alt",
      mac: "option",
    },
    mod: {
      win: "ctrl",
      mac: "⌘",
    },
  };
  const regex = /option|mod/g;
  const replacedStr = str.replace(regex, (match) => {
    return list[match]["win"];
  });

  return replacedStr;
}

export function getShortcut(id, data) {
  const shortcut = defaultShortcut[id] || {};

  if (data) shortcut.data = data;
  if (!shortcut.readable) {
    shortcut.readable = getReadableShortcut(shortcut.combo);
  }

  return shortcut;
}

export function useShortcut(shortcuts, handler) {
  Mousetrap.prototype.stopCallback = () => false;

  const extractedShortcuts = {
    ids: {},
    keys: [],
    data: {},
  };
  const handleShortcut = (event, combo) => {
    const shortcutId = extractedShortcuts.ids[combo];
    const params = {
      event,
      ...extractedShortcuts.data[shortcutId],
    };

    if (shortcutId) event.preventDefault();

    if (typeof params.data === "function") {
      params.data(params, event);
    } else if (handler) {
      handler(params, event);
    }
  };
  const addShortcutData = ({ combo, id, readable, ...rest }) => {
    extractedShortcuts.ids[combo] = id;
    extractedShortcuts.keys.push(combo);
    extractedShortcuts.data[id] = { combo, id, readable, ...rest };
  };

  if (isObject(shortcuts)) {
    addShortcutData(getShortcut(shortcuts.id, shortcuts.data));
  } else if (typeof shortcuts === "string") {
    addShortcutData(getShortcut(shortcuts));
  } else {
    shortcuts.forEach((item) => {
      const currentShortcut =
        typeof item === "string" ? getShortcut(item) : item;

      addShortcutData(currentShortcut);
    });
  }

  onMounted(() => {
    Mousetrap.bind(extractedShortcuts.keys, handleShortcut);
  });
  onUnmounted(() => {
    Mousetrap.unbind(extractedShortcuts.keys);
  });

  return extractedShortcuts.data;
}
