<script lang="ts" setup>
import { EyeIcon } from 'mdi-vue3';
import { Role } from '~~/gen/ts/resources/permissions/permissions';

defineProps<{
    role: Role;
}>();
</script>

<template>
    <tr :key="role.id?.toString()">
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ role.jobLabel }} - {{ role.jobGradeLabel }} ({{ role.grade }})
        </td>
        <td class="whitespace-nowrap py-2 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
            <div class="flex flex-row justify-end">
                <NuxtLink
                    :to="{ name: 'rector-roles-id', params: { id: role.id.toString() } }"
                    class="flex-initial text-primary-500 hover:text-primary-400"
                >
                    <EyeIcon class="h-6 w-6 text-primary-500" aria-hidden="true" />
                </NuxtLink>
            </div>
        </td>
    </tr>
</template>
