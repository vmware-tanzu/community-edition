export interface ProgressMessage {
    message: string
    details?: string

    step?: string
    stepComplete?: boolean
    percentComplete?: number

    error?: boolean
    warning?: boolean
}

export interface ProgressMessenger {
    report: (msg: ProgressMessage) => void
}
