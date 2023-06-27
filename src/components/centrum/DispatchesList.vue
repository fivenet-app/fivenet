<script lang="ts" setup>
const transactions = [
    {
        id: 'AAPS0L',
        company: 'Chase & Co.',
        share: 'CAC',
        commission: '+$4.37',
        price: '$3,509.00',
        quantity: '12.00',
        netAmount: '$4,397.00',
    },
];
</script>

<template>
    <div class="px-4 sm:px-6 lg:px-8">
        <div class="sm:flex sm:items-center">
            <div class="sm:flex-auto">
                <h1 class="text-base font-semibold leading-6 text-gray-100">Dispatches</h1>
            </div>
        </div>
        <div class="mt-8 flow-root">
            <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
                    <table class="min-w-full divide-y divide-gray-300">
                        <thead>
                            <tr>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-100 sm:pl-0"
                                >
                                    Unit
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-2 py-3.5 text-left text-sm font-semibold text-gray-100"
                                >
                                    Company
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-2 py-3.5 text-left text-sm font-semibold text-gray-100"
                                >
                                    Share
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-2 py-3.5 text-left text-sm font-semibold text-gray-100"
                                >
                                    Commision
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-2 py-3.5 text-left text-sm font-semibold text-gray-100"
                                >
                                    Price
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-2 py-3.5 text-left text-sm font-semibold text-gray-100"
                                >
                                    Quantity
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-2 py-3.5 text-left text-sm font-semibold text-gray-100"
                                >
                                    Net amount
                                </th>
                                <th scope="col" class="relative whitespace-nowrap py-3.5 pl-3 pr-4 sm:pr-0">
                                    <span class="sr-only">Edit</span>
                                </th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-gray-200">
                            <tr v-for="transaction in transactions" :key="transaction.id">
                                <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm text-gray-300 sm:pl-0">
                                    {{ transaction.id }}
                                </td>
                                <td class="whitespace-nowrap px-2 py-2 text-sm font-medium text-gray-100">
                                    {{ transaction.company }}
                                </td>
                                <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-100">{{ transaction.share }}</td>
                                <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-300">{{ transaction.commission }}</td>
                                <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-300">{{ transaction.price }}</td>
                                <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-300">{{ transaction.quantity }}</td>
                                <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-300">{{ transaction.netAmount }}</td>
                                <td class="relative whitespace-nowrap py-2 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
                                    <a href="#" class="text-primary-600 hover:text-primary-900"
                                        >Edit<span class="sr-only">, {{ transaction.id }}</span></a
                                    >
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>
</template>
