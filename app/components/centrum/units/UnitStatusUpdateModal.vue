<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import { unitStatusToBadgeColor, unitStatuses } from '~/components/centrum/helpers';
import { useCentrumStore } from '~/stores/centrum';
import { getCentrumCentrumClient } from '~~/gen/ts/clients';
import { type Unit, StatusUnit } from '~~/gen/ts/resources/centrum/units';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { Coordinate } from '~~/shared/types/types';

const props = defineProps<{
    unit: Unit;
    status?: StatusUnit;
    location?: Coordinate;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const centrumStore = useCentrumStore();
const { settings } = storeToRefs(centrumStore);

const notifications = useNotificationsStore();

const centrumCentrumClient = await getCentrumCentrumClient();

const schema = z.object({
    status: z.enum(StatusUnit),
    code: z.union([z.coerce.string().min(1).max(20), z.coerce.string().length(0).optional()]),
    reason: z.union([z.coerce.string().min(3).max(255), z.coerce.string().length(0).optional()]),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    status: props.status ?? props.unit?.status?.status ?? StatusUnit.UNKNOWN,
});

async function updateUnitStatus(id: number, values: Schema): Promise<void> {
    try {
        const call = centrumCentrumClient.updateUnitStatus({
            unitId: id,
            status: values.status,
            code: values.code,
            reason: values.reason,
        });
        await call;

        notifications.add({
            title: { key: 'notifications.centrum.sidebar.unit_status_updated.title', parameters: {} },
            description: { key: 'notifications.centrum.sidebar.unit_status_updated.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emit('close', false);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await updateUnitStatus(props.unit.id, event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

watch(props, () => (state.status = props.status ?? props.unit?.status?.status ?? StatusUnit.UNKNOWN));

function updateReasonField(value: string): void {
    if (value.length === 0) return;

    state.reason = value;
}

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UModal :title="`${$t('components.centrum.update_unit_status.title')}: ${unit.name} (${unit.initials})`">
        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <dl class="divide-y divide-default">
                    <div class="flex flex-col gap-4 px-4 py-3 sm:px-0">
                        <UFormField
                            class="grid grid-cols-2 items-center gap-2"
                            name="status"
                            :label="$t('common.status')"
                            required
                        >
                            <div class="grid w-full grid-cols-2 gap-0.5">
                                <UButton
                                    v-for="item in unitStatuses"
                                    :key="item.name"
                                    class="group my-0.5 flex w-full flex-col items-center rounded-md p-1.5 text-xs font-medium hover:transition-all"
                                    :class="state.status == item.status ? 'bg-neutral-500 hover:bg-neutral-400' : ''"
                                    :color="unitStatusToBadgeColor(item.status)"
                                    :disabled="state.status == item.status"
                                    :icon="item.icon"
                                    :label="
                                        item.status
                                            ? $t(`enums.centrum.StatusUnit.${StatusUnit[item.status ?? 0]}`)
                                            : $t(item.name)
                                    "
                                    @click="state.status = item.status ?? StatusUnit.UNAVAILABLE"
                                />
                            </div>
                        </UFormField>

                        <UFormField class="grid grid-cols-2 items-center gap-2" name="code" :label="$t('common.code')">
                            <UInput
                                v-model="state.code"
                                type="text"
                                class="w-full"
                                name="code"
                                :placeholder="$t('common.code')"
                                :label="$t('common.code')"
                            />
                        </UFormField>

                        <UFormField
                            name="reason"
                            class="grid grid-cols-2 items-center gap-2"
                            :label="$t('common.reason')"
                            required
                        >
                            <UInput v-model="state.reason" type="text" :placeholder="$t('common.reason')" class="w-full" />
                        </UFormField>

                        <UFormField
                            v-if="settings?.predefinedStatus && settings?.predefinedStatus.unitStatus.length > 0"
                            name="unitStatus"
                            class="grid grid-cols-2 items-center gap-2"
                            :label="`${$t('common.predefined', 2)} ${$t('common.reason', 2)}`"
                        >
                            <ClientOnly>
                                <USelectMenu
                                    name="unitStatus"
                                    :items="['&nbsp;', ...settings?.predefinedStatus.unitStatus]"
                                    :search-input="{ placeholder: $t('common.search_field') }"
                                    @update:model-value="($event) => updateReasonField($event)"
                                >
                                    <template #item-label="{ item }">
                                        {{ item !== '' ? item : '&nbsp;' }}
                                    </template>
                                </USelectMenu>
                            </ClientOnly>
                        </UFormField>
                    </div>
                </dl>
            </UForm>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />

                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="$t('common.update')"
                    @click="formRef?.submit()"
                />
            </UFieldGroup>
        </template>
    </UModal>
</template>
