<script lang="ts" setup>
import 'jodit/build/jodit.min.css';
import { JoditEditor } from 'jodit-vue';
import { Jodit } from 'jodit';
import ContentWrapper from '~/components/partials/ContentWrapper.vue';

const content = ref('');

const config = {
    language: 'de',
    spellcheck: true,
    readonly: false,
    defaultActionOnPaste: 'insert_clear_html',
    plugins: [],
    disablePlugins: ['about', 'print', 'classSpan', 'video'],

    // Uploader Plugin
    uploader: {
        insertImageAsBase64URI: true,
    },
    // Clean HTML Plugin
    cleanHTML: {
        denyTags: 'script',
    },

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
        processVideoLink: true,

        /**
         * Wrap inserted link
         */
        processPastedLink: true,

        /**
         * Show `no follow` checkbox in link dialog.
         */
        noFollowCheckbox: true,

        /**
         * Show `Open in new tab` checkbox in link dialog.
         */
        openInNewTabCheckbox: true,

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
</script>

<style>
html.h-full
    body.h-full.bg-body-color.overflow-hidden
    div.jodit.jodit-popup-container.jodit-box.jodit_theme_default
    div.jodit-popup.jodit-popup_strategy_rightbottom
    div.jodit-popup__content
    div.jodit-tabs
    div.jodit-tabs__buttons
    button.jodit-ui-button.jodit-ui-button_size_middle.jodit-ui-button_variant_initial.jodit-ui-button_link.jodit-ui-button_text-icons_true.jodit-tabs__button.jodit-tabs__button_columns_2 {
    display: none;
}
</style>

<template>
    <ContentWrapper>
        <div class="h-full w-full">
            <JoditEditor v-model="content" :config="config" />
        </div>
    </ContentWrapper>
</template>
