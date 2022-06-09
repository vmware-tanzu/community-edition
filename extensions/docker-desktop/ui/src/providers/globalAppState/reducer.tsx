interface ICLUSTER_STATUS {
  type: 'CLUSTER_STATUS';
  payload: IClusterStatus;
}

interface ICLUSTER_LOGS {
  type: 'CLUSTER_LOGS';
  payload: string[];
}

interface ICLUSTER_STARTED {
  type: 'CLUSTER_STARTED';
  payload: boolean;
}

interface IFETCH_KUBECONFIG {
  type: 'FETCH_KUBECONFIG';
  payload: string;
}

export type Actions =
  | ICLUSTER_STATUS
  | ICLUSTER_LOGS
  | ICLUSTER_STARTED
  | IFETCH_KUBECONFIG;

export enum ClusterStates {
  UNKNOWN = 'Unknown',
  NOT_EXISTS = 'NotExists',
  CREATING = 'Creating',
  INITIALIZING = 'Initializing',
  RUNNING = 'Running',
  DELETING = 'Deleting',
  DELETED = 'Deleted',
  ERROR = 'Error',
}

export interface IClusterStatus {
  status?: ClusterStates;
  description?: string;
  isError?: boolean;
  errorMessage?: string;
  output?: string;
}

export const initialUnknownStatus: IClusterStatus = {
  status: ClusterStates.UNKNOWN,
  description: '',
  isError: false,
  errorMessage: '',
  output: '',
};

export interface IAppState {
  clusterStatus: IClusterStatus;
  logs: string[];
  isClusterStarted: boolean;
  kubeconfig: string;
}

export const initialState: IAppState = {
  clusterStatus: initialUnknownStatus,
  logs: [],
  isClusterStarted: false,
  kubeconfig: '',
};

export const globalAppStateReducer = (state: IAppState, action: Actions) => {
  switch (action.type) {
    case 'CLUSTER_STATUS':
      window.localStorage.setItem(
        'appState',
        JSON.stringify({ value: { ...state, clusterStatus: action.payload } }),
      );
      return { ...state, clusterStatus: action.payload };
    case 'CLUSTER_LOGS':
      window.localStorage.setItem(
        'appState',
        JSON.stringify({ value: { ...state, logs: action.payload } }),
      );
      return { ...state, logs: action.payload };
    case 'CLUSTER_STARTED':
      window.localStorage.setItem(
        'appState',
        JSON.stringify({
          value: { ...state, isClusterStarted: action.payload },
        }),
      );
      return { ...state, isClusterStarted: action.payload };
    case 'FETCH_KUBECONFIG':
      window.localStorage.setItem(
        'appState',
        JSON.stringify({ value: { ...state, kubeconfig: action.payload } }),
      );
      return { ...state, kubeconfig: action.payload };
  }
};
