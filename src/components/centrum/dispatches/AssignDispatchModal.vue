<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { CloseIcon } from 'mdi-vue3';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import { useCentrumStore } from '~/store/centrum';
import { Dispatch } from '~~/gen/ts/resources/dispatch/dispatches';
import { UNIT_STATUS, Unit } from '~~/gen/ts/resources/dispatch/units';

const centrumStore = useCentrumStore();
const { units } = storeToRefs(centrumStore);

const props = defineProps<{
    open: boolean;
    dispatch: Dispatch;
}>();

const emits = defineEmits<{
    (e: 'close'): void;
}>();

const { $grpc } = useNuxtApp();

const selectedUnits = ref<bigint[]>(props.dispatch.units.map((du) => du.unitId));

async function assignDispatch(): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const toAdd: bigint[] = [];
            const toRemove: bigint[] = [];
            selectedUnits.value?.forEach((u) => {
                toAdd.push(u);
            });
            props.dispatch.units?.forEach((u) => {
                const idx = selectedUnits.value.findIndex((su) => su === u.unitId);
                if (idx === -1) {
                    toRemove.push(u.unitId);
                }
            });

            const call = $grpc.getCentrumClient().assignDispatch({
                dispatchId: props.dispatch.id,
                toAdd: toAdd,
                toRemove: toRemove,
            });
            await call;

            selectedUnits.value.length = 0;
            emits('close');

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

function selectUnit(item: Unit): void {
    const idx = selectedUnits.value?.findIndex((u) => u === item.id);
    if (idx > -1) {
        delete selectedUnits.value[idx];
    } else {
        selectedUnits.value.push(item.id);
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
                                                <DialogTitle class="inline-flex text-base font-semibold leading-6 text-white">
                                                    {{ $t('components.centrum.assign_dispatch.title') }}:
                                                    <IDCopyBadge class="ml-2" :id="dispatch.id" prefix="DSP" />
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
                                        </div>
                                        <div class="flex flex-1 flex-col justify-between">
                                            <div class="divide-y divide-gray-200 px-4 sm:px-6">
                                                <div class="mt-1">
                                                    <div class="my-2 space-y-24">
                                                        <div class="flex-1 form-control">
                                                            <div class="grid grid-cols-3 gap-4">
                                                                <button
                                                                    v-for="item in units"
                                                                    :key="item.name"
                                                                    type="button"
                                                                    :disabled="item.users.length === 0"
                                                                    class="text-white hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
                                                                    :class="[
                                                                        item.users.length === 0
                                                                            ? 'disabled bg-error-600'
                                                                            : selectedUnits?.findIndex(
                                                                                  (u) => u && u === item.id,
                                                                              ) > -1
                                                                            ? 'bg-success-600'
                                                                            : 'bg-info-600',
                                                                    ]"
                                                                    @click="selectUnit(item)"
                                                                >
                                                                    <span class="mt-1"
                                                                        >{{ item.initials }}: {{ item.name }}</span
                                                                    >
                                                                    <span class="mt-1">
                                                                        {{
                                                                            $t(
                                                                                `enums.centrum.UNIT_STATUS.${
                                                                                    UNIT_STATUS[
                                                                                        item.status?.status ?? (0 as number)
                                                                                    ]
                                                                                }`,
                                                                            )
                                                                        }}
                                                                    </span>
                                                                </button>
                                                            </div>
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                    <div class="flex flex-shrink-0 justify-end px-4 py-4">
                                        <span class="isolate inline-flex rounded-md shadow-sm pr-4 w-full">
                                            <button
                                                type="button"
                                                class="w-full relative inline-flex items-center rounded-l-md bg-primary-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-primary-400"
                                                @click="assignDispatch"
                                            >
                                                {{ $t('common.update') }}
                                            </button>
                                            <button
                                                type="button"
                                                class="w-full relative -ml-px inline-flex items-center rounded-r-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 hover:text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-10"
                                                @click="$emit('close')"
                                            >
                                                {{ $t('common.close', 1) }}
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
