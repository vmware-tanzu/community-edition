import React, { ChangeEvent } from 'react';

export interface RadioButtonProps {
    name: string;
    value: string;

    cdsLayout?: string;
    checked?: boolean;
    className?: string;
    register: any;
    onChange: (evt: ChangeEvent<HTMLSelectElement>) => void;
}

// This is a convenience Component that allows a radio button to be CHECKED or not CHECKED using a property
function RadioButton(props: Partial<RadioButtonProps>) {
    if (props.checked) {
        return (
            <input
                className={props.className}
                cds-layout={props.cdsLayout}
                value={props.value}
                {...props.register(props.name)}
                type="radio"
                onChange={props.onChange}
                checked
            />
        );
    }
    return (
        <input
            className={props.className}
            cds-layout={props.cdsLayout}
            value={props.value}
            {...props.register(props.name)}
            type="radio"
            onChange={props.onChange}
        />
    );
}

export default RadioButton;
