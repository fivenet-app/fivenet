<script lang="ts" setup>
import type { ClientStreamingCall, RpcOptions } from '@protobuf-ts/runtime-rpc';
import type { Extensions, Range } from '@tiptap/core';
import { Blockquote } from '@tiptap/extension-blockquote';
import { Bold } from '@tiptap/extension-bold';
import { BulletList } from '@tiptap/extension-bullet-list';
import CharacterCount from '@tiptap/extension-character-count';
import { Code } from '@tiptap/extension-code';
import { CodeBlock } from '@tiptap/extension-code-block';
import Collaboration from '@tiptap/extension-collaboration';
import CollaborationCursor from '@tiptap/extension-collaboration-cursor';
import { Color } from '@tiptap/extension-color';
import { Document } from '@tiptap/extension-document';
import { Dropcursor } from '@tiptap/extension-dropcursor';
import FontFamily from '@tiptap/extension-font-family';
import { Gapcursor } from '@tiptap/extension-gapcursor';
import { HardBreak } from '@tiptap/extension-hard-break';
import { Heading } from '@tiptap/extension-heading';
import Highlight from '@tiptap/extension-highlight';
import { History } from '@tiptap/extension-history';
import { HorizontalRule } from '@tiptap/extension-horizontal-rule';
import { Italic } from '@tiptap/extension-italic';
import Link from '@tiptap/extension-link';
import { ListItem } from '@tiptap/extension-list-item';
import ListKeymap from '@tiptap/extension-list-keymap';
import { OrderedList } from '@tiptap/extension-ordered-list';
import { Paragraph } from '@tiptap/extension-paragraph';
import Placeholder from '@tiptap/extension-placeholder';
import { Strike } from '@tiptap/extension-strike';
import Subscript from '@tiptap/extension-subscript';
import Superscript from '@tiptap/extension-superscript';
import Table from '@tiptap/extension-table';
import TableCell from '@tiptap/extension-table-cell';
import TableHeader from '@tiptap/extension-table-header';
import TableRow from '@tiptap/extension-table-row';
import TaskItem from '@tiptap/extension-task-item';
import TaskList from '@tiptap/extension-task-list';
import { Text } from '@tiptap/extension-text';
import TextAlign from '@tiptap/extension-text-align';
import TextStyle from '@tiptap/extension-text-style';
import Underline from '@tiptap/extension-underline';
import FontSize from 'tiptap-extension-font-size';
// @ts-expect-error doesn't have types
import UniqueId from 'tiptap-unique-id';
import type * as Y from 'yjs';
import { CheckboxStandalone } from '~/composables/tiptap/extensions/CheckboxStandalone';
import { DeleteImageTracker } from '~/composables/tiptap/extensions/DeleteImageTracker';
import { EnhancedImageResize } from '~/composables/tiptap/extensions/EnhancedImageResize';
import { imageUploadPlugin } from '~/composables/tiptap/extensions/ImageUploadPlugin';
import SearchAndReplace from '~/composables/tiptap/extensions/SearchAndReplace';
import type { UploadNamespaces } from '~/composables/useFileUploader';
import type GrpcProvider from '~/composables/yjs/yjs';
import { fontColors, highlightColors } from '~/types/editor';
import type { Content, Version } from '~/types/history';
import type { File as FileGrpc } from '~~/gen/ts/resources/file/file';
import type { UploadPacket, UploadResponse } from '~~/gen/ts/resources/file/filestore';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import FileListModal from './FileListModal.vue';
import TiptapEditorImagePopover from './TiptapEditorImagePopover.vue';
import TiptapEditorSourceCodeModal from './TiptapEditorSourceCodeModal.vue';
import VersionHistoryModal from './VersionHistoryModal.vue';
import YJSUserPopover from './YJSUserPopover.vue';

const props = withDefaults(
    defineProps<{
        wrapperClass?: string;
        limit?: number;
        fileLimit?: number;
        disabled?: boolean;
        placeholder?: string;
        hideToolbar?: boolean;
        disableImages?: boolean;
        historyType?: string;
        disableCollab?: boolean;

        extensions?: Extensions;

        saving?: boolean;

        targetId?: number;
        filestoreNamespace?: UploadNamespaces;
        filestoreService?: (options?: RpcOptions) => ClientStreamingCall<UploadPacket, UploadResponse>;
    }>(),
    {
        wrapperClass: '',
        limit: undefined,
        fileLimit: 10,
        disabled: false,
        placeholder: undefined,
        hideToolbar: false,
        disableImages: false,
        historyType: undefined,
        disableCollab: false,

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

const { t } = useI18n();

const logger = useLogger('📄 YJS');

const { activeChar } = useAuth();

const notifications = useNotificatorStore();

const modal = useModal();

const modelValue = defineModel<string>({ required: true });
const files = defineModel<FileGrpc[]>('files', { default: () => [] });

const extensions: Extensions = [
    UniqueId.configure({
        attributeName: 'id',
        types: ['heading'],
        createId: () => window.crypto.randomUUID(),
    }),
    // Basics
    Blockquote,
    Bold,
    BulletList,
    Code,
    CodeBlock,
    Color,
    Document,
    Dropcursor,
    FontFamily,
    FontSize,
    Gapcursor,
    HardBreak,
    Heading,
    Highlight.configure({
        multicolor: true,
    }),
    HorizontalRule,
    Italic,
    Link.configure({
        openOnClick: false,
        defaultProtocol: 'https',
        HTMLAttributes: {
            target: null,
        },
    }),
    ListItem,
    ListKeymap,
    OrderedList,
    Paragraph,
    Strike,
    Subscript,
    Superscript,
    Text,
    TextAlign.configure({
        types: ['heading', 'paragraph', 'image'],
    }),
    TextStyle,
    Underline,
    // Table
    Table.configure({
        resizable: true,
        allowTableNodeSelection: true,
        HTMLAttributes: {
            class: 'border border-collapse border-solid border-gray-500',
        },
    }),
    TableRow,
    TableHeader.configure({
        HTMLAttributes: {
            class: 'border border-solid border-gray-600 bg-gray-100 dark:bg-gray-800',
        },
    }),
    TableCell.configure({
        HTMLAttributes: {
            class: 'border border-solid border-gray-500',
        },
    }),
    // Misc
    SearchAndReplace,
    TaskList,
    TaskItem.configure({
        nested: true,
    }),
    CheckboxStandalone,
    CharacterCount.configure({
        limit: props.limit,
    }),
    Placeholder.configure({
        placeholder: props.placeholder ?? '',
    }),
];

const ydoc = inject<Y.Doc | undefined>('yjsDoc', undefined);
const yjsProvider = inject<GrpcProvider | undefined>('yjsProvider', undefined);

const loading = ref(ydoc && yjsProvider && !props.disableCollab);

if (ydoc && yjsProvider && !props.disableCollab) {
    const ourName = `${activeChar.value?.firstname} ${activeChar.value?.lastname}`;
    const user = {
        id: activeChar.value!.userId,
        name: ourName,
        color: stringToColor(ourName),
    };

    const onSync = (synced: boolean) => {
        logger.info('Yjs sync event:', synced);
        if (synced === false) {
            loading.value = true;
            return;
        }

        loading.value = false;
    };
    yjsProvider.on('sync', onSync);
    onBeforeUnmount(() => yjsProvider.off('sync', onSync));
    onMounted(() => yjsProvider.connect());

    extensions.push(
        Collaboration.configure({
            document: ydoc,
            field: 'content',
        }),
        CollaborationCursor.configure({
            provider: yjsProvider,
            user: user,
            // Skip rendering if it's your own cursor
            render: (user): HTMLElement => {
                if (user.id === yjsProvider.ydoc.clientID) {
                    // returns nothing → no widget for your own cursor
                    return new HTMLElement();
                }
                // Otherwise build the “remote” cursor as normal:
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
} else {
    extensions.push(History);
}

function hasFileById(files: FileGrpc[] | undefined | null, id: number): boolean {
    if (!files || !id) return false;
    return files.some((f) => f.id === id);
}

if (!props.disableImages) {
    extensions.push(
        EnhancedImageResize.configure({
            inline: false,
            allowBase64: true,
        }),
        DeleteImageTracker.configure({
            onRemoved: (ids) =>
                ids.forEach((id) => {
                    if (hasFileById(files.value, id)) {
                        const idx = files.value.findIndex((f) => f.id === id);
                        if (idx !== -1) files.value.splice(idx, 1);
                    }
                }),
        }),
    );
}

const disabled = computed(() => props.disabled || loading.value);

let fileUploadHandler: undefined | ((files: File[]) => Promise<void>) = undefined;

const editor = useEditor({
    content: '',
    editable: !disabled.value,
    extensions: [...extensions, ...props.extensions],
    onFocus: () => focusTablet(true),
    onBlur: () => focusTablet(false),
    onUpdate: () => (modelValue.value = unref(editor)?.getHTML() ?? ''),
    editorProps: {
        attributes: {
            class: 'prose prose-sm sm:prose-base lg:prose-lg m-5 focus:outline-none dark:prose-invert max-w-full break-words',
        },
    },
    onCreate: () => {
        if (props.filestoreService && props.filestoreNamespace && fileUploadHandler) {
            unref(editor)?.registerPlugin(imageUploadPlugin(unref(editor)!, fileUploadHandler));
        }
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
                    title: { key: 'components.partials.TiptapEditor.file_limit_reached.title', parameters: {} },
                    description: { key: 'components.partials.TiptapEditor.file_limit_reached.content', parameters: {} },
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

const fonts = [
    {
        label: t('common.default'),
        value: 'DM Sans',
    },
    {
        label: 'Arial',
        value: 'arial, helvetica, sans-serif',
    },
    {
        label: 'Serif',
        value: 'serif',
    },
    {
        label: 'Times New Romain',
        value: 'times new roman, times, serif',
    },
    {
        label: 'Comic Sans',
        value: 'Comic Sans MS, Comic Sans',
    },
    {
        label: 'Monospace',
        value: 'monospace',
    },
];

// If collaboration is enabled, we don't set the content directly
// as it will be handled by the Yjs provider.
const stopWatch = watch(modelValue, (value) => {
    const isSame = unref(editor)?.getHTML() === value;
    // JSON
    // const isSame = JSON.stringify(this.editor.getJSON()) === JSON.stringify(value);

    if (isSame) return;

    // If not authoritative, don't set the content
    if (!props.disableCollab && ydoc && yjsProvider && !yjsProvider.isAuthoritative) return;

    unref(editor)?.commands.setContent(value, true);
    if (!props.disableCollab && ydoc && yjsProvider && yjsProvider.isAuthoritative) {
        stopWatch();
    }
});

watch(disabled, () => unref(editor)?.setEditable(!disabled.value));

const linkState = reactive({
    url: '',
});

function setLink(data: typeof linkState): void {
    const previousUrl = unref(editor)?.getAttributes('link').href;
    const url = data.url.trim() !== '' ? data.url.trim() : previousUrl;

    // Empty URL
    if (url === '') {
        unref(editor)?.chain().focus().extendMarkRange('link').unsetLink().run();
        return;
    }

    // Update link
    unref(editor)
        ?.chain()
        .focus()
        .extendMarkRange('link')
        .setLink({
            href: url,
        })
        .run();
}

const selectedFont = ref<(typeof fonts)[0]>(fonts[0]!);
watch(selectedFont, () => unref(editor)?.chain().focus().setFontFamily(selectedFont.value.value).run());

const selectedFontColor = ref<string | undefined>(undefined);
watch(selectedFontColor, () =>
    selectedFontColor.value
        ? unref(editor)?.chain().focus().setColor(selectedFontColor.value).run()
        : unref(editor)?.chain().focus().unsetColor().run(),
);

const selectedHighlightColor = ref<(typeof highlightColors)[0]>(highlightColors[0]!);
watch(selectedHighlightColor, () =>
    unref(editor)?.chain().focus().toggleHighlight({ color: selectedHighlightColor.value.value }).run(),
);

const tableCreation = reactive({
    rows: 2,
    cols: 3,
    withHeaderRow: true,
});

function createTable(): void {
    unref(editor)
        ?.chain()
        .focus()
        .insertTable({
            rows: tableCreation.rows,
            cols: tableCreation.cols,
            withHeaderRow: tableCreation.withHeaderRow,
        })
        .run();
}

const searchAndReplace = reactive<{
    search: string;
    replace: string;
    caseSensitive: boolean;
}>({
    search: '',
    replace: '',
    caseSensitive: false,
});

const updateSearchReplace = (clearIndex: boolean = false) => {
    if (!editor.value) {
        return;
    }

    if (clearIndex) {
        editor.value.commands.resetIndex();
    }

    unref(editor)?.commands.setSearchTerm(searchAndReplace.search);
    unref(editor)?.commands.setReplaceTerm(searchAndReplace.replace ?? '');
    unref(editor)?.commands.setCaseSensitive(searchAndReplace.caseSensitive);
};

const goToSelection = () => {
    if (!editor.value) {
        return;
    }

    const { results, resultIndex } = editor.value.storage.searchAndReplace;
    const position: Range = results[resultIndex];

    if (!position) {
        return;
    }

    unref(editor)?.commands.setTextSelection(position);

    const { node } = editor.value.view.domAtPos(editor.value.state.selection.anchor);
    node instanceof HTMLElement && node.scrollIntoView({ behavior: 'smooth', block: 'center' });
};

// Search And Replace Modal
watch(
    () => searchAndReplace.search.trim(),
    (val, oldVal) => {
        if (!val) clear();
        if (val !== oldVal) updateSearchReplace(true);
    },
);

watch(
    () => searchAndReplace.replace.trim(),
    (val, oldVal) => (val === oldVal ? null : updateSearchReplace()),
);

watch(
    () => searchAndReplace.caseSensitive,
    (val, oldVal) => (val === oldVal ? null : updateSearchReplace(true)),
);

const replace = () => {
    editor.value?.commands.replace();
    goToSelection();
};

const next = () => {
    editor.value?.commands.nextSearchResult();
    goToSelection();
};

const previous = () => {
    editor.value?.commands.previousSearchResult();
    goToSelection();
};

const clear = () => {
    searchAndReplace.search = searchAndReplace.replace = '';
    editor.value?.commands.resetIndex();
};

const replaceAll = () => editor.value?.commands.replaceAll();

const contentRef = useTemplateRef('contentRef');
watch(contentRef, () => {
    if (!contentRef.value || !contentRef.value.$el) {
        return;
    }

    const element = contentRef.value.$el as HTMLDivElement;
    element.addEventListener('click', (event: MouseEvent) => {
        const element = event.target as HTMLElement;
        if (element.tagName.toLowerCase() !== 'a' && !element.hasAttribute('href')) {
            return;
        }

        event.preventDefault();
    });
});

function applyVersion(version: Version<unknown>): void {
    const v = version as Version<Content>;
    modelValue.value = v.content.content;
    files.value = v.content.files;

    notifications.add({
        title: { key: 'notifications.action_successful.title', parameters: {} },
        description: { key: 'notifications.action_successful.content', parameters: {} },
        type: NotificationType.SUCCESS,
    });
}

onMounted(() => {
    if (ydoc === undefined) {
        unref(editor)?.commands.setContent(modelValue.value);
    }
});

onBeforeUnmount(() => unref(editor)?.destroy());
</script>

<template>
    <div class="relative flex flex-col rounded-none border border-gray-100 dark:border-gray-800">
        <div v-if="editor && !hideToolbar" class="shrink-0 bg-gray-100 p-0.5 dark:bg-gray-800">
            <div class="flex snap-x flex-wrap gap-1">
                <UButtonGroup>
                    <UTooltip :text="$t('components.partials.TiptapEditor.bold')" :popper="{ placement: 'top' }">
                        <UButton
                            :class="{ 'is-active': editor.isActive('bold') }"
                            :disabled="!editor.can().chain().focus().toggleBold().run() || disabled"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-format-bold"
                            @click="editor.chain().focus().toggleBold().run()"
                        />
                    </UTooltip>
                    <UTooltip :text="$t('components.partials.TiptapEditor.italic')" :popper="{ placement: 'top' }">
                        <UButton
                            :class="{ 'is-active': editor.isActive('italic') }"
                            :disabled="!editor.can().chain().focus().toggleItalic().run() || disabled"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-format-italic"
                            @click="editor.chain().focus().toggleItalic().run()"
                        />
                    </UTooltip>
                    <UTooltip :text="$t('components.partials.TiptapEditor.underline')" :popper="{ placement: 'top' }">
                        <UButton
                            :class="{ 'is-active': editor.isActive('underline') }"
                            :disabled="disabled"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-format-underline"
                            @click="editor.chain().focus().toggleUnderline().run()"
                        />
                    </UTooltip>
                    <UTooltip :text="$t('components.partials.TiptapEditor.strike')" :popper="{ placement: 'top' }">
                        <UButton
                            :class="{ 'is-active': editor.isActive('strike') }"
                            :disabled="!editor.can().chain().focus().toggleStrike().run() || disabled"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-format-strikethrough"
                            @click="editor.chain().focus().toggleStrike().run()"
                        />
                    </UTooltip>
                    <UTooltip :text="$t('components.partials.TiptapEditor.clear')" :popper="{ placement: 'top' }">
                        <UButton
                            :disabled="disabled"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-format-clear"
                            @click="editor.chain().focus().unsetAllMarks().run()"
                        />
                    </UTooltip>
                    <UTooltip :text="$t('components.partials.TiptapEditor.superscript')" :popper="{ placement: 'top' }">
                        <UButton
                            :class="{ 'is-active': editor.isActive('superscript') }"
                            :disabled="disabled"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-format-superscript"
                            @click="editor.chain().focus().toggleSuperscript().run()"
                        />
                    </UTooltip>
                    <UTooltip :text="$t('components.partials.TiptapEditor.subscript')" :popper="{ placement: 'top' }">
                        <UButton
                            :class="{ 'is-active': editor.isActive('subscript') }"
                            :disabled="disabled"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-format-subscript"
                            @click="editor.chain().focus().toggleSubscript().run()"
                        />
                    </UTooltip>
                    <UTooltip :text="$t('components.partials.TiptapEditor.code')" :popper="{ placement: 'top' }">
                        <UButton
                            :class="{ 'is-active': editor.isActive('code') }"
                            :disabled="!editor.can().chain().focus().toggleCode().run() || disabled"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-code-braces"
                            @click="editor.chain().focus().toggleCode().run()"
                        />
                    </UTooltip>
                </UButtonGroup>

                <UDivider orientation="vertical" :ui="{ border: { base: 'border-gray-200 dark:border-gray-700' } }" />

                <!-- Text Align -->
                <UButtonGroup>
                    <UTooltip :text="$t('components.partials.TiptapEditor.align_left')" :popper="{ placement: 'top' }">
                        <UButton
                            :class="{ 'is-active': editor.isActive({ textAlign: 'left' }) }"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-format-align-left"
                            :disabled="disabled"
                            @click="editor.chain().focus().setTextAlign('left').run()"
                        />
                    </UTooltip>
                    <UTooltip :text="$t('components.partials.TiptapEditor.align_center')" :popper="{ placement: 'top' }">
                        <UButton
                            :class="{ 'is-active': editor.isActive({ textAlign: 'center' }) }"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-format-align-center"
                            :disabled="disabled"
                            @click="editor.chain().focus().setTextAlign('center').run()"
                        />
                    </UTooltip>
                    <UTooltip :text="$t('components.partials.TiptapEditor.align_right')" :popper="{ placement: 'top' }">
                        <UButton
                            :class="{ 'is-active': editor.isActive({ textAlign: 'right' }) }"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-format-align-right"
                            :disabled="disabled"
                            @click="editor.chain().focus().setTextAlign('right').run()"
                        />
                    </UTooltip>
                    <UTooltip :text="$t('components.partials.TiptapEditor.align_justify')" :popper="{ placement: 'top' }">
                        <UButton
                            :class="{ 'is-active': editor.isActive({ textAlign: 'justify' }) }"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-format-align-justify"
                            :disabled="disabled"
                            @click="editor.chain().focus().setTextAlign('justify').run()"
                        />
                    </UTooltip>
                </UButtonGroup>

                <UDivider orientation="vertical" :ui="{ border: { base: 'border-gray-200 dark:border-gray-700' } }" />

                <!-- Font Family -->
                <UTooltip :text="$t('components.partials.TiptapEditor.font_family')" :popper="{ placement: 'top' }">
                    <UInputMenu
                        v-model="selectedFont"
                        class="max-w-40"
                        name="selectedFont"
                        option-attribute="label"
                        :search-attributes="['label']"
                        :options="fonts"
                        :placeholder="$t('common.font', 1)"
                        search-lazy
                        :search-placeholder="$t('common.search_field')"
                        :disabled="disabled"
                    >
                        <template #label>
                            <span class="truncate" :style="{ fontFamily: selectedFont.value }">{{ selectedFont.label }}</span>
                        </template>

                        <template #option="{ option }">
                            <span class="truncate" :style="{ fontFamily: option.value }">{{ option.label }}</span>
                        </template>

                        <template #option-empty="{ query: search }">
                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                        </template>

                        <template #empty> {{ $t('common.not_found', [$t('common.job', 2)]) }} </template>>
                    </UInputMenu>
                </UTooltip>

                <UButtonGroup>
                    <UPopover>
                        <UTooltip :text="$t('components.partials.TiptapEditor.font_color')" :popper="{ placement: 'top' }">
                            <UButton
                                :class="{ 'is-active': editor.isActive('color', { color: selectedFontColor }) }"
                                color="white"
                                variant="ghost"
                                :style="{ color: selectedFontColor }"
                                icon="i-mdi-format-color-text"
                                :disabled="disabled"
                            />
                        </UTooltip>

                        <template #panel="{ close }">
                            <div class="inline-flex flex-col gap-1 p-4">
                                <UButton
                                    class="rounded-md"
                                    color="white"
                                    variant="outline"
                                    icon="i-mdi-water-off"
                                    :label="$t('common.default')"
                                    :disabled="disabled"
                                    @click="
                                        editor.chain().focus().unsetColor().run();
                                        close();
                                    "
                                />

                                <div v-for="(colors, idx) in fontColors" :key="idx">
                                    <div class="grid grid-cols-10 gap-0.5">
                                        <UButton
                                            v-for="(col, cIdx) in colors"
                                            :key="cIdx"
                                            class="size-6 rounded-none border-0"
                                            :style="{ backgroundColor: col }"
                                            :disabled="disabled"
                                            @click="selectedFontColor = col"
                                        />
                                    </div>
                                </div>
                            </div>
                        </template>
                    </UPopover>

                    <!-- Paragraph + Headers -->
                    <UTooltip :text="$t('components.partials.TiptapEditor.paragraph')" :popper="{ placement: 'top' }">
                        <UButton
                            :class="{ 'is-active': editor.isActive('paragraph') }"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-format-paragraph"
                            :disabled="disabled"
                            @click="editor.chain().focus().setParagraph().run()"
                        />
                    </UTooltip>

                    <UTooltip :text="$t('components.partials.TiptapEditor.header_1')" :popper="{ placement: 'top' }">
                        <UButton
                            :class="{ 'is-active': editor.isActive('heading', { level: 1 }) }"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-format-header-1"
                            :disabled="disabled"
                            @click="editor.chain().focus().toggleHeading({ level: 1 }).run()"
                        />
                    </UTooltip>
                    <UTooltip :text="$t('components.partials.TiptapEditor.header_2')" :popper="{ placement: 'top' }">
                        <UButton
                            :class="{ 'is-active': editor.isActive('heading', { level: 2 }) }"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-format-header-2"
                            :disabled="disabled"
                            @click="editor.chain().focus().toggleHeading({ level: 2 }).run()"
                        />
                    </UTooltip>
                    <UTooltip :text="$t('components.partials.TiptapEditor.header_3')" :popper="{ placement: 'top' }">
                        <UButton
                            :class="{ 'is-active': editor.isActive('heading', { level: 3 }) }"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-format-header-3"
                            :disabled="disabled"
                            @click="editor.chain().focus().toggleHeading({ level: 3 }).run()"
                        />
                    </UTooltip>
                    <UTooltip :text="$t('components.partials.TiptapEditor.header_4')" :popper="{ placement: 'top' }">
                        <UButton
                            :class="{ 'is-active': editor.isActive('heading', { level: 4 }) }"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-format-header-4"
                            :disabled="disabled"
                            @click="editor.chain().focus().toggleHeading({ level: 4 }).run()"
                        />
                    </UTooltip>
                    <UTooltip :text="$t('components.partials.TiptapEditor.header_5')" :popper="{ placement: 'top' }">
                        <UButton
                            :class="{ 'is-active': editor.isActive('heading', { level: 5 }) }"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-format-header-5"
                            :disabled="disabled"
                            @click="editor.chain().focus().toggleHeading({ level: 5 }).run()"
                        />
                    </UTooltip>
                    <UTooltip :text="$t('components.partials.TiptapEditor.header_6')" :popper="{ placement: 'top' }">
                        <UButton
                            :class="{ 'is-active': editor.isActive('heading', { level: 6 }) }"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-format-header-6"
                            :disabled="disabled"
                            @click="editor.chain().focus().toggleHeading({ level: 6 }).run()"
                        />
                    </UTooltip>
                </UButtonGroup>
            </div>

            <div class="flex snap-x flex-wrap gap-1">
                <UButtonGroup>
                    <UTooltip :text="$t('components.partials.TiptapEditor.highlight')" :popper="{ placement: 'top' }">
                        <UButton
                            :class="{ 'is-active': editor.isActive('highlight') }"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-format-color-highlight"
                            :disabled="disabled"
                            @click="editor.chain().focus().toggleHighlight().run()"
                        />
                    </UTooltip>

                    <UPopover>
                        <UTooltip :text="$t('components.partials.TiptapEditor.highlight_color')" :popper="{ placement: 'top' }">
                            <UButton
                                :class="{ 'is-active': editor.isActive('highlight', { color: selectedHighlightColor.value }) }"
                                color="white"
                                variant="ghost"
                                :style="{ color: selectedHighlightColor.value }"
                                icon="i-mdi-format-color-fill"
                                :disabled="disabled"
                            />
                        </UTooltip>

                        <template #panel="{ close }">
                            <div class="inline-flex flex-col gap-1 p-4">
                                <UButton
                                    class="rounded-md"
                                    color="white"
                                    variant="outline"
                                    icon="i-mdi-water-off"
                                    :label="$t('common.reset')"
                                    :disabled="disabled"
                                    @click="
                                        editor.chain().focus().unsetHighlight().run();
                                        close();
                                    "
                                />

                                <div class="grid grid-cols-6 gap-0.5">
                                    <UButton
                                        v-for="(col, idx) in highlightColors"
                                        :key="idx"
                                        class="size-6 rounded-none border-0"
                                        :style="{ backgroundColor: col.value }"
                                        :disabled="disabled"
                                        @click="selectedHighlightColor = col"
                                    />
                                </div>
                            </div>
                        </template>
                    </UPopover>
                </UButtonGroup>

                <UDivider orientation="vertical" :ui="{ border: { base: 'border-gray-200 dark:border-gray-700' } }" />

                <UButtonGroup>
                    <UTooltip :text="$t('components.partials.TiptapEditor.bullet_list')" :popper="{ placement: 'top' }">
                        <UButton
                            :class="{ 'is-active': editor.isActive('bulletList') }"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-format-list-bulleted"
                            :disabled="disabled"
                            @click="editor.chain().focus().toggleBulletList().run()"
                        />
                    </UTooltip>
                    <UTooltip :text="$t('components.partials.TiptapEditor.ordered_list')" :popper="{ placement: 'top' }">
                        <UButton
                            :class="{ 'is-active': editor.isActive('orderedList') }"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-format-list-numbered"
                            :disabled="disabled"
                            @click="editor.chain().focus().toggleOrderedList().run()"
                        />
                    </UTooltip>
                    <UTooltip :text="$t('components.partials.TiptapEditor.task_list')" :popper="{ placement: 'top' }">
                        <UButton
                            :class="{ 'is-active': editor.isActive('taskList') }"
                            icon="i-mdi-format-list-checks"
                            color="white"
                            variant="ghost"
                            :disabled="disabled"
                            @click="editor.chain().focus().toggleTaskList().run()"
                        />
                    </UTooltip>

                    <UTooltip :text="$t('components.partials.TiptapEditor.checkbox')" :popper="{ placement: 'top' }">
                        <UButton
                            :class="{ 'is-active': editor.isActive('checkboxStandalone') }"
                            icon="i-mdi-checkbox-marked-outline"
                            color="white"
                            variant="ghost"
                            :disabled="disabled"
                            @click="editor.chain().focus().addCheckboxStandalone().run()"
                        />
                    </UTooltip>
                </UButtonGroup>

                <UDivider orientation="vertical" :ui="{ border: { base: 'border-gray-200 dark:border-gray-700' } }" />

                <TiptapEditorImagePopover
                    v-if="!disableImages"
                    :editor="editor"
                    :file-limit="fileLimit"
                    :disabled="disabled"
                    :upload-handler="fileUploadHandler"
                    @open-file-list="
                        modal.open(FileListModal, {
                            editor: editor,
                            files: files,
                        })
                    "
                />

                <UPopover>
                    <UTooltip :text="$t('components.partials.TiptapEditor.table')" :popper="{ placement: 'top' }">
                        <UButton
                            :class="{ 'is-active': editor.isActive('table') }"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-table"
                            :disabled="disabled"
                        />
                    </UTooltip>

                    <template #panel>
                        <div class="p-4">
                            <UForm :state="{}" @submit="createTable">
                                <UFormGroup :label="$t('common.rows')">
                                    <UInput v-model="tableCreation.rows" type="text" :disabled="disabled" />
                                </UFormGroup>

                                <UFormGroup :label="$t('common.cols')">
                                    <UInput v-model="tableCreation.cols" type="text" :disabled="disabled" />
                                </UFormGroup>

                                <UFormGroup :label="$t('common.with_header_row')">
                                    <UToggle v-model="tableCreation.withHeaderRow" type="text" :disabled="disabled" />
                                </UFormGroup>

                                <UFormGroup>
                                    <UButton type="submit" :label="$t('common.create')" :disabled="disabled" />
                                </UFormGroup>
                            </UForm>
                        </div>
                    </template>
                </UPopover>

                <UPopover>
                    <UTooltip :text="$t('components.partials.TiptapEditor.link')" :popper="{ placement: 'top' }">
                        <UButton
                            :class="{ 'is-active': editor.isActive('link') }"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-link"
                            :disabled="disabled"
                        />
                    </UTooltip>

                    <template #panel="{ close }">
                        <div class="p-4">
                            <UForm :state="linkState" @submit="($event) => setLink($event.data)">
                                <UFormGroup :label="$t('common.url')">
                                    <UInput v-model="linkState.url" type="text" :disabled="disabled" />
                                </UFormGroup>

                                <slot name="linkModal" :editor="editor" :state="linkState" />

                                <UButtonGroup class="mt-2 w-full">
                                    <UButton
                                        class="flex-1"
                                        type="submit"
                                        icon="i-mdi-link"
                                        :label="$t('common.link')"
                                        :disabled="disabled"
                                    />

                                    <UButton
                                        :disabled="!editor.isActive('link') || disabled"
                                        color="error"
                                        variant="outline"
                                        icon="i-mdi-link-off"
                                        :label="$t('common.unlink')"
                                        @click="
                                            close();
                                            editor.chain().focus().unsetLink().run();
                                            linkState.url = '';
                                        "
                                    />
                                </UButtonGroup>
                            </UForm>
                        </div>
                    </template>
                </UPopover>

                <UButtonGroup>
                    <UTooltip :text="$t('components.partials.TiptapEditor.code_block')" :popper="{ placement: 'top' }">
                        <UButton
                            :class="{ 'is-active': editor.isActive('codeBlock') }"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-code-block-braces"
                            :disabled="disabled"
                            @click="editor.chain().focus().toggleCodeBlock().run()"
                        />
                    </UTooltip>
                    <UTooltip :text="$t('components.partials.TiptapEditor.block_quote')" :popper="{ placement: 'top' }">
                        <UButton
                            :class="{ 'is-active': editor.isActive('blockquote') }"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-format-quote-open"
                            :disabled="disabled"
                            @click="editor.chain().focus().toggleBlockquote().run()"
                        />
                    </UTooltip>
                    <UTooltip :text="$t('components.partials.TiptapEditor.horizontal_rule')" :popper="{ placement: 'top' }">
                        <UButton
                            color="white"
                            variant="ghost"
                            icon="i-mdi-minus"
                            :disabled="disabled"
                            @click="editor.chain().focus().setHorizontalRule().run()"
                        />
                    </UTooltip>
                    <!--
                    <UButton
                    color="white"
                    variant="ghost"
                    icon="i-mdi-format-page-break"
                    @click="editor.chain().focus().setHardBreak().run()"
                    :disabled="disabled"
                    />
                    -->
                </UButtonGroup>

                <div class="flex-1"></div>

                <slot name="toolbar" :editor="editor" :disabled="disabled" />

                <UDivider orientation="vertical" :ui="{ border: { base: 'border-gray-200 dark:border-gray-700' } }" />

                <UPopover>
                    <UTooltip :text="$t('components.partials.TiptapEditor.search_and_replace')" :popper="{ placement: 'top' }">
                        <UButton color="white" variant="ghost" icon="i-mdi-text-search" :disabled="disabled" />
                    </UTooltip>

                    <template #panel>
                        <div class="flex flex-1 gap-0.5 p-4">
                            <UForm :state="searchAndReplace">
                                <UFormGroup name="search" :label="$t('common.search')">
                                    <UInput v-model="searchAndReplace.search" :disabled="disabled" />
                                </UFormGroup>

                                <UFormGroup name="replace" :label="$t('components.partials.TiptapEditor.replace')">
                                    <UInput v-model="searchAndReplace.replace" :disabled="disabled" />
                                </UFormGroup>

                                <UFormGroup name="caseSensitive" :label="$t('common.case_sensitive')">
                                    <UToggle v-model="searchAndReplace.caseSensitive" :disabled="disabled" />
                                </UFormGroup>

                                <UFormGroup class="flex flex-col lg:flex-row">
                                    <UButtonGroup>
                                        <UButton
                                            color="error"
                                            variant="outline"
                                            :label="$t('components.partials.TiptapEditor.clear')"
                                            :disabled="disabled"
                                            @click="clear"
                                        />
                                        <UButton
                                            color="white"
                                            variant="outline"
                                            :label="$t('components.partials.TiptapEditor.previous')"
                                            :disabled="disabled"
                                            @click="previous"
                                        />
                                        <UButton
                                            color="white"
                                            variant="outline"
                                            :label="$t('components.partials.TiptapEditor.next')"
                                            :disabled="disabled"
                                            @click="next"
                                        />
                                        <UButton
                                            color="white"
                                            variant="outline"
                                            :label="$t('components.partials.TiptapEditor.replace')"
                                            :disabled="disabled"
                                            @click="replace"
                                        />
                                        <UButton
                                            color="white"
                                            variant="outline"
                                            :label="$t('components.partials.TiptapEditor.replace_all')"
                                            :disabled="disabled"
                                            @click="replaceAll"
                                        />
                                    </UButtonGroup>

                                    <div class="mt-1 block text-sm font-medium">
                                        {{ $t('common.result', 2) }}:
                                        {{
                                            editor?.storage?.searchAndReplace?.resultIndex > 0
                                                ? editor?.storage?.searchAndReplace?.resultIndex + 1
                                                : 0
                                        }}
                                        /
                                        {{ editor?.storage?.searchAndReplace?.results.length }}
                                    </div>
                                </UFormGroup>
                            </UForm>
                        </div>
                    </template>
                </UPopover>

                <UButtonGroup>
                    <UTooltip :text="$t('components.partials.TiptapEditor.undo')" :popper="{ placement: 'top' }">
                        <UButton
                            :disabled="!editor.can().chain().focus().undo().run() || disabled"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-undo"
                            @click="editor.chain().focus().undo().run()"
                        />
                    </UTooltip>
                    <UTooltip :text="$t('components.partials.TiptapEditor.redo')" :popper="{ placement: 'top' }">
                        <UButton
                            :disabled="!editor.can().chain().focus().redo().run() || disabled"
                            color="white"
                            variant="ghost"
                            icon="i-mdi-redo"
                            @click="editor.chain().focus().redo().run()"
                        />
                    </UTooltip>
                </UButtonGroup>

                <UDivider orientation="vertical" :ui="{ border: { base: 'border-gray-200 dark:border-gray-700' } }" />

                <UButtonGroup>
                    <UTooltip :text="$t('components.partials.TiptapEditor.source_code')" :popper="{ placement: 'top' }">
                        <UButton
                            color="white"
                            variant="ghost"
                            icon="i-mdi-file-code"
                            :disabled="disabled"
                            @click="
                                modal.open(TiptapEditorSourceCodeModal, {
                                    content: modelValue,
                                    'onUpdate:content': ($event) => (modelValue = $event),
                                })
                            "
                        />
                    </UTooltip>

                    <UTooltip
                        v-if="!disableImages && filestoreService"
                        :text="$t('components.partials.TiptapEditor.file_list')"
                        :popper="{ placement: 'top' }"
                    >
                        <UButton
                            color="white"
                            variant="ghost"
                            icon="i-mdi-file-multiple"
                            :disabled="disabled"
                            @click="
                                modal.open(FileListModal, {
                                    editor: editor!,
                                    files: files,
                                })
                            "
                        />
                    </UTooltip>

                    <UTooltip v-if="historyType" :text="$t('common.version_history')" :popper="{ placement: 'top' }">
                        <UButton
                            color="white"
                            variant="ghost"
                            icon="i-mdi-history"
                            :disabled="disabled"
                            @click="
                                modal.open(VersionHistoryModal, {
                                    historyType: historyType,
                                    currentContent: { content: modelValue, files: files },
                                    onApply: applyVersion,
                                })
                            "
                        />
                    </UTooltip>
                </UButtonGroup>
            </div>
        </div>

        <TiptapEditorContent
            ref="contentRef"
            class="min-h-0 w-full min-w-0 max-w-full flex-auto overflow-y-auto"
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
        />

        <div v-if="editor" class="flex w-full flex-none justify-between bg-gray-100 px-1 text-center dark:bg-gray-800">
            <div class="flex" :class="[{ 'flex-1': targetId }]">
                <template v-if="$slots.footer">
                    <slot name="footer" />
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

            <div v-if="targetId" class="inline-flex flex-1 items-center justify-center">
                <YJSUserPopover />
            </div>

            <div class="inline-flex flex-1 items-center justify-end">
                {{ editor.storage.characterCount.characters() }}<template v-if="limit && limit > 0"> / {{ limit }}</template>
                {{ $t('common.chars', editor.storage.characterCount.characters()) }}
                |
                {{ editor.storage.characterCount.words() }} {{ $t('common.word', editor.storage.characterCount.words()) }}
            </div>
        </div>
    </div>
</template>
