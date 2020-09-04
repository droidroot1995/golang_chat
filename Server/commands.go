package main

const (
	CMD_PUB = "/pub"
	CMD_SUB = "/sub"
	CMD_LIST = "/list"
	CMD_HELP = "/help"
	CMD_HIST = "/hist"
	CMD_SET_NAME = "/set_name"
	CMD_UNSUB = "/unsub"
	CMD_QUIT = "/quit"

	SUCCESS_SUB = "Successfully subscribed to room:"
	SUCCESS_UNSUB = "Successfully unsubscribed from room:"
	SUCCESS_NAME_CHANGE = "Name successfully changed:"

	FAIL_PUB_RNE = "Cannot send message to this room. Room doesn't exists.\n"
	FAIL_PUB_NE = "Cannot send message to this room. You aren't subscribed to it.\n"
	FAIL_SUB_RNE = "Cannot subscribe to this room. Room doesn't exist.\n"
	FAIL_SUB_AE = "Cannot subscribe to this room. User with such name already exists.\n"
	FAIL_UNSUB_RNE = "Cannot unsubscribe from this room. Room doesn't exist.\n"
	FAIL_UNSUB_NE = "Cannot unsubscribe from this room. User with such name doesn't exist.\n"
	FAIL_HIST_GET_NS = "Cannot get history from this room. You aren't subscribed to it.\n"
	FAIL_HIST_GET_NE = "Cannot get history from this room. Room doesn't exist.\n"

	MSG_WELCOME = "Welcome to our server. If you want to see the list available commands, type '/help'\nIf you want to connect to any room, set your username, please\n"
	MSG_HELP = `List of available commands:
/sub room_name - subscript to the room
/pub room_name msg - publish msg to the room
/list - show list of rooms
/help - show this list
/hist room_name - get messages history from the room
/unsub room_name - unsubscribe from the room
/set_name name - set your username
/quit - disconnect from the server`
)
