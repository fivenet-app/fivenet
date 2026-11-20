<script lang="ts" setup>
import type { Editor, JSONContent } from '@tiptap/core';
import { z } from 'zod';
import { fontColors, fonts, highlightColors } from '~/types/editor';
import type { HistoryContent, Version } from '~/types/history';
import type { File as FileGrpc } from '~~/gen/ts/resources/file/file';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import FileListModal from './FileListModal.vue';
import MapBlockInsertPopover from './MapBlockInsertPopover.vue';
import ImageSelectPopover from './ImageSelectPopover.vue';
import SearchAndReplacePopover from './SearchAndReplacePopover.vue';
import SourceCodeModal from './SourceCodeModal.vue';
import TablePopover from './TablePopover.vue';
import VersionHistoryModal from './VersionHistoryModal.vue';

const props = defineProps<{
    editor: Editor;
    disabled?: boolean;
    disableImages?: boolean;
    enableImagesJobAssets?: boolean;
    showPenaltyCalculatorButton?: boolean;
    historyType?: string;

    fileLimit?: number;
    fileUploadHandler: undefined | ((files: File[]) => Promise<boolean>);
}>();

const emits = defineEmits<{
    (e: 'update:content', val: JSONContent | string | undefined): void;
}>();

const files = defineModel<FileGrpc[]>('files', {
    required: true,
});

const { can } = useAuth();
const canListJobAssets = can('settings.JobsAssetsService/ListJobAssets');

type TextAlign = 'left' | 'center' | 'right' | 'justify';
type TextBlock = 'paragraph' | 'heading-1' | 'heading-2' | 'heading-3' | 'heading-4' | 'heading-5';
type ListType = 'bullet' | 'ordered' | 'task';

const { t } = useI18n();

const textAlignItems = computed(() => [
    {
        value: 'left' as const,
        icon: 'i-mdi-format-align-left',
        'aria-label': t('components.partials.tiptap_editor.align_left'),
        activeClass: 'bg-accented',
    },
    {
        value: 'center' as const,
        icon: 'i-mdi-format-align-center',
        'aria-label': t('components.partials.tiptap_editor.align_center'),
        activeClass: 'bg-accented',
    },
    {
        value: 'right' as const,
        icon: 'i-mdi-format-align-right',
        'aria-label': t('components.partials.tiptap_editor.align_right'),
        activeClass: 'bg-accented',
    },
    {
        value: 'justify' as const,
        icon: 'i-mdi-format-align-justify',
        'aria-label': t('components.partials.tiptap_editor.align_justify'),
        activeClass: 'bg-accented',
    },
]);

const textBlockItems = computed(() => [
    {
        value: 'paragraph' as const,
        icon: 'i-mdi-format-paragraph',
        'aria-label': t('components.partials.tiptap_editor.paragraph'),
        activeClass: 'bg-accented',
    },
    ...[1, 2, 3, 4, 5].map((level) => ({
        value: `heading-${level}` as TextBlock,
        icon: `i-mdi-format-header-${level}`,
        'aria-label': t(`components.partials.tiptap_editor.heading_${level}`),
        activeClass: 'bg-accented',
    })),
]);

const listItems = computed(() => [
    {
        value: 'bullet' as const,
        icon: 'i-mdi-format-list-bulleted',
        'aria-label': t('components.partials.tiptap_editor.bullet_list'),
        activeClass: 'bg-accented',
    },
    {
        value: 'ordered' as const,
        icon: 'i-mdi-format-list-numbered',
        'aria-label': t('components.partials.tiptap_editor.ordered_list'),
        activeClass: 'bg-accented',
    },
    {
        value: 'task' as const,
        icon: 'i-mdi-format-list-checks',
        'aria-label': t('components.partials.tiptap_editor.task_list'),
        activeClass: 'bg-accented',
    },
]);

const overlay = useOverlay();

const settingsStore = useSettingsStore();
const { editor: editorSettings } = storeToRefs(settingsStore);

const notifications = useNotificationsStore();

const ed = shallowRef<Editor>(toRaw(props.editor));

const { ui, attach, detach } = useTiptapToolbar(() => unref(ed));

const activeTextBlock = computed<TextBlock | undefined>(() => {
    if (ui.paragraph) return 'paragraph';
    if (ui.headingLevel >= 1 && ui.headingLevel <= 5) return `heading-${ui.headingLevel}` as TextBlock;

    return undefined;
});

const activeListType = computed<ListType | undefined>(() => {
    if (ui.bulletList) return 'bullet';
    if (ui.orderedList) return 'ordered';
    if (ui.taskList) return 'task';

    return undefined;
});

onMounted(attach);
onBeforeUnmount(detach);

const fileListModal = overlay.create(FileListModal);
const sourceCodeModal = overlay.create(SourceCodeModal);
const versionHistoryModal = overlay.create(VersionHistoryModal);

function setContent(value: string): void {
    ed.value?.commands.setContent(value, { emitUpdate: true });
}

function setTextAlign(value: unknown): void {
    if (!['left', 'center', 'right', 'justify'].includes(value as TextAlign)) return;

    ed.value
        ?.chain()
        .focus()
        .setTextAlign(value as TextAlign)
        .run();
}

function setTextBlock(value: unknown): void {
    if (value === undefined || value === 'paragraph') {
        ed.value?.chain().focus().setParagraph().run();
        return;
    }

    const headingLevel = Number(String(value).replace('heading-', ''));
    if (!Number.isInteger(headingLevel) || headingLevel < 1 || headingLevel > 5) return;

    ed.value
        ?.chain()
        .focus()
        .toggleHeading({ level: headingLevel as 1 | 2 | 3 | 4 | 5 })
        .run();
}

function setListType(value: unknown): void {
    const listType = value ?? activeListType.value;

    if (listType === 'bullet') {
        ed.value?.chain().focus().toggleBulletList().run();
    } else if (listType === 'ordered') {
        ed.value?.chain().focus().toggleOrderedList().run();
    } else if (listType === 'task') {
        ed.value?.chain().focus().toggleTaskList().run();
    }
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

function applyVersion(version: Version<HistoryContent>): void {
    emits('update:content', version.content.content);
    files.value = version.content.files;

    notifications.add({
        title: { key: 'notifications.action_successful.title', parameters: {} },
        description: { key: 'notifications.action_successful.content', parameters: {} },
        type: NotificationType.SUCCESS,
    });
}

const formRef = useTemplateRef('formRef');

const isLinkOpen = ref<boolean>(false);
</script>

<template>
    <div class="flex snap-x flex-wrap gap-1">
        <UFieldGroup>
            <UTooltip :text="$t('components.partials.tiptap_editor.bold')">
                <UButton
                    :class="{ 'bg-accented': ui.bold }"
                    :disabled="disabled || !ui.canBold"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-bold"
                    @click="() => ed?.chain().focus().toggleBold().run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.italic')">
                <UButton
                    :class="{ 'bg-accented': ui.italic }"
                    :disabled="disabled || !ui.canItalic"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-italic"
                    @click="() => ed?.chain().focus().toggleItalic().run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.underline')">
                <UButton
                    :class="{ 'bg-accented': ui.underline }"
                    :disabled="disabled || !ui.canUnderline"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-underline"
                    @click="() => ed?.chain().focus().toggleUnderline().run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.strike')">
                <UButton
                    :class="{ 'bg-accented': ui.strike }"
                    :disabled="disabled || !ui.canStrike"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-strikethrough"
                    @click="() => ed?.chain().focus().toggleStrike().run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.clear')">
                <UButton
                    :disabled="disabled"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-clear"
                    @click="() => ed?.chain().focus().unsetAllMarks().run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.superscript')">
                <UButton
                    :class="{ 'bg-accented': ui.superscript }"
                    :disabled="disabled"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-superscript"
                    @click="() => ed?.chain().focus().toggleSuperscript().run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.subscript')">
                <UButton
                    :class="{ 'bg-accented': ui.subscript }"
                    :disabled="disabled"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-subscript"
                    @click="() => ed?.chain().focus().toggleSubscript().run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.code')">
                <UButton
                    :class="{ 'bg-accented': ui.code }"
                    :disabled="disabled"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-code-braces"
                    @click="() => ed?.chain().focus().toggleCode().run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.invisible_characters')">
                <UButton
                    :class="{ 'bg-accented': editorSettings.showInvisibleCharacters }"
                    :disabled="disabled"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-pilcrow"
                    @click="editorSettings.showInvisibleCharacters = !editorSettings.showInvisibleCharacters"
                />
            </UTooltip>
        </UFieldGroup>

        <USeparator orientation="vertical" :ui="{ border: 'border-neutral-200 dark:border-neutral-700' }" />

        <!-- Text Align -->
        <UFormField>
            <UTooltip :text="$t('components.partials.tiptap_editor.text_alignment')">
                <UToggleGroup
                    :model-value="ui.textAlign"
                    :items="textAlignItems"
                    :disabled="disabled"
                    :aria-label="$t('components.partials.tiptap_editor.text_alignment')"
                    type="single"
                    attached
                    color="neutral"
                    variant="ghost"
                    active-color="neutral"
                    active-variant="ghost"
                    @update:model-value="setTextAlign"
                />
            </UTooltip>
        </UFormField>

        <USeparator orientation="vertical" :ui="{ border: 'border-neutral-200 dark:border-neutral-700' }" />

        <!-- Font Family -->
        <UFormField name="selectedFont" class="w-48 max-w-full shrink-0">
            <UTooltip :text="$t('components.partials.tiptap_editor.font_family')">
                <USelectMenu
                    v-model="selectedFont"
                    class="w-full"
                    name="selectedFont"
                    leading-icon="i-mdi-format-font"
                    :filter-fields="['label']"
                    :items="fonts"
                    :placeholder="$t('components.partials.tiptap_editor.font_family', 1)"
                    :disabled="disabled"
                    :style="{ fontFamily: selectedFont.value }"
                >
                    <template #default>
                        <span class="truncate" :style="{ fontFamily: selectedFont.value }">{{
                            selectedFont.label.includes('.') ? $t(selectedFont.label) : selectedFont.label
                        }}</span>
                    </template>

                    <template #item-label="{ item }">
                        <span class="truncate" :style="{ fontFamily: item.value }">{{
                            item.label.includes('.') ? $t(item.label) : item.label
                        }}</span>
                    </template>

                    <template #empty>
                        {{ $t('common.not_found', [$t('components.partials.tiptap_editor.font_family')]) }}
                    </template>
                </USelectMenu>
            </UTooltip>
        </UFormField>

        <UPopover>
            <UTooltip :text="$t('components.partials.tiptap_editor.font_color')">
                <UButton
                    :class="{
                        'bg-accented': ui.fontColor === selectedFontColor,
                    }"
                    color="neutral"
                    variant="ghost"
                    :style="{ color: selectedFontColor }"
                    icon="i-mdi-format-color-text"
                    :disabled="disabled"
                />
            </UTooltip>

            <template #content="{ close }">
                <div class="inline-flex flex-col gap-1 p-4">
                    <UButton
                        class="rounded-md"
                        color="neutral"
                        variant="outline"
                        icon="i-mdi-water-off"
                        :label="$t('components.partials.tiptap_editor.default')"
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
        <UTooltip :text="$t('components.partials.tiptap_editor.text_block')">
            <UToggleGroup
                :model-value="activeTextBlock"
                :items="textBlockItems"
                :disabled="disabled"
                :aria-label="$t('components.partials.tiptap_editor.text_block')"
                type="single"
                attached
                color="neutral"
                variant="ghost"
                active-color="neutral"
                active-variant="ghost"
                @update:model-value="setTextBlock"
            />
        </UTooltip>
    </div>

    <div class="flex snap-x flex-wrap gap-1">
        <UFieldGroup>
            <UTooltip :text="$t('components.partials.tiptap_editor.highlight')">
                <UButton
                    :class="{ 'bg-accented': ui.highlight }"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-color-highlight"
                    :disabled="disabled"
                    @click="() => ed?.chain().focus().toggleHighlight().run()"
                />
            </UTooltip>

            <UPopover>
                <UTooltip :text="$t('components.partials.tiptap_editor.highlight_color')">
                    <UButton
                        :class="{
                            'bg-accented': ui.highlight && selectedHighlightColor.value === ui.highlightColor,
                        }"
                        color="neutral"
                        variant="ghost"
                        :style="{ color: selectedHighlightColor.value }"
                        icon="i-mdi-format-color-fill"
                        :disabled="disabled"
                    />
                </UTooltip>

                <template #content="{ close }">
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
        </UFieldGroup>

        <USeparator orientation="vertical" :ui="{ border: 'border-neutral-200 dark:border-neutral-700' }" />

        <UFormField>
            <UTooltip :text="$t('components.partials.tiptap_editor.lists')">
                <UToggleGroup
                    :model-value="activeListType"
                    :items="listItems"
                    :disabled="disabled"
                    :aria-label="$t('components.partials.tiptap_editor.lists')"
                    type="single"
                    attached
                    color="neutral"
                    variant="ghost"
                    active-color="neutral"
                    active-variant="ghost"
                    @update:model-value="setListType"
                />
            </UTooltip>
        </UFormField>

        <USeparator orientation="vertical" :ui="{ border: 'border-neutral-200 dark:border-neutral-700' }" />

        <UTooltip :text="$t('components.partials.tiptap_editor.checkbox')">
            <UButton
                :class="{ 'bg-accented': ui.checkboxStandalone }"
                icon="i-mdi-checkbox-marked-outline"
                color="neutral"
                variant="ghost"
                :disabled="disabled"
                @click="() => ed?.chain().focus().addCheckboxStandalone().run()"
            />
        </UTooltip>

        <USeparator orientation="vertical" :ui="{ border: 'border-neutral-200 dark:border-neutral-700' }" />

        <ImageSelectPopover
            v-if="!disableImages || (enableImagesJobAssets && canListJobAssets)"
            :editor="editor"
            :files="files"
            :file-limit="fileLimit"
            :disable-images="disableImages"
            :enable-images-job-assets="enableImagesJobAssets"
            :disabled="disabled"
            :upload-handler="fileUploadHandler"
            @open-file-list="
                () =>
                    fileListModal.open({
                        editor: unref(editor)!,
                        files: files ?? [],
                    })
            "
        />

        <MapBlockInsertPopover :editor="editor" :disabled="disabled" />

        <TablePopover :editor="unref(editor)" :active="ui.table" :disabled="disabled" />

        <UPopover v-model:open="isLinkOpen">
            <UTooltip :text="$t('components.partials.tiptap_editor.link')">
                <UButton
                    :class="{ 'bg-accented': ui.link }"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-link"
                    :disabled="disabled"
                />
            </UTooltip>

            <template #content>
                <div class="p-4">
                    <UForm ref="formRef" :state="linkState" :schema="linkSchema" @submit="($event) => setLink($event.data)">
                        <UFormField :label="$t('components.partials.tiptap_editor.url')" name="url">
                            <UInput v-model="linkState.url" type="text" :disabled="disabled" />
                        </UFormField>

                        <slot name="linkModal" :editor="editor" :state="linkState" />

                        <UFieldGroup class="mt-2 w-full">
                            <UButton
                                class="flex-1"
                                type="submit"
                                icon="i-mdi-link"
                                :label="$t('components.partials.tiptap_editor.link')"
                                :disabled="disabled"
                                @click="formRef?.submit()"
                            />

                            <UButton
                                :disabled="!ui.link || disabled"
                                color="error"
                                variant="outline"
                                icon="i-mdi-link-off"
                                :label="$t('components.partials.tiptap_editor.unlink')"
                                @click="
                                    () => {
                                        isLinkOpen = false;
                                        ed?.chain().focus().unsetLink().run();
                                        linkState.url = '';
                                    }
                                "
                            />
                        </UFieldGroup>
                    </UForm>
                </div>
            </template>
        </UPopover>

        <UFieldGroup>
            <UTooltip :text="$t('components.partials.tiptap_editor.block_quote')">
                <UButton
                    :class="{ 'bg-accented': ui.blockquote }"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-quote-open"
                    :disabled="disabled"
                    @click="() => ed?.chain().focus().toggleBlockquote().run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.code_block')">
                <UButton
                    :class="{ 'bg-accented': ui.codeBlock }"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-code-block-braces"
                    :disabled="disabled"
                    @click="() => ed?.chain().focus().toggleCodeBlock().run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.horizontal_rule')">
                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-minus"
                    :disabled="disabled"
                    @click="() => ed?.chain().focus().setHorizontalRule().run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.details')">
                <UButton
                    :class="{ 'bg-accented': ui.details }"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-details"
                    :disabled="disabled"
                    @click="
                        () => (!ui.details ? ed?.chain().focus().setDetails().run() : ed?.chain().focus().unsetDetails().run())
                    "
                />
            </UTooltip>

            <UTooltip v-if="showPenaltyCalculatorButton" :text="$t('components.partials.tiptap_editor.penalty_calculator')">
                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-gavel"
                    :disabled="disabled"
                    @click="() => ed?.chain().insertPenaltyCalculator().run()"
                />
            </UTooltip>

            <!--
                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-format-page-break"
                    :disabled="disabled"
                    @click="() => ed?.chain().focus().setHardBreak().run()"
                />
            -->
        </UFieldGroup>

        <div class="flex-1"></div>

        <slot name="toolbar" :editor="editor" :disabled="disabled" />

        <USeparator orientation="vertical" :ui="{ border: 'border-neutral-200 dark:border-neutral-700' }" />

        <UTooltip :text="$t('components.partials.tiptap_editor.focus_mode')">
            <UButton
                :class="{ 'bg-accented': editorSettings.focusMode }"
                :disabled="disabled"
                color="neutral"
                variant="ghost"
                icon="i-mdi-border-radius"
                @click="() => (editorSettings.focusMode = !editorSettings.focusMode)"
            />
        </UTooltip>

        <SearchAndReplacePopover :editor="unref(editor)" :disabled="disabled" />

        <UFieldGroup>
            <UTooltip :text="$t('components.partials.tiptap_editor.undo')">
                <UButton
                    :disabled="disabled || !ui.canUndo"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-undo"
                    @click="() => ed?.chain().focus().undo().run()"
                />
            </UTooltip>

            <UTooltip :text="$t('components.partials.tiptap_editor.redo')">
                <UButton
                    :disabled="disabled || !ui.canRedo"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-redo"
                    @click="() => ed?.chain().focus().redo().run()"
                />
            </UTooltip>
        </UFieldGroup>

        <USeparator orientation="vertical" :ui="{ border: 'border-neutral-200 dark:border-neutral-700' }" />

        <UFieldGroup>
            <UTooltip :text="$t('components.partials.tiptap_editor.source_code')">
                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-file-code"
                    :disabled="disabled"
                    @click="
                        () =>
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
                        () =>
                            fileListModal.open({
                                editor: unref(editor)!,
                                files: files ?? [],
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
                        () =>
                            historyType &&
                            versionHistoryModal.open({
                                historyType: historyType,
                                currentContent: { content: ed?.getJSON(), files: files ?? [] },
                                onApply: applyVersion,
                            })
                    "
                />
            </UTooltip>
        </UFieldGroup>
    </div>
</template>
