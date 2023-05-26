<script lang="ts" setup>
import { useAuthStore } from '~/store/auth';
import { useDocumentEditorStore } from '~/store/documenteditor';
import { useClipboardStore, getUser } from '~/store/clipboard';
import { Quill, QuillEditor } from '@vueup/vue-quill';
import '@vueup/vue-quill/dist/vue-quill.snow.css';
import { AddDocumentRelationRequest, CreateDocumentRequest, GetDocumentRequest, RemoveDocumentRelationRequest, UpdateDocumentRequest, RemoveDocumentReferenceRequest, AddDocumentReferenceRequest, GetTemplateRequest } from '~~/gen/ts/services/docstore/docstore';
import { DocumentAccess, DocumentJobAccess, DocumentReference, DocumentRelation, DocumentUserAccess, DOC_CONTENT_TYPE, DOC_RELATION } from '~~/gen/ts/resources/documents/documents';
import { DocumentCategory } from '~~/gen/ts/resources/documents/category';
import {
    PlusIcon,
    ChevronDownIcon,
    CheckIcon,
} from '@heroicons/vue/20/solid';
import { Job, JobGrade } from '~~/gen/ts/resources/jobs/jobs';
import { UserShort } from '~~/gen/ts/resources/users/users';
import {
    Listbox,
    ListboxButton,
    ListboxOption,
    ListboxOptions,
    Combobox,
    ComboboxButton,
    ComboboxInput,
    ComboboxOption,
    ComboboxOptions
} from '@headlessui/vue';
import { CompleteDocumentCategoriesRequest } from '~~/gen/ts/services/completor/completor';
import { watchDebounced } from '@vueuse/core';
import { ACCESS_LEVEL_Util } from '~~/gen/ts/resources/documents/access.pb_enums';
import DocumentReferenceManager from './DocumentReferenceManager.vue';
import DocumentRelationManager from './DocumentRelationManager.vue';
import DocumentAccessEntry from './DocumentAccessEntry.vue';
import { ArrowPathIcon } from '@heroicons/vue/24/solid';
import { useNotificationsStore } from '~/store/notifications';
import { ACCESS_LEVEL } from '~~/gen/ts/resources/documents/access';

const { $grpc } = useNuxtApp();
const authStore = useAuthStore();
const clipboardStore = useClipboardStore();
const documentStore = useDocumentEditorStore();
const notifications = useNotificationsStore();

const { t } = useI18n();

const route = useRoute();

const props = defineProps({
    id: {
        required: false,
        type: Number
    },
});

const { activeChar } = storeToRefs(authStore);

const maxAccessEntries = 8;

const canEdit = ref(false);

const openclose = [
    { id: 0, label: t('common.open'), closed: false },
    { id: 1, label: t('common.close', 2), closed: true },
];

const doc = ref<{ title: string, content: string, closed: { id: number, label: string, closed: boolean }, state: string }>({
    title: '',
    content: '',
    closed: openclose[0],
    state: '',
});
const isPublic = ref(false);
const access = ref<Map<number, { id: number, type: number, values: { job?: string, char?: number, accessrole?: ACCESS_LEVEL, minimumrank?: number } }>>(new Map());

const relationManagerShow = ref<boolean>(false);
const relationManagerData = ref<Map<number, DocumentRelation>>(new Map());
const currentRelations = ref<Readonly<DocumentRelation>[]>([]);
watch(currentRelations, () => currentRelations.value.forEach(e => relationManagerData.value.set(e.id, e)))

const referenceManagerShow = ref<boolean>(false);
const referenceManagerData = ref<Map<number, DocumentReference>>(new Map());
const currentReferences = ref<Readonly<DocumentReference>[]>([]);
watch(currentReferences, () => currentReferences.value.forEach(e => referenceManagerData.value.set(e.id, e)))

let entriesCategory = [] as DocumentCategory[];
const queryCategory = ref('');
const selectedCategory = ref<DocumentCategory | undefined>(undefined);

const modules = [] as Quill.Module[];

onMounted(async () => {
    await findCategories();

    if (route.query.templateId) {
        const req = new GetTemplateRequest();
        req.setTemplateId(parseInt(route.query.templateId as string));
        req.setRender(true);

        const data = clipboardStore.getTemplateData();
        data.setActivechar(activeChar.value!);
        req.setData(JSON.stringify(data.toObject()));

        try {
            const resp = await $grpc.getDocStoreClient().
                getTemplate(req);

            const template = resp.getTemplate();
            doc.value.title = template?.getContentTitle()!;
            doc.value.content = template?.getContent()!;
            selectedCategory.value = entriesCategory.find(e => e.id === template?.category?.id);

            if (template?.hasContentAccess()) {
                const docAccess = template?.getContentAccess()!;
                let accessId = 0;
                docAccess.getUsersList().forEach(user => {
                    access.value.set(accessId, { id: accessId, type: 0, values: { char: user.userId, accessrole: user.getAccess() } });
                    accessId++;
                });

                docAccess.getJobsList().forEach(job => {
                    access.value.set(accessId, { id: accessId, type: 1, values: { job: job.job, accessrole: job.getAccess(), minimumrank: job.getMinimumgrade() } });
                    accessId++;
                });
            }
        } catch (e) {
        }
    } else if (props.id) {
        const req = new GetDocumentRequest();
        req.setDocumentId(props.id);

        $grpc.getDocStoreClient().getDocument(req).then(async (resp) => {
            const document = resp.document;
            const docAccess = resp.getAccess();

            if (document) {
                doc.value.title = document.title;
                doc.value.content = document.getContent();
                doc.value.closed = openclose.find(e => e.closed === document.closed) as { id: number; label: string; closed: boolean; };
                doc.value.state = document.state;
                selectedCategory.value = entriesCategory.find(e => e.id === document.category?.id);
                isPublic.value = document.getPublic();

                const refs = await $grpc.getDocStoreClient().getDocumentReferences(req);
                currentReferences.value = refs.getReferencesList();
                const rels = await $grpc.getDocStoreClient().getDocumentRelations(req);
                currentRelations.value = rels.getRelationsList();
            };

            if (docAccess) {
                let accessId = 0;

                docAccess.getUsersList().forEach(user => {
                    access.value.set(accessId, { id: accessId, type: 0, values: { char: user.userId, accessrole: user.getAccess() } });
                    accessId++;
                });

                docAccess.getJobsList().forEach(job => {
                    access.value.set(accessId, { id: accessId, type: 1, values: { job: job.job, accessrole: job.getAccess(), minimumrank: job.getMinimumgrade() } });
                    accessId++;
                });
            }

        });
    } else {
        if (documentStore.$state) {
            doc.value.title = documentStore.$state.title;
            doc.value.content = documentStore.$state.content;
            doc.value.state = documentStore.$state.state;
            if (documentStore.$state.closed) {
                doc.value.closed = documentStore.$state.closed;
            }
        }

        access.value.set(0, { id: 0, type: 1, values: { job: activeChar.value?.job, minimumrank: 1, accessrole: ACCESS_LEVEL.EDIT } });
    }

    clipboardStore.users.forEach((user, i) => {
        const rel = new DocumentRelation();
        rel.setId(i);
        rel.setDocumentId(props.id!);
        rel.setTargetUserId(user.id!);
        rel.setTargetUser(getUser(user));
        rel.setSourceUserId(activeChar.value!.userId);
        rel.setSourceUser(activeChar.value!);
        rel.setRelation(DOC_RELATION.CAUSED);

        relationManagerData.value.set(i, rel);
    })

    canEdit.value = true;
});

const saving = ref(false);

function saveToStore(): void {
    if (saving.value) {
        return;
    }
    saving.value = true;

    documentStore.save(doc.value);
    setTimeout(() => {
        saving.value = false;
    }, 850);
}

watchDebounced(doc.value, () => saveToStore(), { debounce: 1250, maxWait: 3500 });

watchDebounced(queryCategory, () => findCategories(), { debounce: 600, maxWait: 1400 });

async function findCategories(): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const req = new CompleteDocumentCategoriesRequest();
            req.setSearch(queryCategory.value);

            const resp = await $grpc.getCompletorClient().completeDocumentCategories(req)
            entriesCategory = resp.getCategoriesList();

            return res();
        } catch (e) {
            return rej(e as RpcError);
        }
    });
}

const accessTypes = [
    { id: 0, name: t('common.citizen', 2) },
    { id: 1, name: t('common.job', 2) },
];

function addAccessEntry(): void {
    if (access.value.size > maxAccessEntries - 1) {
        notifications.dispatchNotification({
            title: t('notifications.max_access_entry.title'),
            content: t('notifications.max_access_entry.content', [maxAccessEntries]),
            type: 'error'
        });
        return;
    }

    let id = access.value.size > 0 ? [...access.value.keys()].pop() as number + 1 : 0;
    access.value.set(id, {
        id,
        type: 1,
        values: {}
    })
}

function removeAccessEntry(event: {
    id: number
}): void {
    access.value.delete(event.id);
}

function updateAccessEntryType(event: {
    id: number,
    type: number
}): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) return;

    accessEntry.type = event.type;
    access.value.set(event.id, accessEntry);
}

function updateAccessEntryName(event: {
    id: number,
    job?: Job,
    char?: UserShort
}): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) return;

    if (event.job) {
        accessEntry.values.job = event.job.name;
        accessEntry.values.char = undefined;
    } else if (event.char) {
        accessEntry.values.job = undefined;
        accessEntry.values.char = event.char.userId;
    }

    access.value.set(event.id, accessEntry);
}

function updateAccessEntryRank(event: {
    id: number,
    rank: JobGrade
}): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) return;

    accessEntry.values.minimumrank = event.rank.getGrade();
    access.value.set(event.id, accessEntry);
}

function updateAccessEntryAccess(event: {
    id: number,
    access: ACCESS_LEVEL
}): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) return;

    accessEntry.values.accessrole = event.access;
    access.value.set(event.id, accessEntry);
}

async function submitForm(): Promise<void> {
    return new Promise(async (res, rej) => {
        // Prepare request
        const req = new CreateDocumentRequest();
        req.setTitle(doc.value.title);
        req.setContent(doc.value.content);
        req.setContentType(DOC_CONTENT_TYPE.HTML);
        req.setClosed(doc.value.closed.closed);
        req.setState(doc.value.state);
        req.setPublic(isPublic.value);
        if (selectedCategory.value != undefined)
            req.setCategoryId(selectedCategory.value.id);

        const reqAccess = new DocumentAccess();
        access.value.forEach(entry => {
            if (entry.values.accessrole === undefined) return;

            if (entry.type === 0) {
                if (!entry.values.char) return;

                const user = new DocumentUserAccess();
                user.setUserId(entry.values.char);
                user.setAccess(ACCESS_LEVEL_Util.fromInt(entry.values.accessrole));

                reqAccess.addUsers(user);
            } else if (entry.type === 1) {
                if (!entry.values.job) return;

                const job = new DocumentJobAccess();
                job.setJob(entry.values.job);
                job.setMinimumgrade(entry.values.minimumrank ? entry.values.minimumrank : 0);
                job.setAccess(ACCESS_LEVEL_Util.fromInt(entry.values.accessrole));

                reqAccess.addJobs(job);
            }
        });
        req.setAccess(reqAccess);

        // Try to submit to server
        try {
            const resp = await $grpc.getDocStoreClient().
                createDocument(req);

            const promises = new Array<Promise<any>>();
            referenceManagerData.value.forEach((ref) => {
                ref.setSourceDocumentId(resp.documentId);

                const req = new AddDocumentReferenceRequest();
                req.setReference(ref);
                promises.push($grpc.getDocStoreClient().addDocumentReference(req));
            });

            relationManagerData.value.forEach((rel) => {
                rel.setDocumentId(resp.documentId);

                const req = new AddDocumentRelationRequest();
                req.setRelation(rel);
                promises.push($grpc.getDocStoreClient().addDocumentRelation(req));
            });
            await Promise.all(promises);

            notifications.dispatchNotification({
                title: t('notifications.document_created.title'),
                content: t('notifications.document_created.content'),
                type: 'success',
            });
            clipboardStore.clearActiveStack();
            documentStore.clear();

            await navigateTo({ name: 'documents-id', params: { id: resp.documentId } });

            return res();
        } catch (e) {
            return rej(e as RpcError);
        }
    });
}

async function editForm(): Promise<void> {
    return new Promise(async (res, rej) => {
        const req = new UpdateDocumentRequest();
        req.setDocumentId(props.id!);
        req.setTitle(doc.value.title);
        req.setContent(doc.value.content);
        req.setContentType(DOC_CONTENT_TYPE.HTML);
        req.setClosed(doc.value.closed.closed);
        req.setState(doc.value.state);
        req.setPublic(isPublic.value);
        if (selectedCategory.value != undefined)
            req.setCategoryId(selectedCategory.value.id);

        const reqAccess = new DocumentAccess();
        access.value.forEach(entry => {
            if (entry.values.accessrole === undefined) return;

            if (entry.type === 0) {
                if (!entry.values.char) return;

                const user = new DocumentUserAccess();
                user.setAccess(ACCESS_LEVEL_Util.fromInt(entry.values.accessrole));
                user.setUserId(entry.values.char);
                if (activeChar.value) user.setCreatorId(activeChar.value.userId);

                reqAccess.addUsers(user);
            } else if (entry.type === 1) {
                if (!entry.values.job) return;

                const job = new DocumentJobAccess();
                job.setJob(entry.values.job);
                job.setMinimumgrade(entry.values.minimumrank ? entry.values.minimumrank : 0);
                job.setAccess(ACCESS_LEVEL_Util.fromInt(entry.values.accessrole));
                if (activeChar.value) job.setCreatorId(activeChar.value.userId);

                reqAccess.addJobs(job);
            }
        });
        req.setAccess(reqAccess);

        try {
            const resp = await $grpc.getDocStoreClient().
                updateDocument(req);

            const referencesToRemove: number[] = [];
            currentReferences.value.forEach((ref) => {
                if (!referenceManagerData.value.has(ref.id)) referencesToRemove.push(ref.id);
            });
            referencesToRemove.forEach((id) => {
                const req = new RemoveDocumentReferenceRequest();
                req.setId(id);
                $grpc.getDocStoreClient().removeDocumentReference(req);
            });

            const relationsToRemove: number[] = [];
            currentRelations.value.forEach((rel) => {
                if (!relationManagerData.value.has(rel.id)) relationsToRemove.push(rel.id);
            });
            relationsToRemove.forEach((id) => {
                const req = new RemoveDocumentRelationRequest();
                req.setId(id);
                $grpc.getDocStoreClient().removeDocumentRelation(req);
            });

            referenceManagerData.value.forEach((ref) => {
                if (currentReferences.value.find(r => r.id === ref.id)) return;
                ref.setSourceDocumentId(resp.documentId);

                const req = new AddDocumentReferenceRequest();
                req.setReference(ref);
                $grpc.getDocStoreClient().addDocumentReference(req);
            });

            relationManagerData.value.forEach((rel) => {
                if (currentRelations.value.find(r => r.id === rel.id)) return;
                rel.setDocumentId(resp.documentId);

                const req = new AddDocumentRelationRequest();
                req.setRelation(rel);
                $grpc.getDocStoreClient().addDocumentRelation(req);
            });

            notifications.dispatchNotification({
                title: t('notifications.document_updated.title'),
                content: t('notifications.document_updated.content'),
                type: 'success'
            });
            clipboardStore.clearActiveStack();
            documentStore.clear();

            await navigateTo({ name: 'documents-id', params: { id: resp.documentId } });
            return res();
        } catch (e) {
            return rej(e as RpcError);
        }
    });
}
</script>

<style>
.ql-container {
    border: none !important;
}
</style>

<template>
    <DocumentRelationManager v-model="relationManagerData" :open="relationManagerShow" :document="$props.id"
        @close="relationManagerShow = false" />
    <DocumentReferenceManager v-model="referenceManagerData" :open="referenceManagerShow" :document="$props.id"
        @close="referenceManagerShow = false" />
    <div class="flex flex-col gap-2 px-3 py-4 rounded-t-lg bg-base-800 text-neutral">
        <div>
            <label for="name" class="block font-medium text-base">
                {{ $t('common.title') }}
            </label>
            <input v-model="doc.title" type="text" name="name"
                class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-3xl sm:leading-6"
                :placeholder="`${$t('common.document', 1)} ${$t('common.title')}`" :disabled="!canEdit" />
        </div>
        <div class="flex flex-row gap-2">
            <div class="flex-1">
                <label for="category" class="block font-medium text-sm">
                    {{ $t('common.category') }}
                </label>
                <Combobox as="div" v-model="selectedCategory" :disabled="!canEdit" nullable>
                    <div class="relative">
                        <ComboboxButton as="div">
                            <ComboboxInput
                                class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                @change="queryCategory = $event.target.value"
                                :display-value="(category: any) => category?.name" />
                        </ComboboxButton>

                        <ComboboxOptions v-if="entriesCategory.length > 0"
                            class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-60 sm:text-sm">
                            <ComboboxOption v-for="category in entriesCategory" :key="category.id" :value="category"
                                as="category" v-slot="{ active, selected }">
                                <li
                                    :class="['relative cursor-default select-none py-2 pl-8 pr-4 text-neutral', active ? 'bg-primary-500' : '']">
                                    <span :class="['block truncate', selected && 'font-semibold']">
                                        {{ category.name }}
                                    </span>

                                    <span v-if="selected"
                                        :class="[active ? 'text-neutral' : 'text-primary-500', 'absolute inset-y-0 left-0 flex items-center pl-1.5']">
                                        <CheckIcon class="w-5 h-5" aria-hidden="true" />
                                    </span>
                                </li>
                            </ComboboxOption>
                        </ComboboxOptions>
                    </div>
                </Combobox>
            </div>
            <div class="flex-1">
                <label for="name" class="block font-medium text-sm">
                    {{ $t('common.state') }}
                </label>
                <input v-model="doc.state" type="text" name="state"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    :placeholder="`${$t('common.document', 1)} ${$t('common.state')}`" :disabled="!canEdit" />
            </div>
            <div class="flex-1">
                <label for="closed" class="block font-medium text-sm">
                    {{ $t('common.close', 2) }}?
                </label>
                <Listbox as="div" v-model="doc.closed">
                    <div class="relative">
                        <ListboxButton :disabled="!canEdit"
                            class="block pl-3 text-left w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6">
                            <span class="block truncate">{{ openclose.find(e => e.closed === doc.closed.closed)?.label
                            }}</span>
                            <span class="absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none">
                                <ChevronDownIcon class="w-5 h-5 text-gray-400" aria-hidden="true" />
                            </span>
                        </ListboxButton>

                        <transition leave-active-class="transition duration-100 ease-in" leave-from-class="opacity-100"
                            leave-to-class="opacity-0">
                            <ListboxOptions
                                class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-60 sm:text-sm">
                                <ListboxOption as="template" v-for="type in openclose" :key="type.id" :value="type"
                                    v-slot="{ active, selected }">
                                    <li
                                        :class="[active ? 'bg-primary-500' : '', 'text-neutral relative cursor-default select-none py-2 pl-8 pr-4']">
                                        <span :class="[selected ? 'font-semibold' : 'font-normal', 'block truncate']">{{
                                            type.label
                                        }}</span>

                                        <span v-if="selected"
                                            :class="[active ? 'text-neutral' : 'text-primary-500', 'absolute inset-y-0 left-0 flex items-center pl-1.5']">
                                            <CheckIcon class="w-5 h-5" aria-hidden="true" />
                                        </span>
                                    </li>
                                </ListboxOption>
                            </ListboxOptions>
                        </transition>
                    </div>
                </Listbox>
            </div>
        </div>
    </div>
    <div class="bg-neutral min-h-[32rem]">
        <QuillEditor v-model:content="doc.content" content-type="html" toolbar="full" theme="snow" :modules="modules" />
    </div>
    <div class="flex flex-row">
        <div class="flex-1">
            <button type="button" :disabled="!canEdit"
                class="rounded-bl-md bg-primary-500 py-2.5 px-3.5 w-full text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                @click="referenceManagerShow = true">
                {{ $t('common.document', 1) }} {{ $t('common.reference', 2) }}
            </button>
        </div>
        <div class="flex-1">
            <button type="button" :disabled="!canEdit"
                class="rounded-br-md bg-primary-500 py-2.5 px-3.5 w-full text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                @click="relationManagerShow = true">
                {{ $t('common.citizen', 1) }} {{ $t('common.relation', 2) }}
            </button>
        </div>
    </div>
    <div class="my-3">
        <h2 class="text-neutral">
            {{ $t('common.access') }}
        </h2>
        <DocumentAccessEntry v-for="entry in access.values()" :key="entry.id" :init="entry" :access-types="accessTypes"
            @typeChange="updateAccessEntryType($event)" @nameChange="updateAccessEntryName($event)"
            @rankChange="updateAccessEntryRank($event)" @accessChange="updateAccessEntryAccess($event)"
            @deleteRequest="removeAccessEntry($event)" />
        <button type="button" :disabled="!canEdit"
            class="p-2 rounded-full bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
            data-te-toggle="tooltip" :title="$t('components.documents.document_editor.add_permission')"
            @click="addAccessEntry()">
            <PlusIcon class="w-5 h-5" aria-hidden="true" />
        </button>
    </div>
    <div class="sm:flex sm:flex-row-reverse">
        <button v-if="!props.id" @click="submitForm()" :disabled="!canEdit"
            class="rounded-md bg-primary-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-primary-400">
            {{ t('common.submit') }}
        </button>
        <button v-if="props.id" @click="editForm()" :disabled="!canEdit"
            class="rounded-md bg-primary-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-primary-400">
            {{ $t('common.edit') }}
        </button>
        <div v-if="saving" class="text-gray-400 mr-4 flex flex-items">
            <ArrowPathIcon class="w-6 h-auto ml-auto mr-2.5 animate-spin" />
            <span class="mt-2">{{ $t('common.save', 2) }}...</span>
        </div>
    </div>
</template>
