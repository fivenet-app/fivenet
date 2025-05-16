<script lang="ts" setup>
import CharSexBadge from '~/components/partials/citizens/CharSexBadge.vue';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import { useAuthStore } from '~/stores/auth';
import { fromSecondsToFormattedDuration } from '~/utils/time';
import type { User } from '~~/gen/ts/resources/users/users';

const authStore = useAuthStore();

const { lastCharID } = storeToRefs(authStore);

const props = withDefaults(
    defineProps<{
        char: User;
        unavailable?: boolean;
        canSubmit?: boolean;
    }>(),
    {
        unavailable: false,
        canSubmit: true,
    },
);

const emit = defineEmits<{
    (e: 'selected', id: number): void;
}>();

function selectChar(): void {
    emit('selected', props.char.userId);
}

const { game } = useAppConfig();
</script>

<template>
    <UCard class="mx-4 flex w-full min-w-[28rem] max-w-md flex-col">
        <template #header>
            <div class="flex flex-col">
                <div class="mx-auto inline-flex items-center gap-2">
                    <ProfilePictureImg :src="char.avatar?.url" :name="`${char.firstname} ${char.lastname}`" :no-blur="true" />

                    <h2 class="text-center text-2xl font-semibold" @click="selectChar">
                        {{ char.firstname }} {{ char.lastname }}
                    </h2>
                </div>
            </div>
        </template>

        <dl class="flex grow flex-col justify-between text-center">
            <dd class="mb-1 flex items-center justify-center gap-2">
                <CharSexBadge :sex="char.sex ?? 'f'" />

                <UBadge v-if="lastCharID === char.userId" class="flex-initial" size="md" variant="subtle">
                    {{ $t('common.last_used') }}
                </UBadge>
            </dd>

            <dt class="text-sm font-semibold">
                {{ $t('common.job') }}
            </dt>
            <dd class="text-sm">
                {{ char.jobLabel }}<template v-if="char.job !== game.unemployedJobName"> ({{ char.jobGradeLabel }})</template>
            </dd>

            <dt class="text-sm font-semibold">
                {{ $t('common.date_of_birth') }}
            </dt>
            <dd class="text-sm">{{ char.dateofbirth }}</dd>

            <dt class="text-sm font-semibold">{{ $t('common.height') }}</dt>
            <dd class="text-sm">{{ char.height }}cm</dd>

            <template v-if="char.visum">
                <dt class="text-sm font-semibold">{{ $t('common.visum') }}</dt>
                <dd class="text-sm">{{ char.visum }}</dd>
            </template>

            <template v-if="char.playtime">
                <dt class="text-sm font-semibold">
                    {{ $t('common.playtime') }}
                </dt>
                <dd class="truncate text-sm">
                    {{ fromSecondsToFormattedDuration(char.playtime!) }}
                </dd>
            </template>
        </dl>

        <template #footer>
            <UButton
                class="inline-flex items-center"
                block
                :disabled="unavailable || !canSubmit"
                :loading="!canSubmit"
                :icon="unavailable ? 'i-mdi-lock' : undefined"
                :label="$t(!unavailable ? 'common.choose' : 'components.auth.CharacterSelectorCard.disabled_char')"
                @click="selectChar"
            />
        </template>
    </UCard>
</template>
