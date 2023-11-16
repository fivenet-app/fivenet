<script lang="ts" setup>
import { EyeIcon } from 'mdi-vue3';
import { Role } from '~~/gen/ts/resources/permissions/permissions';

defineProps<{
    role: Role;
}>();

defineEmits<{
    (e: 'selected', role: Role): void;
}>();
</script>

<template>
    <tr :key="role.id">
        <td class="whitespace-nowrap py-2 pl-2 pr-3 text-sm font-medium text-neutral" :title="`ID: ${role.id}`">
            {{ role.jobLabel }} - {{ role.jobGradeLabel }} ({{ role.grade }})
        </td>
        <td class="whitespace-nowrap py-2 pl-3 pr-2 text-right text-sm font-medium">
            <div class="flex flex-row justify-end">
                <button
                    type="button"
                    class="flex-initial text-primary-500 hover:text-primary-400"
                    @click="$emit('selected', role)"
                >
                    <EyeIcon class="h-6 w-6 text-primary-500" aria-hidden="true" />
                </button>
            </div>
        </td>
    </tr>
</template>
