<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import type { JSONContent } from '@tiptap/core';
import { z } from 'zod';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { enumToAccessLevelEnums, normalizeAccessEntryIds } from '~/components/partials/access/helpers';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import { useMailerStore } from '~/stores/mailer';
import { contentToTiptapValue, tiptapToContent } from '~/utils/content';
import { getMailerMailerClient } from '~~/gen/ts/clients';
import { AccessLevel } from '~~/gen/ts/resources/mailer/access/access';
import type { Email } from '~~/gen/ts/resources/mailer/emails/email';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { GetEmailProposalsResponse } from '~~/gen/ts/services/mailer/mailer';

const props = withDefaults(
    defineProps<{
        personalEmail?: boolean;
        disabled?: boolean;
        hideLabel?: boolean;
    }>(),
    {
        personalEmail: false,
        disabled: false,
        hideLabel: false,
    },
);

const emit = defineEmits<{
    (e: 'refresh'): void;
    (e: 'dirty-change', value: boolean): void;
    (e: 'private-creation-state', value: boolean): void;
}>();

const email = defineModel<Email | undefined>({ default: undefined });
const showCreationSuccess = ref<boolean>(false);
const createdEmailAddress = ref<string>('');

const { t } = useI18n();

const { activeChar, isSuperuser } = useAuth();

const { maxAccessEntries } = useAppConfig();

const notifications = useNotificationsStore();
const overlay = useOverlay();
const confirmModal = overlay.create(ConfirmModal);

const mailerStore = useMailerStore();
const { addressBook, selectedEmail, emails } = storeToRefs(mailerStore);

const mailerMailerClient = await getMailerMailerClient();

const { data: proposals, refresh: refreshProposabls } = useAuthedLazyAsyncData('userState', `emails-proposals`, ({ signal }) =>
    getEmailProposals(signal),
);

async function getEmailProposals(signal: AbortSignal): Promise<GetEmailProposalsResponse> {
    try {
        const call = mailerMailerClient.getEmailProposals(
            {
                input: '',
                job: !props.personalEmail,
                userId: isSuperuser.value ? selectedEmail.value?.userId : undefined,
            },
            { abort: signal },
        );
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
            syncSnapshot();
        }
    }
});

const schema = z.object({
    email: z
        .string()
        .min(3)
        .max(50)
        .refine((email) => (props.personalEmail ? proposals.value?.emails.includes(email) : true), {
            message: t('errors.mailer.MailerService.ErrAddresseInvalid'),
        }),
    domain: z
        .string()
        .min(6)
        .max(50)
        .refine((domain) => proposals.value?.domains.includes(domain), {
            message: t('errors.mailer.MailerService.ErrAddresseInvalid'),
        }),
    label: z.string().max(128).optional(),
    deactivated: z.coerce.boolean(),
    access: z.object({
        jobs: jobsAccessEntries(t).max(maxAccessEntries).default([]),
        users: userAccessEntries(t).max(maxAccessEntries).default([]),
        qualifications: qualificationAccessEntries(t).max(maxAccessEntries).default([]),
    }),
    signature: z.custom<JSONContent | string>().optional(),
    blockedEmails: z.coerce.string().array().max(25).default([]),
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
    signature: '',
    blockedEmails: [],
});

const { hasUnsavedChanges, syncSnapshot } = useSnapshotChanges(state);

function clearAddressBook(): void {
    if (addressBook.value.length === 0) return;

    confirmModal.open({
        title: t('components.mailer.settings.clear_address_book'),
        description: t('components.mailer.settings.clear_address_book_confirm'),
        confirm: () => {
            addressBook.value.length = 0;
        },
    });
}

function setFromProps(): void {
    if (!email.value || !email.value?.email) {
        state.signature = '';
        state.blockedEmails = [];
        syncSnapshot();
        return;
    }

    const split = email.value?.email.split('@');
    if (split[0] && split[1]) {
        state.email = split[0];
        state.domain = split[1];
    }
    state.label = email.value.label;
    state.deactivated = email.value.deactivated;
    if (email.value.access) state.access = email.value.access;
    state.signature = contentToTiptapValue(email.value.settings?.signature);
    state.blockedEmails = email.value.settings?.blockedEmails ?? [];

    syncSnapshot();
}

setFromProps();
watch(email, setFromProps);

// AccessManager can normalize server-returned ACL entries during mount. Take
// the final baseline after those child components have finished hydrating.
onMounted(() => syncSnapshot());

watch(hasUnsavedChanges, (value) => emit('dirty-change', value), { immediate: true, flush: 'sync' });

async function createOrUpdateEmail(values: Schema): Promise<undefined> {
    normalizeAccessEntryIds(values.access.users);
    values.access.users.forEach((user) => {
        user.user = undefined; // Clear user object to avoid sending unnecessary data
    });
    normalizeAccessEntryIds(values.access.jobs);
    normalizeAccessEntryIds(values.access.qualifications);

    const isCreating = !email.value?.id;
    const redirectToMail = emails.value.length === 0;

    if (props.personalEmail && isCreating) emit('private-creation-state', true);

    let response: Awaited<ReturnType<typeof mailerStore.createOrUpdateEmail>>;
    try {
        response = await mailerStore.createOrUpdateEmail({
            email: {
                id: email.value?.id ?? 0,
                email: values.email + '@' + values.domain,
                label: values.label !== '' ? values.label : undefined,
                deactivated: values.deactivated,
                job: !props.personalEmail ? (email.value?.job ?? activeChar.value!.job) : undefined,
                userId: props.personalEmail ? (email.value?.userId ?? activeChar.value!.userId) : undefined,
                access: values.access,
            },
        });
    } catch (error) {
        if (props.personalEmail && isCreating) emit('private-creation-state', false);
        throw error;
    }

    notifications.add({
        title: {
            key:
                props.personalEmail && isCreating
                    ? 'notifications.mailer.private_email_created.title'
                    : 'notifications.action_successful.title',
            parameters: {},
        },
        description: {
            key:
                props.personalEmail && isCreating
                    ? 'notifications.mailer.private_email_created.content'
                    : 'notifications.action_successful.content',
            parameters: props.personalEmail && isCreating ? { email: values.email + '@' + values.domain } : {},
        },
        type: NotificationType.SUCCESS,
    });

    if (response.email) {
        email.value = response.email;

        await mailerStore.setEmailSettings(
            {
                settings: {
                    emailId: response.email.id,
                    signature: tiptapToContent(values.signature),
                    blockedEmails: values.blockedEmails.map((entry) => entry.trim()),
                },
            },
            false,
        );

        setFromProps();

        // The parent page may be navigating away immediately after a private
        // address is created. Clear its dirty state before that navigation.
        emit('dirty-change', false);

        if (props.personalEmail && isCreating) {
            createdEmailAddress.value = response.email.email;
            showCreationSuccess.value = true;
        }

        // Restart notificator stream
        await notifications.restartStream();
    }

    // Job-email forms use this to close their create panel. Private-email
    // creation navigates away after the success state and needs no refresh.
    if (!props.personalEmail) emit('refresh');

    if (redirectToMail && !(props.personalEmail && isCreating)) {
        await navigateTo({ name: 'mail-thread', params: { thread: undefined } });
    }

    if (props.personalEmail && isCreating) {
        await new Promise((resolve) => setTimeout(resolve, 1500));
        await navigateTo({ name: 'mail-thread', params: { thread: undefined } });
    }
}

const { submit, isSubmitting, canSubmit } = useSubmitGuard(async (event: FormSubmitEvent<Schema>) => {
    await createOrUpdateEmail(event.data);
});
</script>

<template>
    <div v-if="showCreationSuccess" class="flex flex-col items-center gap-3 p-6 text-center">
        <UIcon class="h-32 w-32 text-success" name="i-mdi-check-circle" />

        <div>
            <h3 class="text-lg font-bold text-highlighted">
                {{ $t('components.mailer.manage.email.email_created.title') }}
            </h3>
            <p class="text-dimmed">
                {{ $t('components.mailer.manage.email.email_created.content', { email: createdEmailAddress }) }}
            </p>
        </div>
    </div>

    <template v-else>
        <template v-if="personalEmail && !modelValue?.id">
            <UIcon class="h-32 w-32" name="i-mdi-email-multiple" />

            <div class="text-center text-highlighted">
                <h3 class="text-lg font-bold">{{ $t('components.mailer.manage.title') }}</h3>
                <p class="text-bas">{{ $t('components.mailer.manage.subtitle') }}</p>
            </div>
        </template>

        <UForm class="flex flex-col gap-4" :state="state" :schema="schema" @submit="submit">
            <UPageCard
                :title="personalEmail && !modelValue?.id ? undefined : $t('common.mail')"
                :description="
                    personalEmail && !modelValue?.id
                        ? undefined
                        : $t('components.mailer.manage.email.description') +
                          (modelValue?.emailChanged
                              ? ` (${$t('common.last_updated')}: ${$d(toDate(modelValue?.emailChanged))})`
                              : '')
                "
            >
                <div class="flex w-full flex-1 flex-col gap-2">
                    <UFormField class="flex flex-1 flex-col" :label="$t('common.mail')">
                        <div class="flex w-full flex-1 flex-col gap-1 sm:flex-row">
                            <UFormField class="flex-1" name="email">
                                <USelectMenu
                                    v-if="proposals?.emails && proposals.emails.length > 0"
                                    v-model="state.email"
                                    class="w-full"
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
                                    class="w-full"
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
                                    class="w-full"
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
                                    class="w-full"
                                    type="text"
                                    :placeholder="$t('common.mail')"
                                    disabled
                                />
                            </UFormField>
                        </div>
                    </UFormField>

                    <UFormField v-if="!hideLabel" name="label" :label="$t('common.label')">
                        <UInput v-model="state.label" class="w-full" type="text" :disabled="disabled" />
                    </UFormField>

                    <UFormField
                        v-if="modelValue?.id !== undefined && (isSuperuser || state.deactivated)"
                        name="disabled"
                        :label="$t('common.disabled')"
                    >
                        <USwitch v-model="state.deactivated" :disabled="disabled" />
                    </UFormField>
                </div>
            </UPageCard>

            <UPageCard v-if="!personalEmail" :title="$t('common.access')">
                <AccessManager
                    v-model:jobs="state.access!.jobs"
                    v-model:users="state.access!.users"
                    v-model:qualifications="state.access!.qualifications"
                    :target-id="modelValue?.id ?? 0"
                    :access-types="[
                        { label: $t('common.citizen', 2), value: 'user' },
                        { label: $t('common.job', 2), value: 'job' },
                        { label: $t('common.qualification', 2), value: 'qualification' },
                    ]"
                    :access-roles="enumToAccessLevelEnums(AccessLevel, 'enums.mailer.AccessLevel')"
                    :disabled="disabled"
                    default-access-type="job"
                    name="access"
                />
            </UPageCard>

            <div v-if="modelValue?.id !== undefined" class="flex flex-col gap-4">
                <UPageCard
                    :title="$t('common.blocklist')"
                    :description="$t('components.mailer.settings.blocklist_description')"
                >
                    <UFormField class="flex-1" name="blockedEmails">
                        <div class="flex flex-col gap-1">
                            <div v-for="(_, idx) in state.blockedEmails" :key="idx" class="flex items-center gap-1">
                                <UInput
                                    v-model="state.blockedEmails[idx]"
                                    class="w-full"
                                    type="text"
                                    :placeholder="$t('common.mail')"
                                    :disabled="disabled"
                                />

                                <UButton
                                    icon="i-mdi-close"
                                    :disabled="disabled || !canSubmit"
                                    :aria-label="$t('common.remove')"
                                    @click="state.blockedEmails.splice(idx, 1)"
                                />
                            </div>
                        </div>

                        <UButton
                            v-if="!disabled"
                            :class="state.blockedEmails.length ? 'mt-2' : ''"
                            icon="i-mdi-plus"
                            :disabled="!canSubmit || state.blockedEmails.length >= 25"
                            :aria-label="$t('common.add')"
                            @click="state.blockedEmails.push('')"
                        />
                    </UFormField>
                </UPageCard>

                <UPageCard
                    :title="$t('common.address_book')"
                    :description="$t('components.mailer.settings.address_book_description')"
                >
                    <UFormField>
                        <UButton
                            :label="$t('components.mailer.settings.clear_address_book')"
                            color="error"
                            icon="i-mdi-bookmark-remove"
                            :disabled="disabled || !canSubmit"
                            @click="clearAddressBook"
                        />
                    </UFormField>
                </UPageCard>

                <UPageCard
                    :title="$t('common.signature')"
                    :description="$t('components.mailer.settings.signature_description')"
                >
                    <UFormField class="flex-1" name="signature" :ui="{ error: 'hidden' }">
                        <ClientOnly>
                            <TiptapEditor
                                v-model="state.signature"
                                name="signature"
                                :disabled="disabled"
                                :limit="1024"
                                wrapper-class="min-h-44"
                                content-type="json"
                                disable-images
                            />
                        </ClientOnly>
                    </UFormField>
                </UPageCard>
            </div>

            <UFormField>
                <DataErrorBlock
                    v-if="modelValue?.deactivated"
                    :title="$t('errors.mailer.MailerService.ErrEmailDisabled.title')"
                    :message="$t('errors.mailer.MailerService.ErrEmailDisabled.content')"
                />

                <UButton
                    v-if="!disabled"
                    type="submit"
                    block
                    :disabled="!canSubmit"
                    :loading="isSubmitting"
                    :label="modelValue?.id !== undefined ? $t('common.update') : $t('common.create')"
                />
            </UFormField>
        </UForm>
    </template>
</template>
