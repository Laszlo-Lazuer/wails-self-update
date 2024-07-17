export function checkForUpdates(): Promise<void> {
    return window.backend.App.CheckForUpdates();
}