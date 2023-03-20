<script setup lang="ts">
import { useRoute } from 'vue-router/auto';
import { ref, Component, onBeforeMount } from 'vue';
import { CreateDocumentRequest, GetDocumentRequest, RemoveDcoumentReferenceRequest } from '@arpanet/gen/services/docstore/docstore_pb';
import { Document, DocumentAccess, DocumentReference, DocumentRelation } from '@arpanet/gen/resources/documents/documents_pb';
import { getDocStoreClient, handleGRPCError } from '../../grpc';
import { RpcError } from 'grpc-web';
import { getDateLocaleString, getDate } from '../../utils/time';
import { DOC_ACCESS_Util, DOC_REFERENCE_TYPE_Util, DOC_RELATION_TYPE_Util } from '@arpanet/gen/resources/documents/documents.pb_enums';
import {
    TabGroup,
    TabList,
    Tab,
    TabPanels,
    TabPanel,
} from '@headlessui/vue';
import {
    LockOpenIcon,
    PencilIcon,
    ChatBubbleLeftEllipsisIcon,
    CalendarIcon,
    UserIcon,
    DocumentMagnifyingGlassIcon,
    TagIcon,
    ArrowLongRightIcon,
} from '@heroicons/vue/20/solid';

const route = useRoute();

const document = ref<undefined | Document>(undefined)
const access = ref<undefined | DocumentAccess>(undefined)
const comments = ref<Array<Document>>([])
const activeResponse = ref<undefined | Document>(undefined)
const feedReferences = ref<Array<DocumentReference>>([])
const feedRelations = ref<Array<DocumentRelation>>([])
const tabs = ref<{ name: string, href: string, icon: Component }[]>([
    { name: 'References', href: '#', icon: DocumentMagnifyingGlassIcon },
    { name: 'Relations', href: '#', icon: UserIcon },
]);

const props = defineProps({
    documentID: {
        required: true,
        type: Number,
    },
});

function getDocument(): void {
    const req = new GetDocumentRequest();
    req.setDocumentId(props.documentID);

    getDocStoreClient().
        getDocument(req, null).
        then((resp) => {
            document.value = resp.getDocument();
            access.value = resp.getAccess();
        }).
        catch((err: RpcError) => {
            handleGRPCError(err, route);
        });

    // Document References
    getDocStoreClient().
        getDocumentReferences(req, null).
        then((resp) => {
            feedReferences.value = resp.getReferencesList();
        }).
        catch((err: RpcError) => {
            handleGRPCError(err, route);
        });
    // Document Relations
    getDocStoreClient().
        getDocumentRelations(req, null).
        then((resp) => {
            feedRelations.value = resp.getRelationsList();
        }).
        catch((err: RpcError) => {
            handleGRPCError(err, route);
        });
}

function editDocumentTest() {
    const req = new CreateDocumentRequest();
    req.setDocumentId(document?.value?.getId());
    req.setTitle("SCOTT'S DOKUMENTEN WOCHENDSSPAß");
    req.setContent(document?.value?.getContent());
    req.setClosed(document?.value?.getClosed());
    req.setState(document?.value?.getState());
    req.setPublic(document?.value?.getPublic());

    getDocStoreClient().
        createOrUpdateDocument(req, null).then((resp) => {
            console.log(resp);
        }).
        catch((err: RpcError) => {
            handleGRPCError(err, route);
        });
}

function removeDocRefTest() {
    const req = new RemoveDcoumentReferenceRequest();
    req.setId(1);

    getDocStoreClient().
        removeDcoumentReference(req, null).then((resp) => {
            console.log(typeof resp);
        }).catch((err: RpcError) => {
            handleGRPCError(err, route);
        });
}

onBeforeMount(() => {
    getDocument();
});
</script>

<template>
    <div class="mx-auto w-full max-w-7xl flex-grow lg:flex xl:px-8">
        <!-- Left sidebar & main wrapper -->
        <div class="min-w-0 flex-1 bg-white xl:flex">
            <div
                class="border-b border-gray-200 bg-white xl:w-64 xl:flex-shrink-0 xl:border-b-0 xl:border-r xl:border-gray-200 mt-1">
                <div class="h-full py-6 pl-4 pr-6">
                    <div class="flow-root">
                        <h3 class="text-base font-semibold leading-6 text-gray-900 pb-5">Responses</h3>
                        <ul role="list" class="-mb-8">
                            <li v-for="(comment, commentIdx) in comments" :key="comment.getId()" class="mb-5">
                                <div class="relative pb-2 bg-gray-300 rounded-xl">
                                    <span v-if="commentIdx !== comments.length - 1"
                                        class="absolute top-4 left-4 -ml-px h-full w-0.5 bg-gray-200" aria-hidden="true" />
                                    <div class="relative flex space-x-3 cursor-pointer select-none"
                                        @click="activeResponse = comment">
                                        <div>
                                            <span
                                                class="h-8 w-8 rounded-full flex items-center justify-center ring-4 ring-white bg-gray-300">
                                                <UserIcon class="h-5 w-5 text-white" aria-hidden="true" />
                                            </span>
                                        </div>
                                        <div class="flex min-w-0 flex-1 justify-between space-x-4 pt-1.5">
                                            <div>
                                                <p class="font-medium text-gray-900">{{ comment.getTitle() }}</p>
                                                <p class="text-sm text-gray-500">
                                                    by {{ comment.getCreator()?.getFirstname() }} {{
                                                        comment.getCreator()?.getLastname() }}
                                                </p>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </li>
                        </ul>
                    </div>
                </div>
            </div>

            <div class="bg-white lg:min-w-0 lg:flex-1">
                <div class="h-full py-6 px-4 sm:px-6 lg:px-8">
                    <div>
                        <div>
                            <div class="md:flex md:items-center md:justify-between md:space-x-4 xl:border-b xl:pb-6">
                                <div>
                                    <h1 class="text-2xl font-bold text-gray-900">{{ document?.getTitle() }}</h1>
                                    <p class="mt-2 text-sm text-gray-500">
                                        Created by
                                        {{ ' ' }}
                                        <router-link :to="'/citizens/' + document?.getCreator()?.getUserId()"
                                            class="font-medium text-gray-900">
                                            {{ document?.getCreator()?.getFirstname() }}
                                            {{ document?.getCreator()?.getLastname() }}
                                        </router-link>
                                    </p>
                                </div>
                                <div class="mt-4 flex space-x-3 md:mt-0">
                                    <router-link :to="'/documents/' + document?.getId()" type="button"
                                        class="inline-flex justify-center gap-x-1.5 rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50">
                                        <PencilIcon class="-ml-0.5 h-5 w-5 text-gray-400" aria-hidden="true" />
                                        Edit
                                    </router-link>
                                    <button @click="editDocumentTest()">TEST</button>
                                    <button @click="removeDocRefTest()">DELETE REF</button>
                                </div>
                            </div>
                            <aside class="mt-8 xl:hidden">
                                <h2 class="sr-only">Details</h2>
                                <div class="space-y-5">
                                    <div class="flex items-center space-x-2">
                                        <LockOpenIcon class="h-5 w-5 text-green-500" aria-hidden="true" />
                                        <span class="text-sm font-medium text-green-700">Open Issue</span>
                                    </div>
                                    <div class="flex items-center space-x-2">
                                        <ChatBubbleLeftEllipsisIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
                                        <span class="text-sm font-medium text-gray-900">{{ comments.length }}
                                            replies</span>
                                    </div>
                                    <div class="flex items-center space-x-2">
                                        <CalendarIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
                                        <span class="text-sm font-medium text-gray-900"><time
                                                :datetime="getDate(document?.getCreatedAt())?.toLocaleString('de-DE')">{{
                                                    getDate(document?.getCreatedAt())?.toLocaleString('de-DE') }}</time></span>
                                    </div>
                                </div>
                                <div class="mt-6 space-y-8 border-t border-b border-gray-200 py-6">
                                    <div>
                                        <h2 class="text-sm font-medium text-gray-500">Creator</h2>
                                        <ul role="list" class="mt-3 space-y-3">
                                            <li class="flex justify-start">
                                                <router-link :to="'/citizens/' + document?.getCreator()?.getUserId()"
                                                    class="flex items-center space-x-3">
                                                    <div class="text-sm font-medium text-gray-900">{{
                                                        document?.getCreator()?.getFirstname() + ", " +
                                                        document?.getCreator()?.getLastname() }}</div>
                                                </router-link>
                                            </li>
                                        </ul>
                                    </div>
                                    <div>
                                        <h2 class="text-sm font-medium text-gray-500">Tags</h2>
                                        <ul role="list" class="mt-2 leading-8">
                                            <li class="inline">
                                                <a href="#"
                                                    class="relative inline-flex items-center rounded-full px-2.5 py-1 ring-1 ring-inset ring-gray-300 hover:bg-gray-50">
                                                    <div class="absolute flex flex-shrink-0 items-center justify-center">
                                                        <span class="h-1.5 w-1.5 rounded-full bg-rose-500"
                                                            aria-hidden="true" />
                                                    </div>
                                                    <div class="ml-3 text-xs font-semibold text-gray-900">Bug</div>
                                                </a>
                                                {{ ' ' }}
                                            </li>
                                            <li class="inline">
                                                <a href="#"
                                                    class="relative inline-flex items-center rounded-full px-2.5 py-1 ring-1 ring-inset ring-gray-300 hover:bg-gray-50">
                                                    <div class="absolute flex flex-shrink-0 items-center justify-center">
                                                        <span class="h-1.5 w-1.5 rounded-full bg-indigo-500"
                                                            aria-hidden="true" />
                                                    </div>
                                                    <div class="ml-3 text-xs font-semibold text-gray-900">Accessibility
                                                    </div>
                                                </a>
                                                {{ ' ' }}
                                            </li>
                                        </ul>
                                    </div>
                                </div>
                            </aside>
                            <div class="py-3 xl:pt-6 xl:pb-0">
                                <h2 class="sr-only">Description</h2>
                                <div class="prose max-w-none">
                                    <p v-html="document?.getContent()"></p>
                                </div>
                            </div>
                            <div v-if="activeResponse" class="mt-10">
                                <div class="md:flex md:items-center md:justify-between md:space-x-4 xl:border-b xl:pb-6">
                                    <div>
                                        <h1 class="text-2xl font-bold text-gray-900">{{ activeResponse?.getTitle() }}</h1>
                                        <p class="mt-2 text-sm text-gray-500">
                                            Reply by
                                            {{ ' ' }}
                                            <router-link :to="'/citizens/' + activeResponse?.getCreator()?.getUserId()"
                                                class="font-medium text-gray-900">
                                                {{ activeResponse?.getCreator()?.getFirstname() }}
                                                {{ activeResponse?.getCreator()?.getLastname() }}
                                            </router-link>
                                        </p>
                                    </div>
                                </div>
                                <div class="py-3 xl:pt-6 xl:pb-0">
                                    <h2 class="sr-only">Description</h2>
                                    <div class="prose max-w-none">
                                        <p v-html="activeResponse?.getContent()"></p>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="bg-gray-50 pr-4 sm:pr-6 lg:flex-shrink-0 lg:border-l lg:border-gray-200 lg:pr-8 xl:pr-0">
            <div class="h-full py-6 pl-6 lg:w-80">
                <aside class="hidden xl:block xl:pl-8">
                    <h2 class="sr-only">Details</h2>
                    <div class="space-y-5">
                        <div class="flex items-center space-x-2">
                            <LockOpenIcon class="h-5 w-5 text-green-500" aria-hidden="true" />
                            <span class="text-sm font-medium text-green-700">Open Issue</span>
                        </div>
                        <div class="flex items-center space-x-2">
                            <ChatBubbleLeftEllipsisIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
                            <span class="text-sm font-medium text-gray-900">{{ comments.length }}
                                replies</span>
                        </div>
                        <div class="flex items-center space-x-2">
                            <CalendarIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
                            <span class="text-sm font-medium text-gray-900"><time
                                    :datetime="getDate(document?.getCreatedAt())?.toLocaleString('de-DE')">{{
                                        getDate(document?.getCreatedAt())?.toLocaleString('de-DE') }}</time></span>
                        </div>
                    </div>
                    <div class="mt-6 space-y-8 border-t border-gray-200 py-6">
                        <div>
                            <h2 class="text-sm font-medium text-gray-500">Creator</h2>
                            <ul role="list" class="mt-3 space-y-3">
                                <li class="flex justify-start">
                                    <router-link :to="'/citizens/' + document?.getCreator()?.getUserId()"
                                        class="flex items-center space-x-3">
                                        <div class="text-sm font-medium text-gray-900">
                                            {{ document?.getCreator()?.getFirstname() + ", " +
                                                document?.getCreator()?.getLastname() }}
                                        </div>
                                    </router-link>
                                </li>
                            </ul>
                        </div>
                        <div>
                            <h2 class="text-sm font-medium text-gray-500">Access</h2>
                            <ul role="list" class="mt-2 leading-8">
                                <li v-for="ac in access?.getJobsList()" class="inline">
                                    <a href="#"
                                        class="relative inline-flex items-center rounded-full px-2.5 py-1 ring-1 ring-inset ring-gray-300 hover:bg-gray-50">
                                        <div class="absolute flex flex-shrink-0 items-center justify-center">
                                            <span class="h-1.5 w-1.5 rounded-full bg-rose-500" aria-hidden="true" />
                                        </div>
                                        <div class="ml-3 text-xs font-semibold text-gray-900">
                                            {{ ac.getJob() }}<span v-if="ac.getMinimumgrade() > 0">(Rank: {{
                                                ac.getMinimumgrade() }})</span> - {{
        DOC_ACCESS_Util.toEnumKey(ac.getAccess()) }}
                                        </div>
                                    </a>
                                    {{ ' ' }}
                                </li>
                                <li v-for="ac in access?.getUsersList()" class="inline">
                                    <a href="#"
                                        class="relative inline-flex items-center rounded-full px-2.5 py-1 ring-1 ring-inset ring-gray-300 hover:bg-gray-50">
                                        <div class="absolute flex flex-shrink-0 items-center justify-center">
                                            <span class="h-1.5 w-1.5 rounded-full bg-rose-500" aria-hidden="true" />
                                        </div>
                                        <div class="ml-3 text-xs font-semibold text-gray-900">{{ ac.getUserId() }} - {{
                                            DOC_ACCESS_Util.toEnumKey(ac.getAccess()) }}</div>
                                    </a>
                                    {{ ' ' }}
                                </li>
                            </ul>
                        </div>
                    </div>
                </aside>
            </div>
        </div>
    </div>
    <div class="bg-white">
        <TabGroup>
            <TabList>
                <Tab v-for="tab in tabs" :key="tab.name"
                    class="border-indigo-500 text-indigo-600 border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 group inline-flex items-center border-b-2 py-4 px-1 text-sm font-medium">
                    <component :is="tab.icon"
                        class="text-indigo-500 text-gray-400 group-hover:text-gray-500 -ml-0.5 mr-2 h-5 w-5"
                        aria-hidden="true" /> {{ tab.name }}
                </Tab>
            </TabList>
            <TabPanels>
                <TabPanel>
                    <div class="overflow-hidden bg-white shadow sm:rounded-md">
                        <ul role="list" class="divide-y divide-gray-200">
                            <li v-for="item in feedReferences" :key="item.getId()">
                                <a href="#" class="block hover:bg-gray-50">
                                    <div class="px-4 py-4 sm:px-6">
                                        <div class="flex items-center justify-between">
                                            <p class="truncate text-sm font-medium text-indigo-600">{{
                                                item.getSourceDocument()?.getTitle() }}</p>
                                            <div class="ml-2 flex flex-shrink-0">
                                                <p
                                                    class="inline-flex rounded-full bg-green-100 px-2 text-xs font-semibold leading-5 text-green-800">
                                                    {{ DOC_REFERENCE_TYPE_Util.toEnumKey(item.getReference()) }}</p>
                                            </div>
                                        </div>
                                        <div class="mt-2 sm:flex sm:justify-between">
                                            <div class="sm:flex">
                                                <p class="flex items-center text-sm text-gray-500">
                                                    <TagIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-gray-400"
                                                        aria-hidden="true" />
                                                    {{ item.getSourceDocument()?.getCategory()?.getName() }}
                                                </p>
                                            </div>
                                            <div class="mt-2 flex items-center text-sm text-gray-500 sm:mt-0">
                                                <CalendarIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-gray-400"
                                                    aria-hidden="true" />
                                                <p>
                                                    Created on
                                                    {{ ' ' }}
                                                    <time :datetime="getDateLocaleString(item.getCreatedAt())">{{
                                                        getDateLocaleString(item.getCreatedAt()) }}</time>
                                                </p>
                                            </div>
                                        </div>
                                    </div>
                                </a>
                            </li>
                        </ul>
                    </div>
                </TabPanel>
                <TabPanel>
                    <div class="overflow-hidden bg-white shadow sm:rounded-md">
                        <ul role="list" class="divide-y divide-gray-200">
                            <li v-for="item in feedRelations" :key="item.getId()">
                                <a href="#" class="block hover:bg-gray-50">
                                    <div class="px-4 py-4 sm:px-6">
                                        <div class="flex items-center justify-between">
                                            <p class="truncate text-sm font-medium text-indigo-600">
                                                {{ item.getSourceUser()?.getFirstname() }}, {{
                                                    item.getSourceUser()?.getLastname() }}
                                                <ArrowLongRightIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-gray-400"
                                                    aria-hidden="true" /> {{ item.getSourceUser()?.getLastname() }}
                                            </p>
                                            <div class="ml-2 flex flex-shrink-0">
                                                <p
                                                    class="inline-flex rounded-full bg-green-100 px-2 text-xs font-semibold leading-5 text-green-800">
                                                    {{ DOC_RELATION_TYPE_Util.toEnumKey(item.getRelation()) }}</p>
                                            </div>
                                        </div>
                                        <div class="mt-2 sm:flex sm:justify-between">
                                            <div class="mt-2 flex items-center text-sm text-gray-500 sm:mt-0">
                                                <CalendarIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-gray-400"
                                                    aria-hidden="true" />
                                                <p>
                                                    Created on
                                                    {{ ' ' }}
                                                    <time :datetime="getDateLocaleString(item.getCreatedAt())">{{
                                                        getDateLocaleString(item.getCreatedAt()) }}</time>
                                                </p>
                                            </div>
                                        </div>
                                    </div>
                                </a>
                            </li>
                        </ul>
                    </div>
                </TabPanel>
            </TabPanels>
        </TabGroup>
    </div>
</template>
