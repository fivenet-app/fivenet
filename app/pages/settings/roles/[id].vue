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

withDefaults(
    defineProps<{
        roleCount?: number;
    }>(),
    {
        roleCount: 0,
    },
);

defineEmits<{
    (e: 'deleted'): void;
}>();

const roleId = useRoute('settings-roles-id').params.id;
</script>

<template>
    <div v-if="!roleId" class="p-4 sm:p-6">
        <DataNoDataBlock icon="i-mdi-select" :message="$t('common.none_selected', [$t('common.role', 1)])" :padded="false" />
    </div>
    <RoleView
        v-else
        :role-id="parseInt(roleId)"
        :role-count="roleCount"
        @deleted="
            async () => {
                await navigateTo('/settings/roles');
                $emit('deleted');
            }
        "
    />
</template>
