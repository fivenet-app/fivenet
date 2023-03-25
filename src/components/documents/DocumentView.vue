<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { GetDocumentCommentsRequest, GetDocumentRequest } from '@arpanet/gen/services/docstore/docstore_pb';
import { Document, DocumentAccess, DocumentComment, DocumentReference, DocumentRelation } from '@arpanet/gen/resources/documents/documents_pb';
import { getDocStoreClient } from '../../grpc/grpc';
import { getDate } from '../../utils/time';
import { DOC_ACCESS_Util } from '@arpanet/gen/resources/documents/documents.pb_enums';
import {
    TabGroup,
    TabList,
    Tab,
    TabPanels,
    TabPanel,
} from '@headlessui/vue';
import {
    LockOpenIcon,
    LockClosedIcon,
    PencilIcon,
    ChatBubbleLeftEllipsisIcon,
    CalendarIcon,
    UserIcon,
    DocumentMagnifyingGlassIcon,
} from '@heroicons/vue/20/solid';
import DocumentRelations from './DocumentRelations.vue';
import DocumentReferences from './DocumentReferences.vue';
import { toTitleCase } from '../../utils/strings';

const document = ref<undefined | Document>(undefined)
const access = ref<undefined | DocumentAccess>(undefined)
const comments = ref<DocumentComment[]>([])
const activeResponse = ref<undefined | Document>(undefined)
const feedReferences = ref<DocumentReference[]>([])
const feedRelations = ref<DocumentRelation[]>([])
const tabs = ref<{ name: string, icon: typeof LockOpenIcon }[]>([
    { name: 'References', icon: DocumentMagnifyingGlassIcon },
    { name: 'Relations', icon: UserIcon },
]);

const props = defineProps({
    documentId: {
        required: true,
        type: Number,
    },
});

function getDocument(): void {
    const req = new GetDocumentRequest();
    req.setDocumentId(props.documentId);

    getDocStoreClient().
        getDocument(req, null).
        then((resp) => {
            document.value = resp.getDocument();
            access.value = resp.getAccess();
        });

    // Document References
    getDocStoreClient().
        getDocumentReferences(req, null).
        then((resp) => {
            feedReferences.value = resp.getReferencesList();
        });

    // Document Relations
    getDocStoreClient().
        getDocumentRelations(req, null).
        then((resp) => {
            feedRelations.value = resp.getRelationsList();
        });

    // Document Comments
    const creq = new GetDocumentCommentsRequest();
    creq.setDocumentId(props.documentId);

    getDocStoreClient().
        getDocumentComments(creq, null).
        then((resp) => {
            comments.value = resp.getCommentsList();
    });
}
onMounted(() => {
    getDocument();
});
</script>

<template>
    <div>
        <div class="rounded-lg bg-base-850">
            <div class="h-full px-4 py-6 sm:px-6 lg:px-8">
                <div>
                    <div>
                        <div class="pb-2 md:flex md:items-center md:justify-between md:space-x-4">
                            <div>
                                <h1 class="text-2xl font-bold text-neutral">{{ document?.getTitle() }}</h1>
                                <p class="text-sm text-base-300">
                                    Created by
                                    {{ ' ' }}
                                    <router-link
                                        :to="{ name: 'Citizens: Info', params: { id: document?.getCreator()?.getUserId() ?? 0 } }"
                                        class="font-medium text-primary-400">
                                        {{ document?.getCreator()?.getFirstname() }}
                                        {{ document?.getCreator()?.getLastname() }}
                                    </router-link>
                                </p>
                            </div>
                            <div class="flex mt-4 space-x-3 md:mt-0">
                                <router-link :to="{ name: 'Documents: Edit', params: { id: document?.getId() ?? 0 } }"
                                    type="button"
                                    class="inline-flex justify-center gap-x-1.5 rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400"
                                    v-can="'DocStoreService.CreateDocument'">
                                    <PencilIcon class="-ml-0.5 w-5 h-auto" aria-hidden="true" />
                                    Edit
                                </router-link>
                            </div>
                        </div>
                        <div class="flex flex-row gap-2">
                            <div v-if="document?.getClosed()"
                                class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-error-100">
                                <LockClosedIcon class="w-5 h-5 text-error-500" aria-hidden="true" />
                                <span class="text-sm font-medium text-error-700">Closed</span>
                            </div>
                            <div v-else class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-success-100">
                                <LockOpenIcon class="w-5 h-5 text-green-500" aria-hidden="true" />
                                <span class="text-sm font-medium text-green-700">Open</span>
                            </div>
                            <div
                                class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-primary-100 text-primary-500">
                                <ChatBubbleLeftEllipsisIcon class="w-5 h-auto" aria-hidden="true" />
                                <span class="text-sm font-medium text-primary-700">{{ comments.length }}
                                    comments</span>
                            </div>
                            <div class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-base-100 text-base-500">
                                <CalendarIcon class="w-5 h-auto" aria-hidden="true" />
                                <span class="text-sm font-medium text-base-700"><time
                                        :datetime="getDate(document?.getCreatedAt())?.toLocaleString('de-DE')">{{
                                            getDate(document?.getCreatedAt())?.toLocaleString('de-DE')
                                        }}</time></span>
                            </div>
                        </div>
                        <div class="flex flex-row gap-2 pb-3 mt-2 overflow-x-auto sm:pb-0">
                            <div v-for="entry in access?.getJobsList()" :key="entry.getId()"
                                class="flex flex-row items-center flex-initial gap-1 px-2 py-1 rounded-full bg-info-100 whitespace-nowrap">
                                <span class="w-2 h-2 rounded-full bg-info-500" aria-hidden="true" />
                                <span class="text-sm font-medium text-info-800">{{ entry.getJobLabel() }}<span
                                        v-if="entry.getMinimumgrade() > 0"> (Rank: {{ entry.getMinimumgrade() }})</span> -
                                    {{
                                        toTitleCase(DOC_ACCESS_Util.toEnumKey(entry.getAccess())!.toLowerCase()) }}</span>
                            </div>
                            <div v-for="entry in access?.getUsersList()" :key="entry.getId()"
                                class="flex flex-row items-center flex-initial gap-1 px-2 py-1 rounded-full bg-secondary-100 whitespace-nowrap">
                                <span class="w-2 h-2 rounded-full bg-secondary-400" aria-hidden="true" />
                                <span class="text-sm font-medium text-secondary-700">{{ entry.getUser()?.getFirstname() }}
                                    {{ entry.getUser()?.getLastname() }} - {{
                                        toTitleCase(DOC_ACCESS_Util.toEnumKey(entry.getAccess())!.toLowerCase()) }}</span>
                            </div>
                        </div>
                        <div>
                            <h2 class="sr-only">Content</h2>
                            <div class="p-2 mt-4 rounded-lg max-w-none text-neutral bg-base-800">
                                <p v-html="document?.getContent()"></p>
                            </div>
                        </div>
                        <div v-if="activeResponse" class="mt-10">
                            <div class="md:flex md:items-center md:justify-between md:space-x-4 xl:border-b xl:pb-6">
                                <div>
                                    <h1 class="text-2xl font-bold text-gray-900">{{ activeResponse?.getTitle() }}
                                    </h1>
                                    <p class="mt-2 text-sm text-gray-500">
                                        Reply by
                                        {{ ' ' }}
                                        <router-link
                                            :to="{ name: 'Citizens: Info', params: { id: activeResponse?.getCreator()?.getUserId() ?? 0 } }"
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
                        <div class="mt-2">
                            <TabGroup>
                                <TabList class="flex flex-row">
                                    <Tab v-for="tab in tabs" :key="tab.name" v-slot="{ selected }"
                                        class="flex-initial w-full">
                                        <button
                                            :class="[selected ? 'border-primary-500 text-primary-500' : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-400', 'group inline-flex items-center border-b-2 py-4 px-1 text-m font-medium w-full justify-center transition-colors']"
                                            :aria-current="selected ? 'page' : undefined">
                                            <component :is="tab.icon"
                                                :class="[selected ? 'text-primary-500' : 'text-gray-500 group-hover:text-gray-400', '-ml-0.5 mr-2 h-5 w-5 transition-colors']"
                                                aria-hidden="true" />
                                            <span>{{ tab.name }}</span>
                                        </button>
                                    </Tab>
                                </TabList>
                                <TabPanels>
                                    <TabPanel>
                                        <DocumentReferences :references="feedReferences" />
                                    </TabPanel>
                                    <TabPanel>
                                        <DocumentRelations :relations="feedRelations" />
                                    </TabPanel>
                                </TabPanels>
                            </TabGroup>
                        </div>
                    </div>
                </div>
            </div>
        </div>

    </div>
</template>
