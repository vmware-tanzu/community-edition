// App imports
import './Wizard.scss';

import { FormAction, StoreDispatch } from '../../types/types';
// React imports
import React, { ReactElement, useState } from 'react';
import StepWizard, { StepWizardChildProps } from 'react-step-wizard';

// Library imports
import { FieldError } from 'react-hook-form';
import { STATUS } from '../../constants/App.constants';
import StepNav from './StepNav';
import _ from 'lodash';

interface WizardProps {
    tabNames: string[];
    state: { [key: string]: any };
    dispatch: StoreDispatch;
    children: ReactElement<any, any>[];
}
export interface StepProps extends StepWizardChildProps {
    tabStatus: STATUS[];
    setTabStatus: (status: STATUS[]) => void;
    key: number;
    deploy?: () => void;
    submitForm: (data: any | undefined) => void;
    handleValueChange: (
        type: string,
        field: string,
        value: any,
        currentStep: number | undefined,
        errors: { [key: string]: FieldError | undefined },
        locationData?: any
    ) => void;
    provider: string;
    setTabStatus: (status: STATUS[]) => void;
    submitForm: (data: any | undefined) => void;
    tabStatus: STATUS[];
}
function Wizard(props: WizardProps) {
    const { tabNames, children, dispatch } = props;

    const handleValueChange = (
        type: string,
        field: string,
        value: string,
        currentStep: number | undefined,
        errors: { [key: string]: FieldError | undefined },
        locationData?: any
    ) => {
        // update status bar for the wizard tab
        if (errors[field] && currentStep) {
            const status = [...tabStatus];
            status[currentStep - 1] = STATUS.INVALID;
            setTabStatus(status);
        } else if (Object.keys(errors).length === 0 && currentStep) {
            const status = [...tabStatus];
            status[currentStep - 1] = STATUS.VALID;
            setTabStatus(status);
        }
        // update the field in the data store.
        dispatch({
            type,
            field,
            payload: value,
            locationData,
        } as FormAction);
    };

    const [tabStatus, setTabStatus] = useState([STATUS.CURRENT, ..._.times(children.length - 1, () => STATUS.DISABLED)]);

    const submitForm = (currentStep: number) => {
        const status = [...tabStatus];
        status[currentStep] = STATUS.TOUCHED;
        setTabStatus(status);
    };

    return (
        <div className="wizard-container">
            <StepWizard initialStep={1} nav={<StepNav tabStatus={tabStatus} tabNames={tabNames} />}>
                {props.children.map((child: any, index: number) =>
                    React.cloneElement(child, {
                        tabStatus,
                        setTabStatus,
                        key: index,
                        submitForm,
                        handleValueChange,
                    })
                )}
            </StepWizard>
        </div>
    );
}

export default Wizard;
