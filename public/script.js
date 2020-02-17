var signin = document.querySelector("#signin");
var signup = document.querySelector("#signup");
var chat = document.querySelector("#chat");

var ws = new WebSocket("ws://localhost:1337/ws");
ws.addEventListener("message", msg => getMsg(msg));

document.querySelector("#signup-form").addEventListener("submit", e => signUp(e));
document.querySelector("#signin-form").addEventListener("submit", e => signIn(e));
document.querySelector("#chat-form").addEventListener("submit", e => sendMsg(e));
document.querySelector("#signin-link").addEventListener("click", e => toggleSign());
document.querySelector("#signup-link").addEventListener("click", e => toggleSign());

function displaySnackbar(body, color) {
    var snackbar = document.querySelector("#snackbar");
    snackbar.innerHTML = body;
    snackbar.style.backgroundColor = color || "#000";
    snackbar.classList.toggle("show");
    setTimeout(() => snackbar.classList.toggle("show"), 3000);
}

function toggleSign() {
    signup.classList.toggle("hidden");
    signin.classList.toggle("hidden");
}

function sendMsg(e) {
    e.preventDefault();
    var input = document.querySelector("#chat-input");
    ws.send(JSON.stringify({ message: input.value }));
    input.value = "";
    input.focus();
}

function getMsg(msg) {
    var response = JSON.parse(msg.data);
    var messageList = document.querySelector(".messages");
    var div = document.createElement("div");
    div.textContent = response.message;
    messageList.appendChild(div);
}

function signUp(e) {
    e.preventDefault();
    let username = document.querySelector("#signup-user").value;
    let email = document.querySelector("#signup-email").value;
    let password = document.querySelector("#signup-password").value;

    fetch("http://localhost:1337/register", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            username,
            email,
            password
        })
    })
        .then(async res => {
            if (res.status === 200) {
                return res.text();
            } else {
                let err = await res.text();
                throw new Error(err);
            }
        })
        .then(json => {
            signup.classList.toggle("hidden");
            signin.classList.toggle("hidden");
            displaySnackbar(json, "green");
        })
        .catch(err => displaySnackbar(err, "red"));
}

function signIn(e) {
    e.preventDefault();
    let email = document.querySelector("#signin-email").value;
    let password = document.querySelector("#signin-password").value;

    fetch("http://localhost:1337/login", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            email,
            password
        })
    })
        .then(async res => {
            if (res.status === 200) {
                return res.text();
            } else {
                let err = await res.text();
                throw new Error(err);
            }
        })
        .then(json => {
            signin.classList.toggle("hidden");
            chat.classList.toggle("hidden");
            displaySnackbar(json, "green");
            getTicket();
        })
        .catch(err => displaySnackbar(err, "red"));
}

function getTicket() {
    fetch("http://localhost:1337/wsTicket")
        .then(async res => {
            if (res.status === 200) {
                return res.text();
            } else {
                let err = await res.text();
                throw new Error(err);
            }
        })
        .then(message => {
            ws.send(JSON.stringify({ message }));
        })
        .catch(err => console.log(err));
}
