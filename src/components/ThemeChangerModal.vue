<template>
    <!-- Theme Changer Modal -->
    <input type="checkbox" id="modal-navbar-changetheme" class="modal-toggle" />
    <label for="modal-navbar-changetheme" class="modal cursor-pointer">
        <label class="modal-box w-11/12 max-w-5xl" for="">
            <h2 class="text-lg font-bold">Change Theme</h2>
            <div class="rounded-box grid grid-cols-3 gap-4 sm:grid-cols-3 md:grid-cols-3 lg:grid-cols-3">
                <div class="border-base-content/20 hover:border-base-content/40 outline-base-content overflow-hidden rounded-lg border outline-2 outline-offset-2"
                    data-set-theme="light" data-act-class="outline">
                    <div data-theme="light" class="bg-base-100 text-base-content w-full cursor-pointer font-sans">
                        <div class="grid grid-cols-5 grid-rows-3">
                            <div class="bg-base-200 col-start-1 row-span-2 row-start-1"></div>
                            <div class="bg-base-300 col-start-1 row-start-3"></div>
                            <div class="bg-base-100 col-span-4 col-start-2 row-span-3 row-start-1 flex flex-col gap-1 p-2">
                                <div class="font-bold">light</div>
                                <div class="flex flex-wrap gap-1">
                                    <div
                                        class="bg-primary flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-primary-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-secondary flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-secondary-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-accent flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-accent-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-neutral flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-neutral-content text-sm font-bold">A</div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="border-base-content/20 hover:border-base-content/40 outline-base-content overflow-hidden rounded-lg border outline-2 outline-offset-2"
                    data-set-theme="dark" data-act-class="outline">
                    <div data-theme="dark" class="bg-base-100 text-base-content w-full cursor-pointer font-sans">
                        <div class="grid grid-cols-5 grid-rows-3">
                            <div class="bg-base-200 col-start-1 row-span-2 row-start-1"></div>
                            <div class="bg-base-300 col-start-1 row-start-3"></div>
                            <div class="bg-base-100 col-span-4 col-start-2 row-span-3 row-start-1 flex flex-col gap-1 p-2">
                                <div class="font-bold">dark</div>
                                <div class="flex flex-wrap gap-1">
                                    <div
                                        class="bg-primary flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-primary-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-secondary flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-secondary-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-accent flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-accent-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-neutral flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-neutral-content text-sm font-bold">A</div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="border-base-content/20 hover:border-base-content/40 outline-base-content overflow-hidden rounded-lg border outline-2 outline-offset-2"
                    data-set-theme="corporate" data-act-class="outline">
                    <div data-theme="corporate" class="bg-base-100 text-base-content w-full cursor-pointer font-sans">
                        <div class="grid grid-cols-5 grid-rows-3">
                            <div class="bg-base-200 col-start-1 row-span-2 row-start-1"></div>
                            <div class="bg-base-300 col-start-1 row-start-3"></div>
                            <div class="bg-base-100 col-span-4 col-start-2 row-span-3 row-start-1 flex flex-col gap-1 p-2">
                                <div class="font-bold">corporate</div>
                                <div class="flex flex-wrap gap-1">
                                    <div
                                        class="bg-primary flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-primary-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-secondary flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-secondary-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-accent flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-accent-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-neutral flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-neutral-content text-sm font-bold">A</div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="border-base-content/20 hover:border-base-content/40 outline-base-content overflow-hidden rounded-lg border outline-2 outline-offset-2"
                    data-set-theme="business" data-act-class="outline">
                    <div data-theme="business" class="bg-base-100 text-base-content w-full cursor-pointer font-sans">
                        <div class="grid grid-cols-5 grid-rows-3">
                            <div class="bg-base-200 col-start-1 row-span-2 row-start-1"></div>
                            <div class="bg-base-300 col-start-1 row-start-3"></div>
                            <div class="bg-base-100 col-span-4 col-start-2 row-span-3 row-start-1 flex flex-col gap-1 p-2">
                                <div class="font-bold">business</div>
                                <div class="flex flex-wrap gap-1">
                                    <div
                                        class="bg-primary flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-primary-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-secondary flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-secondary-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-accent flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-accent-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-neutral flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-neutral-content text-sm font-bold">A</div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="border-base-content/20 hover:border-base-content/40 outline-base-content overflow-hidden rounded-lg border outline-2 outline-offset-2"
                    data-set-theme="acid" data-act-class="outline">
                    <div data-theme="acid" class="bg-base-100 text-base-content w-full cursor-pointer font-sans">
                        <div class="grid grid-cols-5 grid-rows-3">
                            <div class="bg-base-200 col-start-1 row-span-2 row-start-1"></div>
                            <div class="bg-base-300 col-start-1 row-start-3"></div>
                            <div class="bg-base-100 col-span-4 col-start-2 row-span-3 row-start-1 flex flex-col gap-1 p-2">
                                <div class="font-bold">acid</div>
                                <div class="flex flex-wrap gap-1">
                                    <div
                                        class="bg-primary flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-primary-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-secondary flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-secondary-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-accent flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-accent-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-neutral flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-neutral-content text-sm font-bold">A</div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="border-base-content/20 hover:border-base-content/40 outline-base-content overflow-hidden rounded-lg border outline-2 outline-offset-2"
                    data-set-theme="night" data-act-class="outline">
                    <div data-theme="night" class="bg-base-100 text-base-content w-full cursor-pointer font-sans">
                        <div class="grid grid-cols-5 grid-rows-3">
                            <div class="bg-base-200 col-start-1 row-span-2 row-start-1"></div>
                            <div class="bg-base-300 col-start-1 row-start-3"></div>
                            <div class="bg-base-100 col-span-4 col-start-2 row-span-3 row-start-1 flex flex-col gap-1 p-2">
                                <div class="font-bold">night</div>
                                <div class="flex flex-wrap gap-1">
                                    <div
                                        class="bg-primary flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-primary-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-secondary flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-secondary-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-accent flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-accent-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-neutral flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-neutral-content text-sm font-bold">A</div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </label>
    </label>
</template>
