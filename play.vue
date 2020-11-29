<template>
    <div class="Play">
        <div id="copydiv" style="width:0;height:0;"></div>
        <template v-if="game&&!noSuchRoom">
            <div class="header">第{{game.turn}}局 {{game.fetched}}/{{total}}</div>
            <template v-if="game.manager===me">
                <div class="roles">
                    <div v-for="role in roles.filter(role=>game.info[role.name])" class="role" :key="role.name">
                        <span class="role-name">{{role.chinese}}</span>
                        <span>{{game.info [role.name] }}</span>
                    </div>
                </div>
                <div class="footer">
                    <div>
                        <button @click="createRoom">创建房间</button>
                        <button @click="copyLink">复制链接</button>
                    </div>
                    <div>
                        <button @click="newGame">再来一局</button>
                        <button @click="fetchRoom">刷新</button>
                    </div>
                </div>
            </template>
            <template v-else>
                <div style="display: flex;align-items: center;">
                    <template v-if="role">
                        你是{{space}} <span class="myrole">{{role.chinese}}</span>
                    </template>
                    <template v-else><span class="dangerInfo">身份已被领完</span></template>
                </div>
                <div class="roles">
                    <div v-for="role in roles.filter(role=>game.info[role.name])" class="role" :key="role.name">
                        <span class="role-name">{{role.chinese}}</span>
                        <span>{{game.info [role.name] }}</span>
                    </div>
                </div>
                <div class="footer">
                    <div>
                        <button @click="createRoom">创建房间</button>
                        <button @click="copyLink">复制链接</button>
                    </div>
                    <div>
                        <button @click="fetchRoom">刷新</button>
                    </div>
                </div>
            </template>
        </template>
        <template v-else-if="noSuchRoom">
            <span class="dangerInfo">没有这个房间</span>
        </template>
        <template v-else>
            正在加载
        </template>
    </div>
</template>
<script>
    import * as lib from "./game"
    import * as api from "./api"

    export default {
        data() {
            return {
                me: '',
                game: null,
                room: null,
                roles: lib.roles,
                noSuchRoom: false,
            }
        },
        computed: {
            role() {
                let role = this.game.people[this.me]
                return lib.roleMap[role];
            },
            total() {
                return Object.values(this.game.info).reduce((o, n) => o + n, 0)
            }
        },
        mounted() {
            this.me = lib.getUser();
            console.log(`user=${lib.getUser()}`)
            this.fetchRoom();
        },
        methods: {
            createRoom() {
                location.href = "."
            },
            newGame() {
                const q = lib.parseQuery();
                api.newGame(q.room).then(resp => {
                    console.log('new game')
                    console.log(resp.data)
                    this.game = resp.data;
                })
            },
            fetchRoom() {
                //拉取房间信息
                const q = lib.parseQuery();
                api.fetch(q.room).then(resp => {
                    console.log('刷新信息')
                    console.log(resp.data)
                    if (resp.data === "no such room") {
                        this.noSuchRoom = true
                        return
                    }
                    this.game = resp.data;
                })
            },
            copyLink() {
                this.copyText(location.href)
            },
            copyText(text) {
                let textarea = document.createElement("textarea"); //创建input对象
                let currentFocus = document.activeElement; //当前获得焦点的元素
                let toolBoxwrap = document.querySelector("#copydiv"); //将文本框插入到NewsToolBox这个之后
                toolBoxwrap.appendChild(textarea); //添加元素
                textarea.value = text;
                textarea.focus();
                if (textarea.setSelectionRange) {
                    textarea.setSelectionRange(0, textarea.value.length); //获取光标起始位置到结束位置
                } else {
                    textarea.select();
                }
                let flag = document.execCommand("copy"); //执行复制
                toolBoxwrap.removeChild(textarea); //删除元素
                currentFocus.focus();
                return flag;
            }
        }
    }
</script>
<style lang="less">
    .Play {
        display: flex;
        flex-direction: column;
        align-items: center;

        .header, .footer {
            display: flex;
            font-size: 2.3rem;
            margin: 20px 0;

            button {
                font-size: 1.5rem;
            }
        }

        .footer {
            flex-direction: column;

            & > div {
                display: flex;
                margin: 2px;

                button {
                    flex: 1;
                }
            }
        }

        .roles {
            font-size: 1.3rem;

            input[type="checkbox"] {
                width: 2rem;
                height: 2rem;
            }

            .role {
                display: flex;
                align-items: center;
                justify-content: space-between;

                .role-name {
                    margin-right: 10px;
                }
            }
        }

        .myrole {
            font-size: 2.1rem;
            color: red;
        }

        .dangerInfo {
            font-size: 1.5rem;
            color: red;
        }
    }
</style>