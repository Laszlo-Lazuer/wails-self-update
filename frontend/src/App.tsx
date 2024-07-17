import React from 'react';
import { checkForUpdates } from './backend';

function App() {
    const handleCheckForUpdates = () => {
        checkForUpdates()
            .then(() => {
                console.log('Check for updates completed.');
            })
            .catch((err) => {
                console.error('Error checking for updates:', err);
            });
    };

    return (
        <div className="App">
            <header className="App-header">
                <h1>Welcome to MySelfUpdatingApp</h1>
                <button onClick={handleCheckForUpdates}>Check for Updates</button>
            </header>
        </div>
    );
}

export default App;