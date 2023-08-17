<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { CarEmergencyIcon, CloseIcon } from 'mdi-vue3';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import { Dispatch } from '~~/gen/ts/resources/dispatch/dispatches';
import { Unit } from '~~/gen/ts/resources/dispatch/units';
import { TAKE_DISPATCH_RESP } from '~~/gen/ts/services/centrum/centrum';
import TakeDispatchEntry from './TakeDispatchEntry.vue';

const props = defineProps<{
    open: boolean;
    ownUnit: Unit;
    dispatches: Dispatch[];
}>();

const emits = defineEmits<{
    (e: 'close'): void;
    (e: 'goto', location: { x: number; y: number }): void;
}>();

const { $grpc } = useNuxtApp();

const unselectedDispatches = ref<bigint[]>([]);

async function takeDispatches(resp: TAKE_DISPATCH_RESP): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            if (props.dispatches.length === 0) return;

            const call = $grpc.getCentrumClient().takeDispatch({
                dispatchIds: props.dispatches.filter((d) => !unselectedDispatches.value.includes(d.id)).map((d) => d.id),
                resp: resp,
            });
            await call;

            emits('close');

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

function selectDispatch(id: bigint): void {
    const idx = unselectedDispatches.value.findIndex((did) => did === id);
    if (idx > -1) {
        unselectedDispatches.value.splice(idx, 1);
    } else {
        unselectedDispatches.value.push(id);
    }
}
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-10" @close="$emit('close')">
            <div class="fixed inset-0" />

            <div class="fixed inset-0 overflow-hidden">
                <div class="absolute inset-0 overflow-hidden">
                    <div class="pointer-events-none fixed inset-y-0 right-0 flex max-w-2xl pl-10 sm:pl-16">
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
                                                <DialogTitle class="text-base font-semibold leading-6 text-white">
                                                    Take Dispatch
                                                </DialogTitle>
                                                <div class="ml-3 flex h-7 items-center">
                                                    <button
                                                        type="button"
                                                        class="rounded-md bg-gray-100 text-gray-500 hover:text-gray-400 focus:outline-none focus:ring-2 focus:ring-white"
                                                        @click="$emit('close')"
                                                    >
                                                        <span class="sr-only">{{ $t('common.close') }}</span>
                                                        <CloseIcon class="h-6 w-6" aria-hidden="true" />
                                                    </button>
                                                </div>
                                            </div>
                                            <div class="mt-1">
                                                <p class="text-sm text-primary-300">TODO</p>
                                            </div>
                                        </div>
                                        <div class="flex flex-1 flex-col justify-between">
                                            <div class="divide-y divide-gray-200 px-4 sm:px-6">
                                                <div class="mt-1">
                                                    <dl class="border-b border-white/10 divide-y divide-white/10">
                                                        <DataNoDataBlock
                                                            v-if="dispatches.length === 0"
                                                            :icon="CarEmergencyIcon"
                                                            :type="$t('common.dispatch', 2)"
                                                        />
                                                        <TakeDispatchEntry
                                                            v-else
                                                            v-for="dispatch in dispatches"
                                                            :dispatch="dispatch"
                                                            @selected="selectDispatch(dispatch.id)"
                                                            @goto="$emit('goto', $event)"
                                                        />
                                                    </dl>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                    <div class="flex flex-shrink-0 justify-end px-4 py-4">
                                        <span class="isolate inline-flex rounded-md shadow-sm pr-4 w-full">
                                            <button
                                                :disabled="dispatches.length === 0"
                                                type="button"
                                                class="w-full relative inline-flex items-center rounded-l-md bg-success-500 px-3 py-2 text-sm font-semibold text-white hover:text-white ring-1 ring-inset ring-success-300 hover:bg-success-100 focus:z-10"
                                                :class="dispatches.length === 0 ? 'disabled' : ''"
                                                @click="takeDispatches(TAKE_DISPATCH_RESP.ACCEPTED)"
                                            >
                                                {{ $t('common.accept') }}
                                            </button>
                                            <button
                                                :disabled="dispatches.length === 0"
                                                type="button"
                                                class="w-full relative -ml-px inline-flex items-center bg-error-500 px-3 py-2 text-sm font-semibold text-white ring-1 ring-inset ring-error-300 hover:bg-error-100 focus:z-10"
                                                :class="dispatches.length === 0 ? 'disabled' : ''"
                                                @click="takeDispatches(TAKE_DISPATCH_RESP.DECLINED)"
                                            >
                                                {{ $t('common.decline') }}
                                            </button>
                                            <button
                                                type="button"
                                                class="w-full relative -ml-px inline-flex items-center rounded-r-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 hover:text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-10"
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
