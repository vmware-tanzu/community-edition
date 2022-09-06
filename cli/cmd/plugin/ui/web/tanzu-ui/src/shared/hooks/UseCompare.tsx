import { useEffect, useRef } from 'react';

function useCompare(value: any): boolean {
    const ref = useRef();
    useEffect(() => {
        ref.current = value;
    }, [value]);
    return ref.current !== value;
}
export default useCompare;
