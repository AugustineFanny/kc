<template>
  <div>
    <el-form :inline="true" :model="formInline" class="demo-form-inline">
      <el-form-item label="">
        <el-input v-model="formInline.ident" style="width:300px;" placeholder="ID/用户名/邮箱/手机号"></el-input>
      </el-form-item>
      <el-form-item label="">
      <el-select v-model="formInline.group_id" clearable placeholder="分组">
        <el-option
          v-for="item in groupData"
          :key="item.id"
          :label="item.name"
          :value="item.id">
        </el-option>
      </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">查询</el-button>
      </el-form-item>
    </el-form>
    <el-row>
      <el-col :xs="24" :sm="24" :md="16" :lg="16">
        <el-table :data="tableData" border style="width: 100%; max-width:1371px;" @sort-change="sortChange">
          <el-table-column prop="id" label="ID" sortable="custom" width="100"></el-table-column>
          <el-table-column prop="username" fixed label="用户名" width="150"></el-table-column>
          <el-table-column prop="create_time" label="注册时间" width="180">
            <template slot-scope="scope">
              {{ utilHelper.timeShow(scope.row.create_time) }}
            </template>
          </el-table-column>
          <el-table-column prop="email" label="邮箱" width="200"></el-table-column>
          <el-table-column prop="mobile" label="手机号" width="150"></el-table-column>
          <el-table-column prop="role" label="身份认证" width="120">
            <template slot-scope="scope">{{roleName(scope.row.role)}}</template>
          </el-table-column>
          <el-table-column prop="kyc" label="KYC" width="120">
            <template slot-scope="scope">{{roleName(scope.row.kyc)}}</template>
          </el-table-column>
          <el-table-column label="查看" width="150">
            <template slot-scope="scope">
              <router-link :to="{ path: '/user/wallets/', query: {uid: scope.row.id} }">资金</router-link>
              <router-link v-if="scope.row.role != 0" :to="{ path: '/user/real-name/' + scope.row.id }">身份认证</router-link>
              <router-link v-if="scope.row.kyc != 0" :to="{ path: '/user/kyc/' + scope.row.id }">KYC</router-link>
            </template>
          </el-table-column>
          <el-table-column prop="fund" label="资金密码" width="100">
            <template slot-scope="scope">
              <el-tag :type="scope.row.fund ? 'primary' : 'success'" close-transition>{{scope.row.fund ? '已设置' : '未设置'}}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="tf_opened" label="双重认证" width="100">
            <template slot-scope="scope">
              <el-tag :type="scope.row.tf_opened ? 'primary' : 'success'" close-transition>{{scope.row.tf_opened ? '已设置' : '未设置'}}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="账户状态" width="130">
            <template slot-scope="scope">
              <el-tag :type="scope.row.status == 2 ? 'primary' : 'success'" close-transition>{{scope.row.status == 2 ? '冻结' : '正常'}}</el-tag>
              <el-button type="primary" size="mini" v-if="scope.row.status == 2" @click="onThaw(scope.row.username)">解冻</el-button>
            </template>
          </el-table-column>
          <el-table-column label="分组" width="100">
            <template slot-scope="scope">
              {{ groups[scope.row.group_id] }}
            </template>
          </el-table-column>
          <el-table-column prop="last_time" label="最后登录时间" sortable="custom" width="180">
            <template slot-scope="scope">
              {{ utilHelper.timeShow(scope.row.last_time) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100" fixed="right">
            <template slot-scope="scope">
              <el-button @click="handleClick(scope.row)" type="text" size="small">编辑</el-button>
              <el-button v-if="superAdmin == 1" @click="handleClickToRecharge(scope.row)" type="text" size="small">充币</el-button>
            </template>
          </el-table-column>
        </el-table>
        <div style="position:absolute;">
          <el-pagination layout="total, prev, pager, next, jumper" :page-size="formInline.page_size" :total="totalCount" @current-change="currentChange">
          </el-pagination>
        </div>
      </el-col>
      <el-col :xs="24" :sm="24" :md="8" :lg="8" v-show="show">
        <el-form ref="change" :model="form" :rules="rules" label-width="80px">
          <el-form-item label="用户名" prop="username">
            <el-input v-model="form.username" readonly auto-complete="off"></el-input>
          </el-form-item>
          <el-form-item label="分组" prop="group_id">
            <el-select v-model="form.group_id" placeholder="请选择">
              <el-option
                v-for="item in groupData"
                :key="item.id"
                :label="item.name"
                :value="item.id">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="onChange">提交</el-button>
            <el-button @click="onCancel">取消</el-button>
          </el-form-item>
        </el-form>
      </el-col>
      <el-col :xs="24" :sm="24" :md="8" :lg="8" v-show="showRecharge">
        <el-form ref="recharge" :model="form3" :rules="rules3" label-width="80px">
          <el-form-item label="用户名" prop="username">
            <el-input v-model="form3.username" readonly auto-complete="off"></el-input>
          </el-form-item>
          <el-form-item label="邮箱" prop="email">
            <el-input v-model="form3.email" readonly auto-complete="off"></el-input>
          </el-form-item>
          <el-form-item label="手机" prop="mobile">
            <el-input v-model="form3.mobile" readonly auto-complete="off"></el-input>
          </el-form-item>
          <el-form-item label="币种" prop="currency" required>
            <el-select v-model="form3.currency" placeholder="请选择">
              <el-option
                v-for="item in currencies"
                :key="item.currency"
                :label="item.currency"
                :value="item.currency">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="数量" prop="amount" required>
            <el-input-number controls-position="right" :min="0" v-model="form3.amount"></el-input-number>
          </el-form-item>
          <el-form-item label="密码" prop="password" required>
            <el-input type="password" auto-complete="off" v-model="form3.password"></el-input>
          </el-form-item>
          <el-form-item label="验证码" prop="captcha">
            <el-input v-model="form3.captcha">
              <el-button slot="append" @click="sendCaptcha">发送</el-button>
            </el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="onRecharge">提交</el-button>
            <el-button @click="onCancelRecharge">取消</el-button>
          </el-form-item>
        </el-form>
      </el-col>
    </el-row>
  </div>
</template>

<script>
  export default {
    data() {
      return {
        superAdmin: 0,
        formInline: {
          ident: '',
          group_id: null,
          _sort: '',
          _order: '',
          page: 1,
          page_size: 20,
        },
        totalCount: 0,
        tableData: [],
        groupData: [],
        groups: {},
        show: false,
        form: {
          id: '',
          username: '',
          group_id: null,
        },
        rules: {

        },
        showRecharge: false,
        currencies: [],
        form3: {
          uid: '',
          username: '',
          email: '',
          mobile: '',
          currency: '',
          amount: '',
          password: '',
          captcha: '',
        },
        rules3: {

        }
      };
    },
    created () {
      this.superAdmin = localStorage.super;
      this.getUsers();
      this.getGroups();
      if (this.superAdmin == 1)
        this.getCurrencies();
    },
    methods: {
      img(filename) {
        return '/uphp/name_auth/' + filename;
      },
      roleName(status) {
        if(status == 0) {
          return "未提交认证";
        } else if (status == 1) {
          return "待认证";
        } else if (status == 2) {
          return "认证未通过";
        } else if (status == 3) {
          return "认证通过";
        } else {
          return "未知状态：" + status;
        }
      },
      getUsers(params, pageNo) {
        this.$http.get("/kc/admin/users", {"params": this.formInline})
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data !== null) {
            this.formInline.page = data.page_no;
            this.formInline.page_size = data.page_size;
            this.totalCount = data.total_count;
            this.tableData = data.list;
          }
        }, response => {
          this.$message.error(response.body);
        });
      },
      getGroups() {
        this.$http.get("/kc/admin/groups")
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data !== null) {
            this.groupData = data;
            this.groupData.push({id: 0, name: "default"});
            for(var i in data) {
              var row = data[i];
              this.groups[row.id] = row.name;
            }
            this.groups[0] = "default";
          }
        }, response => {
          this.$message.error(response.body);
        });
      },
      currentChange(pageNo) {
        this.formInline.page = pageNo;
        this.getUsers();
      },
      onSubmit() {
        this.formInline.page = 1;
        this.getUsers();
      },
      onThaw(username) {
        this.$confirm('解除该账户冻结状态?', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          this.$http.post("/kc/admin/user/thaw", {username:username})
          .then(response => {
            this.utilHelper.handle(this, response)
          }, response => {
            this.$message.error(response.body);
          })
          this.getUsers();
        });
      },
      sortChange(a) {
        this.formInline._sort = a.prop;
        this.formInline._order = a.order;
        this.formInline.page = 1;
        this.getUsers();
      },
      handleClick(row) {
        this.showRecharge = false;
        this.show = true;
        this.form.id = row.id;
        this.form.username = row.username;
        this.form.group_id = row.group_id;
      },
      onChange() {
        this.$refs["change"].validate((valid) => {
          if (valid) {
            this.$resource("/kc/admin/super/user/{id}")
            .save({id: this.form.id}, {
              group_id: parseInt(this.form.group_id)
            })
            .then(response => {
              var data = this.utilHelper.handle(this, response);
              if(data !== null) {
                this.$message.success("操作成功");
                this.getUsers();
              };
            }, response => {
              this.$message.error(response.body);
            });
          } else {
            return false;
          }
        });
      },
      onCancel() {
        this.show = false;
      },
      getCurrencies() {
        this.$http.get("/kc/admin/super/currencies")
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data != null) {
            for(var index in data) {
              if(data[index].currency != "BTC" && data[index].currency != "ETH")
                this.currencies.push(data[index]);
            }
          }
        }, response => {
          this.$message.error(response.body);
        });
      },
      handleClickToRecharge(row) {
        this.show = false;
        this.showRecharge = true;
        this.form3.id = row.id;
        this.form3.username = row.username;
        this.form3.email = row.email;
        this.form3.mobile = row.mobile;
      },
      onRecharge() {
        this.$refs["recharge"].validate((valid) => {
          if (valid) {
            this.$resource("/kc/admin/super/user/{id}/recharge")
            .save({id: this.form3.id}, {
              currency: this.form3.currency,
              amount: this.form3.amount,
              password: this.form3.password,
              captcha: this.form3.captcha,
            })
            .then(response => {
              var data = this.utilHelper.handle(this, response);
              if(data !== null) {
                this.$message.success("操作成功");
                this.showRecharge = false;
                this.$refs["recharge"].resetFields();
              };
            }, response => {
              this.$message.error(response.body);
            });
          } else {
            return false;
          }
        });
      },
      onCancelRecharge() {
        this.showRecharge= false;
      },
      sendCaptcha() {
        this.$http.get("/kc/admin/super/captcha", {"params": {"type": "userrecharge"}})
        .then(response => {
          var data = this.utilHelper.handle(this, response);
        }, response => {
          this.$message.error(response.body);
        });
      }
    }
  };
</script>
