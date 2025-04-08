package packets

var PLAYER_IDENTIFICATION = byte(0x00)
var SERVER_IDENTIFICATION = byte(0x00)
var PING = byte(0x01)
var LEVEL_INITIALIZE = byte(0x02)
var LEVEL_DATA_CHUNK = byte(0x03)
var LEVEL_FINALIZE = byte(0x04)
var SET_BLOCK = byte(0x06)
var SPAWN_PLAYER = byte(0x07)
var PLAYER_TELEPORT = byte(0x08)
var PLAYER_MOVEMENT = byte(0x08)
var POSITION_UPDATE = byte(0x09)
var DISCONNECT = byte(0x0e)
