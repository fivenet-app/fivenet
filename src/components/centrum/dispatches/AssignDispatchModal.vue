<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { computedAsync, useThrottleFn } from '@vueuse/core';
import { CancelIcon, CheckIcon, CheckboxBlankOutlineIcon, CloseIcon, LoadingIcon } from 'mdi-vue3';
import { statusOrder, unitStatusToBGColor } from '~/components/centrum/helpers';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import { useCentrumStore } from '~/store/centrum';
import { Dispatch } from '~~/gen/ts/resources/centrum/dispatches';
import { StatusUnit, Unit } from '~~/gen/ts/resources/centrum/units';
import type { GroupedUnits } from '~/components/centrum/helpers';

const centrumStore = useCentrumStore();
const { getSortedUnits } = storeToRefs(centrumStore);

const props = defineProps<{
    open: boolean;
    dispatch: Dispatch;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
}>();

const { $grpc } = useNuxtApp();

const selectedUnits = ref<string[]>(props.dispatch.units.map((du) => du.unitId));

async function assignDispatch(): Promise<void> {
    try {
        const toAdd: string[] = [];
        const toRemove: string[] = [];
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
            toAdd,
            toRemove,
        });
        await call;

        selectedUnits.value.length = 0;
        emit('close');
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

function selectUnit(item: Unit): void {
    const idx = selectedUnits.value?.findIndex((u) => u === item.id);
    if (idx > -1) {
        delete selectedUnits.value[idx];
    } else {
        selectedUnits.value.push(item.id);
    }
}

const grouped = computedAsync(async () => {
    const groups: GroupedUnits = [];
    getSortedUnits.value.forEach((e) => {
        const idx = groups.findIndex((g) => g.key === e.status?.status.toString());
        if (idx === -1) {
            groups.push({
                status: e.status?.status ?? 0,
                units: [e],
                key: e.status?.status.toString() ?? '',
            });
        } else {
            groups[idx].units.push(e);
        }
    });

    groups
        .sort((a, b) => statusOrder.indexOf(a.status) - statusOrder.indexOf(b.status))
        .forEach((group) =>
            group.units.sort((a, b) => {
                if (a.users.length === b.users.length) {
                    return 0;
                } else if (a.users.length === 0) {
                    return 1;
                } else if (b.users.length === 0) {
                    return -1;
                } else {
                    return a.name.localeCompare(b.name);
                }
            }),
        );

    return groups;
});

const canSubmit = ref(true);

const onSubmitThrottle = useThrottleFn(async () => {
    canSubmit.value = false;
    await assignDispatch().finally(() => setTimeout(() => (canSubmit.value = true), 400));
}, 1000);
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
                                                <DialogTitle class="inline-flex text-base font-semibold leading-6 text-neutral">
                                                    {{ $t('components.centrum.assign_dispatch.title') }}:
                                                    <IDCopyBadge :id="dispatch.id" class="ml-2" prefix="DSP" />
                                                </DialogTitle>
                                                <div class="ml-3 flex h-7 items-center">
                                                    <button
                                                        type="button"
                                                        class="rounded-md bg-gray-100 text-gray-500 hover:text-gray-400 focus:outline-none focus:ring-2 focus:ring-neutral"
                                                        @click="$emit('close')"
                                                    >
                                                        <span class="sr-only">{{ $t('common.close') }}</span>
                                                        <CloseIcon class="h-5 w-5" aria-hidden="true" />
                                                    </button>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="flex flex-1 flex-col justify-between">
                                            <div class="divide-y divide-gray-200 px-4">
                                                <div class="mt-1">
                                                    <div class="my-2 space-y-24">
                                                        <div class="flex-1 form-control">
                                                            <template v-for="group in grouped" :key="group.key">
                                                                <p class="text-sm text-neutral">
                                                                    {{
                                                                        $t(
                                                                            `enums.centrum.StatusUnit.${
                                                                                StatusUnit[group.status]
                                                                            }`,
                                                                        )
                                                                    }}
                                                                </p>
                                                                <div class="grid grid-cols-2 lg:grid-cols-3 gap-2">
                                                                    <button
                                                                        v-for="unit in group.units"
                                                                        :key="unit.name"
                                                                        type="button"
                                                                        :disabled="unit.users.length === 0"
                                                                        class="inline-flex flex-row gap-x-1 items-center text-neutral hover:bg-primary-100/10 font-medium hover:transition-all rounded-md p-1.5 text-sm"
                                                                        :class="[
                                                                            unitStatusToBGColor(unit.status?.status),
                                                                            unit.users.length === 0
                                                                                ? 'disabled !bg-error-600'
                                                                                : '',
                                                                        ]"
                                                                        @click="selectUnit(unit)"
                                                                    >
                                                                        <CheckIcon
                                                                            v-if="
                                                                                selectedUnits?.findIndex(
                                                                                    (u) => u && u === unit.id,
                                                                                ) > -1
                                                                            "
                                                                            class="h-5 w-5"
                                                                        />
                                                                        <CheckboxBlankOutlineIcon
                                                                            v-else-if="unit.users.length > 0"
                                                                            class="h-5 w-5"
                                                                        />
                                                                        <CancelIcon v-else class="h-5 w-5" />

                                                                        <div
                                                                            class="ml-0.5 flex flex-col w-full place-items-start"
                                                                        >
                                                                            <span class="font-bold">
                                                                                {{ unit.initials }}
                                                                            </span>
                                                                            <span class="text-xs">
                                                                                {{ unit.name }}
                                                                            </span>
                                                                            <span class="text-xs mt-1">
                                                                                <span class="block">
                                                                                    {{ $t('common.member', unit.users.length) }}
                                                                                </span>
                                                                            </span>
                                                                        </div>
                                                                    </button>
                                                                </div>
                                                            </template>
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
                                                class="w-full relative inline-flex items-center rounded-l-md py-2.5 px-3.5 text-sm font-semibold text-neutral"
                                                :disabled="!canSubmit"
                                                :class="[
                                                    !canSubmit
                                                        ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                        : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                                                ]"
                                                @click="onSubmitThrottle"
                                            >
                                                <template v-if="!canSubmit">
                                                    <LoadingIcon class="animate-spin h-5 w-5 mr-2" />
                                                </template>
                                                {{ $t('common.update') }}
                                            </button>
                                            <button
                                                type="button"
                                                class="w-full relative -ml-px inline-flex items-center rounded-r-md bg-neutral px-3 py-2 text-sm font-semibold text-gray-900 hover:text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50"
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
