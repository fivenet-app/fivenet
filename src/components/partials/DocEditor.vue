<script lang="ts" setup>
import 'jodit/es5/jodit.min.css';
import { Jodit } from 'jodit';
// @ts-ignore jodit-vue has no (detected) types
import { JoditEditor } from 'jodit-vue';
import type { IJodit } from 'jodit/types/types';
import { useSettingsStore } from '~/store/settings';

const props = withDefaults(
    defineProps<{
        modelValue: string;
        disabled?: boolean;
    }>(),
    {
        disabled: false,
    },
);

const emit = defineEmits<{
    (e: 'update:modelValue', content: string): void;
}>();

const editorRef = ref<JoditEditor | null>(null);

const content = useVModel(props, 'modelValue', emit);

const settingsStore = useSettingsStore();
const { design: theme } = storeToRefs(settingsStore);

const config = {
    readOnly: props.disabled,

    language: 'de',
    spellcheck: true,
    minHeight: 475,
    editorClassName: 'prose' + (theme.value.documents.editorTheme === 'dark' ? ' prose-neutral' : ' prose-gray'),
    theme: theme.value.documents.editorTheme,

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
    <div class="documentEditor">
        <JoditEditor ref="editorRef" v-model="content" :config="config" :plugins="plugins" :extra-buttons="extraButtons" />
    </div>
</template>

<style scoped>
.documentEditor:deep(.jodit-wysiwyg) {
    min-width: 100%;

    * {
        margin-top: 4px;
        margin-bottom: 4px;
    }
}
</style>
