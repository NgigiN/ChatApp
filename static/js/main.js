// static/js/main.js
class ChatApp {
  constructor() {
    this.ws = null;
    this.currentRoom = null;
    this.messageInput = document.getElementById('message-input');
    this.messageForm = document.getElementById('message-form');
    this.messagesContainer = document.getElementById('messages');
    this.currentRoomTitle = document.getElementById('current-room');
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;

    this.initializeEventListeners();
    this.connectWebSocket();
  }

  connectWebSocket() {
    this.ws = new WebSocket(`ws://${window.location.host}/ws`);

    this.ws.onopen = () => {
      this.reconnectAttempts = 0;
      this.appendMessage('Connected to server', 'system');

      if (this.currentRoom) {
        this.joinRoom(this.currentRoom);
      }
    };

    this.ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data);
        this.handleMessage(message);
      } catch (e) {
        console.error('Error parsing message:', e);
        this.appendMessage('Error receiving message', 'system');
      }
    };

    this.ws.onclose = () => {
      this.appendMessage('Disconnected from server', 'system');
      this.messageInput.disabled = true;
      this.messageForm.querySelector('button').disabled = true;

      if (this.reconnectAttempts < this.maxReconnectAttempts) {
        const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 10000);
        this.reconnectAttempts++;
        setTimeout(() => this.connectWebSocket(), delay);
      }
    };

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
      this.appendMessage('Connection error', 'system');
    };
  }

  initializeEventListeners() {
    document.querySelectorAll('.room-btn').forEach(button => {
      button.addEventListener('click', () => {
        this.joinRoom(button.dataset.room);
        document.querySelectorAll('.room-btn').forEach(btn => btn.classList.remove('active'));
        button.classList.add('active');
      });
    });

    this.messageForm.addEventListener('submit', (e) => {
      e.preventDefault();
      this.sendMessage();
    });

    this.messageInput.addEventListener('input', () => {
      const isEmpty = !this.messageInput.value.trim();
      this.messageForm.querySelector('button').disabled = isEmpty;
    });
  }

  joinRoom(roomName) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.currentRoom = roomName;
      this.currentRoomTitle.textContent = `Room: ${roomName}`;
      this.messageInput.disabled = false;
      this.messageForm.querySelector('button').disabled = !this.messageInput.value.trim();
      this.messagesContainer.innerHTML = '';

      const joinMessage = {
        type: 'join',
        room: roomName
      };

      this.ws.send(JSON.stringify(joinMessage));
      console.log(`Joining room: ${roomName}`);
    } else {
      this.appendMessage('Cannot join room - connection error', 'system');
    }
  }

  sendMessage() {
    const content = this.messageInput.value.trim();
    if (content && this.currentRoom && this.ws.readyState === WebSocket.OPEN) {
      const message = {
        type: 'message',
        room: this.currentRoom,
        content: content
      };

      try {
        this.ws.send(JSON.stringify(message));
        this.messageInput.value = '';
        this.messageForm.querySelector('button').disabled = true;
      } catch (error) {
        console.error('Error sending message:', error);
        this.appendMessage('Error sending message', 'system');
      }
    }
  }

  handleMessage(message) {
    console.log('Received message:', message);

    if (Array.isArray(message)) {
      this.messagesContainer.innerHTML = '';
      message.reverse().forEach(msg => this.renderMessage(msg));
      return;
    }

    if (message.room === this.currentRoom) {
      this.renderMessage(message);
    }
  }

  renderMessage(message) {
    const sender = message.sender || 'Anonymous';
    const content = `<strong>${sender}:</strong> ${message.content}`;
    this.appendMessage(content, 'user');
  }

  appendMessage(content, type = 'user') {
    const messageElement = document.createElement('div');
    messageElement.className = `message ${type}`;
    messageElement.innerHTML = content; // Use innerHTML to allow formatted content
    this.messagesContainer.appendChild(messageElement);
    this.messagesContainer.scrollTop = this.messagesContainer.scrollHeight;
  }
}

// Initialize the chat application
document.addEventListener('DOMContentLoaded', () => {
  new ChatApp();
});
