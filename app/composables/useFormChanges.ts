import type { Ref } from 'vue';

type FormInstance = {
    dirty?: boolean;
    clear?: () => Promise<void> | void;
    reset?: () => Promise<void> | void;
};

type UseFormChangesOptions = {
    title?: string;
    description?: string;
    confirmationGroup?: symbol;
};

const _useFormChanges = (formRef: Ref<FormInstance | null>, options: UseFormChangesOptions = {}) => {
    const formDirty = computed(() => Boolean(formRef.value?.dirty));
    const { changed, hasUnsavedChanges, markChanged, resetChanged } = useUnsavedChanges({
        title: options.title,
        description: options.description,
        dirty: formDirty,
        confirmationGroup: options.confirmationGroup,
    });

    async function resetForm(): Promise<void> {
        await formRef.value?.reset?.();
        await formRef.value?.clear?.();
        resetChanged();
    }

    return {
        changed,
        hasUnsavedChanges,
        markChanged,
        resetChanged,
        resetForm,
    };
};

export const useFormChanges = _useFormChanges;
