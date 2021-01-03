function logout() {
    var auth2 = gapi.auth2.getAuthInstance();
    auth2.signOut().then(function () {
        console.log('User signed out.');
        fetch("/logout").then(res => {
            window.location.replace("/");
        })
    });
}

window.addEventListener('load', function () {
    document.getElementsByClassName("container")[0].style.display = 'block';
});