<script lang="ts" setup>
import { useThrottleFn } from '@vueuse/core';
import { LoadingIcon } from 'mdi-vue3';
import CharSexBadge from '~/components/partials/citizens/CharSexBadge.vue';
import { useAuthStore } from '~/store/auth';
import { fromSecondsToFormattedDuration } from '~/utils/time';
import { User } from '~~/gen/ts/resources/users/users';
import AvatarImg from '~/components/partials/citizens/AvatarImg.vue';

const authStore = useAuthStore();

const { lastCharID } = storeToRefs(authStore);
const { chooseCharacter } = authStore;

const props = defineProps<{
    char: User;
}>();

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (_) => {
    canSubmit.value = false;
    await chooseCharacter(props.char.userId).finally(() => setTimeout(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <div :key="char.userId" class="flex flex-col content-end rounded-lg bg-base-800 shadow-float">
        <div>
            <div class="flex">
                <div class="mt-4 mx-auto flex flex-row items-center gap-3">
                    <AvatarImg :url="char.avatar?.url" :name="`${char.firstname} ${char.lastname}`" />

                    <h2 class="text-center text-2xl font-medium text-neutral">{{ char.firstname }} {{ char.lastname }}</h2>
                </div>
            </div>
        </div>
        <div class="flex flex-1 flex-col p-4">
            <dl class="flex flex-grow flex-col justify-between text-center">
                <dd class="mb-2 inline-flex items-center justify-center gap-1">
                    <CharSexBadge :sex="char.sex!" />
                    <span
                        class="inline-flex items-center rounded-full bg-gray-100 px-2.5 py-0.5 text-sm font-medium text-gray-800"
                        >{{ char.jobLabel }} ({{ $t('common.rank') }}: {{ char.jobGradeLabel }})</span
                    >
                    <span
                        v-if="lastCharID === char.userId"
                        class="rounded-full bg-success-100 px-3 py-0.5 text-center text-sm font-medium text-success-800"
                    >
                        {{ $t('common.last_used') }}
                    </span>
                </dd>
                <dt class="text-sm font-medium text-neutral">
                    {{ $t('common.date_of_birth') }}
                </dt>
                <dd class="text-sm text-gray-300">{{ char.dateofbirth }}</dd>
                <dt class="text-sm font-medium text-neutral">{{ $t('common.height') }}</dt>
                <dd class="text-sm text-gray-300">{{ char.height }}cm</dd>
                <dt class="text-sm font-medium text-neutral">{{ $t('common.visum') }}</dt>
                <dd class="text-sm text-gray-300">{{ char.visum }}</dd>
                <dt class="text-sm font-medium text-neutral">
                    {{ $t('common.playtime') }}
                </dt>
                <dd class="text-sm text-gray-300 truncate">
                    {{ fromSecondsToFormattedDuration(char.playtime!) }}
                </dd>
            </dl>
        </div>
        <div>
            <div class="-mt-px flex">
                <div class="flex w-0 flex-1">
                    <button
                        type="button"
                        class="relative inline-flex w-0 flex-1 items-center justify-center gap-x-3 rounded-b-lg border border-transparent bg-base-700 py-4 text-sm font-semibold text-neutral transition-colors hover:bg-base-600"
                        :disabled="!canSubmit"
                        @click="onSubmitThrottle(char.userId)"
                    >
                        <template v-if="!canSubmit">
                            <LoadingIcon class="mr-2 h-5 w-5 animate-spin" />
                        </template>
                        {{ $t('common.choose') }}
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>
