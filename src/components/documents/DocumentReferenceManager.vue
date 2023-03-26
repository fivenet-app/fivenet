<script setup lang="ts">
import { PaginationRequest } from '@arpanet/gen/resources/common/database/database_pb';
import { Document, DocumentReference } from '@arpanet/gen/resources/documents/documents_pb';
import { FindDocumentsRequest } from '@arpanet/gen/services/docstore/docstore_pb';
import {
    Dialog,
    DialogPanel,
    DialogTitle,
    TransitionChild,
    TransitionRoot,
    TabGroup,
    TabList,
    Tab,
    TabPanels,
    TabPanel,
} from '@headlessui/vue';
import {
    XMarkIcon,
    ArrowTopRightOnSquareIcon,
    DocumentMinusIcon,
    DocumentPlusIcon,
    MagnifyingGlassIcon,
} from '@heroicons/vue/24/outline';
import { watchDebounced } from '@vueuse/core';
import { onMounted, ref, FunctionalComponent } from 'vue';
import { useRouter } from 'vue-router/auto';
import { getDocStoreClient } from '../../grpc/grpc';
import { useStore } from '../../store/store';
import { getDateLocaleString, getDateRelativeString } from '../../utils/time';

const store = useStore();
const router = useRouter();

const props = defineProps<{
    open: boolean,
    document?: number,
    modelValue: Map<number, DocumentReference>,
}>();

const emit = defineEmits<{
    (e: 'close'): void,
    (e: 'update:modelValue', payload: Map<number, DocumentReference>): void,
}>();

const tabs = ref<{ name: string, icon: FunctionalComponent }[]>([
    { name: 'View current', icon: MagnifyingGlassIcon },
    { name: 'Add new', icon: DocumentPlusIcon },
]);

const entriesDocuments = ref<Document[]>([]);
const queryDoc = ref('');

onMounted(() => {
    findDocuments();
});

watchDebounced(queryDoc, async () => findDocuments(), { debounce: 750, maxWait: 2000 });

function findDocuments(): void {
    const req = new FindDocumentsRequest();
    req.setPagination((new PaginationRequest()).setOffset(0));
    req.setSearch(queryDoc.value);

    getDocStoreClient().findDocuments(req, null).then((resp) => {
        entriesDocuments.value = resp.getDocumentsList().filter(doc => !(Array.from(props.modelValue.values()).find(r => r.getTargetDocumentId() === doc.getId() || doc.getId() === props.document)));
    });
}

function addReference(doc: Document): void {
    const keys = Array.from(props.modelValue.keys());
    const key = !keys.length ? 1 : keys[keys.length - 1] + 1;

    const ref = new DocumentReference();
    ref.setId(key);
    ref.setCreatorId(store.state.auth!.activeChar!.getUserId());
    ref.setCreator(store.state.auth!.activeChar!)
    ref.setTargetDocumentId(doc.getId());
    ref.setTargetDocument(doc);

    props.modelValue.set(key, ref);
    findDocuments();
}

function removeReference(id: number): void {
    props.modelValue.delete(id);
    findDocuments();
}
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-10" @close="emit('close')">
            <TransitionChild as="template" enter="ease-out duration-300" enter-from="opacity-0" enter-to="opacity-100"
                leave="ease-in duration-200" leave-from="opacity-100" leave-to="opacity-0">
                <div class="fixed inset-0 transition-opacity bg-opacity-75 bg-base-900" />
            </TransitionChild>

            <div class="fixed inset-0 z-10 overflow-y-auto">
                <div class="flex items-end justify-center min-h-full p-4 text-center sm:items-center sm:p-0">
                    <TransitionChild as="template" enter="ease-out duration-300"
                        enter-from="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                        enter-to="opacity-100 translate-y-0 sm:scale-100" leave="ease-in duration-200"
                        leave-from="opacity-100 translate-y-0 sm:scale-100"
                        leave-to="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95">
                        <DialogPanel
                            class="relative px-4 pt-5 pb-4 overflow-hidden text-left transition-all transform rounded-lg bg-base-850 text-neutral sm:my-8 sm:w-full sm:max-w-6xl sm:p-6">
                            <div class="absolute top-0 right-0 hidden pt-4 pr-4 sm:block">
                                <button type="button"
                                    class="transition-colors rounded-md hover:text-base-300 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
                                    @click="emit('close')">
                                    <span class="sr-only">Close</span>
                                    <XMarkIcon class="w-6 h-6" aria-hidden="true" />
                                </button>
                            </div>
                            <DialogTitle as="h3" class="text-base font-semibold leading-6">Document References
                            </DialogTitle>
                            <TabGroup>
                                <TabList class="flex flex-row mb-4">
                                    <Tab v-for="tab in tabs" :key="tab.name" v-slot="{ selected }"
                                        class="flex-initial w-full">
                                        <button
                                            :class="[selected ? 'border-primary-500 text-primary-500' : 'border-transparent text-base-300 hover:border-base-300 hover:text-base-200', 'group inline-flex items-center border-b-2 py-4 px-1 text-m font-medium w-full justify-center transition-colors']"
                                            :aria-current="selected ? 'page' : undefined">
                                            <component :is="tab.icon"
                                                :class="[selected ? 'text-primary-500' : 'text-base-300 group-hover:text-base-200', '-ml-0.5 mr-2 h-5 w-5 transition-colors']"
                                                aria-hidden="true" />
                                            <span>{{ tab.name }}</span>
                                        </button>
                                    </Tab>
                                </TabList>
                                <TabPanels>
                                    <div class="px-4 sm:flex sm:items-start sm:px-6 lg:px-8">
                                        <TabPanel class="w-full">
                                            <div class="flow-root">
                                                <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                                                    <div class="inline-block min-w-full py-2 align-middle">
                                                        <table class="min-w-full divide-y divide-base-200 text-neutral">
                                                            <thead>
                                                                <tr>
                                                                    <th scope="col"
                                                                        class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold  sm:pl-6 lg:pl-8">
                                                                        Title</th>
                                                                    <th scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold">
                                                                        Creator</th>
                                                                    <th scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold">
                                                                        State</th>
                                                                    <th scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold">
                                                                        Actions</th>
                                                                </tr>
                                                            </thead>
                                                            <tbody class="divide-y divide-base-500">
                                                                <tr v-for="[key, ref] in $props.modelValue" :key="key">
                                                                    <td
                                                                        class="py-4 pl-4 pr-3 text-sm font-medium truncate whitespace-nowrap sm:pl-6 lg:pl-8">
                                                                        {{ ref.getTargetDocument()?.getTitle() }}</td>
                                                                    <td class="px-3 py-4 text-sm whitespace-nowrap">
                                                                        {{
                                                                            ref.getCreator()?.getFirstname() }}
                                                                        {{ ref.getCreator()?.getLastname() }}
                                                                    </td>
                                                                    <td class="px-3 py-4 text-sm whitespace-nowrap">
                                                                        {{ ref.getTargetDocument()?.getState() }}</td>
                                                                    <td class="px-3 py-4 text-sm whitespace-nowrap">
                                                                        <div class="flex flex-row gap-2">
                                                                            <div class="flex">
                                                                                <a :href="router.resolve({ name: 'Documents: Info', params: { id: ref.getTargetDocumentId() } }).href"
                                                                                    target="_blank" data-te-toggle="tooltip"
                                                                                    title="Open Document">
                                                                                    <ArrowTopRightOnSquareIcon
                                                                                        class="w-6 h-auto text-primary-500 hover:text-primary-300">
                                                                                    </ArrowTopRightOnSquareIcon>
                                                                                </a>
                                                                            </div>
                                                                            <div class="flex">
                                                                                <button role="button"
                                                                                    @click="removeReference(ref.getId())"
                                                                                    data-te-toggle="tooltip"
                                                                                    title="Remove Reference">
                                                                                    <DocumentMinusIcon
                                                                                        class="w-6 h-auto text-error-400 hover:text-error-200">
                                                                                    </DocumentMinusIcon>
                                                                                </button>
                                                                            </div>
                                                                        </div>
                                                                    </td>
                                                                </tr>
                                                            </tbody>
                                                        </table>
                                                    </div>
                                                </div>
                                            </div>
                                        </TabPanel>
                                        <TabPanel class="w-full">
                                            <div>
                                                <label for="title" class="sr-only">Name</label>
                                                <input type="text" name="title" id="title"
                                                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    placeholder="Document Title" v-model="queryDoc" />
                                            </div>
                                            <div class="flow-root mt-2">
                                                <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                                                    <div class="inline-block min-w-full py-2 align-middle">
                                                        <table class="min-w-full divide-y divide-base-200">
                                                            <thead>
                                                                <tr>
                                                                    <th scope="col"
                                                                        class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold sm:pl-6 lg:pl-8">
                                                                        Title</th>
                                                                    <th scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold">
                                                                        Author</th>
                                                                    <th scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold">
                                                                        State</th>
                                                                    <th scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold">
                                                                        Created at</th>
                                                                    <th scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold">
                                                                        Add Reference</th>
                                                                </tr>
                                                            </thead>
                                                            <tbody class="divide-y divide-base-500">
                                                                <tr v-for="doc in entriesDocuments.slice(0, 8)"
                                                                    :key="doc.getId()">
                                                                    <td
                                                                        class="py-4 pl-4 pr-3 text-sm font-medium truncate whitespace-nowrap sm:pl-6 lg:pl-8">
                                                                        {{ doc.getTitle() }}</td>
                                                                    <td class="px-3 py-4 text-sm whitespace-nowrap">
                                                                        {{ doc.getCreator()?.getFirstname() }} {{
                                                                            doc.getCreator()?.getLastname() }}
                                                                    </td>
                                                                    <td class="px-3 py-4 text-sm whitespace-nowrap">
                                                                        {{ doc.getState() }}
                                                                    </td>
                                                                    <td class="px-3 py-4 text-sm whitespace-nowrap">
                                                                        Created <time
                                                                            :datetime="getDateLocaleString(doc.getCreatedAt())">{{
                                                                                getDateRelativeString(doc.getCreatedAt())
                                                                            }}</time>
                                                                    </td>
                                                                    <td class="px-3 py-4 text-sm whitespace-nowrap">
                                                                        <div class="flex flex-row gap-2">
                                                                            <div class="flex">
                                                                                <button role="button"
                                                                                    @click="addReference(doc)"
                                                                                    data-te-toggle="tooltip"
                                                                                    title="References">
                                                                                    <DocumentPlusIcon
                                                                                        class="w-6 h-auto text-success-500 hover:text-success-300">
                                                                                    </DocumentPlusIcon>
                                                                                </button>
                                                                            </div>
                                                                        </div>
                                                                    </td>
                                                                </tr>
                                                            </tbody>
                                                        </table>
                                                    </div>
                                                </div>
                                            </div>
                                        </TabPanel>
                                    </div>
                                </TabPanels>
                            </TabGroup>
                            <div class="gap-2 mt-5 sm:mt-4 sm:flex sm:flex-row-reverse">
                                <button type="button"
                                    class="rounded-md bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                    @click="emit('close')">Close</button>
                            </div>
                        </DialogPanel>
                    </TransitionChild>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>
