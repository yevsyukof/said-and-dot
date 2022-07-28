package db_middleware

import (
	"github.com/google/uuid"
	"said-and-dot-backend/internal/api/entities"
	api_errors "said-and-dot-backend/internal/api/routes/errors"
	"said-and-dot-backend/internal/database"
)

func GetUserByID(id string, db database.Database) (*entities.User, error) {
	var userExists bool
	if err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM Users WHERE id = $1)", id).Scan(&userExists); err != nil {
		return nil, err
	} else if !userExists {
		return nil, api_errors.ErrUserDoesNotExist
	}

	requiredUser := new(entities.User)

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	requiredUser.ID = parsedID

	if err := db.QueryRow(
		"SELECT username, first_name, last_name, email FROM Users WHERE id = $1", requiredUser.ID).Scan(
		&requiredUser.Username, &requiredUser.FirstName, &requiredUser.LastName, &requiredUser.Email); err != nil {
		return nil, api_errors.ErrDatabaseError
	}
	return requiredUser, nil
}

func GetUserTweets(user *entities.User, db database.Database) (*[]entities.Tweet, error) {
	var foundTweetsRecords database.Rows

	userTweets := new([]entities.Tweet)

	foundTweetsRecords, err := db.Query(
		"SELECT id, tweet, created FROM Tweets WHERE user_id = $1", user.ID)
	if err != nil {
		return nil, api_errors.ErrDatabaseError
	}
	defer foundTweetsRecords.Close()

	for foundTweetsRecords.Next() {
		var tweet entities.Tweet

		if err := foundTweetsRecords.Scan(&tweet.ID, &tweet.Tweet, &tweet.Created); err != nil {
			return nil, api_errors.ErrDatabaseError
		}
		tweet.UserID = user.ID

		//////////////
		var foundTweetFavorites database.Rows
		foundTweetFavorites, err := db.Query(
			"SELECT user_id FROM Favorites WHERE tweet_id = $1", tweet.ID)
		if err != nil {
			return nil, api_errors.ErrDatabaseError
		}
		//defer foundTweetFavorites.Close()

		for foundTweetFavorites.Next() {
			var favoritesUser uuid.UUID

			if err := foundTweetFavorites.Scan(&favoritesUser); err != nil {
				return nil, api_errors.ErrDatabaseError
			}
			tweet.Likes = append(tweet.Likes, favoritesUser)
		}
		foundTweetFavorites.Close()
		//////////////

		*userTweets = append(*userTweets, tweet)
	}
	//sort.Slice(userTweets, func(i, j int) (less bool) {
	//	return (*userTweets)[i].Created.Before((*userTweets)[j].Created)
	//})

	return userTweets, nil
}
