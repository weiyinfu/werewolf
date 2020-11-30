<template>
    <div class="Create" v-if="game">
        <div class="header">总人数：{{total}}</div>
        <div class="roles">
            <template v-for="role in roles">
                <div v-if="role.min===role.max" class="role" :key="role.name">
                    {{role.chinese}}<input type="checkbox" v-model="game[role.name]">
                </div>
                <div v-else class="role" :key="role.name">
                    {{role.chinese}} {{game[role.name]}}<input type="range" v-model="game[role.name]" step="1"
                                                               :min="role.min" :max="role.max">
                </div>
            </template>
        </div>
        <div class="footer">
            <button @click="createRoom">创建房间</button>
        </div>
    </div>
</template>
<script>
    import * as lib from "./game"
    import * as api from "./api"

    export default {
        data() {
            return {
                game: null,
                roles: lib.roles,
            }
        },
        computed: {
            total() {
                return Object.values(this.game).reduce((o, n) => o + (typeof (n) === 'boolean' ? (n ? 1 : 0) : parseInt(n)), 0);
            }
        },
        mounted() {
            let game = {};
            for (let role of lib.roles) {
                game[role.name] = role.default;
            }
            this.game = game;
        },
        methods: {
            createRoom() {
                //创建房间
                let game = this.game
                for (let k of Object.keys(this.game)) {
                    if (typeof this.game[k] === 'boolean') {
                        game[k] = this.game[k] ? 1 : 0;
                    } else {
                        game[k] = parseInt(this.game[k])
                    }
                }
                api.createRoom(this.game).then(resp => {
                    location.href = "?room=" + resp.data.id;
                })
            },
        }
    }
</script>
<style lang="less">
    .Create {
        display: flex;
        flex-direction: column;
        align-items: center;

        .header, .footer {
            display: flex;
            font-size: 2.3rem;
            margin: 20px 0;

            button {
                font-size: 2.3rem;
            }
        }


        .roles {
            font-size: 1.8rem;

            input[type="checkbox"] {
                width: 2rem;
                height: 2rem;
            }

            .role {
                display: flex;
                align-items: center;
                justify-content: space-between;
            }
        }

    }
</style>