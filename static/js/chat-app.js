// Enhanced Chat Application
class ChatApp {
  constructor() {
    this.ws = null;
    this.currentRoom = null;
    this.user = null;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.reconnectDelay = 1000;
    this.messageHistory = [];
    this.isTyping = false;
    this.typingUsers = new Set();
    
    this.initializeElements();
    this.loadUserData();
    this.initializeEventListeners();
    this.connectWebSocket();
    this.setupKeyboardShortcuts();
  }

  initializeElements() {
    this.messageInput = Utils.$('#message-input');
    this.messageForm = Utils.$('#message-form');
    this.messagesContainer = Utils.$('#messages');
    this.currentRoomTitle = Utils.$('#current-room');
    this.roomsContainer = Utils.$('.rooms-list');
    this.typingIndicator = Utils.$('#typing-indicator');
    this.connectionStatus = Utils.$('#connection-status');
  }

  loadUserData() {
    this.user = Utils.getStorage('user', {
      username: 'Anonymous',
      id: Date.now()
    });
    
    // Update UI with user info
    const userDisplay = Utils.$('#user-display');
    if (userDisplay) {
      userDisplay.textContent = this.user.username;
    }
  }

  initializeEventListeners() {
    // Room selection
    Utils.on(document, 'click', '.room-btn', (e) => {
      this.selectRoom(e.target.dataset.room);
    });

    // Create room
    Utils.on(document, 'click', '#create-room-btn', () => {
      this.createRoom();
    });

    // Message form
    Utils.on(this.messageForm, 'submit', (e) => {
      e.preventDefault();
      this.sendMessage();
    });

    // Typing indicators
    Utils.on(this.messageInput, 'input', Utils.debounce(() => {
      this.handleTyping();
    }, 500));

    Utils.on(this.messageInput, 'keydown', (e) => {
      if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault();
        this.sendMessage();
      }
    });

    // Auto-resize textarea
    Utils.on(this.messageInput, 'input', () => {
      this.autoResizeTextarea();
    });

    // Connection status
    Utils.on(window, 'online', () => {
      this.updateConnectionStatus('online');
      this.connectWebSocket();
    });

    Utils.on(window, 'offline', () => {
      this.updateConnectionStatus('offline');
    });

    // Page visibility
    Utils.on(document, 'visibilitychange', () => {
      if (document.hidden) {
        this.pauseTyping();
      } else {
        this.resumeTyping();
      }
    });
  }

  setupKeyboardShortcuts() {
    Utils.on(document, 'keydown', (e) => {
      // Ctrl/Cmd + K to focus message input
      if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
        e.preventDefault();
        this.messageInput.focus();
      }
      
      // Escape to clear message input
      if (e.key === 'Escape') {
        this.messageInput.value = '';
        this.messageInput.blur();
      }
    });
  }

  connectWebSocket() {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      return;
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;
    
    this.ws = new WebSocket(wsUrl);
    this.updateConnectionStatus('connecting');

    this.ws.onopen = () => {
      this.reconnectAttempts = 0;
      this.updateConnectionStatus('connected');
      Utils.showNotification('Connected to server', 'success');
      
      if (this.currentRoom) {
        this.joinRoom(this.currentRoom);
      }
    };

    this.ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data);
        this.handleMessage(message);
      } catch (error) {
        console.error('Error parsing message:', error);
        Utils.showNotification('Error receiving message', 'error');
      }
    };

    this.ws.onclose = (event) => {
      this.updateConnectionStatus('disconnected');
      this.messageInput.disabled = true;
      Utils.$('#send-btn').disabled = true;
      
      if (event.code !== 1000) { // Not a normal closure
        this.handleReconnection();
      }
    };

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
      this.updateConnectionStatus('error');
      Utils.showNotification('Connection error', 'error');
    };
  }

  handleReconnection() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      const delay = Math.min(this.reconnectDelay * Math.pow(2, this.reconnectAttempts), 30000);
      this.reconnectAttempts++;
      
      Utils.showNotification(`Reconnecting in ${delay/1000}s... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`, 'warning');
      
      setTimeout(() => {
        this.connectWebSocket();
      }, delay);
    } else {
      Utils.showNotification('Failed to reconnect. Please refresh the page.', 'error');
    }
  }

  updateConnectionStatus(status) {
    if (!this.connectionStatus) return;
    
    const statusMap = {
      connecting: { text: 'Connecting...', class: 'status-connecting' },
      connected: { text: 'Connected', class: 'status-online' },
      disconnected: { text: 'Disconnected', class: 'status-offline' },
      error: { text: 'Error', class: 'status-error' },
      online: { text: 'Online', class: 'status-online' },
      offline: { text: 'Offline', class: 'status-offline' }
    };
    
    const statusInfo = statusMap[status] || statusMap.disconnected;
    this.connectionStatus.textContent = statusInfo.text;
    this.connectionStatus.className = `status-indicator ${statusInfo.class}`;
  }

  selectRoom(roomName) {
    if (this.currentRoom === roomName) return;
    
    // Update UI
    document.querySelectorAll('.room-btn').forEach(btn => {
      btn.classList.remove('active');
    });
    
    const selectedBtn = document.querySelector(`[data-room="${roomName}"]`);
    if (selectedBtn) {
      selectedBtn.classList.add('active');
    }
    
    this.currentRoom = roomName;
    this.currentRoomTitle.textContent = `Room: ${roomName}`;
    this.messageInput.disabled = false;
    this.messageInput.focus();
    
    // Clear messages and join room
    this.messagesContainer.innerHTML = '';
    this.joinRoom(roomName);
  }

  createRoom() {
    const roomName = Utils.$('#new-room-name').value.trim();
    if (!roomName) {
      Utils.showNotification('Please enter a room name', 'warning');
      return;
    }
    
    if (roomName.length < 3) {
      Utils.showNotification('Room name must be at least 3 characters', 'warning');
      return;
    }
    
    // Create room button if it doesn't exist
    let roomBtn = document.querySelector(`[data-room="${roomName}"]`);
    if (!roomBtn) {
      roomBtn = Utils.createElement('button', 'room-btn', roomName);
      roomBtn.dataset.room = roomName;
      this.roomsContainer.appendChild(roomBtn);
    }
    
    this.selectRoom(roomName);
    Utils.$('#new-room-name').value = '';
  }

  joinRoom(roomName) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      const joinMessage = {
        type: 'join',
        room: roomName,
        user: this.user
      };
      
      this.ws.send(JSON.stringify(joinMessage));
      console.log(`Joining room: ${roomName}`);
      
      // Request message history
      this.requestMessageHistory(roomName);
    } else {
      Utils.showNotification('Cannot join room - connection error', 'error');
    }
  }

  requestMessageHistory(roomName) {
    // In a real app, this would fetch from the API
    // For now, we'll use mock data
    const mockHistory = [
      {
        id: 1,
        sender: 'System',
        content: `Welcome to ${roomName}!`,
        timestamp: new Date(Date.now() - 60000),
        type: 'system'
      }
    ];
    
    this.messageHistory = mockHistory;
    this.renderMessages(mockHistory);
  }

  sendMessage() {
    const content = this.messageInput.value.trim();
    if (!content || !this.currentRoom || this.ws.readyState !== WebSocket.OPEN) {
      return;
    }
    
    const message = {
      type: 'message',
      room: this.currentRoom,
      content: content,
      sender: this.user.username,
      timestamp: new Date()
    };
    
    try {
      this.ws.send(JSON.stringify(message));
      this.messageInput.value = '';
      this.autoResizeTextarea();
      this.stopTyping();
    } catch (error) {
      console.error('Error sending message:', error);
      Utils.showNotification('Error sending message', 'error');
    }
  }

  handleMessage(message) {
    console.log('Received message:', message);
    
    if (Array.isArray(message)) {
      // Message history
      this.messageHistory = message;
      this.renderMessages(message);
      return;
    }
    
    if (message.type === 'typing') {
      this.handleTypingMessage(message);
      return;
    }
    
    if (message.room === this.currentRoom) {
      this.messageHistory.push(message);
      this.renderMessage(message);
    }
  }

  renderMessages(messages) {
    this.messagesContainer.innerHTML = '';
    messages.forEach(message => this.renderMessage(message));
    this.scrollToBottom();
  }

  renderMessage(message) {
    const messageElement = Utils.createElement('div', 'message');
    
    if (message.type === 'system') {
      messageElement.classList.add('system');
      messageElement.innerHTML = `
        <div class="message-content">
          <div class="message-text">${this.escapeHtml(message.content)}</div>
        </div>
      `;
    } else {
      const isOwnMessage = message.sender === this.user.username;
      const avatarColor = Utils.generateAvatarColor(message.sender);
      const initials = Utils.getInitials(message.sender);
      
      messageElement.innerHTML = `
        <div class="message-avatar" style="background-color: ${avatarColor}">
          ${initials}
        </div>
        <div class="message-content">
          <div class="message-header">
            <span class="message-sender">${this.escapeHtml(message.sender)}</span>
            <span class="message-time">${Utils.formatTime(new Date(message.timestamp))}</span>
          </div>
          <div class="message-text">${this.formatMessageContent(message.content)}</div>
        </div>
      `;
      
      if (isOwnMessage) {
        messageElement.classList.add('own-message');
      }
    }
    
    this.messagesContainer.appendChild(messageElement);
    this.scrollToBottom();
  }

  formatMessageContent(content) {
    // Convert URLs to links
    const urlRegex = /(https?:\/\/[^\s]+)/g;
    content = content.replace(urlRegex, '<a href="$1" target="_blank" rel="noopener noreferrer">$1</a>');
    
    // Convert newlines to <br>
    content = content.replace(/\n/g, '<br>');
    
    return content;
  }

  escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
  }

  handleTyping() {
    if (!this.isTyping && this.currentRoom) {
      this.isTyping = true;
      this.sendTypingStatus(true);
    }
    
    // Reset typing status after 3 seconds of inactivity
    clearTimeout(this.typingTimeout);
    this.typingTimeout = setTimeout(() => {
      this.stopTyping();
    }, 3000);
  }

  stopTyping() {
    if (this.isTyping) {
      this.isTyping = false;
      this.sendTypingStatus(false);
    }
  }

  sendTypingStatus(isTyping) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN && this.currentRoom) {
      const message = {
        type: 'typing',
        room: this.currentRoom,
        user: this.user.username,
        isTyping: isTyping
      };
      
      this.ws.send(JSON.stringify(message));
    }
  }

  handleTypingMessage(message) {
    if (message.user === this.user.username) return;
    
    if (message.isTyping) {
      this.typingUsers.add(message.user);
    } else {
      this.typingUsers.delete(message.user);
    }
    
    this.updateTypingIndicator();
  }

  updateTypingIndicator() {
    if (!this.typingIndicator) return;
    
    if (this.typingUsers.size === 0) {
      this.typingIndicator.style.display = 'none';
    } else {
      const users = Array.from(this.typingUsers);
      let text = '';
      
      if (users.length === 1) {
        text = `${users[0]} is typing...`;
      } else if (users.length === 2) {
        text = `${users[0]} and ${users[1]} are typing...`;
      } else {
        text = `${users[0]} and ${users.length - 1} others are typing...`;
      }
      
      this.typingIndicator.textContent = text;
      this.typingIndicator.style.display = 'block';
    }
  }

  autoResizeTextarea() {
    if (this.messageInput) {
      this.messageInput.style.height = 'auto';
      this.messageInput.style.height = Math.min(this.messageInput.scrollHeight, 120) + 'px';
    }
  }

  scrollToBottom() {
    this.messagesContainer.scrollTop = this.messagesContainer.scrollHeight;
  }

  pauseTyping() {
    this.stopTyping();
  }

  resumeTyping() {
    // Resume any pending typing operations
  }

  // Public API methods
  getCurrentRoom() {
    return this.currentRoom;
  }

  getUser() {
    return this.user;
  }

  isConnected() {
    return this.ws && this.ws.readyState === WebSocket.OPEN;
  }

  disconnect() {
    if (this.ws) {
      this.ws.close(1000, 'User disconnected');
    }
  }
}

// Initialize the chat application when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
  window.chatApp = new ChatApp();
});

// Handle page unload
window.addEventListener('beforeunload', () => {
  if (window.chatApp) {
    window.chatApp.disconnect();
  }
});
