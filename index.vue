<template>
    <div id="main">
        <PlayComponent v-if="room"></PlayComponent>
        <CreateComponent v-else></CreateComponent>
    </div>
</template>
<script>
    import * as lib from "./game"
    import PlayComponent from "./play.vue"
    import CreateComponent from "./create.vue"

    export default {
        components: {PlayComponent, CreateComponent},
        data() {
            return {
                room: null
            }
        },
        mounted() {
            let q = lib.parseQuery()
            this.room = q.room;
        }
    }
</script>
<style lang="less">
    html,body,#main{
        width:100%;
        height: 100%;
        padding:0;
        margin:0;
    }

    #main {
        display: flex;
        align-items: center;
        justify-content: center;
        flex-direction: column;
        height: fit-content;
    }

</style>