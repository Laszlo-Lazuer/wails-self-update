interface Window {
    backend: {
        App: {
            CheckForUpdates(): Promise<void>;
        };
    };
}