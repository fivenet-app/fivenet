<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, Switch, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { CloseIcon, CogIcon } from 'mdi-vue3';
import { useSettingsStore } from '~/store/settings';

defineProps<{
    open: boolean;
}>();

defineEmits<{
    (e: 'close'): void;
}>();

const settingsStore = useSettingsStore();
const { livemap } = storeToRefs(settingsStore);
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-30" @close="$emit('close')">
            <TransitionChild
                as="template"
                enter="ease-out duration-300"
                enter-from="opacity-0"
                enter-to="opacity-100"
                leave="ease-in duration-200"
                leave-from="opacity-100"
                leave-to="opacity-0"
            >
                <div class="fixed inset-0 bg-gray-500/75 transition-opacity" />
            </TransitionChild>

            <div class="fixed inset-0 z-30 overflow-y-auto">
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
                            class="relative w-full overflow-hidden rounded-lg bg-base-800 px-4 pb-4 pt-5 text-left shadow-xl transition-all sm:my-8 sm:max-w-lg sm:p-6"
                        >
                            <div class="absolute right-0 top-0 block pr-4 pt-4">
                                <button
                                    type="button"
                                    class="rounded-md bg-neutral text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
                                    @click="$emit('close')"
                                >
                                    <span class="sr-only">{{ $t('common.close') }}</span>
                                    <CloseIcon class="size-5" aria-hidden="true" />
                                </button>
                            </div>
                            <div>
                                <div class="mx-auto flex size-12 items-center justify-center rounded-full bg-success-100">
                                    <CogIcon class="size-5 text-success-600" aria-hidden="true" />
                                </div>
                                <div class="mt-3 text-center sm:mt-5">
                                    <DialogTitle as="h3" class="text-base font-semibold leading-6 text-neutral">
                                        {{ $t('common.setting', 2) }}
                                    </DialogTitle>
                                    <div class="mt-2">
                                        <div class="text-sm text-gray-100">
                                            <div class="flex-1">
                                                <label
                                                    for="centerSelectedMarker"
                                                    class="block text-sm font-medium leading-6 text-neutral"
                                                >
                                                    {{ $t('components.livemap.center_selected_marker') }}
                                                </label>
                                                <Switch
                                                    v-model="livemap.centerSelectedMarker"
                                                    :class="[
                                                        livemap.centerSelectedMarker ? 'bg-primary-600' : 'bg-gray-200',
                                                        'relative my-2 inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-primary-600 focus:ring-offset-2',
                                                    ]"
                                                >
                                                    <span
                                                        aria-hidden="true"
                                                        :class="[
                                                            livemap.centerSelectedMarker ? 'translate-x-5' : 'translate-x-0',
                                                            'pointer-events-none inline-block size-5 rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                                                        ]"
                                                    />
                                                </Switch>
                                            </div>
                                            <div class="flex-1">
                                                <label
                                                    for="livemapMarkerSize"
                                                    class="block text-sm font-medium leading-6 text-neutral"
                                                >
                                                    {{ $t('components.livemap.settings.marker_size') }}
                                                </label>
                                                <input
                                                    name="livemapMarkerSize"
                                                    type="range"
                                                    class="my-auto h-1.5 w-full cursor-grab rounded-full"
                                                    min="14"
                                                    max="30"
                                                    step="2"
                                                    :value="livemap.markerSize"
                                                    @change="
                                                        livemap.markerSize = parseInt(($event.target as HTMLInputElement).value)
                                                    "
                                                />
                                                <span class="text-sm text-gray-300">{{ livemap.markerSize }}</span>
                                            </div>
                                            <div class="flex-1 items-center">
                                                <label
                                                    for="showUnitNames"
                                                    class="block text-sm font-medium leading-6 text-neutral"
                                                >
                                                    {{ $t('components.livemap.show_unit_names') }}
                                                </label>
                                                <Switch
                                                    v-model="livemap.showUnitNames"
                                                    :class="[
                                                        livemap.showUnitNames ? 'bg-primary-600' : 'bg-gray-200',
                                                        'relative my-2 inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-primary-600 focus:ring-offset-2',
                                                    ]"
                                                >
                                                    <span
                                                        aria-hidden="true"
                                                        :class="[
                                                            livemap.showUnitNames ? 'translate-x-5' : 'translate-x-0',
                                                            'pointer-events-none inline-block size-5 rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                                                        ]"
                                                    />
                                                </Switch>
                                            </div>
                                            <div class="flex-1 items-center">
                                                <label
                                                    for="showUnitStatus"
                                                    class="block text-sm font-medium leading-6 text-neutral"
                                                >
                                                    {{ $t('components.livemap.show_unit_status') }}
                                                </label>
                                                <Switch
                                                    v-model="livemap.showUnitStatus"
                                                    :class="[
                                                        livemap.showUnitStatus ? 'bg-primary-600' : 'bg-gray-200',
                                                        'relative my-2 inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-primary-600 focus:ring-offset-2',
                                                    ]"
                                                >
                                                    <span
                                                        aria-hidden="true"
                                                        :class="[
                                                            livemap.showUnitStatus ? 'translate-x-5' : 'translate-x-0',
                                                            'pointer-events-none inline-block size-5 rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                                                        ]"
                                                    />
                                                </Switch>
                                            </div>
                                            <div class="flex-1 items-center">
                                                <label
                                                    for="showAllDispatches"
                                                    class="block text-sm font-medium leading-6 text-neutral"
                                                >
                                                    {{ $t('components.livemap.show_all_dispatches') }}
                                                </label>
                                                <Switch
                                                    v-model="livemap.showAllDispatches"
                                                    :class="[
                                                        livemap.showAllDispatches ? 'bg-primary-600' : 'bg-gray-200',
                                                        'relative my-2 inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-primary-600 focus:ring-offset-2',
                                                    ]"
                                                >
                                                    <span
                                                        aria-hidden="true"
                                                        :class="[
                                                            livemap.showAllDispatches ? 'translate-x-5' : 'translate-x-0',
                                                            'pointer-events-none inline-block size-5 rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                                                        ]"
                                                    />
                                                </Switch>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <div class="mt-5 sm:mt-6">
                                <button
                                    type="button"
                                    class="mt-3 inline-flex w-full justify-center rounded-md bg-neutral px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-200 sm:col-start-1 sm:mt-0"
                                    @click="$emit('close')"
                                >
                                    {{ $t('common.close') }}
                                </button>
                            </div>
                        </DialogPanel>
                    </TransitionChild>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>
