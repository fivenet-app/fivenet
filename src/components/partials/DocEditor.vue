<script lang="ts" setup>
import 'jodit/es5/jodit.min.css';
import { Jodit } from 'jodit';
// @ts-ignore jodit-vue has types, but they are not detected
import { JoditEditor } from 'jodit-vue';
import type { IJodit } from 'jodit/types/types';
import { useSettingsStore } from '~/store/settings';

const props = withDefaults(
    defineProps<{
        modelValue: string;
        disabled?: boolean;
        splitScreen?: boolean;
        minHeight?: number;
    }>(),
    {
        disabled: false,
        splitScreen: false,
        minHeight: 475,
    },
);

const emits = defineEmits<{
    (e: 'update:modelValue', content: string): void;
}>();

const editorRef = ref<JoditEditor | null>(null);

const content = useVModel(props, 'modelValue', emits);

const settingsStore = useSettingsStore();
const { design: theme, locale } = storeToRefs(settingsStore);

const config = {
    readOnly: props.disabled,

    language: locale.value,
    spellcheck: true,
    minHeight: props.minHeight,
    editorClassName: 'prose' + (theme.value.documents.editorTheme === 'dark' ? ' prose-neutral' : ' prose-gray'),
    theme: theme.value.documents.editorTheme,
    showPlaceholder: false,

    defaultMode: props.splitScreen ? '3' : '1',

    disablePlugins: ['about', 'poweredByJodit', 'classSpan', 'file', 'video', 'print', 'preview', 'aiAssistant'],
    defaultActionOnPaste: 'insert_clear_html',

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
        formTemplate: (_: Jodit) =>
            `<form class="inline-flex gap-2">
                <input ref="url_input" class="relative block w-full disabled:cursor-not-allowed disabled:opacity-75 focus:outline-none border-0 form-input rounded-md placeholder-gray-400 dark:placeholder-gray-500 text-sm px-2.5 py-1.5 shadow-sm bg-white dark:bg-gray-900 text-gray-900 dark:text-white ring-1 ring-inset ring-gray-300 dark:ring-gray-700 focus:ring-2 focus:ring-primary-500 dark:focus:ring-primary-400" placeholder="Link" />
                <button class="rounded-md focus:outline-none disabled:cursor-not-allowed disabled:opacity-75 flex-shrink-0 font-medium text-sm gap-x-1.5 px-2.5 py-1.5 shadow-sm text-white dark:text-gray-900 bg-primary-500 hover:bg-primary-600 disabled:bg-primary-500 dark:bg-primary-400 dark:hover:bg-primary-500 dark:disabled:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500 dark:focus-visible:outline-primary-400 inline-flex items-center">Apply</button>
            </form>`,
        formClassName: 'some-class',
        followOnDblClick: true,
        processVideoLink: false,
        processPastedLink: true,
        noFollowCheckbox: false,
        openInNewTabCheckbox: false,
        modeClassName: 'input', // 'select'
        /**
         * Allow multiple choises (to use with modeClassName="select")
         */
        selectMultipleClassName: true,
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
    {
        name: 'checkboxes',
        callback: (editor: IJodit) => {
            editor.events.on('afterInit', () => setupCheckboxes(editor.editor));

            // Listen to first change event (it should be the document contents) and stop listening
            function changeListener(): void {
                setupCheckboxes(editor.editor);
            }
            editor.events.one('change', changeListener);
        },
    },
];

const extraButtons = [
    '|',
    {
        name: 'insertCheckbox',
        iconURL: '/images/components/partials/doceditor/format-list-checkbox.svg',
        exec: (editor: IJodit) => {
            const label = document.createElement('label');
            label.setAttribute('contenteditable', 'false');
            const empty = document.createElement('span');
            empty.innerHTML = '&nbsp;';

            const input = document.createElement('input');
            input.setAttribute('type', 'checkbox');
            input.setAttribute('checked', 'true');
            input.onchange = setCheckboxState;

            label.appendChild(input);
            label.appendChild(empty);

            editor.s.insertHTML(label, true);
        },
    },
];

function setupCheckboxes(editor: HTMLElement): void {
    editor.querySelectorAll('.jodit-wysiwyg input[type=checkbox]').forEach((el) => {
        (el as HTMLInputElement).onchange = setCheckboxState;
    });
}

function setCheckboxState(event: Event): void {
    if (event.target === null) {
        return;
    }
    const target = event.target as HTMLInputElement;

    const attr = target.getAttribute('checked');
    const checked = attr !== null ? Boolean(attr) : false;
    if (checked) {
        target.removeAttribute('checked');
    } else {
        target.setAttribute('checked', 'true');
    }
}

watch(
    props,
    () => {
        if (props.modelValue !== '') {
            useTimeoutFn(() => setupCheckboxes(editorRef.value.editor.editor as HTMLElement), 75);
        }
    },
    { once: true },
);

watch(props, () => {
    editorRef.value.editor.setReadOnly(props.disabled);
});
</script>

<template>
    <div class="documentEditor mx-auto max-w-screen-xl">
        <JoditEditor ref="editorRef" v-model="content" :config="config" :plugins="plugins" :extra-buttons="extraButtons" />
    </div>
</template>

<style scoped>
.documentEditor:deep(.jodit-container) {
    .jodit-wysiwyg {
        width: 100% !important;
        max-width: 100%;

        * {
            margin-top: 4px;
            margin-bottom: 4px;
        }
    }
}
</style>
