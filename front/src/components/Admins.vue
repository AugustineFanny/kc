<template>
  <div>
    <el-row><el-button type="primary" @click="showCreate = true">创建管理员</el-button></el-row>
    <el-row>
      <el-col :xs="24" :sm="24" :md="16" :lg="16">
      	<el-table :data="tableData" border style="width: 100%; max-width:1361px;">
          <el-table-column prop="id" label="ID" sortable width="100"></el-table-column>
    	    <el-table-column prop="username" label="用户名" width="200"></el-table-column>
          <el-table-column prop="mobile" label="手机号" width="120"></el-table-column>
          <el-table-column prop="super" label="超级" sortable width="100">
            <template slot-scope="scope">
              <el-tag :type="scope.row.super ? 'primary' : 'success'" close-transition>{{scope.row.super ? '是' : '否'}}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" sortable width="100">
            <template slot-scope="scope">
              <el-tag :type="scope.row.status ? 'primary' : 'success'" close-transition>{{scope.row.status ? '禁用' : '启用'}}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="create_time" label="注册时间" width="180">
            <template slot-scope="scope">
              {{ utilHelper.timeShow(scope.row.create_time) }}
            </template>
          </el-table-column>
          <el-table-column prop="last_time" label="最后登录时间" width="180">
            <template slot-scope="scope">
              {{ utilHelper.timeShow(scope.row.last_time) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100" fixed="right">
            <template slot-scope="scope">
              <el-button @click="handleClick(scope.row)" type="text" size="small">编辑</el-button>
            </template>
          </el-table-column>
    	  </el-table>
      </el-col>
      <el-col :xs="24" :sm="24" :md="8" :lg="8" v-show="showChange">
        <el-form ref="change" :model="form" :rules="rules" label-width="80px">
          <el-form-item label="ID" prop="id">
            <el-input readonly v-model="form.id"></el-input>
          </el-form-item>
          <el-form-item label="用户名" prop="username">
            <el-input v-model="form.username"></el-input>
          </el-form-item>
          <el-form-item label="状态" prop="status">
            <el-select v-model="form.status" placeholder="请选择">
              <el-option key="0" label="启用" value="0" ></el-option>
              <el-option key="1" label="禁用" value="1" ></el-option>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="onChange">提交</el-button>
            <el-button @click="showChange = false">取消</el-button>
          </el-form-item>
        </el-form>
      </el-col>
    </el-row>
    <el-dialog title="创建管理员" :visible.sync="showCreate">
      <el-form ref="create" :model="form2" :rules="rules2" label-width="100px" class="demo-ruleForm">
        <el-form-item label="用户名" prop="ident">
          <el-input v-model="form2.ident" auto-complete="off"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form2.password" auto-complete="off"></el-input>
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
  import utils from '../utils'

  export default {
    data() {
      return {
        showCreate: false,
        showChange: false,
        tableData: [],
        form: {
          id: '',
          username: '',
          status: '',
        },
        rules: {
          username: [
            { required: true, message: "请输入用户名", trigger: 'change'}
          ]
        },
        form2: {
          ident: '',
          password: '',
        },
        rules2: {
          ident: [
            { required: true, message: "请输入用户名", trigger: 'change'}
          ],
          password: [
            { min:6, required: true, message: "请输入大于6位密码", trigger: 'change'}
          ]
        }
      };
    },
    created () {
      this.getAdmins();
    },
    methods: {
      getAdmins() {
        this.$http.get("/kc/admin/super/admins")
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
        this.showChange = true;
        for (var i in row) {
            this.form[i] = row[i];
            this.form.status = this.form.status.toString();
        };
      },
      onChange() {
        this.$refs["change"].validate((valid) => {
          if (valid) {
            this.$resource("/kc/admin/super/admin/{id}")
            .save({id: this.form.id}, {username: this.form.username, status: parseInt(this.form.status)})
            .then(response => {
              var data = this.utilHelper.handle(this, response);
              if(data !== null) {
                this.$message.success("操作成功");
                this.getAdmins();
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
            this.$http.post("/kc/admin/super/admins", this.form2)
            .then(response => {
              var data = this.utilHelper.handle(this, response);
              if(data !== null) {
                this.$message.success("操作成功");
                this.$refs["create"].resetFields();
                this.showCreate = false;
                this.getAdmins();
              };
            }, response => {
              this.$message.error(response.body);
            });
          } else {
            return false;
          }
        });
      },
    }
  }
</script>

<style>
  .el-row {
    margin-bottom: 20px;
    &:last-child {
      margin-bottom: 0;
    }
  }
</style>
