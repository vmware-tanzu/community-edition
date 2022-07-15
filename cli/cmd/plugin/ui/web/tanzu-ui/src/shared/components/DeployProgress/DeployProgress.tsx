// React imports
import React, { useState, useEffect, useContext, ReactElement } from 'react';
import { Link } from 'react-router-dom';

// Library imports
import { CdsButton } from '@cds/react/button';

import { LazyLog } from 'react-lazylog';
import { WebSocketHook } from 'react-use-websocket/dist/lib/types';

// App imports
import { AppFeature, featureAvailable } from '../../services/AppConfiguration.service';
import { DeploymentStates, DeploymentTypes } from '../../constants/Deployment.constants';
import DeployTimeline from './DeployTimeline/DeployTimeline';
import { NavRoutes } from '../../constants/NavRoutes.constants';
import PageNotification, { Notification, NotificationStatus } from '../PageNotification/PageNotification';
import { retrieveProviderInfo, ProviderData } from '../../services/Provider.service';
import { Store } from '../../../state-management/stores/Store';
import { STORE_SECTION_DEPLOYMENT } from '../../../state-management/reducers/Deployment.reducer';
import { useWebsocketService, WsOperations } from '../../services/Websocket.service';
import './DeployProgress.scss';
import { CdsIcon } from '@cds/react/icon';

export const LogTypes = {
    LOG: 'log',
    PROGRESS: 'progress',
};

export interface LogMessage {
    currentPhase: string;
    logType: string;
    message: string;
}

export interface StatusMessageData {
    status: string;
    message: string;
    currentPhase: string;
    totalPhases: Array<string>;
    successfulPhases: Array<string>;
}

function DeployProgress() {
    const { state } = useContext(Store);
    const clusterType: string = state[STORE_SECTION_DEPLOYMENT].deployments['type'];

    const websocketSvc: WebSocketHook = useWebsocketService();

    const [statusMessageHistory, setStatusMessageHistory] = useState<StatusMessageData>();
    const [logMessageHistory, setLogMessageHistory] = useState<Array<string>>(['Awaiting log output...']);
    const [notification, setNotification] = useState<Notification | null>(null);

    // Send message to websocket to start streaming cluster creation logs and status
    useEffect(
        () => {
            websocketSvc.sendJsonMessage({
                operation: WsOperations.LOGS,
            });
        },
        [] // eslint-disable-line react-hooks/exhaustive-deps
    );

    // Processes each message by type ('log' or 'status') and routes to appropriate handlers/state
    useEffect(() => {
        const lastMessage: MessageEvent | null = websocketSvc.lastMessage;
        const logData = lastMessage ? JSON.parse(lastMessage.data) : null;

        if (logData && logData.type === LogTypes.LOG) {
            const logLine = formatLog(logData.data);
            setLogMessageHistory((prev) => prev.concat([logLine]));
        } else if (logData && logData.type === LogTypes.PROGRESS) {
            handleDeploymentProgress(logData.data);
        }
    }, [websocketSvc.lastMessage]); // eslint-disable-line react-hooks/exhaustive-deps

    /**
     * Handler for storing last message indicating deployment progress.
     * Displays page notification once cluster creation success or failure is reached.
     */
    function handleDeploymentProgress(statusData: StatusMessageData) {
        const latestStatusMsg: StatusMessageData = statusData;

        if (statusData.status === DeploymentStates.SUCCESSFUL) {
            setNotification({
                status: NotificationStatus.SUCCESS,
                message: 'Cluster creation has completed successfully.',
            } as Notification);
        } else if (statusData.status === DeploymentStates.FAILED) {
            // failed status message does not include failed phase; substitute with current phase or unknown if missing
            latestStatusMsg.currentPhase = statusMessageHistory?.currentPhase || 'unknown';
            setNotification({
                status: NotificationStatus.DANGER,
                message: 'Cluster creation has failed. See logs for more details.',
            } as Notification);

            console.log(`Cluster creation failed on phase: ${latestStatusMsg.currentPhase}`);
        }

        setStatusMessageHistory((prev) => statusData);
    }

    /**
     * Formats individual log line to prepend log type (INFO, WARNING, ERROR);
     * returns log string
     */
    function formatLog(log: LogMessage): string {
        return `[${log.logType}] ${log.message}`;
    }

    /**
     * Wraps log line in a span with a custom css class if log type is Error or Warning;
     * returns HTML span element with css classname applied
     */
    function setLogLineCssClass(e: string): ReactElement {
        let className = 'info';

        if (e.indexOf('[ERROR]') > -1) {
            className = 'error';
        } else if (e.indexOf('[WARNING]') > -1) {
            className = 'warning';
        }

        return <span className={className} dangerouslySetInnerHTML={{ __html: e }} />;
    }

    // sets notification to null to dismiss alert
    function dismissAlert() {
        setNotification(null);
    }

    const workloadClusterSupport = featureAvailable(AppFeature.WORKLOAD_CLUSTER_SUPPORT);

    return (
        <>
            <div cds-layout="vertical gap:md gap@md:lg col:12">
                <div cds-layout="grid col:12 p:lg gap:md gap@md:lg" className="section-raised">
                    <div cds-text="title" cds-layout="col:12">
                        {displayPageTitle()}
                    </div>
                    <PageNotification notification={notification} closeCallback={dismissAlert}></PageNotification>
                    {displayDeployTimeline()}
                    <div className="log-container" cds-layout="col:9">
                        <LazyLog
                            selectableLines
                            formatPart={(log: string) => {
                                return setLogLineCssClass(log);
                            }}
                            text={logMessageHistory.join('\n')}
                        />
                    </div>
                    <nav cds-layout="col:12">
                        <Link to="/">
                            <CdsButton action="outline">Back to Welcome</CdsButton>
                        </Link>
                    </nav>
                </div>
            </div>
            {displayNextSteps()}
        </>
    );

    // Displays page title; style differs depending on creation of management or unmanaged cluster type
    function displayPageTitle() {
        if (clusterType === DeploymentTypes.UNMANAGED_CLUSTER) {
            return (
                <span>
                    <CdsIcon cds-layout="m-r:sm" shape="computer" size="xl" className="icon-blue"></CdsIcon>
                    Creating Unmanaged Cluster
                </span>
            );
        } else {
            const provider: string = state[STORE_SECTION_DEPLOYMENT].deployments['provider'];
            const providerData: ProviderData = retrieveProviderInfo(provider);

            return (
                <span>
                    <img src={providerData.logo} className="logo logo-42" cds-layout="m-r:md" alt={`${provider} logo`} />
                    Creating Management Cluster on {providerData.name}
                </span>
            );
        }
    }

    // Displays deploy timeline for management cluster creation only
    function displayDeployTimeline() {
        if (clusterType === DeploymentTypes.MANAGEMENT_CLUSTER) {
            return (
                <div cds-layout="col:3 p-b:md">
                    <span cds-text="section">
                        <DeployTimeline data={statusMessageHistory} />
                    </span>
                </div>
            );
        }
        return;
    }

    // Renders markup for next steps section of page
    function displayNextSteps() {
        if (clusterType === DeploymentTypes.MANAGEMENT_CLUSTER) {
            return (
                <>
                    <div cds-layout="vertical gap:md gap@md:lg col:12">
                        <div cds-layout="grid col:12 p:lg gap:lg" className="section-raised next-steps-container">
                            <div cds-text="section" cds-layout="col:12">
                                Next steps
                            </div>
                            <div cds-layout="col:8" cds-text="body">
                                Now that you have created a Management Cluster you can now create a Workload Cluster, then you can take the
                                next steps to deploy and manage application workloads.
                            </div>
                            <div cds-text="subsection" cds-layout="col:12">
                                Creating Workload Clusters
                            </div>
                            <div cds-layout="col:8" cds-text="body">
                                The following Management Cluster configuration file will be used for creating Workload Clusters:
                            </div>
                            <div cds-layout="col:8">
                                <div cds-text="caption semibold" cds-layout="col:12 p-b:sm">
                                    Management Cluster configuration file
                                </div>
                                <div className="code-block" cds-layout="col:12">
                                    <code cds-text="code">{state[STORE_SECTION_DEPLOYMENT].deployments['configPath']}</code>
                                </div>
                                <CdsButton cds-layout="m-t:sm" size="sm" action="outline">
                                    Export configuration
                                </CdsButton>
                            </div>
                            <div cds-layout="col:12">
                                <CdsButton
                                    cds-layout="m-r:md"
                                    status="primary"
                                    onClick={() => {
                                        window.open(
                                            'https://tanzucommunityedition.io/docs/main/getting-started/#deploy-a-workload-cluster',
                                            '_blank'
                                        );
                                    }}
                                >
                                    Get started with Workload Clusters
                                </CdsButton>
                                {workloadClusterSupport && (
                                    <Link to={NavRoutes.WORKLOAD_CLUSTER_WIZARD}>
                                        <CdsButton className="cluster-action-btn" status="neutral">
                                            Create a Workload Cluster
                                        </CdsButton>
                                    </Link>
                                )}
                            </div>
                        </div>
                    </div>
                </>
            );
        }
        return;
    }
}

export default DeployProgress;
