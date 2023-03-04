<script lang="ts">
import { defineComponent } from 'vue';
import { getDocumentsClient } from '../../grpc';


export default defineComponent({
    data() {
        return {
            client: getDocumentsClient(),
            loading: false,
        };
    },
    // TODO
});
</script>

<template>
    <div class="py-2">
        <div class="px-2 sm:px-6 lg:px-8">
            <h1>DOCUMENTS</h1>
        </div>
    </div>
</template>
