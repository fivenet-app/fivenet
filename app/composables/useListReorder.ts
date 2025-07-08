export function useListReorder<T>(list: Ref<T[]>) {
    function move(from: number, to: number) {
        if (from === to) return;
        if (from < 0 || from >= list.value.length) return;
        if (to < 0 || to >= list.value.length) return;

        const temp = list.value[from];
        list.value[from] = list.value[to]!;
        list.value[to] = temp!;
    }

    function moveUp(index: number) {
        if (index <= 0) return; // top item, do nothing
        move(index, index - 1);
    }

    function moveDown(index: number) {
        if (index >= list.value.length - 1) return; // bottom item, do nothing
        move(index, index + 1);
    }

    return { move, moveUp, moveDown };
}
