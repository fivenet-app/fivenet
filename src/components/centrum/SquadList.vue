<script lang="ts" setup>
const squads = [
    { name: 'Charlie 2', initials: 'C2', members: 2, bgColor: 'bg-yellow-600', status: 'Available' },
    { name: 'Overwatch 1', initials: 'O1', members: 3, bgColor: 'bg-purple-600', status: 'Available' },
    { name: 'Delta 2', initials: 'D2', members: 2, bgColor: 'bg-pink-600', status: 'Dispatched' },
    { name: 'Delta 3', initials: 'D3', members: 2, bgColor: 'bg-pink-600', status: 'Dispatched' },
    { name: 'Charlie 1', initials: 'C1', members: 2, bgColor: 'bg-yellow-600', status: 'Break' },
    { name: 'Charlie 3', initials: 'C3', members: 2, bgColor: 'bg-yellow-600', status: 'Break' },
    { name: 'Delta 1', initials: 'D1', members: 2, bgColor: 'bg-pink-600', status: 'Break' },
];
</script>

<template>
    <div class="mx-2">
        <h2 class="text-sm font-medium text-gray-200">Active Squads (non-empty)</h2>
        <ul role="list" class="mt-3 grid grid-cols-1 gap-5 sm:grid-cols-2 sm:gap-2 lg:grid-cols-3">
            <li v-for="squad in squads" :key="squad.name" class="col-span-1 flex rounded-md shadow-sm">
                <div
                    :class="[
                        squad.bgColor,
                        'flex w-9 flex-shrink-0 items-center justify-center rounded-l-md text-sm font-medium text-white',
                    ]"
                >
                    {{ squad.initials }}
                </div>
                <div
                    class="flex flex-1 items-center justify-between truncate rounded-r-md border-b border-r border-t border-gray-200 bg-white"
                >
                    <div class="flex-1 truncate px-4 py-2 text-sm">
                        <span class="font-medium text-gray-900">{{ squad.name }}</span>
                        <p class="text-gray-500">{{ squad.members }} Members</p>
                    </div>
                    <div class="flex-shrink-0 pr-5">
                        <button type="button" class="inline-flex items-center justify-center text-white bg-green-700">
                            {{ squad.status }}
                        </button>
                    </div>
                </div>
            </li>
        </ul>
    </div>
</template>
