// React imports
import React from 'react';

// Library imports
import _ from 'lodash';
import { StepWizardChildProps } from 'react-step-wizard';

// App imports
import { STATUS } from '../../constants/App.constants';
import './StepNav.scss';
export interface ChildProps extends StepWizardChildProps {
    setTabStatus: (tabStatus: string[]) => void;
    tabStatus: string[];
}

function StepNav(props: ChildProps | any) {
    const { totalSteps, currentStep, goToStep, tabStatus, tabNames } = props;
    const hasInvalidStep = tabStatus.indexOf(STATUS.INVALID) !== -1;

    return (
        <div className="wizard-tab-container">
            {_.times(totalSteps, (index) => {
                return (
                    <button
                        cds-layout="p:lg"
                        className={`wizard-tab ${currentStep - 1 === index ? 'current' : ''}`}
                        cds-text="body left"
                        key={index}
                        style={{ width: `${100 / totalSteps}%` }}
                        onClick={() => {
                            if (!hasInvalidStep) {
                                goToStep(index + 1);
                            }
                        }}
                        disabled={tabStatus[index] === STATUS.DISABLED}
                    >
                        <div className={`status-bar ${currentStep - 1 === index ? 'current' : ''} ${tabStatus[index]}`}></div>
                        <span
                            cds-text="title inline"
                            cds-layout="p-r:sm"
                            className={tabStatus[index] === STATUS.DISABLED ? 'tab-number disabled' : 'tab-number'}
                        >
                            {index + 1}
                        </span>{' '}
                        {tabNames[index]}
                    </button>
                );
            })}
        </div>
    );
}

export default StepNav;
