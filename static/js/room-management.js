class RoomManagement {
  constructor() {
    this.roomId = this.getRoomIdFromUrl();
    this.room = null;
    this.members = [];
    this.selectedMembers = new Set();
    this.init();
  }

  getRoomIdFromUrl() {
    const urlParams = new URLSearchParams(window.location.search);
    return urlParams.get('id');
  }

  async init() {
    if (!this.roomId) {
      this.showError('Invalid room ID');
      return;
    }

    this.setupEventListeners();
    await this.loadRoomData();
    await this.loadMembers();
  }

  setupEventListeners() {
    // Tab switching
    document.querySelectorAll('.tab-button').forEach(button => {
      button.addEventListener('click', (e) => {
        this.switchTab(e.target.id.replace('-tab', ''));
      });
    });

    // Member management
    document.getElementById('remove-selected-btn').addEventListener('click', () => {
      this.removeSelectedMembers();
    });

    document.getElementById('reset-room-btn').addEventListener('click', () => {
      this.resetRoom();
    });

    // Settings
    document.getElementById('room-settings-form').addEventListener('submit', (e) => {
      e.preventDefault();
      this.updateRoomSettings();
    });

    // Invites
    document.getElementById('invite-single-btn').addEventListener('click', () => {
      this.inviteSingleUser();
    });

    document.getElementById('invite-bulk-btn').addEventListener('click', () => {
      this.inviteBulkUsers();
    });

    // Modal
    document.getElementById('modal-cancel').addEventListener('click', () => {
      this.hideModal();
    });

    document.getElementById('modal-confirm').addEventListener('click', () => {
      this.confirmAction();
    });
  }

  async loadRoomData() {
    try {
      // Mock data - in real app, fetch from /api/v1/rooms/{id}
      this.room = {
        id: parseInt(this.roomId),
        name: 'Math 101 - Calculus',
        description: 'Introduction to Calculus and Differential Equations',
        is_private: false,
        created_by: 1,
        member_count: 25
      };

      this.updateRoomInfo();
    } catch (error) {
      console.error('Error loading room data:', error);
      this.showError('Failed to load room data');
    }
  }

  async loadMembers() {
    try {
      // Mock data - in real app, fetch from /api/v1/rooms/{id}/members
      this.members = [
        { id: 1, username: 'john_doe', email: 'john@example.com', joined_at: '2024-01-15' },
        { id: 2, username: 'jane_smith', email: 'jane@example.com', joined_at: '2024-01-16' },
        { id: 3, username: 'bob_wilson', email: 'bob@example.com', joined_at: '2024-01-17' },
        { id: 4, username: 'alice_brown', email: 'alice@example.com', joined_at: '2024-01-18' },
        { id: 5, username: 'charlie_davis', email: 'charlie@example.com', joined_at: '2024-01-19' }
      ];

      this.renderMembers();
    } catch (error) {
      console.error('Error loading members:', error);
      this.showError('Failed to load room members');
    }
  }

  updateRoomInfo() {
    document.getElementById('room-title').textContent = this.room.name;
    document.getElementById('room-description').textContent = this.room.description;
    document.getElementById('room-name-display').textContent = this.room.name;
    document.getElementById('member-count').textContent = this.room.member_count;

    // Update form fields
    document.getElementById('room-name-input').value = this.room.name;
    document.getElementById('room-description-input').value = this.room.description;
    document.getElementById('room-private-input').checked = this.room.is_private;
  }

  renderMembers() {
    const membersList = document.getElementById('members-list');
    membersList.innerHTML = this.members.map(member => `
            <div class="flex items-center justify-between p-3 border border-gray-200 rounded-lg hover:bg-gray-50">
                <div class="flex items-center">
                    <input type="checkbox" class="member-checkbox rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                           data-member-id="${member.id}">
                    <div class="ml-3">
                        <div class="flex items-center">
                            <span class="font-medium text-gray-900">${member.username}</span>
                            <span class="ml-2 text-sm text-gray-500">${member.email}</span>
                        </div>
                        <div class="text-sm text-gray-500">
                            Joined: ${new Date(member.joined_at).toLocaleDateString()}
                        </div>
                    </div>
                </div>
                <div class="flex items-center space-x-2">
                    <button onclick="roomManagement.removeMember(${member.id})"
                            class="text-red-600 hover:text-red-800 text-sm">
                        <i class="fas fa-user-minus"></i>
                        Remove
                    </button>
                </div>
            </div>
        `).join('');

    // Add checkbox event listeners
    document.querySelectorAll('.member-checkbox').forEach(checkbox => {
      checkbox.addEventListener('change', () => {
        this.updateSelectedMembers();
      });
    });
  }

  updateSelectedMembers() {
    this.selectedMembers.clear();
    document.querySelectorAll('.member-checkbox:checked').forEach(checkbox => {
      this.selectedMembers.add(parseInt(checkbox.dataset.memberId));
    });

    const removeBtn = document.getElementById('remove-selected-btn');
    removeBtn.disabled = this.selectedMembers.size === 0;
  }

  switchTab(tabName) {
    // Update tab buttons
    document.querySelectorAll('.tab-button').forEach(button => {
      button.classList.remove('active', 'border-blue-500', 'text-blue-600');
      button.classList.add('border-transparent', 'text-gray-500');
    });

    const activeButton = document.getElementById(`${tabName}-tab`);
    activeButton.classList.add('active', 'border-blue-500', 'text-blue-600');
    activeButton.classList.remove('border-transparent', 'text-gray-500');

    // Update tab content
    document.querySelectorAll('.tab-content').forEach(content => {
      content.classList.add('hidden');
    });

    document.getElementById(`${tabName}-content`).classList.remove('hidden');
  }

  async removeSelectedMembers() {
    if (this.selectedMembers.size === 0) return;

    this.showConfirmation(
      'Remove Members',
      `Are you sure you want to remove ${this.selectedMembers.size} member(s) from this room?`,
      () => {
        this.performRemoveMembers();
      }
    );
  }

  async performRemoveMembers() {
    try {
      // In real app, POST to /api/v1/rooms/{id}/moderation/remove for each member
      console.log('Removing members:', Array.from(this.selectedMembers));
      this.showSuccess('Members removed successfully!');
      await this.loadMembers();
    } catch (error) {
      console.error('Error removing members:', error);
      this.showError('Failed to remove members');
    }
  }

  async removeMember(memberId) {
    this.showConfirmation(
      'Remove Member',
      'Are you sure you want to remove this member from the room?',
      async () => {
        try {
          // In real app, POST to /api/v1/rooms/{id}/moderation/remove
          console.log('Removing member:', memberId);
          this.showSuccess('Member removed successfully!');
          await this.loadMembers();
        } catch (error) {
          console.error('Error removing member:', error);
          this.showError('Failed to remove member');
        }
      }
    );
  }

  resetRoom() {
    this.showConfirmation(
      'Reset Room',
      'This will remove ALL members from the room except you. This action cannot be undone. Are you sure?',
      () => {
        this.performResetRoom();
      }
    );
  }

  async performResetRoom() {
    try {
      // In real app, POST to /api/v1/rooms/{id}/moderation/reset
      console.log('Resetting room:', this.roomId);
      this.showSuccess('Room reset successfully! All members have been removed.');
      await this.loadMembers();
    } catch (error) {
      console.error('Error resetting room:', error);
      this.showError('Failed to reset room');
    }
  }

  async updateRoomSettings() {
    const formData = new FormData(document.getElementById('room-settings-form'));
    const updates = {
      name: formData.get('name'),
      description: formData.get('description'),
      is_private: formData.get('is_private') === 'on'
    };

    try {
      // In real app, PUT to /api/v1/rooms/{id}
      console.log('Updating room settings:', updates);
      this.showSuccess('Room settings updated successfully!');
      await this.loadRoomData();
    } catch (error) {
      console.error('Error updating room settings:', error);
      this.showError('Failed to update room settings');
    }
  }

  async inviteSingleUser() {
    const username = document.getElementById('invite-username').value.trim();
    if (!username) {
      this.showError('Please enter a username');
      return;
    }

    try {
      // In real app, POST to /api/v1/rooms/{id}/invites
      console.log('Inviting user:', username);
      this.showSuccess(`User ${username} invited successfully!`);
      document.getElementById('invite-username').value = '';
    } catch (error) {
      console.error('Error inviting user:', error);
      this.showError('Failed to invite user');
    }
  }

  async inviteBulkUsers() {
    const usernames = document.getElementById('invite-usernames').value.trim();
    if (!usernames) {
      this.showError('Please enter usernames');
      return;
    }

    const usernameList = usernames.split('\n').filter(u => u.trim());
    if (usernameList.length === 0) {
      this.showError('Please enter at least one username');
      return;
    }

    try {
      // In real app, POST to /api/v1/rooms/{id}/invites/bulk
      console.log('Bulk inviting users:', usernameList);
      this.showSuccess(`${usernameList.length} users invited successfully!`);
      document.getElementById('invite-usernames').value = '';
    } catch (error) {
      console.error('Error bulk inviting users:', error);
      this.showError('Failed to invite users');
    }
  }

  showConfirmation(title, message, onConfirm) {
    document.getElementById('modal-title').textContent = title;
    document.getElementById('modal-message').textContent = message;
    document.getElementById('confirmation-modal').classList.remove('hidden');

    // Store the confirm action
    this.pendingAction = onConfirm;
  }

  confirmAction() {
    if (this.pendingAction) {
      this.pendingAction();
    }
    this.hideModal();
  }

  hideModal() {
    document.getElementById('confirmation-modal').classList.add('hidden');
    this.pendingAction = null;
  }

  showError(message) {
    alert('Error: ' + message);
  }

  showSuccess(message) {
    alert('Success: ' + message);
  }
}

// Initialize the room management when the page loads
document.addEventListener('DOMContentLoaded', () => {
  window.roomManagement = new RoomManagement();
});
