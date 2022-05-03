// React imports
import React, { useEffect, useState } from 'react';

// Library imports
import { CdsIcon } from '@cds/react/icon';
import { circleIcon, ClarityIcons, dotCircleIcon, successStandardIcon } from '@cds/core/icon';

// App imports
import { StatusMessageData } from '../DeployProgress';
import './DeployTimeline.scss';

interface CurrentStatus {
    msg: string;
    status: string;
    currentPhase: string;
    finishedCount: number;
    totalCount: number;
}

interface PropsData {
    data?: StatusMessageData
}

ClarityIcons.addIcons(circleIcon, dotCircleIcon, successStandardIcon);

function DeployTimeline(props:PropsData) {

    const [phases, setPhases] = useState<Array<string>>([]);
    const [currentPhaseIdx, setCurrentPhaseIdx] = useState<number>(0);
    const [currentStatus, setCurrentStatus] = useState<CurrentStatus>({
        msg: '',
        status: '',
        currentPhase: '',
        finishedCount: 0,
        totalCount: 0,
    });

    useEffect(()=>{
        if (props.data) {
            parseStatusMsg(props.data);
        }
    },[props]); // eslint-disable-line react-hooks/exhaustive-deps

    /**
     * @method parseStatusMsg
     * @param msg - the latest deployment status message returned from the deployment progress websocket
     * Parses message data and compares current phase status to list of total deployment
     * phases to calculate overl
     */
    const parseStatusMsg = (msg: StatusMessageData) => {
        // TODO: show a default spinner if status not yet set
        if (msg.status) {
            setPhases(msg.totalPhases);

            setCurrentStatus(prevState => ({
                ...prevState,
                msg: msg.message,
                status: msg.status
            }));

            if (msg.currentPhase && phases.length) {
                setCurrentStatus(prevState => ({
                    ...prevState,
                    currentPhase: msg.currentPhase
                }));
                setCurrentPhaseIdx(phases.indexOf(currentStatus.currentPhase));
            }

            if (currentStatus.status === 'successful') {
                // currentStatus.finishedCount = currentStatus.totalCount;
                setCurrentStatus(prevState => ({
                    ...prevState,
                    finishedCount: prevState.totalCount
                }));
                setCurrentPhaseIdx(phases.length);
            } else if (currentStatus.status !== 'failed') {
                setCurrentStatus(prevState => ({
                    ...prevState,
                    finishedCount: Math.max(0, msg.totalPhases.indexOf(currentStatus.currentPhase))
                }));
            }

            setCurrentStatus(prevState => ({
                ...prevState,
                totalCount: msg.totalPhases ? msg.totalPhases.length : 0
            }));
        }
    };

    /**
     * @method getStepState
     * @param idx - the index of each step in the list of phases
     * Determines which state should be displayed for each step of
     * the timeline component by returning the appropriate CdsIcon
     */
    const getStepState = (idx: number) => {
        if (idx === currentPhaseIdx && currentStatus.status === 'failed') {
            return (
                <CdsIcon shape="error-standard" aria-label="Error"></CdsIcon>
            );
        } else if (idx < currentPhaseIdx || currentStatus.status === 'successful') {
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
                {phases.length &&
                    phases.map((phase: string, index: number) => (
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