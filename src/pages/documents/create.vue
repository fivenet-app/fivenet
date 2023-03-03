<script lang="ts">
import { defineComponent } from 'vue';
import Editor from '../../components/Editor.vue';

export default defineComponent({
    name: 'app',
    components: {
        Editor
    },
});
</script>

<route lang="json">
{
    "name": "documents-new",
    "meta": {
        "requiresAuth": true,
        "permission": "documents-createdocument"
    }
}
</route>

<template>
    <Navbar />
    <NavPageHeader title="Documents" />
    <ContentWrapper>
        <Editor />
    </ContentWrapper>
    <Footer />
</template>
