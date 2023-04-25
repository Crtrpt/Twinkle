<template>
    <div class=" w-full  h-screen flex  flex-col">
        <top-nav>
            <template v-slot:right>
                <li><label for="my-modal" >添加</label></li>
            </template>
        </top-nav>

        <input type="checkbox" id="my-modal" class="modal-toggle" />
                <div class="modal">
                <div class="modal-box ">
                    <h3 class="font-bold text-lg">添加代理</h3>
                    <proxy-edit :proxy=proxy></proxy-edit>
                    <div class="modal-action">
                    <label for="my-modal" class="btn btn-sm">确认</label>
                    <label for="my-modal" class="btn btn-error  text-white btn-sm">删除</label>
                    </div>
                </div>
        </div>

        <div class="h-full  flex flex-col justify-center items-center">
            <div class=" grid
             grid-cols-4 gap-4">
                <div v-for="proxy, idx in list" :key="idx" >
                    <div class="card w-96 bg-base-100 shadow hover:shadow-md" for="my-modal">
                        <div class="card-body">
                            <h2 class="card-title">{{ proxy?.Name }}</h2>
                            <p>{{ proxy?.Desc }}</p>
                            <div class="card-actions justify-end">
                                <input type="checkbox" class="toggle  text-sm" :checked="proxy.Enable" />
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>


<script setup lang="ts">
import api from "../api/index"

</script>

<script lang="ts">
import TopNav from '../components/TopNav.vue'
import ProxyEdit from './ProxyEdit.vue'
export default {
    components: { TopNav ,ProxyEdit},
    data() {
        return {
            proxy:{}  as tw.TWProxy, 
            list: {} as  { [key: string]: tw.TWProxy }
        }
    },
    mounted() {
        api.get("/api", {}).then((res) => {
            this.list = res;
        })
    }
}
</script>

<style></style>