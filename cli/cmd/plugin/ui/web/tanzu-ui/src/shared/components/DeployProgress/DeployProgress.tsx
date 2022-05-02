// React imports
import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';

// Library imports
import { CdsButton } from '@cds/react/button';
import { LazyLog } from 'react-lazylog';
import { WebSocketHook } from 'react-use-websocket/dist/lib/types';

// App imports
import { useWebsocketService, WsOperations } from '../../services/Websocket.service';
import { NavRoutes } from '../../constants/NavRoutes.constants';
import DeployTimeline from './DeployTimeline/DeployTimeline';
import './DeployProgress.scss';

export const LogTypes = {
    LOG: 'log',
    STATUS: 'status'
};

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
            const logLine = formatLog(logData);
            setLogMessageHistory((prev) => prev.concat([logLine]));
        } else if (logData && logData.type === LogTypes.STATUS) {
            // Note: deployment status not yet being used
            setStatusMessageHistory((prev) => logData.data);
        }

    }, [websocketSvc.lastMessage]);

    // Formats a log line to include pre-pended log type (INFO, WARNING, ERROR)
    const formatLog = (log: any): string => {
        return `[${log.data.logType}] ${log.data.message}`;
    };

    // Wraps log line in a span with a custom css class if log type is Error or Warning
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
                        Creating Management Cluster on ...provider name from store
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
                        <div className="code" cds-layout="col:12">~/.config/tanzu/tkg/clusterconfigs...</div>
                    </div>
                    <div cds-layout="col:8">
                        <div cds-text="caption semibold" cds-layout="col:12 p-b:sm">Create your workload cluster</div>
                        <div className="code" cds-layout="col:12">~/.config/tanzu/tkg/clusterconfigs...</div>
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
