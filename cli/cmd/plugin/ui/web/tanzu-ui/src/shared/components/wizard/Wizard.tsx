// React imports
import React, { ReactElement, useState } from 'react';

// Library imports
import StepWizard, { StepWizardChildProps } from 'react-step-wizard';
import _ from 'lodash';

// App imports
import { STATUS } from '../../constants/App.constants';
import StepNav from './StepNav';
import './Wizard.scss';

interface WizardProps {
    tabNames: string[];
    children: ReactElement<any, any>[];
}
export interface StepProps extends StepWizardChildProps {
    deploy?: () => void;
    key: number;
    provider: string;
    tabStatus: STATUS[];
    setTabStatus: (status: STATUS[]) => void;
    submitForm: (data: any | undefined) => void;
    updateTabStatus: (currentStep: number | undefined, validForm: boolean) => void;
}
function Wizard(props: WizardProps) {
    const { tabNames, children } = props;

    const [tabStatus, setTabStatus] = useState([STATUS.CURRENT, ..._.times(children.length - 1, () => STATUS.DISABLED)]);

    const updateTabStatus = (currentStep: number, validForm: boolean) => {
        const status = [...tabStatus];
        status[currentStep - 1] = validForm ? STATUS.VALID : STATUS.INVALID;
        setTabStatus(status);
    };
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
                        updateTabStatus,
                    })
                )}
            </StepWizard>
        </div>
    );
}

export default Wizard;
