<script lang="ts" setup>
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import type { User } from '~~/gen/ts/resources/users/users';
defineProps<{
    user: User;
}>();

defineEmits<{
    (e: 'toggle-actions'): void;
}>();

const { game } = useAppConfig();
</script>

<template>
    <div class="flex items-center gap-2 p-4">
        <ProfilePictureImg
            :src="user?.props?.mugshot?.filePath"
            :name="`${user.firstname} ${user.lastname}`"
            :alt="$t('common.mugshot')"
            :enable-popup="true"
            size="3xl"
        />

        <div class="flex-1">
            <div class="flex snap-x flex-row flex-wrap justify-between gap-2 overflow-x-auto">
                <h2 class="flex-1 px-0.5 py-1 text-4xl font-bold break-words sm:pl-1">
                    {{ user?.firstname }} {{ user?.lastname }}
                </h2>
            </div>

            <div class="inline-flex flex-col gap-2 lg:flex-row">
                <UBadge>
                    {{ user.jobLabel }}
                    <template v-if="user.job !== game.unemployedJobName">
                        ({{ $t('common.rank') }}: {{ user.jobGradeLabel }})
                    </template>
                    {{ user.props?.jobName || user.props?.jobGradeNumber ? '*' : '' }}
                </UBadge>

                <UBadge v-if="user?.props?.wanted" color="error">
                    {{ $t('common.wanted').toUpperCase() }}
                </UBadge>
            </div>
        </div>

        <div class="flex flex-col gap-1 sm:flex-row">
            <UButton :label="$t('common.action', 2)" class="lg:hidden" icon="i-mdi-menu" @click="$emit('toggle-actions')" />
        </div>
    </div>
</template>
