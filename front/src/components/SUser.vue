<template>
  <div>
    <el-row>
      <el-col :xs="24" :sm="24" :md="16" :lg="16">
        <el-button type="primary" @click="showCreate = true">添加</el-button>
        <el-upload style="display: inline-block" action="/kc/simulation/batch/users" :show-file-list="false" accept=".csv" :on-success="onSuccess">
          <el-button type="primary">批量添加</el-button>
        </el-upload>
        <el-button type="primary" @click="onDeleteAll">全部删除</el-button>
        <router-link :to="{ path: '/manage/simulation'}"><el-button type="primary" style="float: right">返回</el-button></router-link>
      </el-col>
    </el-row>
    <el-row>
      <el-col :xs="24" :sm="24" :md="16" :lg="16">
        <el-table :data="tableData" border style="width: 100%; max-width: 901px">
          <el-table-column prop="username" label="label" width="100"></el-table-column>
          <el-table-column prop="inviter" label="推荐人" width="100"></el-table-column>
          <el-table-column prop="parents_name" label="parents" width="600"></el-table-column>
          <el-table-column label="操作" width="100">
            <template slot-scope="scope">
              <el-button @click="handleClick(scope.row)" type="text" size="small" v-show="false">编辑</el-button>
              <el-button @click="onDelete(scope.row.id)" type="text" size="small">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-col>
      <el-col :xs="24" :sm="24" :md="8" :lg="8" v-show="show">
        <el-form ref="change" :model="form" :rules="rules" label-width="120px">
          <el-form-item label="label" prop="username">
            <el-input v-model="form.username" auto-complete="off" readonly></el-input>
          </el-form-item>
          <el-form-item label="邀请人" prop="inviter">
            <el-select v-model="form.inviter" clearable placeholder="请选择">
              <el-option
                v-for="item in tableData"
                :key="item.username"
                :label="item.username"
                :value="item.username">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="onChange">提交</el-button>
            <el-button @click="onCancel">取消</el-button>
          </el-form-item>
        </el-form>
      </el-col>
    </el-row>
    <el-dialog title="添加" :visible.sync="showCreate">
      <el-form ref="create" :model="form2" :rules="rules2" label-width="100px" class="demo-ruleForm">
        <el-form-item label="label" prop="username">
          <el-input v-model="form2.username" auto-complete="off"></el-input>
        </el-form-item>
        <el-form-item label="邀请人" prop="inviter">
          <el-select v-model="form2.inviter" clearable placeholder="请选择">
            <el-option filterable
              v-for="item in tableData"
              :key="item.username"
              :label="item.username"
              :value="item.username">
            </el-option>
          </el-select>
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
        show: false,
        form: {
          username: "",
          inviter: "",
        },
        rules: {

        },
        showCreate: false,
        form2: {
          username: "",
          inviter: "",
        },
        rules2: {

        }
      };
    },
    created () {
      this.getSUsers();
    },
    methods: {
      getSUsers() {
        this.$http.get("/kc/simulation/users")
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data != null) {
            this.tableData = data;
          }
        }, response => {
          this.$message.error(response.body);
        });
      },
      handleClick(row) {
        this.show = true;
        for (var i in row) {
            this.form[i] = row[i];
        };
      },
      onChange() {
        this.$refs["change"].validate((valid) => {
          if (valid) {
            this.$resource("/kc/simulation/user/{id}")
            .save({id: this.form.id}, {
              username: this.form.username,
              parent: this.form.inviter,
            })
            .then(response => {
              var data = this.utilHelper.handle(this, response);
              if(data !== null) {
                this.$message.success("操作成功");
                this.$refs["change"].resetFields();
                this.show = false;
                this.getSUsers();
              };
            }, response => {
              this.$message.error(response.body);
            });
          } else {
            return false;
          }
        });
      },
      onCreate() {
        this.$refs["create"].validate((valid) => {
          if (valid) {
            this.$http.post("/kc/simulation/users", {
              username: this.form2.username,
              parent: this.form2.inviter,
            })
            .then(response => {
              var data = this.utilHelper.handle(this, response);
              if(data !== null) {
                this.$message.success("操作成功");
                this.$refs["create"].resetFields();
                this.showCreate = false;
                this.getSUsers();
              };
            }, response => {
              this.$message.error(response.body);
            });
          } else {
            return false;
          }
        });
      },
      onSuccess(response) {
        if(response.code == "100103") {
          router.push("/login");
        } else if(response.code == "100200") {
          this.getSUsers();
        } else {
          this.$message.error(response.msg);
          this.getSUsers();
        }
      },
      onCancel() {
        this.show = false;
      },
      onDelete(uid) {
        this.$http.delete("/kc/simulation/user/" + uid)
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data !== null) {
            this.getSUsers();
          }
        }, response => {
          this.$message.error(response.body);
        })
      },
      onDeleteAll() {
        this.$http.delete("/kc/simulation/users")
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data !== null) {
            this.getSUsers();
          }
        }, response => {
          this.$message.error(response.body);
        })
      },
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
