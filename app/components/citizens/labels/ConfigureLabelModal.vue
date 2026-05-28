<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { addSeconds } from 'date-fns';
import { z } from 'zod';
import InputDatePicker from '~/components/partials/InputDatePicker.vue';
import type { Label } from '~~/gen/ts/resources/citizens/labels/labels';
import LabelBadge from './LabelBadge.vue';

const props = defineProps<{
    label: Label;
}>();

const emits = defineEmits<{
    (e: 'close', v: Label | undefined): void;
}>();

const { t } = useI18n();

const schema = z.object({
    id: z.coerce.number(),
    sortOrder: z.number().min(0).default(0),
    name: z.coerce.string().min(1),
    color: z.coerce.string().length(7),
    icon: z.coerce.string().max(255).optional(),
    expiresAt: z
        .date()
        .min(new Date())
        .optional()
        .superRefine((data, ctx) => {
            if (props.label.settings?.requiresExpiration && !data) {
                ctx.addIssue({
                    code: 'custom',
                    message: t('common.required'),
                });
            }
        }),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    id: props.label.id,
    sortOrder: props.label.sortOrder,
    name: props.label.name,
    color: props.label.color,
    icon: props.label.icon,
    expiresAt: props.label.expiresAt ? toDate(props.label.expiresAt) : undefined,
});

const minExpiresAt = computed(() =>
    dateToCalendarDateTime(
        props.label.settings?.minDuration
            ? addSeconds(new Date(), durationToSeconds(props.label.settings.minDuration))
            : new Date(),
    ),
);

const maxExpiresAt = computed(() =>
    dateToCalendarDateTime(
        props.label.settings?.maxDuration
            ? addSeconds(new Date(), durationToSeconds(props.label.settings.maxDuration))
            : undefined,
    ),
);

const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    emits('close', {
        id: event.data.id,
        sortOrder: event.data.sortOrder,
        name: event.data.name,
        color: event.data.color,
        icon: event.data.icon,
        expiresAt: event.data.expiresAt ? toTimestamp(event.data.expiresAt) : undefined,
        settings: props.label.settings,
    });
}, 1000);

const formatDuration = useDurationFormatter();

const formRef = useTemplateRef('formRef');

watch(formRef, () => {
    if (props.label.settings?.requiresExpiration) formRef.value?.submit();
});
</script>

<template>
    <UModal>
        <template #title>
            <div class="inline-flex gap-2">
                <span>{{ $t('common.label') }}</span>

                <LabelBadge :label="label" />
            </div>
        </template>

        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <UFormField
                    name="expiresAt"
                    :label="$t('common.expires_at')"
                    :required="label.settings?.requiresExpiration"
                    :ui="{ container: 'flex flex-col gap-2' }"
                >
                    <InputDatePicker
                        v-model="state.expiresAt"
                        :clearable="!label.settings?.requiresExpiration"
                        time
                        :min-value="minExpiresAt"
                        :max-value="maxExpiresAt"
                    />

                    <div
                        v-if="label.settings?.minDuration || label.settings?.maxDuration"
                        class="inline-flex w-full flex-1 gap-2"
                    >
                        <UBadge
                            v-if="label.settings?.minDuration"
                            :label="`${$t('components.citizens.citizen_labels.settings.min_duration')}: ${formatDuration(label.settings?.minDuration)}`"
                        />

                        <UBadge
                            v-if="label.settings?.maxDuration"
                            :label="`${$t('components.citizens.citizen_labels.settings.max_duration')}: ${formatDuration(label.settings?.maxDuration)}`"
                        />
                    </div>
                </UFormField>
            </UForm>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton
                    class="flex-1"
                    color="neutral"
                    block
                    :label="$t('common.close', 1)"
                    @click="$emit('close', undefined)"
                />

                <UButton class="flex-1" block :label="$t('common.save')" @click="formRef?.submit()" />
            </UFieldGroup>
        </template>
    </UModal>
</template>
