// jest-dom adds custom jest matchers for asserting on DOM nodes.
// allows you to do things like:
// expect(element).toHaveTextContent(/react/i)
// learn more: https://github.com/testing-library/jest-dom
import '@testing-library/jest-dom';

if (typeof window !== 'undefined') {
    // jsdom doesn't implment IntersectionObserver
    // https://github.com/jsdom/jsdom/issues/2032
    const mockObserverAPI: any = () => {
        return function () {
            return {
                observe: jest.fn(),
                unobserve: jest.fn(),
                disconnect: jest.fn(),
            };
        };
    };

    window.ResizeObserver = mockObserverAPI();
    window.IntersectionObserver = mockObserverAPI();

    // jsdom implements getBoundingClientRect, but not DOMRect
    if (!window.DOMRect) {
        const domRectMock = () => jest.fn().mockReturnValue({});
        window.DOMRect = domRectMock() as unknown as typeof window.DOMRect;
    }
}
