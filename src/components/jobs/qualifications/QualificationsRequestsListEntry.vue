<script lang="ts" setup>
import { ListStatusIcon } from 'mdi-vue3';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { QualificationRequest, RequestStatus } from '~~/gen/ts/resources/qualifications/qualifications';

defineProps<{
    request: QualificationRequest;
}>();

defineEmits<{
    (e: 'selected'): void;
}>();
</script>

<template>
    <tr>
        <td>{{ request.qualification?.abbreviation }}: {{ request.qualification?.title }}</td>
        <td>
            <span v-if="request.userComment">{{ request.userComment }}</span>
        </td>
        <td>
            <div class="flex flex-initial flex-row gap-1 rounded-full px-2 py-1">
                <ListStatusIcon class="h-5 w-5 text-info-400" aria-hidden="true" />
                <template v-if="request.status !== undefined">
                    <span class="text-sm font-medium text-info-400">
                        <span class="font-semibold">{{
                            $t(`enums.qualifications.RequestStatus.${RequestStatus[request.status]}`)
                        }}</span>
                    </span>
                </template>
            </div>
        </td>
        <td>
            <p v-if="request.createdAt" class="mt-1 text-sm leading-5 text-gray-300">
                {{ $t('common.created_at') }} <GenericTime :value="request.createdAt" />
            </p>
        </td>
        <td>APPROVE / DENY</td>
    </tr>
</template>
