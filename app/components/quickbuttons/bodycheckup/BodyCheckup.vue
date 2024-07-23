<script lang="ts" setup>
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import BMICalculator from '~/components/quickbuttons/bodycheckup/BMICalculator.vue';
import { useNotificatorStore } from '~/store/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const notifications = useNotificatorStore();

const { t } = useI18n();

type Pin = {
    x: number;
    y: number;
    selected?: boolean;
    description?: string;
};

const pins = useState<Pin[]>('quickButton:bodyCheckup:summary', () => []);

const svgRef = ref<SVGSVGElement | null>(null);

function addPin(event: MouseEvent): void {
    if (svgRef.value === null) {
        return;
    }

    const pt = svgRef.value.createSVGPoint();
    pt.x = event.clientX;
    pt.y = event.clientY;
    const cursorPt = pt.matrixTransform(svgRef.value.getScreenCTM()!.inverse());

    pins.value.push({
        x: cursorPt.x,
        y: cursorPt.y,
    });
}

const selectedPin = ref<Pin | undefined>();
function selectPin(p: Pin): void {
    if (selectedPin.value !== undefined) {
        selectedPin.value.selected = false;
    }

    p.selected = true;
    selectedPin.value = p;
}

function removePin(idx: number): void {
    pins.value.splice(idx, 1);
}

async function copyToClipboard(): Promise<void> {
    let text = t('components.bodycheckup.title') + `\n`;

    pins.value.forEach((p, idx) => {
        text += (idx + 1).toString() + '. ' + p.description + '\n';
    });

    notifications.add({
        title: { key: 'notifications.bodycheckup.title', parameters: {} },
        description: { key: 'notifications.bodycheckup.content', parameters: {} },
        type: NotificationType.INFO,
    });

    return copyToClipboardWrapper(text);
}

function reset(): void {
    pins.value = [];
}
</script>

<template>
    <div class="grid w-full grid-cols-1 gap-4 sm:grid-cols-2">
        <div>
            <svg
                ref="svgRef"
                version="1.1"
                xmlns="http://www.w3.org/2000/svg"
                xmlns:xlink="http://www.w3.org/1999/xlink"
                x="0px"
                y="0px"
                width="420px"
                height="780px"
                viewBox="0 0 420 780"
                enable-background="new 0 0 420 780"
                xml:space="preserve"
            >
                <g>
                    <path
                        fill="#FFDEC7"
                        d="M375.86,438.039c-4.383-1.051-14.305-5.084-17.721-7.705c-3.414-2.625-9.121-5.25-10.98-10.928
		c-8.719-26.635-6.174-42.584-12.33-80.583c-3.279-20.259-10.594-38.978-12.291-49.546c-5.002-31.152,0.314-68.261-9.955-100.304
		c-8.758-27.333-23.908-41.438-41.4-44.698c-20.547-3.037-37.541-14.717-39.795-21.196c0.059-1.203,0.465-5.679,0.793-9.207
		c4.123-3.525,6.537-7.708,7.357-11.811c1.736-8.662,4.16-12.907,4.16-12.907s3.189,1.093,4.959-1.936
		c0.781-1.342,0.545-4.525,1.527-9.028c1.264-5.775,4.736-10.806-0.355-12.834c-2.934-1.165-4.52,0.494-4.52,0.494
		c0.332-1.69,1.822-4.07,2.139-5.757c4.262-22.743-0.867-47.178-22.947-48.343c-8.394-0.442-10.479-0.458-24.623,1.443
		c-23.828,3.204-25.584,26.232-23.307,40.744c0.81,5.156,2.608,9.326,2.658,12.186c0.001,0.001,0.002,0.002,0.004,0.002
		c-1.595-1.072-3.429-1.382-4.965-0.77c-5.09,2.029-1.616,7.059-0.354,12.834c0.984,4.503,0.745,7.687,1.527,9.028
		c1.769,3.029,4.958,1.936,4.958,1.936s2.399,4.245,4.133,12.907c0.822,4.107,3.251,8.295,7.385,11.823
		c0.337,3.668,0.754,8.354,0.795,9.194c-2.256,6.479-19.213,18.159-39.761,21.196c-17.492,3.261-29.013,12.168-36.785,32.724
		c-12.462,32.957-8.649,82.364-14.569,112.278c-2.078,10.5-9.012,29.288-12.291,49.546c-6.157,38-3.612,53.949-12.33,80.583
		c-1.859,5.678-7.566,8.303-10.982,10.928c-3.414,2.617-13.337,6.654-17.719,7.705c-4.387,1.053-12.249,2.311-10.652,5.324
		c1.592,3.008,10.973,3.93,15.565,3.357c4.594-0.572,8.902-2.443,12.412-0.27c3.508,2.172,2.364,6.092,1.015,9.525
		c-1.356,3.438-2.306,5.842-3.521,8.936c-1.221,3.094-4.678,10.863-6.029,14.295c-1.354,3.439-6.797,11.219-1.717,12.424
		c5.086,1.211,7.301-7.443,9.34-10.611c2.04-3.164,6.176-13.77,7.956-13.139c1.779,0.633-0.239,6.652-1.112,9.883
		c-0.876,3.225-3.418,12.158-4.699,16.418c-1.282,4.262-3.459,7.77,0.665,9.393c4.126,1.625,5.53-4.969,5.53-4.969
		s2.207-7.078,3.833-11.203c1.955-4.959,3.715-17.967,5.722-17.484c1.845,0.445,1.681,6.135,1.803,9.328
		c0.125,3.189-0.082,10.957-0.033,15.205c0.046,4.25-0.938,8.068,3.287,8.375c4.221,0.301,3.621-6.111,3.621-6.111
		s0.428-5.977,0.736-10.201c1.229-17.053,1.15-19.08,2.185-19.158c1.035-0.078,2.409,4.799,2.895,7.924
		c0.472,3.031,1.129,10.461,1.644,14.508c0.512,4.051-0.004,7.803,4.057,7.629c4.061-0.178,2.782-6.229,2.782-6.229
		s-0.248-5.75-0.419-9.816c-0.161-3.729-1.881-10.793-2.289-14.807c-0.382-3.178,2.616-10.561,3.377-16.92
		c0.903-10.361-3.785-20.211-3.23-26.406c1.803-20.283,6.973-28.092,14.745-49.844c17.097-47.856,13.28-65.657,15.886-80.78
		c2.605-15.122,10.261-24.323,14.427-40.698c5.112,20.145,7.377,54.859,6.403,76.693c-0.985,22.071-5.784,33.415-8.874,67.759
		c-5.461,29.082,9.767,111.285,9.767,143.623c0,20.951-6.761,36.402-2.704,68.162c12.218,95.697,12.934,104.977,4.644,121.693
		c-3.639,7.344-14.124,12.658-18.891,14.695s-11.323,6.408-17.261,8.039c-5.94,1.631-11.073,5.115-13.702,5.324
		c-2.634,0.213-5.73,0.064-7.277,2.088c-1.55,2.027-1.46,4.814,1.659,7.197c3.114,2.383,8.34,1.688,8.34,1.688
		c-1,0.877,6.705,4.803,11.274,2.691c0,0,2.694,1.32,7.279,2.111c4.58,0.787,15.077-1.311,20.319-3.973
		c5.242-2.654,15.64-6.559,21.285-6.191c5.647,0.365,14.962,3.043,20.338,2.461c5.383-0.574,8.903-3.562,10.314-7.67
		c1.411-4.105-1.256-12.566-3.064-16.664c-1.804-4.1-4.258-14.123-5.387-23.629c-1.456-12.289-2.07-54.711,2.704-72.016
		c5.406-19.596,8.317-48.104,6.963-73.783c-1.35-25.682-1.963-22.494-1.227-34.037c0.442-6.928,6.283-64.115,11.111-101.697
		c4.826,37.582,10.668,94.77,11.109,101.697c0.736,11.543,0.123,8.355-1.227,34.037c-1.354,25.68,1.557,54.188,6.963,73.783
		c4.775,17.305,4.16,59.727,2.705,72.016c-1.129,9.506-3.584,19.529-5.389,23.629c-1.807,4.098-4.475,12.559-3.062,16.664
		c1.41,4.107,4.932,7.096,10.314,7.67c5.375,0.582,14.691-2.096,20.338-2.461c5.646-0.367,16.043,3.537,21.285,6.191
		c5.242,2.662,15.74,4.76,20.32,3.973c4.584-0.791,7.279-2.111,7.279-2.111c4.57,2.111,12.273-1.814,11.273-2.691
		c0,0,5.227,0.695,8.34-1.688c3.119-2.383,3.209-5.17,1.66-7.197c-1.547-2.023-4.645-1.875-7.277-2.088
		c-2.629-0.209-7.762-3.693-13.703-5.324c-5.938-1.631-12.494-6.002-17.26-8.039c-4.768-2.037-15.252-7.352-18.891-14.695
		c-8.291-16.717-7.574-25.996,4.643-121.693c4.057-31.76-2.703-47.211-2.703-68.162c0-32.338,15.229-114.541,9.766-143.623
		c-3.09-34.344-7.889-45.688-8.873-67.759c-0.975-21.834,1.291-56.548,6.402-76.693c4.166,16.375,11.822,25.576,14.428,40.698
		c2.605,15.123-1.211,32.925,15.885,80.781c7.773,21.75,12.943,29.559,14.746,49.842c0.555,6.195-4.135,16.045-3.23,26.406
		c0.76,6.359,3.758,13.744,3.377,16.92c-0.408,4.014-2.129,11.078-2.291,14.809c-0.17,4.064-0.418,9.814-0.418,9.814
		s-1.279,6.051,2.783,6.229c4.059,0.174,3.543-3.578,4.057-7.629c0.514-4.047,1.172-11.477,1.643-14.508
		c0.486-3.125,1.859-8.002,2.895-7.924s0.955,2.105,2.186,19.158c0.309,4.225,0.734,10.201,0.734,10.201s-0.598,6.412,3.623,6.111
		c4.223-0.305,3.24-4.125,3.287-8.375c0.049-4.248-0.158-12.016-0.035-15.205c0.123-3.193-0.041-8.883,1.805-9.328
		c2.006-0.482,3.766,12.525,5.721,17.484c1.627,4.125,3.834,11.203,3.834,11.203s1.402,6.594,5.529,4.969
		c4.123-1.623,1.947-5.131,0.664-9.393c-1.279-4.26-3.822-13.193-4.697-16.418c-0.873-3.23-2.893-9.25-1.113-9.883
		c1.781-0.631,5.916,9.975,7.957,13.139c2.039,3.168,4.254,11.822,9.34,10.611c5.078-1.205-0.363-8.984-1.717-12.424
		c-1.352-3.432-4.809-11.201-6.031-14.295c-1.215-3.094-2.164-5.498-3.52-8.936c-1.35-3.434-2.494-7.354,1.014-9.525
		c3.51-2.174,7.818-0.303,12.412,0.27s10.975-0.35,12.566-3.357C385.108,440.35,380.246,439.092,375.86,438.039z"
                        @click="addPin"
                    />

                    <image
                        overflow="visible"
                        width="130"
                        height="416"
                        xlink:href="/images/components/quickbuttons/bodycheckup/human-organs.png"
                        transform="matrix(1 0 0 1 150.0002 12.8901)"
                        @click="addPin"
                    ></image>

                    <template v-for="(pin, idx) in pins" :key="idx">
                        <circle
                            :cx="pin.x"
                            :cy="pin.y"
                            r="12"
                            :class="pin.selected ? 'animate-pulse' : ''"
                            :data-x="pin.x - 50"
                            :data-y="pin.y - 50"
                            @click="selectPin(pin)"
                        />
                        <text
                            :x="pin.x"
                            :y="pin.y"
                            text-anchor="middle"
                            stroke="#ffffff"
                            stroke-width="1.5px"
                            dy=".3em"
                            @click="selectPin(pin)"
                        >
                            {{ idx + 1 }}
                        </text>
                    </template>
                </g>
            </svg>
        </div>
        <div class="flex flex-col">
            <div class="flex-1">
                <h3 class="text-xl font-semibold leading-6">
                    {{ $t('common.summary') }}
                </h3>

                <DataNoDataBlock
                    v-if="pins.length === 0"
                    :message="$t('components.bodycheckup.no_points')"
                    icon="i-mdi-vector-point-select"
                    class="mt-2"
                />
                <ol v-else>
                    <li v-for="(pin, idx) in pins" :key="idx" class="my-2 inline-flex w-full items-center gap-1">
                        <span class="w-4 text-base" :class="pin.selected ? 'underline' : ''"> {{ idx + 1 }}. </span>

                        <UInput
                            v-model="pin.description"
                            type="text"
                            block
                            class="flex-1"
                            @focusin="
                                focusTablet(true);
                                selectPin(pin);
                            "
                            @focusout="focusTablet(false)"
                        />

                        <UButton variant="link" icon="i-mdi-trash-can" class="ml-1" @click="removePin(idx)" />
                    </li>
                </ol>
            </div>

            <div class="my-2 flow-root">
                <UButtonGroup class="inline-flex w-full">
                    <UButton icon="i-mdi-content-copy" class="flex-1" @click="copyToClipboard()">
                        {{ $t('common.copy') }}
                    </UButton>
                    <UButton trailing-icon="i-mdi-clear-outline" color="red" @click="reset()">
                        {{ $t('common.reset') }}
                    </UButton>
                </UButtonGroup>
            </div>

            <BMICalculator />
        </div>
    </div>
</template>

<style scoped>
circle {
    width: 5vh;
    height: 5vh;
    position: absolute;
    z-index: 8;
    display: block;
    stroke-width: 5px;
    stroke: rgba(255, 255, 255, 0.9);
    fill: rgba(240, 0, 0, 0.8);
}
</style>
