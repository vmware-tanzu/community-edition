import { addErrorInfo, removeErrorInfo } from './Error.util';

describe('Array utilities', () => {
    test('removeErrorInfo(foo) removes foo but leaves bar', () => {
        const obj = { foo: 'some value', bar: 'some other value' };
        const result = removeErrorInfo(obj, 'foo');
        expect(result['foo']).toBeUndefined();
        expect(result['bar']).toEqual('some other value');
    });
    test('addErrorInfo(foo) add foo and leaves bar', () => {
        const obj = { bar: 'some other value' };
        const value = { x: 'y' };
        const result = addErrorInfo(obj, value, 'foo');
        expect(result['foo']).toEqual(value);
        expect(result['bar']).toEqual('some other value');
    });
});
