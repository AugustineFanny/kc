<template>
  <div>
    <el-row>
    <el-col :span="12" :offset="6">
      <el-form ref="form" :model="form" :label-position="labelPosition" label-width="80px">
      <el-form-item label="账号">
        <el-input v-model="form.ident"></el-input>
      </el-form-item>
      <el-form-item label="密码">
        <el-input type="password" v-model="form.password"></el-input>
      </el-form-item>
      <el-form-item label="验证码">
        <el-col :span="15">
          <el-input v-model="form.code"></el-input>
        </el-col>
        <el-col :span="9">
          <img :src="imageUrl" @click="onRefresh" class="avatar">
        </el-col>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">登录</el-button>
      </el-form-item>
    </el-form>
    </el-col>
  </el-row>
  </div>
</template>

<script>
  import router from '../router'
  import utils from '../utils'

  export default {
    data() {
      return {
        labelPosition: 'left',
        imageUrl: '/kc/captcha/code?_=' + Date.now(),
        form: {
          ident: '',
          password: '',
          code: ''
        }
      }
    },
    methods: {
      onRefresh() {
        this.imageUrl = '/kc/captcha/code?_=' + Date.now();
      },
      onSubmit() {
        this.$http.post("/kc/admin-login", this.form).then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data !== null) {
            localStorage.user = data.username;
            localStorage.mobile = data.mobile;
            localStorage.super = data.super;
            router.push('/');
          }
        }, response => {
          console.log(response);
          this.$message.error(response.body);
        })
      }
    }
  }
</script>
