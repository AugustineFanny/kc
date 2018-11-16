<template>
  <div>
    <el-row><el-button type="primary" @click="showCreate = true">创建分组</el-button></el-row>
    <el-row>
      <el-col :xs="24" :sm="24" :md="16" :lg="16">
        <el-table :data="tableData" border style="width: 100%; max-width: 281px;">
          <el-table-column prop="id" label="" width="80"></el-table-column>
          <el-table-column prop="name" label="组名" width="200"></el-table-column>
        </el-table>
      </el-col>
    </el-row>
    <el-dialog title="创建分组" :visible.sync="showCreate">
      <el-form ref="create" :model="form2" :rules="rules2" label-width="100px" class="demo-ruleForm">
        <el-form-item label="分组" prop="name">
          <el-input v-model="form2.name" auto-complete="off"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="showCreate = false">取 消</el-button>
        <el-button type="primary" @click="onCreate">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
  export default {
    data() {
      return {
        tableData: [],
        showCreate: false,
        form2: {
          name: '',
        },
        rules2: {
          name: [
            { required: true, message: "请输入组名", trigger: 'blur'}
          ]
        }
      };
    },
    created () {
      this.getGroups();
	  },
    methods: {
      getGroups() {
        this.$http.get("/kc/admin/groups")
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data != null) {
            this.tableData = data;
          }
        }, response => {
          this.$message.error(response.body);
        });
      },
      onCreate() {
        this.$refs["create"].validate((valid) => {
          if (valid) {
            this.$http.post("/kc/admin/super/groups", {
              name: this.form2.name,
            })
            .then(response => {
              var data = this.utilHelper.handle(this, response);
              if(data !== null) {
                this.$message.success("操作成功");
                this.$refs["create"].resetFields();
                this.showCreate = false;
                this.getGroups();
              };
            }, response => {
              this.$message.error(response.body);
            });
          } else {
            return false;
          }
        });
      }
    }
  };
</script>

<style>
  .el-row {
    margin-bottom: 20px;
    &:last-child {
      margin-bottom: 0;
    }
  }
</style>
