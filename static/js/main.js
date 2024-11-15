// static/js/main.js
class ChatApp {
  constructor () {
    this.ws = null;
    this.currentRoom = null;
    this.messageInput = document.getElementById('message-input');
    this.messageForm = document.getElementById('message-form');
    this.messagesContainer = document.getElementById('messages');
    this.currentRoomTitle = document.getElementById('current-room');

    this.initializeEventListeners();
    this.connectWebSocket();
  }

  connectWebSocket () {
    this.ws = new window.WebSocket(`ws://${window.location.host}/ws`);

    this.ws.onopen = () => {
      this.appendMessage('Connected to server', 'system');
    };

    this.ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      this.handleMessage(message);
    };

    this.ws.onclose = () => {
      this.appendMessage('Disconnected from server', 'system');
      this.messageInput.disabled = true;
      this.messageForm.querySelector('button').disabled = true;
      setTimeout(() => this.connectWebSocket(), 5000);
    };
  }

  initializeEventListeners () {
    // Room button listeners
    document.querySelectorAll('.room-btn').forEach(button => {
      button.addEventListener('click', () => {
        this.joinRoom(button.dataset.room);
        document.querySelectorAll('.room-btn').forEach(btn => btn.classList.remove('active'));
        button.classList.add('active');
      });
    });

    // Message form listener
    this.messageForm.addEventListener('submit', (e) => {
      e.preventDefault();
      this.sendMessage();
    });
  }

  joinRoom (roomName) {
    if (this.ws && this.ws.readyState === window.WebSocket.OPEN) {
      this.currentRoom = roomName;
      this.currentRoomTitle.textContent = `Room: ${roomName}`;
      this.messageInput.disabled = false;
      this.messageForm.querySelector('button').disabled = false;
      this.messagesContainer.innerHTML = '';

      this.ws.send(JSON.stringify({
        type: 'join',
        room: roomName
      }));
    }
  }

  sendMessage () {
    const content = this.messageInput.value.trim();
    if (content && this.currentRoom) {
      const message = {
        type: 'message',
        room: this.currentRoom,
        content
      };

      this.ws.send(JSON.stringify(message));
      this.messageInput.value = '';
    }
  }

  handleMessage (message) {
    if (message.type === 'system') {
      this.appendMessage(message.content, 'system');
    } else {
      this.appendMessage(`${message.sender || 'Anonymous'}: ${message.content}`);
    }
  }

  appendMessage (content, type = 'user') {
    const messageElement = document.createElement('div');
    messageElement.className = `message ${type}`;
    messageElement.textContent = content;
    this.messagesContainer.appendChild(messageElement);
    this.messagesContainer.scrollTop = this.messagesContainer.scrollHeight;
  }
}

// Initialize the chat application
document.addEventListener('DOMContentLoaded', () => {
  const chatApp = new ChatApp();
});
