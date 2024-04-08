<script lang="ts" setup>
import { type User, type UserShort } from '~~/gen/ts/resources/users/users';
import { ClipboardUser } from '~/store/clipboard';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';

withDefaults(
    defineProps<{
        user: ClipboardUser | User | UserShort | undefined;
        noPopover?: boolean;
        textClass?: unknown;
        buttonClass?: unknown;
        showAvatar?: boolean;
        trailing?: boolean;
    }>(),
    {
        textClass: '' as any,
        buttonClass: '' as any,
        showAvatar: undefined,
        trailing: true,
    },
);
</script>

<template>
    <template v-if="!user">
        <span class="inline-flex items-center">
            <slot name="before" />
            <span>N/A</span>
            <slot name="after" />
        </span>
    </template>
    <template v-else-if="noPopover">
        <span class="inline-flex items-center">
            <slot name="before" />
            <UButton variant="link" :padded="false" :to="{ name: 'citizens-id', params: { id: user.userId ?? 0 } }">
                {{ user.firstname }} {{ user.lastname }}
            </UButton>
            <span v-if="user.phoneNumber">
                <PhoneNumberBlock v-if="user.phoneNumber" :number="user.phoneNumber" :hide-number="true" :show-label="false" />
            </span>
            <slot name="after" />
        </span>
    </template>
    <UPopover v-else>
        <UButton
            variant="link"
            :padded="false"
            class="inline-flex items-center"
            :class="buttonClass"
            :trailing-icon="trailing ? 'i-mdi-chevron-down' : undefined"
        >
            <slot name="before" />
            <span class="truncate" :class="textClass"> {{ user.firstname }} {{ user.lastname }} </span>
            <slot name="after" />
        </UButton>

        <template #panel>
            <div class="flex flex-col gap-2 p-4">
                <UButtonGroup v-if="can('CitizenStoreService.ListCitizens') || user.phoneNumber" class="inline-flex w-full">
                    <UButton
                        v-if="can('CitizenStoreService.ListCitizens')"
                        variant="link"
                        icon="i-mdi-account"
                        :to="{ name: 'citizens-id', params: { id: user.userId ?? 0 } }"
                    >
                        {{ $t('common.profile') }}
                    </UButton>

                    <PhoneNumberBlock
                        v-if="user.phoneNumber"
                        :number="user.phoneNumber"
                        :hide-number="true"
                        :show-label="true"
                    />
                </UButtonGroup>

                <div class="inline-flex flex-row gap-2">
                    <div v-if="showAvatar === undefined || showAvatar">
                        <ProfilePictureImg :url="user.avatar?.url" :name="`${user.firstname} ${user.lastname}`" />
                    </div>
                    <div>
                        <UButton variant="link" :padded="false" :to="{ name: 'citizens-id', params: { id: user.userId ?? 0 } }">
                            {{ user.firstname }} {{ user.lastname }}
                        </UButton>

                        <p v-if="user.jobLabel" class="text-sm font-normal">
                            <span class="font-semibold">{{ $t('common.job') }}:</span>
                            {{ user.jobLabel }}
                            <span v-if="(user.jobGrade ?? 0) > 0 && user.jobGradeLabel"> ({{ user.jobGradeLabel }})</span>
                        </p>

                        <p v-if="user.dateofbirth" class="text-sm font-normal">
                            <span class="font-semibold">{{ $t('common.date_of_birth') }}:</span>
                            {{ user.dateofbirth }}
                        </p>
                    </div>
                </div>
            </div>
        </template>
    </UPopover>
</template>
