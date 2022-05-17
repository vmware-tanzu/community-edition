// React imports
import React, { useState, useEffect, useContext } from 'react';
import { Link } from 'react-router-dom';

// Library imports
import { CdsButton } from '@cds/react/button';
import { LazyLog } from 'react-lazylog';
import { WebSocketHook } from 'react-use-websocket/dist/lib/types';

// App imports
import DeployTimeline from './DeployTimeline/DeployTimeline';
import { NavRoutes } from '../../constants/NavRoutes.constants';
import { retrieveProviderInfo, ProviderData } from '../../services/Provider.service';
import { Store } from '../../../state-management/stores/Store';
import { useWebsocketService, WsOperations } from '../../services/Websocket.service';
import './DeployProgress.scss';

export const LogTypes = {
    LOG: 'log',
    STATUS: 'status'
};

export interface LogMessage {
    currentPhase: string,
    logType: string,
    message: string
}

export interface StatusMessage {
    type: string,
    data: StatusMessageData
}

export interface StatusMessageData {
    status: string,
    message: string,
    currentPhase: string,
    totalPhases: Array<string>,
    successfulPhases: Array<string>
}

function DeployProgress() {
    const { state } = useContext(Store);
    const provider: string = state.data.deployments['provider'];
    const providerData: ProviderData = retrieveProviderInfo(provider);

    let websocketSvc: WebSocketHook = useWebsocketService();

    const [statusMessageHistory, setStatusMessageHistory] = useState<StatusMessageData>();
    const [logMessageHistory, setLogMessageHistory] = useState<Array<string>>(['Awaiting log output...']);

    // Send message to websocket to start streaming cluster creation logs and status
    useEffect(
        () => {
            websocketSvc.sendJsonMessage({
                operation: WsOperations.LOGS
            });
        }, [] // eslint-disable-line react-hooks/exhaustive-deps
    );

    // Processes each message by type ('log' or 'status') and routes to appropriate handlers/state
    useEffect(() => {
        const lastMessage: MessageEvent | null = websocketSvc.lastMessage;
        const logData = (lastMessage) ? JSON.parse(lastMessage.data) : null;

        if (logData && logData.type === LogTypes.LOG) {
            const logLine = formatLog(logData.data);
            setLogMessageHistory((prev) => prev.concat([logLine]));
        } else if (logData && logData.type === LogTypes.STATUS) {
            setStatusMessageHistory((prev) => logData.data);
        }

    }, [websocketSvc.lastMessage]);

    /**
     * @method formatLog
     * @param log - individual log object
     * Formats individual log line to prepend log type (INFO, WARNING, ERROR);
     * returns log string
     */
    const formatLog = (log: LogMessage): string => {
        return `[${log.logType}] ${log.message}`;
    };

    // Wraps log line in a span with a custom css class if log type is Error or Warning
    /**
     * @method setLogLineCssClass
     * @param e - log line HTML element (in string format)
     * Applies a custom css class to a log line HTML element for the purpose of color formatting;
     * returns HTML span element with css classname applied
     */
    const setLogLineCssClass = (e: string) => {
        let className: string = 'info';

        if (e.indexOf('[ERROR]') > -1)  {
            className = 'error';
        } else if (e.indexOf('[WARNING]') > -1) {
            className = 'warning';
        }

        return <span className={className} dangerouslySetInnerHTML={{ __html: e }} />;
    };

    return (
        <>
            <div cds-layout="vertical gap:md gap@md:lg col:12">
                <div cds-layout="grid col:12 p:lg gap:md gap@md:lg" className="section-raised">
                    <div cds-text="title" cds-layout="col:12">
                        <img src={providerData.logo} className="logo logo-42" cds-layout="m-r:md" alt={`${provider} logo`}/>
                        Creating Management Cluster on {providerData.name}
                    </div>
                    <div cds-layout="col:3 p-b:md">
                        <span cds-text="section">
                            <DeployTimeline data={statusMessageHistory}/>
                        </span>
                    </div>
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

            <div cds-layout="vertical gap:md gap@md:lg col:12" >
                <div cds-layout="grid col:12 p:lg gap:lg" className="section-raised next-steps-container">
                    <div cds-text="section" cds-layout="col:12">
                        Next Steps
                    </div>
                    <div cds-layout="col:8">
                        Now that you have created a Management Cluster you can now create a Workload Cluster, then you
                        can take the next steps to deploy and manage application workloads.
                    </div>
                    <div cds-text="section" cds-layout="col:12">
                        Creating Workload Clusters
                    </div>
                    <div cds-layout="col:8">
                        <div cds-text="caption semibold" cds-layout="col:12 p-b:sm">Management Cluster configuration file</div>
                        <div className="code" cds-layout="col:12">{state.data.deployments['configPath']}</div>
                    </div>
                    <div cds-layout="col:8">
                        <div cds-text="caption semibold" cds-layout="col:12 p-b:sm">Create your workload cluster</div>
                        <div className="code" cds-layout="col:12">WC CLI Command here...</div>
                    </div>
                    <div cds-layout="col:12">
                        <Link to={NavRoutes.WORKLOAD_CLUSTER_WIZARD}>
                            <CdsButton
                                className="cluster-action-btn"
                                status="neutral">
                                Create a Workload Cluster
                            </CdsButton>
                        </Link>
                        <CdsButton
                            cds-layout="m-l:md"
                            className="cluster-action-btn"
                            action="outline">
                            Download Kubeconfig
                        </CdsButton>
                    </div>
                </div>
            </div>
        </>
    );
}

export default DeployProgress;
