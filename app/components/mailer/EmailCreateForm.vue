<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { useMailerStore } from '~/store/mailer';
import { useNotificatorStore } from '~/store/notificator';
import { type Access, AccessLevel } from '~~/gen/ts/resources/mailer/access';
import type { Email } from '~~/gen/ts/resources/mailer/email';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { GetEmailProposalsResponse } from '~~/gen/ts/services/mailer/mailer';
import AccessManager from '../partials/access/AccessManager.vue';
import { enumToAccessLevelEnums } from '../partials/access/helpers';
import DataErrorBlock from '../partials/data/DataErrorBlock.vue';

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

const notifications = useNotificatorStore();

const mailerStore = useMailerStore();
const { selectedEmail, emails } = storeToRefs(mailerStore);

const { data: proposals, refresh: refreshProposabls } = useLazyAsyncData(`emails-proposals`, () => getEmailProposals());

async function getEmailProposals(): Promise<GetEmailProposalsResponse> {
    try {
        const call = getGRPCMailerClient().getEmailProposals({
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
    deactivated: z.boolean(),
    internal: z.boolean(),
    access: z.custom<Access>(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    email: '',
    domain: '',
    deactivated: false,
    internal: false,
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
    state.internal = props.modelValue.internal;
    if (props.modelValue.access) {
        state.access = props.modelValue.access;
    }
}

setFromProps();
watch(props, setFromProps);

async function createOrUpdateEmail(values: Schema): Promise<undefined> {
    const redirectToMail = emails.value.length === 0;

    const response = await mailerStore.createOrUpdateEmail({
        email: {
            id: props.modelValue?.id ?? 0,
            email: values.email + '@' + values.domain,
            internal: values.internal,
            label: values.label !== '' ? values.label : undefined,
            deactivated: values.deactivated,
            job: !props.personalEmail ? (props.modelValue?.job ?? activeChar.value!.job) : undefined,
            userId: props.personalEmail ? (props.modelValue?.userId ?? activeChar.value!.userId) : undefined,
            access: values.access,
        },
    });

    notifications.add({
        title: { key: 'notifications.action_successfull.title', parameters: {} },
        description: { key: 'notifications.action_successfull.content', parameters: {} },
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
    <UForm :state="state" :schema="schema" class="flex flex-col gap-y-2" @submit="onSubmitThrottle">
        <UFormGroup
            :label="$t('common.mail')"
            :description="
                $t('components.mailer.manage.email.description') +
                (modelValue?.emailChanged ? ` (${$t('common.last_updated')}: ${$d(toDate(modelValue?.emailChanged))})` : '')
            "
            class="flex flex-1 flex-col"
        >
            <div class="flex w-full flex-1 flex-col gap-1 sm:flex-row">
                <UFormGroup name="email" class="flex-1">
                    <USelectMenu
                        v-if="proposals?.emails && proposals.emails.length > 0"
                        v-model="state.email"
                        :options="proposals?.emails"
                        :disabled="disabled"
                        class="flex-1"
                    >
                        <template #empty>
                            {{ $t('common.not_found', [$t('common.mail')]) }}
                        </template>
                    </USelectMenu>
                    <UInput
                        v-else
                        v-model="state.email"
                        type="text"
                        :placeholder="$t('common.mail')"
                        :disabled="disabled"
                        class="flex-1"
                    />
                </UFormGroup>

                <span class="flex-initial font-semibold">@</span>

                <UFormGroup name="domain" class="flex-1">
                    <USelectMenu
                        v-if="proposals?.domains && proposals.domains.length > 1"
                        v-model="state.domain"
                        :options="proposals?.domains"
                        :disabled="disabled"
                        class="flex-1"
                    >
                        <template #empty>
                            {{ $t('common.not_found', [$t('common.mail')]) }}
                        </template>
                    </USelectMenu>
                    <UInput
                        v-else
                        v-model="state.domain"
                        type="text"
                        :placeholder="$t('common.mail')"
                        disabled
                        class="flex-1"
                    />
                </UFormGroup>
            </div>
        </UFormGroup>

        <UFormGroup v-if="!hideLabel" name="label" :label="$t('common.label')">
            <UInput v-model="state.label" type="text" :disabled="disabled" />
        </UFormGroup>

        <UFormGroup
            v-if="modelValue?.id !== undefined && (isSuperuser || state.deactivated)"
            name="disabled"
            :label="$t('common.disabled')"
        >
            <UToggle v-model="state.deactivated" :disabled="disabled" />
        </UFormGroup>

        <UFormGroup v-if="!personalEmail" name="access" :label="$t('common.access')">
            <AccessManager
                v-model:jobs="state.access!.jobs"
                v-model:users="state.access!.users"
                v-model:qualifications="state.access!.qualifications"
                :target-id="modelValue?.id ?? 0"
                :access-types="[
                    { type: 'user', name: $t('common.citizen', 2) },
                    { type: 'job', name: $t('common.job', 2) },
                    { type: 'qualification', name: $t('common.qualification', 2) },
                ]"
                :access-roles="enumToAccessLevelEnums(AccessLevel, 'enums.mailer.AccessLevel')"
                :disabled="disabled"
                default-access-type="job"
            />
        </UFormGroup>

        <UFormGroup>
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
        </UFormGroup>
    </UForm>
</template>
