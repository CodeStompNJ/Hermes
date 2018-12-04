new Vue({
    el: '#app',

    data: {
        ws: null, // Our websocket
        newMsg: '', // Holds new messages to be sent to the server
        chatContent: '', // A running list of chat messages displayed on the screen
        email: null, // Email address used for grabbing an avatar
        username: null, // Our username
        group: null, //group we want to join
        joined: false // True if email and username have been filled in
    },

created: function() {

        var self = this;
        this.ws = new WebSocket('ws://' + window.location.host + '/ws');
        this.ws.addEventListener('message', self.addMessageWS);
        this.getAndAddChatroomHistory();
    },

methods: {
        addMessageWS: function(e) {
            var msg = JSON.parse(e.data);
            this.chatContent += '<div class="chip">'
                    + '<img src="' + this.gravatarURL(msg.email) + '">' // Avatar
                    + msg.username
                + '</div>'
                + emojione.toImage(msg.message) + '<br/>'; // Parse emojis

            var element = document.getElementById('chat-messages');
            element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
        },

        addMessages: function (messages) {
            messages.forEach(message => {
                this.chatContent += '<div class="chip">'
                    + '<img src="http://www.gravatar.com/avatar/7b064dad507c266a161ffc73c53dcdc5">' // Avatar
                    + message.UserID
                + '</div>'
                + emojione.toImage(message.Text) + '<br/>'; // Parse emojis
            });

            var element = document.getElementById('chat-messages');
            element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
        },

        send: function () {
            if (this.newMsg != '') {
                this.ws.send(
                    JSON.stringify({
                        email: this.email,
                        username: this.username,
                        group: this.group,
                        message: $('<p>').html(this.newMsg).text() // Strip out html
                    }
                ));
                this.newMsg = ''; // Reset newMsg
            }
        },

        getAndAddChatroomHistory: function () {
            axios.get('/history')
                .then((response) => {
                    console.log(response);
                    this.addMessages(response.data);
                }).catch((error) => {
                    console.log("[ERRRRRRR] Something happened getting history!");
                    console.log(error);
                });
        },

        history: function() {
            axios.get('/history')
                .then((response) => {
                    debugger;
                })
                .catch((error) => {
                    debugger;
                });
        },

join: function () {
            if (!this.email) {
                Materialize.toast('You must enter an email', 2000);
                return
            }
            if (!this.username) {
                Materialize.toast('You must choose a username', 2000);
                return
            }
            this.email = $('<p>').html(this.email).text();
            this.username = $('<p>').html(this.username).text();
            this.group = $('<p>').html(this.group).text();
            if (this.joined = true);
        },
gravatarURL: function(email) {
            return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email);
        }
    }
});
