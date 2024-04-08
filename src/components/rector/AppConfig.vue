<script lang="ts" setup>
import { z } from 'zod';
import type { FormSubmitEvent } from '#ui/types';
import { useSettingsStore } from '~/store/settings';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { type GetAppConfigResponse } from '~~/gen/ts/services/rector/config';
import { useNotificatorStore } from '~/store/notificator';
import { useCompletorStore } from '~/store/completor';
import { toDuration } from '~/utils/duration';

const { $grpc } = useNuxtApp();

const settingsStore = useSettingsStore();
const { streamerMode } = storeToRefs(settingsStore);

const notifications = useNotificatorStore();

const { data, pending, refresh, error } = useLazyAsyncData(`rector-appconfig`, () => getAppConfig());

async function getAppConfig(): Promise<GetAppConfigResponse> {
    try {
        const call = $grpc.getRectorConfigClient().getAppConfig({});
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const completorStore = useCompletorStore();
const { listJobs } = completorStore;

const { data: jobs } = useLazyAsyncData(`rector-appconfig-jobs`, () => listJobs());

const schema = z.object({
    auth: z.object({
        signupEnabled: z.boolean(),
    }),
    perms: z.object({
        default: z
            .object({
                category: z.string().min(1).max(48),
                name: z.string().min(1).max(48),
            })
            .array()
            .max(25),
    }),
    website: z.object({
        links: z.object({
            privacyPolicy: z.union([z.string().min(1).max(255).url().startsWith('https://'), z.string().length(0).optional()]),
            imprint: z.union([z.string().min(1).max(255).url().startsWith('https://'), z.string().length(0).optional()]),
        }),
    }),
    jobInfo: z.object({
        unemployedJob: z.object({
            name: z.string().min(1).max(20),
            grade: z.coerce.number().min(1).max(99),
        }),
        publicJobs: z.string().array().max(99),
        hiddenJobs: z.string().array().max(99),
    }),
    userTracker: z.object({
        refreshTime: zodDurationSchema,
        dbRefreshTime: zodDurationSchema,
        livemapJobs: z.string().array().max(99),
    }),
    // Discord
    discord: z.object({
        enabled: z.boolean(),
        syncInterval: zodDurationSchema,
        inviteUrl: z.union([
            z.string().min(1).max(255).url().startsWith('https://discord.com/'),
            z.string().length(0).optional(),
        ]),
    }),
});

const state = reactive<Schema>({
    auth: {
        signupEnabled: false,
    },
    perms: {
        default: [],
    },
    website: {
        links: {},
    },
    jobInfo: {
        hiddenJobs: [],
        publicJobs: [],
        unemployedJob: {
            name: '',
            grade: 1,
        },
    },
    userTracker: {
        dbRefreshTime: '1s',
        refreshTime: '3.35s',
        livemapJobs: [],
    },
    // Discord
    discord: {
        enabled: false,
        syncInterval: '9s',
        inviteUrl: '',
    },
});

type Schema = z.output<typeof schema>;

async function updateAppConfig(values: Schema): Promise<void> {
    if (!data.value || !data.value?.config) {
        return;
    }

    // Update local version of retrieved config
    data.value.config.auth = values.auth;
    data.value.config.perms = values.perms;
    data.value.config.website = values.website;
    data.value.config.jobInfo = values.jobInfo;
    data.value.config.userTracker = {
        livemapJobs: values.userTracker.livemapJobs,
        dbRefreshTime: toDuration(values.userTracker.dbRefreshTime),
        refreshTime: toDuration(values.userTracker.refreshTime),
    };
    data.value.config.discord = {
        enabled: values.discord.enabled,
        inviteUrl: values.discord.inviteUrl,
        syncInterval: toDuration(values.discord.syncInterval),
    };

    try {
        const { response } = await $grpc.getRectorConfigClient().updateAppConfig({
            config: data.value?.config,
        });

        notifications.add({
            title: { key: 'notifications.rector.app_config.title', parameters: {} },
            description: { key: 'notifications.rector.app_config.content', parameters: {} },
            type: 'success',
        });

        if (response.config) {
            data.value = response;
        } else {
            refresh();
        }
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

function setSettingsValues(): void {
    if (!data.value || !data.value.config) {
        return;
    }

    if (data.value.config.auth) {
        state.auth = data.value.config.auth;
    }
    if (data.value.config.perms) {
        state.perms = data.value.config.perms;
    }
    if (data.value.config.website) {
        if (data.value.config.website.links) {
            state.website.links = data.value.config.website.links;
        }
    }
    if (data.value.config.jobInfo) {
        if (data.value.config.jobInfo.unemployedJob) {
            state.jobInfo.unemployedJob = data.value.config.jobInfo.unemployedJob;
        }
        state.jobInfo.hiddenJobs = data.value.config.jobInfo.hiddenJobs;
        state.jobInfo.publicJobs = data.value.config.jobInfo.publicJobs;
    }
    if (data.value.config.userTracker) {
        if (data.value.config.userTracker.dbRefreshTime) {
            state.userTracker.dbRefreshTime = fromDuration(data.value.config.userTracker.dbRefreshTime);
        }
        if (data.value.config.userTracker.refreshTime) {
            state.userTracker.refreshTime = fromDuration(data.value.config.userTracker.refreshTime);
        }
        state.userTracker.livemapJobs = data.value.config.userTracker.livemapJobs;
    }
    if (data.value.config.discord) {
        state.discord.enabled = data.value.config.discord.enabled;
        if (data.value.config.discord.syncInterval) {
            state.discord.syncInterval = fromDuration(data.value.config.discord.syncInterval);
        }
        state.discord.inviteUrl = data.value.config.discord.inviteUrl;
    }
}

watchOnce(data, () => setSettingsValues());

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await updateAppConfig(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <div>
        <template v-if="streamerMode">
            <UDashboardPanelContent class="pb-2">
                <UDashboardSection
                    :title="$t('system.streamer_mode.title')"
                    :description="$t('system.streamer_mode.description')"
                />
            </UDashboardPanelContent>
        </template>
        <template v-else>
            <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.setting', 2)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.setting', 2)])"
                :retry="refresh"
            />
            <DataNoDataBlock v-else-if="data === null" icon="i-mdi-office-building-cog" :type="$t('common.setting', 2)" />

            <template v-else>
                <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
                    <UDashboardNavbar :title="$t('pages.rector.settings.title')">
                        <template #right>
                            <UButton
                                type="submit"
                                class="flex w-full justify-center rounded-md px-3 py-2 text-sm font-semibold transition-colors focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                                :disabled="!canSubmit"
                                :loading="!canSubmit"
                            >
                                {{ $t('common.save', 1) }}
                            </UButton>
                        </template>
                    </UDashboardNavbar>

                    <UDashboardPanelContent class="pb-2">
                        <UDashboardSection
                            :title="$t('components.rector.app_config.auth.title')"
                            :description="$t('components.rector.app_config.auth.description')"
                        >
                            <UFormGroup
                                name="auth.signupEnabled"
                                :label="$t('components.rector.app_config.auth.sign_up')"
                                class="grid grid-cols-2 items-center gap-2"
                                :ui="{ container: '' }"
                            >
                                <UToggle v-model="state.auth.signupEnabled">
                                    <span class="sr-only">
                                        {{ $t('components.rector.app_config.auth.sign_up') }}
                                    </span>
                                </UToggle>
                            </UFormGroup>
                        </UDashboardSection>

                        <UDashboardSection
                            :title="$t('components.rector.app_config.perms.title')"
                            :description="$t('components.rector.app_config.perms.description')"
                        >
                            <UFormGroup
                                name="perms.default"
                                :label="$t('components.rector.app_config.perms.default_perms')"
                                class="grid grid-cols-2 items-center gap-2"
                                :ui="{ container: '' }"
                            >
                                <div class="flex flex-col gap-1">
                                    <div v-for="(perm, idx) in state.perms.default" :key="idx" class="flex items-center gap-1">
                                        <UFormGroup :name="`perms.default.${idx}.category`" class="flex-1">
                                            <UInput
                                                v-model="state.perms.default[idx].category"
                                                type="text"
                                                :placeholder="$t('common.category')"
                                                @focusin="focusTablet(true)"
                                                @focusout="focusTablet(false)"
                                            />
                                        </UFormGroup>

                                        <UFormGroup :name="`perms.default.${idx}.name`" class="flex-1">
                                            <UInput
                                                v-model="state.perms.default[idx].name"
                                                type="text"
                                                :placeholder="$t('common.name')"
                                                @focusin="focusTablet(true)"
                                                @focusout="focusTablet(false)"
                                            />
                                        </UFormGroup>

                                        <UButton
                                            :ui="{ rounded: 'rounded-full' }"
                                            icon="i-mdi-close"
                                            @click="state.perms.default.splice(idx, 1)"
                                        />
                                    </div>
                                </div>

                                <UButton
                                    class="mt-2"
                                    :ui="{ rounded: 'rounded-full' }"
                                    :disabled="!canSubmit"
                                    icon="i-mdi-plus"
                                    @click="state.perms.default.push({ category: '', name: '' })"
                                />
                            </UFormGroup>
                        </UDashboardSection>

                        <UDashboardSection
                            :title="$t('components.rector.app_config.website.title')"
                            :description="$t('components.rector.app_config.website.description')"
                        >
                            <UFormGroup
                                name="website.links.privacyPolicy"
                                :label="$t('common.privacy_policy')"
                                class="grid grid-cols-2 items-center gap-2"
                                :ui="{ container: '' }"
                            >
                                <UInput
                                    v-model="state.website.links.privacyPolicy"
                                    type="text"
                                    :placeholder="$t('common.privacy_policy')"
                                    maxlength="255"
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                />
                            </UFormGroup>

                            <UFormGroup
                                name="website.links.imprint"
                                :label="$t('common.imprint')"
                                class="grid grid-cols-2 items-center gap-2"
                                :ui="{ container: '' }"
                            >
                                <UInput
                                    v-model="state.website.links.imprint"
                                    type="text"
                                    :placeholder="$t('common.imprint')"
                                    maxlength="255"
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                />
                            </UFormGroup>
                        </UDashboardSection>

                        <UDashboardSection
                            :title="$t('components.rector.app_config.job_info.title')"
                            :description="$t('components.rector.app_config.job_info.description')"
                        >
                            <UFormGroup
                                name="jobInfo.unemployedJob.name"
                                :label="`${$t('common.job')} ${$t('common.name')}`"
                                class="grid grid-cols-2 items-center gap-2"
                                :ui="{ container: '' }"
                            >
                                <UInput
                                    v-model="state.jobInfo.unemployedJob.name"
                                    type="text"
                                    :placeholder="$t('common.job')"
                                    maxlength="255"
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                />
                            </UFormGroup>

                            <UFormGroup
                                name="jobInfo.unemployedJob.grade"
                                :label="$t('common.rank')"
                                class="grid grid-cols-2 items-center gap-2"
                                :ui="{ container: '' }"
                            >
                                <UInput
                                    v-model="state.jobInfo.unemployedJob.grade"
                                    type="number"
                                    min="1"
                                    max="99"
                                    name="jobInfoUnemployedGrade"
                                    :placeholder="$t('common.rank')"
                                    :label="$t('common.rank')"
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                />
                            </UFormGroup>

                            <UFormGroup
                                name="jobInfo.publicJobs"
                                :label="$t('components.rector.app_config.job_info.public_jobs')"
                                class="grid grid-cols-2 items-center gap-2"
                                :ui="{ container: '' }"
                            >
                                <USelectMenu
                                    v-model="state.jobInfo.publicJobs"
                                    multiple
                                    :options="jobs ?? []"
                                    value-attribute="name"
                                    by="label"
                                >
                                    <template #label>
                                        <template v-if="state.jobInfo.publicJobs">
                                            <span class="truncate">{{ state.jobInfo.publicJobs.join(',') }}</span>
                                        </template>
                                        <template v-else>
                                            <span class="truncate">{{ $t('common.none_selected', [$t('common.job')]) }}</span>
                                        </template>
                                    </template>
                                    <template #option="{ option: job }">
                                        <span class="truncate">{{ job.label }} ({{ job.name }})</span>
                                    </template>
                                </USelectMenu>
                            </UFormGroup>

                            <UFormGroup
                                name="jobInfo.hiddenJobs"
                                :label="$t('components.rector.app_config.job_info.hidden_jobs')"
                                class="grid grid-cols-2 items-center gap-2"
                                :ui="{ container: '' }"
                            >
                                <USelectMenu
                                    v-model="state.jobInfo.hiddenJobs"
                                    multiple
                                    :options="jobs ?? []"
                                    value-attribute="name"
                                    by="label"
                                >
                                    <template #label>
                                        <template v-if="state.jobInfo.hiddenJobs">
                                            <span class="truncate">{{ state.jobInfo.hiddenJobs.join(',') }}</span>
                                        </template>
                                        <template v-else>
                                            <span class="truncate">{{ $t('common.none_selected', [$t('common.job')]) }}</span>
                                        </template>
                                    </template>
                                    <template #option="{ option: job }">
                                        <span class="truncate">{{ job.label }} ({{ job.name }})</span>
                                    </template>
                                </USelectMenu>
                            </UFormGroup>
                        </UDashboardSection>

                        <UDashboardSection
                            :title="$t('components.rector.app_config.user_tracker.title')"
                            :description="$t('components.rector.app_config.user_tracker.description')"
                        >
                            <UFormGroup
                                name="userTracker.refreshTime"
                                :label="$t('components.rector.app_config.user_tracker.refresh_time')"
                                class="grid grid-cols-2 items-center gap-2"
                                :ui="{ container: '' }"
                            >
                                <UInput
                                    v-model="state.userTracker.refreshTime"
                                    type="text"
                                    :placeholder="$t('common.duration')"
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                />
                            </UFormGroup>

                            <UFormGroup
                                name="userTracker.dbRefreshTime"
                                :label="$t('components.rector.app_config.user_tracker.db_refresh_time')"
                                class="grid grid-cols-2 items-center gap-2"
                                :ui="{ container: '' }"
                            >
                                <UInput
                                    v-model="state.userTracker.dbRefreshTime"
                                    type="text"
                                    :placeholder="$t('common.duration')"
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                />
                            </UFormGroup>

                            <UFormGroup
                                name="userTracker.livemapJobs"
                                :label="$t('components.rector.app_config.user_tracker.livemap_jobs')"
                                class="grid grid-cols-2 items-center gap-2"
                                :ui="{ container: '' }"
                            >
                                <USelectMenu
                                    v-model="state.userTracker.livemapJobs"
                                    multiple
                                    :options="jobs ?? []"
                                    value-attribute="name"
                                    by="label"
                                >
                                    <template #label>
                                        <template v-if="state.userTracker.livemapJobs">
                                            <span class="truncate">{{ state.userTracker.livemapJobs.join(',') }}</span>
                                        </template>
                                        <template v-else>
                                            <span class="truncate">{{ $t('common.none_selected', [$t('common.job')]) }}</span>
                                        </template>
                                    </template>
                                    <template #option="{ option: job }">
                                        <span class="truncate">{{ job.label }} ({{ job.name }})</span>
                                    </template>
                                </USelectMenu>
                            </UFormGroup>
                        </UDashboardSection>

                        <!-- Discord -->
                        <UDashboardSection
                            :title="$t('common.discord')"
                            :description="$t('components.rector.app_config.discord.description')"
                        >
                            <UFormGroup
                                name="discordEnabled"
                                :label="$t('common.enabled')"
                                class="grid grid-cols-2 items-center gap-2"
                                :ui="{ container: '' }"
                            >
                                <UToggle v-model="state.discord.enabled">
                                    <span class="sr-only">
                                        {{ $t('common.enabled') }}
                                    </span>
                                </UToggle>
                            </UFormGroup>

                            <UFormGroup
                                name="discord.syncInterval"
                                :label="$t('components.rector.app_config.discord.sync_interval')"
                                class="grid grid-cols-2 items-center gap-2"
                                :ui="{ container: '' }"
                            >
                                <UInput
                                    v-model="state.discord.syncInterval"
                                    name="discord.syncInterval"
                                    type="text"
                                    :placeholder="$t('common.duration')"
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                />
                            </UFormGroup>

                            <UFormGroup
                                name="discord.inviteUrl"
                                :label="$t('components.rector.app_config.discord.bot_invite_url')"
                                class="grid grid-cols-2 items-center gap-2"
                                :ui="{ container: '' }"
                            >
                                <UInput
                                    v-model="state.discord.inviteUrl"
                                    type="text"
                                    :placeholder="$t('components.rector.app_config.discord.bot_invite_url')"
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                />
                            </UFormGroup>
                        </UDashboardSection>
                    </UDashboardPanelContent>
                </UForm>
            </template>
        </template>
    </div>
</template>
