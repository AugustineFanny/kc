<template>
  <div>
    <el-form ref="create" style="width:500px;" :model="form" :rules="rules" label-width="80px">
      <el-form-item>
        <el-switch
          v-model="batch"
          active-text="平台全部用户"
          inactive-text="单个手机号">
        </el-switch>
      </el-form-item>
      <el-form-item v-if="batch == false" label="手机号">
        <el-input v-model="form.mobile" placeholder="单个手机号"></el-input>
      </el-form-item>
      <el-form-item label="发送内容" prop="content" required>
        <el-input type="textarea" v-model="form.content" placeholder="发送内容"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" :loading="loading" @click="onSubmit">发送</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
  export default {
    data() {
      return {
        loading: false,
        batch: false,
        form: {
          mobile: "",
          content: "",
        },
        rules: {

        },
      };
    },
    methods: {
      onSubmit() {
        this.$refs["create"].validate((valid) => {
          if (valid) {
            this.loading = true;
            this.$http.post("/kc/admin/super/batch-sms", {
              batch: this.batch,
              mobile: this.form.mobile,
              content: this.form.content,
            })
            .then(response => {
              var data = this.utilHelper.handle(this, response);
              if(data !== null) {
                this.$message.success("操作成功");
              };
              this.loading = false;
            }, response => {
              this.$message.error(response.body);
              this.loading = false;
            });
          } else {
            return false;
          }
        });
      }
    }
  };
</script>
