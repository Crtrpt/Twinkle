<template>
     <div class="my-5">

        <div class="tabs">
        <a class="tab tab-lg tab-lifted tab-active">本机代理</a> 
        <a class="tab tab-lg tab-lifted ">远程隧道</a> 
        <a class="tab tab-lg tab-lifted flex-grow "></a> 
        </div>
        <div class="form-control w-full max-w-xs mt-2" >
        <input type="text" placeholder="名称" class="input input-bordered w-full max-w-xs" />
        </div>
        <div class="form-control w-full max-w-xs mt-2">
        <textarea class="textarea input-bordered w-full" placeholder="描述信息"></textarea>
        </div>

        <div class="form-control w-36">
        <label class="cursor-pointer label">
            <span class="label-text">是否启用:</span> 
            <input type="checkbox" class="toggle" checked />
        </label>
        </div>

        <div class="form-control w-full max-w-xs mt-2">
        <input type="text" placeholder="本地地址" class="input input-bordered w-full max-w-xs" />
        </div>
        <div class="form-control w-full max-w-xs mt-2">
        <input type="text" placeholder="代理地址" class="input input-bordered w-full max-w-xs" />
        </div>

        </div>
</template>

<script lang="ts">
export default {

}
</script>

<style>

</style>