async function login() {
    try {
        let obj = {
            nome: document.getElementById("username").value,
            password: document.getElementById("password").value
        }
        let user = await $.ajax({
            url: '/api/users/login',
            method: 'post',
            dataType: 'json',
            data: JSON.stringify(obj),
            contentType: 'application/json'
        });
        window.location.href="files.html"
    } catch (err) {
        document.getElementById("msg").innerHTML = html;
    }
}

async function addUser() {
    try {
        let users = {        
            nome: document.getElementById("new-username").value,
            email: document.getElementById("email").value,
            password: document.getElementById("new-password").value,               
        } 
        console.log(users);  
        let result = await $.ajax({
            url: "/api/users/adduser",
            method: "post",
            dataType: "json",
            data: JSON.stringify(users),
            contentType: "application/json"
        });

    } catch (err) {
        document.getElementById("msg1").innerHTML = html;
    }
}