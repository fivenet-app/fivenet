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
import { Perms, type Auth, Website, Discord, UserTracker, JobInfo } from '~~/gen/ts/resources/rector/config';

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

// TODO add custom validation and transformers for durations
const schema = z.object({
    auth: z.custom<Auth>(),
    perms: z.custom<Perms>(),
    website: z.custom<Website>(),
    jobInfo: z.custom<JobInfo>(),
    userTracker: z.custom<UserTracker>(),
    discord: z.custom<Discord>(),
});

type Schema = z.output<typeof schema>;

async function updateAppConfig(values: Schema): Promise<void> {
    if (!data.value?.config) {
        return;
    }

    data.value.config = values;

    try {
        const { response } = await $grpc.getRectorConfigClient().updateAppConfig({
            config: data.value.config,
        });

        notifications.add({
            title: { key: 'notifications.rector.app_config.title', parameters: {} },
            description: { key: 'notifications.rector.app_config.content', parameters: {} },
            type: 'success',
        });

        if (response.config) {
            data.value.config = response.config;
        } else {
            refresh();
        }
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const { setValues } = useForm<FormData>({
    validationSchema: {
        permsDefault: { size: 25 },

        websitePrivacyPolicy: { required: false, max: 255, url: 'https://.*' },
        websiteLinksImprint: { required: false, max: 255, url: 'https://.*' },

        jobInfoUnemployedName: { required: true, max: 20 },
        jobInfoUnemployedGrade: { required: true, numeric: true, min_value: 1, max_value: 99 },

        userTrackerRefreshTime: { required: true, max: 5, regex: /^\d+(\.\d+)?s$/ },
        userTrackerDbRefreshTime: { required: true, max: 5, regex: /^\d+(\.\d+)?s$/ },

        discordBotInviteUrl: { required: false, url: 'https://discord.com/.*' },
        discordSyncInterval: { required: true, max: 5, regex: /^\d+(\.\d+)?s$/ },
    },
});

function setSettingsValues(): void {
    if (!data.value) {
        return;
    }

    setValues({
        permsDefault: data.value.config?.perms?.default,

        websiteLinksPrivacyPolicy: data.value.config?.website?.links?.privacyPolicy,
        websiteLinksImprint: data.value.config?.website?.links?.imprint,

        jobInfoUnemployedName: data.value.config?.jobInfo?.unemployedJob?.name,
        jobInfoUnemployedGrade: data.value.config?.jobInfo?.unemployedJob?.grade,

        userTrackerRefreshTime: data.value.config?.userTracker?.refreshTime
            ? parseFloat(
                  data.value.config?.userTracker?.refreshTime?.seconds.toString() +
                      '.' +
                      (data.value.config?.userTracker?.refreshTime?.nanos ?? 0) / 1000000,
              ).toString() + 's'
            : undefined,
        userTrackerDbRefreshTime: data.value.config?.userTracker?.dbRefreshTime
            ? parseFloat(
                  data.value.config?.userTracker?.dbRefreshTime?.seconds.toString() +
                      '.' +
                      (data.value.config?.userTracker?.dbRefreshTime?.nanos ?? 0) / 1000000,
              ).toString() + 's'
            : undefined,

        discordBotInviteUrl: data.value.config?.discord?.inviteUrl,
        discordSyncInterval: data.value.config?.discord?.syncInterval
            ? parseFloat(
                  data.value.config?.discord?.syncInterval?.seconds.toString() +
                      '.' +
                      (data.value.config?.discord?.syncInterval?.nanos ?? 0) / 1000000,
              ).toString() + 's'
            : undefined,
    });
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
                <UDashboardNavbar :title="$t('pages.rector.settings.title')">
                    <template #right>
                        <UButton
                            class="flex w-full justify-center rounded-md px-3 py-2 text-sm font-semibold transition-colors focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                            :disabled="!canSubmit"
                            :loading="!canSubmit"
                            @click="onSubmitThrottle"
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
                            name="authSignupEnabled"
                            :label="$t('components.rector.app_config.auth.sign_up')"
                            class="grid grid-cols-2 items-center gap-2"
                            :ui="{ container: '' }"
                        >
                            <UToggle v-model="data.config!.auth!.signupEnabled">
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
                            name="permsDefaultPerms"
                            :label="$t('components.rector.app_config.perms.default_perms')"
                            class="grid grid-cols-2 items-center gap-2"
                            :ui="{ container: '' }"
                        >
                            <div class="flex flex-col gap-1">
                                <div
                                    v-for="(perm, idx) in data.config!.perms!.default"
                                    :key="idx"
                                    class="flex items-center gap-1"
                                >
                                    <div class="flex-1">
                                        <UInput
                                            v-model="perm.category"
                                            type="text"
                                            :placeholder="$t('common.category')"
                                            :label="$t('common.category')"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        />
                                    </div>
                                    <div class="flex-1">
                                        <UInput
                                            v-model="perm.name"
                                            type="text"
                                            :placeholder="$t('common.name')"
                                            :label="$t('common.name')"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        />
                                    </div>

                                    <UButton
                                        :ui="{ rounded: 'rounded-full' }"
                                        icon="i-mdi-close"
                                        @click="data.config!.perms!.default.splice(idx, 1)"
                                    />
                                </div>
                            </div>

                            <UButton
                                class="mt-2"
                                :ui="{ rounded: 'rounded-full' }"
                                :disabled="!canSubmit"
                                icon="i-mdi-plus"
                                @click="data.config!.perms!.default.push({ category: '', name: '' })"
                            >
                            </UButton>
                        </UFormGroup>
                    </UDashboardSection>

                    <UDashboardSection
                        :title="$t('components.rector.app_config.website.title')"
                        :description="$t('components.rector.app_config.website.description')"
                    >
                        <UFormGroup
                            name="websiteLinksPrivacyPolicy"
                            :label="$t('common.privacy_policy')"
                            class="grid grid-cols-2 items-center gap-2"
                            :ui="{ container: '' }"
                        >
                            <UInput
                                type="text"
                                name="websiteLinksPrivacyPolicy"
                                :value="data.config!.website!.links!.privacyPolicy"
                                :placeholder="$t('common.privacy_policy')"
                                :label="$t('common.privacy_policy')"
                                maxlength="128"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                        </UFormGroup>

                        <UFormGroup
                            name="websiteLinksImprint"
                            :label="$t('common.imprint')"
                            class="grid grid-cols-2 items-center gap-2"
                            :ui="{ container: '' }"
                        >
                            <UInput
                                type="text"
                                name="websiteLinksImprint"
                                :placeholder="$t('common.imprint')"
                                :value="data.config!.website!.links!.imprint"
                                maxlength="128"
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
                            name="jobInfoUnemployedName"
                            :label="`${$t('common.job')} ${$t('common.name')}`"
                            class="grid grid-cols-2 items-center gap-2"
                            :ui="{ container: '' }"
                        >
                            <UInput
                                type="text"
                                name="jobInfoUnemployedName"
                                :value="data.config!.jobInfo!.unemployedJob!.name"
                                :placeholder="$t('common.job')"
                                :label="$t('common.job')"
                                maxlength="128"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                        </UFormGroup>

                        <UFormGroup
                            name="jobInfoUnemployedGrade"
                            :label="$t('common.rank')"
                            class="grid grid-cols-2 items-center gap-2"
                            :ui="{ container: '' }"
                        >
                            <UInput
                                type="number"
                                min="1"
                                max="99"
                                :value="data.config!.jobInfo!.unemployedJob!.grade"
                                name="jobInfoUnemployedGrade"
                                :placeholder="$t('common.rank')"
                                :label="$t('common.rank')"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                        </UFormGroup>

                        <UFormGroup
                            name="jobInfoPublicJobs"
                            :label="$t('components.rector.app_config.job_info.public_jobs')"
                            class="grid grid-cols-2 items-center gap-2"
                            :ui="{ container: '' }"
                        >
                            <UFormGroup class="flex-1" name="job" :label="$t('common.job')">
                                <USelectMenu
                                    v-model="data.config!.jobInfo!.publicJobs"
                                    multiple
                                    :options="jobs"
                                    value-attribute="name"
                                    by="label"
                                >
                                    <template #label>
                                        <template v-if="data.config!.jobInfo!.publicJobs">
                                            <span class="truncate">{{ data.config!.jobInfo!.publicJobs.join(',') }}</span>
                                        </template>
                                    </template>
                                    <template #option="{ option: job }">
                                        <span class="truncate">{{ job.label }} ({{ job.name }})</span>
                                    </template>
                                </USelectMenu>
                            </UFormGroup>
                        </UFormGroup>

                        <UFormGroup
                            name="jobInfoHiddenJobs"
                            :label="$t('components.rector.app_config.job_info.hidden_jobs')"
                            class="grid grid-cols-2 items-center gap-2"
                            :ui="{ container: '' }"
                        >
                            <UFormGroup class="flex-1" name="job" :label="$t('common.job')">
                                <USelectMenu
                                    v-model="data.config!.jobInfo!.hiddenJobs"
                                    multiple
                                    :options="jobs"
                                    value-attribute="name"
                                    by="label"
                                >
                                    <template #label>
                                        <template v-if="data.config!.jobInfo!.hiddenJobs">
                                            <span class="truncate">{{ data.config!.jobInfo!.hiddenJobs.join(',') }}</span>
                                        </template>
                                    </template>
                                    <template #option="{ option: job }">
                                        <span class="truncate">{{ job.label }} ({{ job.name }})</span>
                                    </template>
                                </USelectMenu>
                            </UFormGroup>
                        </UFormGroup>
                    </UDashboardSection>

                    <UDashboardSection
                        :title="$t('components.rector.app_config.user_tracker.title')"
                        :description="$t('components.rector.app_config.user_tracker.description')"
                    >
                        <UFormGroup
                            name="userTrackerRefreshTime"
                            :label="$t('components.rector.app_config.user_tracker.refresh_time')"
                            class="grid grid-cols-2 items-center gap-2"
                            :ui="{ container: '' }"
                        >
                            <UInput
                                name="userTrackerRefreshTime"
                                type="text"
                                :value="
                                    parseFloat(
                                        data.config?.userTracker?.refreshTime?.seconds.toString() +
                                            '.' +
                                            (data.config?.userTracker?.refreshTime?.nanos ?? 0) / 1000000,
                                    ).toString() + 's'
                                "
                                :placeholder="$t('common.duration')"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                        </UFormGroup>

                        <UFormGroup
                            name="userTrackerDbRefreshTime"
                            :label="$t('components.rector.app_config.user_tracker.db_refresh_time')"
                            class="grid grid-cols-2 items-center gap-2"
                            :ui="{ container: '' }"
                        >
                            <UInput
                                name="userTrackerDbRefreshTime"
                                type="text"
                                :value="
                                    parseFloat(
                                        data.config?.userTracker?.dbRefreshTime?.seconds.toString() +
                                            '.' +
                                            (data.config?.userTracker?.dbRefreshTime?.nanos ?? 0) / 1000000,
                                    ).toString() + 's'
                                "
                                :placeholder="$t('common.duration')"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                        </UFormGroup>

                        <UFormGroup
                            name="livemapJobs"
                            :label="$t('components.rector.app_config.user_tracker.livemap_jobs')"
                            class="grid grid-cols-2 items-center gap-2"
                            :ui="{ container: '' }"
                        >
                            <USelectMenu
                                v-model="data.config!.userTracker!.livemapJobs"
                                multiple
                                :options="jobs"
                                value-attribute="name"
                                by="label"
                            >
                                <template #label>
                                    <template v-if="data.config!.userTracker!.livemapJobs">
                                        <span class="truncate">{{ data.config!.userTracker!.livemapJobs.join(',') }}</span>
                                    </template>
                                </template>
                                <template #option="{ option: job }">
                                    <span class="truncate">{{ job.label }} ({{ job.name }})</span>
                                </template>
                            </USelectMenu>
                        </UFormGroup>
                    </UDashboardSection>

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
                            <UToggle v-model="data.config!.discord!.enabled">
                                <span class="sr-only">
                                    {{ $t('common.enabled') }}
                                </span>
                            </UToggle>
                        </UFormGroup>

                        <UFormGroup
                            name="discordSyncInterval"
                            :label="$t('components.rector.app_config.discord.sync_interval')"
                            class="grid grid-cols-2 items-center gap-2"
                            :ui="{ container: '' }"
                        >
                            <UInput
                                name="discordSyncInterval"
                                type="text"
                                :value="
                                    parseFloat(
                                        data.config?.discord?.syncInterval?.seconds.toString() +
                                            '.' +
                                            (data.config?.discord?.syncInterval?.nanos ?? 0) / 1000000,
                                    ).toString() + 's'
                                "
                                :placeholder="$t('common.duration')"
                                :label="$t('common.duration')"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                        </UFormGroup>

                        <UFormGroup
                            name="discordBotInviteUrl"
                            :label="$t('components.rector.app_config.discord.bot_invite_url')"
                            class="grid grid-cols-2 items-center gap-2"
                            :ui="{ container: '' }"
                        >
                            <UInput
                                type="url"
                                name="discordBotInviteUrl"
                                :value="data.config!.discord!.inviteUrl"
                                :placeholder="$t('components.rector.app_config.discord.bot_invite_url')"
                                :label="$t('components.rector.app_config.discord.bot_invite_url')"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                        </UFormGroup>
                    </UDashboardSection>
                </UDashboardPanelContent>
            </template>
        </template>
    </div>
</template>
