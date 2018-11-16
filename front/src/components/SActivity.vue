<template>
  <div>
    <el-row>
      <el-col :xs="24" :sm="24" :md="16" :lg="16">
        <el-button type="primary" @click="showCreate = true">添加</el-button>
        <el-upload style="display: inline-block" action="/kc/simulation/batch/activities" :show-file-list="false" accept=".csv" :on-success="onSuccess">
          <el-button type="primary">批量添加</el-button>
        </el-upload>
        <el-button type="primary" @click="onDeleteAll">全部删除</el-button>
        <router-link :to="{ path: '/manage/simulation'}"><el-button type="primary" style="float: right">返回</el-button></router-link>
      </el-col>
    </el-row>
    <el-row>
      <el-col :xs="24" :sm="24" :md="16" :lg="16">
        <el-table :data="tableData" border style="width: 100%;">
          <el-table-column prop="username" label="label" width="150" sortable></el-table-column>
          <el-table-column prop="subscription" label="认购" width="150"></el-table-column>
          <el-table-column prop="lock" label="转移" width="150"></el-table-column>
          <el-table-column prop="in" label="转入" width="150"></el-table-column>
          <el-table-column prop="in_source" label="转入来源用户" width="150"></el-table-column>
          <el-table-column prop="out" label="转出" width="150"></el-table-column>
          <el-table-column prop="out_dest" label="转出目标用户" width="150"></el-table-column>
          <el-table-column prop="date" label="时间" sortable width="180">
            <template slot-scope="scope">
              {{ utilHelper.timeShow(scope.row.date) }}
            </template>
          </el-table-column>
          <el-table-column prop="expire_dates" label="解锁日期" width="900"></el-table-column>
          <el-table-column label="操作" width="100" fixed="right">
            <template slot-scope="scope">
              <el-button @click="handleClick(scope.row)" type="text" size="small">编辑</el-button>
              <el-button @click="onDelete(scope.row.id)" type="text" size="small">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-col>
      <el-col :xs="24" :sm="24" :md="8" :lg="8" v-show="show">
        <el-form ref="change" :model="form" :rules="rules" label-width="120px">
          <el-form-item label="用户" prop="username">
            <el-input controls-position="right" v-model="form.username" readonly></el-input>
          </el-form-item>
          <el-form-item label="认购" prop="subscription" required>
            <el-input-number controls-position="right" v-model="form.subscription" :min="0" :debounce="500"></el-input-number>
          </el-form-item>
          <el-form-item label="转移" prop="lock" required>
            <el-input-number controls-position="right" v-model="form.lock" :min="0" :debounce="500"></el-input-number>
          </el-form-item>
          <el-form-item label="转入" prop="in" required>
            <el-input-number controls-position="right" v-model="form.in" :min="0" :debounce="500"></el-input-number>
            <el-select v-model="form.in_source" filterable placeholder="来源用户" clearable>
            <el-option
              v-for="item in users"
              :key="item.username"
              :label="item.username"
              :value="item.username">
            </el-option>
          </el-select>
          </el-form-item>
          <el-form-item label="转出" prop="out" required>
            <el-input-number controls-position="right" v-model="form.out" :min="0" :debounce="500"></el-input-number>
            <el-select v-model="form.out_dest" filterable placeholder="目标用户" clearable>
            <el-option
              v-for="item in users"
              :key="item.username"
              :label="item.username"
              :value="item.username">
            </el-option>
          </el-select>
          </el-form-item>
          <el-form-item label="日期" prop="date" required>
            <el-date-picker v-model="form.date" type="date" placeholder="选择日期" value-format="yyyy-MM-dd" required></el-date-picker>
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
        <el-form-item label="用户" prop="username" required>
          <el-select v-model="form2.username" filterable placeholder="请选择">
            <el-option
              v-for="item in users"
              :key="item.username"
              :label="item.username"
              :value="item.username">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="认购" prop="subscription" required>
          <el-input-number controls-position="right" v-model="form2.subscription" :min="0" :debounce="500"></el-input-number>
        </el-form-item>
        <el-form-item label="转移" prop="lock" required>
          <el-input-number controls-position="right" v-model="form2.lock" :min="0" :debounce="500"></el-input-number>
        </el-form-item>
        <el-form-item label="转入" prop="in" required>
          <el-input-number controls-position="right" v-model="form2.in" :min="0" :debounce="500"></el-input-number> 
          <el-select v-model="form2.in_source" filterable placeholder="来源用户" clearable>
            <el-option
              v-for="item in users"
              :key="item.username"
              :label="item.username"
              :value="item.username">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="转出" prop="out" required>
          <el-input-number controls-position="right" v-model="form2.out" :min="0" :debounce="500"></el-input-number>
          <el-select v-model="form2.out_dest" filterable placeholder="目标用户" clearable>
            <el-option
              v-for="item in users"
              :key="item.username"
              :label="item.username"
              :value="item.username">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="日期" prop="date" required>
          <el-date-picker v-model="form2.date" type="date" placeholder="选择日期" value-format="yyyy-MM-dd" required></el-date-picker>
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
        users: [],
        show: false,
        form: {
          username: "",
          subscription: 0,
          lock: 0,
          in: 0,
          in_source: "",
          out: 0,
          out_dest: "",
          date: "",
        },
        rules: {

        },
        showCreate: false,
        form2: {
          username: "",
          subscription: 0,
          lock: 0,
          in: 0,
          in_source: "",
          out: 0,
          out_dest: "",
          date: "",
        },
        rules2: {

        }
      };
    },
    created () {
      this.getSActivities();
      this.getSUsers();
    },
    methods: {
      getSUsers() {
        this.$http.get("/kc/simulation/users")
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data != null) {
            this.users = data;
          }
        }, response => {
          this.$message.error(response.body);
        });
      },
      getSActivities() {
        this.$http.get("/kc/simulation/activities")
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
            this.$resource("/kc/simulation/activity/{id}")
            .save({id: this.form.id}, {
              subscription: this.form.subscription,
              lock: this.form.lock,
              in: this.form.in,
              in_source: this.form.in_source,
              out: this.form.out,
              out_dest: this.form.out_dest,
              date: this.form.date.slice(0, 10),
            })
            .then(response => {
              var data = this.utilHelper.handle(this, response);
              if(data !== null) {
                this.$message.success("操作成功");
                this.$refs["change"].resetFields();
                this.show = false;
                this.getSActivities();
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
            this.$http.post("/kc/simulation/activities", {
              username: this.form2.username,
              subscription: this.form2.subscription,
              lock: this.form2.lock,
              in: this.form2.in,
              in_source: this.form2.in_source,
              out: this.form2.out,
              out_dest: this.form2.out_dest,
              date: this.form2.date.slice(0, 10),
            })
            .then(response => {
              var data = this.utilHelper.handle(this, response);
              if(data !== null) {
                this.$message.success("操作成功");
                this.$refs["create"].resetFields();
                this.showCreate = false;
                this.getSActivities();
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
          this.getSActivities();
        } else {
          this.$message.error(response.msg);
          this.getSActivities();
        }
      },
      onCancel() {
        this.show = false;
      },
      onDelete(id) {
        this.$http.delete("/kc/simulation/activity/" + id)
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data !== null) {
            this.getSActivities();
          }
        }, response => {
          this.$message.error(response.body);
        })
      },
      onDeleteAll() {
        this.$http.delete("/kc/simulation/activities")
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data !== null) {
            this.getSActivities();
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
