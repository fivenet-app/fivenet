<script lang="ts" setup>
import { isFuture } from 'date-fns';
import EmailInfoPopover from '~/components/mailer/EmailInfoPopover.vue';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues/colleagues';
import ColleagueName from './ColleagueName.vue';

const props = withDefaults(
    defineProps<{
        colleague?: Colleague;
        userId?: number;
        active?: boolean;
        showContact?: boolean;
        showLabels?: boolean;
        showAbsence?: boolean;
        compact?: boolean;
    }>(),
    {
        colleague: undefined,
        userId: undefined,
        active: false,
        showContact: true,
        showLabels: true,
        showAbsence: true,
        compact: false,
    },
);

const emit = defineEmits<{
    label: [label: NonNullable<NonNullable<Colleague['props']>['labels']>['list'][number]];
}>();

const { t } = useI18n();
const { game } = useAppConfig();

const displayName = computed(() =>
    props.colleague ? `${props.colleague.firstname} ${props.colleague.lastname}` : `${t('common.id')}: ${props.userId ?? 0}`,
);

const gradeLabel = computed(() => {
    if (!props.colleague) return undefined;

    const label = props.colleague.jobGradeLabel || t('common.rank');
    if (props.colleague.job === game.unemployedJobName) return label;

    return `${label} (${props.colleague.jobGrade})`;
});
</script>

<template>
    <UCard
        :class="active ? 'ring-2 ring-primary' : ''"
        :ui="{
            root: 'flex flex-col',
            body: compact ? 'p-3 sm:p-3 flex-1' : 'p-4 sm:p-4 flex-1',
            footer: 'w-full',
        }"
    >
        <div class="flex min-w-0 gap-2" :class="compact ? 'items-start' : 'flex-col items-center sm:items-start'">
            <div class="flex flex-col items-center justify-center overflow-hidden" :class="compact ? '' : 'w-full'">
                <ProfilePictureImg
                    :src="colleague?.profilePicture"
                    :name="displayName"
                    :size="compact ? 'md' : '3xl'"
                    :enable-popup="!!colleague"
                    :alt="$t('common.profile_picture')"
                    :img-class="compact ? '' : 'size-42'"
                />
            </div>

            <div class="min-w-0 flex-1" :class="compact ? 'grid grid-cols-2 gap-2' : 'flex w-full flex-col'">
                <div class="flex min-w-0 flex-col gap-0.5">
                    <div class="flex min-w-0 flex-row flex-wrap gap-2">
                        <ColleagueName
                            v-if="colleague"
                            class="min-w-0 truncate font-medium text-highlighted"
                            :colleague="colleague"
                        />
                        <span v-else class="min-w-0 truncate font-medium text-highlighted">{{ displayName }}</span>

                        <slot name="badges" />
                    </div>

                    <div v-if="gradeLabel" class="truncate text-sm text-muted">{{ gradeLabel }}</div>
                </div>

                <div
                    class="mt-1 flex min-w-0 flex-col gap-1 overflow-x-hidden text-sm text-muted"
                    :class="compact ? '' : 'w-full'"
                >
                    <template v-if="showContact && colleague">
                        <USeparator v-if="!compact" class="my-0.5" />

                        <PhoneNumberBlock :number="colleague.phoneNumber" />

                        <div v-if="colleague.dateofbirth" class="inline-flex min-w-0 items-center gap-1">
                            <UIcon class="size-4 shrink-0" name="i-mdi-birthday-cake" />
                            <span class="truncate">{{ colleague.dateofbirth }}</span>
                        </div>

                        <div v-if="colleague.email" class="inline-flex min-w-0 items-center gap-1">
                            <UIcon class="size-4 shrink-0" name="i-mdi-email" />
                            <EmailInfoPopover
                                :email="colleague.email"
                                variant="link"
                                :trailing="false"
                                :ui="{ base: 'px-1 sm:px-1 py-0 sm:py-0' }"
                            />
                        </div>
                    </template>

                    <div v-if="showLabels" class="flex min-w-0 items-start gap-1 overflow-x-hidden">
                        <UIcon class="size-4 shrink-0 self-start" name="i-mdi-tag" />

                        <span v-if="!colleague?.props?.labels || !colleague?.props.labels.list.length">
                            {{ $t('common.none', [$t('common.label', 2)]) }}
                        </span>
                        <div v-else class="flex min-w-0 flex-1 basis-0 flex-row flex-wrap gap-1 overflow-hidden">
                            <UButton
                                v-for="label in colleague.props.labels.list"
                                :key="label.name"
                                class="max-w-full min-w-0 cursor-pointer overflow-hidden"
                                :class="isColorBright(hexToRgb(label.color, rgbBlack)!) ? 'text-black!' : 'text-white!'"
                                size="xs"
                                :label="label.name"
                                :icon="
                                    label.icon && label.icon !== '' ? convertComponentIconNameToDynamic(label.icon) : undefined
                                "
                                :style="{ backgroundColor: label.color }"
                                :ui="{ label: 'block min-w-0 truncate' }"
                                @click="emit('label', label)"
                            />
                        </div>
                    </div>

                    <div
                        v-if="showAbsence && colleague?.props?.absenceEnd && isFuture(toDate(colleague.props.absenceEnd))"
                        class="inline-flex min-w-0 items-center gap-1"
                    >
                        <UIcon class="size-4 shrink-0" name="i-mdi-island" />
                        <GenericTime :value="colleague.props.absenceBegin" type="shortDate" />
                        <span>{{ $t('common.to') }}</span>
                        <GenericTime :value="colleague.props.absenceEnd" type="date" />
                    </div>

                    <slot />
                </div>
            </div>
        </div>

        <template v-if="$slots.footer" #footer>
            <slot name="footer" />
        </template>
    </UCard>
</template>
