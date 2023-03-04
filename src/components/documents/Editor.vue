<script lang="ts">
import { defineComponent } from 'vue';
import CKEditor from '@ckeditor/ckeditor5-vue';
import ClassicEditor from '@ckeditor/ckeditor5-build-classic';

export default defineComponent({
    components: {
        ckeditor: CKEditor.component
    },
    data() {
        return {
            editor: ClassicEditor,
            editorData: '<p>Content of the editor.</p>',
            // The configuration of the editor.
            editorConfig: {
                mention: {
                    feeds: [
                        {
                            marker: '@',
                            feed: ['@Barney', '@Lily', '@Marry Ann', '@Marshall', '@Robin', '@Ted'],
                            minimumCharacters: 1
                        },
                    ],
                },
            },
        };
    },
});
</script>

<route lang="json">
{
    "name": "documents-new",
    "meta": {
        "requiresAuth": true
    }
}
</route>

<template>
    <ckeditor :editor="editor" v-model="editorData" :config="editorConfig"></ckeditor>
</template>
