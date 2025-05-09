<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import AttrView from '~/components/rector/attrs/AttrView.vue';

definePageMeta({
    requiresAuth: true,
    permission: 'SuperUser',
    validate: async (route) => {
        route = route as TypedRouteFromName<'rector-limiter-id'>;
        // Check if the id is made up of digits
        if (typeof route.params.id !== 'string') {
            return false;
        }
        return !!(route.params.id && !isNaN(Number(route.params.id))) && Number(route.params.id) > -1;
    },
});

const roleId = useRoute('rector-roles-id').params.id;
</script>

<template>
    <div>
        <DataNoDataBlock v-if="!roleId" icon="i-mdi-select" :message="$t('common.none_selected', [$t('common.job', 2)])" />
        <AttrView v-else :role-id="parseInt(roleId)" @deleted="navigateTo({ name: 'rector-limiter' })" />
    </div>
</template>
