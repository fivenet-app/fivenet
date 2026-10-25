<script lang="ts" setup>
import FileGallery from '~/components/partials/files/FileGallery.vue';
import ConfirmModalWithReason from '~/components/partials/ConfirmModalWithReason.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import UploadFileReasonDrawer from '~/components/partials/files/UploadFileReasonDrawer.vue';
import type { Vehicle } from '~~/gen/ts/resources/vehicles/vehicles';
import type { File as FilestoreFile } from '~~/gen/ts/resources/file/file';
import { getVehiclesVehiclesClient } from '~~/gen/ts/clients';

const vehicle = defineModel<Vehicle>({ required: true });
const maxFiles = 6;
const { can, attr } = useAuth();
const { t } = useI18n();
const client = await getVehiclesVehiclesClient();
const overlay = useOverlay();
const uploadDrawer = overlay.create(UploadFileReasonDrawer);
const confirmModal = overlay.create(ConfirmModalWithReason);
const { data, refresh, error } = useAuthedLazyAsyncData('userState', `vehicle-files-${vehicle.value.plate}`, async () => {
    try {
        const { response } = await client.listVehicleFiles({ plate: vehicle.value.plate });
        return response.files;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
});
const files = computed(() => data.value ?? []);
const canUpload = attr('vehicles.VehiclesService/SetVehicleProps', 'Properties', 'vehicle.files');
const canDelete = computed(
    () =>
        can('vehicles.VehiclesService/DeleteVehicleFile').value &&
        attr('vehicles.VehiclesService/SetVehicleProps', 'Properties', 'vehicle.files').value,
);

function openUpload(): void {
    uploadDrawer.open({
        namespace: 'vehicle-files',
        parentId: vehicle.value.plate,
        uploadFunc: (opts) => client.uploadVehicleFile(opts),
        onUploaded: () => refresh(),
    });
}
function deleteFile(file: FilestoreFile): void {
    confirmModal.open({
        title: t('common.delete'),
        confirm: async (reason) => {
            await client.deleteVehicleFile({ plate: vehicle.value.plate, fileId: file.id, reason });
            await refresh();
        },
    });
}
</script>

<template>
    <UDashboardPanel :ui="{ root: 'min-h-0 h-full', body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardToolbar>
                <div class="my-2 flex w-full items-center justify-between gap-2">
                    <UButton
                        v-if="canUpload"
                        :disabled="files.length >= maxFiles"
                        icon="i-mdi-upload"
                        :label="$t('common.upload')"
                        @click="openUpload"
                    />
                    <UBadge :label="`${files.length} / ${maxFiles}`" color="neutral" variant="subtle" />
                </div>
            </UDashboardToolbar>
        </template>

        <template #body>
            <DataErrorBlock
                v-if="error"
                :title="$t('common.not_found', [`${$t('common.vehicle', 1)} ${$t('common.file', 2)}`])"
                :error="error"
                :retry="refresh"
            />
            <div class="p-4 sm:p-6">
                <FileGallery v-if="!error && files.length" :files="files" :can-delete="canDelete" @delete="deleteFile" />
                <DataNoDataBlock v-else-if="!error" :message="$t('common.not_found', [$t('common.file', 2)])" />
            </div>
        </template>
    </UDashboardPanel>
</template>
