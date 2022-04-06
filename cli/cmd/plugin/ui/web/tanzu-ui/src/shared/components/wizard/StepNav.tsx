// React imports
import React from 'react';

// Library imports
import _ from 'lodash';
import { StepWizardChildProps } from 'react-step-wizard';
import styled from 'styled-components';

// App imports
import { STATUS } from '../../constants/App.constants';
interface StatusProps {
    status?: STATUS;
}

interface TabProps {
    disabled?: boolean
}

export interface ChildProps extends StepWizardChildProps {
    setTabStatus: (tabStatus: string[]) => void,
    tabStatus: string[]
}

const Tab = styled.button<TabProps>`
    position: relative;
    padding: 20px;
    min-width: 200px;
    border: none;
    border-right: 1px solid #0f181c;
    border-bottom: 1px solid #0f181c;
    display: inline-block;
    background-color: #1c2b32;
    margin-bottom: 30px;
    font-size: 12px;
    font-family: 'Metropolis';
    color: white;
    &:hover {
        cursor: pointer;
        background-color: #304250;
    }
    &:last-child {
        border-right: none;
    }
`;
const StatusBar = styled.div<StatusProps>`
    position: absolute;
    background-color: ${props => {
        if (props.status === STATUS.CURRENT) {
            return '#49AFD9';
        } else if (props.status === STATUS.VALID) {
            return '#60B515';
        } else if (props.status === STATUS.INVALID){
            return 'red';
        } else {
            return '#324F61';
        }
    }};
    height: ${props => props.status === STATUS.CURRENT ? '2px' : '1px'};
    left: 0;
    right: 0;
    top: 0;
`;
const TabNumber = styled.span`
    font-size: 28px;
    padding-right: 10px;
    line-height: 36px;
`;


const TAB_NAMES = ['AWS Credentials', 'Cluster settings', 'Regions and resources', 'Configuration', 'Go!'];

function StepNav(props: ChildProps | any) {
    const { totalSteps, goToStep, currentStep, tabStatus } = props;
    const hasInvalidStep = tabStatus.indexOf(STATUS.INVALID) !== -1;
    return (
        <div>
            {
                _.times(totalSteps, (index) => {
                    return <Tab key={index} onClick={() => {
                        if (!hasInvalidStep) {
                            goToStep(index + 1);
                        }
                
                    }} disabled={tabStatus[index] === STATUS.DISABLED }>
                        <StatusBar status={
                            tabStatus[index] === STATUS.INVALID ? 
                                STATUS.INVALID : (index + 1 === currentStep ? 
                                    STATUS.CURRENT : tabStatus[index])}/>
                        <TabNumber>{ index+1 }</TabNumber> {TAB_NAMES[index]}
                    </Tab>;
                })
            }
        </div>
    );
}

export default StepNav;
