export interface ProgressMessage {
    message: string
    details?: string

    installStarting?: boolean,
    installComplete?: boolean,

    step?: string
    stepComplete?: boolean
    stepStarting?: boolean
    percentComplete?: number

    error?: boolean
    warning?: boolean
}

export interface ProgressMessenger {
    report: (msg: ProgressMessage) => void
}
