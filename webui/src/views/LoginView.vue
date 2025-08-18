<script>
    export default {
        data: function() {
            return {
                errormsg: null,
                username: null
            }
        },
        methods: {
            async login() {
                this.errormsg = null;
                try {
                    let response = await this.$axios.post("/session", {username: this.username.trim()});
                    let userinfo = response.data;
                    localStorage.setItem('username', this.username);
                    localStorage.setItem('userid', userinfo.userid);
                    localStorage.setItem('userphoto',userinfo.photo);
                    this.$router.push({path: "/chats"});
                } catch (e) {
                    this.errormsg = e.response.data;
                }
            }
        }
    }
</script>

<template>
        <div class="blurred-box">
            <h1 style="position: absolute; top:10%;">WASAText</h1>
            <div class="input-field">
                <input class="username-input" v-model="username" placeholder="Username" @keyup.enter="login">
            </div>
            <ErrorMsg v-if="errormsg" :msg="errormsg" style="position: absolute; top: 60%;"></ErrorMsg>
            <button class="login-button" @click="login">Login</button>
        </div>
</template>

<style>
    body {
        background-color: #0e1621; /* Colore di sfondo scuro simile a Telegram */
        font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; /* Font moderno e leggibile */
        color: white;
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100vh;
        margin: 0;
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
    }

    .input-field {
        position: absolute;
        top: 45%;
        border-radius: 20px;
        width: 70%;
        padding: 10px;
        background-color: #1e2a32; /* Colore di sfondo scuro per il campo di input */
        border: 1px solid #2a3942; /* Bordo sottile come Telegram */
        color: white;
        font-size: 16px;
    }

    .username-input {
        background: none;
        border: none;
        outline: none;
        width: 100%;
        color: white; /* Testo bianco per contrasto */
    }

    .login-button {
        position: absolute;
        top: 65%;
        padding: 10px 20px;
        border-radius: 20px;
        width: 70%;
        border: none;
        outline: none;
        background-color: #25d366; /* Verde di WhatsApp per il pulsante */
        color: white;
        font-size: 16px;
        cursor: pointer;
        transition: background-color 0.3s ease;
    }

    .login-button:hover {
        background-color: #128c7e; /* Verde più scuro al passaggio del mouse */
    }

    .logo {
        width: 100px;
        height: 100px;
        margin-bottom: 20px;
    }

    .title {
        font-size: 24px;
        font-weight: bold;
        margin-bottom: 10px;
        color: #25d366; /* Colore verde di WhatsApp per il titolo */
    }

    .subtitle {
        font-size: 14px;
        color: #a8b2b9; /* Colore grigio chiaro per il sottotitolo */
    }
</style>