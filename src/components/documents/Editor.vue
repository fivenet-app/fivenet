<script lang="ts" setup>
import {
    Combobox,
    ComboboxButton,
    ComboboxInput,
    ComboboxOption,
    ComboboxOptions,
    Listbox,
    ListboxButton,
    ListboxOption,
    ListboxOptions,
} from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { max, min, required } from '@vee-validate/rules';
import { QuillEditor } from '@vueup/vue-quill';
import '@vueup/vue-quill/dist/vue-quill.snow.css';
import { useDebounceFn, useThrottleFn, watchDebounced, watchOnce } from '@vueuse/core';
import { ImageActions } from '@xeger/quill-image-actions';
import { ImageFormats } from '@xeger/quill-image-formats';
import {
    AccountMultipleIcon,
    CheckIcon,
    ChevronDownIcon,
    ContentSaveIcon,
    FileDocumentIcon,
    LoadingIcon,
    PlusIcon,
} from 'mdi-vue3';
import htmlEditButton from 'quill-html-edit-button';
import ImageCompress from 'quill-image-compress';
import MagicUrl from 'quill-magic-url';
// @ts-expect-error
import QuillPasteSmart from 'quill-paste-smart';
import { defineRule } from 'vee-validate';
import { TranslateItem } from '~/composables/i18n';
import '~/composables/quill/divider/divider';
import '~/composables/quill/divider/quill-divider';
import DividerToolbar from '~/composables/quill/divider/quill-divider';
import { useAuthStore } from '~/store/auth';
import { getDocument, getUser, useClipboardStore } from '~/store/clipboard';
import { useCompletorStore } from '~/store/completor';
import { useDocumentEditorStore } from '~/store/documenteditor';
import { useNotificatorStore } from '~/store/notificator';
import { AccessLevel } from '~~/gen/ts/resources/documents/access';
import { Category } from '~~/gen/ts/resources/documents/category';
import {
    DocContentType,
    DocReference,
    DocRelation,
    DocumentAccess,
    DocumentReference,
    DocumentRelation,
} from '~~/gen/ts/resources/documents/documents';
import { Job, JobGrade } from '~~/gen/ts/resources/users/jobs';
import { UserShort } from '~~/gen/ts/resources/users/users';
import { CreateDocumentRequest, UpdateDocumentRequest } from '~~/gen/ts/services/docstore/docstore';
import AccessEntry from './AccessEntry.vue';
import ReferenceManager from './ReferenceManager.vue';
import RelationManager from './RelationManager.vue';

const props = defineProps<{
    id?: bigint;
}>();

const { $grpc } = useNuxtApp();

const authStore = useAuthStore();
const clipboardStore = useClipboardStore();
const documentStore = useDocumentEditorStore();
const notifications = useNotificatorStore();
const completorStore = useCompletorStore();

const { t } = useI18n();

const route = useRoute();

const { activeChar } = storeToRefs(authStore);

const maxAccessEntries = 8;

const canEdit = ref(false);

const openclose = [
    { id: 0, label: t('common.open', 2), closed: false },
    { id: 1, label: t('common.close', 2), closed: true },
];

const doc = ref<{
    content: string;
    closed: { id: number; label: string; closed: boolean };
    public: boolean;
}>({
    content: '',
    closed: openclose[0],
    public: false,
});
const access = ref<
    Map<
        bigint,
        {
            id: bigint;
            type: number;
            values: {
                job?: string;
                char?: number;
                accessRole?: AccessLevel;
                minimumGrade?: number;
            };
        }
    >
>(new Map());

const relationManagerShow = ref<boolean>(false);
const relationManagerData = ref<Map<bigint, DocumentRelation>>(new Map());
const currentRelations = ref<Readonly<DocumentRelation>[]>([]);
watch(currentRelations, () => currentRelations.value.forEach((e) => relationManagerData.value.set(e.id!, e)));

const referenceManagerShow = ref<boolean>(false);
const referenceManagerData = ref<Map<bigint, DocumentReference>>(new Map());
const currentReferences = ref<Readonly<DocumentReference>[]>([]);
watch(currentReferences, () => currentReferences.value.forEach((e) => referenceManagerData.value.set(e.id!, e)));

const entriesCategories = ref<Category[]>([]);
const queryCategories = ref('');
const selectedCategory = ref<Category | undefined>(undefined);

// Quill Editor modules and options
const modules = [
    {
        name: 'clipboard',
        module: QuillPasteSmart,
        options: {
            allowed: {
                tags: ['a', 'b', 'strong', 'u', 's', 'i', 'p', 'br', 'ul', 'ol', 'li', 'span'],
                attributes: ['href', 'rel', 'target', 'class'],
            },
            keepSelection: true,
            substituteBlockElements: true,
            magicPasteLinks: false,
        },
    },
    {
        name: 'imageFormats',
        module: ImageFormats,
    },
    {
        name: 'imageActions',
        module: ImageActions,
    },
    {
        name: 'imageCompress',
        module: ImageCompress,
        options: {
            quality: 0.7,
            maxWidth: 1250,
            maxHeight: 1250,
            imageType: 'image/png',
            keepImageTypes: ['image/jpeg', 'image/png'],
            debug: false,
            suppressErrorLogging: false,
            insertIntoEditor: undefined,
        },
    },
    {
        name: 'magicUrl',
        module: MagicUrl,
        options: {
            normalizeUrlOptions: {
                stripHash: true,
            },
        },
    },
    {
        name: 'htmlEditButton',
        module: htmlEditButton,
        options: {
            debug: false,
            msg: t('components.documents.document_editor.quill.msg'),
            okText: t('components.documents.document_editor.quill.okText'),
            cancelText: t('common.cancel'),
        },
    },
    {
        name: 'divider',
        module: DividerToolbar,
        options: {},
    },
];

const formats = [
    'bold',
    'italic',
    'underline',
    'strike',
    'blockquote',
    'code-block',
    'code',
    'header',
    'list',
    'script',
    'indent',
    'direction',
    'size',
    'color',
    'background',
    'font',
    'align',
    'float',
    'link',
    'video',
    'image',
    'height',
    'width',
    'divider',
];
const toolbarOptions = [
    ['bold', 'italic', 'underline', 'strike'], // toggled buttons
    ['blockquote', 'code-block'],

    [{ header: 1 }, { header: 2 }], // custom button values
    [{ list: 'ordered' }, { list: 'bullet' }],
    [{ script: 'sub' }, { script: 'super' }], // superscript/subscript
    [{ indent: '-1' }, { indent: '+1' }], // outdent/indent
    [{ direction: 'rtl' }], // text direction

    [{ size: ['small', false, 'large', 'huge'] }], // custom dropdown
    [{ header: [1, 2, 3, 4, 5, 6, false] }],

    [{ color: [] }, { background: [] }], // dropdown with defaults from theme
    [{ font: [] }],
    [{ align: [] }, 'divider'],

    ['link', 'video', 'image'],

    ['clean'], // remove formatting button
];
const options = {
    readOnly: false,
    contentType: 'html',
    theme: 'snow',
    formats: formats,
    toolbar: toolbarOptions,
};

onMounted(async () => {
    if (route.query.templateId) {
        const data = clipboardStore.getTemplateData();
        data.activeChar = activeChar.value!;
        console.debug('Clipboard Template Data', data, JSON.stringify(data));

        try {
            const call = $grpc.getDocStoreClient().getTemplate({
                templateId: BigInt(route.query.templateId as string),
                data: JSON.stringify(data),
                render: true,
            });
            const { response } = await call;

            const template = response.template;
            setFieldValue('title', template?.contentTitle!);
            setFieldValue('state', template?.state!);
            doc.value.content = template?.content!;
            selectedCategory.value = template?.category;

            if (template?.contentAccess) {
                const docAccess = template?.contentAccess!;
                let accessId = 0n;
                docAccess.users.forEach((user) => {
                    access.value.set(accessId, {
                        id: accessId,
                        type: 0,
                        values: { char: user.userId, accessRole: user.access },
                    });
                    accessId++;
                });

                docAccess.jobs.forEach((job) => {
                    access.value.set(accessId, {
                        id: accessId,
                        type: 1,
                        values: {
                            job: job.job,
                            accessRole: job.access,
                            minimumGrade: job.minimumGrade,
                        },
                    });
                    accessId++;
                });
            }
        } catch (e) {
            $grpc.handleError(e as RpcError);

            await navigateTo({ name: 'documents' });

            return;
        }
    } else if (props.id) {
        try {
            const req = { documentId: props.id };
            const call = $grpc.getDocStoreClient().getDocument(req);
            const { response } = await call;
            const document = response.document;
            const docAccess = response.access;

            if (document) {
                setFieldValue('title', document.title);
                setFieldValue('state', document.state);
                doc.value.content = document.content;
                doc.value.closed = openclose.find((e) => e.closed === document.closed) as {
                    id: number;
                    label: string;
                    closed: boolean;
                };
                selectedCategory.value = document.category;
                doc.value.public = document.public;

                const refs = await $grpc.getDocStoreClient().getDocumentReferences(req);
                currentReferences.value = refs.response.references;
                const rels = await $grpc.getDocStoreClient().getDocumentRelations(req);
                currentRelations.value = rels.response.relations;
            }

            if (docAccess) {
                let accessId = 0n;

                docAccess.users.forEach((user) => {
                    access.value.set(accessId, {
                        id: accessId,
                        type: 0,
                        values: { char: user.userId, accessRole: user.access },
                    });
                    accessId++;
                });

                docAccess.jobs.forEach((job) => {
                    access.value.set(accessId, {
                        id: accessId,
                        type: 1,
                        values: {
                            job: job.job,
                            accessRole: job.access,
                            minimumGrade: job.minimumGrade,
                        },
                    });
                    accessId++;
                });
            }
        } catch (e) {
            $grpc.handleError(e as RpcError);

            await navigateTo({ name: 'documents' });

            return;
        }
    } else {
        if (documentStore.$state) {
            setFieldValue('title', documentStore.$state.title);
            setFieldValue('state', documentStore.$state.state);
            doc.value.content = documentStore.$state.content;
            if (documentStore.$state.closed) {
                doc.value.closed = documentStore.$state.closed;
            }
        }

        let accessId = 0n;
        access.value.set(accessId, {
            id: accessId,
            type: 1,
            values: {
                job: activeChar.value?.job,
                minimumGrade: 1,
                accessRole: AccessLevel.EDIT,
            },
        });
    }

    clipboardStore.activeStack.documents.forEach((doc, i) => {
        const id = BigInt(i);
        referenceManagerData.value.set(id, {
            id: id,
            sourceDocumentId: props.id ?? 0n,
            targetDocumentId: doc.id!,
            targetDocument: getDocument(doc),
            creatorId: activeChar.value!.userId,
            creator: activeChar.value!,
            reference: DocReference.SOLVES,
        });
    });
    clipboardStore.activeStack.users.forEach((user, i) => {
        const id = BigInt(i);
        relationManagerData.value.set(id, {
            id: id,
            documentId: props.id ?? 0n,
            targetUserId: user.userId!,
            targetUser: getUser(user),
            sourceUserId: activeChar.value!.userId,
            sourceUser: activeChar.value!,
            relation: DocRelation.CAUSED,
        });
    });

    canEdit.value = true;

    findCategories();
});

const saving = ref(false);

async function saveToStore(values: FormData): Promise<void> {
    if (saving.value) {
        return;
    }
    saving.value = true;

    documentStore.save({
        title: values.title,
        content: doc.value.content,
        state: values.state,
        closed: doc.value.closed,
    });
    setTimeout(() => {
        saving.value = false;
    }, 1250);
}

async function findCategories(): Promise<void> {
    entriesCategories.value = await completorStore.completeDocumentCategories(queryCategories.value);
    if (selectedCategory.value && entriesCategories.value.findIndex((c) => c.id === selectedCategory.value?.id) !== 0)
        entriesCategories.value.push(selectedCategory.value);
}

const changed = ref(false);
watchDebounced(
    doc.value,
    async () => {
        if (changed.value) {
            saveToStore(values);
        }
    },
    {
        debounce: 1350,
        maxWait: 4000,
    },
);

watchDebounced(queryCategories, async () => findCategories(), {
    debounce: 600,
    maxWait: 1400,
});

const accessTypes = [
    { id: 0, name: t('common.citizen', 2) },
    { id: 1, name: t('common.job', 2) },
];

function addAccessEntry(): void {
    if (access.value.size > maxAccessEntries - 1) {
        notifications.dispatchNotification({
            title: { key: 'notifications.max_access_entry.title', parameters: {} },
            content: {
                key: 'notifications.max_access_entry.content',
                parameters: { max: maxAccessEntries.toString() },
            } as TranslateItem,
            type: 'error',
        });
        return;
    }

    const id = access.value.size > 0 ? ([...access.value.keys()].pop() as bigint) + 1n : 0n;
    access.value.set(id, {
        id,
        type: 1,
        values: {},
    });
}

function removeAccessEntry(event: { id: bigint }): void {
    access.value.delete(event.id);
}

function updateAccessEntryType(event: { id: bigint; type: number }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) return;

    accessEntry.type = event.type;
    access.value.set(event.id, accessEntry);
}

function updateAccessEntryName(event: { id: bigint; job?: Job; char?: UserShort }): void {
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

function updateAccessEntryRank(event: { id: bigint; rank: JobGrade }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) return;

    accessEntry.values.minimumGrade = event.rank.grade;
    access.value.set(event.id, accessEntry);
}

function updateAccessEntryAccess(event: { id: bigint; access: AccessLevel }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) return;

    accessEntry.values.accessRole = event.access;
    access.value.set(event.id, accessEntry);
}

async function createDocument(values: FormData, content: string, closed: boolean): Promise<void> {
    return new Promise(async (res, rej) => {
        // Prepare request
        const req: CreateDocumentRequest = {
            title: values.title,
            content: content,
            contentType: DocContentType.HTML,
            closed: closed,
            state: values.state,
            public: doc.value.public,
        };
        if (selectedCategory.value !== undefined) req.categoryId = selectedCategory.value.id;

        const reqAccess: DocumentAccess = {
            jobs: [],
            users: [],
        };
        access.value.forEach((entry) => {
            if (entry.values.accessRole === undefined) return;

            if (entry.type === 0) {
                if (!entry.values.char) return;

                reqAccess.users.push({
                    id: 0n,
                    documentId: 0n,
                    userId: entry.values.char,
                    access: entry.values.accessRole,
                });
            } else if (entry.type === 1) {
                if (!entry.values.job) return;

                reqAccess.jobs.push({
                    id: 0n,
                    documentId: 0n,
                    job: entry.values.job,
                    minimumGrade: entry.values.minimumGrade ? entry.values.minimumGrade : 0,
                    access: entry.values.accessRole,
                });
            }
        });
        req.access = reqAccess;

        // Try to submit to server
        try {
            const call = $grpc.getDocStoreClient().createDocument(req);
            const { response } = await call;

            const promises: Promise<any>[] = [];
            referenceManagerData.value.forEach((ref) => {
                ref.sourceDocumentId = response.documentId;

                const prom = $grpc.getDocStoreClient().addDocumentReference({
                    reference: ref,
                });
                promises.push(prom.response);
            });

            relationManagerData.value.forEach((rel) => {
                rel.documentId = response.documentId;

                const prom = $grpc.getDocStoreClient().addDocumentRelation({
                    relation: rel,
                });
                promises.push(prom.response);
            });
            await Promise.all(promises);

            notifications.dispatchNotification({
                title: { key: 'notifications.document_created.title', parameters: {} },
                content: { key: 'notifications.document_created.content', parameters: {} },
                type: 'success',
            });
            clipboardStore.clear();
            documentStore.clear();

            await navigateTo({
                name: 'documents-id',
                params: { id: response.documentId.toString() },
            });

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function updateDocument(id: bigint, values: FormData, content: string, closed: boolean): Promise<void> {
    return new Promise(async (res, rej) => {
        const req: UpdateDocumentRequest = {
            documentId: id,
            title: values.title,
            content: content,
            contentType: DocContentType.HTML,
            closed: closed,
            state: values.state,
            public: doc.value.public,
        };
        if (selectedCategory.value !== undefined) req.categoryId = selectedCategory.value.id;

        const reqAccess: DocumentAccess = {
            jobs: [],
            users: [],
        };
        access.value.forEach((entry) => {
            if (entry.values.accessRole === undefined) return;

            if (entry.type === 0) {
                if (!entry.values.char) return;

                reqAccess.users.push({
                    id: 0n,
                    documentId: id,
                    userId: entry.values.char,
                    access: entry.values.accessRole,
                });
            } else if (entry.type === 1) {
                if (!entry.values.job) return;

                reqAccess.jobs.push({
                    id: 0n,
                    documentId: id,
                    job: entry.values.job,
                    minimumGrade: entry.values.minimumGrade ? entry.values.minimumGrade : 0,
                    access: entry.values.accessRole,
                });
            }
        });
        req.access = reqAccess;

        try {
            const call = $grpc.getDocStoreClient().updateDocument(req);
            const { response } = await call;

            const referencesToRemove: bigint[] = [];
            currentReferences.value.forEach((ref) => {
                if (!referenceManagerData.value.has(ref.id!)) referencesToRemove.push(ref.id!);
            });
            referencesToRemove.forEach((id) => {
                $grpc.getDocStoreClient().removeDocumentReference({
                    id: id,
                });
            });

            const relationsToRemove: bigint[] = [];
            currentRelations.value.forEach((rel) => {
                if (!relationManagerData.value.has(rel.id!)) relationsToRemove.push(rel.id!);
            });
            relationsToRemove.forEach((id) => {
                $grpc.getDocStoreClient().removeDocumentRelation({
                    id: id,
                });
            });

            referenceManagerData.value.forEach((ref) => {
                if (currentReferences.value.find((r) => r.id === ref.id!)) return;
                ref.sourceDocumentId = response.documentId;

                $grpc.getDocStoreClient().addDocumentReference({
                    reference: ref,
                });
            });

            relationManagerData.value.forEach((rel) => {
                if (currentRelations.value.find((r) => r.id === rel.id!)) return;
                rel.documentId = response.documentId;

                $grpc.getDocStoreClient().addDocumentRelation({
                    relation: rel,
                });
            });

            notifications.dispatchNotification({
                title: { key: 'notifications.document_updated.title', parameters: {} },
                content: { key: 'notifications.document_updated.content', parameters: {} },
                type: 'success',
            });
            clipboardStore.clear();
            documentStore.clear();

            await navigateTo({
                name: 'documents-id',
                params: { id: response.documentId.toString() },
            });

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

type Stats = {
    words: number;
    chars: number;
};

const stats = ref<Stats>({
    words: 0,
    chars: 0,
});

const quillEditorRef = ref<null | InstanceType<typeof QuillEditor>>(null);

function calculate(content: string): Stats {
    const stats: Stats = {
        words: 0,
        chars: 0,
    };

    if (content.length > 0) {
        stats.words = content.split(/\s+/).length - 1;
        stats.chars = content.split(/\S/).length - 1;
    }

    return stats;
}

const debounced = useDebounceFn(
    () => {
        if (quillEditorRef.value === null) return;
        changed.value = true;

        stats.value = calculate(quillEditorRef.value.getQuill()?.getText());
    },
    750,
    { maxWait: 1500 },
);

watchOnce(quillEditorRef, () => {
    if (quillEditorRef.value === null) return;

    quillEditorRef.value.getQuill().on('text-change', debounced);
});

defineRule('required', required);
defineRule('max', max);
defineRule('min', min);

interface FormData {
    title: string;
    state: string;
    public: boolean;
}

const { handleSubmit, values, setFieldValue, meta } = useForm<FormData>({
    validationSchema: {
        title: { required: true, min: 3, max: 21845 },
        state: { required: false, min: 2, max: 24 },
    },
    initialValues: {
        public: false,
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(async (values): Promise<void> => {
    let prom: Promise<void>;
    if (props.id === undefined) {
        prom = createDocument(values, doc.value.content, doc.value.closed.closed);
    } else {
        prom = updateDocument(props.id, values, doc.value.content, doc.value.closed.closed);
    }

    await prom.finally(() => setTimeout(() => (canSubmit.value = true), 350));
});
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);
</script>

<style>
.ql-container {
    border: none !important;
    height: auto !important;
}

.ql-editor {
    height: 100%;
    width: 100%;
    min-height: 400px;
}

.ql-hidden {
    display: none;
}

.ql-editor img,
.ql-editor svg,
.ql-editor video,
.ql-editor canvas,
.ql-editor audio,
.ql-editor iframe,
.ql-editor embed,
.ql-editor object {
    display: inline !important;
}

.ql-editor hr {
    color: black;
}
</style>

<template>
    <div class="m-2">
        <form @submit.prevent="onSubmitThrottle">
            <RelationManager
                v-model="relationManagerData"
                :open="relationManagerShow"
                :document="$props.id"
                @close="relationManagerShow = false"
            />
            <ReferenceManager
                v-model="referenceManagerData"
                :open="referenceManagerShow"
                :document="$props.id"
                @close="referenceManagerShow = false"
            />

            <div class="flex flex-col gap-2 px-3 py-4 rounded-t-lg bg-base-800 text-neutral">
                <div>
                    <label for="title" class="block font-medium text-base">
                        {{ $t('common.title') }}
                    </label>
                    <VeeField
                        name="title"
                        type="text"
                        :placeholder="$t('common.title')"
                        :label="$t('common.title')"
                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-3xl sm:leading-6"
                        :disabled="!canEdit"
                    />
                    <VeeErrorMessage name="title" as="p" class="mt-2 text-sm text-error-400" />
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
                                        @change="queryCategories = $event.target.value"
                                        :display-value="(category: any) => category?.name"
                                    />
                                </ComboboxButton>

                                <ComboboxOptions
                                    v-if="entriesCategories.length > 0"
                                    class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-44 sm:text-sm"
                                >
                                    <ComboboxOption
                                        v-for="category in entriesCategories"
                                        :key="category.id?.toString()"
                                        :value="category"
                                        as="category"
                                        v-slot="{ active, selected }"
                                    >
                                        <li
                                            :class="[
                                                'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                                                active ? 'bg-primary-500' : '',
                                            ]"
                                        >
                                            <span :class="['block truncate', selected && 'font-semibold']">
                                                {{ category.name }}
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
                    <div class="flex-1">
                        <label for="state" class="block font-medium text-sm">
                            {{ $t('common.state') }}
                        </label>
                        <VeeField
                            name="state"
                            type="text"
                            class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                            :placeholder="`${$t('common.document', 1)} ${$t('common.state')}`"
                            :label="`${$t('common.document', 1)} ${$t('common.state')}`"
                            :disabled="!canEdit"
                        />
                        <VeeErrorMessage name="state" as="p" class="mt-2 text-sm text-error-400" />
                    </div>
                    <div class="flex-1">
                        <label for="closed" class="block font-medium text-sm"> {{ $t('common.close', 2) }}? </label>
                        <Listbox as="div" v-model="doc.closed">
                            <div class="relative">
                                <ListboxButton
                                    :disabled="!canEdit"
                                    class="block pl-3 text-left w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                >
                                    <span class="block truncate">
                                        {{ openclose.find((e) => e.closed === doc.closed.closed)?.label }}</span
                                    >
                                    <span class="absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none">
                                        <ChevronDownIcon class="w-5 h-5 text-gray-400" aria-hidden="true" />
                                    </span>
                                </ListboxButton>

                                <transition
                                    leave-active-class="transition duration-100 ease-in"
                                    leave-from-class="opacity-100"
                                    leave-to-class="opacity-0"
                                >
                                    <ListboxOptions
                                        class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-44 sm:text-sm"
                                    >
                                        <ListboxOption
                                            as="template"
                                            v-for="st in openclose"
                                            :key="st.closed?.toString()"
                                            :value="st"
                                            v-slot="{ active, selected }"
                                        >
                                            <li
                                                :class="[
                                                    active ? 'bg-primary-500' : '',
                                                    'text-neutral relative cursor-default select-none py-2 pl-8 pr-4',
                                                ]"
                                            >
                                                <span :class="[selected ? 'font-semibold' : 'font-normal', 'block truncate']">{{
                                                    st.label
                                                }}</span>

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
                                        </ListboxOption>
                                    </ListboxOptions>
                                </transition>
                            </div>
                        </Listbox>
                    </div>
                </div>
            </div>
            <div class="bg-neutral">
                <QuillEditor
                    ref="quillEditorRef"
                    v-model:content="doc.content"
                    content-type="html"
                    :toolbar="toolbarOptions"
                    :modules="modules"
                    :options="options"
                />
                <div class="grid grid-cols-2 text-base text-gray-600 h-7 mx-2">
                    <div class="flex flex-items items-center">
                        <template v-if="saving">
                            <ContentSaveIcon class="w-6 h-auto mr-2 animate-spin" />
                            {{ $t('common.save', 2) }}...
                        </template>
                    </div>
                    <div class="text-end">
                        {{ $t('common.char', 2) }}: {{ stats.chars }}, {{ $t('common.word', 2) }}: {{ stats.words }}
                    </div>
                </div>
            </div>
            <div class="flex flex-row">
                <div class="flex-1 inline-flex rounded-md shadow-sm" role="group">
                    <button
                        type="button"
                        :disabled="!canEdit"
                        class="inline-flex justify-center rounded-bl-md bg-primary-500 py-2.5 px-3.5 w-full text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                        @click="relationManagerShow = true"
                    >
                        <div class="flex justify-center">
                            <AccountMultipleIcon
                                class="text-base-300 group-hover:text-base-200 -ml-0.5 mr-2 h-5 w-5 transition-colors"
                                aria-hidden="true"
                            />
                            {{ $t('common.citizen', 1) }} {{ $t('common.relation', 2) }}
                        </div>
                    </button>
                    <button
                        type="button"
                        :disabled="!canEdit"
                        class="inline-flex justify-center rounded-br-md bg-primary-500 py-2.5 px-3.5 w-full text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                        @click="referenceManagerShow = true"
                    >
                        <div class="flex justify-center">
                            <FileDocumentIcon
                                class="text-base-300 group-hover:text-base-200 -ml-0.5 mr-2 h-5 w-5 transition-colors"
                                aria-hidden="true"
                            />
                            {{ $t('common.document', 1) }} {{ $t('common.reference', 2) }}
                        </div>
                    </button>
                </div>
            </div>
            <div class="my-3">
                <h2 class="text-neutral">
                    {{ $t('common.access') }}
                </h2>
                <AccessEntry
                    v-for="entry in access.values()"
                    :key="entry.id?.toString()"
                    :init="entry"
                    :access-types="accessTypes"
                    @typeChange="updateAccessEntryType($event)"
                    @nameChange="updateAccessEntryName($event)"
                    @rankChange="updateAccessEntryRank($event)"
                    @accessChange="updateAccessEntryAccess($event)"
                    @deleteRequest="removeAccessEntry($event)"
                />
                <button
                    type="button"
                    :disabled="!canEdit"
                    class="p-2 rounded-full bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                    data-te-toggle="tooltip"
                    :title="$t('components.documents.document_editor.add_permission')"
                    @click="addAccessEntry()"
                >
                    <PlusIcon class="w-5 h-5" aria-hidden="true" />
                </button>
            </div>
            <div class="flex pb-14">
                <button
                    type="submit"
                    :disabled="!meta.valid || !canEdit || !canSubmit"
                    class="flex justify-center rounded-md py-2.5 px-3.5 text-sm font-semibold text-neutral w-full"
                    :class="[
                        !canEdit || !meta.valid || !canSubmit
                            ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                            : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                    ]"
                >
                    <template v-if="!canSubmit">
                        <LoadingIcon class="animate-spin h-5 w-5 mr-2" />
                    </template>
                    <template v-if="!props.id">
                        {{ t('common.create') }}
                    </template>
                    <template v-else>
                        {{ $t('common.save') }}
                    </template>
                </button>
            </div>
        </form>
    </div>
</template>
