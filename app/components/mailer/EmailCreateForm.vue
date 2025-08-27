<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { enumToAccessLevelEnums } from '~/components/partials/access/helpers';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { useMailerStore } from '~/stores/mailer';
import { getMailerMailerClient } from '~~/gen/ts/clients';
import { type Access, AccessLevel } from '~~/gen/ts/resources/mailer/access';
import type { Email } from '~~/gen/ts/resources/mailer/email';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { GetEmailProposalsResponse } from '~~/gen/ts/services/mailer/mailer';

const props = withDefaults(
    defineProps<{
        modelValue?: Email | undefined;
        personalEmail?: boolean;
        disabled?: boolean;
        hideLabel?: boolean;
    }>(),
    {
        modelValue: undefined,
        personalEmail: false,
        disabled: false,
        hideLabel: false,
    },
);

const emit = defineEmits<{
    (e: 'update:modelValue', email: Email | undefined): void;
    (e: 'refresh'): void;
}>();

const { t } = useI18n();

const { activeChar, isSuperuser } = useAuth();

const notifications = useNotificationsStore();

const mailerStore = useMailerStore();
const { selectedEmail, emails } = storeToRefs(mailerStore);

const mailerMailerClient = await getMailerMailerClient();

const { data: proposals, refresh: refreshProposabls } = useLazyAsyncData(`emails-proposals`, () => getEmailProposals());

async function getEmailProposals(): Promise<GetEmailProposalsResponse> {
    try {
        const call = mailerMailerClient.getEmailProposals({
            input: '',
            job: !props.personalEmail,
            userId: isSuperuser.value ? selectedEmail.value?.userId : undefined,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(
    () => props.personalEmail,
    async () => refreshProposabls(),
);

watch(proposals, () => {
    if (state.domain === '') {
        if (proposals.value?.domains[0]) {
            state.domain = proposals.value?.domains[0];
        }
    }
});

const schema = z.object({
    email: z
        .string()
        .min(3)
        .max(50)
        .refine((email) => (props.personalEmail ? proposals.value?.emails.includes(email) : true), {
            message: t('errors.MailerService.ErrAddresseInvalid'),
        }),
    domain: z
        .string()
        .min(6)
        .max(50)
        .refine((domain) => proposals.value?.domains.includes(domain), {
            message: t('errors.MailerService.ErrAddresseInvalid'),
        }),
    label: z.string().max(128).optional(),
    deactivated: z.coerce.boolean(),
    access: z.custom<Access>(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    email: '',
    domain: '',
    deactivated: false,
    access: {
        jobs: [],
        users: [],
        qualifications: [],
    },
});

function setFromProps(): void {
    if (!props.modelValue || !props.modelValue?.email) {
        return;
    }

    const split = props.modelValue?.email.split('@');
    if (split[0] && split[1]) {
        state.email = split[0];
        state.domain = split[1];
    }

    state.deactivated = props.modelValue.deactivated;
    if (props.modelValue.access) {
        state.access = props.modelValue.access;
    }
}

setFromProps();
watch(props, setFromProps);

async function createOrUpdateEmail(values: Schema): Promise<undefined> {
    values.access.users.forEach((user) => {
        if (user.id < 0) user.id = 0;
        user.user = undefined; // Clear user object to avoid sending unnecessary data
    });
    values.access.jobs.forEach((job) => job.id < 0 && (job.id = 0));
    values.access.qualifications.forEach((quali) => quali.id < 0 && (quali.id = 0));

    const redirectToMail = emails.value.length === 0;

    const response = await mailerStore.createOrUpdateEmail({
        email: {
            id: props.modelValue?.id ?? 0,
            email: values.email + '@' + values.domain,
            label: values.label !== '' ? values.label : undefined,
            deactivated: values.deactivated,
            job: !props.personalEmail ? (props.modelValue?.job ?? activeChar.value!.job) : undefined,
            userId: props.personalEmail ? (props.modelValue?.userId ?? activeChar.value!.userId) : undefined,
            access: values.access,
        },
    });

    notifications.add({
        title: { key: 'notifications.action_successful.title', parameters: {} },
        description: { key: 'notifications.action_successful.content', parameters: {} },
        type: NotificationType.SUCCESS,
    });

    if (redirectToMail) {
        await navigateTo({ name: 'mail' });
    }

    if (response.email) {
        emit('update:modelValue', response.email);

        // Restart notificator stream
        await notifications.restartStream();
    }

    emit('refresh');
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;

    await createOrUpdateEmail(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UForm class="flex flex-col gap-y-2" :state="state" :schema="schema" @submit="onSubmitThrottle">
        <UFormField
            class="flex flex-1 flex-col"
            :label="$t('common.mail')"
            :description="
                $t('components.mailer.manage.email.description') +
                (modelValue?.emailChanged ? ` (${$t('common.last_updated')}: ${$d(toDate(modelValue?.emailChanged))})` : '')
            "
        >
            <div class="flex w-full flex-1 flex-col gap-1 sm:flex-row">
                <UFormField class="flex-1" name="email">
                    <USelectMenu
                        v-if="proposals?.emails && proposals.emails.length > 0"
                        v-model="state.email"
                        class="flex-1"
                        :items="proposals?.emails"
                        :disabled="disabled"
                    >
                        <template #empty>
                            {{ $t('common.not_found', [$t('common.mail')]) }}
                        </template>
                    </USelectMenu>
                    <UInput
                        v-else
                        v-model="state.email"
                        class="flex-1"
                        type="text"
                        :placeholder="$t('common.mail')"
                        :disabled="disabled"
                    />
                </UFormField>

                <span class="flex-initial font-semibold">@</span>

                <UFormField class="flex-1" name="domain">
                    <USelectMenu
                        v-if="proposals?.domains && proposals.domains.length > 1"
                        v-model="state.domain"
                        class="flex-1"
                        :items="proposals?.domains"
                        :disabled="disabled"
                    >
                        <template #empty>
                            {{ $t('common.not_found', [$t('common.mail')]) }}
                        </template>
                    </USelectMenu>
                    <UInput
                        v-else
                        v-model="state.domain"
                        class="flex-1"
                        type="text"
                        :placeholder="$t('common.mail')"
                        disabled
                    />
                </UFormField>
            </div>
        </UFormField>

        <UFormField v-if="!hideLabel" name="label" :label="$t('common.label')">
            <UInput v-model="state.label" type="text" :disabled="disabled" />
        </UFormField>

        <UFormField
            v-if="modelValue?.id !== undefined && (isSuperuser || state.deactivated)"
            name="disabled"
            :label="$t('common.disabled')"
        >
            <USwitch v-model="state.deactivated" :disabled="disabled" />
        </UFormField>

        <UFormField v-if="!personalEmail" name="access" :label="$t('common.access')">
            <AccessManager
                v-model:jobs="state.access!.jobs"
                v-model:users="state.access!.users"
                v-model:qualifications="state.access!.qualifications"
                :target-id="modelValue?.id ?? 0"
                :access-types="[
                    { type: 'user', label: $t('common.citizen', 2) },
                    { type: 'job', label: $t('common.job', 2) },
                    { type: 'qualification', label: $t('common.qualification', 2) },
                ]"
                :access-roles="enumToAccessLevelEnums(AccessLevel, 'enums.mailer.AccessLevel')"
                :disabled="disabled"
                default-access-type="job"
            />
        </UFormField>

        <UFormField>
            <DataErrorBlock
                v-if="modelValue?.deactivated"
                :title="$t('errors.MailerService.ErrEmailDisabled.title')"
                :message="$t('errors.MailerService.ErrEmailDisabled.content')"
            />

            <UButton
                v-if="!disabled"
                type="submit"
                block
                :label="modelValue?.id !== undefined ? $t('common.update') : $t('common.create')"
            />
        </UFormField>
    </UForm>
</template>
