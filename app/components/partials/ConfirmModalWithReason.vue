<script lang="ts" setup>
import type { ButtonProps, FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';

const props = withDefaults(
    defineProps<{
        title?: string;
        description?: string;
        cancel?: () => Promise<unknown> | unknown;
        confirm: (reason: string) => Promise<unknown>;
        icon?: string;
        color?: ButtonProps['color'];
        iconClass?: string;
    }>(),
    {
        title: undefined,
        description: undefined,
        cancel: undefined,
        icon: 'i-mdi-warning-circle',
        color: 'error',
        iconClass: 'text-red-500 dark:text-red-400',
    },
);

const emit = defineEmits<{
    (e: 'close'): void;
}>();

const schema = z.object({
    reason: z.string().min(3).max(255),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    reason: '',
});

const canSubmit = ref(true);

const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await props.confirm(event.data.reason).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
    emit('close');
}, 1000);
</script>

<template>
    <UModal
        :title="title ?? $t('components.partials.confirm_dialog.title')"
        :description="description ?? $t('components.partials.confirm_dialog.description')"
        @update:model-value="cancel && cancel()"
    >
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UFormField class="sm:px-4" name="reason" :label="$t('common.reason')">
                <UInput v-model="state.reason" :placeholder="$t('common.reason')" :ui="{ base: 'w-full' }" />
            </UFormField>

            <div class="flex shrink-0 items-center gap-x-1.5 px-4 py-4 sm:px-4">
                <UButton type="submit" :color="color" :label="$t('common.confirm')" />
                <UButton
                    color="neutral"
                    :label="$t('common.cancel')"
                    @click="
                        if (cancel) {
                            cancel();
                        }
                        $emit('close');
                    "
                />
            </div>
        </UForm>
    </UModal>
</template>
