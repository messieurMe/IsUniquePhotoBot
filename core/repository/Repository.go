package repository

type Repository interface {
	AddImage(userId int64, messageId int, hash uint64) error

	FindExisting(userId int64, hash uint64) (int, error)

	ClearUserData(userId string) error

	GetGroupIdByUser(userId int64) (int64, error)

	GetLastGroupMessageByUser(userId int64) (int, error)

	SetUserGroupAndLastMessage(userId int64, groupId int64, lastMessage int) error
}
