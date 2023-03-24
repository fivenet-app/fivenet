<script setup lang="ts">
import { OrderBy, PaginationRequest } from '@arpanet/gen/resources/common/database/database_pb';
import { DOC_RELATION_Util } from '@arpanet/gen/resources/documents/documents.pb_enums';
import { DocumentRelation, DOC_RELATION } from '@arpanet/gen/resources/documents/documents_pb';
import { User } from '@arpanet/gen/resources/users/users_pb';
import { FindUsersRequest } from '@arpanet/gen/services/citizenstore/citizenstore_pb';
import { GetDocumentRequest, RemoveDocumentRelationRequest, AddDocumentRelationRequest } from '@arpanet/gen/services/docstore/docstore_pb';
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
    UserPlusIcon,
    MagnifyingGlassIcon,
    ChatBubbleBottomCenterTextIcon,
    ExclamationTriangleIcon,
    ShieldExclamationIcon,
} from '@heroicons/vue/24/outline';
import { watchDebounced } from '@vueuse/core';
import { onMounted, ref, FunctionalComponent } from 'vue';
import { useRouter } from 'vue-router/auto';
import { getCitizenStoreClient, getDocStoreClient } from '../../grpc/grpc';
import { useStore } from '../../store/store';
import { toTitleCase } from '../../utils/strings';

const store = useStore();
const router = useRouter();

const props = defineProps<{
    open: boolean,
    document: number | undefined,
}>();

const emit = defineEmits<{
    (e: 'close'): void,
}>();

const relations = ref<DocumentRelation[]>([])
const tabs = ref<{ name: string, icon: FunctionalComponent }[]>([
    { name: 'View current', icon: MagnifyingGlassIcon },
    { name: 'Add new', icon: UserPlusIcon },
]);

const entriesUsers = ref<User[]>([]);
const queryChar = ref('');

onMounted(async () => {
    await findRelations();
    await findUsers();
});

watchDebounced(queryChar, async () => await findUsers(), { debounce: 750, maxWait: 2000 });

async function findUsers(): Promise<void> {
    const req = new FindUsersRequest();
    req.setPagination((new PaginationRequest()).setOffset(0));
    req.setSearchname(queryChar.value);

    const resp = await getCitizenStoreClient().findUsers(req, null)
    entriesUsers.value = resp.getUsersList().filter(user => !relations.value.find(r => r.getTargetUserId() === user.getUserId()));
}

async function findRelations(): Promise<void> {
    if (!props.document) return;

    const req = new GetDocumentRequest();
    req.setDocumentId(props.document);

    const resp = await getDocStoreClient().getDocumentRelations(req, null)
    relations.value = resp.getRelationsList();
}

function addRelation(user: User, relation: number): void {
    const rel = new DocumentRelation();
    rel.setDocumentId(props.document!)
    rel.setSourceUserId(store.state.auth!.lastCharID)
    rel.setTargetUserId(user.getUserId())
    rel.setRelation(DOC_RELATION_Util.fromInt(relation));

    const req = new AddDocumentRelationRequest();
    req.setRelation(rel);

    getDocStoreClient().addDocumentRelation(req, null).then(async () => {
        await findRelations();
        await findUsers();
    });
}

function removeRelation(id: number): void {
    const req = new RemoveDocumentRelationRequest();
    req.setId(id);

    getDocStoreClient().removeDocumentRelation(req, null).then(() => {
        findRelations();
    });
}
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-10" @close="emit('close')">
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
                            class="relative transform overflow-hidden rounded-lg bg-white px-4 pt-5 pb-4 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-6xl sm:p-6">
                            <div class="absolute top-0 right-0 hidden pt-4 pr-4 sm:block">
                                <button type="button"
                                    class="rounded-md bg-white text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
                                    @click="emit('close')">
                                    <span class="sr-only">Close</span>
                                    <XMarkIcon class="h-6 w-6" aria-hidden="true" />
                                </button>
                            </div>
                            <DialogTitle as="h3" class="text-base font-semibold leading-6 text-gray-900">Citizen Relations
                            </DialogTitle>
                            <TabGroup>
                                <TabList>
                                    <Tab v-for="tab in tabs" :key="tab.name" v-slot="{ selected }">
                                        <button
                                            :class="[selected ? 'border-indigo-500 text-indigo-600' : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700', 'group inline-flex items-center border-b-2 py-4 px-1 text-sm font-medium']"
                                            :aria-current="selected ? 'page' : undefined">
                                            <component :is="tab.icon"
                                                :class="[selected ? 'text-indigo-500' : 'text-gray-400 group-hover:text-gray-500', '-ml-0.5 mr-2 h-5 w-5']"
                                                aria-hidden="true" />
                                            <span>{{ tab.name }}</span>
                                        </button>
                                    </Tab>
                                </TabList>
                                <TabPanels>
                                    <div class="sm:flex sm:items-start px-4 sm:px-6 lg:px-8">
                                        <TabPanel class="w-full">
                                            <div class="flow-root">
                                                <div class="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-8">
                                                    <div class="inline-block min-w-full py-2 align-middle">
                                                        <table class="min-w-full divide-y divide-gray-300">
                                                            <thead>
                                                                <tr>
                                                                    <th scope="col"
                                                                        class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 sm:pl-6 lg:pl-8">
                                                                        Name</th>
                                                                    <th scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
                                                                        Creator</th>
                                                                    <th scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
                                                                        Relation</th>
                                                                    <th scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
                                                                        Actions</th>
                                                                </tr>
                                                            </thead>
                                                            <tbody class="divide-y divide-gray-200 bg-white">
                                                                <tr v-for="ref in relations" :key="ref.getId()">
                                                                    <td
                                                                        class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-6 lg:pl-8 truncate">
                                                                        {{ ref.getTargetUser()?.getFirstname() }} {{
                                                                            ref.getTargetUser()?.getLastname() }}</td>
                                                                    <td
                                                                        class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                                                                        {{
                                                                            ref.getSourceUser()?.getFirstname() }}
                                                                        {{ ref.getSourceUser()?.getLastname() }}
                                                                    </td>
                                                                    <td
                                                                        class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                                                                        {{
                                                                            toTitleCase(DOC_RELATION_Util.toEnumKey(ref.getRelation())!.toLowerCase())
                                                                            ?? ref.getRelation() }}</td>
                                                                    <td
                                                                        class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                                                                        <div class="flex flex-row gap-2">
                                                                            <div class="flex">
                                                                                <a :href="router.resolve({ name: 'Citizens: Info', params: { id: ref.getTargetUserId() } }).href"
                                                                                    target="_blank">
                                                                                    <ArrowTopRightOnSquareIcon
                                                                                        class="w-6 h-auto text-indigo-700 hover:text-indigo-500">
                                                                                    </ArrowTopRightOnSquareIcon>
                                                                                </a>
                                                                            </div>
                                                                            <div class="flex">
                                                                                <button role="button"
                                                                                    @click="removeRelation(ref.getId())">
                                                                                    <DocumentMinusIcon
                                                                                        class="w-6 h-auto text-red-700 hover:text-red-500">
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
                                                <label for="name" class="sr-only">Name</label>
                                                <input type="name" name="name" id="name"
                                                    class="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                                                    placeholder="Citizen Name" v-model="queryChar" />
                                            </div>
                                            <div class="flow-root">
                                                <div class="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-8">
                                                    <div class="inline-block min-w-full py-2 align-middle">
                                                        <table class="min-w-full divide-y divide-gray-300">
                                                            <thead>
                                                                <tr>
                                                                    <th scope="col"
                                                                        class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 sm:pl-6 lg:pl-8">
                                                                        Name</th>
                                                                    <th scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
                                                                        Job</th>
                                                                    <th scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
                                                                        Sex</th>
                                                                    <th scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
                                                                        Add Relation</th>
                                                                </tr>
                                                            </thead>
                                                            <tbody class="divide-y divide-gray-200 bg-white">
                                                                <tr v-for="user in entriesUsers" :key="user.getUserId()">
                                                                    <td
                                                                        class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-6 lg:pl-8 truncate">
                                                                        {{ user.getFirstname() }} {{
                                                                            user.getLastname() }}</td>
                                                                    <td
                                                                        class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                                                                        {{ user.getJobLabel() }}
                                                                    </td>
                                                                    <td
                                                                        class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                                                                        {{ user.getSex() }}
                                                                    </td>
                                                                    <td
                                                                        class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                                                                        <div class="flex flex-row gap-2">
                                                                            <div class="flex">
                                                                                <button role="button"
                                                                                    @click="addRelation(user, 0)"
                                                                                    data-te-toggle="tooltip"
                                                                                    title="Mentioned">
                                                                                    <ChatBubbleBottomCenterTextIcon
                                                                                        class="w-6 h-auto text-green-700 hover:text-green-500">
                                                                                    </ChatBubbleBottomCenterTextIcon>
                                                                                </button>
                                                                            </div>
                                                                            <div class="flex">
                                                                                <button role="button"
                                                                                    @click="addRelation(user, 1)"
                                                                                    data-te-toggle="tooltip"
                                                                                    title="Targets">
                                                                                    <ExclamationTriangleIcon
                                                                                        class="w-6 h-auto text-yellow-700 hover:text-yellow-500">
                                                                                    </ExclamationTriangleIcon>
                                                                                </button>
                                                                            </div>
                                                                            <div class="flex">
                                                                                <button role="button"
                                                                                    @click="addRelation(user, 2)"
                                                                                    data-te-toggle="tooltip" title="Caused">
                                                                                    <ShieldExclamationIcon
                                                                                        class="w-6 h-auto text-red-700 hover:text-red-500">
                                                                                    </ShieldExclamationIcon>
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
                            <div class="mt-5 sm:mt-4 sm:flex sm:flex-row-reverse gap-2">
                                <button type="button"
                                    class="mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:mt-0 sm:w-auto"
                                    @click="emit('close')">Close</button>
                        </div>
                    </DialogPanel>
                </TransitionChild>
            </div>
        </div>
    </Dialog>
</TransitionRoot></template>
