// React imports
import { useContext } from 'react';

// Library imports
import useWebSocket from 'react-use-websocket';
import { WebSocketHook } from 'react-use-websocket/dist/lib/types';

// App imports
import { AppEnvironment, WebsocketAddress } from '../constants/App.constants';
import { Store } from '../../state-management/stores/Store';
import { STORE_SECTION_APP } from '../../state-management/reducers/App.reducer';

export const WsOperations = {
    LOGS: 'logs',
};

// Open websocket connection through @react-use-websocket
const useWsConnect = () => {
    const { state } = useContext(Store);

    // Initialize websocket protocol and address
    let protocol = window.location.protocol;
    let host = window.location.host;

    if (state[STORE_SECTION_APP].appEnv === AppEnvironment.DEV) {
        protocol = WebsocketAddress.DEFAULT_PROTOCOL;
        host = WebsocketAddress.DEV_LOCATION;
    }

    const socketUrl: string | null =
        (protocol === 'https' ? WebsocketAddress.SECURE_PROTOCOL : WebsocketAddress.DEFAULT_PROTOCOL) + `://${host}/ws`;

    const wsConnection = useWebSocket(socketUrl, {
        onOpen: () => console.log('websocket opened'),
        onClose: () => console.log('websocket closed'),
    });
    return wsConnection;
};

// Exports websocket service
export const useWebsocketService = (): WebSocketHook => {
    // open websocket connection
    const wsObj = useWsConnect();

    return wsObj;
};
