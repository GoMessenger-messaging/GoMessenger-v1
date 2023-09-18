let name = "";
let channel = "";
let state = 0; //0=Name, 1=Channel, 2=Operational
let timer = 5;
let updateRunning = false;
const help = [
    "<p> - You may change your username at any time<br> - You may choose a username that is still available on the server<br> - Still need help? Write an email to <a href=mailto:contact@gomessenger.link>contact@gomessenger.link</a></p>",
    "<p> - You may join any channel that exists on the server<br> - You may create a new channel by entering a name that is not taken<br> - Still need help? Write an email to <a href=mailto:contact@gomessenger.link>contact@gomessenger.link</a></p>",
    "<p> - ':cu' to get the current username<br> - ':cc' to get the current channel<br> - ':u' to change your username<br> - ':c' to change the channel<br> - ':r' to manually check for new messages<br> - Still need help? Write an email to <a href=mailto:contact@gomessenger.link>contact@gomessenger.link</a></p>"
];

function ping(){
    let ping = new XMLHttpRequest();
    ping.open("GET", "../ping?name=" + name);
    ping.send();

    setTimeout(null, 300000);
}

function getMessages(){
    let get = new XMLHttpRequest();
    get.onreadystatechange = function() {
        if (get.readyState === 4) {
            if (get.status !== 200) {
                document.getElementById("messages").innerHTML = "HTTP request error! Please try again<br>Error code: " + get.status;
            } else {
                document.getElementById("messages").innerHTML = get.responseText.replace(/(\r\n|\r|\n)/g, '<br>');
            }
        }
    }
    get.open("GET", "../get_messages?channel=" + channel);
    get.send();
}

function send() {
    let input = document.getElementById("input").value;
    document.getElementById("input").value = "";
    //Username
    if (state === 0) {
        if (input[0] === ":") {
            if (input === ":h") {
                document.getElementById("messages").innerHTML = help[0];
            }
        } else {
            name = input;
            let check = new XMLHttpRequest();
            check.onreadystatechange = function () {
                if (check.readyState === 4) {
                    if (check.status !== 200) {
                        document.getElementById("messages").innerHTML = "HTTP request error! Please try again<br>Error code: " + check.status;
                    } else if (check.responseText !== "Registered successfully") {
                        document.getElementById("messages").innerHTML = check.responseText;
                    } else {
                        document.getElementById("messages").innerHTML = "Input channel";
                        state = 1;
                    }
                }
            }
            check.open("GET", "../register?name=" + name);
            check.send();
        }
    }
    //Channel
    else if (state === 1) {
        if (input[0] === ":") {
            if (input === ":h") {
                document.getElementById("messages").innerHTML = help[1];
            } else {
                document.getElementById("messages").innerHTML = "Not a valid command";
            }
        } else {
            channel = input;
            state = 2;
            if (!updateRunning) {
                updateRunning = true;
                update();
            }
            else {
                getMessages();
            }
        }
    }
    //Messages
    else if (state === 2) {
        if (input[0] === ":") {
            if (input === ":h") {
                document.getElementById("messages").innerHTML = help[2];
                timer = 0;
            } else if (input === ":cu") {
                document.getElementById("messages").innerHTML = name;
                timer = 0;
            } else if (input === ":cc") {
                document.getElementById("messages").innerHTML = channel;
                timer = 0;
            } else if (input === ":u") {
                document.getElementById("messages").innerHTML = "Input username";
                state = 0;
            } else if (input === ":c") {
                document.getElementById("messages").innerHTML = "Input channel";
                state = 1;
            } else if (input === ":r") {
                getMessages();
            } else {
                document.getElementById("messages").innerHTML = "Not a valid command";
                timer = 0;
            }
        } else {
            let send = new XMLHttpRequest();
            send.onreadystatechange = function () {
                if (send.readyState === 4) {
                    if (send.status !== 200) {
                        document.getElementById("messages").innerHTML = "HTTP request error! Please try again<br>Error code: " + send.status;
                        timer = 0;
                    }
                    else {
                        getMessages();
                        timer = 0;
                    }
                }
            }
            send.open("GET", "../send?name=" + name + "&channel=" + channel + "&message=" + input);
            send.send();
        }
    }
}

function update() {
    if (state === 2) {
        if (timer === 5) {
            getMessages();
            ping();
            timer = 0;
        } else {
            timer++;
        }
    }
    setTimeout(update, 1000);
}
