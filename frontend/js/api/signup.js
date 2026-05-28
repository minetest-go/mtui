export const signup = data => fetch("api/signup", {
    method: "POST",
    body: JSON.stringify(data)
}).then(r => r.text());


export const signup_captcha = () => fetch("api/signup/captcha").then(r => r.text());