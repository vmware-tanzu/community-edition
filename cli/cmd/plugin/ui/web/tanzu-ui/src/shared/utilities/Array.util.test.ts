import '@testing-library/jest-dom';
import { first, last, middle } from './Array.util';

describe('Array utilities', () => {
    test('first(empty)', () => {
        const result = first<string>([]);
        expect(result).toBeUndefined();
    });
    test('first(data)', () => {
        const result = first<string>(['a', 'b', 'c']);
        expect(result).toEqual('a');
    });
    test('last(empty)', () => {
        const result = last<string>([]);
        expect(result).toBeUndefined();
    });
    test('last(data)', () => {
        const result = last<string>(['a', 'b', 'c']);
        expect(result).toEqual('c');
    });
    test('middle(empty)', () => {
        const result = middle<string>([]);
        expect(result).toBeUndefined();
    });
    test('middle(even array)', () => {
        const result = middle<string>(['a', 'b', 'c', 'd']);
        // NOTE: with an even array either of the middle two elements satisfies the "middle" condition
        expect(result).toMatch(/[b|c]/);
    });
    test('middle(odd array)', () => {
        const result = middle<string>(['a', 'b', 'c']);
        expect(result).toEqual('b');
    });
});
