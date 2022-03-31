// React imports
import { useContext } from 'react';

// Library imports
import useWebSocket from 'react-use-websocket';
import { SendJsonMessage } from 'react-use-websocket/src/lib/types';

// App imports
import { Store } from '../../state-management/stores/store';
import { AppEnvironment, WebsocketAddress } from '../constants/App.constants';

export interface WebsocketService {
    wsSendMessage: (msg: OutboundMessage) => void;
    wsLastMessage: MessageEvent | null;
}

interface OutboundMessage {
    operation: string
}

export const WsOperations = {
    LOGS: 'logs'
};

let sendJsonMessage: SendJsonMessage, lastMessage: MessageEvent | null, socketUrl: string | null;

// Exports websocket service
export const useWebsocketService = ():WebsocketService => {
    // open websocket connection
    useWsConnect();

    return {
        wsSendMessage: wsSendMessage,
        wsLastMessage: lastMessage
    };
};

// Open websocket connection through @react-use-websocket
const useWsConnect = () => {
    const { state } = useContext(Store);

    // Initialize websocket protocol and address
    let protocol = window.location.protocol;
    let host = window.location.host;

    if (state.app.appEnv === AppEnvironment.DEV) {
        protocol = WebsocketAddress.DEFAULT_PROTOCOL;
        host = WebsocketAddress.DEV_LOCATION;
    }

    socketUrl = (protocol === 'https' ? WebsocketAddress.SECURE_PROTOCOL : WebsocketAddress.DEFAULT_PROTOCOL) + `://${host}/ws`;

    return { sendJsonMessage, lastMessage } = useWebSocket(socketUrl, {
        onOpen: () => console.log('websocket opened'),
        onClose: () => console.log('websocket closed')
    });
};

// Send message to server via websocket
const wsSendMessage = (msg: OutboundMessage) => {
    sendJsonMessage(msg);
};
