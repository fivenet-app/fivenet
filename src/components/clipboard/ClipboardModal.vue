<script lang="ts" setup>
import { useStore } from '../../store/store';
import { computed } from 'vue';
import {
    Dialog,
    DialogPanel,
    DialogTitle,
    TransitionChild,
    TransitionRoot,
} from '@headlessui/vue';
import { ClipboardDocumentListIcon, TrashIcon } from '@heroicons/vue/24/solid';
import { DocumentTextIcon, TruckIcon, UsersIcon } from '@heroicons/vue/20/solid';

const store = useStore();

const users = computed(() => store.state.clipboard?.users);
const documents = computed(() => store.state.clipboard?.documents);
const vehicles = computed(() => store.state.clipboard?.vehicles);

defineProps({
    open: {
        required: true,
        type: Boolean,
    },
});

defineEmits<{
    (e: 'close'): void,
}>();
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-10" @close="$emit('close')">
            <TransitionChild as="template" enter="ease-out duration-300" enter-from="opacity-0" enter-to="opacity-100"
                leave="ease-in duration-200" leave-from="opacity-100" leave-to="opacity-0">
                <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" />
            </TransitionChild>

            <div class="fixed inset-0 z-10 overflow-y-auto">
                <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
                    <TransitionChild as="template" enter="ease-out duration-300"
                        enter-from="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                        enter-to="opacity-100 translate-y-0 sm:scale-100" leave="ease-in duration-200"
                        leave-from="opacity-100 translate-y-0 sm:scale-100"
                        leave-to="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95">
                        <DialogPanel
                            class="relative transform overflow-hidden rounded-lg bg-gray-800 px-4 pt-5 pb-4 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-4xl sm:p-6">
                            <div>
                                <div class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-gray-100">
                                    <ClipboardDocumentListIcon class="h-6 w-6 text-indigo-600" aria-hidden="true" />
                                </div>
                                <div class="mt-3 text-center sm:mt-5">
                                    <DialogTitle as="h3" class="text-base font-semibold leading-6 text-white">
                                        Your Clipboard Contents
                                    </DialogTitle>
                                    <div class="mt-2 text-white">
                                        <h3 class="font-medium pt-1 pb-1">Users</h3>
                                        <button v-if="users?.length == 0" type="button"
                                            class="relative block w-full p-4 text-center border-2 border-dashed rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2"
                                            disabled>
                                            <UsersIcon class="w-12 h-12 mx-auto text-neutral" />
                                            <span class="block mt-2 text-sm font-semibold text-gray-300">No Users in
                                                Clipboard.</span>
                                        </button>
                                        <table v-else class="min-w-full divide-y divide-gray-700">
                                            <thead>
                                                <tr>
                                                    <th scope="col"
                                                        class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-white sm:pl-0">
                                                        Name</th>
                                                    <th scope="col"
                                                        class="py-3.5 px-3 text-left text-sm font-semibold text-white">Job
                                                    </th>
                                                    <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-0">
                                                        Actions
                                                    </th>
                                                </tr>
                                            </thead>
                                            <tbody class="divide-y divide-gray-800">
                                                <tr v-for="user in users" :key="user.id">
                                                    <td
                                                        class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-white sm:pl-0">
                                                        {{ user.firstname }}, {{ user.lastname }}
                                                    </td>
                                                    <td class="whitespace-nowrap py-4 px-3 text-sm text-gray-300">
                                                        {{ user.jobLabel }} (Rank: {{ user.jobGradeLabel }})
                                                    </td>
                                                    <td
                                                        class="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
                                                        <button @click="store.dispatch('clipboard/removeUser', user.id)">
                                                            <TrashIcon class="w-6 h-6 mx-auto text-neutral" />
                                                        </button>
                                                    </td>
                                                </tr>
                                            </tbody>
                                        </table>
                                        <h3 class="font-medium pt-1 pb-1">Documents</h3>
                                        <button v-if="documents?.length == 0" type="button"
                                            class="relative block w-full p-4 text-center border-2 border-dashed rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2"
                                            disabled>
                                            <DocumentTextIcon class="w-12 h-12 mx-auto text-neutral" />
                                            <span class="block mt-2 text-sm font-semibold text-gray-300">No Documents in
                                                Clipboard.</span>
                                        </button>
                                        <table v-else class="min-w-full divide-y divide-gray-700">
                                            <thead>
                                                <tr>
                                                    <th scope="col"
                                                        class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-white sm:pl-0">
                                                        Title</th>
                                                    <th scope="col"
                                                        class="py-3.5 px-3 text-left text-sm font-semibold text-white">
                                                        Creator
                                                    </th>
                                                    <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-0">
                                                        <span class="sr-only">Actions</span>
                                                    </th>
                                                </tr>
                                            </thead>
                                            <tbody class="divide-y divide-gray-800">
                                                <tr v-for="document in documents" :key="document.id">
                                                    <td
                                                        class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-white sm:pl-0">
                                                        {{ document.title }}
                                                    </td>
                                                    <td class="whitespace-nowrap py-4 px-3 text-sm text-gray-300">
                                                        {{ document.creator.firstname }}, {{ document.creator.lastname }}
                                                    </td>
                                                    <td
                                                        class="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
                                                        <button
                                                            @click="store.dispatch('clipboard/removeDocument', document.id)">
                                                            <TrashIcon class="w-6 h-6 mx-auto text-neutral" />
                                                        </button>
                                                    </td>
                                                </tr>
                                            </tbody>
                                        </table>
                                        <h3 class="font-medium pt-2 pb-1">Vehicles</h3>
                                        <button v-if="vehicles?.length == 0" type="button"
                                            class="relative block w-full p-4 text-center border-2 border-dashed rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2"
                                            disabled>
                                            <TruckIcon class="w-12 h-12 mx-auto text-neutral" />
                                            <span class="block mt-2 text-sm font-semibold text-gray-300">No Vehicles in
                                                Clipboard.</span>
                                        </button>
                                        <table v-else class="min-w-full divide-y divide-gray-700">
                                            <thead>
                                                <tr>
                                                    <th scope="col"
                                                        class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-white sm:pl-0">
                                                        Plate</th>
                                                    <th scope="col"
                                                        class="py-3.5 px-3 text-left text-sm font-semibold text-white">Model
                                                    </th>
                                                    <th scope="col"
                                                        class="py-3.5 px-3 text-left text-sm font-semibold text-white">Owner
                                                    </th>
                                                    <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-0">
                                                        <span class="sr-only">Actions</span>
                                                    </th>
                                                </tr>
                                            </thead>
                                            <tbody class="divide-y divide-gray-800">
                                                <tr v-for="vehicle in vehicles" :key="vehicle.plate">
                                                    <td
                                                        class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-white sm:pl-0">
                                                        {{ vehicle.plate }}
                                                    </td>
                                                    <td class="whitespace-nowrap py-4 px-3 text-sm text-gray-300">
                                                        {{ vehicle.model }}
                                                    </td>
                                                    <td class="whitespace-nowrap py-4 px-3 text-sm text-gray-300">
                                                        {{ vehicle.owner.firstname }}, {{ vehicle.owner.lastname }}
                                                    </td>
                                                    <td
                                                        class="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
                                                        <button
                                                            @click="store.dispatch('clipboard/removeVehicle', vehicle.plate)">
                                                            <TrashIcon class="w-6 h-6 mx-auto text-neutral" />
                                                        </button>
                                                    </td>
                                                </tr>
                                            </tbody>
                                        </table>
                                    </div>
                                </div>
                            </div>
                            <div class="mt-5 sm:mt-6 sm:grid sm:grid-flow-row-dense sm:grid-cols-2 sm:gap-3">
                                <button type="button"
                                    class="mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:col-start-1 sm:mt-0"
                                    @click="$emit('close')" ref="cancelButtonRef">Close</button>
                                <button type="button"
                                    class="inline-flex w-full justify-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 sm:col-start-2"
                                    @click="store.dispatch('clipboard/clear')">Clear Clipboard</button>
                            </div>
                        </DialogPanel>
                    </TransitionChild>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>
