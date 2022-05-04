// React imports
import React from 'react';

// Library imports
import { CdsIcon } from '@cds/react/icon';
import { circleIcon, ClarityIcons, dotCircleIcon, successStandardIcon } from '@cds/core/icon';

// App imports
import { StatusMessageData } from '../DeployProgress';
import './DeployTimeline.scss';

export const DeploymentStates = {
    FAILED: 'failed',
    RUNNING: 'running',
    SUCCESSFUL: 'successful'
};

interface PropsData {
    data?: StatusMessageData
}

ClarityIcons.addIcons(circleIcon, dotCircleIcon, successStandardIcon);

function DeployTimeline(props:PropsData) {

    let currentPhaseIdx = 0;

    if (props.data?.totalPhases.length && props.data?.currentPhase) {
        currentPhaseIdx = props.data?.totalPhases.indexOf(props.data?.currentPhase);
    }

    /**
     * @method getStepState
     * @param idx - the index of each step in the list of phases
     * Determines which state should be displayed for each step of
     * the timeline component by returning the appropriate CdsIcon
     */
    const getStepState = (idx: number) => {
        if (idx === currentPhaseIdx && props.data?.status === DeploymentStates.FAILED) {
            return (
                <CdsIcon shape="error-standard" aria-label="Error"></CdsIcon>
            );
        } else if (idx < currentPhaseIdx || props.data?.status === DeploymentStates.SUCCESSFUL) {
            return (
                <CdsIcon shape="success-standard" aria-label="Completed"></CdsIcon>
            );
        } else if (idx === currentPhaseIdx) {
            return (
                <CdsIcon shape="dot-circle" aria-current="true" aria-label="Current"></CdsIcon>
            );
        } else {
            return (
                <CdsIcon shape="circle" aria-label="Not started"></CdsIcon>
            );
        }
    };

    return (
        <div>
            <ul className="cds-timeline cds-timeline-vertical">
                {props.data?.totalPhases.length &&
                    props.data?.totalPhases.map((phase: string, index: number) => (
                        <li className="cds-timeline-step" key={index}>
                            {getStepState(index)}
                            <div className="cds-timeline-step-body">
                                <span className="cds-timeline-step-title">{phase}</span>
                            </div>
                        </li>
                    ))
                }
            </ul>
        </div>
    );
}

export default DeployTimeline;