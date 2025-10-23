<script lang="ts" setup>
import type { Range } from '@tiptap/core';
import type { Editor } from '@tiptap/vue-3';
import z from 'zod';
import { fontColors, fonts, highlightColors } from '~/types/editor';
import type { Content, Version } from '~/types/history';
import type { File as FileGrpc } from '~~/gen/ts/resources/file/file';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import FileListModal from './FileListModal.vue';
import ImageSelectPopover from './ImageSelectPopover.vue';
import SourceCodeModal from './SourceCodeModal.vue';
import TablePopover from './TablePopover.vue';
import VersionHistoryModal from './VersionHistoryModal.vue';

const props = defineProps<{
    editor: Editor;
    disabled?: boolean;
    disableImages?: boolean;
    historyType?: string;

    fileLimit?: number;
    fileUploadHandler: undefined | ((files: File[]) => Promise<void>);
}>();

const emits = defineEmits<{
    (e: 'update:content', val: string): void;
}>();

const files = defineModel<FileGrpc[]>('files', { default: () => [] });

const overlay = useOverlay();

const settingsStore = useSettingsStore();
const { editor: editorSettings } = storeToRefs(settingsStore);

const notifications = useNotificationsStore();

const ed = shallowRef<Editor>(toRaw(props.editor));

const { ui, attach, detach } = useTiptapToolbar(() => unref(ed));

onMounted(attach);
onBeforeUnmount(detach);

const fileListModal = overlay.create(FileListModal);
const sourceCodeModal = overlay.create(SourceCodeModal);
const versionHistoryModal = overlay.create(VersionHistoryModal);

function setContent(value: string): void {
    ed.value?.commands.setContent(value, { emitUpdate: true });
}

watch(
    () => props.disabled,
    () => {
        ed.value?.setEditable(!props.disabled);
        ed.value?.view.updateState(ed.value!.view.state);
    },
);

const linkSchema = z.object({
    url: z.url().or(z.literal('')),
});

type LinkSchema = z.output<typeof linkSchema>;

const linkState = reactive<LinkSchema>({
    url: '',
});

function setLink(data: LinkSchema): void {
    if (data.url.trim() === '') return;

    const previousUrl = ed.value?.getAttributes('link').href;
    const url = data.url.trim() !== '' ? data.url.trim() : previousUrl;

    // Empty URL
    if (url === '') {
        ed.value?.chain().focus().extendMarkRange('link').unsetLink().run();
        return;
    }

    // Update link
    ed.value
        ?.chain()
        .focus()
        .extendMarkRange('link')
        .setLink({
            href: url,
        })
        .run();
}

const selectedFont = ref<(typeof fonts)[0]>(fonts[0]!);
watch(selectedFont, () => ed.value?.chain().focus().setFontFamily(selectedFont.value.value).run());

const selectedFontColor = ref<string | undefined>(undefined);
watch(selectedFontColor, () =>
    selectedFontColor.value
        ? ed.value?.chain().focus().setColor(selectedFontColor.value).run()
        : ed.value?.chain().focus().unsetColor().run(),
);

const selectedHighlightColor = ref<(typeof highlightColors)[0]>(highlightColors[0]!);
watch(selectedHighlightColor, () =>
    ed.value?.chain().focus().toggleHighlight({ color: selectedHighlightColor.value.value }).run(),
);

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
    if (!ed.value) return;

    if (clearIndex) ed.value.commands.resetIndex();

    ed.value?.commands.setSearchTerm(searchAndReplace.search);
    ed.value?.commands.setReplaceTerm(searchAndReplace.replace ?? '');
    ed.value?.commands.setCaseSensitive(searchAndReplace.caseSensitive);
};

const goToSelection = () => {
    if (!ed.value) return;

    const { results, resultIndex } = ed.value!.storage.searchAndReplace;
    const position: Range | undefined = results[resultIndex];

    if (!position) return;

    ed.value?.commands.setTextSelection(position);

    const { node } = ed.value!.view.domAtPos(ed.value!.state.selection.anchor);
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
    ed.value?.commands.replace();
    goToSelection();
};

const next = () => {
    ed.value?.commands.nextSearchResult();
    goToSelection();
};

const previous = () => {
    ed.value?.commands.previousSearchResult();
    goToSelection();
};

const clear = () => {
    searchAndReplace.search = searchAndReplace.replace = '';
    ed.value?.commands.resetIndex();
};

const replaceAll = () => ed.value?.commands.replaceAll();

function applyVersion(version: Version<unknown>): void {
    const v = version as Version<Content>;
    emits('update:content', v.content.content);
    files.value = v.content.files;

    notifications.add({
        title: { key: 'notifications.action_successful.title', parameters: {} },
        description: { key: 'notifications.action_successful.content', parameters: {} },
        type: NotificationType.SUCCESS,
    });
}

const formRef = useTemplateRef('formRef');

const isLinkOpen = ref(false);
</script>

<template>
    <div class="flex snap-x flex-wrap gap-1">
        <UButtonGroup>
            <UTooltip :text="$t('components.partials.tiptap_editor.bold')">
                <UButton
                    :class="{ 'bg-neutral-300 dark:bg-neutral-900': ui.bold }"
                    :disabled="disabled || !ui.canBold"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-bold"
                    @click="ed?.chain().focus().toggleBold().run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.italic')">
                <UButton
                    :class="{ 'bg-neutral-300 dark:bg-neutral-900': ui.italic }"
                    :disabled="disabled || !ui.canItalic"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-italic"
                    @click="ed?.chain().focus().toggleItalic().run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.underline')">
                <UButton
                    :class="{ 'bg-neutral-300 dark:bg-neutral-900': ui.underline }"
                    :disabled="disabled || !ui.canUnderline"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-underline"
                    @click="ed?.chain().focus().toggleUnderline().run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.strike')">
                <UButton
                    :class="{ 'bg-neutral-300 dark:bg-neutral-900': ui.strike }"
                    :disabled="disabled || !ui.canStrike"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-strikethrough"
                    @click="ed?.chain().focus().toggleStrike().run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.clear')">
                <UButton
                    :disabled="disabled"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-clear"
                    @click="ed?.chain().focus().unsetAllMarks().run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.superscript')">
                <UButton
                    :class="{ 'bg-neutral-300 dark:bg-neutral-900': ui.superscript }"
                    :disabled="disabled"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-superscript"
                    @click="ed?.chain().focus().toggleSuperscript().run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.subscript')">
                <UButton
                    :class="{ 'bg-neutral-300 dark:bg-neutral-900': ui.subscript }"
                    :disabled="disabled"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-subscript"
                    @click="ed?.chain().focus().toggleSubscript().run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_ed?.code')">
                <UButton
                    :class="{ 'bg-neutral-300 dark:bg-neutral-900': ui.code }"
                    :disabled="disabled"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-code-braces"
                    @click="ed?.chain().focus().toggleCode().run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.invisible_characters')">
                <UButton
                    :class="{ 'bg-neutral-300 dark:bg-neutral-900': editorSettings.showInvisibleCharacters }"
                    :disabled="disabled"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-pilcrow"
                    @click="editorSettings.showInvisibleCharacters = !editorSettings.showInvisibleCharacters"
                />
            </UTooltip>
        </UButtonGroup>

        <USeparator orientation="vertical" :ui="{ border: 'border-neutral-200 dark:border-neutral-700' }" />

        <!-- Text Align -->
        <UButtonGroup>
            <UTooltip :text="$t('components.partials.tiptap_editor.align_left')">
                <UButton
                    :class="{ 'bg-neutral-300 dark:bg-neutral-900': ui.textAlign === 'left' }"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-align-left"
                    :disabled="disabled"
                    @click="ed?.chain().focus().setTextAlign('left').run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.align_center')">
                <UButton
                    :class="{ 'bg-neutral-300 dark:bg-neutral-900': ui.textAlign === 'center' }"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-align-center"
                    :disabled="disabled"
                    @click="ed?.chain().focus().setTextAlign('center').run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.align_right')">
                <UButton
                    :class="{ 'bg-neutral-300 dark:bg-neutral-900': ui.textAlign === 'right' }"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-align-right"
                    :disabled="disabled"
                    @click="ed?.chain().focus().setTextAlign('right').run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.align_justify')">
                <UButton
                    :class="{ 'bg-neutral-300 dark:bg-neutral-900': ui.textAlign === 'justify' }"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-align-justify"
                    :disabled="disabled"
                    @click="ed?.chain().focus().setTextAlign('justify').run()"
                />
            </UTooltip>
        </UButtonGroup>

        <USeparator orientation="vertical" :ui="{ border: 'border-neutral-200 dark:border-neutral-700' }" />

        <!-- Font Family -->
        <UTooltip :text="$t('components.partials.tiptap_editor.font_family')">
            <UInputMenu
                v-model="selectedFont"
                class="w-full max-w-44"
                name="selectedFont"
                :filter-fields="['label']"
                :items="fonts"
                :placeholder="$t('common.font', 1)"
                :disabled="disabled"
                :style="{ fontFamily: selectedFont.value }"
            >
                <template #item-label="{ item }">
                    <span class="truncate" :style="{ fontFamily: item.value }">{{
                        item.label.includes('.') ? $t(item.label) : item.label
                    }}</span>
                </template>

                <template #empty>
                    {{ $t('common.not_found', [$t('components.partials.tiptap_editor.font_family')]) }}
                </template>
            </UInputMenu>
        </UTooltip>

        <UButtonGroup>
            <UPopover>
                <UTooltip :text="$t('components.partials.tiptap_editor.font_color')">
                    <UButton
                        :class="{
                            'bg-neutral-300 dark:bg-neutral-900': ui.fontColor === selectedFontColor,
                        }"
                        color="neutral"
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
                            color="neutral"
                            variant="outline"
                            icon="i-mdi-water-off"
                            :label="$t('common.default')"
                            :disabled="disabled"
                            @click="
                                ed?.chain().focus().unsetColor().run();
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
            <UTooltip :text="$t('components.partials.tiptap_editor.paragraph')">
                <UButton
                    :class="{ 'bg-neutral-300 dark:bg-neutral-900': ui.paragraph }"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-paragraph"
                    :disabled="disabled"
                    @click="ed?.chain().focus().setParagraph().run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.header_1')">
                <UButton
                    :class="{ 'bg-neutral-300 dark:bg-neutral-900': ui.headingLevel === 1 }"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-header-1"
                    :disabled="disabled"
                    @click="ed?.chain().focus().toggleHeading({ level: 1 }).run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.header_2')">
                <UButton
                    :class="{ 'bg-neutral-300 dark:bg-neutral-900': ui.headingLevel === 2 }"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-header-2"
                    :disabled="disabled"
                    @click="ed?.chain().focus().toggleHeading({ level: 2 }).run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.header_3')">
                <UButton
                    :class="{ 'bg-neutral-300 dark:bg-neutral-900': ui.headingLevel === 3 }"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-header-3"
                    :disabled="disabled"
                    @click="ed?.chain().focus().toggleHeading({ level: 3 }).run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.header_4')">
                <UButton
                    :class="{ 'bg-neutral-300 dark:bg-neutral-900': ui.headingLevel === 4 }"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-header-4"
                    :disabled="disabled"
                    @click="ed?.chain().focus().toggleHeading({ level: 4 }).run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.header_5')">
                <UButton
                    :class="{ 'bg-neutral-300 dark:bg-neutral-900': ui.headingLevel === 5 }"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-header-5"
                    :disabled="disabled"
                    @click="ed?.chain().focus().toggleHeading({ level: 5 }).run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.header_6')">
                <UButton
                    :class="{ 'bg-neutral-300 dark:bg-neutral-900': ui.headingLevel === 6 }"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-header-6"
                    :disabled="disabled"
                    @click="ed?.chain().focus().toggleHeading({ level: 6 }).run()"
                />
            </UTooltip>
        </UButtonGroup>
    </div>

    <div class="flex snap-x flex-wrap gap-1">
        <UButtonGroup>
            <UTooltip :text="$t('components.partials.tiptap_editor.highlight')">
                <UButton
                    :class="{ 'bg-neutral-300 dark:bg-neutral-900': ui.highlight }"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-color-highlight"
                    :disabled="disabled"
                    @click="ed?.chain().focus().toggleHighlight().run()"
                />
            </UTooltip>

            <UPopover>
                <UTooltip :text="$t('components.partials.tiptap_editor.highlight_color')">
                    <UButton
                        :class="{
                            'bg-neutral-300 dark:bg-neutral-900':
                                ui.highlight && selectedHighlightColor.value === ui.highlightColor,
                        }"
                        color="neutral"
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
                            color="neutral"
                            variant="outline"
                            icon="i-mdi-water-off"
                            :label="$t('common.reset')"
                            :disabled="disabled"
                            @click="
                                ed?.chain().focus().unsetHighlight().run();
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

        <USeparator orientation="vertical" :ui="{ border: 'border-neutral-200 dark:border-neutral-700' }" />

        <UButtonGroup>
            <UTooltip :text="$t('components.partials.tiptap_editor.bullet_list')">
                <UButton
                    :class="{ 'bg-neutral-300 dark:bg-neutral-900': ui.bulletList }"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-list-bulleted"
                    :disabled="disabled"
                    @click="ed?.chain().focus().toggleBulletList().run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.ordered_list')">
                <UButton
                    :class="{ 'bg-neutral-300 dark:bg-neutral-900': ui.orderedList }"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-list-numbered"
                    :disabled="disabled"
                    @click="ed?.chain().focus().toggleOrderedList().run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.task_list')">
                <UButton
                    :class="{ 'bg-neutral-300 dark:bg-neutral-900': ui.taskList }"
                    icon="i-mdi-format-list-checks"
                    color="neutral"
                    variant="ghost"
                    :disabled="disabled"
                    @click="ed?.chain().focus().toggleTaskList().run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_ed?.checkbox')">
                <UButton
                    :class="{ 'bg-neutral-300 dark:bg-neutral-900': ui.checkboxStandalone }"
                    icon="i-mdi-checkbox-marked-outline"
                    color="neutral"
                    variant="ghost"
                    :disabled="disabled"
                    @click="ed?.chain().focus().addCheckboxStandalone().run()"
                />
            </UTooltip>
        </UButtonGroup>

        <USeparator orientation="vertical" :ui="{ border: 'border-neutral-200 dark:border-neutral-700' }" />

        <ImageSelectPopover
            v-if="!disableImages"
            :editor="editor"
            :file-list="files"
            :file-limit="fileLimit"
            :disabled="disabled"
            :upload-handler="fileUploadHandler"
            @open-file-list="
                fileListModal.open({
                    editor: unref(editor)!,
                    files: files,
                })
            "
        />

        <TablePopover :editor="unref(editor)" :disabled="disabled" />

        <UPopover v-model:open="isLinkOpen">
            <UTooltip :text="$t('components.partials.tiptap_editor.link')">
                <UButton
                    :class="{ 'bg-neutral-300 dark:bg-neutral-900': ui.link }"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-link"
                    :disabled="disabled"
                />
            </UTooltip>

            <template #content>
                <div class="p-4">
                    <UForm ref="formRef" :state="linkState" :schema="linkSchema" @submit="($event) => setLink($event.data)">
                        <UFormField :label="$t('common.url')" name="url">
                            <UInput v-model="linkState.url" type="text" :disabled="disabled" />
                        </UFormField>

                        <slot name="linkModal" :editor="editor" :state="linkState" />

                        <UButtonGroup class="mt-2 w-full">
                            <UButton
                                class="flex-1"
                                type="submit"
                                icon="i-mdi-link"
                                :label="$t('common.link')"
                                :disabled="disabled"
                                @click="formRef?.submit()"
                            />

                            <UButton
                                :disabled="!editor.isActive('link') || disabled"
                                color="error"
                                variant="outline"
                                icon="i-mdi-link-off"
                                :label="$t('common.unlink')"
                                @click="
                                    isLinkOpen = false;
                                    ed?.chain().focus().unsetLink().run();
                                    linkState.url = '';
                                "
                            />
                        </UButtonGroup>
                    </UForm>
                </div>
            </template>
        </UPopover>

        <UButtonGroup>
            <UTooltip :text="$t('components.partials.tiptap_ed?.code_block')">
                <UButton
                    :class="{ 'bg-neutral-300 dark:bg-neutral-900': ui.codeBlock }"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-code-block-braces"
                    :disabled="disabled"
                    @click="ed?.chain().focus().toggleCodeBlock().run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.block_quote')">
                <UButton
                    :class="{ 'bg-neutral-300 dark:bg-neutral-900': ui.blockquote }"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-quote-open"
                    :disabled="disabled"
                    @click="ed?.chain().focus().toggleBlockquote().run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.horizontal_rule')">
                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-minus"
                    :disabled="disabled"
                    @click="ed?.chain().focus().setHorizontalRule().run()"
                />
            </UTooltip>
            <!--
                    <UButton
                        color="neutral"
                        variant="ghost"
                        icon="i-mdi-format-page-break"
                        :disabled="disabled"
                        @click="ed?.chain().focus().setHardBreak().run()"
                    />
                    -->
        </UButtonGroup>

        <div class="flex-1"></div>

        <slot name="toolbar" :editor="editor" :disabled="disabled" />

        <USeparator orientation="vertical" :ui="{ border: 'border-neutral-200 dark:border-neutral-700' }" />

        <UPopover>
            <UTooltip :text="$t('components.partials.tiptap_editor.search_and_replace')">
                <UButton color="neutral" variant="ghost" icon="i-mdi-text-search" :disabled="disabled" />
            </UTooltip>

            <template #content>
                <div class="flex flex-1 gap-0.5 p-4">
                    <UForm :state="searchAndReplace">
                        <UFormField name="search" :label="$t('common.search')">
                            <UInput v-model="searchAndReplace.search" :disabled="disabled" />
                        </UFormField>

                        <UFormField name="replace" :label="$t('components.partials.tiptap_editor.replace')">
                            <UInput v-model="searchAndReplace.replace" :disabled="disabled" />
                        </UFormField>

                        <UFormField name="caseSensitive" :label="$t('common.case_sensitive')">
                            <USwitch v-model="searchAndReplace.caseSensitive" :disabled="disabled" />
                        </UFormField>

                        <UFormField class="flex flex-col lg:flex-row">
                            <UButtonGroup>
                                <UButton
                                    color="error"
                                    variant="outline"
                                    :label="$t('components.partials.tiptap_editor.clear')"
                                    :disabled="disabled"
                                    @click="() => clear()"
                                />
                                <UButton
                                    color="neutral"
                                    variant="outline"
                                    :label="$t('components.partials.tiptap_editor.previous')"
                                    :disabled="disabled"
                                    @click="() => previous()"
                                />
                                <UButton
                                    color="neutral"
                                    variant="outline"
                                    :label="$t('components.partials.tiptap_editor.next')"
                                    :disabled="disabled"
                                    @click="() => next()"
                                />
                                <UButton
                                    color="neutral"
                                    variant="outline"
                                    :label="$t('components.partials.tiptap_editor.replace')"
                                    :disabled="disabled"
                                    @click="() => replace()"
                                />
                                <UButton
                                    color="neutral"
                                    variant="outline"
                                    :label="$t('components.partials.tiptap_editor.replace_all')"
                                    :disabled="disabled"
                                    @click="
                                        () => {
                                            replaceAll();
                                        }
                                    "
                                />
                            </UButtonGroup>

                            <div class="mt-1 block text-sm font-medium">
                                {{ $t('common.result', 2) }}:
                                {{
                                    ed?.storage?.searchAndReplace?.resultIndex > 0
                                        ? ed?.storage?.searchAndReplace?.resultIndex + 1
                                        : 0
                                }}
                                /
                                {{ ed?.storage?.searchAndReplace?.results.length }}
                            </div>
                        </UFormField>
                    </UForm>
                </div>
            </template>
        </UPopover>

        <UButtonGroup>
            <UTooltip :text="$t('components.partials.tiptap_editor.undo')">
                <UButton
                    :disabled="disabled || !ui.canUndo"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-undo"
                    @click="ed?.chain().focus().undo().run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.redo')">
                <UButton
                    :disabled="disabled || !ui.canRedo"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-redo"
                    @click="ed?.chain().focus().redo().run()"
                />
            </UTooltip>
        </UButtonGroup>

        <USeparator orientation="vertical" :ui="{ border: 'border-neutral-200 dark:border-neutral-700' }" />

        <UButtonGroup>
            <UTooltip :text="$t('components.partials.tiptap_editor.source_code')">
                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-file-code"
                    :disabled="disabled"
                    @click="
                        sourceCodeModal.open({
                            content: ed?.getHTML() || '',
                            disabled: disabled,
                            'onUpdate:content': ($event) => setContent($event),
                        })
                    "
                />
            </UTooltip>

            <UTooltip v-if="!disableImages && fileUploadHandler" :text="$t('components.partials.tiptap_editor.file_list')">
                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-file-multiple"
                    :disabled="disabled"
                    @click="
                        fileListModal.open({
                            editor: unref(editor)!,
                            files: files,
                        })
                    "
                />
            </UTooltip>

            <UTooltip v-if="historyType" :text="$t('common.version_history')">
                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-history"
                    :disabled="disabled"
                    @click="
                        versionHistoryModal.open({
                            historyType: historyType,
                            currentContent: { content: ed?.getHTML() || '', files: files },
                            onApply: applyVersion,
                        })
                    "
                />
            </UTooltip>
        </UButtonGroup>
    </div>
</template>
