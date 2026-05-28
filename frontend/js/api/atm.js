import { protected_fetch } from "./util.js";

export const atm_transfer = tx => protected_fetch(`api/atm/transfer`, {
    method: "POST",
    body: JSON.stringify(tx)
});

export const get_balance = playername => protected_fetch(`api/atm/balance/${playername}`);
