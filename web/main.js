let username = '';
let chatID = '';
let ws = null;

function joinChat() {
    username = document.getElementById('username').value.trim();
    chatID = document.getElementById('chatid').value.trim();
    if (username && chatID) {
        document.getElementById('join').style.display = 'none';
        document.getElementById('chat').style.display = 'block';
        startWebSocket();
    } else {
        alert('Please enter both username and chat room.');
    }
}

function startWebSocket() {
    const messages = document.getElementById('messages');
    const form = document.getElementById('form');
    const input = document.getElementById('input');

    ws = new WebSocket(`ws://localhost:8080/${chatID}/username?username=${username}`);

    ws.onopen = () => {
        messages.innerHTML += `<div>Joined room ${chatID}</div>`;
        console.log('Received message:');
    };

    ws.onmessage = (event) => {
        console.log('Received message:');
        try {
            const msg = JSON.parse(event.data);
            console.log('Received message:', msg);
            if (!msg.username || !msg.text || msg.chatid !== chatID) {
                console.log('Invalid message received:', msg);
                //return;
            }
            const div = document.createElement('div');
            div.textContent = `${msg.username}: ${msg.content}`;
            messages.appendChild(div);
            messages.scrollTop = messages.scrollHeight;
        } catch (err) {
            console.error('Failed to parse message:', err, 'Data:', event.data);
        }
    };

    ws.onclose = () => {
        messages.innerHTML += '<div>Disconnected from chat</div>';
    };

    form.addEventListener('submit', (event) => {
        event.preventDefault();
        const text = input.value.trim();
        if (text) {
            ws.send(JSON.stringify({ username: username, text: text, chatid: chatID }));
            input.value = '';
        }
    });
}
