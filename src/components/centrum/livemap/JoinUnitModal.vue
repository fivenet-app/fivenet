<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { useThrottleFn } from '@vueuse/core';
import { CloseIcon, LoadingIcon } from 'mdi-vue3';
import { useCentrumStore } from '~/store/centrum';
import { Unit } from '~~/gen/ts/resources/dispatch/units';

defineProps<{
    open: boolean;
}>();

const emit = defineEmits<{
    (e: 'joined', unit: Unit): void;
    (e: 'left'): void;
    (e: 'close'): void;
}>();

const { $grpc } = useNuxtApp();

const centrumStore = useCentrumStore();
const { ownUnitId, getSortedUnits } = storeToRefs(centrumStore);

async function joinOrLeaveUnit(unitId?: bigint): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCentrumClient().joinUnit({
                unitId: unitId,
            });
            const { response } = await call;

            emit('close');

            if (response.unit) {
                emit('joined', response.unit);
            } else {
                emit('left');
            }

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            throw e;
        }
    });
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (unitID?: bigint) => {
    canSubmit.value = false;
    await joinOrLeaveUnit(unitID).finally(() => setTimeout(() => (canSubmit.value = true), 400));
}, 1000);

const queryUnit = ref('');

const sortedUnits = computed(() => getSortedUnits.value);
const filteredUnits = computed(() =>
    sortedUnits.value
        .filter(
            (u) =>
                u.name.toLowerCase().includes(queryUnit.value.toLowerCase()) ||
                u.initials.toLowerCase().includes(queryUnit.value.toLowerCase()),
        )
        .sort((a, b) => a.name.localeCompare(b.name)),
);
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-30" @close="$emit('close')">
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
                                                <DialogTitle class="text-base font-semibold leading-6 text-neutral">
                                                    {{ $t('common.leave_unit') }}
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
                                        </div>
                                        <div class="flex flex-1 flex-col justify-between">
                                            <div class="divide-y divide-gray-200 px-4 sm:px-6">
                                                <div class="mt-1">
                                                    <dl class="">
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
                                                                        v-model="queryUnit"
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
                                                    </dl>
                                                    <div class="my-2 space-y-24">
                                                        <div class="flex-1 form-control">
                                                            <div class="grid grid-cols-2 gap-4">
                                                                <button
                                                                    v-for="unit in filteredUnits"
                                                                    :key="unit.name"
                                                                    type="button"
                                                                    @click="onSubmitThrottle(unit.id)"
                                                                    :disabled="!canSubmit"
                                                                    class="text-neutral hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-1 text-xs my-0.5"
                                                                    :class="[
                                                                        !canSubmit
                                                                            ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                                            : ownUnitId !== undefined && ownUnitId === unit.id
                                                                            ? 'bg-warn-500 hover:bg-warn-100/10'
                                                                            : 'bg-primary-500 hover:bg-primary-100/10',
                                                                    ]"
                                                                >
                                                                    <span class="mt-0.5 text-base">
                                                                        <span class="font-semibold">{{ unit.initials }}</span
                                                                        >: {{ unit.name }}
                                                                    </span>
                                                                    <span class="mt-1 text-xs">
                                                                        {{ $t('common.member', unit.users.length) }}
                                                                    </span>
                                                                    <span
                                                                        v-if="unit.description && unit.description.length > 0"
                                                                        class="text-xs"
                                                                    >
                                                                        <span class="font-semibold">{{
                                                                            $t('common.description')
                                                                        }}</span
                                                                        >: {{ unit.description }}
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
                                                v-if="ownUnitId !== undefined"
                                                type="button"
                                                @click="onSubmitThrottle()"
                                                :disabled="!canSubmit"
                                                class="w-full relative inline-flex items-center rounded-l-md py-2.5 px-3.5 text-sm font-semibold text-neutral"
                                                :class="[
                                                    !canSubmit
                                                        ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                        : 'bg-error-500 hover:bg-primary-400',
                                                ]"
                                            >
                                                <template v-if="!canSubmit">
                                                    <LoadingIcon class="animate-spin h-5 w-5 mr-2" />
                                                </template>
                                                {{ $t('common.leave') }}
                                            </button>
                                            <button
                                                type="button"
                                                class="w-full relative inline-flex items-center bg-neutral px-3 py-2 text-sm font-semibold text-gray-900 hover:text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50"
                                                :class="ownUnitId !== undefined ? '-ml-px rounded-r-md' : 'rounded-md'"
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
