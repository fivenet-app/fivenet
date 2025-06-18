<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import RoleView from '~/components/settings/roles/RoleView.vue';

definePageMeta({
    requiresAuth: true,
    permission: 'settings.SettingsService/GetRoles',
    validate: async (route) => {
        route = route as TypedRouteFromName<'settings-roles-id'>;
        // Check if the id is made up of digits
        if (typeof route.params.id !== 'string') {
            return false;
        }
        return !!(route.params.id && !isNaN(Number(route.params.id))) && Number(route.params.id) > -1;
    },
});

defineEmits<{
    (e: 'deleted'): void;
}>();

const roleId = useRoute('settings-roles-id').params.id;
</script>

<template>
    <div>
        <DataNoDataBlock v-if="!roleId" icon="i-mdi-select" :message="$t('common.none_selected', [$t('common.role', 1)])" />
        <RoleView
            v-else
            :role-id="parseInt(roleId)"
            @deleted="
                navigateTo({ name: 'settings-roles' });
                $emit('deleted');
            "
        />
    </div>
</template>
