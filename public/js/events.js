import EventEmitter from './util/eventemitter.js';

// on startup
export const EVENT_STARTUP = "startup";
// on login (or on startup with existing session/credentials)
export const EVENT_LOGGED_IN = "logged_in";

export const EVENT_STATS = "stats";
export const EVENT_PLAYER_STATS = "player_stats";
export const EVENT_PLAYER_STATS_EXTRA = "player_stats_extra";

export default new EventEmitter();