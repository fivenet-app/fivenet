<script lang="ts" setup>
import { vMaska } from 'maska/vue';

const height = ref<string>('165');
const mass = ref<string>('83');
const bmi = ref(0);

watchDebounced(height, () => bmiCalculate(), {
    debounce: 200,
});

watchDebounced(mass, () => bmiCalculate(), {
    debounce: 200,
});

function bmiCalculate(): void {
    bmi.value = parseInt(mass.value) / (parseInt(height.value) / 100) ** 2;
}

onBeforeMount(() => bmiCalculate());
</script>

<template>
    <div>
        <h3 class="text-xl font-semibold leading-6">
            {{ $t('components.bodycheckup.bmi_calculator') }}
        </h3>
        <div>
            <div>
                <label for="height" class="block text-sm font-medium leading-6">
                    {{ $t('components.bodycheckup.height') }}
                </label>
                <div>
                    <div class="relative rounded-md shadow-sm">
                        <UInput
                            v-model="height"
                            v-maska
                            data-maska="###"
                            name="height"
                            type="text"
                            :placeholder="$t('components.bodycheckup.height')"
                            :label="$t('components.bodycheckup.height')"
                        />
                        <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-3">
                            <span class="sm:text-sm">cm</span>
                        </div>
                    </div>
                </div>
            </div>

            <div>
                <label for="mass" class="block text-sm font-medium leading-6">
                    {{ $t('components.bodycheckup.mass') }}
                </label>
                <div>
                    <div class="relative rounded-md shadow-sm">
                        <UInput
                            v-model="mass"
                            v-maska
                            data-maska="###,##"
                            name="mass"
                            type="text"
                            :placeholder="$t('components.bodycheckup.mass')"
                            :label="$t('components.bodycheckup.mass')"
                        />
                        <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-3">
                            <span class="sm:text-sm">kg</span>
                        </div>
                    </div>
                </div>
            </div>
            <p class="block text-sm font-medium leading-6">
                BMI: <span class="font-semibold">{{ bmi.toFixed(1) }}</span>
            </p>
        </div>
    </div>
</template>
