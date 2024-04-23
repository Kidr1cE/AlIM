package store

import "AlIM/pkg/tcp"

var roomHistoryCache = make(map[int][]*tcp.Message)

func AddMessageToRoomHistory(roomID int, message *tcp.Message) {
	roomHistoryCache[roomID] = append(roomHistoryCache[roomID], message)
}

func GetRoomHistory(roomID int) []*tcp.Message {
	return roomHistoryCache[roomID]
}
