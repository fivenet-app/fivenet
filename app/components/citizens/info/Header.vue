<script lang="ts" setup>
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import type { User } from '~~/gen/ts/resources/users/user';
defineProps<{
    user: User;
}>();

defineEmits<{
    (e: 'toggle-actions'): void;
}>();

const { game } = useAppConfig();
</script>

<template>
    <div class="flex min-w-0 flex-1 items-center gap-2">
        <ProfilePictureImg
            class="shrink-0"
            :src="user?.props?.mugshot?.filePath"
            :name="`${user.firstname} ${user.lastname}`"
            :alt="$t('common.mugshot')"
            enable-popup
            size="3xl"
        />

        <div class="min-w-0 flex-1">
            <div class="flex min-w-0 flex-row flex-wrap gap-2">
                <h1 class="min-w-0 flex-1 px-0.5 py-1 text-2xl font-bold break-words sm:text-3xl">
                    {{ user?.firstname }} {{ user?.lastname }}
                </h1>
            </div>

            <div class="flex flex-row flex-wrap gap-2">
                <UBadge>
                    {{ user.jobLabel }}
                    <template v-if="user.job !== game.unemployedJobName">
                        ({{ $t('common.rank') }}: {{ user.jobGradeLabel }})
                    </template>
                    {{ user.props?.jobName || user.props?.jobGradeNumber ? '*' : '' }}
                </UBadge>

                <UBadge v-if="user?.props?.wanted" color="error" :label="$t('common.wanted').toUpperCase()" />
            </div>
        </div>

        <div class="flex shrink-0 flex-col gap-1 sm:flex-row">
            <UButton class="lg:hidden" :label="$t('common.action', 2)" icon="i-mdi-menu" @click="$emit('toggle-actions')" />
        </div>
    </div>
</template>
