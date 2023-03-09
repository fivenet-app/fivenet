<script lang="ts">
import { defineComponent } from 'vue';
import { Document, GetDocumentRequest } from '@arpanet/gen/documents/documents_pb';
import { getDocumentsClient, handleGRPCError } from '../../grpc';
import { RpcError } from 'grpc-web';
import { getDateLocaleString } from '../../utils/time';
import {
    Menu,
    MenuButton,
    MenuItem,
    MenuItems,
} from '@headlessui/vue';
import {
    ArchiveBoxIcon,
    ArrowUturnLeftIcon,
    ChevronDownIcon,
    ChevronUpIcon,
    EllipsisVerticalIcon,
    FolderArrowDownIcon,
    PencilIcon,
    UserPlusIcon,
} from '@heroicons/vue/20/solid';

export default defineComponent({
    components: {
        Menu,
        MenuButton,
        MenuItem,
        MenuItems,
        ArchiveBoxIcon,
        ArrowUturnLeftIcon,
        ChevronDownIcon,
        ChevronUpIcon,
        EllipsisVerticalIcon,
        FolderArrowDownIcon,
        PencilIcon,
        UserPlusIcon,
    },
    data() {
        return {
            document: undefined as undefined | Document,
            responses: [] as Array<Document>,
        };
    },
    props: {
        documentID: {
            required: true,
            type: Number,
        },
    },
    mounted() {
        this.getDocument();
    },
    methods: {
        getDocument(): void {
            const req = new GetDocumentRequest();
            req.setId(this.documentID);
            req.setResponses(true);

            getDocumentsClient().
                getDocument(req, null).
                then((resp) => {
                    this.document = resp.getDocument();
                    this.responses = resp.getResponsesList();
                }).
                catch((err: RpcError) => {
                    handleGRPCError(err, this.$route);
                });
        },
        getDateLocaleString,
    },
});
</script>

<template>
    <div class="flex h-full flex-col">
        <!-- Main area -->
        <main class="min-w-0 flex-1 border-t border-gray-200 xl:flex">
            <section aria-labelledby="message-heading"
                class="flex h-full min-w-0 flex-1 flex-col overflow-hidden xl:order-last">
                <!-- Top section -->
                <div class="flex-shrink-0 border-b border-gray-200 bg-white">
                    <!-- Toolbar-->
                    <div class="flex h-16 flex-col justify-center">
                        <div class="px-4 sm:px-6 lg:px-8">
                            <div class="flex justify-between py-3">
                                <!-- Left buttons -->
                                <div>
                                    <div class="isolate inline-flex rounded-md shadow-sm sm:space-x-3 sm:shadow-none">
                                        <span class="inline-flex sm:shadow-sm">
                                            <button type="button"
                                                class="relative inline-flex items-center gap-x-1.5 rounded-l-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 ring-1 ring-inset ring-gray-300 hover:z-10 hover:bg-gray-50 focus:z-10">
                                                <ArrowUturnLeftIcon class="-ml-0.5 h-5 w-5 text-gray-400"
                                                    aria-hidden="true" />
                                                Reply
                                            </button>
                                            <button type="button"
                                                class="relative -ml-px hidden items-center gap-x-1.5 bg-white px-3 py-2 text-sm font-semibold text-gray-900 ring-1 ring-inset ring-gray-300 hover:z-10 hover:bg-gray-50 focus:z-10 sm:inline-flex">
                                                <PencilIcon class="-ml-0.5 h-5 w-5 text-gray-400" aria-hidden="true" />
                                                Note
                                            </button>
                                            <button type="button"
                                                class="relative -ml-px hidden items-center gap-x-1.5 rounded-r-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 ring-1 ring-inset ring-gray-300 hover:z-10 hover:bg-gray-50 focus:z-10 sm:inline-flex">
                                                <UserPlusIcon class="-ml-0.5 h-5 w-5 text-gray-400" aria-hidden="true" />
                                                Assign
                                            </button>
                                        </span>

                                        <span class="hidden space-x-3 lg:flex">
                                            <button type="button"
                                                class="relative -ml-px hidden items-center gap-x-1.5 rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 ring-1 ring-inset ring-gray-300 hover:z-10 hover:bg-gray-50 focus:z-10 sm:inline-flex">
                                                <ArchiveBoxIcon class="-ml-0.5 h-5 w-5 text-gray-400" aria-hidden="true" />
                                                Archive
                                            </button>
                                            <button type="button"
                                                class="relative -ml-px hidden items-center gap-x-1.5 rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 ring-1 ring-inset ring-gray-300 hover:z-10 hover:bg-gray-50 focus:z-10 sm:inline-flex">
                                                <FolderArrowDownIcon class="-ml-0.5 h-5 w-5 text-gray-400"
                                                    aria-hidden="true" />
                                                Move
                                            </button>
                                        </span>

                                        <Menu as="div" class="relative -ml-px block sm:shadow-sm lg:hidden">
                                            <div>
                                                <MenuButton
                                                    class="relative inline-flex items-center gap-x-1.5 rounded-r-md bg-white px-2 py-2 text-sm font-semibold text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-10 sm:rounded-md sm:px-3">
                                                    <span class="sr-only sm:hidden">More</span>
                                                    <span class="hidden sm:inline">More</span>
                                                    <ChevronDownIcon class="h-5 w-5 text-gray-400 sm:-mr-1"
                                                        aria-hidden="true" />
                                                </MenuButton>
                                            </div>

                                            <transition enter-active-class="transition ease-out duration-100"
                                                enter-from-class="transform opacity-0 scale-95"
                                                enter-to-class="transform opacity-100 scale-100"
                                                leave-active-class="transition ease-in duration-75"
                                                leave-from-class="transform opacity-100 scale-100"
                                                leave-to-class="transform opacity-0 scale-95">
                                                <MenuItems
                                                    class="absolute right-0 z-10 mt-2 w-36 origin-top-right rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                                                    <div class="py-1">
                                                        <MenuItem v-slot="{ active }">
                                                        <a href="#"
                                                            :class="[active ? 'bg-gray-100 text-gray-900' : 'text-gray-700', 'block px-4 py-2 text-sm sm:hidden']">Note</a>
                                                        </MenuItem>
                                                        <MenuItem v-slot="{ active }">
                                                        <a href="#"
                                                            :class="[active ? 'bg-gray-100 text-gray-900' : 'text-gray-700', 'block px-4 py-2 text-sm sm:hidden']">Assign</a>
                                                        </MenuItem>
                                                        <MenuItem v-slot="{ active }">
                                                        <a href="#"
                                                            :class="[active ? 'bg-gray-100 text-gray-900' : 'text-gray-700', 'block px-4 py-2 text-sm']">Archive</a>
                                                        </MenuItem>
                                                        <MenuItem v-slot="{ active }">
                                                        <a href="#"
                                                            :class="[active ? 'bg-gray-100 text-gray-900' : 'text-gray-700', 'block px-4 py-2 text-sm']">Move</a>
                                                        </MenuItem>
                                                    </div>
                                                </MenuItems>
                                            </transition>
                                        </Menu>
                                    </div>
                                </div>

                                <!-- Right buttons -->
                                <nav aria-label="Pagination">
                                    <span class="isolate inline-flex rounded-md shadow-sm">
                                        <a href="#"
                                            class="relative inline-flex items-center rounded-l-md bg-white px-3 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:z-10 hover:bg-gray-50 focus:z-10">
                                            <span class="sr-only">Next</span>
                                            <ChevronUpIcon class="h-5 w-5" aria-hidden="true" />
                                        </a>
                                        <a href="#"
                                            class="relative -ml-px inline-flex items-center rounded-r-md bg-white px-3 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:z-10 hover:bg-gray-50 focus:z-10">
                                            <span class="sr-only">Previous</span>
                                            <ChevronDownIcon class="h-5 w-5" aria-hidden="true" />
                                        </a>
                                    </span>
                                </nav>
                            </div>
                        </div>
                    </div>
                    <!-- Message header -->
                </div>

                <div class="min-h-0 flex-1 overflow-y-auto">
                    <div class="bg-white pt-5 pb-6 shadow">
                        <div class="px-4 sm:flex sm:items-baseline sm:justify-between sm:px-6 lg:px-8">
                            <div class="sm:w-0 sm:flex-1">
                                <h1 id="message-heading" class="text-lg font-medium text-gray-900">{{ document?.getTitle()
                                }}
                                </h1>
                                <p class="mt-1 truncate text-sm text-gray-500">{{ document?.getCreator() }}</p>
                            </div>

                            <div
                                class="mt-4 flex items-center justify-between sm:mt-0 sm:ml-6 sm:flex-shrink-0 sm:justify-start">
                                <span
                                    class="inline-flex items-center rounded-full bg-cyan-100 px-3 py-0.5 text-sm font-medium text-cyan-800">{{
                                        document?.getState() }}</span>
                                <Menu as="div" class="relative ml-3 inline-block text-left">
                                    <div>
                                        <MenuButton
                                            class="-my-2 flex items-center rounded-full bg-white p-2 text-gray-400 hover:text-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-600">
                                            <span class="sr-only">Open options</span>
                                            <EllipsisVerticalIcon class="h-5 w-5" aria-hidden="true" />
                                        </MenuButton>
                                    </div>

                                    <transition enter-active-class="transition ease-out duration-100"
                                        enter-from-class="transform opacity-0 scale-95"
                                        enter-to-class="transform opacity-100 scale-100"
                                        leave-active-class="transition ease-in duration-75"
                                        leave-from-class="transform opacity-100 scale-100"
                                        leave-to-class="transform opacity-0 scale-95">
                                        <MenuItems
                                            class="absolute right-0 z-10 mt-2 w-56 origin-top-right rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                                            <div class="py-1">
                                                <MenuItem v-slot="{ active }">
                                                <button type="button"
                                                    :class="[active ? 'bg-gray-100 text-gray-900' : 'text-gray-700', 'flex w-full justify-between px-4 py-2 text-sm']">
                                                    <span>Copy email address</span>
                                                </button>
                                                </MenuItem>
                                                <MenuItem v-slot="{ active }">
                                                <a href="#"
                                                    :class="[active ? 'bg-gray-100 text-gray-900' : 'text-gray-700', 'flex justify-between px-4 py-2 text-sm']">
                                                    <span>Previous conversations</span>
                                                </a>
                                                </MenuItem>
                                                <MenuItem v-slot="{ active }">
                                                <a href="#"
                                                    :class="[active ? 'bg-gray-100 text-gray-900' : 'text-gray-700', 'flex justify-between px-4 py-2 text-sm']">
                                                    <span>View original</span>
                                                </a>
                                                </MenuItem>
                                            </div>
                                        </MenuItems>
                                    </transition>
                                </Menu>
                            </div>
                        </div>
                    </div>
                </div>
            </section>

            <!-- Message list-->
            <aside class="hidden xl:order-first xl:block xl:flex-shrink-0">
                <div class="relative flex h-full w-96 flex-col border-r border-gray-200 bg-gray-100">
                    <nav aria-label="Message list" class="min-h-0 flex-1 overflow-y-auto">
                        <ul role="list" class="divide-y divide-gray-200 border-b border-gray-200">
                            <li v-for="response in responses" :key="response.getId()"
                                class="relative bg-white py-5 px-6 focus-within:ring-2 focus-within:ring-inset focus-within:ring-blue-600 hover:bg-gray-50">
                                <div class="flex justify-between space-x-3">
                                    <div class="min-w-0 flex-1">
                                        <a href="#" class="block focus:outline-none">
                                            <span class="absolute inset-0" aria-hidden="true" />
                                            <p class="truncate text-sm font-medium text-gray-900">{{ response.getCreator()
                                            }}
                                            </p>
                                            <p class="truncate text-sm text-gray-500">{{ response.getTitle() }}</p>
                                        </a>
                                    </div>
                                    <time :datetime="getDateLocaleString(response.getCreatedAt())"
                                        class="flex-shrink-0 whitespace-nowrap text-sm text-gray-500">{{
                                            getDateLocaleString(response.getCreatedAt()) }}</time>
                                </div>
                                <div class="mt-1">
                                    <p class="text-sm text-gray-600 line-clamp-2">{{ response.getContent() }}</p>
                                </div>
                            </li>
                        </ul>
                    </nav>
                </div>
            </aside>
        </main>
    </div>
</template>
