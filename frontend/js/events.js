import EventEmitter from './util/eventemitter.js';

// on startup
export const EVENT_STARTUP = "startup";
// on login (or on startup with existing session/credentials)
export const EVENT_LOGGED_IN = "logged_in";

export default new EventEmitter();