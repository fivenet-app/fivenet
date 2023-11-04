<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { useDebounceFn } from '@vueuse/core';
import { useSound } from '@vueuse/sound';
import { CarEmergencyIcon, CloseIcon } from 'mdi-vue3';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import { useCentrumStore } from '~/store/centrum';
import { Dispatch, StatusDispatch } from '~~/gen/ts/resources/dispatch/dispatches';
import { CentrumMode } from '~~/gen/ts/resources/dispatch/settings';
import { TakeDispatchResp } from '~~/gen/ts/services/centrum/centrum';
import TakeDispatchEntry from '~/components/centrum/livemap/TakeDispatchEntry.vue';
import { isStatusDispatchCompleted } from '~/components/centrum/helpers';

defineProps<{
    open: boolean;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
    (e: 'goto', loc: Coordinate): void;
}>();

const { $grpc } = useNuxtApp();

const centrumStore = useCentrumStore();
const { dispatches, pendingDispatches, getCurrentMode } = storeToRefs(centrumStore);

const selectedDispatches = ref<bigint[]>([]);
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

        emit('close');
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

function selectDispatch(id: bigint, state: boolean): void {
    const idx = selectedDispatches.value.findIndex((did) => did === id);
    if (idx > -1 && !state) {
        selectedDispatches.value.splice(idx, 1);
    } else if (idx === -1 && state) {
        selectedDispatches.value.push(id);
    }
}

const newDispatchSound = useSound('/sounds/centrum/message-incoming.mp3', {
    volume: 0.15,
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

const filteredDispatches = computed(() => {
    const filtered: Dispatch[] = [];
    dispatches.value.forEach((d) => {
        if (d.id.toString().includes(queryDispatches.value) || d.message.includes(queryDispatches.value)) {
            if (d.status === undefined || d.status.status < StatusDispatch.COMPLETED) filtered.push(d);
        }
    });
    return filtered.sort((a, b) => (a.status?.status ?? 0) - (b.status?.status ?? 0)).map((d) => d.id);
});
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-30" @close="$emit('close')">
            <div class="fixed inset-0" />

            <div class="fixed inset-0 overflow-hidden">
                <div class="absolute inset-0 overflow-hidden">
                    <div class="pointer-events-none fixed inset-y-0 right-0 flex max-w-xl pl-10 sm:pl-16">
                        <TransitionChild
                            as="template"
                            enter="transform transition ease-in-out duration-150 sm:duration-300"
                            enter-from="translate-x-full"
                            enter-to="translate-x-0"
                            leave="transform transition ease-in-out duration-150 sm:duration-300"
                            leave-from="translate-x-0"
                            leave-to="translate-x-full"
                        >
                            <DialogPanel class="pointer-events-auto w-screen max-w-3xl">
                                <form class="flex h-full flex-col divide-y divide-gray-200 bg-gray-900 shadow-xl">
                                    <div class="h-0 flex-1 overflow-y-auto">
                                        <div class="bg-primary-700 px-4 py-6 sm:px-6">
                                            <div class="flex items-center justify-between">
                                                <DialogTitle class="text-base font-semibold leading-6 text-neutral">
                                                    {{ $t('components.centrum.take_dispatch.title') }}
                                                </DialogTitle>
                                                <div class="ml-3 flex h-7 items-center">
                                                    <button
                                                        type="button"
                                                        class="rounded-md bg-gray-100 text-gray-500 hover:text-gray-400 focus:outline-none focus:ring-2 focus:ring-neutral"
                                                        @click="$emit('close')"
                                                    >
                                                        <span class="sr-only">{{ $t('common.close') }}</span>
                                                        <CloseIcon class="h-6 w-6" aria-hidden="true" />
                                                    </button>
                                                </div>
                                            </div>
                                            <div class="mt-1">
                                                <p class="text-sm text-primary-300">
                                                    {{ $t('components.centrum.take_dispatch.sub_title') }}
                                                </p>
                                            </div>
                                        </div>
                                        <div class="flex flex-1 flex-col justify-between">
                                            <div class="divide-y divide-gray-200 px-4 sm:px-6">
                                                <div class="mt-1">
                                                    <dl class="border-b border-neutral/10 divide-y divide-neutral/10">
                                                        <template v-if="getCurrentMode === CentrumMode.SIMPLIFIED">
                                                            <DataNoDataBlock
                                                                v-if="dispatches.size === 0"
                                                                :icon="CarEmergencyIcon"
                                                                :type="$t('common.dispatch', 2)"
                                                            />
                                                            <template v-else>
                                                                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                                    <dt class="text-sm font-medium leading-6 text-neutral">
                                                                        <div class="flex h-6 items-center">
                                                                            {{ $t('common.search') }}
                                                                        </div>
                                                                    </dt>
                                                                    <dd
                                                                        class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                                    >
                                                                        <div class="relative flex items-center">
                                                                            <input
                                                                                v-model="queryDispatches"
                                                                                type="text"
                                                                                name="search"
                                                                                :placeholder="$t('common.search')"
                                                                                class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                                @focusin="focusTablet(true)"
                                                                                @focusout="focusTablet(false)"
                                                                            />
                                                                        </div>
                                                                    </dd>
                                                                </div>

                                                                <TakeDispatchEntry
                                                                    v-for="pd in filteredDispatches"
                                                                    :key="pd.toString()"
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
                                                                :icon="CarEmergencyIcon"
                                                                :type="$t('common.dispatch', 2)"
                                                            />
                                                            <template v-else>
                                                                <TakeDispatchEntry
                                                                    v-for="pd in pendingDispatches"
                                                                    :key="pd.toString()"
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
                                    <div class="flex flex-shrink-0 justify-end px-4 py-4">
                                        <span class="isolate inline-flex rounded-md shadow-sm pr-4 w-full">
                                            <button
                                                :disabled="!canTakeDispatch"
                                                type="button"
                                                class="w-full relative inline-flex items-center rounded-l-md bg-success-500 px-3 py-2 text-sm font-semibold text-neutral hover:text-neutral ring-1 ring-inset ring-success-300 hover:bg-success-100"
                                                :class="!canTakeDispatch ? 'disabled' : ''"
                                                @click="takeDispatches(TakeDispatchResp.ACCEPTED)"
                                            >
                                                {{ $t('common.accept') }}
                                            </button>
                                            <button
                                                :disabled="!canTakeDispatch"
                                                type="button"
                                                class="w-full relative -ml-px inline-flex items-center bg-error-500 px-3 py-2 text-sm font-semibold text-neutral ring-1 ring-inset ring-error-300 hover:bg-error-100"
                                                :class="!canTakeDispatch ? 'disabled' : ''"
                                                @click="takeDispatches(TakeDispatchResp.DECLINED)"
                                            >
                                                {{ $t('common.decline') }}
                                            </button>
                                            <button
                                                type="button"
                                                class="w-full relative -ml-px inline-flex items-center rounded-r-md bg-neutral px-3 py-2 text-sm font-semibold text-gray-900 hover:text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50"
                                                @click="$emit('close')"
                                            >
                                                {{ $t('common.close') }}
                                            </button>
                                        </span>
                                    </div>
                                </form>
                            </DialogPanel>
                        </TransitionChild>
                    </div>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>
