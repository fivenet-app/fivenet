type ReorderListRef<T> = Ref<T[]> | Ref<T[] | undefined>;

export function useListReorder<T>(list: ReorderListRef<T>, options?: { onMove?: () => void }) {
    function move(from: number, to: number) {
        const items = list.value;
        if (!items) return;
        if (from === to) return;
        if (from < 0 || from >= items.length) return;
        if (to < 0 || to >= items.length) return;

        const temp = items[from];
        items[from] = items[to]!;
        items[to] = temp!;
        options?.onMove?.();
    }

    function moveUp(index: number) {
        if (index <= 0) return; // top item, do nothing
        move(index, index - 1);
    }

    function moveDown(index: number) {
        if (!list.value || index >= list.value.length - 1) return; // bottom item, do nothing
        move(index, index + 1);
    }

    return { move, moveUp, moveDown };
}
