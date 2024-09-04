<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { unitStatusToBGColor, unitStatuses } from '~/components/centrum/helpers';
import { useCentrumStore } from '~/store/centrum';
import { useNotificatorStore } from '~/store/notificator';
import { StatusUnit, Unit } from '~~/gen/ts/resources/centrum/units';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    unit: Unit;
    status?: StatusUnit;
    location?: Coordinate;
}>();

const { isOpen } = useModal();

const centrumStore = useCentrumStore();
const { settings } = storeToRefs(centrumStore);

const notifications = useNotificatorStore();

const schema = z.object({
    status: z.nativeEnum(StatusUnit),
    code: z.union([z.string().min(1).max(20), z.string().length(0).optional()]),
    reason: z.union([z.string().min(3).max(255), z.string().length(0).optional()]),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    status: props.status ?? props.unit?.status?.status ?? StatusUnit.UNKNOWN,
});

interface FormData {
    status: number;
    code?: string;
    reason?: string;
}

async function updateUnitStatus(id: string, values: FormData): Promise<void> {
    try {
        const call = getGRPCCentrumClient().updateUnitStatus({
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

        isOpen.value = false;
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
    if (value.length === 0) {
        return;
    }

    state.reason = value;
}
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard
                class="flex flex-1 flex-col"
                :ui="{
                    body: {
                        padding: 'px-1 py-2 sm:p-2',
                    },
                    ring: '',
                    divide: 'divide-y divide-gray-100 dark:divide-gray-800',
                }"
            >
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ $t('components.centrum.update_unit_status.title') }}: {{ unit.name }} ({{ unit.initials }})
                        </h3>

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <dl class="divide-neutral/10 divide-y">
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
                                <label for="status" class="block text-sm font-medium leading-6">
                                    {{ $t('common.status') }}
                                </label>
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <UFormGroup name="status">
                                    <div class="grid w-full grid-cols-2 gap-0.5">
                                        <UButton
                                            v-for="item in unitStatuses"
                                            :key="item.name"
                                            class="hover:bg-primary-100/10 group my-0.5 flex w-full flex-col items-center rounded-md p-1.5 text-xs font-medium hover:transition-all"
                                            :class="[
                                                state.status == item.status
                                                    ? 'bg-base-500 hover:bg-base-400'
                                                    : item.status
                                                      ? unitStatusToBGColor(item.status)
                                                      : '',
                                                ,
                                            ]"
                                            :disabled="state.status == item.status"
                                            @click="state.status = item.status ?? StatusUnit.UNAVAILABLE"
                                        >
                                            <UIcon :name="item.icon" class="size-5 shrink-0" />
                                            <span class="mt-1">
                                                {{
                                                    item.status
                                                        ? $t(`enums.centrum.StatusUnit.${StatusUnit[item.status ?? 0]}`)
                                                        : $t(item.name)
                                                }}
                                            </span>
                                        </UButton>
                                    </div>
                                </UFormGroup>
                            </dd>
                        </div>
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
                                <label for="code" class="block text-sm font-medium leading-6">
                                    {{ $t('common.code') }}
                                </label>
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <UFormGroup name="code" class="flex-1">
                                    <UInput
                                        v-model="state.code"
                                        type="text"
                                        name="code"
                                        :placeholder="$t('common.code')"
                                        :label="$t('common.code')"

                                        
                                    />
                                </UFormGroup>
                            </dd>
                        </div>
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
                                <label for="reason" class="block text-sm font-medium leading-6">
                                    {{ $t('common.reason') }}
                                </label>
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <UFormGroup name="reason" class="flex-1">
                                    <UInput
                                        v-model="state.reason"
                                        type="text"
                                        name="reason"
                                        :placeholder="$t('common.reason')"

                                        
                                    />
                                </UFormGroup>
                            </dd>
                        </div>
                        <div
                            v-if="settings?.predefinedStatus && settings?.predefinedStatus.unitStatus.length > 0"
                            class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0"
                        >
                            <dt class="text-sm font-medium leading-6">
                                <label for="unitStatus" class="block text-sm font-medium leading-6">
                                    {{ $t('common.predefined', 2) }}
                                    {{ $t('common.reason', 2) }}
                                </label>
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <USelectMenu
                                    name="unitStatus"
                                    :options="['&nbsp;', ...settings?.predefinedStatus.unitStatus]"
                                    :searchable-placeholder="$t('common.search_field')"

                                    
                                    @change="updateReasonField($event)"
                                >
                                    <template #option="{ option }">
                                        <span class="truncate">
                                            {{ option !== '' ? option : '&nbsp;' }}
                                        </span>
                                    </template>
                                </USelectMenu>
                            </dd>
                        </div>
                    </dl>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton color="black" block class="flex-1" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton type="submit" block class="flex-1" :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.update') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
