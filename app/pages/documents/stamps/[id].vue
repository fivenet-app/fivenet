<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
import { getDocumentsStampsClient } from '~~/gen/ts/clients';

useHead({
    title: 'pages.documents.stamps.update',
});

definePageMeta({
    title: 'pages.documents.stamps.update',
    requiresAuth: true,
    permission: 'documents.StampsService/UpsertStampPerm',
    validate: async (route) => {
        route = route as TypedRouteFromName<'documents-stamps-id'>;
        // Check if the id is made up of digits
        if (typeof route.params.id !== 'string') {
            return false;
        }
        return !!(route.params.id && !isNaN(Number(route.params.id))) && Number(route.params.id) > -1;
    },
});

const route = useRoute('documents-stamps-id');

const { data, status, error } = useLazyAsyncData('stamp', async () => {
    const stampsClient = await getDocumentsStampsClient();
    const { response } = await stampsClient.getStamp({
        id: Number.parseInt(route.params.id as string),
    });
    return response;
});
</script>

<template>
    <div>
        TODO

        {{ data }}
    </div>
</template>
