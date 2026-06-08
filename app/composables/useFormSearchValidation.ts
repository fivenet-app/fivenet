import type { Form, FormData, FormSchema } from '@nuxt/ui';
import { watchDebounced } from '@vueuse/shared';
import type { MaybeRefOrGetter, Ref } from 'vue';

function cloneSnapshot<T>(value: T): T {
    if (typeof value !== 'object' || value === null) return value;

    return structuredClone(toRaw(value as object)) as T;
}

export function useFormSearchValidation<TSchema extends FormSchema>(
    source: MaybeRefOrGetter<FormData<TSchema>>,
    formRef: MaybeRefOrGetter<Form<TSchema> | null | undefined>,
    options: {
        debounce?: number;
        maxWait?: number;
    } = {},
): {
    validatedQuery: Ref<FormData<TSchema>>;
    commitValidatedQuery: () => Promise<void>;
} {
    const { debounce = 200, maxWait = 1250 } = options;

    const validatedQuery = shallowRef<FormData<TSchema>>(cloneSnapshot(toValue(source)) as FormData<TSchema>);

    async function commitValidatedQuery(): Promise<void> {
        const form = toValue(formRef);
        if (!form) return;

        try {
            const valid = await form.validate({});
            if (!valid) return;

            validatedQuery.value = valid;
        } catch {
            return;
        }
    }

    watchDebounced(
        () => toValue(source),
        async () => {
            await commitValidatedQuery();
        },
        {
            debounce,
            maxWait,
            deep: true,
        },
    );

    return {
        validatedQuery,
        commitValidatedQuery,
    };
}
