package store

import (
	"github.com/aryanA101a/legoshichat-backend/model"
	"gofr.dev/pkg/datastore"
	"gofr.dev/pkg/gofr"
)

type friend struct {
}

type FriendStore interface {
	AddFriend(ctx *gofr.Context,senderId,recieverId string)error
	GetFriends(ctx *gofr.Context,userId string)(*[]model.User,error)

}

func NewFriendStore(db *datastore.SQLClient) FriendStore {
	f:= friend{}
	f.init(db)
	return f
}

func (f friend) init(db *datastore.SQLClient) {
	createFriendsTable(db)
}

func (f friend)AddFriend(ctx *gofr.Context,senderId,recieverId string)error{
	_, err := ctx.DB().ExecContext(ctx, `INSERT INTO friends (account_id1, account_id2)
	VALUES (LEAST($1, $2)::UUID, GREATEST($1, $2)::UUID)
	ON CONFLICT DO NOTHING`, senderId,recieverId)
	return err
}

func (f friend)GetFriends(ctx *gofr.Context,userId string)(*[]model.User,error){
	query:=`SELECT u.id, u.name,u.phoneNumber
	FROM friends f
	JOIN accounts u ON u.id = f.account_id2
	WHERE f.account_id1 = $1
	UNION
	SELECT u.id, u.name,u.phoneNumber
	FROM friends f
	JOIN accounts u ON u.id = f.account_id1
	WHERE f.account_id2 = $1`
	
	rows, err := ctx.DB().QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]model.User, 0)

	for rows.Next() {
		var user model.User

		err = rows.Scan(&user.ID, &user.Name, &user.PhoneNumber,)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func createFriendsTable(db *datastore.SQLClient) error {
	query := `CREATE TABLE IF NOT EXISTS friends (
		account_id1 UUID NOT NULL,
		account_id2 UUID NOT NULL,
		FOREIGN KEY (account_id1) REFERENCES accounts(id),
		FOREIGN KEY (account_id2) REFERENCES accounts(id),
		PRIMARY KEY (account_id1, account_id2)
	);`
	_, err := db.Exec(query)
	return err
}