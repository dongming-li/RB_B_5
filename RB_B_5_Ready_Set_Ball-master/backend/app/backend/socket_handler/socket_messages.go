package sockethandler

type socketMessage string

// Socket Messages
const (
	Ping                 socketMessage = ""
	ReceivedFriendInvite               = "receivedFriendInvite"
	CancelFriendInvite                 = "canceledFriendInvite"
	ResponseFriendInvite               = "responseFriendInvite"
	UnfriendUser                       = "unfriendUser"
	ReceivedGameInvite                 = "receivedGameInvite"
	NewGameMember                      = "newGameMember"
	GameMemberLeave                    = "gameMemberLeave"
	GameHostLeave                      = "gameHostLeave"
)
