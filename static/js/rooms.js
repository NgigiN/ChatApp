class RoomManager {
    constructor() {
        this.rooms = [];
        this.currentUser = null;
        this.init();
    }

    async init() {
        this.setupEventListeners();
        await this.loadUserProfile();
        await this.loadRooms();
    }

    setupEventListeners() {
        // Create room modal
        document.getElementById('create-room-btn').addEventListener('click', () => {
            this.showCreateRoomModal();
        });

        document.getElementById('create-first-room-btn').addEventListener('click', () => {
            this.showCreateRoomModal();
        });

        document.getElementById('cancel-create-room').addEventListener('click', () => {
            this.hideCreateRoomModal();
        });

        document.getElementById('create-room-form').addEventListener('submit', (e) => {
            e.preventDefault();
            this.createRoom();
        });

        // Logout
        document.getElementById('logout-btn').addEventListener('click', () => {
            this.logout();
        });
    }

    async loadUserProfile() {
        try {
            // In a real app, this would fetch from /api/v1/profile
            this.currentUser = {
                id: 1,
                username: 'john_doe',
                email: 'john@example.com'
            };
            document.getElementById('username-display').textContent = `Welcome, ${this.currentUser.username}`;
        } catch (error) {
            console.error('Error loading user profile:', error);
            this.showError('Failed to load user profile');
        }
    }

    async loadRooms() {
        try {
            this.showLoading();
            
            // Mock data - in real app, fetch from /api/v1/rooms
            this.rooms = [
                {
                    id: 1,
                    name: 'Math 101 - Calculus',
                    description: 'Introduction to Calculus and Differential Equations',
                    is_private: false,
                    created_by: 1,
                    member_count: 25,
                    is_creator: true
                },
                {
                    id: 2,
                    name: 'Physics Lab',
                    description: 'Advanced Physics Laboratory Sessions',
                    is_private: true,
                    created_by: 1,
                    member_count: 12,
                    is_creator: true
                },
                {
                    id: 3,
                    name: 'Chemistry Study Group',
                    description: 'Organic Chemistry Study and Discussion',
                    is_private: false,
                    created_by: 2,
                    member_count: 18,
                    is_creator: false
                }
            ];

            this.renderRooms();
            this.updateStatistics();
        } catch (error) {
            console.error('Error loading rooms:', error);
            this.showError('Failed to load rooms');
        }
    }

    renderRooms() {
        const roomsList = document.getElementById('rooms-list');
        const loadingState = document.getElementById('loading-state');
        const emptyState = document.getElementById('empty-state');

        if (this.rooms.length === 0) {
            loadingState.classList.add('hidden');
            emptyState.classList.remove('hidden');
            return;
        }

        loadingState.classList.add('hidden');
        emptyState.classList.add('hidden');

        roomsList.innerHTML = this.rooms.map(room => `
            <div class="p-6 hover:bg-gray-50">
                <div class="flex items-center justify-between">
                    <div class="flex-1">
                        <div class="flex items-center">
                            <h4 class="text-lg font-medium text-gray-900">${room.name}</h4>
                            ${room.is_private ? '<i class="fas fa-lock text-gray-400 ml-2"></i>' : ''}
                            ${room.is_creator ? '<span class="ml-2 bg-blue-100 text-blue-800 text-xs px-2 py-1 rounded-full">Creator</span>' : ''}
                        </div>
                        <p class="text-gray-600 mt-1">${room.description}</p>
                        <div class="flex items-center mt-2 text-sm text-gray-500">
                            <i class="fas fa-users mr-1"></i>
                            <span>${room.member_count} members</span>
                            <span class="mx-2">•</span>
                            <i class="fas fa-calendar mr-1"></i>
                            <span>Created ${new Date().toLocaleDateString()}</span>
                        </div>
                    </div>
                    <div class="flex items-center space-x-2">
                        <button onclick="roomManager.joinRoom(${room.id})" 
                                class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-sm">
                            <i class="fas fa-sign-in-alt mr-1"></i>
                            Join
                        </button>
                        ${room.is_creator ? `
                            <button onclick="roomManager.manageRoom(${room.id})" 
                                    class="bg-gray-600 hover:bg-gray-700 text-white px-4 py-2 rounded-lg text-sm">
                                <i class="fas fa-cog mr-1"></i>
                                Manage
                            </button>
                        ` : `
                            <button onclick="roomManager.leaveRoom(${room.id})" 
                                    class="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-lg text-sm">
                                <i class="fas fa-sign-out-alt mr-1"></i>
                                Leave
                            </button>
                        `}
                    </div>
                </div>
            </div>
        `).join('');
    }

    updateStatistics() {
        const totalRooms = this.rooms.length;
        const privateRooms = this.rooms.filter(room => room.is_private).length;
        const activeMembers = this.rooms.reduce((sum, room) => sum + room.member_count, 0);

        document.getElementById('total-rooms').textContent = totalRooms;
        document.getElementById('private-rooms').textContent = privateRooms;
        document.getElementById('active-members').textContent = activeMembers;
    }

    showCreateRoomModal() {
        document.getElementById('create-room-modal').classList.remove('hidden');
    }

    hideCreateRoomModal() {
        document.getElementById('create-room-modal').classList.add('hidden');
        document.getElementById('create-room-form').reset();
    }

    async createRoom() {
        const formData = new FormData(document.getElementById('create-room-form'));
        const roomData = {
            name: formData.get('name'),
            description: formData.get('description'),
            is_private: formData.get('is_private') === 'on'
        };

        try {
            // In real app, POST to /api/v1/rooms
            console.log('Creating room:', roomData);
            
            // Mock success
            this.showSuccess('Room created successfully!');
            this.hideCreateRoomModal();
            await this.loadRooms(); // Reload rooms
        } catch (error) {
            console.error('Error creating room:', error);
            this.showError('Failed to create room');
        }
    }

    async joinRoom(roomId) {
        try {
            // In real app, POST to /api/v1/rooms/{id}/join
            console.log('Joining room:', roomId);
            this.showSuccess('Successfully joined room!');
            await this.loadRooms(); // Reload rooms
        } catch (error) {
            console.error('Error joining room:', error);
            this.showError('Failed to join room');
        }
    }

    async leaveRoom(roomId) {
        if (!confirm('Are you sure you want to leave this room?')) {
            return;
        }

        try {
            // In real app, DELETE /api/v1/rooms/{id}/leave
            console.log('Leaving room:', roomId);
            this.showSuccess('Successfully left room!');
            await this.loadRooms(); // Reload rooms
        } catch (error) {
            console.error('Error leaving room:', error);
            this.showError('Failed to leave room');
        }
    }

    manageRoom(roomId) {
        // Redirect to room management page
        window.location.href = `/room-management.html?id=${roomId}`;
    }

    logout() {
        // In real app, clear auth token and redirect
        window.location.href = '/login.html';
    }

    showLoading() {
        document.getElementById('loading-state').classList.remove('hidden');
        document.getElementById('empty-state').classList.add('hidden');
    }

    showError(message) {
        // Simple error display - in real app, use a proper notification system
        alert('Error: ' + message);
    }

    showSuccess(message) {
        // Simple success display - in real app, use a proper notification system
        alert('Success: ' + message);
    }
}

// Initialize the room manager when the page loads
document.addEventListener('DOMContentLoaded', () => {
    window.roomManager = new RoomManager();
});
