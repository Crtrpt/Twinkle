<template>
  <div class="navbar bg-base-100">
    <div class="flex-1">
      <a class="btn btn-ghost normal-case text-xl">Twinkle</a>
    </div>
    <div class="flex-none">
      <ul class="menu menu-horizontal px-1">

       
       <slot name="right"></slot>
        <!-- <li tabindex="0">
          <a>
            Parent
            <svg class="fill-current" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24">
              <path d="M7.41,8.58L12,13.17L16.59,8.58L18,10L12,16L6,10L7.41,8.58Z" />
            </svg>
          </a>
          <ul class="p-2 bg-base-100">
            <li><a>Submenu 1</a></li>
            <li><a>Submenu 2</a></li>
          </ul>
        </li> -->
        <li><a href="./doc">文档</a></li>
        <!-- <li> <router-link :to="{name:'login'}">登录</router-link></li> -->
       
      </ul>
    </div>
  </div>
</template>

<script lang="ts">
export default {
  name:"TopNav"
}
</script>

<style></style>