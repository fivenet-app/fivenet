<script lang="ts">
import { defineComponent } from 'vue';
import { Document } from '@arpanet/gen/documents/documents_pb';

export default defineComponent({
    data() {
        return {
            'activities': [] as Array<Document>,
        };
    },
    props: {
        userID: {
            required: true,
            type: Number,
        },
    },
});
</script>

<template>
    <!-- TODO -->
</template>
