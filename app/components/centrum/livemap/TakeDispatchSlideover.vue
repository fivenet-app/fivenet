<script lang="ts" setup>
import { isStatusDispatchCompleted } from '~/components/centrum/helpers';
import TakeDispatchEntry from '~/components/centrum/livemap/TakeDispatchEntry.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import { useCentrumStore } from '~/stores/centrum';
import { getCentrumCentrumClient } from '~~/gen/ts/clients';
import { type Dispatch, StatusDispatch, TakeDispatchResp } from '~~/gen/ts/resources/centrum/dispatches';
import { CentrumMode } from '~~/gen/ts/resources/centrum/settings';

const { isOpen } = useSlideover();

const centrumStore = useCentrumStore();
const { dispatches, pendingDispatches, getCurrentMode } = storeToRefs(centrumStore);

const centrumCentrumClient = await getCentrumCentrumClient();

const selectedDispatches = ref<number[]>([]);
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
        const call = centrumCentrumClient.takeDispatch({
            dispatchIds,
            resp,
        });
        await call;

        selectedDispatches.value.length = 0;

        isOpen.value = false;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function selectDispatch(id: number, state: boolean): void {
    const idx = selectedDispatches.value.findIndex((did) => did === id);
    if (idx > -1 && !state) {
        selectedDispatches.value.splice(idx, 1);
    } else if (idx === -1 && state) {
        selectedDispatches.value.push(id);
    }
}

const canTakeDispatch = computed(
    () =>
        selectedDispatches.value.length > 0 &&
        (pendingDispatches.value.length > 0 || (getCurrentMode.value === CentrumMode.SIMPLIFIED && dispatches.value.size > 0)),
);

const filteredDispatches = computedAsync(async () => {
    const filtered: Dispatch[] = [];
    dispatches.value.forEach((d) => {
        if (d.id.toString().includes(queryDispatches.value) || d.message.includes(queryDispatches.value)) {
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
    <USlideover :ui="{ width: 'w-screen max-w-xl' }">
        <UCard
            class="flex flex-1 flex-col"
            :ui="{
                body: {
                    base: 'flex-1 min-h-[calc(100dvh-(2*var(--header-height)))] max-h-[calc(100dvh-(2*var(--header-height)))] overflow-y-auto',
                    padding: 'px-1 py-2 sm:p-2',
                },
                ring: '',
                divide: 'divide-y divide-gray-100 dark:divide-gray-800',
            }"
        >
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('components.centrum.take_dispatch.title') }}
                    </h3>

                    <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                </div>
            </template>

            <div>
                <dl class="divide-neutral/10 divide-y">
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
                                <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                    <UFormGroup name="search">
                                        <UInput
                                            v-model="queryDispatches"
                                            type="text"
                                            name="search"
                                            :placeholder="$t('common.search')"
                                            leading-icon="i-mdi-search"
                                        />
                                    </UFormGroup>
                                </dd>
                            </div>

                            <TakeDispatchEntry
                                v-for="pd in filteredDispatches"
                                :key="pd"
                                :dispatch="dispatches.get(pd)!"
                                :preselected="false"
                                @selected="selectDispatch(pd, $event)"
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
                            />
                        </template>
                    </template>
                </dl>
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton
                        class="flex-1"
                        color="green"
                        :disabled="!canTakeDispatch || !canSubmit"
                        :loading="!canSubmit"
                        @click="onSubmitThrottle(TakeDispatchResp.ACCEPTED)"
                    >
                        {{ $t('common.accept') }}
                    </UButton>

                    <UButton
                        class="flex-1"
                        color="error"
                        :disabled="!canTakeDispatch || !canSubmit"
                        :loading="!canSubmit"
                        @click="onSubmitThrottle(TakeDispatchResp.DECLINED)"
                    >
                        {{ $t('common.decline') }}
                    </UButton>

                    <UButton class="flex-1" @click="isOpen = false">
                        {{ $t('common.close') }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </USlideover>
</template>
