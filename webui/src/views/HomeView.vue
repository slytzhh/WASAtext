<script>
    export default {
        data: function() {
            return {
                errormsg: null,
                userid: localStorage.getItem("userid"),
                username: localStorage.getItem("username"),
                userphoto: localStorage.getItem("userphoto"),
                
                // for searching users
                searcheduser: null,
                users: [],

                emojiList: ['😀', '😂', '😍', '👍', '👏', '🙌', '🔥', '❤️', '🎉', '🤔', '😢', '👎'],

                // variable to make appear the blurred boxes
                boxshown: 0,
                /* 
                    2 = change group name and photo
                    3 = create group name and photo
                    4 = create group select users
                    5 = add users to group
                    6 = list of users in group
                */

                // for changing user info
                newusername: null,
                newuserphoto: null,
                changedinfo: false,

                // for chats preview
                chats: [],

                // for main chat
                mainchat: null,
                chatshown: false,
                messagetext: null,
                messagephoto: null,

                // for comments
                commentshown: 0,
                commentemoji: null,

                // for forwarding message
                messageToforward: 0,

                // for creating new group
                createdgroupname: null,
                createdgroupphoto: null,
                userscreategroup: new Set(),
                allusers: [],

                // for adding users to a group
                usersnotinchat: [],
                userstoadd: new Set(),

                // for group members
                groupMembers: [],

                // for replying to message
                replyid: null,
                replyuser: null,
                replytext: null,
                replyphoto: null,

                // to refresh every n seconds
                intervalid: null
            }
        },
        methods: {
            handleClickOutside(event) {
                // Check if the click is outside the search box to search users
                if (this.$refs.boxsearchuser && !this.$refs.boxsearchuser.contains(event.target)) {
                    this.users = [];
                    this.searcheduser = null;
                }
                // Check if the click is outside the sidebar when forwarding a message
                if (this.messageToforward!=0 && event.target.id != "forwardbutton" && !this.$refs?.chatlist.contains(event.target)){
                    this.messageToforward = 0;
                    this.refresh();
                    this.intervalid=setInterval(this.refresh,5000);
                }
            },
            async searchUser() {
                this.errormsg = null;
                this.users = []; // Resetta l'array prima di una nuova ricerca
                this.usersnotinchat = [];
                
                if (!this.searcheduser || this.searcheduser.trim() === '') {
                    return; // Esci se la stringa di ricerca è vuota
                }

                try {
                    let response = await this.$axios.get("/users", { 
                        params: { username: this.searcheduser } 
                    });
                    
                    // Aggiorna sia users che usersnotinchat
                    response.data.userlist.forEach(user => {
                        if (user.username != this.username) {
                            this.users.push(user); // Aggiungi all'array principale dei risultati
                            if (!this.mainchat?.usernamelist?.includes(user.username)) {
                                this.usersnotinchat.push(user);
                            }
                        }
                    });
                } catch (e) {
                    this.errormsg = e.response?.status + ": " + e.response?.data || "Errore nella ricerca";
                }
            },
            goToProfileView(){
                this.$router.push({ path: '/profile' });
            },
            async buildChatPreview(){
                this.errormsg = null;
                this.chats=[];
                try {
                    let response = await this.$axios.get("/chats",{headers:{"Authorization": `Bearer ${this.userid}`}});
                    response.data.forEach(chat => {
                        if (chat.groupname.length>16){
                            chat.groupname = chat.groupname.slice(0,16)+"...";
                        }
                        if (chat.lastmessage.text.length>18){
                            chat.lastmessage.text = chat.lastmessage.text.slice(0,18)+"...";
                        }
                        if (chat.lastmessage.photo.length>0 && chat.lastmessage.text.length===0){
                            chat.lastmessage.text="Photo";
                        }
                        if (chat.lastmessage.username==this.username){
                            chat.lastmessage.username="You";
                        }
                        chat.lastmessage.timestamp = chat.lastmessage.timestamp.slice(11,16);
                        chat.timecreated = chat.timecreated.slice(11,16);
                        this.chats.push(chat);
                    });
                } catch (e) {
                    this.errormsg = e.response.status + ": " + e.response.data;
                }
            },
            async buildMainChat(chatid){
                this.errormsg = null;
                if(this.messageToforward != 0){
                    if(chatid>=0){
                        this.forwardMessage(chatid);
                    }else{
                        await this.getAllUsers();
                        var name;
                        this.allusers.forEach(user =>{
                            if(user.userid==-chatid){
                                name=user.username;
                            }
                        });
                        try {
                            let response = await this.$axios.post("/newchat",{usernamelist:[this.username,name], forwardedid: this.messageToforward},{headers:{"Authorization": `Bearer ${this.userid}`}});
                            this.messageToforward = 0;
                            this.buildMainChat(response.data.chatid);
                        } catch (e) {
                            this.errormsg = e.response.status + ": " + e.response.data;
                        }
                    }
                    this.intervalid=setInterval(this.refresh,5000);
                }else{
                    try {
                        let response = await this.$axios.get("/chats/"+chatid,{headers:{"Authorization": `Bearer ${this.userid}`}});
                        this.mainchat=response.data;
                        if (this.mainchat.groupname.length>16){
                            this.mainchat.groupname = this.mainchat.groupname.slice(0,16)+"...";
                        }
                        this.mainchat.messagelist.forEach( message => {
                            message.timestamp = message.timestamp.slice(11,16);
                        });
                        this.chatshown = true;
                    } catch (e) {
                        this.errormsg = e.response.status + ": " + e.response.data;
                    }
                    this.buildChatPreview();
                }
            },
            async closeMainChat(){
                this.errormsg = null;
                this.mainchat = null;
                this.messagephoto = null;
                this.messagetext = null;
                this.commentshown = 0;
                this.commentemoji = null;
                this.chatshown = false;
                this.replyid = null;
                this.replyphoto = null;
                this.replytext = null;
                this.replyuser = null;
            },
            async getGroupMembers(chat) {
                this.errormsg = null;
                try {
                    let response = await this.$axios.get("/chats/"+chat.chatid+"/users", { headers: { Authorization: `Bearer ${this.userid}` } });
                    this.groupMembers = response.data;
                } catch (e) {
                    this.errormsg = e.response.status + ": " + e.response.data;
                }
            },
            // button to send a photo handler
            sendPhotoFileSelect(){
                const file = this.$refs.sendPhotoInput.files[0];
                if (file) {
                    const reader = new FileReader();
                    reader.onload = (e) => {
                        this.messagephoto = e.target.result;
                    };
                    reader.readAsDataURL(file);
                }
            },
            sendPhotoButton(){
                this.$refs.sendPhotoInput.click();
            },
            async sendMessage(){
                if(this.replyid){
                    try{
                        let response = await this.$axios.post("/chats/"+this.mainchat.chatid+"/repliedmessages",{text: this.messagetext,photo: this.messagephoto, replyid:this.replyid},{headers:{"Authorization": `Bearer ${this.userid}`}});
                        this.messagetext = null;
                        this.messagephoto = null;
                        this.resetReplyMessage();
                        this.buildMainChat(this.mainchat.chatid);
                    }catch (e) {
                        this.errormsg = e.response.status + ": " + e.response.data;
                    }
                }else{
                    try {
                        let response = await this.$axios.post("/chats/"+this.mainchat.chatid+"/messages",{text: this.messagetext,photo: this.messagephoto},{headers:{"Authorization": `Bearer ${this.userid}`}});
                        this.messagetext = null;
                        this.messagephoto = null;
                        this.buildMainChat(this.mainchat.chatid);
                    } catch (e) {
                        this.errormsg = e.response.status + ": " + e.response.data;
                    }
                }
            },

            // function called when sending a message that check if the chat is temp
            async sendMessageorCreateChat(){
                if(this.messagetext==null && this.messagephoto==null){
                    return
                }
                if(this.mainchat.chatid==-1){
                    try {
                        let response = await this.$axios.post("/newchat",{usernamelist:[this.username,this.mainchat.groupname],firstmessage:{text:this.messagetext, photo:this.messagephoto}},{headers:{"Authorization": `Bearer ${this.userid}`}});
                        this.messagetext = null;
                        this.messagephoto = null;
                        this.buildMainChat(response.data.chatid);
                    } catch (e) {
                        this.errormsg = e.response.status + ": " + e.response.data;
                    }
                }else{
                    this.sendMessage();
                }
            },
            async forwardMessage(chatid){
                try {
                    let response = await this.$axios.post("/chats/"+chatid+"/forwardedmessages",{messageid: this.messageToforward},{headers:{"Authorization": `Bearer ${this.userid}`}});
                    this.messageToforward = 0;
                    this.buildMainChat(chatid);
                } catch (e) {
                    this.errormsg = e.response.status + ": " + e.response.data;
                }
            },
            async openChatFromUser(user){
                await this.buildChatPreview();
                var chatid = -1;
                for(let i=0;i<this.chats.length;i++){
                    if(this.chats[i].groupname==user.username && !this.chats[i].isgroup){
                        chatid = this.chats[i].chatid;
                        break;
                    }
                }
                if(chatid==-1){
                    this.mainchat={
                        chatid:-1,
                        groupname:user.username,
                        groupphoto:user.photo,
                        isgroup: false,
                        messagelist:[]
                    }
                    this.chatshown = true;
                }else{
                    this.buildMainChat(chatid);
                }
                this.searcheduser = null;
                this.users = [];
            },
            async deleteMessage(message){
                try{
                    let response = await this.$axios.delete("/chats/"+this.mainchat.chatid+"/messages/"+message.messageid,{headers:{"Authorization": `Bearer ${this.userid}`}});
                    if(this.mainchat.messagelist.length>1 || this.mainchat.isgroup){
                        this.buildMainChat(message.chatid);
                    }else{
                        this.mainchat = null;
                        this.buildChatPreview();
                    }
                } catch (e) {
                    this.errormsg = e.response.status + ": " + e.response.data;
                }
            },
            async showComments(message){
                if(this.commentshown!=message.messageid){
                    this.commentshown=message.messageid;
                }else{
                    this.commentshown=0;
                }
                this.errormsg = null;
            },
            async commentMessage(message){
                this.errormsg = null;
                const emojiRegex = /[\u{1F600}-\u{1F64F}|\u{1F300}-\u{1F5FF}|\u{1F680}-\u{1F6FF}|\u{1F700}-\u{1F77F}|\u{1F780}-\u{1F7FF}|\u{1F800}-\u{1F8FF}|\u{1F900}-\u{1F9FF}|\u{1FA00}-\u{1FA6F}|\u{2600}-\u{26FF}|\u{2700}-\u{27BF}|\u{FE00}-\u{FE0F}]/gu;
                if(emojiRegex.test(this.commentemoji)){
                    try{
                        let response = await this.$axios.put("/chats/"+this.mainchat.chatid+"/messages/"+message.messageid+"/comments",{emoji: this.commentemoji},{headers:{"Authorization": `Bearer ${this.userid}`}});
                        this.buildMainChat(message.chatid);
                    } catch (e) {
                        this.errormsg = e.response.status + ": " + e.response.data;
                    }
                }else{
                    this.errormsg = "the comment must consist of a single emoji";
                }
                this.commentemoji = null;
            },
            async selectEmoji(emoji, message) {
                this.commentemoji = emoji;
                this.commentMessage(message);
            },
            async deleteComment(message){
                try{
                    let response = await this.$axios.delete("/chats/"+this.mainchat.chatid+"/messages/"+message.messageid+"/comments",{headers:{"Authorization": `Bearer ${this.userid}`}});
                    this.buildMainChat(message.chatid);
                } catch (e) {
                    this.errormsg = e.response.status + ": " + e.response.data;
                }
            },
            async startForwardingMessage(message) {
                clearInterval(this.intervalid);
                this.intervalid = null;
                this.messageToforward = message.messageid;
                await this.getAllUsers();
                this.allusers.forEach(user => {
                    var found = false;
                    for (var i = 0; i < this.chats.length; i++) {
                        if (user.username == this.chats[i].groupname && !this.chats[i].isgroup) {
                            found = true;
                            break;
                        }
                    }
                    if (!found) {
                        this.chats.push({ chatid: -user.userid, groupname: user.username, groupphoto: user.photo });
                    }
                });
            },
            // Generic file input handler
            handleFileSelect(ref, property) {
                const file = this.$refs[ref].files[0];
                if (file) {
                    const reader = new FileReader();
                    reader.onload = (e) => {
                        this[property] = e.target.result;
                    };
                    reader.readAsDataURL(file);
                }
            },

            // Button to trigger file input for group photo
            triggerFileInput(ref) {
                this.$refs[ref].click();
            },
            // Main method to change group name and photo
            async changeGroupNamePhoto() {
                this.errormsg = null;
                let hasError = false;

                if (this.newgroupname) {
                    try {
                        await this.$axios.put(
                            `/chats/${this.mainchat.chatid}/name`,
                            { groupname: this.newgroupname.trim() },
                            { headers: { Authorization: `Bearer ${this.userid}` } }
                        );
                        this.changedgroupinfo = true;
                        // Aggiorna il nome del gruppo nella chat principale
                        this.mainchat.groupname = this.newgroupname.trim();
                    } catch (e) {
                        hasError = true;
                        this.errormsg = e.response
                            ? `${e.response.status}: ${e.response.data}`
                            : "Failed to update group name. Please check your network connection.";
                    }
                }

                if (this.newgroupphoto) {
                    try {
                        await this.$axios.put(
                            `/chats/${this.mainchat.chatid}/photo`,
                            { photo: this.newgroupphoto },
                            { headers: { Authorization: `Bearer ${this.userid}` } }
                        );
                        this.changedgroupinfo = true;
                        // Aggiorna la foto del gruppo nella chat principale
                        this.mainchat.groupphoto = this.newgroupphoto;
                    } catch (e) {
                        hasError = true;
                        this.errormsg = e.response
                            ? `${e.response.status}: ${e.response.data}`
                            : "Failed to update group photo. Please check your network connection.";
                    }
                }

                this.newgroupname = null;
                this.newgroupphoto = null;

                if (this.changedgroupinfo && !hasError) {
                    this.boxshown = 0;
                    this.errormsg = null;
                    this.buildMainChat(this.mainchat.chatid); // Ricarica i dati della chat principale
                }
            },
            createGroupPhotoFileSelect() {
                const file = this.$refs.createGroupPhotoInput.files[0];
                if (file) {
                    const reader = new FileReader();
                    reader.onload = (e) => {
                        this.createdgroupphoto = e.target.result;
                        this.newgroupphoto = e.target.result;
                    };
                    reader.readAsDataURL(file);
                }
            },
            createGroupPhotoButton() {
                this.$refs.createGroupPhotoInput.click();
            },
            async resetChangeGroupPrompt(){
                this.newgroupname = null;
                this.newgroupphoto = null;
                this.boxshown = 0;
                this.errormsg = null;
            },
            async leaveGroup(){
                try{
                    let response = await this.$axios.delete("/chats/"+this.mainchat.chatid+"/users/"+this.userid,{headers:{"Authorization": `Bearer ${this.userid}`}});
                    this.mainchat = null;
                    this.buildChatPreview();
                } catch (e) {
                    this.errormsg = e.response.status + ": " + e.response.data;
                }
            },
            async resetCreateGroupPrompt(){
                this.userscreategroup = new Set();
                this.createdgroupname = null;
                this.createdgroupphoto = null;
                this.boxshown = 0;
                this.errormsg = null;
            },
            async createGroup(){
                // default photo
	            const default_photo = "data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0idXRmLTgiPz4NCg0KPCFET0NUWVBFIHN2ZyBQVUJMSUMgIi0vL1czQy8vRFREIFNWRyAxLjEvL0VOIiAiaHR0cDovL3d3dy53My5vcmcvR3JhcGhpY3MvU1ZHLzEuMS9EVEQvc3ZnMTEuZHRkIj4NCjwhLS0gVXBsb2FkZWQgdG86IFNWRyBSZXBvLCB3d3cuc3ZncmVwby5jb20sIEdlbmVyYXRvcjogU1ZHIFJlcG8gTWl4ZXIgVG9vbHMgLS0+CjxzdmcgZmlsbD0iIzAwMDAwMCIgdmVyc2lvbj0iMS4xIiBpZD0iTGF5ZXJfMSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiB4bWxuczp4bGluaz0iaHR0cDovL3d3dy53My5vcmcvMTk5OS94bGluayIgDQoJIHdpZHRoPSI4MDBweCIgaGVpZ2h0PSI4MDBweCIgdmlld0JveD0iNzk2IDc5NiAyMDAgMjAwIiBlbmFibGUtYmFja2dyb3VuZD0ibmV3IDc5NiA3OTYgMjAwIDIwMCIgeG1sOnNwYWNlPSJwcmVzZXJ2ZSI+DQo8cGF0aCBkPSJNODk2LDc5NmMtNTUuMTQsMC05OS45OTksNDQuODYtOTkuOTk5LDEwMGMwLDU1LjE0MSw0NC44NTksMTAwLDk5Ljk5OSwxMDBjNTUuMTQxLDAsOTkuOTk5LTQ0Ljg1OSw5OS45OTktMTAwDQoJQzk5NS45OTksODQwLjg2LDk1MS4xNDEsNzk2LDg5Niw3OTZ6IE04OTYuNjM5LDgyNy40MjVjMjAuNTM4LDAsMzcuMTg5LDE5LjY2LDM3LjE4OSw0My45MjFjMCwyNC4yNTctMTYuNjUxLDQzLjkyNC0zNy4xODksNDMuOTI0DQoJcy0zNy4xODctMTkuNjY3LTM3LjE4Ny00My45MjRDODU5LjQ1Miw4NDcuMDg1LDg3Ni4xMDEsODI3LjQyNSw4OTYuNjM5LDgyNy40MjV6IE04OTYsOTgzLjg2DQoJYy0yNC42OTIsMC00Ny4wMzgtMTAuMjM5LTYzLjAxNi0yNi42OTVjLTIuMjY2LTIuMzM1LTIuOTg0LTUuNzc1LTEuODQtOC44MmM1LjQ3LTE0LjU1NiwxNS43MTgtMjYuNzYyLDI4LjgxNy0zNC43NjENCgljMi44MjgtMS43MjgsNi40NDktMS4zOTMsOC45MSwwLjgyOGM3LjcwNiw2Ljk1OCwxNy4zMTYsMTEuMTE0LDI3Ljc2NywxMS4xMTRjMTAuMjQ5LDAsMTkuNjktNC4wMDEsMjcuMzE4LTEwLjcxOQ0KCWMyLjQ4OC0yLjE5MSw2LjEyOC0yLjQ3OSw4LjkzMi0wLjcxMWMxMi42OTcsOC4wMDQsMjIuNjE4LDIwLjAwNSwyNy45NjcsMzQuMjUzYzEuMTQ0LDMuMDQ3LDAuNDI1LDYuNDgyLTEuODQyLDguODE3DQoJQzk0My4wMzcsOTczLjYyMSw5MjAuNjkxLDk4My44Niw4OTYsOTgzLjg2eiIvPg0KPC9zdmc+"
                
                if(!this.createdgroupphoto){
                    this.createdgroupphoto = default_photo;
                }
                try {
                    const userslist = [this.username];
                    this.userscreategroup.forEach( user =>{
                        userslist.push(user);
                    });
                    let response = await this.$axios.post("/newchat",{usernamelist:userslist, groupname: this.createdgroupname, groupphoto: this.createdgroupphoto},{headers:{"Authorization": `Bearer ${this.userid}`}});
                    this.createdgroupname = null;
                    this.createdgroupphoto = null;
                    this.buildMainChat(response.data.chatid);
                } catch (e) {
                    this.errormsg = e.response.status + ": " + e.response.data;
                }
                this.chatshown = true;
                this.boxshown = 0;
                this.userscreategroup = new Set();
            },
            async getAllUsers() {
                this.errormsg = null;
                this.allusers=[];
                try {
                    let response = await this.$axios.get("/users", {params: {username: ""}});
                    response.data.userlist.forEach(user => {
                        if (user.username != this.username){
                            this.allusers.push(user);
                        }
                    });
                } catch (e) {
                    this.errormsg = e.response.status + ": " + e.response.data;;
                }
            },

            // functions to handle the list of users when creating a group
            isSelectedCreation(user){
                return this.userscreategroup.has(user.username);
            },
            insertUserCreation(user){
                if (this.userscreategroup.has(user.username)) {
                    this.userscreategroup.delete(user.username);
                } else {
                    this.userscreategroup.add(user.username);
                }
            },
            // functions to add users to a group
            async resetAddUsersPrompt(){
                this.userstoadd = new Set();
                this.boxshown = 0;
                this.errormsg = null;
            },
            isSelectedAdding(user){
                return this.userstoadd.has(user.username);
            },
            insertUserAdding(user){
                if (this.userstoadd.has(user.username)) {
                    this.userstoadd.delete(user.username);
                } else {
                    this.userstoadd.add(user.username);
                }
            },
            async addToGroup(){
                try {
                    const userslist = [];
                    this.userstoadd.forEach( user =>{
                        userslist.push(user);
                    });
                    let response = await this.$axios.put("/chats/"+this.mainchat.chatid+"/users",{usernamelist:userslist},{headers:{"Authorization": `Bearer ${this.userid}`}});
                    this.buildMainChat(this.mainchat.chatid);
                } catch (e) {
                    this.errormsg = e.response.status + ": " + e.response.data;
                }
                this.boxshown = 0;
                this.userstoadd = new Set();
            },
            async getUsersNotInChat() {
                this.errormsg = null;
                this.usersnotinchat=[];
                try {
                    let response = await this.$axios.get("/users", {params: {username: ""}});
                    response.data.userlist.forEach(user => {
                        if (!this.mainchat.usernamelist.includes(user.username)){
                            this.usersnotinchat.push(user);
                        }
                    });
                } catch (e) {
                    this.errormsg = e.response.status + ": " + e.response.data;;
                }
            },
            async setReplyMessage(message){
                this.replyid=message.messageid;
                this.replyuser=message.username;
                this.replytext=message.text;
                this.replyphoto=message.photo;
            },
            async resetReplyMessage(){
                this.replyid=null;
                this.replyuser=null;
                this.replytext=null;
                this.replyphoto=null;
            },
            async logout(){
                this.$router.push({path: "/"});
            },
            // function to refresh the views
            async refresh(){
                this.messageToforward = 0;
                if(this.mainchat){
                    this.buildMainChat(this.mainchat.chatid);
                }else{
                    this.buildChatPreview();
                }
            }
        },
        mounted(){
            this.refresh();
            document.addEventListener('click', this.handleClickOutside);
            this.intervalid=setInterval(this.refresh,5000);
        },
        beforeUnmount(){
            clearInterval(this.intervalid);
            this.intervalid = null;
        },
        /* Updater to check if the messagelist is shown and make the list start from bottom */
        updated(){
            const div = document.querySelector('#messagelist');
            if (div) {
                div.scrollTop=div.scrollHeight;
            }
        }
    }
</script>


<template>
    <div class="all-screen">
        <div class="navbar-dark">
            <!-- User photo, name and button to change those -->
            <div class="user-info">
                <img class="img-circular" :src="userphoto" style="width: 32px; height: 32px; margin-left: 2px;"/>
                <h3 style="margin-left: 10px; margin-bottom: 0; margin-right: 10px;">{{username}}</h3>
                <img @click="goToProfileView" src="/assets/pencil.svg" style="width: 16px; height: 16px; cursor: pointer; margin-right: 10px;" v-if="boxshown != 1"/>
            </div>

            <!-- Searchbox to search users -->
            <div class="searchbox" ref="boxsearchuser">
                <img src="/assets/search.svg" style="width: 32px; height: 32px; margin-top: 2px; margin-left: 2px;">
                <div class="searchbox-userlist">
                    <input class="searchbox-user" v-model="searcheduser" 
                        placeholder="Search user" @input="searchUser">
                    <div v-if="users.length > 0" class="searched-dropdown">
                        <ul>
                            <li v-for="user in users" :key="user.username" 
                                @click="openChatFromUser(user)">
                                <img :src="user.photo" class="img-circular" style="width: 24px; height: 24px;">
                                {{ user.username }}
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
            <!-- Logout button -->
            <button class="leavebutton" @click="logout">
                <img src="/assets/leave.svg" style="height: 32px; width: 32px;">
                Log Out
            </button>
        </div>


        <div class="main-screen">
            <!-- Sidebar with chats -->
            <div class="sidebar-chats" ref="sidebar">
                <div class="sidebar-buttons">
                    <img src="/assets/add-square.svg" style="width: 32px; height: 32px; margin-left: 10px; margin-right: 10px; cursor: pointer;" @click="boxshown = 3">
                    <div v-if="messageToforward!=0" style="color: whitesmoke;">
                        Forward to...<br>
                        (click ouside sidebar to cancel)
                    </div>
                    <img src="/assets/refresh.svg" style="width: 32px; height: 32px; margin-left: 10px; margin-right: 10px; cursor: pointer;" @click="refresh">
                </div>
                <div v-if="chats.length>0" class="chats-dropdown" ref="chatlist">
                    <ul>
                        <li v-for="chat in chats" :key="chat.chatid" @click="buildMainChat(chat.chatid)">
                            <div class="chatpreview">
                                <div class="chatpreviewname">
                                    <img class="img-circular" :src="chat.groupphoto" style="width: 32px; height: 32px;">
                                    <h4 style="margin-left: 10px; margin-bottom: 0;">{{chat.groupname}}</h4>
                                    <div v-if="chat.chatid<0"></div>
                                    <div v-else-if="chat.lastmessage.messageid!=0" class="timepreview">{{chat.lastmessage.timestamp}}</div>
                                    <div v-else class="timepreview">{{chat.timecreated}}</div>
                                </div>
                                <div v-if="chat.chatid<0"></div>
                                <div class="messagepreview" v-else-if="chat.lastmessage.messageid!=0">
                                    <b>{{chat.lastmessage.username}}: </b>
                                    <img v-if="chat.lastmessage.photo.length>0" src="/assets/photo-icon.svg" style="height: 24px; width: 24px; margin-left: 5px;">
                                    &nbsp;{{chat.lastmessage.text}}
                                    <img class="checkmark" v-if="chat.lastmessage.isallread && chat.lastmessage.userid==this.userid" src="/assets/double-check-blue.svg" style="height: 24px; width: 24px;">
                                    <img class="checkmark" v-else-if="chat.lastmessage.isallreceived && chat.lastmessage.userid==this.userid" src="/assets/double-check.svg" style="height: 24px; width: 24px;">
                                    <img class="checkmark" v-else-if="chat.lastmessage.userid==this.userid" src="/assets/single-check.svg" style="height: 24px; width: 24px;">
                                </div>
                                <div class="messagepreview" v-else>
                                    <b>Group created </b>
                                </div>
                            </div>
                        </li>
                    </ul>
                </div>
            </div>
            <!-- Main chat screen -->
            <div v-if="!mainchat" class="no-chat-message">
                <h2 style="color: whitesmoke; text-align: center; margin-top: 20px;">No chat selected. Please select a chat to start messaging.</h2>
            </div>
            <div class="main-chat" v-if="chatshown && mainchat">
                <!-- Topbar in mainchat -->
                <div class="topbar-chat">
                    <img class="backarrow" src="/assets/back-arrow.svg" style="width: 32px; height: 32px; cursor: pointer;" @click="closeMainChat"/>
                    <div class="user-info">
                        <img class="img-circular" :src="mainchat.groupphoto" style="width: 32px; height: 32px; margin-left: 2px;"/>
                        <h3 style="margin-left: 10px; margin-bottom: 0; margin-right: 10px;">{{mainchat.groupname}}</h3>
                        <img v-if="mainchat.isgroup && boxshown != 2" src="/assets/pencil.svg" style="width: 16px; height: 16px; cursor: pointer; margin-right: 10px;" @click="boxshown = 2"/>
                    </div>
                    <div v-if="mainchat.isgroup && mainchat.chatid != -2" @click="boxshown = 5; this.getUsersNotInChat();" class="addusers-button">
                        <img src="/assets/add-users.svg" style="height: 32px; width: 32px;">
                    </div>
                    <button v-if="mainchat.isgroup" class="leavebutton" @click="leaveGroup">
                        <img src="/assets/leave.svg" style="height: 32px; width: 32px;">
                        Leave
                    </button>
                    <!-- Group Members button --> 
                    <button v-if="mainchat && mainchat.isgroup" class="group-members" @click="boxshown = 6; this.getGroupMembers(mainchat);">
                        Group Info
                    </button>
                </div>
                <!-- Message screen in mainchat -->
                <div class="message-screen">
                    <div class="messagelist" id="messagelist">
                        <ul>
                            <li v-for="message in mainchat.messagelist" :key="message.messageid">
                                <span v-if="message.userid==this.userid" style="display:flex; flex-direction: row-reverse; width: calc(100vw - 360px); height: 100%; ">
                                    <div class="messagebox-you">
                                        <div v-if="message.isforwarded" class="forwarded-info" style="display: flex; justify-content: right;">
                                            <img src="/assets/forward.svg" style="width: 24px; height: 24px; margin-right: 5px;">
                                            Forwarded
                                        </div>
                                        <div class="messagebox-username" style="text-align: right;">
                                            <b><h3 style="margin-bottom: 0;">You</h3></b>
                                        </div>
                                        <div v-if="message.replymessage.messageid!=0" style="display: flex; flex-direction: column; margin-left: 15px; margin-right: 15px; background-color: #695d5d; width: calc(100% - 30px); align-items: end;">
                                            <div v-if="message.replymessage.username==this.username">
                                                <h5 style="margin-bottom: 0; margin-right: 10px;">You</h5>
                                            </div>
                                            <div v-else>
                                                <h5 style="margin-bottom: 0; margin-right: 10px;">{{message.replymessage.username}}</h5>
                                            </div>
                                            <img v-if="message.replymessage.photo" :src="message.replymessage.photo" style="max-width: 100px; max-height: 100px; margin: 10px;">
                                            <div style="margin-right: 10px; word-break: break-word;">{{message.replymessage.text}}</div>
                                        </div>
                                        <img v-if="message.photo" :src="message.photo" style="max-width: 200px; max-height: 200px; margin: 10px;">
                                        <div class="messagebox-text">
                                            {{message.text}}
                                        </div>
                                        <div class="messagebox-time">
                                            <img class="messagebox-checkmark" v-if="message.isallread" src="/assets/double-check-blue.svg" style="height: 24px; width: 24px;">
                                            <img class="messagebox-checkmark" v-else-if="message.isallreceived" src="/assets/double-check.svg" style="height: 24px; width: 24px;">
                                            <img class="messagebox-checkmark" v-else src="/assets/single-check.svg" style="height: 24px; width: 24px;">
                                            {{message.timestamp}}
                                        </div>
                                        <div class="messagebox-buttons">
                                            <img src="/assets/reply.svg" style="height: 24px; width: 24px; cursor: pointer;" @click="setReplyMessage(message)">
                                            <img src="/assets/forward.svg" style="height: 24px; width: 24px; cursor: pointer;" @click="startForwardingMessage(message)" id="forwardbutton">
                                            <img src="/assets/trashcan.svg" style="height: 24px; width: 24px; cursor: pointer;" @click="deleteMessage(message)">
                                            <div>
                                                <img src="/assets/comment.svg" style="height: 24px; width: 24px; cursor: pointer;" @click="showComments(message)">
                                                {{message.commentlist.length}}
                                            </div>
                                        </div>
                                        <div v-if="commentshown==message.messageid" class="messagebox-comment">
                                            <div class="emoji-picker">
                                                <span v-for="emoji in emojiList" 
                                                    :key="emoji" 
                                                    class="emoji-option"
                                                    @click="selectEmoji(emoji, message)">
                                                    {{emoji}}
                                                </span>
                                            </div>
                                            <ErrorMsg v-if="errormsg" :msg="errormsg" style="word-break: break-word; width: 200px;"></ErrorMsg>
                                            <div class="commentlist">
                                                <ul>
                                                    <li v-for="comment in message.commentlist" :key="comment.userid">
                                                        {{comment.username}}: {{comment.emoji}}
                                                        <img v-if="comment.userid==this.userid" src="/assets/trashcan.svg" style="height: 16px; width: 16px; cursor: pointer;" @click="deleteComment(message)">
                                                    </li>
                                                </ul>
                                            </div>
                                        </div>
                                    </div>
                                </span>
                                <span v-else style="display:flex; width: calc(100vw - 360px); height: 100%;">
                                    <div class="messagebox-other">
                                        <div v-if="message.isforwarded" class="forwarded-info">
                                            <img src="/assets/forward.svg" style="width: 24px; height: 24px; margin-right: 5px;">
                                            Forwarded
                                        </div>
                                        <div class="messagebox-username">
                                            <b><h3 style="margin-bottom: 0;">{{message.username}}</h3></b>
                                        </div>
                                        <div v-if="message.replymessage.messageid!=0" style="display: flex; flex-direction: column; margin-left: 15px; margin-right: 15px; background-color: #695d5d; width: calc(100% - 30px);">
                                            <div v-if="message.replymessage.username==this.username">
                                                <h5 style="margin-bottom: 0; margin-left: 10px;">You</h5>
                                            </div>
                                            <div v-else>
                                                <h5 style="margin-bottom: 0; margin-left: 10px;">{{message.replymessage.username}}</h5>
                                            </div>
                                            <img v-if="message.replymessage.photo" :src="message.replymessage.photo" style="max-width: 100px; max-height: 100px; margin: 10px;">
                                            <div style="margin-left: 10px; word-break: break-word;">{{message.replymessage.text}}</div>
                                        </div>
                                        <img v-if="message.photo" :src="message.photo" style="max-width: 200px; max-height: 200px; margin: 10px;">
                                        <div class="messagebox-text">
                                            {{message.text}}
                                        </div>
                                        <div class="messagebox-time">
                                            {{message.timestamp}}
                                        </div>
                                        <div class="messagebox-buttons">
                                            <div>
                                                <img src="/assets/comment.svg" style="height: 24px; width: 24px; cursor: pointer;" @click="showComments(message)">
                                                {{message.commentlist.length}}
                                            </div>
                                            <img src="/assets/forward.svg" style="height: 24px; width: 24px; cursor: pointer;" @click="startForwardingMessage(message)" id="forwardbutton">
                                            <img src="/assets/reply.svg" style="height: 24px; width: 24px; cursor: pointer;" @click="setReplyMessage(message)">
                                        </div>
                                        <div v-if="commentshown==message.messageid" class="messagebox-comment">
                                            <div class="emoji-picker">
                                                <span v-for="emoji in emojiList" 
                                                    :key="emoji" 
                                                    class="emoji-option"
                                                    @click="selectEmoji(emoji, message)">
                                                    {{emoji}}
                                                </span>
                                            </div>
                                            <ErrorMsg v-if="errormsg" :msg="errormsg" style="word-break: break-word; width: 200px;"></ErrorMsg>
                                            <div class="commentlist">
                                                <ul>
                                                    <li v-for="comment in message.commentlist" :key="comment.userid">
                                                        {{comment.username}}: {{comment.emoji}}
                                                        <img v-if="comment.userid==this.userid" src="/assets/trashcan.svg" style="height: 16px; width: 16px; cursor: pointer;" @click="deleteComment(message)">
                                                    </li>
                                                </ul>
                                            </div>
                                        </div>
                                    </div>
                                </span>
                            </li>
                        </ul>
                    </div>
                </div>
                <!-- Bottom bar in mainchat -->
                <div class="bottombar-chat">
                    <input type="file" accept="image/*" ref="sendPhotoInput" style="display: none;" @change="sendPhotoFileSelect"/>
                    <img v-if="!messagephoto" src="/assets/photo-icon.svg" @click="sendPhotoButton" style="width: 32px; height: 32px; cursor: pointer; margin-left: 10px; margin-right: 10px;">
                    <img v-else src="/assets/cross.svg" @click="messagephoto=null" style="width: 32px; height: 32px; cursor: pointer; margin-left: 10px; margin-right: 10px;">
                    <div class="message-text">
                        <input class="message-textbox" v-model="messagetext" placeholder="Write a message" @keyup.enter="sendMessageorCreateChat">
                    </div>
                    <img v-if="messagetext || messagephoto" src="/assets/send.svg" style="width: 32px; height: 32px; cursor: pointer; margin-left: 10px; margin-right: 10px;" @click="sendMessageorCreateChat">
                    <div v-if="messagephoto || replyid" class="messagephoto-preview">
                        <div v-if="replyid" style="color: whitesmoke; margin: 25px; background-color: #2a3942;">
                            <div class="replyname" style="display: flex; justify-content: space-between; align-items: center;">
                                <div v-if="this.replyuser==this.username">
                                    <h3 style="margin: 0; margin-left: 10px;">You</h3>
                                </div>
                                <div v-else>
                                    <h3 style="margin: 0; margin-left: 10px;">{{replyuser}}</h3>
                                </div>
                                <img src="/assets/cross.svg" style="width: 24px; height: 24px; cursor: pointer; margin-left: 10px; margin-right: 10px;" @click="resetReplyMessage()">
                            </div>
                            <img v-if="replyphoto" :src="replyphoto" style="margin: 10px; max-height: 100px;">
                            <div v-if="replytext" style="margin-left: 10px;">{{replytext}}</div>
                        </div>
                        <img v-if="messagephoto" :src="messagephoto" style="max-width: 250px; max-height: 250px; margin: 25px;">
                    </div>
                </div>
            </div>
        </div>
    </div>

    <!-- Box to change group name and photo -->
    <div class="box-container" v-if="boxshown == 2">
        <div class="blurred-box">
            <h1 style="margin-top: 20px;">Group Info</h1>
            <div class="new-username-box">
                Enter a new group name:
                <div class="new-username-container">
                    <input class="new-username" v-model="newgroupname" placeholder="New group name" />
                </div>
            </div>
            <input type="file" accept="image/*" ref="createGroupPhotoInput" style="display: none;" @change="createGroupPhotoFileSelect"/>
            <button class="selectphoto-button" @click="createGroupPhotoButton">Select Photo</button>
            <div v-if="newgroupphoto" style="display: flex; flex-direction: column; align-items: center;">
                Preview group photo:
                <img class="img-circular" :src="newgroupphoto" style="width: 64px; height: 64px; background-color: #695d5d;" />
            </div>
            <button class="confirm-button" @click="changeGroupNamePhoto">Confirm</button>
            <button class="cancel-button" @click="resetChangeGroupPrompt">Cancel</button>
            <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
        </div>
    </div>

    <!-- Box to create a group -->
    <div class="box-container" v-if="boxshown == 3">
        <div class="blurred-box">
            <h1 style="margin-top: 20px;">New Group</h1>
            <div class="new-username-box">
                Enter the group name:
                <div class="new-username-container">
                    <input class="new-username" v-model="createdgroupname" placeholder="Group name">
                </div>
            </div>
            <input type="file" accept="image/*" ref="createGroupPhotoInput" style="display: none;" @change="createGroupPhotoFileSelect"/>
            <button class="selectphoto-button" @click="createGroupPhotoButton">Select Photo</button>
            <div v-if="createdgroupphoto" style="display: flex; flex-direction: column; align-items: center;">
                Preview group photo:
                <img class="img-circular" :src="createdgroupphoto" style="width: 64px; height: 64px; background-color: #695d5d;"/>
            </div>
            <button v-if="createdgroupname" class="confirm-button" @click="boxshown = 4;this.getAllUsers();">Next</button>
            <button class="cancel-button" @click="resetCreateGroupPrompt">Cancel</button>
            <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
        </div>
    </div>

    <!-- Box to add people when creating a group -->
    <div class="box-container" v-if="boxshown == 4">
        <div class="blurred-box">
            <h1 style="margin-top: 20px;">New Group</h1>
            <h4>Select users</h4>
            <input v-model="searcheduser" placeholder="Search user" @keyup.enter="searchUser" />
            <div class="users-checkbox">
                <ul>
                    <li v-for="user in allusers" :key="user.userid" @click="insertUserCreation(user)">
                        <label>
                            <input type="checkbox" :checked="isSelectedCreation(user)" @change.prevent />
                            {{ user.username }}
                        </label>
                    </li>
                </ul>
            </div>
            <button v-if="userscreategroup.size>=1" class="confirm-button" @click="createGroup">Confirm</button>
            <button class="cancel-button" @click="boxshown = 3">Back</button>
            <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
        </div>
    </div>

    <!-- Box to add people to a group -->
    <div class="box-container" v-if="boxshown == 5">
        <div class="blurred-box">
            <h1 style="margin-top: 20px;">Add to group</h1>
            <h4>Select users</h4>
            <input v-model="searcheduser" placeholder="Search user" @keyup.enter="searchUser" />
            <div class="users-checkbox">
                <ul>
                    <li v-for="user in usersnotinchat" :key="user.userid" @click="insertUserAdding(user)">
                        <label>
                            <input type="checkbox" :checked="isSelectedAdding(user)" @change.prevent />
                            {{ user.username }}
                        </label>
                    </li>
                </ul>
            </div>
            <button v-if="userstoadd.size>=1" class="confirm-button" @click="addToGroup">Confirm</button>
            <button class="cancel-button" @click="resetAddUsersPrompt">Cancel</button>
            <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
        </div>
    </div>

    <!-- Box to show group members -->
    <div class="box-container" v-if="boxshown == 6">
        <div class="blurred-box">
            <h1 style="margin-top: 20px;">Group Members</h1>
            <div class="group-members-list">
                <ul>
                    <li v-for="member in groupMembers" :key="member.userid">
                        {{ member.username }}
                    </li>
                </ul>
            </div>
            <button class="cancel-button" @click="boxshown = 0">Close</button>
            <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
        </div>
    </div>
</template>

<style>
    body {
        background-color: #0e1621;
        font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
        color: white;
        margin: 0;
    }

    /* Circular images */
    .img-circular {
        border-radius: 50%;
        object-fit: cover;
        box-shadow: black 0 0 5px;
    }

    /* All page */
    .all-screen {
        width: 100vw;
        height: 100vh;
    }

    /* Navbar */
    .navbar-dark {
        background-color: #1f2c3d;
        position: relative;
        display: flex;
        align-items: center;
        height: 60px;
        width: 100%;
        box-shadow: black 0 0 1px;
    }

    /* Info of user */
    .user-info {
        margin-left: 10px;
        color: whitesmoke;
        background-color: #2a394200;
        box-shadow: rgb(0, 0, 0) 0 0 8px;
        border-radius: 20px;
        height: 36px;
        display: flex;
        justify-content: center;
        align-items: center;
        padding: 0 15px;
    }

    /* Searchbox */
    .searchbox {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        margin: 0;
        height: 36px;
        border-radius: 20px;
        background-color: #2a3942;
        display: flex;
        flex-direction: row;
        z-index: 100;
        padding: 0 10px;
        box-shadow: rgb(0, 0, 0) 0 0 8px;
    }
    .searchbox-userlist {
        position: relative;
        margin-right: 34px;
        margin-top: 6.5px;
        width: 200px;
    }
    .searchbox-user {
        color: whitesmoke;
        width: 100%;
        border: none;
        background-color: transparent;
        margin-bottom: 6.5px;
    }
    .searchbox-user:focus {
        outline: none;
    }
    .searched-dropdown {
        position: relative;
        background-color: #1f2c3d;
        width: 200px;
        max-height: 120px;
        overflow-y: auto;
        border-radius: 10px;
        margin-top: 5px;
        box-shadow: rgb(0, 0, 0) 0 0 8px;
    }
    .searched-dropdown ul {
        list-style: none;
        padding: 0;
        margin: 0;
        color: whitesmoke;
    }
    .searched-dropdown li {
        height: 40px;
        padding: 10px;
        cursor: pointer;
        box-shadow: black 0 0 1px;
    }
    .searched-dropdown li:hover {
        background-color: #2a3942;
    }

    /* Main screen */
    .main-screen {
        display: flex;
        flex-direction: row;
        height: calc(100vh - 60px);
        width: 100%;
    }

    /* Sidebar */
    .sidebar-chats {
        background-color: #1f2c3d;
        height: 100%;
        width: 350px;
        box-shadow: rgb(0, 0, 0) 0 0 15px;
    }
    .sidebar-buttons {
        width: 100%;
        height: 60px;
        display: flex;
        align-items: center;
        justify-content: space-around;
        padding: 0 10px;
        box-shadow: black 0 0 1px;
    }
    .chats-dropdown {
        width: 100%;
        height: calc(100% - 36px);
        overflow-y: auto;
    }
    .chats-dropdown ul {
        list-style: none;
        padding: 0;
        margin: 0;
        color: whitesmoke;
    }
    .chats-dropdown li {
        padding: 10px;
        cursor: pointer;
        height: 80px;
        display: flex;
        align-items: center;
        box-shadow: black 0 0 1px;
    }
    .chats-dropdown li:hover {
        background-color: #2a3942;
    }
    .chatpreview {
        width: 100%;
    }
    .chatpreviewname {
        display: flex;
        flex-direction: row;
        position: relative;
        width: 100%;
    }
    .messagepreview {
        display: flex;
        flex-direction: row;
        align-items: center;
        position: relative;
        width: 100%;
    }
    .checkmark {
        position: absolute;
        left: 95%;
        top: 50%;
        transform: translate(-50%, -50%);
    }
    .timepreview {
        position: absolute;
        left: 95%;
        top: 50%;
        transform: translate(-50%, -50%);
    }

    /* Main chat */
    .main-chat {
        height: 100%;
        width: calc(100% - 350px);
        background-color: #0e1621;
    }

    /* no chat message */
    .no-chat-message {
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100%;
        width: 100%;
        color: #ffffff;
        font-size: 1rem;
        font-style: italic;
        border-radius: 10px;
        padding: 20px;
        margin: 20px;
    }

    /* Topbar */
    .topbar-chat {
        display: flex;
        align-items: center;
        background-color: #1f2c3d;
        height: 60px;
        width: 100%;
        position: relative;
        box-shadow: black 0 0 1px;
    }
    .backarrow {
        margin-left: 10px;
    }
    .addusers-button {
        position: left;
        left: 60%;
        cursor: pointer;
    }
    .leavebutton {
        position: absolute;
        right: 25px;
        background-color: #f44336;
        color: whitesmoke;
        width: 120px;
        border-radius: 20px;
        border-color: #000000;
    }
    .leavebutton:hover {
        background-color: #b40c00;
    }

    /* Group Members button */
    .group-members {
        position: absolute;
        left: 300px;
        background-color: #2a3942;
        color: whitesmoke;
        width: 120px;
        border-radius: 20px;
        border-color: #000000;
    }
    .group-members:hover {
        background-color: #1f2c3d;
    }

    /* Stile per la lista dei membri del gruppo */
    .group-members-list {
        max-height: 200px;
        overflow-y: auto;
        margin-top: 10px;
    }

    .group-members-list ul {
        list-style: none;
        padding: 0;
        margin: 0;
        color: whitesmoke;
    }

    .group-members-list li {
        padding: 10px;
        cursor: pointer;
        box-shadow: black 0 0 1px;
    }

    .group-members-list li:hover {
        background-color: #2a3942;
    }

    /* Message screen */
    .message-screen {
        height: calc(100% - 120px);
        width: 100%;
        overflow-y: auto;
        overflow-x: hidden;
    }
    .messagelist {
        height: 100%;
        width: 100%;
        overflow-y: auto;
        overflow-x: hidden;
        display: flex;
        flex-direction: column;
    }
    .messagelist ul {
        list-style: none;
        padding: 0;
        color: whitesmoke;
    }
    .messagelist li {
        padding: 5px;
        position: relative;
    }

    /* Box containing each message */
    .messagebox-you {
        width: max-content;
        background-color: #25d366;
        border-radius: 20px;
        border-top-right-radius: 0;
        margin-right: 20px;
        max-width: 600px;
        display: flex;
        flex-direction: column;
        align-items: end;
        padding: 10px;
    }
    .messagebox-other {
        width: max-content;
        background-color: #2a3942;
        border-radius: 20px;
        border-top-left-radius: 0;
        margin-left: 20px;
        max-width: 600px;
        display: flex;
        flex-direction: column;
        align-items: start;
        padding: 10px;
    }
    .forwarded-info {
        margin-left: 15px;
        margin-right: 15px;
    }
    .messagebox-text {
        margin-left: 15px;
        margin-right: 15px;
        word-break: break-word;
        font-size: 0.875rem;
        font-family: sans-serif;
    }
    .messagebox-username {
        margin-left: 15px;
        margin-right: 15px;
    }
    .messagebox-time {
        display: flex;
        justify-content: right;
        align-items: center;
        width: 100%;
        padding-right: 15px;
    }
    .messagebox-checkmark {
        margin-right: 5px;
    }
    .messagebox-buttons {
        display: flex;
        align-items: center;
        justify-content: space-between;
        min-width: 150px;
        width: calc(100% - 30px);
        margin-bottom: 5px;
        margin-left: 15px;
        margin-right: 15px;
    }

    /* Box containing comments */
    .messagebox-comment {
        display: flex;
        flex-direction: column;
        align-items: center;
        width: 100%;
        border-top: 2px solid #2a3942;
        padding: 10px;
        background-color: #1f2c3d;
        border-radius: 0 0 20px 20px;
    }
    .emoji-picker {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;
        padding: 8px;
        background: #2a3942;
        border-radius: 8px;
        max-width: 200px;
        margin: 10px;
    }
    .emoji-option {
        font-size: 24px;
        cursor: pointer;
        transition: transform 0.2s;
    }
    .emoji-option:hover {
        transform: scale(1.2);
    }
    .commentlist {
        width: 100%;
        overflow-y: auto;
        display: flex;
        flex-direction: column;
        max-height: 150px;
    }
    .commentlist ul {
        list-style: none;
        padding: 0;
        color: whitesmoke;
    }
    .commentlist li {
        padding: 5px;
        padding-left: 10px;
        padding-right: 10px;
        height: 30px;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    /* Bottom bar */
    .bottombar-chat {
        display: flex;
        align-items: center;
        background-color: #1f2c3d;
        height: 60px;
        width: 100%;
        position: relative;
        box-shadow: black 0 0 1px;
    }
    .bottombar-chat img {
        filter: brightness(0) invert(1);
    }
    .message-text {
        height: 36px;
        border-radius: 20px;
        background-color: #2a3942;
        display: flex;
        align-items: center;
        width: 80%;
        padding: 0 10px;
    }
    .message-textbox {
        margin-left: 10px;
        margin-right: 10px;
        color: whitesmoke;
        background: transparent;
        border: none;
        width: 100%;
        height: 23px;
    }
    .message-textbox:focus {
        outline: none;
    }
    .messagephoto-preview {
        position: absolute;
        top: 0%;
        transform: translateY(-100%);
        background-color: #1f2c3d;
        border-top-right-radius: 20px;
        display: flex;
        flex-direction: column;
        justify-content: center;
    }

    /* List with users checkbox */
    .users-checkbox {
        position: relative;
        background-color: #1f2c3d;
        width: 200px;
        max-height: 200px;
        overflow-y: auto;
        border-radius: 10px;
        margin-top: 5px;
    }
    .users-checkbox ul {
        list-style: none;
        padding: 0;
        margin: 0;
        color: whitesmoke;
    }
    .users-checkbox li {
        height: 40px;
        padding: 10px;
        cursor: pointer;
    }
    .users-checkbox li:hover {
        background-color: #2a3942;
    }
</style>