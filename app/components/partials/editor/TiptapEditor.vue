<script lang="ts" setup>
import type { EditorEmojiMenuItem, FormError } from '@nuxt/ui';
import type { ClientStreamingCall, RpcOptions } from '@protobuf-ts/runtime-rpc';
import { generateJSON, getSchema, type Extensions, type JSONContent } from '@tiptap/core';
import Collaboration from '@tiptap/extension-collaboration';
import CollaborationCaret from '@tiptap/extension-collaboration-caret';
import { gitHubEmojis } from '@tiptap/extension-emoji';
import { UndoRedo } from '@tiptap/extensions';
import type { Schema } from '@tiptap/pm/model';
import { initProseMirrorDoc, prosemirrorJSONToYDoc } from '@tiptap/y-tiptap';
import * as Y from 'yjs';
import { DeleteImageTracker } from '~/composables/tiptap/extensions/DeleteImageTracker';
import { imageUploadPlugin } from '~/composables/tiptap/extensions/ImageUploadPlugin';
import type { UploadNamespaces } from '~/composables/useFileUploader';
import type GrpcProvider from '~/composables/yjs/yjs';
import type { File as FileGrpc } from '~~/gen/ts/resources/file/file';
import type { UploadFileRequest, UploadFileResponse } from '~~/gen/ts/resources/file/filestore';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import TiptapToolbar from './TiptapToolbar.vue';
import YJSUserPopover from './YJSUserPopover.vue';

const props = withDefaults(
    defineProps<{
        name?: string;
        wrapperClass?: string;
        limit?: number;
        fileLimit?: number;
        disabled?: boolean;
        placeholder?: string;
        hideToolbar?: boolean;
        disableImages?: boolean;
        historyType?: string;
        enableCollab?: boolean;

        extensions?: Extensions;

        saving?: boolean;

        targetId?: number;
        filestoreNamespace?: UploadNamespaces;
        filestoreService?: (options?: RpcOptions) => ClientStreamingCall<UploadFileRequest, UploadFileResponse>;
    }>(),
    {
        name: undefined,
        wrapperClass: '',
        limit: 0,
        fileLimit: 10,
        disabled: false,
        placeholder: '',
        hideToolbar: false,
        disableImages: false,
        historyType: undefined,
        enableCollab: false,

        extensions: () => [],

        saving: false,

        targetId: undefined,
        filestoreNamespace: undefined,
        filestoreService: undefined,
    },
);

const emits = defineEmits<{
    (e: 'file-uploaded', file: FileGrpc): void;
}>();

defineOptions({
    inheritAttrs: false,
});

const modelValue = defineModel<JSONContent | string | undefined>({ required: true });
const files = defineModel<FileGrpc[]>('files', { default: () => [] });

const logger = useLogger('📄 Editor' + (props.name ? ` ${props.name}` : ''));

const { activeChar } = useAuth();

const settingsStore = useSettingsStore();
const { editor: editorSettings } = storeToRefs(settingsStore);

const notifications = useNotificationsStore();

const extensions = useTiptapEditor(toRef(props, 'limit'), toRef(props, 'placeholder'));

if (!props.disableImages) {
    extensions.push(
        DeleteImageTracker.configure({
            onRemoved: (ids) =>
                ids.forEach((id) => {
                    if (hasFileById(files.value, id)) {
                        const idx = files.value.findIndex((f) => f.id === id);
                        if (idx > -1) files.value.splice(idx, 1);
                    }
                    logger.info('Removed file:', id);
                }),
        }),
    );
}

const ydoc = inject<Y.Doc | undefined>('yjsDoc', undefined);
const yjsProvider = inject<GrpcProvider | undefined>('yjsProvider', undefined);

const loading = ref(props.enableCollab && ydoc !== undefined && yjsProvider !== undefined);

function seedDocument(schema: Schema, value: JSONContent | string): void {
    let seedDoc: Y.Doc;
    if (typeof value === 'string') {
        if (value === '') return;

        // HTML -> ProseMirror JSON
        const json = generateJSON(value, extensions);
        // ProseMirror JSON -> Yjs update in-place
        seedDoc = prosemirrorJSONToYDoc(schema, json, 'content');
    } else {
        if (!value.content || value.content.length === 0) return;

        logger.info('Seeding document content into Yjs document');

        // ProseMirror JSON -> Yjs update in-place
        seedDoc = prosemirrorJSONToYDoc(schema, value, 'content');
    }

    // Merge that doc's state into the live document
    Y.applyUpdate(ydoc!, Y.encodeStateAsUpdate(seedDoc));
}

let yjsSchema: Schema | undefined = undefined;

if (props.enableCollab && ydoc && yjsProvider) {
    const ourName = `${activeChar.value?.firstname} ${activeChar.value?.lastname}`;
    const user = {
        id: activeChar.value!.userId,
        name: ourName,
        color: stringToColor(ourName),
    };

    yjsSchema = getSchema(extensions);

    const yXml = ydoc.getXmlFragment('content');
    const { mapping } = initProseMirrorDoc(yXml, yjsSchema!);

    extensions.push(
        Collaboration.configure({
            document: ydoc,
            field: 'content',
            ySyncOptions: { mapping },
        }),
        CollaborationCaret.configure({
            provider: yjsProvider,
            user: user,
            // Skip rendering if it's your own cursor
            render: (user): HTMLElement => {
                if (user.id === yjsProvider.ydoc.clientID) {
                    // returns nothing -> no widget for your own cursor
                    return new HTMLElement();
                }
                // Otherwise build the "remote" cursor as normal:
                const cursor = document.createElement('span');
                cursor.classList.add('collaboration-cursor__caret');
                cursor.setAttribute('style', `border-color: ${user.color}`);

                const label = document.createElement('div');
                label.classList.add('collaboration-cursor__label');
                label.setAttribute('style', `background-color: ${user.color}`);
                label.insertBefore(document.createTextNode(user.name), null);

                cursor.insertBefore(label, null);
                return cursor;
            },
            // Same for text selections
            selectionRender: (user) => {
                if (user.id === yjsProvider.ydoc.clientID) {
                    return {};
                }
                return {
                    nodeName: 'span',
                    class: 'collaboration-cursor__selection',
                    style: `background-color: ${user.color}`,
                    'data-user': user.name,
                };
            },
        }),
    );

    const onSync = (synced: boolean) => {
        logger.info('Yjs sync event:', synced);
        if (!synced) {
            loading.value = true;
            return;
        }

        // Only set initial content if authoritative and Yjs doc is empty
        if (yjsProvider.isAuthoritative && modelValue.value) {
            seedDocument(yjsSchema!, modelValue.value);
        }

        useTimeoutFn(() => (loading.value = false), 250);
    };
    yjsProvider.on('sync', onSync);

    const onLoading = (state: boolean) => (loading.value = state);
    yjsProvider.on('loading', onLoading);

    onBeforeUnmount(() => {
        yjsProvider.off('sync', onSync);
        yjsProvider.off('loading', onLoading);
    });
    onMounted(() => yjsProvider.connect());
} else {
    extensions.push(UndoRedo);
}

function hasFileById(files: FileGrpc[] | undefined | null, id: number): boolean {
    if (!files || !id) return false;
    return files.some((f) => f.id === id);
}

const disabled = computed(() => props.disabled || loading.value);

let fileUploadHandler: undefined | ((files: File[]) => Promise<void>) = undefined;

const editor = useEditor({
    content: '',
    editorProps: {
        attributes: {
            class: 'prose prose-sm sm:prose-base lg:prose-lg m-5 focus:outline-hidden dark:prose-invert max-w-full break-words',
        },
    },
    editable: !disabled.value,
    extensions: [...extensions, ...markRaw(props.extensions)],
    onFocus: () => focusTablet(true),
    onBlur: () => focusTablet(false),
    onCreate: () => {
        if (props.filestoreService && props.filestoreNamespace && fileUploadHandler) {
            unref(editor)?.registerPlugin(imageUploadPlugin(unref(editor)!, fileUploadHandler));
        }
        logger.info('Editor created');
    },
    onUpdate: ({ editor }) => {
        modelValue.value = unref(editor)?.getJSON();
    },
});

if (props.filestoreService && props.filestoreNamespace && props.targetId) {
    const { resizeAndUpload } = useFileUploader(props.filestoreService, props.filestoreNamespace, props.targetId);

    async function handleFiles(fs: File[]): Promise<void> {
        for (const f of fs) {
            if (!f.type.startsWith('image/')) continue;

            if (files.value && files.value.length >= props.fileLimit) {
                logger.warn('File limit reached, cannot upload more files');
                notifications.add({
                    title: { key: 'components.partials.tiptap_editor.file_limit_reached.title', parameters: {} },
                    description: { key: 'components.partials.tiptap_editor.file_limit_reached.content', parameters: {} },
                    type: NotificationType.ERROR,
                });
                return;
            }

            try {
                const resp = await resizeAndUpload(f);

                unref(editor)!
                    .chain()
                    .focus()
                    .setEnhancedImage({ src: resp.url, alt: resp.file?.filePath, fileId: resp.file?.id })
                    .run();

                resp.file && emits('file-uploaded', resp.file);

                files.value.push(resp.file!);
            } catch (e) {
                logger.warn('Image resize failed, uploading original image', e);
            }
        }
    }
    fileUploadHandler = handleFiles;
}

// If collaboration is enabled, we don't set the content directly
// as it will be handled by the Yjs provider.
const stopWatch = watch(modelValue, (value) => {
    const editorJSON = unref(editor)?.getJSON();
    if (!editorJSON || !value) return;

    if (typeof value === 'string') {
        value = generateJSON(value, extensions) as JSONContent;
    }
    const isSame = isSameDoc(editorJSON, value, extensions);
    if (isSame) return;

    // If not authoritative, don't set the content
    if (props.enableCollab && ydoc && yjsProvider && !yjsProvider.isAuthoritative) return;

    if (value) {
        if (props.enableCollab && ydoc && yjsProvider) {
            seedDocument(yjsSchema!, value);
        } else {
            unref(editor)?.commands.setContent(value, { emitUpdate: true });
        }
    }

    if (props.enableCollab && ydoc && yjsProvider && yjsProvider.isAuthoritative) stopWatch();
});

const emojiItems: EditorEmojiMenuItem[] = gitHubEmojis.filter((emoji) => !emoji.name.startsWith('regional_indicator_'));

const contentRef = useTemplateRef('contentRef');

const openLink = ref(false);
const openLinkPopover = refDebounced(openLink, 200);
const anchor = ref({ x: 0, y: 0 });

const selectedLink = ref('');
const selectedAnchor = ref({ x: 0, y: 0 });
const reference = computed(() => ({
    getBoundingClientRect: () =>
        ({
            width: 0,
            height: 0,
            left: selectedAnchor.value.x,
            right: selectedAnchor.value.x,
            top: selectedAnchor.value.y,
            bottom: selectedAnchor.value.y,
            ...selectedAnchor.value,
        }) as DOMRect,
}));

function onClickContent(event: MouseEvent): void {
    let element: HTMLElement | null = event.target as HTMLElement;
    if (element.tagName.toLowerCase() !== 'a' && !element.hasAttribute('href')) {
        element = element.parentElement as HTMLElement;
        if (!element || (element.tagName.toLowerCase() !== 'a' && !element.hasAttribute('href'))) return;
    }
    event.preventDefault();

    selectedAnchor.value = { ...anchor.value };
    selectedLink.value = element.getAttribute('href') || '';
    openLink.value = true;

    element.addEventListener(
        'pointerleave',
        () => {
            openLink.value = false;
        },
        { once: true },
    );
}

watchOnce(contentRef, () => {
    if (!contentRef.value || !contentRef.value.$el) return;

    const element = contentRef.value.$el as HTMLDivElement;
    element.addEventListener('click', onClickContent);
});

const formErrors = inject<Ref<FormError[]> | null>(formErrorsInjectionKey, null);

const error = computed(() => formErrors?.value?.find((error) => error.name && error.name === props.name)?.message);

watch(
    editorSettings,
    () => {
        if (editorSettings.value.showInvisibleCharacters) {
            unref(editor)?.chain().focus().showInvisibleCharacters().run();
        } else {
            unref(editor)?.chain().focus().hideInvisibleCharacters().run();
        }

        if (
            editor.value &&
            typeof editor.value.options.editorProps.attributes === 'object' &&
            'class' in editor.value.options.editorProps.attributes
        ) {
            let c = editor.value.options.editorProps.attributes['class'] || '';
            if (editorSettings.value.focusMode) {
                if (c.includes('editor-focus')) return;

                c += ' editor-focus';
            } else {
                c = c.replace('editor-focus', '').trim();
            }

            unref(editor)!.setOptions({
                editorProps: {
                    attributes: {
                        ...unref(editor)!.options.editorProps.attributes,
                        class: c,
                    },
                },
            });
        }
    },
    { deep: true },
);

onMounted(() => {
    if (props.enableCollab) return;

    logger.info('Setting initial content for Tiptap editor (collab is disabled)');
    if (modelValue.value) unref(editor)?.commands.setContent(modelValue.value, { emitUpdate: false });
});

onBeforeUnmount(() => {
    if (contentRef.value?.$el) {
        const element = contentRef.value.$el as HTMLDivElement;
        element.removeEventListener('click', onClickContent);
    }
});

onBeforeRouteLeave(() => yjsProvider?.destroy());

defineExpose<{
    editor: typeof editor;
}>({
    editor,
});
</script>

<template>
    <UCard
        class="relative flex min-h-0 flex-1 flex-col overflow-y-hidden"
        :ui="{
            header: 'p-0 sm:px-2 sticky inset-x-0 top-0 z-[1] shrink-0 bg-neutral-100/75 p-0.5 backdrop-blur dark:bg-neutral-800/75',
            body: 'p-0 sm:p-0 overflow-y-auto flex-1 border-x border-neutral-100/75 dark:border-neutral-800/75',
            footer: 'p-0 sm:px-2 sticky inset-x-0 bottom-0 z-[1] flex w-full flex-none justify-between bg-neutral-100 px-1 text-center dark:bg-neutral-800',
        }"
        v-bind="$attrs"
    >
        <template v-if="editor && !hideToolbar" #header>
            <TiptapToolbar
                :editor="markRaw(editor)"
                :disabled="disabled"
                :disable-images="disableImages"
                :history-type="historyType"
                :file-limit="fileLimit"
                :file-upload-handler="fileUploadHandler"
                @update:content="($event) => editor?.commands.setContent($event ?? '', { emitUpdate: true })"
            >
                <template #toolbar>
                    <slot name="toolbar" :editor="editor" :disabled="disabled" />
                </template>
            </TiptapToolbar>
        </template>

        <!-- Nuxt UI Editor parts  that work with our wonderful "TiptapEditor" component -->
        <UEditorEmojiMenu :editor="editor" :items="emojiItems" />

        <UEditorDragHandle v-if="editor" :editor="editor" />

        <UPopover
            :open="openLinkPopover"
            :reference="reference"
            :content="{ side: 'top', sideOffset: 16, updatePositionStrategy: 'always' }"
        >
            <TiptapEditorContent
                ref="contentRef"
                class="min-h-0 w-full max-w-full min-w-0 flex-1 flex-auto overflow-y-auto px-6 py-2 sm:px-10"
                :class="[
                    wrapperClass,
                    'hover:prose-a:text-blue-500',
                    'dark:hover:prose-a:text-blue-300',
                    'prose-headings:my-0.5',
                    'prose-lead:my-0.5',
                    'prose-h1:my-0.5',
                    'prose-h2:my-0.5',
                    'prose-h3:my-0.5',
                    'prose-h4:my-0.5',
                    'prose-p:my-0.5',
                    'prose-a:my-0.5',
                    'prose-blockquote:my-0.5',
                    'prose-figure:my-0.5',
                    'prose-figcaption:my-0.5',
                    'prose-strong:my-0.5',
                    'prose-em:my-0.5',
                    'prose-kbd:my-0.5',
                    'prose-code:my-0.5',
                    'prose-pre:my-0.5',
                    'prose-ol:my-0.5',
                    'prose-ul:my-0.5',
                    'prose-li:my-0.5',
                    'prose-table:my-0.5',
                    'prose-thead:my-0.5',
                    'prose-tr:my-0.5',
                    'prose-th:my-0.5',
                    'prose-td:my-0.5',
                    'prose-img:my-0.5',
                    'prose-video:my-0.5',
                    'prose-hr:my-0.5',
                ]"
                :editor="editor"
                @pointerleave="openLink = false"
                @pointermove="
                    (ev: PointerEvent) => {
                        anchor.x = ev.clientX;
                        anchor.y = ev.clientY;
                    }
                "
            />

            <template #content>
                <div class="p-2" @pointerenter="openLink = true" @pointerleave="openLink = false">
                    <UButton variant="link" :to="selectedLink" external target="_blank" :label="selectedLink" />
                </div>
            </template>
        </UPopover>

        <template v-if="editor" #footer>
            <div class="flex w-full flex-1 flex-col gap-1">
                <div v-if="error" class="mb-2 flex items-start">
                    <div v-if="typeof error === 'string'" :id="`${$props.name}-error`" class="text-error">{{ error }}</div>
                </div>

                <div class="flex flex-1 flex-row flex-wrap justify-between gap-2">
                    <div class="inline-flex">
                        <template v-if="$slots.footer">
                            <slot name="footer" :saving="saving" />
                        </template>
                        <div v-else-if="saving" class="inline-flex items-center gap-1">
                            <UIcon class="h-4 w-4 animate-spin" name="i-mdi-content-save" />
                            <span>{{ $t('common.save', 2) }}...</span>
                        </div>

                        <div v-if="loading" class="inline-flex items-center gap-1">
                            <UIcon class="size-5 animate-spin" name="i-mdi-refresh" />
                            {{ $t('common.loading') }}
                        </div>
                    </div>

                    <YJSUserPopover v-if="enableCollab && targetId" />

                    <div class="flex items-center justify-end gap-1">
                        {{ unref(editor).storage.characterCount.characters()
                        }}<template v-if="limit && limit > 0"> / {{ limit }}</template>
                        {{ $t('common.chars', unref(editor).storage.characterCount.characters()) }}
                        |
                        {{ unref(editor).storage.characterCount.words() }}
                        {{ $t('common.word', unref(editor).storage.characterCount.words()) }}
                    </div>
                </div>
            </div>
        </template>
    </UCard>
</template>

<style lang="scss">
.ProseMirror {
    > * {
        margin-left: 0.75rem;
    }

    .ProseMirror-widget * {
        margin-top: auto;
    }

    ul,
    ol {
        padding: 0 1rem;
    }
}

/* Tiptap Editor Drag Handle Style */
.tiptap-drag-handle {
    &::after {
        margin-right: 0.1rem;
        width: 1rem;
        height: 1.25rem;
        font-weight: 700;
        background: #0d0d0d10;
        color: #0d0d0d50;
        border-radius: 0.25rem;
    }
}
</style>
