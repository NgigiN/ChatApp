// Mobile Navigation Component
class MobileNavigation {
  constructor() {
    this.isOpen = false;
    this.navToggle = null;
    this.navMenu = null;
    this.overlay = null;
    this.init();
  }

  init() {
    this.createMobileNav();
    this.setupEventListeners();
    this.handleResize();
  }

  createMobileNav() {
    // Create mobile nav toggle button
    this.navToggle = Utils.createElement('button', 'mobile-nav-toggle');
    this.navToggle.innerHTML = '<i class="fas fa-bars"></i>';
    this.navToggle.setAttribute('aria-label', 'Toggle navigation menu');

    // Create overlay
    this.overlay = Utils.createElement('div', 'mobile-nav-overlay');
    this.overlay.style.cssText = `
            position: fixed;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: rgba(0, 0, 0, 0.5);
            z-index: 999;
            opacity: 0;
            visibility: hidden;
            transition: all 0.3s ease;
        `;

    // Create mobile menu
    this.navMenu = Utils.createElement('div', 'mobile-nav-menu');
    this.navMenu.style.cssText = `
            position: fixed;
            top: 0;
            left: -300px;
            width: 300px;
            height: 100vh;
            background: white;
            z-index: 1000;
            transition: left 0.3s ease;
            box-shadow: 2px 0 10px rgba(0, 0, 0, 0.1);
            overflow-y: auto;
        `;

    // Add to page
    document.body.appendChild(this.overlay);
    document.body.appendChild(this.navMenu);
  }

  setupEventListeners() {
    // Toggle button click
    Utils.on(this.navToggle, 'click', () => {
      this.toggle();
    });

    // Overlay click
    Utils.on(this.overlay, 'click', () => {
      this.close();
    });

    // Close on escape key
    Utils.on(document, 'keydown', (e) => {
      if (e.key === 'Escape' && this.isOpen) {
        this.close();
      }
    });

    // Handle window resize
    Utils.on(window, 'resize', Utils.debounce(() => {
      this.handleResize();
    }, 250));
  }

  handleResize() {
    const isMobile = window.innerWidth <= 768;

    if (isMobile) {
      this.showMobileNav();
    } else {
      this.hideMobileNav();
      this.close();
    }
  }

  showMobileNav() {
    // Find the navbar and add the toggle button
    const navbar = document.querySelector('.navbar .container .flex');
    if (navbar && !navbar.querySelector('.mobile-nav-toggle')) {
      navbar.appendChild(this.navToggle);
    }

    // Move navigation items to mobile menu
    this.moveNavItemsToMobile();
  }

  hideMobileNav() {
    // Remove toggle button
    if (this.navToggle && this.navToggle.parentNode) {
      this.navToggle.parentNode.removeChild(this.navToggle);
    }

    // Move navigation items back to navbar
    this.moveNavItemsToNavbar();
  }

  moveNavItemsToMobile() {
    const navbar = document.querySelector('.navbar');
    const navItems = navbar.querySelectorAll('.nav-item, .btn');

    // Clear mobile menu
    this.navMenu.innerHTML = `
            <div style="padding: 2rem 1.5rem; border-bottom: 1px solid #e2e8f0;">
                <h3 style="font-size: 1.25rem; font-weight: 600; color: #1f2937; margin-bottom: 0.5rem;">
                    <i class="fas fa-graduation-cap mr-2"></i>
                    School Chat
                </h3>
                <p style="color: #6b7280; font-size: 0.875rem;">Connect & Learn Together</p>
            </div>
            <div style="padding: 1rem 0;">
        `;

    // Add navigation items
    navItems.forEach(item => {
      if (item.classList.contains('mobile-nav-toggle')) return;

      const mobileItem = item.cloneNode(true);
      mobileItem.style.cssText = `
                display: block;
                width: 100%;
                padding: 0.75rem 1.5rem;
                border: none;
                background: none;
                text-align: left;
                color: #374151;
                font-size: 0.875rem;
                border-bottom: 1px solid #f3f4f6;
                transition: background-color 0.2s;
            `;

      mobileItem.addEventListener('mouseenter', () => {
        mobileItem.style.backgroundColor = '#f9fafb';
      });

      mobileItem.addEventListener('mouseleave', () => {
        mobileItem.style.backgroundColor = 'transparent';
      });

      this.navMenu.querySelector('div:last-child').appendChild(mobileItem);
    });

    // Add close button
    const closeBtn = Utils.createElement('button', 'mobile-nav-close');
    closeBtn.innerHTML = '<i class="fas fa-times"></i>';
    closeBtn.style.cssText = `
            position: absolute;
            top: 1rem;
            right: 1rem;
            width: 2rem;
            height: 2rem;
            border: none;
            background: #f3f4f6;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            color: #6b7280;
            cursor: pointer;
            transition: all 0.2s;
        `;

    closeBtn.addEventListener('click', () => this.close());
    this.navMenu.appendChild(closeBtn);
  }

  moveNavItemsToNavbar() {
    // This would move items back to the main navbar
    // Implementation depends on your specific navbar structure
  }

  toggle() {
    if (this.isOpen) {
      this.close();
    } else {
      this.open();
    }
  }

  open() {
    this.isOpen = true;
    this.navMenu.style.left = '0';
    this.overlay.style.opacity = '1';
    this.overlay.style.visibility = 'visible';
    document.body.style.overflow = 'hidden';

    // Update toggle button
    this.navToggle.innerHTML = '<i class="fas fa-times"></i>';
  }

  close() {
    this.isOpen = false;
    this.navMenu.style.left = '-300px';
    this.overlay.style.opacity = '0';
    this.overlay.style.visibility = 'hidden';
    document.body.style.overflow = '';

    // Update toggle button
    this.navToggle.innerHTML = '<i class="fas fa-bars"></i>';
  }

  // Public API
  isMobileNavVisible() {
    return window.innerWidth <= 768;
  }

  destroy() {
    if (this.navToggle && this.navToggle.parentNode) {
      this.navToggle.parentNode.removeChild(this.navToggle);
    }
    if (this.overlay && this.overlay.parentNode) {
      this.overlay.parentNode.removeChild(this.overlay);
    }
    if (this.navMenu && this.navMenu.parentNode) {
      this.navMenu.parentNode.removeChild(this.navMenu);
    }
  }
}

// Auto-initialize on pages that need it
document.addEventListener('DOMContentLoaded', () => {
  if (document.querySelector('.navbar')) {
    window.mobileNav = new MobileNavigation();
  }
});

// Make it available globally
window.MobileNavigation = MobileNavigation;
