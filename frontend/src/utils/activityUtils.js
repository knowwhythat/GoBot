import { nanoid } from "nanoid";
export function deleteSelected(children, id) {
  const index = children.findIndex(
    (item) => item.id === id && !item.hasOwnProperty("deleteable")
  );
  if (index != -1) {
    children.splice(index, 1);
  } else {
    children.forEach((ele) => {
      if (ele.children) {
        deleteSelected(ele.children, id);
      }
    });
  }
}

export function cutSelected(children, id, cutBlock) {
  const index = children.findIndex(
    (item) => item.id === id && !item.hasOwnProperty("deleteable")
  );
  if (index != -1) {
    cutBlock.push(children.splice(index, 1)[0]);
  } else {
    children.forEach((ele) => {
      if (ele.children) {
        cutSelected(ele.children, id, cutBlock);
      }
    });
  }
}

export function copySelected(children, id, copyBlock) {
  const block = children.find(
    (item) => item.id === id && !item.hasOwnProperty("deleteable")
  );
  if (block) {
    const rawBlock = JSON.parse(JSON.stringify(block));
    copyBlock.push({ ...rawBlock, id: nanoid(16) });
  } else {
    children.forEach((ele) => {
      if (ele.children) {
        copySelected(ele.children, id, copyBlock);
      }
    });
  }
}

export function innerPaste(children, id, copiedBlocks) {
  const index = children.findIndex(
    (item) => item.id === id && !item.hasOwnProperty("deleteable")
  );
  if (index != -1) {
    copiedBlocks.forEach((block, innerIndex) => {
      children.splice(index + innerIndex + 1, 0, { ...block, id: nanoid(16) });
    });
  } else {
    children.forEach((ele) => {
      if (ele.children) {
        innerPaste(ele.children, id, copiedBlocks);
      }
    });
  }
}
