<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { OfficeBuildingCogIcon } from 'mdi-vue3';
import { useSettingsStore } from '~/store/settings';
import GenericContainerPanel from '~/components/partials/elements/GenericContainerPanel.vue';
import GenericContainerPanelEntry from '~/components/partials/elements/GenericContainerPanelEntry.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { type GetAppConfigResponse } from '~~/gen/ts/services/rector/config';

const { $grpc } = useNuxtApp();

const settingsStore = useSettingsStore();
const { streamerMode } = storeToRefs(settingsStore);

const { data, pending, refresh, error } = useLazyAsyncData(`rector-jobprops`, () => getAppConfig());

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

// TODO
</script>

<template>
    <div class="mx-auto max-w-5xl py-2">
        <template v-if="streamerMode">
            <GenericContainerPanel>
                <template #title>
                    {{ $t('system.streamer_mode.title') }}
                </template>
                <template #description>
                    {{ $t('system.streamer_mode.description') }}
                </template>
            </GenericContainerPanel>
        </template>
        <template v-else>
            <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.setting', 2)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.setting', 2)])"
                :retry="refresh"
            />
            <DataNoDataBlock v-else-if="data === null" :icon="OfficeBuildingCogIcon" :type="$t('common.setting', 2)" />

            <template v-else>
                <GenericContainerPanel>
                    <template #title> Auth </template>
                    <template #default>
                        <GenericContainerPanelEntry>
                            <template #title>Sign-up Enabled</template>
                            <template #default>
                                {{ jsonStringify(data.config?.auth?.signupEnabled) }}
                            </template>
                        </GenericContainerPanelEntry>
                    </template>
                </GenericContainerPanel>
                <GenericContainerPanel>
                    <template #title> Website </template>
                    <template #default>
                        <GenericContainerPanelEntry>
                            <template #title>Links</template>
                            <template #default>
                                {{ jsonStringify(data.config?.website?.links) }}
                            </template>
                        </GenericContainerPanelEntry>
                    </template>
                </GenericContainerPanel>
                <GenericContainerPanel>
                    <template #title> Job Info </template>
                    <template #default>
                        <GenericContainerPanelEntry>
                            <template #title>Public Jobs</template>
                            <template #default>
                                {{ jsonStringify(data.config?.jobInfo?.publicJobs) }}
                            </template>
                        </GenericContainerPanelEntry>
                        <GenericContainerPanelEntry>
                            <template #title>Hidden Jobs</template>
                            <template #default>
                                {{ jsonStringify(data.config?.jobInfo?.hiddenJobs) }}
                            </template>
                        </GenericContainerPanelEntry>
                    </template>
                </GenericContainerPanel>
                <GenericContainerPanel>
                    <template #title> User Tracker / Livemap </template>
                    <template #default>
                        <GenericContainerPanelEntry>
                            <template #title>Refresh Times</template>
                            <template #default>
                                {{ jsonStringify(data.config?.userTracker?.refreshTime) }}
                            </template>
                        </GenericContainerPanelEntry>
                        <GenericContainerPanelEntry>
                            <template #title>DB Refresh Times</template>
                            <template #default>
                                {{ jsonStringify(data.config?.userTracker?.dbRefreshTime) }}
                            </template>
                        </GenericContainerPanelEntry>
                        <GenericContainerPanelEntry>
                            <template #title>Livemap Jobs</template>
                            <template #default>
                                {{ jsonStringify(data.config?.userTracker?.livemapJobs) }}
                            </template>
                        </GenericContainerPanelEntry>
                        <GenericContainerPanelEntry>
                            <template #title>Timeclock Jobs</template>
                            <template #default>
                                {{ jsonStringify(data.config?.userTracker?.timeclockJobs) }}
                            </template>
                        </GenericContainerPanelEntry>
                    </template>
                </GenericContainerPanel>
                <GenericContainerPanel>
                    <template #title> Discord </template>
                    <template #default>
                        <GenericContainerPanelEntry>
                            <template #title>Enabled</template>
                            <template #default>
                                {{ jsonStringify(data.config?.discord?.enabled) }}
                            </template>
                        </GenericContainerPanelEntry>
                    </template>
                </GenericContainerPanel>

                <!-- TODO -->
            </template>
        </template>
    </div>
</template>
