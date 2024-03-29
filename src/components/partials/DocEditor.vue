<script lang="ts" setup>
import { useTimeoutFn, useVModel } from '@vueuse/core';
import 'jodit/es5/jodit.min.css';
import { Jodit } from 'jodit';
// @ts-ignore jodit-vue has no (detected) types
import { JoditEditor } from 'jodit-vue';
import type { IJodit } from 'jodit/types/types';
import { useSettingsStore } from '~/store/settings';

const props = defineProps<{
    modelValue: string;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', content: string): void;
}>();

const content = useVModel(props, 'modelValue', emit);

const settingsStore = useSettingsStore();
const { documents } = storeToRefs(settingsStore);

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
        iconURL: '/images/partials/doceditor/format-list-checkbox.svg',
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

watch(
    props,
    () => {
        if (props.modelValue !== '') {
            useTimeoutFn(setupCheckboxes, 50);
        }
    },
    {
        once: true,
    },
);

onBeforeUnmount(() => {
    // Remove event listeners on unmount
    const checkboxes: NodeListOf<HTMLInputElement> = document.querySelectorAll('.jodit-wysiwyg input[type=checkbox]');
    checkboxes.forEach((el) => (el.onchange = null));
});
</script>

<template>
    <JoditEditor ref="editorRef" v-model="content" :config="config" :plugins="plugins" :extra-buttons="extraButtons" />
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
