// Taken from: https://github.com/kubeapps/kubeapps-dd/blob/main/client/src/globals.d.ts

interface Window {
  ddClient: {
    docker: Docker;
    extension: Extension;
    host: Host;
    desktopUI: DesktopUI;
  }
}

interface Docker {
  readonly cli: DockerCommand;
  listContainers(options?: any): Promise<any>;
  listImages(options?: any): Promise<any>;
}

interface DockerCommand {
  exec: Exec;
}

interface Extension {
  readonly vm: ExtensionVM;
  readonly host: ExtensionHost;
}

interface DesktopUI{
  readonly toast: Toast;
  readonly navigate: NavigationIntents;
}

interface ExtensionVM {
  readonly cli: ExtensionCli;
  readonly service: HttpService;
}

interface ExtensionHost {
  readonly cli: ExtensionCli;
}

interface ExtensionCli {
  exec: Exec;
}

interface HttpService {
  get(url: string): Promise<unknown>;
  post(url: string, data: any): Promise<unknown>;
  put(url: string, data: any): Promise<unknown>;
  patch(url: string, data: any): Promise<unknown>;
  delete(url: string): Promise<unknown>;
  head(url: string): Promise<unknown>;
  request(config: RequestConfig): Promise<unknown>;
}

interface Exec {
  (cmd: string, args?: string[]): Promise<ExecResult>
  (cmd: string, args: string[], options: Object): void
}

interface RequestConfig {
  url: string;
  method: string;
  headers: Record<string, string>;
  data: any;
}

interface RawExecResult {
  readonly cmd?: string;
  readonly killed?: boolean;
  readonly signal?: string;
  readonly code?: number;
  readonly stdout: string;
  readonly stderr: string;
}

interface ExecResult extends RawExecResult {
  lines(): string[];
  parseJsonLines(): any[];
  parseJsonObject(): any;
}

interface Toast {
  success(msg): void;
  warning(msg): void;
  error(msg): void;
}

interface Host {
  openExternal(url): void;
}

interface NavigationIntents {
  viewContainers(): Promise<void>;
  viewContainer(id: string): Promise<void>;
  viewContainerLogs(id: string): Promise<void>;
  viewContainerInspect(id: string): Promise<void>;
  viewContainerStats(id: string): Promise<void>;
  viewImages(): Promise<void>;
  viewImage(id: string, tag: string): Promise<void>;
  viewDevEnvironments(): Promise<void>;
  viewVolumes(): Promise<void>;
  viewVolume(volume: string): Promise<void>;
}
