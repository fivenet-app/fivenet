<script lang="ts" setup>
import type { ButtonColor, FormSubmitEvent } from '#ui/types';
import { z } from 'zod';

const props = withDefaults(
    defineProps<{
        title?: string;
        description?: string;
        cancel?: () => Promise<unknown> | unknown;
        confirm: (reason: string) => Promise<unknown>;
        icon?: string;
        color?: ButtonColor;
        iconClass?: string;
    }>(),
    {
        title: undefined,
        description: undefined,
        cancel: undefined,
        icon: 'i-mdi-warning-circle',
        color: 'red',
        iconClass: 'text-red-500 dark:text-red-400',
    },
);

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
    isOpen.value = false;
}, 1000);

const { isOpen } = useModal();
</script>

<template>
    <UDashboardModal
        :title="title ?? $t('components.partials.confirm_dialog.title')"
        :description="description ?? $t('components.partials.confirm_dialog.description')"
        :icon="icon"
        :ui="{
            icon: { base: iconClass },
            body: { base: 'sm:p-0 sm:px-6' },
        }"
        @update:model-value="cancel && cancel()"
    >
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UFormGroup name="reason" :label="$t('common.reason')" class="sm:px-4">
                <UInput v-model="state.reason" :placeholder="$t('common.reason')" :ui="{ base: 'w-full' }" />
            </UFormGroup>

            <div class="flex flex-shrink-0 items-center gap-x-1.5 px-4 py-4 sm:px-4">
                <UButton type="submit" :color="color" :label="$t('common.confirm')" />
                <UButton
                    color="white"
                    :label="$t('common.cancel')"
                    @click="
                        if (cancel) {
                            cancel();
                        }
                        isOpen = false;
                    "
                />
            </div>
        </UForm>
    </UDashboardModal>
</template>
