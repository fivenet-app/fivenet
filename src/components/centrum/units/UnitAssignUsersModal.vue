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
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { useThrottleFn, watchDebounced } from '@vueuse/core';
import { CheckIcon, CloseIcon, LoadingIcon } from 'mdi-vue3';
import { useCompletorStore } from '~/store/completor';
import { Unit } from '~~/gen/ts/resources/centrum/units';
import { UserShort } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    open: boolean;
    unit: Unit;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
}>();

const { $grpc } = useNuxtApp();

const completorStore = useCompletorStore();

const entriesCitizens = ref<UserShort[]>([]);
const selectedCitizens = ref<UserShort[]>(props.unit.users.filter((u) => u !== undefined).map((u) => u.user!));
const queryCitizens = ref('');

async function assignUnit(): Promise<void> {
    try {
        const toAdd: number[] = [];
        const toRemove: number[] = [];
        selectedCitizens.value?.forEach((u) => {
            toAdd.push(u.userId);
        });
        props.unit.users?.forEach((u) => {
            const idx = selectedCitizens.value.findIndex((su) => su.userId === u.userId);
            if (idx === -1) {
                toRemove.push(u.userId);
            }
        });

        const call = $grpc.getCentrumClient().assignUnit({
            unitId: props.unit.id,
            toAdd,
            toRemove,
        });
        await call;

        emit('close');
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function findCitizens(): Promise<void> {
    entriesCitizens.value = await completorStore.completeCitizens({
        search: queryCitizens.value,
        currentJob: true,
        onDuty: true,
    });
    entriesCitizens.value.unshift(...selectedCitizens.value);
}

function charsGetDisplayValue(chars: UserShort[]): string {
    const cs: string[] = [];
    chars.forEach((c) => cs.push(`${c?.firstname} ${c?.lastname}`));

    return cs.join(', ');
}

watchDebounced(queryCitizens, async () => await findCitizens(), {
    debounce: 500,
    maxWait: 1250,
});

onMounted(async () => {
    findCitizens();
});

const canSubmit = ref(true);

const onSubmitThrottle = useThrottleFn(async () => {
    canSubmit.value = false;
    await assignUnit().finally(() => setTimeout(() => (canSubmit.value = true), 400));
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
                                                <DialogTitle class="text-base font-semibold leading-6 text-neutral">
                                                    {{ $t('components.centrum.assign_unit.title') }}: {{ unit.name }} ({{
                                                        unit.initials
                                                    }})
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
                                                    <div class="my-2 space-y-24">
                                                        <div class="flex-1 form-control">
                                                            <Combobox v-model="selectedCitizens" as="div" multiple nullable>
                                                                <div class="relative">
                                                                    <ComboboxButton as="div">
                                                                        <ComboboxInput
                                                                            autocomplete="off"
                                                                            class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                            :display-value="
                                                                                (chars: any) =>
                                                                                    chars
                                                                                        ? charsGetDisplayValue(chars)
                                                                                        : $t('common.na')
                                                                            "
                                                                            :placeholder="$t('common.user', 2)"
                                                                            @change="queryCitizens = $event.target.value"
                                                                            @focusin="focusTablet(true)"
                                                                            @focusout="focusTablet(false)"
                                                                        />
                                                                    </ComboboxButton>

                                                                    <ComboboxOptions
                                                                        v-if="entriesCitizens.length > 0"
                                                                        class="absolute z-30 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-44 sm:text-sm"
                                                                    >
                                                                        <ComboboxOption
                                                                            v-for="user in entriesCitizens"
                                                                            v-slot="{ active, selected }"
                                                                            :key="user?.userId"
                                                                            :value="user"
                                                                            as="char"
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
                                                                                    {{ user?.firstname }}
                                                                                    {{ user?.lastname }}
                                                                                </span>

                                                                                <span
                                                                                    v-if="selected"
                                                                                    :class="[
                                                                                        active
                                                                                            ? 'text-neutral'
                                                                                            : 'text-primary-500',
                                                                                        'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                                                                    ]"
                                                                                >
                                                                                    <CheckIcon
                                                                                        class="w-5 h-5"
                                                                                        aria-hidden="true"
                                                                                    />
                                                                                </span>
                                                                            </li>
                                                                        </ComboboxOption>
                                                                    </ComboboxOptions>
                                                                </div>
                                                            </Combobox>
                                                            <div class="mt-4">
                                                                <ul
                                                                    class="text-sm font-medium max-w-md space-y-1 text-gray-100 list-disc list-inside dark:text-gray-300"
                                                                >
                                                                    <li v-for="user in selectedCitizens" :key="user.userId">
                                                                        {{ user?.firstname }}
                                                                        {{ user?.lastname }}
                                                                    </li>
                                                                </ul>
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
