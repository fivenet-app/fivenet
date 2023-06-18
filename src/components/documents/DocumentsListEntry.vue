<script lang="ts" setup>
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiAccount, mdiBriefcase, mdiCalendar, mdiLock, mdiLockOpenVariant, mdiTrashCan, mdiUpdate } from '@mdi/js';
import { DocumentShort } from '~~/gen/ts/resources/documents/documents';
import IDCopyBadge from '../partials/IDCopyBadge.vue';

defineProps<{
    doc: DocumentShort;
}>();
</script>

<template>
    <li
        :key="doc.id?.toString()"
        :class="[
            doc.deletedAt ? 'hover:bg-warn-800 bg-warn-800' : 'hover:bg-base-800 bg-base-850',
            'flex-initial my-1 rounded-lg',
        ]"
    >
        <NuxtLink
            :to="{
                name: 'documents-id',
                params: { id: doc.id.toString() },
            }"
        >
            <div class="mx-2 mt-1 mb-4">
                <div class="flex flex-row">
                    <p class="py-2 pl-4 pr-3 text-lg font-medium text-neutral sm:pl-0 truncate max-w-3xl">
                        <span
                            v-if="doc.category"
                            class="inline-flex items-center rounded-md bg-primary-400/10 px-2 py-1 text-xs font-medium text-primary-400 ring-1 ring-inset ring-primary-400/30"
                        >
                            {{ doc.category.name }}
                        </span>
                        {{ doc.title }}
                    </p>
                    <p
                        class="inline-flex px-2 text-xs font-semibold leading-5 rounded-full bg-primary-100 text-primary-700 my-auto"
                        v-if="doc.state"
                    >
                        {{ doc.state }}
                    </p>
                    <div class="flex flex-row items-center justify-end flex-1 text-base-200">
                        <div v-if="doc?.closed" class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-error-100">
                            <SvgIcon class="w-5 h-5 text-error-400" aria-hidden="true" type="mdi" :path="mdiLock" />
                            <span class="text-sm font-medium text-error-700">
                                {{ $t('common.close', 2) }}
                            </span>
                        </div>
                        <div v-else class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-success-100">
                            <SvgIcon class="w-5 h-5 text-green-500" aria-hidden="true" type="mdi" :path="mdiLockOpenVariant" />
                            <span class="text-sm font-medium text-green-700">
                                {{ $t('common.open') }}
                            </span>
                        </div>
                    </div>
                </div>
                <div class="flex flex-row gap-2 text-base-300 truncate max-w-5xl">
                    <div class="flex flex-row items-center justify-start flex-1">
                        <IDCopyBadge :id="doc.id" prefix="DOC" />
                    </div>
                    <div v-if="doc.deletedAt" class="flex flex-row items-center justify-center flex-1 text-base-100">
                        <SvgIcon
                            class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400"
                            aria-hidden="true"
                            type="mdi"
                            :path="mdiTrashCan"
                        />
                        {{ $t('common.deleted') }}
                    </div>
                    <div class="flex flex-row items-center justify-end flex-1">
                        <SvgIcon
                            class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400"
                            aria-hidden="true"
                            type="mdi"
                            :path="mdiUpdate"
                        />
                        <p>
                            {{ $t('common.updated_at') }}
                            <time :datetime="$d(toDate(doc.deletedAt)!, 'short')">
                                {{ useLocaleTimeAgo(toDate(doc.deletedAt)!).value }}
                            </time>
                        </p>
                    </div>
                </div>
                <div class="flex flex-row gap-2 text-base-200">
                    <div class="flex flex-row items-center justify-start flex-1">
                        <SvgIcon
                            class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400"
                            aria-hidden="true"
                            type="mdi"
                            :path="mdiAccount"
                        />
                        {{ doc.creator?.firstname }},
                        {{ doc.creator?.lastname }}
                    </div>
                    <div class="flex flex-row items-center justify-center flex-1">
                        <SvgIcon
                            class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400"
                            aria-hidden="true"
                            type="mdi"
                            :path="mdiBriefcase"
                        />
                        {{ doc.creator?.jobLabel }}
                    </div>
                    <div class="flex flex-row items-center justify-end flex-1">
                        <SvgIcon
                            class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400"
                            aria-hidden="true"
                            type="mdi"
                            :path="mdiCalendar"
                        />
                        <p>
                            {{ $t('common.created_at') }}
                            <time :datetime="$d(toDate(doc.createdAt)!, 'short')">
                                {{ $d(toDate(doc.createdAt)!, 'short') }}
                            </time>
                        </p>
                    </div>
                </div>
            </div>
        </NuxtLink>
    </li>
</template>
