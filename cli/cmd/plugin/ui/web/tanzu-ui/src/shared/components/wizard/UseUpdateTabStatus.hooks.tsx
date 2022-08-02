import { useEffect, useRef } from 'react';
import { FieldErrors } from 'react-hook-form';

type UpdateTabStatus = (currentStep: number | undefined, validForm: boolean) => void;

function UseUpdateTabStatus(errors: FieldErrors, currentStep: number | undefined, updateTabStatus: UpdateTabStatus) {
    const hasError = Object.keys(errors).length > 0;
    const prevHasErrorRef = useRef<boolean>(false);

    useEffect(() => {
        if (updateTabStatus && prevHasErrorRef.current !== hasError) {
            updateTabStatus(currentStep, !hasError);
            prevHasErrorRef.current = hasError;
        }
    }, [currentStep, hasError, updateTabStatus]);
}

export default UseUpdateTabStatus;
