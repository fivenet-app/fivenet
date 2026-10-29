<script lang="ts" setup>
import FileGallery from '~/components/partials/files/FileGallery.vue';
import ConfirmModalWithReason from '~/components/partials/ConfirmModalWithReason.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import UploadFileReasonDrawer from '~/components/partials/files/UploadFileReasonDrawer.vue';
import { getCitizensCitizensClient } from '~~/gen/ts/clients';
import type { File as FilestoreFile } from '~~/gen/ts/resources/file/file';
import type { CitizenFile } from '~~/gen/ts/services/citizens/citizens';
import type { User } from '~~/gen/ts/resources/users/user';

const user = defineModel<User>('user', { required: true });
const route = useRoute('citizens-id-files');
const { can, attr } = useAuth();
const { t } = useI18n();
const client = await getCitizensCitizensClient();
const overlay = useOverlay();
const deleteConfirmModal = overlay.create(ConfirmModalWithReason);
const uploadDrawer = overlay.create(UploadFileReasonDrawer);

const { data, refresh } = useAuthedLazyAsyncData('userState', `citizen-files-${route.params.id}`, async () => {
    const { response } = await client.listCitizenFiles({ userId: Number(route.params.id) });
    return response.files;
});
const files = computed<CitizenFile[]>(() => data.value ?? []);
const canUpload = attr('citizens.CitizensService/SetUserProps', 'Properties', 'user.files');
const canDelete = computed(
    () =>
        can('citizens.CitizensService/DeleteCitizenFile').value &&
        attr('citizens.CitizensService/SetUserProps', 'Properties', 'user.files').value,
);

function openUploadDrawer(): void {
    uploadDrawer.open({
        namespace: 'citizen-files',
        parentId: user.value.userId,
        uploadFunc: (opts) => client.uploadCitizenFile(opts),
        onUploaded: () => refresh(),
    });
}

function deleteFile(file: FilestoreFile): void {
    deleteConfirmModal.open({
        title: t('common.delete'),
        confirm: async (reason) => {
            try {
                await client.deleteCitizenFile({ userId: user.value.userId, fileId: file.id, reason });
                await refresh();
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },
    });
}
</script>

<template>
    <UDashboardPanel :ui="{ root: 'min-h-0', body: 'p-4 sm:p-6' }">
        <template #header>
            <UDashboardToolbar>
                <div class="my-2 flex w-full items-center justify-between gap-2">
                    <UButton
                        v-if="canUpload"
                        :disabled="files.length >= 6"
                        icon="i-mdi-upload"
                        :label="$t('common.upload')"
                        @click="openUploadDrawer"
                    />
                    <UBadge :label="`${files.length} / 6`" color="neutral" variant="subtle" />
                </div>
            </UDashboardToolbar>
        </template>

        <template #body>
            <FileGallery v-if="files.length" :files="files" :can-delete="canDelete" @delete="deleteFile" />
            <DataNoDataBlock v-else :message="$t('common.not_found', [$t('common.file', 2)])" />
        </template>
    </UDashboardPanel>
</template>
