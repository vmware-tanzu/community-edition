// React imports
import React, { ReactElement, useState } from 'react';

// Library imports
import { FieldError } from 'react-hook-form';
import _ from 'lodash';
import StepWizard, { StepWizardChildProps } from 'react-step-wizard';

// App imports
import { INPUT_CHANGE } from '../../../state-management/actions/Form.actions';
import { STATUS } from '../../constants/App.constants';
import { StoreDispatch } from '../../types/types';
import StepNav from './StepNav';
import './Wizard.scss';

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
    submitForm: (data: any | undefined) => void;
    handleValueChange: (
        field: string,
        value: string,
        currentStep: number | undefined,
        errors: { [key: string]: FieldError | undefined }
    ) => void;
}
function Wizard(props: WizardProps) {
    const { tabNames, children, dispatch } = props;

    const [tabStatus, setTabStatus] = useState([
        STATUS.CURRENT,
        ..._.times(children.length - 1, () => STATUS.DISABLED),
    ]);

    const submitForm = (currentStep: number) => {
        const status = [...tabStatus];
        status[currentStep] = STATUS.TOUCHED;
        setTabStatus(status);
    };

    const handleValueChange = (
        field: string,
        value: string,
        currentStep: number | undefined,
        errors: { [key: string]: FieldError | undefined }
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
            type: INPUT_CHANGE,
            field,
            payload: value,
        });
    };

    return (
        <div className="wizard-container">
            <StepWizard
                initialStep={1}
                nav={
                    <StepNav
                        tabStatus={tabStatus}
                        tabNames={tabNames}
                    />
                }
            >
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
