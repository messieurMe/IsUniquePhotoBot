package sqlite

import (
	"IsUniquePhotoBot/core/util"
	"database/sql"
)

type SQLiteRepository struct {
	db         *sql.DB
	hashHelper *util.HashHelper
}

func NewSQLiteRepository(
	db *sql.DB,
	hashHelper *util.HashHelper,
) (*SQLiteRepository, error) {
	r := &SQLiteRepository{db: db, hashHelper: hashHelper}
	if err := r.init(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *SQLiteRepository) init() error {
	schema := `
	CREATE TABLE IF NOT EXISTS userImages (
		group_id INTEGER NOT NULL,
		message_id INTEGER NOT NULL,
		hash INTEGER NOT NULL,
		PRIMARY KEY (group_id, message_id)
	);

	CREATE INDEX IF NOT EXISTS idx_images_user_hash
		ON userImages(group_id, message_id);

	CREATE TABLE IF NOT EXISTS userGroups (
		user_id INTEGER NOT NULL,
		group_id INTEGER NOT NULL,
		last_message_id INTEGER NOT NULL,
		PRIMARY KEY (user_id)
	);

	CREATE INDEX IF NOT EXISTS idx_user_id
		ON userGroups(user_id);
	`
	_, err := r.db.Exec(schema)
	return err
}

// region: methods

func (r *SQLiteRepository) AddImage(userId int64, messageId int, hash uint64) error {
	int64Hash := int64(hash)

	_, err := r.db.Exec(`
		INSERT INTO userImages (group_id, message_id, hash)
		VALUES (?, ?, ?)
		ON CONFLICT(group_id, message_id) DO UPDATE SET
			hash = excluded.hash
		;
		`,
		userId,
		messageId,
		int64Hash,
	)
	return err
}

func (r *SQLiteRepository) FindExisting(userId int64, hash uint64) (int, error) {
	rows, err := r.db.Query(
		`
		SELECT group_id, message_id, hash
		FROM userImages
		`,
	)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var channelID int64
		var messageID int
		var rawHash int64
		if err := rows.Scan(&channelID, &messageID, &rawHash); err != nil {
			return -1, err
		}
		storedHash := uint64(rawHash)

		areSimilar := r.hashHelper.AreSimilar(hash, storedHash)
		if areSimilar {
			return messageID, nil
		}
	}
	return -1, nil
}

func (r *SQLiteRepository) ClearUserData(userId string) error {
	_, err := r.db.Exec(`
		DELETE FROM userImages
		WHERE user_id = ?;
		`,
		userId,
	)

	return err
}

func (r *SQLiteRepository) GetGroupIdByUser(userId int64) (int64, error) {
	var groupId int64

	err := r.db.QueryRow(`
		SELECT group_id
		FROM userGroups
		WHERE user_id = ?
		LIMIT 1
		`,
		userId,
	).Scan(&groupId)

	if err != nil {
		return -1, err
	}

	return groupId, nil
}

func (r *SQLiteRepository) GetLastGroupMessageByUser(userId int64) (int, error) {
	var lastMessageId int

	err := r.db.QueryRow(`
		SELECT last_message_id
		FROM userGroups
		WHERE user_id = ?
		LIMIT 1
		`,
		userId,
	).Scan(&lastMessageId)

	if err == sql.ErrNoRows {
		return 0, nil
	}

	if err != nil {
		return -1, err
	}

	return lastMessageId, nil

}

func (r *SQLiteRepository) SetUserGroupAndLastMessage(
	userId int64,
	groupId int64,
	lastMessage int,
) error {

	_, err := r.db.Exec(`
		INSERT INTO userGroups
		(user_id, group_id, last_message_id)
		VALUES
		(?, ?, ?)
		ON CONFLICT(user_id) DO UPDATE SET
			group_id = excluded.group_id,
			last_message_id = excluded.last_message_id
		`,
		userId,
		groupId,
		lastMessage,
	)
	if err != nil {
		return err
	}
	return nil
}

// end region
