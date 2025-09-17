package migrations

import (
	"database/sql"
	"fmt"
	"log"
)

type Migration struct {
	Version int
	Name    string
	Up      func(*sql.DB) error
	Down    func(*sql.DB) error
}

var migrations = []Migration{
	{
		Version: 1,
		Name:    "create_users_table",
		Up:      createUsersTable,
		Down:    dropUsersTable,
	},
	{
		Version: 2,
		Name:    "create_user_sessions_table",
		Up:      createUserSessionsTable,
		Down:    dropUserSessionsTable,
	},
	{
		Version: 3,
		Name:    "create_rooms_table",
		Up:      createRoomsTable,
		Down:    dropRoomsTable,
	},
	{
		Version: 4,
		Name:    "create_messages_table",
		Up:      createMessagesTable,
		Down:    dropMessagesTable,
	},
	{
		Version: 5,
		Name:    "create_room_members_table",
		Up:      createRoomMembersTable,
		Down:    dropRoomMembersTable,
	},
	{
		Version: 6,
		Name:    "create_indexes",
		Up:      createIndexes,
		Down:    dropIndexes,
	},
}

func RunMigrations(db *sql.DB) error {
	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get current version
	currentVersion, err := getCurrentVersion(db)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	// Run pending migrations
	for _, migration := range migrations {
		if migration.Version > currentVersion {
			log.Printf("Running migration %d: %s", migration.Version, migration.Name)
			if err := migration.Up(db); err != nil {
				return fmt.Errorf("failed to run migration %d: %w", migration.Version, err)
			}
			if err := updateVersion(db, migration.Version); err != nil {
				return fmt.Errorf("failed to update version: %w", err)
			}
			log.Printf("Migration %d completed successfully", migration.Version)
		}
	}

	return nil
}

func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			version INT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`
	_, err := db.Exec(query)
	return err
}

func getCurrentVersion(db *sql.DB) (int, error) {
	query := `SELECT COALESCE(MAX(version), 0) FROM migrations`
	var version int
	err := db.QueryRow(query).Scan(&version)
	return version, err
}

func updateVersion(db *sql.DB, version int) error {
	query := `INSERT INTO migrations (version, name) VALUES (?, ?)`
	_, err := db.Exec(query, version, fmt.Sprintf("migration_%d", version))
	return err
}

func createUsersTable(db *sql.DB) error {
	query := `
		CREATE TABLE users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			is_active BOOLEAN DEFAULT TRUE
		)`
	_, err := db.Exec(query)
	return err
}

func dropUsersTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS users")
	return err
}

func createUserSessionsTable(db *sql.DB) error {
	query := `
		CREATE TABLE user_sessions (
			id VARCHAR(255) PRIMARY KEY,
			user_id INT NOT NULL,
			token VARCHAR(255) UNIQUE NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			is_active BOOLEAN DEFAULT TRUE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`
	_, err := db.Exec(query)
	return err
}

func dropUserSessionsTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS user_sessions")
	return err
}

func createRoomsTable(db *sql.DB) error {
	query := `
		CREATE TABLE rooms (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(100) UNIQUE NOT NULL,
			description TEXT,
			is_private BOOLEAN DEFAULT FALSE,
			created_by INT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			is_active BOOLEAN DEFAULT TRUE,
			FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
		)`
	_, err := db.Exec(query)
	return err
}

func dropRoomsTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS rooms")
	return err
}

func createMessagesTable(db *sql.DB) error {
	query := `
		CREATE TABLE messages (
			id INT AUTO_INCREMENT PRIMARY KEY,
			room_id INT NOT NULL,
			user_id INT NOT NULL,
			username VARCHAR(50) NOT NULL,
			content TEXT NOT NULL,
			type VARCHAR(50) DEFAULT 'message',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`
	_, err := db.Exec(query)
	return err
}

func dropMessagesTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS messages")
	return err
}

func createRoomMembersTable(db *sql.DB) error {
	query := `
		CREATE TABLE room_members (
			id INT AUTO_INCREMENT PRIMARY KEY,
			room_id INT NOT NULL,
			user_id INT NOT NULL,
			joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			is_active BOOLEAN DEFAULT TRUE,
			UNIQUE KEY unique_room_user (room_id, user_id),
			FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`
	_, err := db.Exec(query)
	return err
}

func dropRoomMembersTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS room_members")
	return err
}

func createIndexes(db *sql.DB) error {
	indexes := []string{
		"CREATE INDEX idx_users_username ON users(username)",
		"CREATE INDEX idx_users_email ON users(email)",
		"CREATE INDEX idx_users_is_active ON users(is_active)",
		"CREATE INDEX idx_user_sessions_token ON user_sessions(token)",
		"CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id)",
		"CREATE INDEX idx_user_sessions_expires_at ON user_sessions(expires_at)",
		"CREATE INDEX idx_rooms_name ON rooms(name)",
		"CREATE INDEX idx_rooms_created_by ON rooms(created_by)",
		"CREATE INDEX idx_rooms_is_active ON rooms(is_active)",
		"CREATE INDEX idx_messages_room_id ON messages(room_id)",
		"CREATE INDEX idx_messages_user_id ON messages(user_id)",
		"CREATE INDEX idx_messages_created_at ON messages(created_at)",
		"CREATE INDEX idx_room_members_room_id ON room_members(room_id)",
		"CREATE INDEX idx_room_members_user_id ON room_members(user_id)",
		"CREATE INDEX idx_room_members_is_active ON room_members(is_active)",
	}

	for _, index := range indexes {
		if _, err := db.Exec(index); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	return nil
}

func dropIndexes(db *sql.DB) error {
	indexes := []string{
		"DROP INDEX IF EXISTS idx_users_username ON users",
		"DROP INDEX IF EXISTS idx_users_email ON users",
		"DROP INDEX IF EXISTS idx_users_is_active ON users",
		"DROP INDEX IF EXISTS idx_user_sessions_token ON user_sessions",
		"DROP INDEX IF EXISTS idx_user_sessions_user_id ON user_sessions",
		"DROP INDEX IF EXISTS idx_user_sessions_expires_at ON user_sessions",
		"DROP INDEX IF EXISTS idx_rooms_name ON rooms",
		"DROP INDEX IF EXISTS idx_rooms_created_by ON rooms",
		"DROP INDEX IF EXISTS idx_rooms_is_active ON rooms",
		"DROP INDEX IF EXISTS idx_messages_room_id ON messages",
		"DROP INDEX IF EXISTS idx_messages_user_id ON messages",
		"DROP INDEX IF EXISTS idx_messages_created_at ON messages",
		"DROP INDEX IF EXISTS idx_room_members_room_id ON room_members",
		"DROP INDEX IF EXISTS idx_room_members_user_id ON room_members",
		"DROP INDEX IF EXISTS idx_room_members_is_active ON room_members",
	}

	for _, index := range indexes {
		if _, err := db.Exec(index); err != nil {
			return fmt.Errorf("failed to drop index: %w", err)
		}
	}

	return nil
}
