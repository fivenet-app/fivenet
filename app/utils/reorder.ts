export type NeighborMovePayload = {
    id: number;
    beforeId?: number;
    afterId?: number;
};

export function resolveNeighborMovePayload<T extends { id: number }>(
    entries: T[],
    oldIndex: number | undefined,
    newIndex: number | undefined,
): NeighborMovePayload | undefined {
    if (oldIndex === undefined || newIndex === undefined || oldIndex === newIndex) return undefined;

    const entry = entries[newIndex];
    if (!entry) return undefined;

    if (newIndex < oldIndex) {
        const beforeId = entries[newIndex + 1]?.id;
        return beforeId ? { id: entry.id, beforeId } : undefined;
    }

    const afterId = entries[newIndex - 1]?.id;
    if (!afterId) return undefined;

    return { id: entry.id, afterId };
}
