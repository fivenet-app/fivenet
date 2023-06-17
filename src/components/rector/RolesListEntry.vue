<script lang="ts" setup>
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiEye } from '@mdi/js';
import { RoutesNamedLocations } from '~~/.nuxt/typed-router/__routes';
import { Role } from '~~/gen/ts/resources/permissions/permissions';

const props = defineProps<{
    role: Role;
    to?: RoutesNamedLocations;
}>();

if (props.to && (props.to.name === 'rector-limiter-id' || props.to.name === 'rector-roles-id')) {
    props.to.params!.id = props.role.id.toString();
}
</script>

<template>
    <tr :key="role.id?.toString()">
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ role.jobLabel }} - {{ role.jobGradeLabel }}
        </td>
        <td class="whitespace-nowrap py-2 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
            <div class="flex flex-row justify-end">
                <NuxtLink
                    :to="to ?? { name: 'rector-roles-id', params: { id: role.id.toString() } }"
                    class="flex-initial text-primary-500 hover:text-primary-400"
                >
                    <SvgIcon class="h-6 w-6 text-primary-500" aria-hidden="true" type="mdi" :path="mdiEye" />
                </NuxtLink>
            </div>
        </td>
    </tr>
</template>
