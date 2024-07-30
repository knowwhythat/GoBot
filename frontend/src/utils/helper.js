import { typeBuilder } from "@/utils/shared.js";
import { nanoid } from "nanoid";

export function isXPath(str) {
  const regex = /^([(/@]|id\()/;

  return regex.test(str);
}

export function visibleInViewport(element) {
  const { top, left, bottom, right, height, width } =
    element.getBoundingClientRect();

  if (height === 0 || width === 0) return false;

  return (
    top >= 0 &&
    left >= 0 &&
    bottom <= (window.innerHeight || document.documentElement.clientHeight) &&
    right <= (window.innerWidth || document.documentElement.clientWidth)
  );
}

export function sleep(timeout = 500) {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve();
    }, timeout);
  });
}

export function throttle(callback, limit) {
  let waiting = false;

  return (...args) => {
    if (!waiting) {
      callback.apply(this, args);
      waiting = true;
      setTimeout(() => {
        waiting = false;
      }, limit);
    }
  };
}

export function convertArrObjTo2DArr(arr) {
  const keyIndex = new Map();
  const values = [[]];

  arr.forEach((obj) => {
    const keys = Object.keys(obj);
    const row = [];

    keys.forEach((key) => {
      if (!keyIndex.has(key)) {
        keyIndex.set(key, keyIndex.size);
        values[0].push(key);
      }

      const value = obj[key];

      const rowIndex = keyIndex.get(key);
      row[rowIndex] = typeof value === "object" ? JSON.stringify(value) : value;
    });

    values.push([...row]);
  });

  return values;
}

export function convert2DArrayToArrayObj(values) {
  let keyIndex = 0;
  const keys = values.shift();
  const result = [];

  for (let columnIndex = 0; columnIndex < values.length; columnIndex += 1) {
    const currentColumn = {};

    for (
      let rowIndex = 0;
      rowIndex < values[columnIndex].length;
      rowIndex += 1
    ) {
      let key = keys[rowIndex];

      if (!key) {
        keyIndex += 1;
        key = `_row${keyIndex}`;
        keys.push(key);
      }

      currentColumn[key] = values[columnIndex][rowIndex];
    }

    result.push(currentColumn);
  }

  return result;
}

export function parseJSON(data, def) {
  try {
    return JSON.parse(data);
  } catch (error) {
    return def;
  }
}

export function parseFlow(flow) {
  return typeof flow === "string" ? parseJSON(flow, {}) : flow;
}

export function replaceMustache(str, replacer) {
  /* eslint-disable-next-line */
  return str.replace(/\{\{(.*?)}}/g, replacer);
}

export function openFilePicker(acceptedFileTypes = [], attrs = {}) {
  return new Promise((resolve) => {
    const input = document.createElement("input");
    input.type = "file";
    input.accept = Array.isArray(acceptedFileTypes)
      ? acceptedFileTypes.join(",")
      : acceptedFileTypes;

    Object.entries(attrs).forEach(([key, value]) => {
      input[key] = value;
    });

    input.onchange = (event) => {
      const { files } = event.target;
      const validFiles = [];

      Array.from(files).forEach((file) => {
        if (!acceptedFileTypes.includes(file.type)) return;

        validFiles.push(file);
      });

      resolve(validFiles);
    };

    input.click();
  });
}

export function fileSaver(filename, data) {
  const anchor = document.createElement("a");
  anchor.download = filename;
  anchor.href = data;

  anchor.dispatchEvent(new MouseEvent("click"));
  anchor.remove();
}

export function countDuration(started, ended) {
  const duration = Math.round((ended - started) / 1000);
  const minutes = parseInt((duration / 60) % 60, 10);
  const seconds = parseInt(duration % 60, 10);

  const getText = (num, suffix) => (num > 0 ? `${num}${suffix}` : "");

  return `${getText(minutes, "m")} ${seconds}s`;
}

export function toCamelCase(str, capitalize = false) {
  const result = str.replace(/(?:^\w|[A-Z]|\b\w)/g, (letter, index) => {
    return index === 0 && !capitalize
      ? letter.toLowerCase()
      : letter.toUpperCase();
  });

  return result.replace(/\s+|[-]/g, "");
}

export function isObject(obj) {
  return typeof obj === "object" && obj !== null && !Array.isArray(obj);
}

export function objectHasKey(obj, key) {
  return Object.prototype.hasOwnProperty.call(obj, key);
}

export function isWhitespace(str) {
  return !/\S/.test(str);
}

export function debounce(callback, time = 200) {
  let interval;

  return (...args) => {
    clearTimeout(interval);

    return new Promise((resolve) => {
      interval = setTimeout(() => {
        interval = null;

        callback(...args);
        resolve();
      }, time);
    });
  };
}

export function arraySorter({ data, key, order = "asc" }) {
  let runCounts = {};
  const copyData = data.slice();

  if (key === "mostUsed") {
    runCounts = parseJSON(localStorage.getItem("runCounts"), {}) || {};
  }

  return copyData.sort((a, b) => {
    let comparison = 0;
    let itemA = a[key] || a;
    let itemB = b[key] || b;

    if (key === "mostUsed") {
      itemA = runCounts[a.id] || 0;
      itemB = runCounts[b.id] || 0;
    }

    if (itemA > itemB) {
      comparison = 1;
    } else if (itemA < itemB) {
      comparison = -1;
    }

    return order === "desc" ? comparison * -1 : comparison;
  });
}

export function getIconPath(path) {
  if (path && path.startsWith("path")) {
    const { 1: iconPath } = path.split(":");
    return { viewBox: "0 0 1024 1024", path: iconPath };
  }

  return { name: path };
}

export function getVariableDefine(name, type, scope = "local") {
  const typeProperty =
    type in typeBuilder ? typeBuilder[type] : typeBuilder["default"];
  const children = [];
  let label = name;
  if (scope === "global") {
    label = `${name}(全局变量)`;
    name = `glv['${name}']`;
  }
  if (typeProperty["property"]) {
    typeProperty["property"].forEach((property) => {
      children.push({
        key: nanoid(8),
        label: property.label,
        name: property.name,
        type: property.type,
        parent: name,
      });
    });
  }
  return {
    key: nanoid(8),
    label: label,
    name: name,
    type: type,
    icon: typeProperty["icon"],
    children: children,
  };
}
