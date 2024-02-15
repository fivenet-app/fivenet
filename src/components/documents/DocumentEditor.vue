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
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { max, min, required } from '@vee-validate/rules';
import { useThrottleFn, watchDebounced, watchOnce } from '@vueuse/core';
import 'jodit/es5/jodit.min.css';
import { Jodit } from 'jodit';
// @ts-ignore jodit-vue has no (detected) types
import { JoditEditor } from 'jodit-vue';
import type { IJodit } from 'jodit/types/types';
import {
    AccountMultipleIcon,
    CheckIcon,
    ChevronDownIcon,
    ContentSaveIcon,
    FileDocumentIcon,
    LoadingIcon,
    PlusIcon,
} from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import { type TranslateItem } from '~/composables/i18n';
import { useAuthStore } from '~/store/auth';
import { getDocument, getUser, useClipboardStore } from '~/store/clipboard';
import { useCompletorStore } from '~/store/completor';
import { useDocumentEditorStore } from '~/store/documenteditor';
import { useNotificatorStore } from '~/store/notificator';
import { useSettingsStore } from '~/store/settings';
import { AccessLevel, DocumentAccess } from '~~/gen/ts/resources/documents/access';
import { Category } from '~~/gen/ts/resources/documents/category';
import {
    DocContentType,
    DocReference,
    DocRelation,
    DocumentReference,
    DocumentRelation,
} from '~~/gen/ts/resources/documents/documents';
import { Job, JobGrade } from '~~/gen/ts/resources/users/jobs';
import { UserShort } from '~~/gen/ts/resources/users/users';
import { CreateDocumentRequest, UpdateDocumentRequest } from '~~/gen/ts/services/docstore/docstore';
import DocumentAccessEntry from '~/components/documents/DocumentAccessEntry.vue';
import DocumentReferenceManager from '~/components/documents/DocumentReferenceManager.vue';
import DocumentRelationManager from '~/components/documents/DocumentRelationManager.vue';
import { checkDocAccess } from '~/components/documents/helpers';

const props = defineProps<{
    id?: string;
}>();

const { $grpc } = useNuxtApp();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const clipboardStore = useClipboardStore();
const completorStore = useCompletorStore();
const documentStore = useDocumentEditorStore();
const notifications = useNotificatorStore();

const settingsStore = useSettingsStore();
const { documents } = storeToRefs(settingsStore);

const { t } = useI18n();

const route = useRoute();

const maxAccessEntries = 10;

const canEdit = ref(false);

interface FormData {
    title: string;
    state: string;
    public: boolean;
}

const openclose = [
    { id: 0, label: t('common.open', 2), closed: false },
    { id: 1, label: t('common.close', 2), closed: true },
];

const content = ref('');
const doc = ref<{
    closed: { id: number; label: string; closed: boolean };
    public: boolean;
}>({
    closed: openclose[0],
    public: false,
});
const access = ref<
    Map<
        string,
        {
            id: string;
            type: number;
            values: {
                job?: string;
                char?: number;
                accessRole?: AccessLevel;
                minimumGrade?: number;
            };
            required?: boolean;
        }
    >
>(new Map());
const docAccess = ref<DocumentAccess>();
const docCreator = ref<UserShort | undefined>();

const relationManagerShow = ref<boolean>(false);
const relationManagerData = ref<Map<string, DocumentRelation>>(new Map());
const currentRelations = ref<Readonly<DocumentRelation>[]>([]);
watch(currentRelations, () => currentRelations.value.forEach((e) => relationManagerData.value.set(e.id!, e)));

const referenceManagerShow = ref<boolean>(false);
const referenceManagerData = ref<Map<string, DocumentReference>>(new Map());
const currentReferences = ref<Readonly<DocumentReference>[]>([]);
watch(currentReferences, () => currentReferences.value.forEach((e) => referenceManagerData.value.set(e.id!, e)));

const entriesCategories = ref<Category[]>([]);
const queryCategories = ref('');
const selectedCategory = ref<Category | undefined>(undefined);

const templateId = ref<undefined | string>();

onMounted(async () => {
    if (route.query.templateId) {
        const data = clipboardStore.getTemplateData();
        data.activeChar = activeChar.value!;
        console.debug('Documents: Editor - Clipboard Template Data', data);

        templateId.value = route.query.templateId as string;

        try {
            const call = $grpc.getDocStoreClient().getTemplate({
                templateId: templateId.value as string,
                data,
                render: true,
            });
            const { response } = await call;

            if (response.template === undefined) {
                throw new Error('failed to get template from server response');
            }

            const template = response.template;
            setFieldValue('title', template.contentTitle!);
            setFieldValue('state', template.state!);
            content.value = template.content.replace(/\s+/g, ' ')!;
            selectedCategory.value = template?.category;

            if (template?.contentAccess) {
                if (authStore.activeChar !== null) {
                    docCreator.value = authStore.activeChar;
                }
                const docAccess = template.contentAccess!;
                let accessId = 0;
                docAccess.users.forEach((user) => {
                    const id = accessId.toString();
                    access.value.set(id, {
                        id,
                        type: 0,
                        values: { char: user.userId, accessRole: user.access },
                        required: user.required,
                    });
                    accessId++;
                });

                docAccess.jobs.forEach((job) => {
                    const id = accessId.toString();
                    access.value.set(id, {
                        id,
                        type: 1,
                        values: {
                            job: job.job,
                            accessRole: job.access,
                            minimumGrade: job.minimumGrade,
                        },
                        required: job.required,
                    });
                    accessId++;
                });
            }
        } catch (e) {
            $grpc.handleError(e as RpcError);
            console.error('Documents: Editor - Template Error', e);

            await navigateTo({ name: 'documents' });

            return;
        }
    } else if (props.id) {
        try {
            const req = { documentId: props.id };
            const call = $grpc.getDocStoreClient().getDocument(req);
            const { response } = await call;
            const document = response.document;
            docAccess.value = response.access;
            docCreator.value = document?.creator;

            if (document) {
                setFieldValue('title', document.title);
                setFieldValue('state', document.state);
                content.value = document.content;
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

            if (response.access) {
                let accessId = 0;

                response.access.users.forEach((user) => {
                    const id = accessId.toString();
                    access.value.set(id, {
                        id,
                        type: 0,
                        values: { char: user.userId, accessRole: user.access },
                    });
                    accessId++;
                });

                response.access.jobs.forEach((job) => {
                    const id = accessId.toString();
                    access.value.set(id, {
                        id,
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
            content.value = documentStore.$state.content;
            if (documentStore.$state.closed) {
                doc.value.closed = documentStore.$state.closed;
            }
        }

        const accessId = 0;
        access.value.set(accessId.toString(), {
            id: accessId.toString(),
            type: 1,
            values: {
                job: activeChar.value?.job,
                minimumGrade: 1,
                accessRole: AccessLevel.EDIT,
            },
        });
    }

    clipboardStore.activeStack.documents.forEach((doc, i) => {
        const id = i.toString();
        referenceManagerData.value.set(id, {
            id,
            sourceDocumentId: props.id ?? '0',
            targetDocumentId: doc.id!,
            targetDocument: getDocument(doc),
            creatorId: activeChar.value!.userId,
            creator: activeChar.value!,
            reference: DocReference.SOLVES,
        });
    });
    clipboardStore.activeStack.users.forEach((user, i) => {
        const id = i.toString();
        relationManagerData.value.set(id, {
            id,
            documentId: props.id ?? '0',
            targetUserId: user.userId!,
            targetUser: getUser(user),
            sourceUserId: activeChar.value!.userId,
            sourceUser: activeChar.value!,
            relation: DocRelation.CAUSED,
        });
    });

    setTimeout(() => {
        setupCheckboxes();
    }, 25);
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
        content: content.value,
        state: values.state,
        closed: doc.value.closed,
        category: selectedCategory.value,
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

watchOnce(content, () => (changed.value = true));
watchDebounced(
    doc.value,
    async () => {
        if (changed.value) {
            saveToStore(values);
        }
    },
    {
        debounce: 1000,
        maxWait: 2500,
    },
);
watchDebounced(
    content,
    async () => {
        if (changed.value) {
            saveToStore(values);
        }
    },
    {
        debounce: 1000,
        maxWait: 3500,
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

function addDocumentAccessEntry(): void {
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

    const id = access.value.size > 0 ? parseInt([...access.value.keys()]?.pop() ?? '1', 10) + 1 : 0;
    access.value.set(id.toString(), {
        id: id.toString(),
        type: 1,
        values: {},
    });
}

function removeDocumentAccessEntry(event: { id: string }): void {
    access.value.delete(event.id);
}

function updateDocumentAccessEntryType(event: { id: string; type: number }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.type = event.type;
    access.value.set(event.id, accessEntry);
}

function updateDocumentAccessEntryName(event: { id: string; job?: Job; char?: UserShort }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    if (event.job) {
        accessEntry.values.job = event.job.name;
        accessEntry.values.char = undefined;
    } else if (event.char) {
        accessEntry.values.job = undefined;
        accessEntry.values.char = event.char.userId;
    }

    access.value.set(event.id, accessEntry);
}

function updateDocumentAccessEntryRank(event: { id: string; rank: JobGrade }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.values.minimumGrade = event.rank.grade;
    access.value.set(event.id, accessEntry);
}

function updateDocumentAccessEntryAccess(event: { id: string; access: AccessLevel }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.values.accessRole = event.access;
    access.value.set(event.id, accessEntry);
}

async function createDocument(values: FormData, content: string, closed: boolean): Promise<void> {
    // Prepare request
    const req: CreateDocumentRequest = {
        title: values.title,
        content,
        contentType: DocContentType.HTML,
        closed,
        state: values.state,
        public: doc.value.public,
        templateId: templateId.value,
    };
    if (selectedCategory.value !== undefined) req.categoryId = selectedCategory.value.id;

    const reqAccess: DocumentAccess = {
        jobs: [],
        users: [],
    };
    access.value.forEach((entry) => {
        if (entry.values.accessRole === undefined) {
            return;
        }

        if (entry.type === 0) {
            if (!entry.values.char) {
                return;
            }

            reqAccess.users.push({
                id: '0',
                documentId: '0',
                userId: entry.values.char,
                access: entry.values.accessRole,
            });
        } else if (entry.type === 1) {
            if (!entry.values.job) {
                return;
            }

            reqAccess.jobs.push({
                id: '0',
                documentId: '0',
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
        if (canDo.value.references) {
            referenceManagerData.value.forEach((ref) => {
                ref.sourceDocumentId = response.documentId;

                const prom = $grpc.getDocStoreClient().addDocumentReference({
                    reference: ref,
                });
                promises.push(prom.response);
            });
        }

        if (canDo.value.relations) {
            relationManagerData.value.forEach((rel) => {
                rel.documentId = response.documentId;

                const prom = $grpc.getDocStoreClient().addDocumentRelation({
                    relation: rel,
                });
                promises.push(prom.response);
            });
        }
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
            params: { id: response.documentId },
        });
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function updateDocument(id: string, values: FormData, content: string, closed: boolean): Promise<void> {
    const req: UpdateDocumentRequest = {
        documentId: id,
        title: values.title,
        content,
        contentType: DocContentType.HTML,
        closed,
        state: values.state,
        public: doc.value.public,
    };
    if (selectedCategory.value !== undefined) req.categoryId = selectedCategory.value.id;

    const reqAccess: DocumentAccess = {
        jobs: [],
        users: [],
    };
    access.value.forEach((entry) => {
        if (entry.values.accessRole === undefined) {
            return;
        }

        if (entry.type === 0) {
            if (!entry.values.char) {
                return;
            }

            reqAccess.users.push({
                id: '0',
                documentId: id,
                userId: entry.values.char,
                access: entry.values.accessRole,
            });
        } else if (entry.type === 1) {
            if (!entry.values.job) {
                return;
            }

            reqAccess.jobs.push({
                id: '0',
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

        if (canDo.value.references) {
            const referencesToRemove: string[] = [];
            currentReferences.value.forEach((ref) => {
                if (!referenceManagerData.value.has(ref.id!)) referencesToRemove.push(ref.id!);
            });
            referencesToRemove.forEach((id) => {
                $grpc.getDocStoreClient().removeDocumentReference({ id });
            });
            referenceManagerData.value.forEach((ref) => {
                if (currentReferences.value.find((r) => r.id === ref.id!)) {
                    return;
                }
                ref.sourceDocumentId = response.documentId;

                $grpc.getDocStoreClient().addDocumentReference({
                    reference: ref,
                });
            });
        }

        if (canDo.value.relations) {
            const relationsToRemove: string[] = [];
            currentRelations.value.forEach((rel) => {
                if (!relationManagerData.value.has(rel.id!)) relationsToRemove.push(rel.id!);
            });
            relationsToRemove.forEach((id) => {
                $grpc.getDocStoreClient().removeDocumentRelation({ id });
            });
            relationManagerData.value.forEach((rel) => {
                if (currentRelations.value.find((r) => r.id === rel.id!)) {
                    return;
                }
                rel.documentId = response.documentId;

                $grpc.getDocStoreClient().addDocumentRelation({
                    relation: rel,
                });
            });
        }

        notifications.dispatchNotification({
            title: { key: 'notifications.document_updated.title', parameters: {} },
            content: { key: 'notifications.document_updated.content', parameters: {} },
            type: 'success',
        });
        clipboardStore.clear();
        documentStore.clear();

        await navigateTo({
            name: 'documents-id',
            params: { id: response.documentId },
        });
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const canDo = computed(() => ({
    edit:
        props.id === undefined
            ? true
            : checkDocAccess(docAccess.value, docCreator.value, AccessLevel.EDIT, 'DocStoreService.UpdateDocument'),
    access:
        props.id === undefined
            ? true
            : checkDocAccess(docAccess.value, docCreator.value, AccessLevel.ACCESS, 'DocStoreService.UpdateDocument'),
    references: can('DocStoreService.AddDocumentReference'),
    relations: can('DocStoreService.AddDocumentRelation'),
}));

console.info(
    'Documents Editor - Can Do: Edit',
    canDo.value.edit,
    'Access',
    canDo.value.access,
    'References',
    canDo.value.references,
    'Relations',
    canDo.value.relations,
);

defineRule('required', required);
defineRule('max', max);
defineRule('min', min);

const { handleSubmit, values, setFieldValue, meta } = useForm<FormData>({
    validationSchema: {
        title: { required: true, min: 3, max: 255 },
        state: { required: false, min: 2, max: 32 },
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
        prom = createDocument(values, content.value, doc.value.closed.closed);
    } else {
        prom = updateDocument(props.id, values, content.value, doc.value.closed.closed);
    }

    await prom.finally(() => setTimeout(() => (canSubmit.value = true), 400));
});
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);

const config = {
    language: 'de',
    spellcheck: true,
    minHeight: 475,
    editorClassName: 'prose' + (documents.value.editorTheme === 'dark' ? ' prose-neutral' : ' prose-gray'),
    theme: documents.value.editorTheme,

    readonly: false,
    defaultActionOnPaste: 'insert_clear_html',
    disablePlugins: ['about', 'poweredByJodit', 'classSpan', 'file', 'video', 'print'],
    // Uploader Plugin
    uploader: {
        insertImageAsBase64URI: true,
    },
    // Clean HTML Plugin
    cleanHTML: {
        denyTags: 'script,iframe,form,button,svg',
        fillEmptyParagraph: false,
    },
    nl2brInPlainText: true,
    // Inline Plugin
    toolbarInline: true,
    toolbarInlineForSelection: true,
    toolbarInlineDisableFor: [],
    toolbarInlineDisabledButtons: ['source'],
    popup: {
        a: Jodit.atom(['link', 'unlink']),
    },
    // Link Plugin
    link: {
        /**
         * Template for the link dialog form
         */
        formTemplate: (_: Jodit) => `<form><input ref="url_input"><button>Apply</button></form>`,
        formClassName: 'some-class',
        /**
         * Follow link address after dblclick
         */
        followOnDblClick: true,
        /**
         * Replace inserted youtube/vimeo link to `iframe`
         */
        processVideoLink: false,
        /**
         * Wrap inserted link
         */
        processPastedLink: true,
        /**
         * Show `no follow` checkbox in link dialog.
         */
        noFollowCheckbox: false,
        /**
         * Show `Open in new tab` checkbox in link dialog.
         */
        openInNewTabCheckbox: false,
        /**
         * Use an input text to ask the classname or a select or not ask
         */
        modeClassName: 'input', // 'select'
        /**
         * Allow multiple choises (to use with modeClassName="select")
         */
        selectMultipleClassName: true,
        /**
         * The size of the select (to use with modeClassName="select")
         */
        selectSizeClassName: 10,
        /**
         * The list of the option for the select (to use with modeClassName="select")
         */
        selectOptionsClassName: [],
    },
};

const plugins = [
    {
        name: 'focus',
        callback: (editor: IJodit) => {
            editor.e
                .on('blur', () => {
                    focusTablet(false);
                })
                .on('focus', () => {
                    focusTablet(true);
                });
        },
    },
];

const extraButtons = [
    '|',
    {
        name: 'insertCheckbox',
        iconURL: '/images/icons/format-list-checkbox.svg',
        exec: function (editor: IJodit) {
            const label = document.createElement('label');
            label.setAttribute('contenteditable', 'false');
            const empty = document.createElement('span');
            empty.innerHTML = '&nbsp;';

            const input = document.createElement('input');
            input.setAttribute('type', 'checkbox');
            input.setAttribute('checked', 'true');
            input.onchange = (ev) => {
                if (ev.target === null) {
                    return;
                }
                setCheckboxState(ev.target as HTMLInputElement);
            };

            label.appendChild(input);
            label.appendChild(empty);

            editor.s.insertHTML(label, true);
        },
    },
];

function setCheckboxState(target: HTMLInputElement): void {
    const attr = target.getAttribute('checked');
    const checked = attr !== null ? Boolean(attr) : false;
    if (checked) {
        target.removeAttribute('checked');
    } else {
        target.setAttribute('checked', 'true');
    }
}

function setupCheckboxes(): void {
    const checkboxes: NodeListOf<HTMLInputElement> = document.querySelectorAll('.jodit-wysiwyg input[type=checkbox]');
    checkboxes.forEach(
        (el) =>
            (el.onchange = (ev) => {
                if (ev.target === null) {
                    return;
                }
                setCheckboxState(ev.target as HTMLInputElement);
            }),
    );
}
</script>

<template>
    <div class="m-2">
        <form @submit.prevent="onSubmitThrottle">
            <DocumentRelationManager
                v-model="relationManagerData"
                :open="relationManagerShow"
                :document="id"
                @close="relationManagerShow = false"
            />
            <DocumentReferenceManager
                v-model="referenceManagerData"
                :open="referenceManagerShow"
                :document-id="id"
                @close="referenceManagerShow = false"
            />

            <div
                class="flex flex-col gap-2 rounded-t-lg bg-base-800 px-3 py-4 text-neutral"
                :class="!(canDo.edit && canDo.relations && canDo.references) ? 'rounded-b-md' : ''"
            >
                <div>
                    <label for="title" class="block text-base font-medium">
                        {{ $t('common.title') }}
                    </label>
                    <VeeField
                        name="title"
                        type="text"
                        :placeholder="$t('common.title')"
                        :label="$t('common.title')"
                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-3xl sm:leading-6"
                        :disabled="!canEdit || !canDo.edit"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                    <VeeErrorMessage name="title" as="p" class="mt-2 text-sm text-error-400" />
                </div>
                <div class="flex flex-row gap-2">
                    <div class="flex-1">
                        <label for="category" class="block text-sm font-medium">
                            {{ $t('common.category') }}
                        </label>
                        <Combobox v-model="selectedCategory" as="div" :disabled="!canEdit || !canDo.edit" nullable>
                            <div class="relative">
                                <ComboboxButton as="div">
                                    <ComboboxInput
                                        autocomplete="off"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        :display-value="(category: any) => category?.name"
                                        @change="queryCategories = $event.target.value"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </ComboboxButton>

                                <ComboboxOptions
                                    v-if="entriesCategories.length > 0"
                                    class="absolute z-10 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                                >
                                    <ComboboxOption
                                        v-for="category in entriesCategories"
                                        :key="category.id"
                                        v-slot="{ active, selected }"
                                        :value="category"
                                        as="category"
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
                                                <CheckIcon class="h-5 w-5" aria-hidden="true" />
                                            </span>
                                        </li>
                                    </ComboboxOption>
                                </ComboboxOptions>
                            </div>
                        </Combobox>
                    </div>
                    <div class="flex-1">
                        <label for="state" class="block text-sm font-medium">
                            {{ $t('common.state') }}
                        </label>
                        <VeeField
                            name="state"
                            type="text"
                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                            :placeholder="`${$t('common.document', 1)} ${$t('common.state')}`"
                            :label="`${$t('common.document', 1)} ${$t('common.state')}`"
                            :disabled="!canEdit || !canDo.edit"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                        <VeeErrorMessage name="state" as="p" class="mt-2 text-sm text-error-400" />
                    </div>
                    <div class="flex-1">
                        <label for="closed" class="block text-sm font-medium"> {{ $t('common.close', 2) }}? </label>
                        <Listbox v-model="doc.closed" as="div" :disabled="!canEdit || !canDo.edit">
                            <div class="relative">
                                <ListboxButton
                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 pl-3 text-left text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                >
                                    <span class="block truncate">
                                        {{ openclose.find((e) => e.closed === doc.closed.closed)?.label }}</span
                                    >
                                    <span class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
                                        <ChevronDownIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
                                    </span>
                                </ListboxButton>

                                <transition
                                    leave-active-class="transition duration-100 ease-in"
                                    leave-from-class="opacity-100"
                                    leave-to-class="opacity-0"
                                >
                                    <ListboxOptions
                                        class="absolute z-10 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                                    >
                                        <ListboxOption
                                            v-for="st in openclose"
                                            :key="st.closed.toString()"
                                            v-slot="{ active, selected }"
                                            as="template"
                                            :value="st"
                                        >
                                            <li
                                                :class="[
                                                    active ? 'bg-primary-500' : '',
                                                    'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
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
                                                    <CheckIcon class="h-5 w-5" aria-hidden="true" />
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
            <div v-if="canDo.edit" class="bg-neutral">
                <JoditEditor v-model="content" :config="config" :plugins="plugins" :extra-buttons="extraButtons" />
                <template v-if="saving">
                    <div class="flex animate-pulse justify-center">
                        <ContentSaveIcon class="mr-2 h-auto w-4 animate-spin" />
                        {{ $t('common.save', 2) }}...
                    </div>
                </template>
            </div>
            <div v-if="canDo.edit" class="flex flex-row">
                <div class="inline-flex flex-1 rounded-md shadow-sm" role="group">
                    <button
                        v-if="canDo.relations"
                        type="button"
                        :disabled="!canEdit || !canDo.edit"
                        class="inline-flex w-full justify-center rounded-bl-md bg-primary-500 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                        :class="canDo.references ? '' : 'rounded-br-md'"
                        @click="relationManagerShow = true"
                    >
                        <div class="flex justify-center">
                            <AccountMultipleIcon
                                class="-ml-0.5 mr-2 h-5 w-5 text-base-300 transition-colors group-hover:text-accent-200"
                                aria-hidden="true"
                            />
                            {{ $t('common.citizen', 1) }} {{ $t('common.relation', 2) }}
                        </div>
                    </button>
                    <button
                        v-if="canDo.references"
                        type="button"
                        :disabled="!canEdit || !canDo.edit"
                        class="inline-flex w-full justify-center rounded-br-md bg-primary-500 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                        :class="canDo.relations ? '' : 'rounded-bl-md'"
                        @click="referenceManagerShow = true"
                    >
                        <div class="flex justify-center">
                            <FileDocumentIcon
                                class="-ml-0.5 mr-2 h-5 w-5 text-base-300 transition-colors group-hover:text-accent-200"
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
                <DocumentAccessEntry
                    v-for="entry in access.values()"
                    :key="entry.id"
                    :init="entry"
                    :access-types="accessTypes"
                    :read-only="!canDo.access || entry.required === true"
                    @type-change="updateDocumentAccessEntryType($event)"
                    @name-change="updateDocumentAccessEntryName($event)"
                    @rank-change="updateDocumentAccessEntryRank($event)"
                    @access-change="updateDocumentAccessEntryAccess($event)"
                    @delete-request="removeDocumentAccessEntry($event)"
                />
                <button
                    type="button"
                    :disabled="!canEdit || !canDo.access"
                    class="rounded-full bg-primary-500 p-2 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                    data-te-toggle="tooltip"
                    :title="$t('components.documents.document_editor.add_permission')"
                    @click="addDocumentAccessEntry()"
                >
                    <PlusIcon class="h-5 w-5" aria-hidden="true" />
                </button>
            </div>
            <div class="flex pb-14">
                <button
                    type="submit"
                    :disabled="!meta.valid || !canEdit || !canSubmit"
                    class="flex w-full justify-center rounded-md px-3.5 py-2.5 text-sm font-semibold text-neutral"
                    :class="[
                        !canEdit || !meta.valid || !canSubmit
                            ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                            : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                    ]"
                >
                    <template v-if="!canSubmit">
                        <LoadingIcon class="mr-2 h-5 w-5 animate-spin" />
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

<style>
.jodit-wysiwyg {
    min-width: 100%;

    * {
        margin-top: 4px;
        margin-bottom: 4px;
    }
}
</style>
