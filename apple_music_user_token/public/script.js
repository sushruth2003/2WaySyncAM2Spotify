document.addEventListener('DOMContentLoaded', function() {
    const authButton = document.getElementById('authButton');
    const statusMessage = document.getElementById('status');

    // Function to load MusicKit library asynchronously
    function loadMusicKit() {
        return new Promise((resolve, reject) => {
            const script = document.createElement('script');
            script.src = 'https://js-cdn.music.apple.com/musickit/v1/musickit.js';
            script.onload = resolve;
            script.onerror = reject;
            document.head.appendChild(script);
            console.log('MusicKit library loaded')
        });
    }

    // Initialize MusicKit instance after loading
    loadMusicKit().then(() => {
        // Fetch developer token from the server
        fetch('/developer-token')
            .then(response => response.text())
            .then(token => {
                // Initialize MusicKit instance with the fetched token
                MusicKit.configure({
                    developerToken: token,
                    app: {
                        name: 'Your App Name',
                        build: '1.0.0'
                    }
                });
            })
            .catch(error => {
                console.error('Failed to fetch developer token:', error);
                statusMessage.textContent = 'Failed to fetch developer token. Please try again.';
            });

        authButton.addEventListener('click', async function() {
            try {
                // Request user authorization
                music =  MusicKit.getInstance();
                const userToken = await music.authorize();
                fetch('/user-token', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ userToken })
                });
                statusMessage.textContent = 'Authentication successful! You can now access Apple Music.';
            } catch (error) {
                console.error('Authentication failed:', error);
                statusMessage.textContent = 'Authentication failed. Please try again.';
            }
        });
    }).catch(error => {
        console.error('Failed to load MusicKit library:', error);
        statusMessage.textContent = 'Failed to load MusicKit library. Please try again later.';
    });
});
