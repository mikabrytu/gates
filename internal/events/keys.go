package events

import "github.com/Papiermond/eventbus"

const GAME_OVER_EVENT eventbus.EventType = "game:over"
const GAME_RESTART_EVENT eventbus.EventType = "game:restart"
const PLAYER_ATTACK_EVENT eventbus.EventType = "player:attack"
const PLAYER_BREAK_SPELL_EVENT = "player:break-spell"
const ENEMY_ATTACK_EVENT eventbus.EventType = "enemy:attack"
const ENEMY_DEAD_EVENT eventbus.EventType = "enemy:dead"
const ENEMY_BREAK_PARALYSIS_EVENT = "enemy:break-paralysis"
