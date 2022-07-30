// React imports
import React, { ChangeEvent } from 'react';

// Library imports
import { CdsInput } from '@cds/react/input';
import { useFormContext } from 'react-hook-form';
import { CdsControlMessage } from '@cds/react/forms';

interface Props {
    className?: string;
    defaultValue?: string;
    label: string;
    name: string;
    placeholder?: string;
    handleInputChange: (field: any, value: string) => void;
    maskText?: boolean;
}
function TextInputWithError(props: Props) {
    const { className, defaultValue, label, name, handleInputChange, placeholder, maskText } = props;
    const {
        register,
        formState: { errors },
    } = useFormContext();
    return (
        <CdsInput layout="vertical" control-width="shrink">
            <label>{label}</label>
            <input
                {...register(name, {
                    onChange: (e: ChangeEvent<HTMLInputElement>) => {
                        handleInputChange(name, e.target.value);
                    },
                })}
                name={name}
                placeholder={placeholder || label}
                type={maskText ? 'password' : 'text'}
                className={className}
                defaultValue={defaultValue}
            ></input>
            {errors[name] && <CdsControlMessage status="error">{errors[name]?.message}</CdsControlMessage>}
        </CdsInput>
    );
}

export default TextInputWithError;
