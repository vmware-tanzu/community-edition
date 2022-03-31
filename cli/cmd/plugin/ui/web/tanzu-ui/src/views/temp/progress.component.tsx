// React imports
import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';

// Library imports
import { CdsButton } from '@cds/react/button';
import { LazyLog } from 'react-lazylog';
import styled from 'styled-components';

// App imports
import { useWebsocketService, WebsocketService, WsOperations } from '../../shared/services/Websocket.service';
import './progress.component.scss';

const LogViewContainer = styled.section`
    height: 500px;
    width: 900px;
`;

export const LogTypes = {
    LOG: 'log',
    STATUS: 'status'
};

interface StatusMessage {
    type: string,
    data: StatusMessageData
}

interface StatusMessageData {
    data: {
        status: string,
        currentPhase: string,
        totalPhases: Array<string>,
        successfulPhases: Array<string>
    }
}

function WelcomeComponent(props: any) {
    let websocketSvc: WebsocketService = useWebsocketService();

    const [statusMessageHistory, setStatusMessageHistory] = useState<Array<StatusMessage>>([]);
    const [logMessageHistory, setLogMessageHistory] = useState<Array<string>>(['Displaying logs']);

    // Send message to websocket to start streaming cluster creation logs and status
    useEffect(
        () => {
            websocketSvc.wsSendMessage({
                operation: WsOperations.LOGS
            });
        }, [] // eslint-disable-line react-hooks/exhaustive-deps
    );

    // Processes each message by type ('log' or 'status') and routes to appropriate handlers/state
    useEffect(() => {
        const lastMessage: MessageEvent | null = websocketSvc.wsLastMessage;
        const logData = (lastMessage) ? JSON.parse(lastMessage.data) : null;

        if (logData && logData.type === LogTypes.LOG) {
            const logLine = formatLog(logData);
            setLogMessageHistory((prev) => prev.concat([logLine]));
        } else if (logData && logData.type === LogTypes.STATUS) {
            // Note: deployment status not yet being used
            setStatusMessageHistory((prev) => [logData.data]);
        }

    }, [websocketSvc.wsLastMessage]);

    console.log(`TODO: output statusMessageHistory as steps complete in UI ${statusMessageHistory}`);

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
            <p cds-text="title">
                Stream logs to UI demo
            </p>
            <LogViewContainer>
                <LazyLog
                    selectableLines
                    formatPart={(log: string) => {
                        return setLogLineCssClass(log);
                    }}
                    text={logMessageHistory.join('\n')}
                />
            </LogViewContainer>
            <nav>
                <Link to="/">
                    <CdsButton>Back to Welcome</CdsButton>
                </Link>
            </nav>
        </>
    );
}

export default WelcomeComponent;