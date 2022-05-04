interface ICLUSTER_STATUS {
  type: "CLUSTER_STATUS";
  payload: IClusterStatus;
}

interface IERROR_MESSAGE {
  type: "ERROR";
  payload: string;
}

interface ICLUSTER_LOGS {
  type: "CLUSTER_LOGS";
  payload: string[];
}

interface ICLUSTER_STARTED {
  type: "CLUSTER_STARTED";
  payload: boolean;
}

interface IFETCH_KUBECONFIG {
  type: "FETCH_KUBECONFIG";
  payload: string;
}

interface IPROVISION_INGRESS {
  type: "PROVISION_INGRESS";
  payload: boolean;
}

interface ICLUSTER_STATS {
  type: "CLUSTER_STATS";
  payload: IClusterResourceStats;
}

export type Actions =
  | ICLUSTER_STATUS
  | IERROR_MESSAGE
  | ICLUSTER_LOGS
  | ICLUSTER_STARTED
  | IFETCH_KUBECONFIG
  | IPROVISION_INGRESS
  | ICLUSTER_STATS;

export enum ClusterStates {
  UNKNOWN = "Unknown",
  NOT_EXISTS = "NotExists",
  CREATING = "Creating",
  INITIALIZING = "Initializing",
  RUNNING = "Running",
  DELETING = "Deleting",
  DELETED = "Deleted",
  ERROR = "Error",
};

export interface IClusterStatus {
  status?: ClusterStates;
  description?: string;
  isError?: boolean;
  errorMessage?: string;
  output?: string;
  stats?: IClusterResourceStats;
}

export const initialUnknownStatus: IClusterStatus = {
  status: ClusterStates.UNKNOWN,
  description: "",
  isError: false, 
  errorMessage: "",
  output: "",
}

export interface IClusterResourceStats {
  id?: string;
  memory?: IClusterMemoryStats;
  cpu?: IClusterCpuStats;
}

export interface IClusterMemoryStats {
  used?: number;
  total?: number;
  usage?: number;
}

export interface IClusterCpuStats {
  cpu_delta?: number;
  system_cpu_delta?: number;
  number_cpus?: number;
  usage?: number;
}

export interface IAppState {
  clusterStatus: IClusterStatus;
  error: string;
  logs: string[];
  isClusterStarted: boolean;
  kubeconfig: string;
  isIngressProvisioned: boolean;
  stats?: IClusterResourceStats;
}

export const emptyClusterStats: IClusterResourceStats = {
}

export const initialState: IAppState = {
  clusterStatus: initialUnknownStatus,
  error: "",
  logs: [],
  isClusterStarted: false,
  kubeconfig: "",
  isIngressProvisioned: false,
  stats: emptyClusterStats,
};

export const globalAppStateReducer = (state: IAppState, action: Actions) => {
  switch (action.type) {
    case "CLUSTER_STATUS":
      window.localStorage.setItem("appState", JSON.stringify({ value: { ...state, clusterStatus: action.payload } }));
      return { ...state, clusterStatus: action.payload };
    case "ERROR":
      window.localStorage.setItem("appState", JSON.stringify({ value: { ...state, error: action.payload } }));
      return { ...state, error: action.payload };
    case "CLUSTER_LOGS":
      window.localStorage.setItem("appState", JSON.stringify({ value: { ...state, logs: action.payload } }));
      return { ...state, logs: action.payload };
    case "CLUSTER_STARTED":
      window.localStorage.setItem("appState", JSON.stringify({ value: { ...state, isClusterStarted: action.payload } }));
      return { ...state, isClusterStarted: action.payload };
    case "FETCH_KUBECONFIG":
      window.localStorage.setItem("appState", JSON.stringify({ value: { ...state, kubeconfig: action.payload } }));
      return { ...state, kubeconfig: action.payload };
    case "PROVISION_INGRESS":
      window.localStorage.setItem("appState", JSON.stringify({ value: { ...state, isIngressProvisioned: action.payload } }));
      return { ...state, isIngressProvisioned: action.payload };
    case "CLUSTER_STATS":
      window.localStorage.setItem("appState", JSON.stringify({ value: { ...state, stats: action.payload } }));
      return { ...state, stats: action.payload };
    }
};