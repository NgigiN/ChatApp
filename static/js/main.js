// Legacy main.js - now redirects to new chat app
document.addEventListener('DOMContentLoaded', () => {
  // Check if we're on the old index.html
  if (document.querySelector('.container .rooms-container')) {
    // Load the new chat app
    const script = document.createElement('script');
    script.src = '/static/js/chat-app.js';
    document.head.appendChild(script);
  }
});