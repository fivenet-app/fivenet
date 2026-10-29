<script lang="ts" setup>
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import ProfilePropertyBadges from '~/components/partials/properties/ProfilePropertyBadges.vue';
import type { User } from '~~/gen/ts/resources/users/user';
defineProps<{
    user: User;
}>();

defineEmits<{
    (e: 'toggle-actions'): void;
}>();
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

            <ProfilePropertyBadges
                :properties="user.properties?.properties"
                :user="user"
                :key-order="['user.job', 'user.wanted']"
            />
        </div>

        <div class="flex shrink-0 flex-col gap-1 sm:flex-row">
            <UButton
                data-tour="citizen-profile-actions"
                class="lg:hidden"
                :label="$t('common.action', 2)"
                icon="i-mdi-menu"
                @click="$emit('toggle-actions')"
            />
        </div>
    </div>
</template>
