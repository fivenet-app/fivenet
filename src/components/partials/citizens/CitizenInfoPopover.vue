<script lang="ts" setup>
import { type User, type UserShort } from '~~/gen/ts/resources/users/users';
import { ClipboardUser } from '~/store/clipboard';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import IDCopyBadge from '../IDCopyBadge.vue';
import DataErrorBlock from '../data/DataErrorBlock.vue';

const props = withDefaults(
    defineProps<{
        userId?: number;
        user?: ClipboardUser | User | UserShort | undefined;
        textClass?: unknown;
        showAvatar?: boolean;
        showAvatarInName?: boolean;
        trailing?: boolean;
    }>(),
    {
        textClass: '' as any,
        showAvatar: undefined,
        showAvatarInName: false,
        trailing: true,
    },
);

const { popover } = useAppConfig();

const {
    data,
    refresh,
    pending: loading,
    error,
} = useLazyAsyncData(
    `citizen-info-${props.userId ?? props.user?.userId}`,
    () => getCitizen(props.userId ?? props.user?.userId ?? 0),
    { immediate: false },
);

async function getCitizen(id: number): Promise<User> {
    try {
        const call = getGRPCCitizenStoreClient().getUser({
            userId: id,
            infoOnly: true,
        });
        const { response } = await call;

        return response.user!;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const user = computed(() => data.value || props.user);

const opened = ref(false);
watchOnce(opened, async () => {
    if (!props.user) {
        refresh();
    } else {
        useTimeoutFn(async () => refresh(), popover.waitTime);
    }
});
</script>

<template>
    <template v-if="!user && !userId">
        <span class="inline-flex items-center">
            {{ $t('common.na') }}
        </span>
    </template>
    <UPopover v-else>
        <UButton
            v-bind="$attrs"
            variant="link"
            :padded="false"
            class="inline-flex items-center gap-1 p-px"
            :trailing-icon="trailing ? 'i-mdi-chevron-down' : undefined"
            @click="opened = true"
        >
            <template #leading v-if="showAvatarInName">
                <USkeleton v-if="!user && loading" class="h-6 w-6" :ui="{ rounded: 'rounded-full' }" />
                <ProfilePictureImg v-else :src="user?.avatar?.url" :name="`${user?.firstname} ${user?.lastname}`" size="3xs" />
            </template>

            <USkeleton v-if="!user && loading" class="h-8 w-[125px]" />
            <span v-else class="truncate" :class="textClass"> {{ user?.firstname }} {{ user?.lastname }} </span>
        </UButton>

        <template #panel>
            <div class="flex flex-col gap-2 p-4">
                <div class="inline-flex w-full gap-1">
                    <IDCopyBadge
                        :id="userId ?? user?.userId ?? 0"
                        prefix="CIT"
                        :title="{ key: 'notifications.document_view.copy_document_id.title', parameters: {} }"
                        :content="{ key: 'notifications.document_view.copy_document_id.content', parameters: {} }"
                        size="xs"
                        variant="link"
                    />

                    <UButton
                        v-if="can('CitizenStoreService.ListCitizens')"
                        variant="link"
                        icon="i-mdi-account"
                        :to="{ name: 'citizens-id', params: { id: userId ?? user?.userId ?? 0 } }"
                    >
                        {{ $t('common.profile') }}
                    </UButton>

                    <PhoneNumberBlock
                        v-if="user?.phoneNumber"
                        :number="user.phoneNumber"
                        :hide-number="true"
                        :show-label="true"
                    />
                </div>

                <div v-if="error">
                    <DataErrorBlock :title="$t('common.unable_to_load', [$t('common.citizen', 2)])" :retry="refresh" />
                </div>

                <div v-else-if="loading && !user" class="flex flex-col gap-2 text-gray-900 dark:text-white">
                    <USkeleton class="h-8 w-[250px]" />

                    <div class="flex flex-row items-center gap-2">
                        <USkeleton class="h-7 w-[60px]" />
                        <USkeleton class="h-6 w-[215px]" />
                    </div>
                </div>

                <div v-else-if="user" class="flex flex-col gap-2 text-gray-900 dark:text-white">
                    <div class="inline-flex flex-row gap-2">
                        <ProfilePictureImg
                            v-if="showAvatar === undefined || showAvatar"
                            :src="user.avatar?.url"
                            :name="`${user.firstname} ${user.lastname}`"
                        />

                        <UButton variant="link" :padded="false" :to="{ name: 'citizens-id', params: { id: user.userId ?? 0 } }">
                            <span>{{ user.firstname }} {{ user.lastname }}</span>
                        </UButton>
                    </div>

                    <div class="flex flex-col gap-1">
                        <p v-if="user.jobLabel" class="text-sm font-normal">
                            <span class="font-semibold">{{ $t('common.job') }}:</span>
                            {{ user.jobLabel }}
                            <span v-if="(user.jobGrade ?? 0) > 0 && user.jobGradeLabel"> ({{ user.jobGradeLabel }})</span>
                        </p>

                        <p v-if="user.dateofbirth">
                            <span class="font-semibold">{{ $t('common.date_of_birth') }}:</span>
                            {{ user.dateofbirth }}
                        </p>
                    </div>
                </div>
            </div>
        </template>
    </UPopover>
</template>
