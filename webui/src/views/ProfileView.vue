<script>
export default {
    data(){
        return {
            newusername: null,
            newuserphoto: null,
            errormsg: null,
            changedinfo: false,
            userid: localStorage.getItem('userid'), // or another method to get the user id
        }
    },
    methods: {
        async changeUsernamePhoto(){
                    this.errormsg = null;
                    if (this.newusername){
                        try {
                            await this.$axios.put("/users/"+this.userid+"/name", {username: this.newusername.trim()},{headers:{"Authorization": `Bearer ${this.userid}`}});
                            localStorage.setItem('username', this.newusername);
                            this.username = this.newusername;
                            this.changedinfo = true;
                        } catch (e) {
                            this.errormsg = e.response.data;
                        }
                    }
                    if (this.newuserphoto){
                        try {
                            await this.$axios.put("/users/"+this.userid+"/photo",{photo: this.newuserphoto},{headers:{"Authorization": `Bearer ${this.userid}`}});
                            localStorage.setItem('userphoto', this.newuserphoto);
                            this.userphoto = this.newuserphoto;
                            this.changedinfo = true;
                        } catch (e) {
                            this.errormsg = e.response.data;
                        }
                    }
                    this.newusername = null;
                    this.newuserphoto = null;
                    if (this.changedinfo){
                        this.boxshown = 0;
                        this.errormsg = null;
                        this.changedinfo = false;
                        this.navigateToHome();
                    }
                },
                // Button to change the photo of user handler
                changePhotoFileSelect(){
                    const file = this.$refs.changePhotoInput.files[0];
                    if (file) {
                        const reader = new FileReader();
                        reader.onload = (e) => {
                            this.newuserphoto = e.target.result;
                        };
                        reader.readAsDataURL(file);
                    }
                },
                changePhotoButton(){
                    this.$refs.changePhotoInput.click();
                },
                async resetChangeUsernamePrompt(){
                    this.newusername = null;
                    this.newuserphoto = null;
                    this.boxshown = 0;
                    this.errormsg = null;
                    this.navigateToHome()
                },
                navigateToHome() {
                    this.$router.push({path: "/chats"});
                },
            }
        }
</script>

<template>
    <!-- Box to change username and photo -->
    <div class="box-container">
        <div class="blurred-box">
            <h1 style="margin-top: 20px;">Profile Info</h1>
            <div class="new-username-box">
                Enter a new username:
                <div class="new-username-container">
                    <input class="new-username" v-model="newusername" placeholder="New username">
                </div>
            </div>
            <input type="file" accept="image/*" ref="changePhotoInput" style="display: none;" @change="changePhotoFileSelect"/>
            <button class="selectphoto-button" @click="changePhotoButton">Select Photo</button>
            <div v-if="newuserphoto" style="display: flex; flex-direction: column; align-items: center;">
                Preview profile pic:
                <img class="img-circular" :src="newuserphoto" style="width: 64px; height: 64px; background-color: #695d5d;"/>
            </div>
            <button class="confirm-button" @click="changeUsernamePhoto">Confirm</button>
            <button class="cancel-button" @click="resetChangeUsernamePrompt">Cancel</button>
            <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
        </div>
    </div>
</template>

<style>
    /* Box change username and photo */
    .box-container {
        display: flex;
        justify-content: center;
        align-items: center;
        width: 100%;
        height: 100%;
        position: absolute;
    }
    .blurred-box {
        display: flex;
        flex-direction: column;
        align-items: center;
        width: 400px;
        height: 400px;
        background-color: rgba(32, 44, 51, 0.8); /* Sfondo semi-trasparente simile a Telegram */
        border-radius: 15px; /* Bordi arrotondati come WhatsApp */
        backdrop-filter: blur(10px); /* Effetto di sfocatura */
        box-shadow: 0 4px 10px rgba(0, 0, 0, 0.2); /* Ombreggiatura per profondità */
        padding: 20px;
        position: relative;
        color: whitesmoke;
    }
    .new-username-box {
        margin: 20px;
        width: 80%;
    }
    .new-username-container {
        width: 100%;
        padding: 10px;
        background-color: #1f2c3d; /* Colore di sfondo scuro per il campo di input */
        border-radius: 20px;
        border: 1px solid #2a3942; /* Bordo sottile come Telegram */
    }
    .new-username {
        background: none;
        border: none;
        outline: none;
        width: 100%;
        color: whitesmoke;
    }
    .selectphoto-button {
        padding: 10px 20px;
        border-radius: 20px;
        border: none;
        outline: none;
        background-color: #25d366; /* Verde di WhatsApp per il pulsante */
        color: white;
        cursor: pointer;
        transition: background-color 0.3s ease;
        margin-bottom: 10px;
    }
    .selectphoto-button:hover {
        background-color: #128c7e; /* Verde più scuro al passaggio del mouse */
    }
    .confirm-button {
        position: absolute;
        top: 90%;
        left: 80%;
        transform: translate(-50%, -50%);
        padding: 10px 20px;
        border-radius: 20px;
        border: none;
        outline: none;
        background-color: #25d366; /* Verde di WhatsApp per il pulsante */
        color: white;
        cursor: pointer;
        transition: background-color 0.3s ease;
    }
    .confirm-button:hover {
        background-color: #128c7e; /* Verde più scuro al passaggio del mouse */
    }
    .cancel-button {
        position: absolute;
        top: 90%;
        left: 20%;
        transform: translate(-50%, -50%);
        padding: 10px 20px;
        border-radius: 20px;
        border: none;
        outline: none;
        background-color: #f44336; /* Rosso per il pulsante di cancellazione */
        color: white;
        cursor: pointer;
        transition: background-color 0.3s ease;
    }
    .cancel-button:hover {
        background-color: #d32f2f; /* Rosso più scuro al passaggio del mouse */
    }
</style>
