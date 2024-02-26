<script lang="ts" setup>
import { EyeIcon } from 'mdi-vue3';
import AvatarImg from '~/components/partials/citizens/AvatarImg.vue';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import { User } from '~~/gen/ts/resources/users/users';

defineProps<{
    user: User;
}>();
</script>

<template>
    <tr :key="user.userId" class="transition-colors even:bg-base-800 hover:bg-neutral/5">
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-base font-medium text-neutral sm:pl-1">
            <AvatarImg :url="user.avatar?.url" :name="`${user.firstname} ${user.lastname}`" size="sm" :rounded="false" />
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-base font-medium text-neutral sm:pl-1">
            {{ user.firstname }} {{ user.lastname }}
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-accent-200">{{ user.jobGradeLabel }} ({{ user.jobGrade }})</td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-accent-200">
            <PhoneNumberBlock :number="user.phoneNumber" />
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-accent-200">
            {{ user.dateofbirth }}
        </td>
        <td v-if="can('SuperUser')" class="whitespace-nowrap px-1 py-1 text-left text-accent-200">
            <div class="flex flex-row justify-end">
                <NuxtLink
                    :to="{
                        name: 'jobs-colleagues-id',
                        params: { id: user.userId ?? 0 },
                    }"
                    class="flex-initial text-primary-500 hover:text-primary-400"
                >
                    <EyeIcon class="ml-auto mr-2.5 w-5 h-auto" />
                </NuxtLink>
            </div>
        </td>
    </tr>
</template>
