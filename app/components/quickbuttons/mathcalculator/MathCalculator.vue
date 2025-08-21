<script lang="ts" setup>
import type Big from 'big.js';
import { evaluateTokens } from './helpers';

const powerOn = ref(true);

const histroryOperation = useState<
    {
        operation: string;
        result: string;
    }[]
>('quickButton:mathCalculator:historyOps', () => []);

const { n } = useI18n();

const currentCalculate = ref('0');
const lastOperation = ref('0');
const lastStatus = ref('');

const numberInput = useTemplateRef('numberInput');
const resultElement = useTemplateRef('resultElement');

function getPercent() {
    if (!powerOn.value) {
        return;
    }

    const percent = parseFloat(currentCalculate.value.replace(',', '.')) / 100;
    currentCalculate.value = percent.toString().replace('.', ',');
}

function allClear() {
    if (!powerOn.value) {
        return;
    }

    currentCalculate.value = '0';
    lastOperation.value = '0';
    lastStatus.value = '';
}

function doDelete() {
    if (!powerOn.value) {
        return;
    }

    if (currentCalculate.value !== '0') {
        currentCalculate.value = currentCalculate.value.slice(0, -1);
    }
    if (currentCalculate.value === '') {
        currentCalculate.value = '0';
    }
}

function inputValue(val: string) {
    if (!powerOn.value) {
        return;
    }

    // Clear the error status when a new character is inputted
    lastStatus.value = '';

    if (val === '0' || val === '00') {
        if (currentCalculate.value === '0') {
            currentCalculate.value = '0';
        } else {
            currentCalculate.value += maskedIn(val);
        }
    } else if (val === ',') {
        if (!currentCalculate.value.includes(',')) {
            if (lastIsOperand()) {
                return currentCalculate.value;
            } else {
                currentCalculate.value += maskedIn(val);
            }
        }
    } else if (val === '/' || val === '*' || val === '-' || val === '+') {
        if (lastIsOperand()) {
            const vars = currentCalculate.value.slice(0, -1) + val;
            currentCalculate.value = maskedIn(vars);
        } else {
            currentCalculate.value += maskedIn(val);
        }
    } else {
        if (currentCalculate.value == '0' || currentCalculate.value === 'NaN') {
            currentCalculate.value = maskedIn(val);
        } else {
            currentCalculate.value += maskedIn(val);
        }
    }

    nextTick(() => {
        scrollRight();
    });
}

function scrollRight() {
    if (numberInput.value?.scrollTop === undefined) {
        return;
    }
    numberInput.value.scrollLeft = numberInput.value?.scrollWidth;
}

function scrollDown() {
    if (resultElement.value?.scrollTop === undefined) {
        return;
    }
    resultElement.value.scrollTop = resultElement.value.scrollHeight;
}

function calculate() {
    if (!powerOn.value) {
        return;
    }

    if (lastIsOperand()) {
        // Trigger an error for incomplete expressions
        lastStatus.value = 'ERROR';
        console.error('Calculation error: Incomplete expression');
        return;
    }

    try {
        counted();
        nextTick(() => {
            scrollDown();
        });
    } catch (error) {
        console.error('Calculation error:', error);
        lastStatus.value = 'ERROR';
    }
}

function counted() {
    try {
        const rep = currentCalculate.value
            .replace('x', '*')
            .replace(',', '.')
            .replace(/([*+/-])0(?!\.)/, '$1');
        const lastCalculate = currentCalculate.value;

        // Parse and calculate using Big.js
        const tokens = tokenize(rep);
        const result = evaluateTokens(tokens);

        addHistroy(lastCalculate, formatResult(result));
        currentCalculate.value = formatResult(result);
    } catch (error) {
        console.error('Calculation error:', error);
        lastStatus.value = 'ERROR';
    }
}

function tokenize(expression: string): string[] {
    const regex = /([+\-*/])|([0-9]+(?:\.[0-9]*)?)/g;
    const tokens: string[] = [];
    let match;
    while ((match = regex.exec(expression)) !== null) {
        tokens.push(match[0]);
    }
    return tokens;
}

function formatResult(result: Big): string {
    return n(result.toNumber());
}

function maskedIn(val: string) {
    return val.replace(/\*/g, 'x');
}

function lastIsOperand() {
    const operand = currentCalculate.value.slice(-1);
    if (operand === '+' || operand === '-' || operand === 'x' || operand === '/') {
        return true;
    } else {
        return false;
    }
}

function checkContainsOperand(str: string) {
    const operand = ['+', '-', 'x', '/'];
    return operand.some((v) => str.includes(v));
}

function addHistroy(operation: string, result: string) {
    if (checkContainsOperand(currentCalculate.value)) {
        histroryOperation.value.push({
            operation: operation,
            result: result,
        });
    }
}

function deleteHistory() {
    histroryOperation.value = [];
}

const digitKeys = ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', ','];
const operatorKeys = ['/', '*', '+', '-'];
const resultKeys = ['=', 'Enter'];
const eraseKeys = ['Backspace', 'Delete'];
const clearKeys = ['Escape'];

onKeyStroke([...digitKeys, ...operatorKeys, ...resultKeys, ...clearKeys, ...eraseKeys], (e) => {
    const key = e.key === ',' ? '.' : e.key;

    if (digitKeys.includes(key)) {
        inputValue(key);
    } else if (operatorKeys.includes(key)) {
        inputValue(key);
    } else if (resultKeys.includes(key)) {
        calculate();
    } else if (eraseKeys.includes(key)) {
        doDelete();
    } else if (clearKeys.includes(key)) {
        allClear();
    }
});
</script>

<template>
    <div class="flex flex-1 flex-col gap-4 overflow-hidden">
        <div
            :class="[
                'relative flex h-48 max-h-48 flex-col items-end justify-end rounded-md text-right text-white transition-colors duration-300',
                { 'bg-red-800': lastStatus === 'ERROR', 'bg-gray-800': lastStatus !== 'ERROR' },
            ]"
        >
            <template v-if="powerOn">
                <UTooltip :text="$t('common.delete')">
                    <UButton
                        class="absolute left-0 top-0 m-2 p-0.5"
                        :padded="false"
                        size="xs"
                        icon="i-mdi-delete"
                        color="error"
                        @click="deleteHistory"
                    />
                </UTooltip>

                <div ref="resultElement" class="space-y-3 overflow-y-auto p-1.5">
                    <div v-for="(history, idx) in histroryOperation" :key="idx" class="-space-y-1">
                        <p class="text-lg font-light">{{ history.operation }}</p>
                        <p class="text-2xl font-medium">&#x3D; {{ $n(parseInt(history.result)) }}</p>
                    </div>
                </div>

                <div class="flex w-full flex-col items-end p-1.5">
                    <p class="text-2xl">{{ lastOperation != '0' ? lastOperation : '' }}</p>
                    <div ref="numberInput" class="scroll z-0 w-full overflow-x-auto overflow-y-hidden">
                        <p class="text-5xl font-semibold">{{ lastOperation != '0' ? '&#x3D;' : '' }}{{ currentCalculate }}</p>
                    </div>
                </div>
            </template>
            <div v-else class="flex w-full flex-col items-end">
                <p class="text-5xl font-semibold">123'456'789'1234</p>
            </div>
        </div>

        <div class="grid grid-cols-1 gap-2">
            <div class="grid grid-cols-4 gap-2">
                <UButton block :color="powerOn ? 'success' : 'error'" icon="i-mdi-power-standby" @click="powerOn = !powerOn" />
                <UButton block color="white" @click="allClear()">
                    <span class="font-semibold">AC</span>
                </UButton>
                <UButton block color="white" icon="i-mdi-percent" @click="getPercent()" />
                <UButton block color="white" icon="i-mdi-division" @click="inputValue('/')" />
            </div>

            <div class="grid grid-cols-4 gap-2">
                <UButton block color="black" icon="i-mdi-numeric-7" @click="inputValue('7')" />
                <UButton block color="black" icon="i-mdi-numeric-8" @click="inputValue('8')" />
                <UButton block color="black" icon="i-mdi-numeric-9" @click="inputValue('9')" />
                <UButton block color="white" icon="i-mdi-multiplication" @click="inputValue('*')" />
            </div>

            <div class="grid grid-cols-4 gap-2">
                <UButton block color="black" icon="i-mdi-numeric-4" @click="inputValue('4')" />
                <UButton block color="black" icon="i-mdi-numeric-5" @click="inputValue('5')" />
                <UButton block color="black" icon="i-mdi-numeric-6" @click="inputValue('6')" />
                <UButton block color="white" icon="i-mdi-minus" @click="inputValue('-')" />
            </div>

            <div class="grid grid-cols-4 gap-2">
                <UButton block color="black" icon="i-mdi-numeric-1" @click="inputValue('1')" />
                <UButton block color="black" icon="i-mdi-numeric-2" @click="inputValue('2')" />
                <UButton block color="black" icon="i-mdi-numeric-3" @click="inputValue('3')" />
                <UButton block color="white" icon="i-mdi-plus" @click="inputValue('+')" />
            </div>

            <div class="grid grid-cols-4 gap-2">
                <UButton block color="black" icon="i-mdi-numeric-0" @click="inputValue('0')" />
                <UButton block color="black" icon="i-mdi-comma" @click="inputValue(',')" />
                <UButton block color="white" @click="doDelete()">
                    <span class="font-semibold">DEL</span>
                </UButton>
                <UButton block color="primary" icon="i-mdi-equal" @click="calculate()" />
            </div>
        </div>
    </div>
</template>

<style scoped>
.crt::before {
    content: ' ';
    display: block;
    position: absolute;
    top: 0;
    left: 0;
    bottom: 0;
    right: 0;
    background:
        linear-gradient(rgba(18, 16, 16, 0) 50%, rgba(0, 0, 0, 0.25) 50%),
        linear-gradient(90deg, rgba(255, 0, 0, 0.06), rgba(0, 255, 0, 0.02), rgba(0, 0, 255, 0.06));
    z-index: 2;
    background-size:
        100% 2px,
        3px 100%;
    pointer-events: none;
}
</style>
