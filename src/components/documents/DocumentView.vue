<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { GetDocumentCommentsRequest, GetDocumentRequest } from '@fivenet/gen/services/docstore/docstore_pb';
import { Document, DocumentAccess, DocumentComment, DocumentReference, DocumentRelation } from '@fivenet/gen/resources/documents/documents_pb';
import { toDate } from '~/utils/time';
import { DOC_ACCESS_Util } from '@fivenet/gen/resources/documents/documents.pb_enums';
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
import DocumentComments from './DocumentComments.vue';
import { toTitleCase } from '~/utils/strings';
import { PaginationRequest } from '@fivenet/gen/resources/common/database/database_pb';
import { useClipboardStore } from '~/store/clipboard';
import { PlusIcon } from '@heroicons/vue/24/solid';
import { RpcError } from 'grpc-web';

const { $grpc } = useNuxtApp();
const clipboardStore = useClipboardStore();

const document = ref<undefined | Document>(undefined);
const access = ref<undefined | DocumentAccess>(undefined);
const comments = ref<DocumentComment[]>([]);
const commentCount = ref(0);
const feedReferences = ref<DocumentReference[]>([]);
const feedRelations = ref<DocumentRelation[]>([]);
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

async function getDocument(): Promise<void> {
    const req = new GetDocumentRequest();
    req.setDocumentId(props.documentId);

    $grpc.getDocStoreClient().
        getDocument(req, null).
        then((resp) => {
            document.value = resp.getDocument();
            access.value = resp.getAccess();
        }).catch((err: RpcError) => $grpc.handleRPCError(err));

    // Document References
    $grpc.getDocStoreClient().
        getDocumentReferences(req, null).
        then((resp) => {
            feedReferences.value = resp.getReferencesList();
        }).catch((err: RpcError) => $grpc.handleRPCError(err));

    // Document Relations
    $grpc.getDocStoreClient().
        getDocumentRelations(req, null).
        then((resp) => {
            feedRelations.value = resp.getRelationsList();
        }).catch((err: RpcError) => $grpc.handleRPCError(err));

    // Document Comments
    const creq = new GetDocumentCommentsRequest();
    creq.setPagination((new PaginationRequest()).setOffset(0));
    creq.setDocumentId(props.documentId);

    $grpc.getDocStoreClient().
        getDocumentComments(creq, null).
        then((resp) => {
            comments.value = resp.getCommentsList();
            commentCount.value = resp.getPagination()!.getTotalCount();
        }).catch((err: RpcError) => $grpc.handleRPCError(err));
}

function addToClipboard(): void {
    if (document.value) {
        clipboardStore.addDocument(document.value);
    }
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
                                    <NuxtLink
                                        :to="{ name: 'citizens-id', params: { id: document?.getCreator()?.getUserId() ?? 0 } }"
                                        class="font-medium text-primary-400 hover:text-primary-300">
                                        {{ document?.getCreator()?.getFirstname() }}
                                        {{ document?.getCreator()?.getLastname() }}
                                    </NuxtLink>
                                </p>
                            </div>
                            <div class="flex mt-4 space-x-3 md:mt-0">
                                <NuxtLink :to="{ name: 'documents-edit-id', params: { id: document?.getId() ?? 0 } }"
                                    type="button"
                                    class="inline-flex justify-center gap-x-1.5 rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400"
                                    v-can="'DocStoreService.CreateDocument'">
                                    <PencilIcon class="-ml-0.5 w-5 h-auto" aria-hidden="true" />
                                    Edit
                                </NuxtLink>
                            </div>
                        </div>
                        <div class="flex flex-row gap-2">
                            <div v-if="document?.getClosed()"
                                class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-error-100">
                                <LockClosedIcon class="w-5 h-5 text-error-400" aria-hidden="true" />
                                <span class="text-sm font-medium text-error-700">Closed</span>
                            </div>
                            <div v-else class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-success-100">
                                <LockOpenIcon class="w-5 h-5 text-green-500" aria-hidden="true" />
                                <span class="text-sm font-medium text-green-700">Open</span>
                            </div>
                            <div
                                class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-primary-100 text-primary-500">
                                <ChatBubbleLeftEllipsisIcon class="w-5 h-auto" aria-hidden="true" />
                                <span class="text-sm font-medium text-primary-700">{{ commentCount }}
                                    comments</span>
                            </div>
                            <div class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-base-100 text-base-500">
                                <CalendarIcon class="w-5 h-auto" aria-hidden="true" />
                                <span class="text-sm font-medium text-base-700"><time
                                        :datetime="toDate(document?.getCreatedAt())?.toLocaleString('de-DE')">{{
                                            toDate(document?.getCreatedAt())?.toLocaleString('de-DE')
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
                            <div class="p-2 mt-4 rounded-lg text-neutral bg-base-800 break-words">
                                <p v-html="document?.getContent()"></p>
                            </div>
                        </div>
                        <div>
                            <TabGroup>
                                <TabList class="flex flex-row">
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
                                    <TabPanel>
                                        <DocumentReferences :references="feedReferences" :show-source="false" />
                                    </TabPanel>
                                    <TabPanel>
                                        <DocumentRelations :relations="feedRelations" :show-document="false" />
                                    </TabPanel>
                                </TabPanels>
                            </TabGroup>
                        </div>
                        <div class="mt-4" v-can="'DocStoreService.GetDocumentComments'">
                            <h2 class="text-lg font-semibold text-neutral">Comments</h2>
                            <DocumentComments :document-id="documentId" :comments="comments" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
    <button title="Add to Clipboard" @click="addToClipboard()"
        class="fixed flex items-center justify-center w-12 h-12 rounded-full z-90 bottom-24 right-8 bg-primary-500 shadow-float text-neutral hover:bg-primary-400">
        <PlusIcon class="w-10 h-auto" />
    </button>
</template>
