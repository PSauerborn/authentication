<template>
    <v-col cols=4>
        <v-sheet>
            <v-card>
                <v-card-title>
                    <v-row align="center" justify="center">
                        <v-col cols=12 align="center" justify="center">
                            Welcome to Monty!
                        </v-col>
                    </v-row>
                </v-card-title>
                <v-divider class="mx-4"></v-divider>
                <br>
                <v-tabs centered=true>
                    <v-tab>Login</v-tab>
                    <v-tab-item>
                        <v-card-text>
                            <v-row align="center" justify="center" dense>
                                <v-col cols=6>
                                    <v-text-field v-model="username" prepend-icon="mdi-account" label="Username" dense></v-text-field>
                                </v-col>
                            </v-row>
                            <v-row align="center" justify="center" dense>
                                <v-col cols=6>
                                    <v-text-field v-model="password" prepend-icon="mdi-key-variant" label="Password" type="password" dense></v-text-field>
                                </v-col>
                            </v-row>
                        </v-card-text>
                        <v-card-text>
                            <v-row align="center" justify="center" dense>
                                <v-col cols=6>
                                    <v-btn class="login-button" block large color='info' @click="login">Login</v-btn>
                                </v-col>
                            </v-row>
                        </v-card-text>
                    </v-tab-item>

                    <v-tab>Sign Up</v-tab>
                    <v-tab-item>
                        <v-card-text>
                            <v-row align="center" justify="center" dense>
                                <v-col cols=6>
                                    <v-text-field v-model="newUsername" prepend-icon="mdi-account" label="Username" dense></v-text-field>
                                </v-col>
                            </v-row>
                            <v-row align="center" justify="center" dense>
                                <v-col cols=6>
                                    <v-text-field v-model="newEmail" prepend-icon="mdi-at" label="Email" dense></v-text-field>
                                </v-col>
                            </v-row>
                            <v-row align="center" justify="center" dense>
                                <v-col cols=6>
                                    <v-text-field v-model="newPassword" prepend-icon="mdi-key-variant" label="Password" type="password" dense></v-text-field>
                                </v-col>
                            </v-row>
                            <v-row align="center" justify="center" dense>
                                <v-col cols=6>
                                    <v-text-field v-model="newPasswordRepeat" prepend-icon="mdi-key-variant" label="Repeat Password" type="password" dense></v-text-field>
                                </v-col>
                            </v-row>
                        </v-card-text>
                        <v-card-text>
                            <v-row align="center" justify="center" dense>
                                <v-col cols=6>
                                    <v-btn class="login-button" block large color='info' @click="signup">Sign Up</v-btn>
                                </v-col>
                            </v-row>
                        </v-card-text>
                    </v-tab-item>
                </v-tabs>
            </v-card>
        </v-sheet>
    </v-col>
</template>


<script>
import axios from 'axios';

export default {
    name: "Login",
    methods: {
        /**
         * Function used to create a new user with the identity
         * provider. Users are created by making and POST request
         * to the identity provider route with the UID, email and
         * password in the request body
         */
        signup: function() {
            // retrieve application URL from environment variables
            const url = process.env.VUE_APP_IDP_URL + '/signup'
            let vm = this;

            // send error message if passwords dont match
            if (vm.newPassword != vm.newPasswordRepeat) {
                vm.$notify({
                    group: 'main',
                    title: ' monty backend',
                    type: 'error',
                    text: 'passwords must match'
                })
                return
            }

            axios({
                method: 'post',
                url: url,
                data: {uid: vm.newUsername, password: vm.newPassword, email: vm.newEmail, admin: false}
            }).then(function (response) {
                // parse payload and display notification
                console.log(JSON.stringify(response))
                 vm.$notify({
                    group: 'main',
                    title: 'Identity Provider',
                    type: 'success',
                    text: 'successfully logged in user ' + vm.username
                })
                // set username and password to current values and retrieve token
                vm.username = vm.newUsername
                vm.password = vm.newPassword
                vm.login()
            }).catch(function (error) {
                console.log(error)
                vm.$notify({
                    group: 'main',
                    title: 'Identity Provider',
                    type: 'error',
                    text: 'failed login user: ' + error.response.data.message
                })
            })
        },
        /**
         * Function used to login users and retrieve JWToken required
         * for authentication. All JWTokens are stored in a users local
         * storage
         */
        login: function() {
            const url = process.env.VUE_APP_IDP_URL + '/token'
            let vm = this;

            axios({
                method: 'post',
                url: url,
                data: {uid: vm.username, password: vm.password}
            }).then(function (response) {
                // parse payload and display notification
                 vm.$notify({
                    group: 'main',
                    title: ' idP backend',
                    type: 'success',
                    text: 'successfully logged in user ' + vm.username
                })
                localStorage.setItem('userToken', response.data.token)
                // get redirect URI from request path if present and redirect client
                const redirect = vm.getRedirect()
                if (redirect) {
                    window.location.href = redirect
                }
            }).catch(function (error) {
                console.log(error)
                vm.$notify({
                    group: 'main',
                    title: ' monty backend',
                    type: 'error',
                    text: 'failed login user: ' + error.response.data.message
                })
            })
        },
        /**
         * Function used to retrieve redirect url from query parameters
         * if present in request path
         */
        getRedirect: function() {
            const urlParams = new URLSearchParams(window.location.search)
            return urlParams.get('redirect_uri')
        }
    },
    data() {
        return {
            username: "",
            password: "",
            newUsername: "",
            newPassword: "",
            newPasswordRepeat: "",
            newEmail: ""
        }
    }
}
</script>

<style scoped>
.blue-span {
    font-weight: bold;
}

.blue-span:hover {
    cursor: pointer;
    color:#2196F3;
}
</style>