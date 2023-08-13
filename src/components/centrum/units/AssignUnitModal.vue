<script lang="ts" setup>
import {
    Combobox,
    ComboboxButton,
    ComboboxInput,
    ComboboxOption,
    ComboboxOptions,
    Dialog,
    DialogPanel,
    DialogTitle,
    TransitionChild,
    TransitionRoot,
} from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { watchDebounced } from '@vueuse/core';
import { CarEmergencyIcon, CheckIcon } from 'mdi-vue3';
import { Unit } from '~~/gen/ts/resources/dispatch/units';
import { UserShort } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    open: boolean;
    unit: Unit;
}>();

const emits = defineEmits<{
    (e: 'close'): void;
}>();

const { $grpc } = useNuxtApp();

const entriesUsers = ref<UserShort[]>([]);
const selectedUsers = ref<UserShort[]>(props.unit.users.filter((u) => u !== undefined).map((u) => u.user!));
const queryUser = ref('');

async function assignUnit(): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const toAdd: number[] = [];
            const toRemove: number[] = [];
            selectedUsers.value?.forEach((u) => {
                toAdd.push(u.userId);
            });
            props.unit.users?.forEach((u) => {
                const idx = selectedUsers.value.findIndex((su) => su.userId === u.userId);
                if (idx === -1) {
                    toRemove.push(u.userId);
                }
            });

            const call = $grpc.getCentrumClient().assignUnit({
                unitId: props.unit.id,
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

async function findChars(): Promise<void> {
    const call = $grpc.getCompletorClient().completeCitizens({
        search: queryUser.value,
        currentJob: true,
        onDuty: true,
    });
    const { response } = await call;

    entriesUsers.value = response.users;
    entriesUsers.value.push(...selectedUsers.value);
}

function charsGetDisplayValue(chars: UserShort[]): string {
    const cs: string[] = [];
    chars.forEach((c) => cs.push(`${c?.firstname} ${c?.lastname}`));

    return cs.join(', ');
}

watchDebounced(queryUser, async () => await findChars(), {
    debounce: 500,
    maxWait: 1250,
});
onMounted(async () => {
    findChars();
});
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
                                    <DialogTitle as="h3" class="text-base font-semibold leading-6">
                                        Assign Users to Unit
                                    </DialogTitle>
                                    <div class="mt-2">
                                        <div class="my-2 space-y-24">
                                            <div class="flex-1 form-control">
                                                <label for="message" class="block text-sm font-medium leading-6 text-neutral">
                                                    {{ $t('common.unit', 2) }}
                                                </label>
                                                <div class="grid grid-cols-4 gap-4">
                                                    <Combobox as="div" v-model="selectedUsers" multiple nullable>
                                                        <div class="relative">
                                                            <ComboboxButton as="div">
                                                                <ComboboxInput
                                                                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                    @change="queryUser = $event.target.value"
                                                                    :display-value="
                                                                        (chars: any) =>
                                                                            chars ? charsGetDisplayValue(chars) : 'N/A'
                                                                    "
                                                                    :placeholder="$t('common.user', 2)"
                                                                />
                                                            </ComboboxButton>

                                                            <ComboboxOptions
                                                                v-if="entriesUsers.length > 0"
                                                                class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-44 sm:text-sm"
                                                            >
                                                                <ComboboxOption
                                                                    v-for="char in entriesUsers"
                                                                    :key="char?.userId"
                                                                    :value="char"
                                                                    as="char"
                                                                    v-slot="{ active, selected }"
                                                                >
                                                                    <li
                                                                        :class="[
                                                                            'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                                                                            active ? 'bg-primary-500' : '',
                                                                        ]"
                                                                    >
                                                                        <span
                                                                            :class="[
                                                                                'block truncate',
                                                                                selected && 'font-semibold',
                                                                            ]"
                                                                        >
                                                                            {{ char?.firstname }}
                                                                            {{ char?.lastname }}
                                                                        </span>

                                                                        <span
                                                                            v-if="selected"
                                                                            :class="[
                                                                                active ? 'text-neutral' : 'text-primary-500',
                                                                                'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                                                            ]"
                                                                        >
                                                                            <CheckIcon class="w-5 h-5" aria-hidden="true" />
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
                                    @click="assignUnit"
                                >
                                    {{ $t('common.update') }}
                                </button>
                            </div>
                        </DialogPanel>
                    </TransitionChild>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>
