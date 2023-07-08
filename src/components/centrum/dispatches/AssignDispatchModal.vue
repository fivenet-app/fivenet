<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiCarEmergency, mdiCheck } from '@mdi/js';
import { Dispatch } from '~~/gen/ts/resources/dispatch/dispatch';
import { UNIT_STATUS, Unit } from '~~/gen/ts/resources/dispatch/units';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';

const props = defineProps<{
    open: boolean;
    dispatch: Dispatch;
    units: Unit[] | null;
}>();

const emits = defineEmits<{
    (e: 'close'): void;
}>();

const { $grpc } = useNuxtApp();

const queryUnit = ref('');
const query = computed(() => queryUnit.value.toLowerCase());
const filteredUnits = computed(() =>
    (props.units ?? []).filter(
        (v) => v.name.toLowerCase().includes(query.value) || v.initials.toLowerCase().includes(query.value),
    ),
);
const selectedUnits = ref<undefined | Unit[]>(props.dispatch.units.map((du) => du.unit!));

async function assignDispatch(): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const toAdd: bigint[] = [];
            const toRemove: bigint[] = [];
            selectedUnits.value?.forEach((u) => {
                const idx = props.dispatch.units.findIndex((s) => s.unitId === u.id);
                if (idx > -1) {
                    toRemove.push(u.id);
                } else {
                    toAdd.push(u.id);
                }
            });

            const call = $grpc.getCentrumClient().assignDispatch({
                dispatchId: props.dispatch.id,
                toAdd: toAdd,
                toRemove: toRemove,
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
                                    <SvgIcon
                                        class="h-6 w-6 text-primary-500"
                                        aria-hidden="true"
                                        type="mdi"
                                        :path="mdiCarEmergency"
                                    />
                                </div>
                                <div class="mt-3 text-center sm:mt-5">
                                    <DialogTitle as="h3" class="text-base font-semibold leading-6">
                                        Assign Dispatch
                                    </DialogTitle>
                                    <div class="mt-2">
                                        <div class="my-2 space-y-24">
                                            <div class="flex-1 form-control">
                                                <label for="message" class="block text-sm font-medium leading-6 text-neutral">
                                                    {{ $t('common.unit', 2) }}
                                                </label>
                                                <Combobox as="div" v-model="selectedUnits" multiple nullable>
                                                    <div class="relative">
                                                        <ComboboxButton as="div">
                                                            <ComboboxInput
                                                                class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                @change="queryUnit = $event.target.value"
                                                                :display-value="
                                                                    (units: any) =>
                                                                        units
                                                                            ? units.map((u: Unit) => u.initials).join(', ')
                                                                            : ''
                                                                "
                                                                :placeholder="$t('common.unit', 2)"
                                                            />
                                                        </ComboboxButton>

                                                        <ComboboxOptions
                                                            v-if="filteredUnits.length > 0"
                                                            class="absolute z-50 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-44 sm:text-sm"
                                                        >
                                                            <ComboboxOption
                                                                v-for="unit in filteredUnits"
                                                                :key="unit?.id.toString()"
                                                                :value="unit"
                                                                v-slot="{ active, selected }"
                                                            >
                                                                <li
                                                                    :class="[
                                                                        'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                                                                        active ? 'bg-primary-500' : '',
                                                                    ]"
                                                                >
                                                                    <span
                                                                        :class="['block truncate', selected && 'font-semibold']"
                                                                    >
                                                                        {{ unit?.initials }} - {{ unit?.name }} ({{
                                                                            UNIT_STATUS[unit?.status?.status ?? 0]
                                                                        }})
                                                                    </span>

                                                                    <span
                                                                        v-if="selected"
                                                                        :class="[
                                                                            active ? 'text-neutral' : 'text-primary-500',
                                                                            'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                                                        ]"
                                                                    >
                                                                        <SvgIcon
                                                                            class="w-5 h-5"
                                                                            aria-hidden="true"
                                                                            type="mdi"
                                                                            :path="mdiCheck"
                                                                        />
                                                                    </span>
                                                                </li>
                                                            </ComboboxOption>
                                                        </ComboboxOptions>
                                                    </div>
                                                </Combobox>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <div class="gap-2 mt-5 sm:mt-4 sm:flex">
                                <button
                                    type="button"
                                    class="flex-1 rounded-md bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                    @click="$emit('close')"
                                    ref="cancelButtonRef"
                                >
                                    {{ $t('common.close', 1) }}
                                </button>
                                <button
                                    type="button"
                                    class="flex-1 rounded-md bg-primary-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-primary-400"
                                    @click="assignDispatch"
                                >
                                    {{ $t('common.create') }}
                                </button>
                            </div>
                        </DialogPanel>
                    </TransitionChild>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>
