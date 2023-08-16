<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { CarEmergencyIcon } from 'mdi-vue3';
import { Unit } from '~~/gen/ts/resources/dispatch/units';
import { Dispatch } from '../../../../gen/ts/resources/dispatch/dispatches';
import { TAKE_DISPATCH_RESP } from '../../../../gen/ts/services/centrum/centrum';

defineProps<{
    open: boolean;
    ownUnit: Unit;
    dispatch: Dispatch;
}>();

const emits = defineEmits<{
    (e: 'close'): void;
}>();

const { $grpc } = useNuxtApp();

async function takeDispatch(dispatch: Dispatch, resp: TAKE_DISPATCH_RESP): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCentrumClient().takeDispatch({
                dispatchId: dispatch.id,
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
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-10" @close="$emit('close')">
            <TransitionChild
                as="template"
                enter="ease-out duration-300"
                enter-from="opacity-0"
                enter-to="opacity-100"
                leave="ease-in duration-200"
                leave-from="opacity-100"
                leave-to="opacity-0"
            >
                <div class="fixed inset-0 transition-opacity bg-opacity-75 bg-base-900" />
            </TransitionChild>

            <div class="fixed inset-0 z-10 overflow-y-auto">
                <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
                    <TransitionChild
                        as="template"
                        enter="ease-out duration-300"
                        enter-from="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                        enter-to="opacity-100 translate-y-0 sm:scale-100"
                        leave="ease-in duration-200"
                        leave-from="opacity-100 translate-y-0 sm:scale-100"
                        leave-to="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                    >
                        <DialogPanel
                            class="relative px-4 pt-5 pb-4 overflow-hidden text-left transition-all transform rounded-lg bg-base-850 text-neutral sm:my-8 sm:w-full sm:max-w-6xl sm:p-6"
                        >
                            <div>
                                <div class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-base-800">
                                    <CarEmergencyIcon class="h-6 w-6 text-primary-500" aria-hidden="true" />
                                </div>
                                <div class="mt-3 text-center sm:mt-5">
                                    <div v-if="ownUnit">
                                        <DialogTitle as="h3" class="text-base font-semibold leading-6">
                                            Take Dispatch
                                        </DialogTitle>
                                    </div>
                                    <div v-else>
                                        <DialogTitle as="h3" class="text-base font-semibold leading-6"> Join Unit </DialogTitle>
                                        <div class="mt-2">
                                            <div class="my-2 space-y-24">
                                                <div class="flex-1 form-control">
                                                    <label
                                                        for="message"
                                                        class="block text-sm font-medium leading-6 text-neutral"
                                                    >
                                                        {{ $t('common.unit', 2) }}
                                                    </label>
                                                    <div class="grid grid-cols-4 gap-4">TODO</div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <div class="gap-2 mt-5 sm:mt-4 sm:flex">
                                <button
                                    type="button"
                                    class="flex-1 rounded-md bg-error-600 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                    @click="takeDispatch(dispatch, TAKE_DISPATCH_RESP.ACCEPTED)"
                                >
                                    {{ $t('common.accept') }}
                                </button>
                                <button
                                    type="button"
                                    class="flex-1 rounded-md bg-error-600 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                    @click="takeDispatch(dispatch, TAKE_DISPATCH_RESP.DECLINED)"
                                >
                                    {{ $t('common.decline') }}
                                </button>
                                <button
                                    type="button"
                                    class="flex-1 rounded-md bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                    @click="takeDispatch(dispatch, TAKE_DISPATCH_RESP.DECLINED)"
                                    ref="cancelButtonRef"
                                >
                                    {{ $t('common.cancel') }}
                                </button>
                            </div>
                        </DialogPanel>
                    </TransitionChild>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>
