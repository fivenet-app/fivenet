import type { MaybeRefOrGetter } from 'vue';

type UseSnapshotChangesOptions<T> = {
    serializer?: (value: T) => string;
    dirty?: MaybeRefOrGetter<boolean>;
    title?: string;
    description?: string;
    confirmationGroup?: symbol;
};

export function useSnapshotChanges<T>(source: MaybeRefOrGetter<T>, options: UseSnapshotChangesOptions<T> = {}) {
    const serialize = options.serializer ?? ((value: T) => JSON.stringify(value));
    const baselineSnapshot = ref(serialize(toValue(source)));
    const currentSnapshot = computed(() => serialize(toValue(source)));
    const snapshotDirty = computed(() => currentSnapshot.value !== baselineSnapshot.value);
    const extraDirty = computed(() => Boolean(toValue(options.dirty ?? false)));
    const isDirty = computed(() => snapshotDirty.value || extraDirty.value);

    const { hasUnsavedChanges, confirmLeave } = useUnsavedChanges({
        title: options.title,
        description: options.description,
        dirty: isDirty,
        confirmationGroup: options.confirmationGroup,
    });

    function syncSnapshot(value: T = toValue(source)): void {
        baselineSnapshot.value = serialize(value);
    }

    return {
        hasUnsavedChanges,
        snapshotDirty,
        isDirty,
        syncSnapshot,
        confirmLeave,
    };
}
