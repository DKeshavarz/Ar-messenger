import { useEffect, useRef, useState } from 'react';
import './App.css';

interface ChatMessage {
  username: string;
  content: string;
  room_name: string;
}

interface JoinFormData {
  username: string;
  chatid: string;
}

function App() {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [input, setInput] = useState('');
  const [joinData, setJoinData] = useState<JoinFormData>({ username: '', chatid: '' });
  const [isJoined, setIsJoined] = useState(false);
  const ws = useRef<WebSocket | null>(null);
  const [connected, setConnected] = useState(false);

  const joinChat = () => {
    const { username, chatid } = joinData;
    if (username.trim() && chatid.trim()) {
      setIsJoined(true);
      startWebSocket(username.trim(), chatid.trim());
    } else {
      alert('Please enter both username and chat room.');
    }
  };

  const startWebSocket = (username: string, chatid: string) => {
    const baseurl =  `${import.meta.env.VITE_SERVER_URL}`;
    const baseport = `${import.meta.env.VITE_SERVER_PORT}`;
    console.log("urlllllll: ", baseurl, ":", baseport)
    const wsUrl = `ws://${baseurl}:${baseport}/${chatid}/username?username=${encodeURIComponent(username)}`;
    ws.current = new WebSocket(wsUrl);
    
    ws.current.onopen = () => {
      setConnected(true);
      setMessages(prev => [...prev, { 
        username: 'System', 
        content: `Joined room ${chatid}`, 
        room_name: chatid 
      }]);
    };
    
    ws.current.onclose = () => {
      setConnected(false);
      setMessages(prev => [...prev, { 
        username: 'System', 
        content: 'Disconnected from chat', 
        room_name: chatid 
      }]);
    };
    
    ws.current.onerror = () => {
      setConnected(false);
      setMessages(prev => [...prev, { 
        username: 'System', 
        content: 'Connection error', 
        room_name: chatid 
      }]);
    };
    
    ws.current.onmessage = (event) => {
      try {
        const msg: ChatMessage = JSON.parse(event.data);
        setMessages(prev => [...prev, msg]);
      } catch (err) {
        console.error('Failed to parse message:', err, 'Data:', event.data);
        // Fallback for plain text
        setMessages(prev => [...prev, { 
          username: 'System', 
          content: event.data, 
          room_name: chatid 
        }]);
      }
    };
  };

  const sendMessage = () => {
    if (input.trim() && ws.current && connected && isJoined) {
      const message = {
        username: joinData.username,
        text: input.trim(),
        chatid: joinData.chatid
      };
      ws.current.send(JSON.stringify(message));
      console.log(message)
      setInput('');
    }
  };

  const leaveChat = () => {
    if (ws.current) {
      ws.current.close();
    }
    setIsJoined(false);
    setConnected(false);
    setMessages([]);
    setJoinData({ username: '', chatid: '' });
  };

  useEffect(() => {
    return () => {
      ws.current?.close();
    };
  }, []);

  return (
    <div className="chat-container">
      <h2>Multi-Room Chat</h2>
      
      {!isJoined ? (
        <div className="join-form">
          <div className="form-group">
            <label htmlFor="username">Username:</label>
            <input
              type="text"
              id="username"
              value={joinData.username}
              onChange={e => setJoinData(prev => ({ ...prev, username: e.target.value }))}
              placeholder="Enter username"
              autoComplete="off"
            />
          </div>
          <div className="form-group">
            <label htmlFor="chatid">Chat Room:</label>
            <input
              type="text"
              id="chatid"
              value={joinData.chatid}
              onChange={e => setJoinData(prev => ({ ...prev, chatid: e.target.value }))}
              placeholder="Enter room (e.g., room1)"
              autoComplete="off"
            />
          </div>
          <button onClick={joinChat}>Join</button>
        </div>
      ) : (
        <div className="chat-interface">
          <div className="chat-header">
            <span>Room: {joinData.chatid} | User: {joinData.username}</span>
            <button onClick={leaveChat} className="leave-btn">Leave</button>
          </div>
          
          <div className="chat-messages">
            {messages.map((msg, idx) => (
              <div 
                key={idx} 
                className={`message ${msg.username === joinData.username ? 'my-message' : 'other-message'}`}
              >
                <span className="username">{msg.username}:</span>
                <span className="content">{msg.content}</span>
              </div>
            ))}
          </div>
          
          <div className="chat-input">
            <input
              type="text"
              value={input}
              onChange={e => setInput(e.target.value)}
              onKeyDown={e => e.key === 'Enter' && sendMessage()}
              placeholder={connected ? 'Type a message...' : 'Connecting...'}
              disabled={!connected}
            />
            <button onClick={sendMessage} disabled={!connected || !input.trim()}>
              Send
            </button>
          </div>
          
          {!connected && <div className="chat-status">Connecting to server...</div>}
        </div>
      )}
    </div>
  );
}

export default App;