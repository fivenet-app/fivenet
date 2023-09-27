<script lang="ts" setup>
import { useThrottleFn } from '@vueuse/core';
import { LoadingIcon } from 'mdi-vue3';
import CharSexBadge from '~/components/citizens/CharSexBadge.vue';
import { useAuthStore } from '~/store/auth';
import { fromSecondsToFormattedDuration } from '~/utils/time';
import { User } from '~~/gen/ts/resources/users/users';

const authStore = useAuthStore();

const { lastCharID } = storeToRefs(authStore);
const { chooseCharacter } = authStore;

const props = defineProps<{
    char: User;
}>();

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (_) => {
    canSubmit.value = false;
    await chooseCharacter(props.char.userId).finally(() => setTimeout(() => (canSubmit.value = true), 350));
}, 1000);
</script>

<template>
    <div :key="char.userId" class="divide-y rounded-lg bg-base-800 shadow-float flex flex-col content-end">
        <div class="flex flex-col flex-1 p-8">
            <div class="flex flex-row items-center gap-3 mx-auto">
                <h2 class="text-2xl font-medium text-center text-neutral">{{ char.firstname }}, {{ char.lastname }}</h2>
                <CharSexBadge :sex="char.sex!" />
                <div v-if="lastCharID === char.userId">
                    <span
                        class="inline-flex items-center text-center rounded-full bg-success-100 px-3 py-0.5 text-sm font-medium text-success-800"
                    >
                        {{ $t('common.last_used') }}
                    </span>
                </div>
            </div>
            <dl class="flex flex-col justify-between flex-grow mt-2 text-center">
                <dd class="mt-2 mb-2">
                    <span
                        class="inline-flex items-center rounded-full bg-gray-100 px-2.5 py-0.5 text-sm font-medium text-gray-800"
                        >{{ char.jobLabel }} ({{ $t('common.rank') }}: {{ char.jobGradeLabel }})</span
                    >
                </dd>
                <dt class="text-sm text-neutral font-medium">
                    {{ $t('common.date_of_birth') }}
                </dt>
                <dd class="text-sm text-gray-300">{{ char.dateofbirth }}</dd>
                <dt class="text-sm text-neutral font-medium">{{ $t('common.height') }}</dt>
                <dd class="text-sm text-gray-300">{{ char.height }}cm</dd>
                <dt class="text-sm text-neutral font-medium">{{ $t('common.visum') }}</dt>
                <dd class="text-sm text-gray-300">{{ char.visum }}</dd>
                <dt class="text-sm text-neutral font-medium">
                    {{ $t('common.playtime') }}
                </dt>
                <dd class="text-sm text-gray-300">
                    {{ fromSecondsToFormattedDuration(char.playtime!) }}
                </dd>
            </dl>
        </div>
        <div>
            <div class="flex -mt-px">
                <div class="flex flex-1 w-0">
                    <button
                        type="button"
                        @click="onSubmitThrottle(char.userId)"
                        class="relative inline-flex items-center justify-center flex-1 w-0 py-4 text-sm font-semibold transition-colors border border-transparent rounded-b-lg gap-x-3 text-neutral bg-base-700 hover:bg-base-600"
                        :disabled="!canSubmit"
                    >
                        <template v-if="!canSubmit">
                            <LoadingIcon class="animate-spin h-5 w-5 mr-2" />
                        </template>
                        {{ $t('common.choose') }}
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>
