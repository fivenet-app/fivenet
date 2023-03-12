<script lang="ts">
import { User } from '@arpanet/gen/resources/users/users_pb';
import { OrderBy } from '@arpanet/gen/resources/common/database/database_pb';
import { defineComponent } from 'vue';
import { watchDebounced } from '@vueuse/core'
import { getUsersClient, handleGRPCError } from '../../grpc';
import { RpcError } from 'grpc-web';
import { FindUsersRequest } from '@arpanet/gen/services/users/users_pb';
import TablePagination from '../partials/TablePagination.vue';
import CitizenListEntry from './CitizensListEntry.vue';
import { Switch } from '@headlessui/vue';

export default defineComponent({
    components: {
        TablePagination,
        CitizenListEntry,
        Switch,
    },
    data() {
        return {
            loading: false,
            search: {
                name: '',
                wanted: false
            },
            orderBys: [] as Array<OrderBy>,
            offset: 0,
            totalCount: 0,
            listEnd: 0,
            users: [] as Array<User>,
        };
    },
    methods: {
        findUsers: function (offset: number) {
            if (offset < 0) {
                return;
            }
            if (this.loading) {
                return;
            }
            this.loading = true;

            const req = new FindUsersRequest();
            req.setOffset(offset);
            req.setSearchname(this.search.name);
            req.setWanted(this.search.wanted);
            req.setOrderbyList(this.orderBys);

            getUsersClient().
                findUsers(req, null).
                then((resp) => {
                    this.totalCount = resp.getTotalcount();
                    this.offset = resp.getOffset();
                    this.listEnd = resp.getEnd();
                    this.users = resp.getUsersList();

                    this.loading = false;
                }).
                catch((err: RpcError) => {
                    this.loading = false;
                    handleGRPCError(err, this.$route);
                });
        },
        toggleOrderBy: function (column: string): void {
            const index = this.orderBys.findIndex((o) => {
                return o.getColumn() == column;
            });
            let orderBy: OrderBy;
            if (index > -1) {
                //@ts-ignore I just checked if it exists, so it should exist
                orderBy = this.orderBys.at(index);
                if (orderBy.getDesc()) {
                    this.orderBys.splice(index);
                }
                else {
                    orderBy.setDesc(true);
                }
            }
            else {
                orderBy = new OrderBy();
                orderBy.setColumn(column);
                orderBy.setDesc(false);
                this.orderBys.push(orderBy);
            }
            this.findUsers(this.offset);
        },
    },
    mounted: function () {
        watchDebounced(this.search, () => this.findUsers(0), { debounce: 750, maxWait: 1000 });
    },
});
</script>

<template>
    <div class="py-2">
        <div class="px-2 sm:px-6 lg:px-8">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <form @submit.prevent="findUsers(0)">
                        <div class="grid grid-cols-5 gap-4">
                            <div class="col-span-4 form-control">
                                <label for="search" class="block text-sm font-medium leading-6 text-white">Name</label>
                                <div class="relative mt-2 flex items-center">
                                    <input v-model="search.name" v-on:keyup.enter="findUsers(0)" type="text"
                                        name="search" id="search"
                                        class="block w-full rounded-md border-0 py-1.5 pr-14 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6" />
                                </div>
                            </div>
                            <div class="form-control">
                                <label for="search" class="block text-sm font-medium leading-6 text-white">Only
                                    Wanted</label>
                                <div class="relative mt-2 flex items-center">
                                    <Switch v-model="search.wanted"
                                        :class="[search.wanted ? 'bg-indigo-600' : 'bg-gray-200', 'relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-indigo-600 focus:ring-offset-2']">
                                        <span class="sr-only">Wanted</span>
                                        <span aria-hidden="true"
                                            :class="[search.wanted ? 'translate-x-5' : 'translate-x-0', 'pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out']" />
                                    </Switch>
                                </div>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
            <div class="mt-2 flow-root">
                <div class="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-8">
                    <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
                        <table class="min-w-full divide-y divide-gray-700">
                            <thead>
                                <tr>
                                    <th v-on:click="toggleOrderBy('firstname')" scope="col"
                                        class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-white sm:pl-0">Name
                                    </th>
                                    <th v-on:click="toggleOrderBy('job')" scope="col"
                                        class="py-3.5 px-2 text-left text-sm font-semibold text-white">Job
                                    </th>
                                    <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-white">Sex
                                    </th>
                                    <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-white">Date of
                                        Birth
                                    </th>
                                    <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-white">Height
                                    </th>
                                    <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-0">
                                        <span class="sr-only">Edit</span>
                                    </th>
                                </tr>
                            </thead>
                            <tbody class="divide-y divide-gray-800">
                                <CitizenListEntry v-for="user in users" :key="user.getUserid()" :user="user" />
                            </tbody>
                            <thead>
                                <tr>
                                    <th v-on:click="toggleOrderBy('firstname')" scope="col"
                                        class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-white sm:pl-0">Name
                                    </th>
                                    <th v-on:click="toggleOrderBy('job')" scope="col"
                                        class="py-3.5 px-2 text-left text-sm font-semibold text-white">Job
                                    </th>
                                    <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-white">Sex
                                    </th>
                                    <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-white">Date of
                                        Birth
                                    </th>
                                    <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-white">Height
                                    </th>
                                    <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-0">
                                        <span class="sr-only">Edit</span>
                                    </th>
                                </tr>
                            </thead>
                        </table>

                        <TablePagination :offset="offset" :entries="users.length" :end="listEnd" :total="totalCount"
                            :callback="findUsers" />
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
