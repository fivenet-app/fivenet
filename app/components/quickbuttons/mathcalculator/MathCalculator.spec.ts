import { describe, expect, it } from 'vitest';
import { evaluateTokens } from './helpers';

// Mocked functions from MathCalculator.vue
function tokenize(expression: string): string[] {
    const regex = /([+\-*/])|([0-9]+(?:\.[0-9]*)?)/g;
    const tokens: string[] = [];
    let match;
    while ((match = regex.exec(expression)) !== null) {
        tokens.push(match[0]);
    }
    return tokens;
}

describe('MathCalculator', () => {
    it('should tokenize a valid expression', () => {
        const expression = '12+34-5*6/2';
        const tokens = tokenize(expression);
        expect(tokens).toEqual(['12', '+', '34', '-', '5', '*', '6', '/', '2']);
    });

    it('should evaluate a simple addition', () => {
        const tokens = tokenize('12+34');
        const result = evaluateTokens(tokens);
        expect(result.toString()).toBe('46');
    });

    it('should evaluate a mixed operation', () => {
        const tokens = tokenize('12+34-5*6/2');
        const result = evaluateTokens(tokens);
        expect(result.toString()).toBe('31');
    });

    it('should throw an error for invalid expressions', () => {
        const tokens = tokenize('12+');
        expect(() => evaluateTokens(tokens)).toThrow('Invalid expression');
    });

    it('should handle division by zero gracefully', () => {
        const tokens = tokenize('12/0');
        expect(() => evaluateTokens(tokens)).toThrow();
    });
});
