<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';
import { MinusSmallIcon, PlusSmallIcon } from '@heroicons/vue/24/outline';
import { CloudArrowUpIcon, LockClosedIcon, ServerIcon } from '@heroicons/vue/20/solid';
import Footer from '../components/partials/Footer.vue';

useHead({
    title: 'About',
});
definePageMeta({
    title: 'About',
    requiresAuth: false,
});

const faqs = [
    {
        question: "What's the best thing about Switzerland?",
        answer:
            "I don't know, but the flag is a big plus. Lorem ipsum dolor sit amet consectetur adipisicing elit. Quas cupiditate laboriosam fugiat.",
    },
]
</script>

<template>
    <div class="bg-gray-900 px-6 py-24 sm:py-32 lg:px-8">
        <div class="mx-auto max-w-2xl text-center">
            <h2 class="text-4xl font-bold tracking-tight text-white sm:text-6xl">About</h2>
            <p class="mt-6 text-lg leading-8 text-gray-300">
                About FiveNet and it's creators.
            </p>
        </div>
    </div>

    <div class="relative isolate overflow-hidden px-6 py-24 sm:py-32 lg:overflow-visible lg:px-0">
        <div class="absolute inset-0 -z-10 overflow-hidden"></div>
        <div
            class="mx-auto grid max-w-2xl grid-cols-1 gap-x-8 gap-y-16 lg:mx-0 lg:max-w-none lg:grid-cols-2 lg:items-start lg:gap-y-10">
            <div
                class="lg:col-span-2 lg:col-start-1 lg:row-start-1 lg:mx-auto lg:grid lg:w-full lg:max-w-7xl lg:grid-cols-2 lg:gap-x-8 lg:px-8">
                <div class="lg:pr-4">
                    <div class="lg:max-w-lg">
                        <p class="text-base font-semibold leading-7 text-primary-400">Deploy faster</p>
                        <h1 class="mt-2 text-3xl font-bold tracking-tight text-gray-400 sm:text-4xl">A better workflow
                        </h1>
                        <p class="mt-6 text-xl leading-8 text-gray-100">Aliquet nec orci mattis amet quisque ullamcorper
                            neque, nibh sem. At arcu, sit dui mi, nibh dui, diam eget aliquam. Quisque id at vitae
                            feugiat
                            egestas.</p>
                    </div>
                </div>
            </div>
            <div
                class="-ml-12 -mt-12 p-12 lg:sticky lg:top-4 lg:col-start-2 lg:row-span-2 lg:row-start-1 lg:overflow-hidden">
                <img class="w-[48rem] max-w-none rounded-xl bg-gray-900 shadow-xl ring-1 ring-gray-400/10 sm:w-[57rem]"
                    src="https://tailwindui.com/img/component-images/dark-project-app-screenshot.png" alt="" />
            </div>
            <div
                class="lg:col-span-2 lg:col-start-1 lg:row-start-2 lg:mx-auto lg:grid lg:w-full lg:max-w-7xl lg:grid-cols-2 lg:gap-x-8 lg:px-8">
                <div class="lg:pr-4">
                    <div class="max-w-xl text-base leading-7 text-gray-100 lg:max-w-lg">
                        <p>Faucibus commodo massa rhoncus, volutpat. Dignissim sed eget risus enim. Mattis mauris semper
                            sed
                            amet vitae sed turpis id. Id dolor praesent donec est. Odio penatibus risus viverra tellus
                            varius sit neque erat velit. Faucibus commodo massa rhoncus, volutpat. Dignissim sed eget
                            risus
                            enim. Mattis mauris semper sed amet vitae sed turpis id.</p>
                        <ul role="list" class="mt-8 space-y-8 text-gray-400">
                            <li class="flex gap-x-3">
                                <CloudArrowUpIcon class="mt-1 h-5 w-5 flex-none text-primary-300" aria-hidden="true" />
                                <span><strong class="font-semibold text-gray-200">Push to deploy.</strong> Lorem ipsum,
                                    dolor sit amet consectetur adipisicing elit. Maiores impedit perferendis suscipit
                                    eaque,
                                    iste dolor cupiditate blanditiis ratione.</span>
                            </li>
                            <li class="flex gap-x-3">
                                <LockClosedIcon class="mt-1 h-5 w-5 flex-none text-primary-300" aria-hidden="true" />
                                <span><strong class="font-semibold text-gray-200">SSL certificates.</strong> Anim aute
                                    id
                                    magna aliqua ad ad non deserunt sunt. Qui irure qui lorem cupidatat commodo.</span>
                            </li>
                            <li class="flex gap-x-3">
                                <ServerIcon class="mt-1 h-5 w-5 flex-none text-primary-300" aria-hidden="true" />
                                <span><strong class="font-semibold text-gray-200">Database backups.</strong> Ac
                                    tincidunt
                                    sapien vehicula erat auctor pellentesque rhoncus. Et magna sit morbi
                                    lobortis.</span>
                            </li>
                        </ul>
                        <p class="mt-8">Et vitae blandit facilisi magna lacus commodo. Vitae sapien duis odio id et. Id
                            blandit molestie auctor fermentum dignissim. Lacus diam tincidunt ac cursus in vel. Mauris
                            varius vulputate et ultrices hac adipiscing egestas. Iaculis convallis ac tempor et ut. Ac
                            lorem
                            vel integer orci.</p>
                        <h2 class="mt-16 text-2xl font-bold tracking-tight text-gray-200">No server? No problem.</h2>
                        <p class="mt-6">Id orci tellus laoreet id ac. Dolor, aenean leo, ac etiam consequat in.
                            Convallis
                            arcu ipsum urna nibh. Pharetra, euismod vitae interdum mauris enim, consequat vulputate
                            nibh.
                            Maecenas pellentesque id sed tellus mauris, ultrices mauris. Tincidunt enim cursus ridiculus
                            mi.
                            Pellentesque nam sed nullam sed diam turpis ipsum eu a sed convallis diam.</p>
                    </div>
                </div>
            </div>
        </div>
    </div>
    <div class="bg-gray-900">
        <div class="mx-auto max-w-7xl px-6 py-24 sm:py-32 lg:px-8 lg:py-40">
            <div class="mx-auto max-w-4xl divide-y divide-white/10">
                <h2 class="text-2xl font-bold leading-10 tracking-tight text-white">Frequently asked questions</h2>
                <dl class="mt-10 space-y-6 divide-y divide-white/10">
                    <Disclosure as="div" v-for="faq in faqs" :key="faq.question" class="pt-6" v-slot="{ open }">
                        <dt>
                            <DisclosureButton class="flex w-full items-start justify-between text-left text-white">
                                <span class="text-base font-semibold leading-7">{{ faq.question }}</span>
                                <span class="ml-6 flex h-7 items-center">
                                    <PlusSmallIcon v-if="!open" class="h-6 w-6" aria-hidden="true" />
                                    <MinusSmallIcon v-else class="h-6 w-6" aria-hidden="true" />
                                </span>
                            </DisclosureButton>
                        </dt>
                        <DisclosurePanel as="dd" class="mt-2 pr-12">
                            <p class="text-base leading-7 text-gray-300">{{ faq.answer }}</p>
                        </DisclosurePanel>
                    </Disclosure>
                </dl>
            </div>
        </div>
    </div>
    <div class="relative bg-gray-900">
        <div class="relative h-80 overflow-hidden bg-primary-600 md:absolute md:left-0 md:h-full md:w-1/3 lg:w-1/2">
            <img class="h-full w-full object-cover"
                src="https://images.unsplash.com/photo-1525130413817-d45c1d127c42?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=1920&q=60&blend=6366F1&sat=-100&blend-mode=multiply"
                alt="" />
            <svg viewBox="0 0 926 676" aria-hidden="true"
                class="absolute -bottom-24 left-24 w-[57.875rem] transform-gpu blur-[118px]">
                <path fill="url(#60c3c621-93e0-4a09-a0e6-4c228a0116d8)" fill-opacity=".4"
                    d="m254.325 516.708-90.89 158.331L0 436.427l254.325 80.281 163.691-285.15c1.048 131.759 36.144 345.144 168.149 144.613C751.171 125.508 707.17-93.823 826.603 41.15c95.546 107.978 104.766 294.048 97.432 373.585L685.481 297.694l16.974 360.474-448.13-141.46Z" />
                <defs>
                    <linearGradient id="60c3c621-93e0-4a09-a0e6-4c228a0116d8" x1="926.392" x2="-109.635" y1=".176"
                        y2="321.024" gradientUnits="userSpaceOnUse">
                        <stop stop-color="#776FFF" />
                        <stop offset="1" stop-color="#FF4694" />
                    </linearGradient>
                </defs>
            </svg>
        </div>
        <div class="relative mx-auto max-w-7xl py-24 sm:py-32 lg:px-8 lg:py-40">
            <div class="pl-6 pr-6 md:ml-auto md:w-2/3 md:pl-16 lg:w-1/2 lg:pl-24 lg:pr-0 xl:pl-32">
                <p class="mt-2 text-3xl font-bold tracking-tight text-white sm:text-4xl">Got any issues?</p>
                <p class="mt-6 text-base leading-7 text-gray-300">
                    Be sure to join the Discord server should you encounter any issues or have questions regarding FiveNet.
                </p>
                <div class="mt-8">
                    <a href="https://discord.gg/sWvkHuVQA5"
                        class="inline-flex rounded-md bg-white/10 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-white/20 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white">
                        Join the Discord
                    </a>
                </div>
            </div>
        </div>
    </div>
    <div class="mt-24"></div>
    <Footer />
</template>
