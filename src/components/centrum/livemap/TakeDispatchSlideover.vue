<script lang="ts" setup>
import { useSound } from '@raffaelesgarro/vue-use-sound';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import { useCentrumStore } from '~/store/centrum';
import { Dispatch, StatusDispatch, TakeDispatchResp } from '~~/gen/ts/resources/centrum/dispatches';
import { CentrumMode } from '~~/gen/ts/resources/centrum/settings';
import TakeDispatchEntry from '~/components/centrum/livemap/TakeDispatchEntry.vue';
import { isStatusDispatchCompleted } from '~/components/centrum/helpers';
import { useSettingsStore } from '~/store/settings';

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useSlideover();

const centrumStore = useCentrumStore();
const { dispatches, pendingDispatches, getCurrentMode } = storeToRefs(centrumStore);

const settingsStore = useSettingsStore();
const { audio: audioSettings } = storeToRefs(settingsStore);

const selectedDispatches = ref<string[]>([]);
const queryDispatches = ref('');

async function takeDispatches(resp: TakeDispatchResp): Promise<void> {
    try {
        if (!canTakeDispatch.value) {
            return;
        }

        const dispatchIds = selectedDispatches.value.filter((sd) => {
            const dsp = dispatches.value.get(sd);
            if (dsp === undefined) {
                return false;
            }

            // Dispatch has no status or is not completed?
            return dsp.status === undefined ? true : !isStatusDispatchCompleted(dsp.status.status);
        });

        if (dispatchIds.length === 0) {
            return;
        }

        // Make sure all selected dispatches are still existing and not in a "completed"
        const call = $grpc.getCentrumClient().takeDispatch({
            dispatchIds,
            resp,
        });
        await call;

        selectedDispatches.value.length = 0;

        isOpen.value = false;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

function selectDispatch(id: string, state: boolean): void {
    const idx = selectedDispatches.value.findIndex((did) => did === id);
    if (idx > -1 && !state) {
        selectedDispatches.value.splice(idx, 1);
    } else if (idx === -1 && state) {
        selectedDispatches.value.push(id);
    }
}

const newDispatchSound = useSound('/sounds/centrum/message-incoming.mp3', {
    volume: audioSettings.value.notificationsVolume,
});

const debouncedPlay = useDebounceFn(() => newDispatchSound.play(), 2000, { maxWait: 5000 });

const previousLength = ref(0);
watch(pendingDispatches.value, () => {
    if (getCurrentMode.value !== CentrumMode.SIMPLIFIED) {
        if (previousLength.value <= pendingDispatches.value.length && pendingDispatches.value.length !== 0) {
            debouncedPlay();
        }
    }

    previousLength.value = pendingDispatches.value.length;
});

const canTakeDispatch = computed(
    () =>
        selectedDispatches.value.length > 0 &&
        (pendingDispatches.value.length > 0 || (getCurrentMode.value === CentrumMode.SIMPLIFIED && dispatches.value.size > 0)),
);

const filteredDispatches = computedAsync(async () => {
    const filtered: Dispatch[] = [];
    dispatches.value.forEach((d) => {
        if (d.id.includes(queryDispatches.value) || d.message.includes(queryDispatches.value)) {
            if (d.status === undefined || d.status.status < StatusDispatch.COMPLETED) filtered.push(d);
        }
    });
    return filtered.sort((a, b) => (a.status?.status ?? 0) - (b.status?.status ?? 0)).map((d) => d.id);
});

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (resp: TakeDispatchResp) => {
    canSubmit.value = false;
    await takeDispatches(resp).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <USlideover>
        <UCard
            class="flex flex-col flex-1"
            :ui="{ body: { base: 'flex-1' }, ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }"
        >
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('components.centrum.take_dispatch.title') }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <div>
                <div class="flex flex-1 flex-col justify-between">
                    <div class="divide-y divide-gray-200 px-2 sm:px-6">
                        <div class="mt-1">
                            <dl class="divide-y divide-neutral/10 border-b border-neutral/10">
                                <template v-if="getCurrentMode === CentrumMode.SIMPLIFIED">
                                    <DataNoDataBlock
                                        v-if="dispatches.size === 0"
                                        icon="i-mdi-car-emergency"
                                        :type="$t('common.dispatch', 2)"
                                    />
                                    <template v-else>
                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                            <dt class="text-sm font-medium leading-6">
                                                <div class="flex h-6 items-center">
                                                    {{ $t('common.search') }}
                                                </div>
                                            </dt>
                                            <dd class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0">
                                                <div class="relative flex items-center">
                                                    <UInput
                                                        v-model="queryDispatches"
                                                        type="text"
                                                        name="search"
                                                        :placeholder="$t('common.search')"
                                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 pr-14 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                        @focusin="focusTablet(true)"
                                                        @focusout="focusTablet(false)"
                                                    />
                                                </div>
                                            </dd>
                                        </div>

                                        <TakeDispatchEntry
                                            v-for="pd in filteredDispatches"
                                            :key="pd"
                                            :dispatch="dispatches.get(pd)!"
                                            :preselected="false"
                                            @selected="selectDispatch(pd, $event)"
                                            @goto="$emit('goto', $event)"
                                        />
                                    </template>
                                </template>

                                <template v-else>
                                    <DataNoDataBlock
                                        v-if="pendingDispatches.length === 0"
                                        icon="i-mdi-car-emergency"
                                        :type="$t('common.dispatch', 2)"
                                    />
                                    <template v-else>
                                        <TakeDispatchEntry
                                            v-for="pd in pendingDispatches"
                                            :key="pd"
                                            :dispatch="dispatches.get(pd)!"
                                            @selected="selectDispatch(pd, $event)"
                                            @goto="$emit('goto', $event)"
                                        />
                                    </template>
                                </template>
                            </dl>
                        </div>
                    </div>
                </div>
            </div>

            <template #footer>
                <div class="inline-flex w-full">
                    <UButton @click="isOpen = false">
                        {{ $t('common.close') }}
                    </UButton>
                    <UButton
                        class="relative inline-flex w-full items-center rounded-l-md bg-success-500 px-3 py-2 text-sm font-semibold ring-1 ring-inset ring-success-300 hover:bg-success-100"
                        :disabled="!canTakeDispatch || !canSubmit"
                        :loading="!canSubmit"
                        @click="onSubmitThrottle(TakeDispatchResp.ACCEPTED)"
                    >
                        {{ $t('common.accept') }}
                    </UButton>
                    <UButton
                        class="relative -ml-px inline-flex w-full items-center bg-error-500 px-3 py-2 text-sm font-semibold ring-1 ring-inset ring-error-300 hover:bg-error-100"
                        :disabled="!canTakeDispatch || !canSubmit"
                        :loading="!canSubmit"
                        @click="onSubmitThrottle(TakeDispatchResp.DECLINED)"
                    >
                        {{ $t('common.decline') }}
                    </UButton>
                </div>
            </template>
        </UCard>
    </USlideover>
</template>
