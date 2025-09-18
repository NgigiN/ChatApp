// Utility Functions for Chat App
class Utils {
  // DOM Utilities
  static $(selector) {
    return document.querySelector(selector);
  }

  static $$(selector) {
    return document.querySelectorAll(selector);
  }

  static createElement(tag, className = '', content = '') {
    const element = document.createElement(tag);
    if (className) element.className = className;
    if (content) element.textContent = content;
    return element;
  }

  // Local Storage Utilities
  static setStorage(key, value) {
    try {
      localStorage.setItem(key, JSON.stringify(value));
      return true;
    } catch (error) {
      console.error('Error saving to localStorage:', error);
      return false;
    }
  }

  static getStorage(key, defaultValue = null) {
    try {
      const item = localStorage.getItem(key);
      return item ? JSON.parse(item) : defaultValue;
    } catch (error) {
      console.error('Error reading from localStorage:', error);
      return defaultValue;
    }
  }

  static removeStorage(key) {
    try {
      localStorage.removeItem(key);
      return true;
    } catch (error) {
      console.error('Error removing from localStorage:', error);
      return false;
    }
  }

  // API Utilities
  static async apiRequest(url, options = {}) {
    const defaultOptions = {
      headers: {
        'Content-Type': 'application/json',
      },
    };

    const token = this.getStorage('authToken');
    if (token) {
      defaultOptions.headers.Authorization = `Bearer ${token}`;
    }

    const config = {
      ...defaultOptions,
      ...options,
      headers: {
        ...defaultOptions.headers,
        ...options.headers,
      },
    };

    try {
      const response = await fetch(url, config);
      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || `HTTP error! status: ${response.status}`);
      }

      return { success: true, data };
    } catch (error) {
      console.error('API request failed:', error);
      return { success: false, error: error.message };
    }
  }

  // Notification System
  static showNotification(message, type = 'info', duration = 3000) {
    const notification = this.createElement('div', `notification notification-${type}`, message);
    notification.style.cssText = `
      position: fixed;
      top: 20px;
      right: 20px;
      padding: 12px 20px;
      border-radius: 8px;
      color: white;
      font-weight: 500;
      z-index: 10000;
      animation: slideInRight 0.3s ease-out;
      max-width: 300px;
      word-wrap: break-word;
    `;

    const colors = {
      success: '#10b981',
      error: '#ef4444',
      warning: '#f59e0b',
      info: '#3b82f6'
    };

    notification.style.backgroundColor = colors[type] || colors.info;

    document.body.appendChild(notification);

    setTimeout(() => {
      notification.style.animation = 'slideOutRight 0.3s ease-in';
      setTimeout(() => {
        if (notification.parentNode) {
          notification.parentNode.removeChild(notification);
        }
      }, 300);
    }, duration);
  }

  // Loading States
  static showLoading(element, text = 'Loading...') {
    if (typeof element === 'string') {
      element = this.$(element);
    }

    if (element) {
      element.disabled = true;
      element.dataset.originalText = element.textContent;
      element.innerHTML = `<span class="loading"></span> ${text}`;
    }
  }

  static hideLoading(element) {
    if (typeof element === 'string') {
      element = this.$(element);
    }

    if (element) {
      element.disabled = false;
      element.textContent = element.dataset.originalText || 'Submit';
      delete element.dataset.originalText;
    }
  }

  // Form Utilities
  static serializeForm(form) {
    const formData = new FormData(form);
    const data = {};
    for (let [key, value] of formData.entries()) {
      data[key] = value;
    }
    return data;
  }

  static validateForm(form, rules = {}) {
    const errors = {};
    const formData = this.serializeForm(form);

    for (const [field, rule] of Object.entries(rules)) {
      const value = formData[field];

      if (rule.required && (!value || value.trim() === '')) {
        errors[field] = `${field} is required`;
        continue;
      }

      if (value && rule.minLength && value.length < rule.minLength) {
        errors[field] = `${field} must be at least ${rule.minLength} characters`;
        continue;
      }

      if (value && rule.maxLength && value.length > rule.maxLength) {
        errors[field] = `${field} must be no more than ${rule.maxLength} characters`;
        continue;
      }

      if (value && rule.pattern && !rule.pattern.test(value)) {
        errors[field] = rule.message || `${field} format is invalid`;
        continue;
      }
    }

    return {
      isValid: Object.keys(errors).length === 0,
      errors
    };
  }

  static showFormErrors(form, errors) {
    // Clear previous errors
    form.querySelectorAll('.error-message').forEach(el => el.remove());
    form.querySelectorAll('.form-control').forEach(el => {
      el.classList.remove('error');
    });

    // Show new errors
    for (const [field, message] of Object.entries(errors)) {
      const input = form.querySelector(`[name="${field}"]`);
      if (input) {
        input.classList.add('error');
        const errorDiv = this.createElement('div', 'error-message', message);
        errorDiv.style.cssText = 'color: #ef4444; font-size: 0.875rem; margin-top: 0.25rem;';
        input.parentNode.appendChild(errorDiv);
      }
    }
  }

  // Time Utilities
  static formatTime(date) {
    const now = new Date();
    const diff = now - date;
    const seconds = Math.floor(diff / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);
    const days = Math.floor(hours / 24);

    if (seconds < 60) return 'just now';
    if (minutes < 60) return `${minutes}m ago`;
    if (hours < 24) return `${hours}h ago`;
    if (days < 7) return `${days}d ago`;

    return date.toLocaleDateString();
  }

  static formatDateTime(date) {
    return new Date(date).toLocaleString();
  }

  // String Utilities
  static truncate(str, length = 50) {
    if (str.length <= length) return str;
    return str.substring(0, length) + '...';
  }

  static capitalize(str) {
    return str.charAt(0).toUpperCase() + str.slice(1);
  }

  static slugify(str) {
    return str
      .toLowerCase()
      .replace(/[^\w\s-]/g, '')
      .replace(/[\s_-]+/g, '-')
      .replace(/^-+|-+$/g, '');
  }

  // URL Utilities
  static getQueryParam(name) {
    const urlParams = new URLSearchParams(window.location.search);
    return urlParams.get(name);
  }

  static setQueryParam(name, value) {
    const url = new URL(window.location);
    url.searchParams.set(name, value);
    window.history.replaceState({}, '', url);
  }

  // Debounce and Throttle
  static debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
      const later = () => {
        clearTimeout(timeout);
        func(...args);
      };
      clearTimeout(timeout);
      timeout = setTimeout(later, wait);
    };
  }

  static throttle(func, limit) {
    let inThrottle;
    return function (...args) {
      if (!inThrottle) {
        func.apply(this, args);
        inThrottle = true;
        setTimeout(() => inThrottle = false, limit);
      }
    };
  }

  // Event Utilities
  static on(element, event, handler) {
    if (typeof element === 'string') {
      element = this.$(element);
    }
    if (element) {
      element.addEventListener(event, handler);
    }
  }

  static off(element, event, handler) {
    if (typeof element === 'string') {
      element = this.$(element);
    }
    if (element) {
      element.removeEventListener(event, handler);
    }
  }

  // Animation Utilities
  static fadeIn(element, duration = 300) {
    if (typeof element === 'string') {
      element = this.$(element);
    }
    if (element) {
      element.style.opacity = '0';
      element.style.display = 'block';

      let start = performance.now();
      const animate = (timestamp) => {
        const elapsed = timestamp - start;
        const progress = Math.min(elapsed / duration, 1);
        element.style.opacity = progress;

        if (progress < 1) {
          requestAnimationFrame(animate);
        }
      };
      requestAnimationFrame(animate);
    }
  }

  static fadeOut(element, duration = 300) {
    if (typeof element === 'string') {
      element = this.$(element);
    }
    if (element) {
      let start = performance.now();
      const animate = (timestamp) => {
        const elapsed = timestamp - start;
        const progress = Math.min(elapsed / duration, 1);
        element.style.opacity = 1 - progress;

        if (progress < 1) {
          requestAnimationFrame(animate);
        } else {
          element.style.display = 'none';
        }
      };
      requestAnimationFrame(animate);
    }
  }

  // Modal Utilities
  static showModal(modalId) {
    const modal = this.$(modalId);
    if (modal) {
      modal.classList.remove('hidden');
      modal.style.display = 'flex';
      document.body.style.overflow = 'hidden';
    }
  }

  static hideModal(modalId) {
    const modal = this.$(modalId);
    if (modal) {
      modal.classList.add('hidden');
      modal.style.display = 'none';
      document.body.style.overflow = '';
    }
  }

  // WebSocket Utilities
  static createWebSocket(url, options = {}) {
    const ws = new WebSocket(url);

    ws.onopen = options.onOpen || (() => console.log('WebSocket connected'));
    ws.onclose = options.onClose || (() => console.log('WebSocket disconnected'));
    ws.onerror = options.onError || ((error) => console.error('WebSocket error:', error));
    ws.onmessage = options.onMessage || ((event) => console.log('WebSocket message:', event.data));

    return ws;
  }

  // Device Detection
  static isMobile() {
    return window.innerWidth <= 768;
  }

  static isTablet() {
    return window.innerWidth > 768 && window.innerWidth <= 1024;
  }

  static isDesktop() {
    return window.innerWidth > 1024;
  }

  // Copy to Clipboard
  static async copyToClipboard(text) {
    try {
      await navigator.clipboard.writeText(text);
      this.showNotification('Copied to clipboard!', 'success');
      return true;
    } catch (error) {
      console.error('Failed to copy to clipboard:', error);
      this.showNotification('Failed to copy to clipboard', 'error');
      return false;
    }
  }

  // File Utilities
  static formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  }

  static isValidImageFile(file) {
    const validTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'];
    return validTypes.includes(file.type);
  }

  // Color Utilities
  static generateAvatarColor(name) {
    const colors = [
      '#ef4444', '#f97316', '#f59e0b', '#eab308',
      '#84cc16', '#22c55e', '#10b981', '#14b8a6',
      '#06b6d4', '#0ea5e9', '#3b82f6', '#6366f1',
      '#8b5cf6', '#a855f7', '#d946ef', '#ec4899'
    ];

    let hash = 0;
    for (let i = 0; i < name.length; i++) {
      hash = name.charCodeAt(i) + ((hash << 5) - hash);
    }

    return colors[Math.abs(hash) % colors.length];
  }

  static getInitials(name) {
    return name
      .split(' ')
      .map(word => word.charAt(0))
      .join('')
      .toUpperCase()
      .substring(0, 2);
  }

  static escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
  }

  // Modal utilities
  static showModal(modalId) {
    const modal = this.$(modalId);
    if (modal) {
      modal.classList.add('show');
      document.body.style.overflow = 'hidden';
    }
  }

  static hideModal(modalId) {
    const modal = this.$(modalId);
    if (modal) {
      modal.classList.remove('show');
      document.body.style.overflow = '';
    }
  }
}

// Add CSS animations
const style = document.createElement('style');
style.textContent = `
  @keyframes slideInRight {
    from { transform: translateX(100%); opacity: 0; }
    to { transform: translateX(0); opacity: 1; }
  }

  @keyframes slideOutRight {
    from { transform: translateX(0); opacity: 1; }
    to { transform: translateX(100%); opacity: 0; }
  }

  .form-control.error {
    border-color: #ef4444;
    box-shadow: 0 0 0 3px rgb(239 68 68 / 0.1);
  }

  .notification {
    box-shadow: 0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1);
  }
`;
document.head.appendChild(style);

// Make Utils available globally
window.Utils = Utils;
