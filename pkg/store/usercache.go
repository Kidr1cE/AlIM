package store

import (
	"crypto/sha256"
	"fmt"
	"sync/atomic"
)

var (
	friendShipCache = make(map[int][]int)
	CurrentUserID   = new(int32)
)

func GetNextUserID() int {
	return int(atomic.AddInt32(CurrentUserID, 1))
}

func IsFriend(userID1, userID2 int) bool {
	// Check if userID2 is in the list of userID1's friends
	for _, friendID := range friendShipCache[userID1] {
		if friendID == userID2 {
			return true
		}
	}

	// If we didn't find userID2 in the list of userID1's friends, they are not friends
	return false
}

func GetFriends(userID int) []int {
	// Return the list of friends for the given userID
	return friendShipCache[userID]
}

func AddFriend(userID1, userID2 int) {
	// Add userID2 to the list of userID1's friends
	friendShipCache[userID1] = append(friendShipCache[userID1], userID2)
}

func GenerateRoomID(userID1, userID2 int) int {
	// Ensure the smaller userID always comes first
	if userID2 < userID1 {
		userID1, userID2 = userID2, userID1
	}

	// Generate a unique roomID for the two users
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%d%d", userID1, userID2)))
	return int(hash.Sum(nil)[0])
}
