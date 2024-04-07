<script lang="ts" setup>
import { z } from 'zod';
import type { FormSubmitEvent } from '#ui/types';
import { CentrumMode, Settings } from '~~/gen/ts/resources/centrum/settings';

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const { data: settings, refresh } = useLazyAsyncData('rector-centrum-settings', () => getCentrumSettings());

async function getCentrumSettings(): Promise<Settings> {
    try {
        const call = $grpc.getCentrumClient().getSettings({});
        const { response } = await call;

        return response.settings!;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const modes = ref<{ mode: CentrumMode; selected?: boolean }[]>([
    { mode: CentrumMode.MANUAL },
    { mode: CentrumMode.SIMPLIFIED },
    { mode: CentrumMode.CENTRAL_COMMAND },
    { mode: CentrumMode.AUTO_ROUND_ROBIN },
]);

const schema = z.object({
    enabled: z.boolean(),
    mode: z.custom<CentrumMode>(),
    fallbackMode: z.custom<CentrumMode>(),
    unitStatus: z.string().array().max(10),
    dispatchStatus: z.string().array().max(10),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    enabled: false,
    mode: CentrumMode.MANUAL,
    fallbackMode: CentrumMode.AUTO_ROUND_ROBIN,
    unitStatus: [],
    dispatchStatus: [],
});

async function updateSettings(values: Schema): Promise<void> {
    try {
        const call = $grpc.getCentrumClient().updateSettings({
            settings: {
                job: '',
                enabled: values.enabled,
                mode: values.mode,
                fallbackMode: values.fallbackMode,
                predefinedStatus: {
                    dispatchStatus: values.dispatchStatus,
                    unitStatus: values.unitStatus,
                },
            },
        });
        await call;

        refresh();

        isOpen.value = false;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

function setSettingsValues(): void {
    if (!settings.value) {
        return;
    }

    state.enabled = settings.value.enabled;
    state.mode = settings.value.mode;
    state.fallbackMode = settings.value.fallbackMode;
    state.dispatchStatus = settings.value.predefinedStatus?.dispatchStatus ?? [];
    state.unitStatus = settings.value.predefinedStatus?.unitStatus ?? [];
}

watch(settings, () => setSettingsValues());

setSettingsValues();

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await updateSettings(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ $t('components.centrum.units.update_settings') }}
                        </h3>

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <div class="text-sm text-gray-100">
                        <UFormGroup name="enabled" :label="$t('common.enabled')" class="flex-1">
                            <UCheckbox
                                v-model="state.enabled"
                                name="enabled"
                                :disabled="!can('SuperUser')"
                                :placeholder="$t('common.enabled')"
                            />
                        </UFormGroup>

                        <UFormGroup name="mode" :label="$t('common.mode')" class="flex-1">
                            <USelectMenu
                                v-model="state.mode"
                                :options="modes"
                                value-attribute="mode"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            >
                                <template #label>
                                    <span class="truncate">{{
                                        $t(`enums.centrum.CentrumMode.${CentrumMode[state.mode ?? 0]}`)
                                    }}</span>
                                </template>
                                <template #option="{ option }">
                                    <span class="truncate">{{
                                        $t(`enums.centrum.CentrumMode.${CentrumMode[option.mode ?? 0]}`)
                                    }}</span>
                                </template>
                            </USelectMenu>
                        </UFormGroup>

                        <UFormGroup name="fallbackMode" :label="$t('common.fallback_mode')" class="flex-1">
                            <USelectMenu
                                v-model="state.fallbackMode"
                                :options="modes"
                                value-attribute="mode"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            >
                                <template #label>
                                    <span class="truncate">{{
                                        $t(`enums.centrum.CentrumMode.${CentrumMode[state.mode ?? 0]}`)
                                    }}</span>
                                </template>
                                <template #option="{ option }">
                                    <span class="truncate">{{
                                        $t(`enums.centrum.CentrumMode.${CentrumMode[option.mode ?? 0]}`)
                                    }}</span>
                                </template>
                            </USelectMenu>
                        </UFormGroup>

                        <!-- Predefined Unit Status Reason -->
                        <UFormGroup name="unitStatus" :label="`${$t('common.units')} ${$t('common.status')}`" class="flex-1">
                            <div class="flex flex-col gap-1">
                                <div v-for="(_, idx) in state.unitStatus" :key="idx" class="flex items-center gap-1">
                                    <UInput
                                        type="text"
                                        class="w-full flex-1"
                                        :placeholder="$t('common.reason')"
                                        :label="$t('common.reason')"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />

                                    <UButton
                                        :ui="{ rounded: 'rounded-full' }"
                                        icon="i-mdi-close"
                                        @click="state.unitStatus.splice(idx, 1)"
                                    />
                                </div>
                            </div>

                            <UButton
                                :ui="{ rounded: 'rounded-full' }"
                                icon="i-mdi-plus"
                                :disabled="!canSubmit || state.unitStatus.length >= 8"
                                @click="state.unitStatus.push('')"
                            />
                        </UFormGroup>

                        <!-- Predefined Dispatch Status Reason -->
                        <UFormGroup
                            name="dispatchStatus"
                            :label="`${$t('common.dispatches')} ${$t('common.status')}`"
                            class="flex-1"
                        >
                            <div class="flex flex-col gap-1">
                                <div v-for="(_, idx) in state.dispatchStatus" :key="idx" class="flex items-center gap-1">
                                    <UInput
                                        v-model="state.dispatchStatus[idx]"
                                        type="text"
                                        class="w-full flex-1"
                                        :placeholder="$t('common.reason')"
                                        :label="$t('common.reason')"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />

                                    <UButton
                                        :ui="{ rounded: 'rounded-full' }"
                                        icon="i-mdi-close"
                                        @click="state.dispatchStatus.splice(idx, 1)"
                                    />
                                </div>
                            </div>

                            <UButton
                                :ui="{ rounded: 'rounded-full' }"
                                icon="i-mdi-plus"
                                :disabled="!canSubmit || state.dispatchStatus.length >= 8"
                                @click="state.dispatchStatus.push('')"
                            />
                        </UFormGroup>
                    </div>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton color="black" block class="flex-1" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>
                        <UButton
                            v-if="can('CentrumService.UpdateSettings')"
                            type="submit"
                            block
                            class="flex-1"
                            :disabled="!canSubmit"
                            :loading="!canSubmit"
                        >
                            {{ $t('common.update') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
