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
        background: linear-gradient(135deg, #ffffff, #e6f0ff); /* sfondo chiaro contrastato */
        font-family: 'Trebuchet MS', Tahoma, Geneva, Verdana, sans-serif; /* font usato prima */
        color: #1a1a1a;
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
        background: rgba(255, 255, 255, 0.85); /* box semitrasparente più leggibile */
        border-radius: 20px;
        backdrop-filter: blur(15px);
        box-shadow: 0 8px 20px rgba(0, 0, 0, 0.25);
        padding: 20px;
        position: relative;
    }

    .input-field {
        position: absolute;
        top: 45%;
        border-radius: 25px;
        width: 70%;
        padding: 12px;
        background-color: #f0f6ff; 
        border: 1px solid #bcd6ff;
        color: #1a1a1a;
        font-size: 16px;
    }

    .username-input {
        background: none;
        border: none;
        outline: none;
        width: 100%;
        color: #1a1a1a;
        font-size: 16px;
    }

    .login-button {
        position: absolute;
        top: 65%;
        padding: 12px 20px;
        border-radius: 25px;
        width: 70%;
        border: none;
        outline: none;
        background-color: #2f65d9; /* blu acceso */
        color: white;
        font-size: 16px;
        cursor: pointer;
        font-weight: bold;
        letter-spacing: 1px;
        transition: all 0.3s ease;
    }

    .login-button:hover {
        background-color: #204a9c; /* blu più scuro */
        transform: scale(1.05);
    }

    .logo {
        width: 100px;
        height: 100px;
        margin-bottom: 20px;
    }

    .title {
        font-size: 26px;
        font-weight: bold;
        margin-bottom: 10px;
        color: #4ade80; /* verde acceso */
        text-shadow: 1px 1px 4px rgba(0,0,0,0.2);
    }

    .subtitle {
        font-size: 14px;
        color: #444;
    }
</style>
