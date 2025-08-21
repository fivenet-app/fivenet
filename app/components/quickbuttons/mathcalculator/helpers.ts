import Big from 'big.js';

export function evaluateTokens(tokens: string[]): Big {
    const operators: Record<string, (a: Big, b: Big) => Big> = {
        '+': (a: Big, b: Big) => a.plus(b),
        '-': (a: Big, b: Big) => a.minus(b),
        '*': (a: Big, b: Big) => a.times(b),
        '/': (a: Big, b: Big) => a.div(b),
    };

    const precedence: Record<string, number> = {
        '+': 1,
        '-': 1,
        '*': 2,
        '/': 2,
    };

    const outputQueue: (Big | string)[] = [];
    const operatorStack: string[] = [];

    tokens.forEach((token) => {
        if (!isNaN(Number(token))) {
            outputQueue.push(new Big(token));
        } else if (['+', '-', '*', '/'].includes(token)) {
            while (operatorStack.length) {
                const topOperator = operatorStack[operatorStack.length - 1];
                const topPrecedence = topOperator ? precedence[topOperator] : undefined;
                const currentPrecedence = precedence[token];

                if (topPrecedence !== undefined && currentPrecedence !== undefined && topPrecedence >= currentPrecedence) {
                    outputQueue.push(operatorStack.pop()!);
                } else {
                    break;
                }
            }
            operatorStack.push(token);
        } else {
            throw new Error('Invalid token in expression');
        }
    });

    while (operatorStack.length) {
        outputQueue.push(operatorStack.pop()!);
    }

    const evaluationStack: Big[] = [];

    outputQueue.forEach((item) => {
        if (item instanceof Big) {
            evaluationStack.push(item);
        } else if (typeof item === 'string' && operators[item]) {
            const b = evaluationStack.pop();
            const a = evaluationStack.pop();
            if (!a || !b) {
                throw new Error('Invalid expression');
            }
            evaluationStack.push(operators[item](a, b));
        }
    });

    if (evaluationStack.length !== 1) {
        throw new Error('Invalid expression');
    }

    const result = evaluationStack[0];
    if (!result) {
        throw new Error('Invalid expression: Result is undefined');
    }

    return result;
}
