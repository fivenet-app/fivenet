<script lang="ts" setup>
import { Role } from '@fivenet/gen/resources/permissions/permissions_pb';
import { EyeIcon } from '@heroicons/vue/24/solid';

defineProps({
    role: {
        required: true,
        type: Role,
    },
});
</script>

<template>
    <tr :key="role.getId()">
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ role.getDescription() }}
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ $d(role.hasUpdatedAt() ? role.getUpdatedAt() : role.getCreatedAt()?.getTimestamp()?.toDate()!, 'short') }}
        </td>
        <td class="whitespace-nowrap py-2 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
            <div class="flex flex-row justify-end">
                <NuxtLink :to="{ name: 'rector-roles-id', params: { id: role.getId() } }"
                    class="flex-initial text-primary-500 hover:text-primary-400">
                    <EyeIcon class="h-6 w-6 text-primary-500" aria-hidden="true" />
                </NuxtLink>
            </div>
        </td>
    </tr>
</template>
