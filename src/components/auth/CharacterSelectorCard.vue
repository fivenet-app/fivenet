<script lang="ts" setup>
import { useThrottleFn, useTimeoutFn } from '@vueuse/core';
import CharSexBadge from '~/components/partials/citizens/CharSexBadge.vue';
import { useAuthStore } from '~/store/auth';
import { fromSecondsToFormattedDuration } from '~/utils/time';
import { User } from '~~/gen/ts/resources/users/users';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';

const authStore = useAuthStore();

const { lastCharID } = storeToRefs(authStore);
const { chooseCharacter } = authStore;

const props = defineProps<{
    char: User;
}>();

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (_) => {
    canSubmit.value = false;
    await chooseCharacter(props.char.userId).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UCard :key="char.userId" class="min-w-[28rem] mx-4 flex w-full max-w-md flex-col rounded-lg bg-base-800">
        <template #header>
            <div class="flex">
                <div class="mx-auto flex flex-row items-center gap-2">
                    <ProfilePictureImg :url="char.avatar?.url" :name="`${char.firstname} ${char.lastname}`" :no-blur="true" />

                    <h2 class="text-center text-2xl font-medium">{{ char.firstname }} {{ char.lastname }}</h2>
                </div>
            </div>
        </template>

        <dl class="flex grow flex-col justify-between text-center">
            <dd class="mb-2 inline-flex items-center justify-center gap-1">
                <CharSexBadge :sex="char.sex!" />
                <span
                    class="inline-flex items-center rounded-full bg-gray-100 px-2.5 py-0.5 text-sm font-medium text-gray-800 truncate"
                    >{{ char.jobLabel }} ({{ char.jobGradeLabel }})</span
                >
                <span
                    v-if="lastCharID === char.userId"
                    class="rounded-full bg-success-100 px-3 py-0.5 text-center text-sm font-medium text-success-800 truncate"
                >
                    {{ $t('common.last_used') }}
                </span>
            </dd>
            <dt class="text-sm font-medium">
                {{ $t('common.date_of_birth') }}
            </dt>
            <dd class="text-sm text-gray-300">{{ char.dateofbirth }}</dd>
            <dt class="text-sm font-medium">{{ $t('common.height') }}</dt>
            <dd class="text-sm text-gray-300">{{ char.height }}cm</dd>
            <template v-if="char.visum">
                <dt class="text-sm font-medium">{{ $t('common.visum') }}</dt>
                <dd class="text-sm text-gray-300">{{ char.visum }}</dd>
            </template>
            <template v-if="char.playtime">
                <dt class="text-sm font-medium">
                    {{ $t('common.playtime') }}
                </dt>
                <dd class="truncate text-sm text-gray-300">
                    {{ fromSecondsToFormattedDuration(char.playtime!) }}
                </dd>
            </template>
        </dl>

        <template #footer>
            <UButton
                block
                class="inline-flex items-center"
                :disabled="!canSubmit"
                :loading="!canSubmit"
                @click="onSubmitThrottle(char.userId)"
            >
                {{ $t('common.choose') }}
            </UButton>
        </template>
    </UCard>
</template>
