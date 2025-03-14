<script lang="ts" setup>
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import type { ClipboardUser } from '~/stores/clipboard';
import type { ClassProp } from '~/typings';
import type { User, UserShort } from '~~/gen/ts/resources/users/users';
import EmailBlock from './EmailBlock.vue';

const props = withDefaults(
    defineProps<{
        userId?: number;
        user?: ClipboardUser | User | UserShort;
        textClass?: ClassProp;
        showAvatar?: boolean;
        showAvatarInName?: boolean;
        trailing?: boolean;
        showBirthdate?: boolean;
    }>(),
    {
        userId: undefined,
        user: undefined,
        textClass: '',
        showAvatar: undefined,
        showAvatarInName: false,
        trailing: true,
        showBirthdate: false,
    },
);

const { $grpc } = useNuxtApp();

const { can, activeChar } = useAuth();

const { popover } = useAppConfig();

const userId = computed(() => props.userId ?? props.user?.userId ?? 0);

const {
    data,
    refresh,
    pending: loading,
    error,
} = useLazyAsyncData(`citizen-info-${userId.value}`, () => getCitizen(userId.value), { immediate: !props.user });

async function getCitizen(id: number): Promise<User | undefined> {
    try {
        const call = $grpc.citizenstore.citizenStore.getUser({
            userId: id,
            infoOnly: true,
        });
        const { response } = await call;

        if (response.user!.phoneNumber && props.user?.phoneNumber) {
            response.user!.phoneNumber = props.user.phoneNumber;
        }
        if (response.user!.avatar && props.user?.avatar) {
            response.user!.avatar = props.user.avatar;
        }

        return response.user!;
    } catch (_) {
        return undefined;
    }
}

const user = computed(() => ({ ...data.value, ...props.user }) as User);

const { game } = useAppConfig();

const opened = ref(false);
watchOnce(opened, async () => {
    if (props.user) {
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
    <UPopover v-else :ui="{ trigger: 'inline-flex w-auto', wrapper: 'inline-block' }">
        <UButton
            v-bind="$attrs"
            variant="link"
            :padded="false"
            class="inline-flex items-center gap-1 p-px"
            :trailing-icon="trailing ? 'i-mdi-chevron-down' : undefined"
            @click="opened = true"
        >
            <template v-if="showAvatarInName" #leading>
                <USkeleton v-if="!user && loading" class="h-6 w-6" :ui="{ rounded: 'rounded-full' }" />
                <ProfilePictureImg v-else :src="user?.avatar?.url" :name="`${user?.firstname} ${user?.lastname}`" size="3xs" />
            </template>

            <USkeleton v-if="!user && loading" class="h-8 w-[125px]" />
            <span v-else class="truncate" :class="textClass">
                <slot name="name" :user="user">
                    {{ user?.firstname }} {{ user?.lastname }}
                    <template v-if="showBirthdate && user.dateofbirth">({{ user.dateofbirth }})</template>
                </slot>
            </span>
        </UButton>

        <template #panel>
            <div class="flex flex-col gap-2 p-4">
                <div class="grid w-full grid-cols-3 gap-2">
                    <IDCopyBadge
                        :id="userId ?? user?.userId ?? 0"
                        prefix="CIT"
                        :title="{ key: 'notifications.citizen_info.copy_citizen_id.title', parameters: {} }"
                        :content="{ key: 'notifications.citizen_info.copy_citizen_id.content', parameters: {} }"
                        size="xs"
                    />

                    <UButton
                        v-if="can('CitizenStoreService.ListCitizens').value"
                        variant="link"
                        icon="i-mdi-account"
                        :to="{ name: 'citizens-id', params: { id: userId ?? user?.userId ?? 0 } }"
                    >
                        {{ $t('common.profile') }}
                    </UButton>

                    <UButton
                        v-if="can('JobsService.GetColleague').value && user?.job === activeChar?.job"
                        variant="link"
                        icon="i-mdi-briefcase"
                        :to="{ name: 'jobs-colleagues-id-info', params: { id: userId ?? user?.userId ?? 0 } }"
                    >
                        <span class="truncate">
                            {{ $t('common.colleague') }}
                        </span>
                    </UButton>

                    <PhoneNumberBlock
                        v-if="user?.phoneNumber"
                        :number="user.phoneNumber"
                        :hide-number="true"
                        :show-label="true"
                    />

                    <EmailBlock v-if="user?.props && user.props?.email" :email="user.props?.email" hide-email />
                </div>

                <div v-if="error">
                    <DataErrorBlock
                        :title="$t('common.unable_to_load', [$t('common.citizen', 2)])"
                        :error="error"
                        :retry="refresh"
                    />
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

                    <div class="flex flex-col gap-1 text-sm font-normal">
                        <p v-if="user.jobLabel">
                            <span class="font-semibold">{{ $t('common.job') }}:</span>
                            {{ user.jobLabel }}
                            <span v-if="(user.jobGrade ?? 0) > 0 && user.job !== game.unemployedJobName">
                                ({{ user.jobGradeLabel }})</span
                            >
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
