<template>
  <div id="app">
    <el-container style="height: 100vh;">
      <el-header style="text-align: right; font-size: 12px">
        <el-dropdown>
          <i class="el-icon-user"></i>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item>个人中心</el-dropdown-item>
              <el-dropdown-item>退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </el-header>
      <el-container>
        <el-aside width="200px" style="background-color: rgb(238, 241, 246)">
          <el-menu :default-openeds="['1']">
            <el-sub-menu index="1">
              <template #title>
                <i class="el-icon-menu"></i>
                <span>音乐库</span>
              </template>
              <el-menu-item index="1-1">全部歌曲</el-menu-item>
              <el-menu-item index="1-2">我的歌单</el-menu-item>
            </el-sub-menu>
            <el-sub-menu index="2">
              <template #title>
                <i class="el-icon-setting"></i>
                <span>设置</span>
              </template>
              <el-menu-item index="2-1">系统设置</el-menu-item>
            </el-sub-menu>
          </el-menu>
        </el-aside>
        <el-main>
          <router-view />
        </el-main>
      </el-container>
      <el-footer style="text-align: center;">
        HomeMusic © 2024 家庭音乐播放器
      </el-footer>
    </el-container>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}
</style>
