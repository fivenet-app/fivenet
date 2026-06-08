import type { MaybeRefOrGetter } from 'vue';
import FormChangedModal from '~/components/partials/FormChangedModal.vue';
import { DEFAULT_UNSAVED_CHANGES_CONFIRMATION_GROUP, runSharedUnsavedChangesConfirmation } from '~/utils/unsavedChanges';

type UseUnsavedChangesOptions = {
    title?: string;
    description?: string;
    dirty?: MaybeRefOrGetter<boolean>;
    confirmationGroup?: symbol;
};

export function useUnsavedChanges(options: UseUnsavedChangesOptions = {}) {
    const overlay = useOverlay();
    const formChangedModal = overlay.create(FormChangedModal);
    const changed = ref<boolean>(false);
    let stopBeforeUnloadWatch: (() => void) | undefined;

    const hasUnsavedChanges = computed(() => Boolean(changed.value || toValue(options.dirty ?? false)));
    const handleBeforeUnload = (event: BeforeUnloadEvent) => {
        if (!hasUnsavedChanges.value) return;

        event.preventDefault();
        event.returnValue = '';
    };

    function markChanged(): void {
        changed.value = true;
    }

    function resetChanged(): void {
        changed.value = false;
    }

    async function confirmLeave(): Promise<boolean> {
        return runSharedUnsavedChangesConfirmation(
            options.confirmationGroup ?? DEFAULT_UNSAVED_CHANGES_CONFIRMATION_GROUP,
            async () => {
                const response = await formChangedModal.open({
                    title: options.title ?? undefined,
                    description: options.description ?? undefined,
                });

                return response === true;
            },
        );
    }

    onMounted(() => {
        stopBeforeUnloadWatch = watch(
            hasUnsavedChanges,
            (isDirty) => {
                if (isDirty) {
                    window.addEventListener('beforeunload', handleBeforeUnload);
                } else {
                    window.removeEventListener('beforeunload', handleBeforeUnload);
                }
            },
            { immediate: true },
        );
    });

    onBeforeUnmount(() => {
        stopBeforeUnloadWatch?.();
        window.removeEventListener('beforeunload', handleBeforeUnload);
    });

    onBeforeRouteLeave(async () => {
        if (!hasUnsavedChanges.value) return;

        const confirmed = await confirmLeave();
        if (!confirmed) return false;
    });

    return {
        changed,
        hasUnsavedChanges,
        markChanged,
        resetChanged,
        confirmLeave,
    };
}
