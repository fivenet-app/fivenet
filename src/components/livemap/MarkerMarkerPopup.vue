<script lang="ts" setup>
import { LPopup } from '@vue-leaflet/vue-leaflet';
import { TrashCanIcon } from 'mdi-vue3';
import type { Marker } from '~~/gen/ts/resources/livemap/livemap';
import PhoneNumber from '~/components/partials/citizens/PhoneNumber.vue';

defineProps<{
    marker: Marker;
}>();

defineEmits<{
    (e: 'delete'): void;
}>();
</script>

<template>
    <LPopup :options="{ closeButton: true }">
        <ul>
            <li class="inline-flex items-center">
                {{ marker.info?.name }}
                <template v-if="can('LivemapperService.DeleteMarker')">
                    <button
                        type="button"
                        :title="$t('common.delete')"
                        class="flex flex-row items-center"
                        @click="$emit('delete')"
                    >
                        <TrashCanIcon class="w-6 h-6" />
                        <span class="sr-only">{{ $t('common.delete') }}</span>
                    </button>
                </template>
            </li>
            <li v-if="marker.info?.description">{{ $t('common.description') }}: {{ marker.info?.description }}</li>
            <li class="italic">
                <span class="font-semibold">{{ $t('common.sent_by') }}</span
                >:
                <span v-if="marker.creator">
                    {{ marker.creator?.firstname }}, {{ marker.creator?.lastname }} (<PhoneNumber
                        :number="marker.creator.phoneNumber"
                    />)
                </span>
                <span v-else>
                    {{ $t('common.unknown') }}
                </span>
            </li>
        </ul>
    </LPopup>
</template>
