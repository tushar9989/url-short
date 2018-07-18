var signedIn = false;
function signIn () {
    var auth2 = gapi.auth2.getAuthInstance();
    if (signedIn) {
        auth2.signOut()
            .then(() => {
                signedIn = false;
                adjustButtonText();
                document.cookie = 'X-Auth-Token=;expires=Thu, 01 Jan 1970 00:00:01 GMT;';
            });
    } else {
        auth2.signIn()
            .then((googleUser) => {
                signedIn = true;
                adjustButtonText();
                document.cookie = "X-Auth-Token=" + googleUser.getAuthResponse().id_token + ";max-age=300;path=/";
            }).catch(function (error) {
                console.log(error);
                signedIn = false;
                adjustButtonText();
            });
    }
}

function adjustButtonText() {
    var text;
    if(signedIn) {
        text = "Sign Out";
    } else {
        text = "Sign In";
    }
    document.getElementById("button-text").textContent = text;
}

gapi.load('auth2', function () {
    gapi.auth2.init().then(function (auth2) {
        if (auth2.isSignedIn.get()) {
            var googleUser = auth2.currentUser.get();
            document.cookie = "X-Auth-Token=" + googleUser.getAuthResponse().id_token + ";max-age=300;path=/";
            signedIn = true;
        } else {
            signedIn = false;
        }
        adjustButtonText();
    });
});