<script lang="ts" setup>
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import { User } from '~~/gen/ts/resources/users/users';

defineProps<{
    user: User;
}>();
</script>

<template>
    <tr :key="user.userId" class="transition-colors even:bg-base-800 hover:bg-neutral/5">
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-base font-medium text-neutral sm:pl-1">
            {{ user.firstname }}, {{ user.lastname }}
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200">{{ user.jobGradeLabel }} ({{ user.jobGrade }})</td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200">
            <PhoneNumberBlock :number="user.phoneNumber" />
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200">
            {{ user.dateofbirth }}
        </td>
    </tr>
</template>
