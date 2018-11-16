<template>
  <div>
      <el-menu theme="dark" class="el-menu-demo" mode="horizontal" background-color="#545c64" text-color="#fff" active-text-color="#ffd04b">
        <el-menu-item index="1" style="font-size:25px;">后台管理系统</el-menu-item>
        <el-submenu index="2" style="float:right;">
          <template slot="title">{{ user }}</template>
          <el-menu-item index="2-1" @click="dialogFormVisible = true">修改密码</el-menu-item>
          <el-menu-item index="2-2" v-if="superAdmin == 1" @click="setSuperPassword = true">设置超级密码</el-menu-item>
          <el-menu-item index="2-3" @click="setMobile = true">绑定手机号</el-menu-item>
          <el-menu-item index="2-4" @click="logout">退出</el-menu-item>
        </el-submenu>
      </el-menu>
      <el-col class="sidebar">
        <el-menu router style="height: 100%;" class="el-menu-vertical-demo">
          <el-submenu index="1">
            <template slot="title">
              <i class="el-icon-message"></i>
              <span slot="title">用户管理</span>
            </template>
              <el-menu-item index="/user">用户</el-menu-item>
              <el-menu-item index="/user/wallets">资金</el-menu-item>
              <el-menu-item index="/user/addresses">地址</el-menu-item>
              <el-menu-item index="/user/real-names">实名认证</el-menu-item>
              <el-menu-item index="/user/kycs">KYC</el-menu-item>
              <el-menu-item index="/user/transfers">充提记录</el-menu-item>
              <el-menu-item index="/user/transfer-outs">提现审核</el-menu-item>
              <el-menu-item index="/user/fund-changes">资金变动</el-menu-item>
          </el-submenu>
          <el-submenu index="2">
            <template slot="title">
              <i class="el-icon-message"></i>
              <span slot="title">管理</span>
            </template>
              <el-menu-item index="/manage/subscription">认购记录</el-menu-item>
              <el-menu-item index="/manage/simulation">模拟</el-menu-item>
              <el-menu-item index="/manage/profit-date">收益(日)</el-menu-item>
              <el-menu-item index="/manage/profit-month">收益(月)</el-menu-item>
              <el-menu-item index="/manage/predistribution">未分配</el-menu-item>
          </el-submenu>
          <el-submenu index="3">
            <template slot="title">
              <i class="el-icon-menu"></i>
              <span slot="title">统计</span>
            </template>
              <el-menu-item index="/stat/statistics">总览</el-menu-item>
              <el-menu-item index="/stat/invites">邀请</el-menu-item>
              <el-menu-item index="/stat/child/amounts">下级统计</el-menu-item>
          </el-submenu>
          <el-submenu index="5" v-if="superAdmin == 1">
            <template slot="title">
              <i class="el-icon-message"></i>
              <span slot="title">后台管理</span>
            </template>
              <el-menu-item index="/admin/admins">管理员</el-menu-item>
              <el-menu-item index="/admin/coins">币种</el-menu-item>
              <el-menu-item index="/admin/groups">分组</el-menu-item>
              <el-menu-item index="/admin/batch-sms">短信发送</el-menu-item>
              <el-menu-item index="/admin/batch-email">邮箱发送</el-menu-item>
              <el-menu-item index="/admin/address-pool">地址池</el-menu-item>
              <el-menu-item index="/admin/logs">日志</el-menu-item>
              <el-menu-item index="/admin/ledgerETH">LedgerETH</el-menu-item>
          </el-submenu>
          <el-menu-item></el-menu-item>
        </el-menu>
      </el-col>
      <el-col class="main-container">
          <router-view></router-view>
      </el-col>
    <el-dialog title="修改密码" :visible.sync="dialogFormVisible">
      <el-form :model="form" :rules="rules" ref="form" label-width="100px" class="demo-ruleForm">
        <el-form-item label="当前密码" prop="password">
          <el-input type="password" v-model="form.password" auto-complete="off"></el-input>
        </el-form-item>
        <el-form-item label="新密码" prop="new_password">
          <el-input type="password" v-model="form.new_password" auto-complete="off"></el-input>
        </el-form-item>
        <el-form-item label="确认密码" prop="check_password">
          <el-input type="password" v-model="form.check_password" auto-complete="off"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">取 消</el-button>
        <el-button type="primary" @click="submitForm('form')">确 定</el-button>
      </div>
    </el-dialog>
    <el-dialog title="设置超级密码" :visible.sync="setSuperPassword">
      <el-form :model="form2" :rules="rules2" ref="form2" label-width="100px" class="demo-ruleForm">
        <el-form-item label="验证码" prop="captcha" required>
          <el-input v-model="form2.captcha" auto-complete="off">
            <el-button slot="append" @click="sendCaptcha">发送</el-button>
          </el-input>
        </el-form-item>
        <el-form-item label="新密码" prop="password" required>
          <el-input type="password" v-model="form2.password" auto-complete="off"></el-input>
        </el-form-item>
        <el-form-item label="确认密码" prop="check_password" required>
          <el-input type="password" v-model="form2.check_password" auto-complete="off"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="setSuperPassword = false">取 消</el-button>
        <el-button type="primary" @click="submitSuperPasswordForm('form2')">确 定</el-button>
      </div>
    </el-dialog>
    <el-dialog title="绑定手机号" :visible.sync="setMobile">
      <el-form :model="form3" :rules="rules3" ref="form3" label-width="100px" class="demo-ruleForm">
        <el-form-item label="手机号" prop="mobile" required>
          <el-input v-model="form3.mobile" auto-complete="off"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="setMobile = false">取 消</el-button>
        <el-button type="primary" @click="submitMobileForm('form3')">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
  import router from '../router'

  export default {
    data() {
      var validatePass = (rule, value, callback) => {
        if (value === '') {
          callback(new Error('请输入密码'));
        } else {
          if (this.form.check_password !== '') {
            this.$refs.form.validateField('check_password');
          }
          callback();
        }
      };
      var validatePass2 = (rule, value, callback) => {
        if (value === '') {
          callback(new Error('请再次输入密码'));
        } else if (value !== this.form.new_password) {
          callback(new Error('两次输入密码不一致!'));
        } else {
          callback();
        }
      };
      var validatePassword = (rule, value, callback) => {
        if (value === '') {
          callback(new Error('请输入密码'));
        } else {
          if (this.form2.check_password !== '') {
            this.$refs.form2.validateField('check_password');
          }
          callback();
        }
      };
      var validatePassword2 = (rule, value, callback) => {
        if (value === '') {
          callback(new Error('请再次输入密码'));
        } else if (value !== this.form2.password) {
          callback(new Error('两次输入密码不一致!'));
        } else {
          callback();
        }
      };
      var validateMobile = (rule, value, callback) => {
        if (value === '') {
          callback(new Error('请输入手机号'));
        } else if (value.length != 11 || value[0] != '1') {
          callback(new Error('请输入正确的手机号'));
        } else {
          callback();
        }
      }
      return {
        user: "",
        superAdmin: 0,
        leftSpan: 4,
        rightSpan: 19,
        dialogFormVisible: false,
        setSuperPassword: false,
        setMobile: false,
        form: {
          password: '',
          new_password: '',
          check_password: ''
        },
        rules: {
          new_password: [
            { validator: validatePass, trigger: 'blur' }
          ],
          check_password: [
            { validator: validatePass2, trigger: 'blur' }
          ]
        },
        formLabelWidth: '120px',
        form2: {
          captcha: '',
          password: '',
          check_password: '',
        },
        rules2: {
          password: [
            { validator: validatePassword, trigger: 'blur' }
          ],
          check_password: [
            { validator: validatePassword2, trigger: 'blur' }
          ]
        },
        form3: {
          mobile: '',
        },
        rules3: {
          mobile: [
            { validator: validateMobile, trigger: 'blur' }
          ],
        },
      };
    },
    created() {
      this.user = localStorage.user;
      this.form3.mobile = localStorage.mobile;
      this.superAdmin = localStorage.super;
    },
    methods: {
      submitForm(formName) {
        this.$refs[formName].validate((valid) => {
          if (valid) {
            this.$http.post("/kc/admin/change-password", this.form)
            .then(response => {
              var data = this.utilHelper.handle(this, response);
              if(data !== null) {
                this.$refs[formName].resetFields();
                this.$message.success("修改成功,请重新登录");
                this.dialogFormVisible = false;
                this.logout();
              };
            }, response => {
              this.$message.error(response.body);
            });
          } else {
            return false;
          }
        });
      },
      submitSuperPasswordForm(formName) {
        this.$refs[formName].validate((valid) => {
          if (valid) {
            this.$http.post("/kc/admin/super/super-password", this.form2)
            .then(response => {
              var data = this.utilHelper.handle(this, response);
              if(data !== null) {
                this.$refs[formName].resetFields();
                this.$message.success("修改成功");
                this.setSuperPassword = false;
              };
            }, response => {
              this.$message.error(response.body);
            });
          } else {
            return false;
          }
        });
      },
      sendCaptcha() {
        this.$http.get("/kc/admin/super/captcha", {"params": {"type": "superpassword"}})
        .then(response => {
          var data = this.utilHelper.handle(this, response);
        }, response => {
          this.$message.error(response.body);
        });
      },
      submitMobileForm(formName) {
        this.$refs[formName].validate((valid) => {
          if (valid) {
            this.$http.post("/kc/admin/set-mobile", this.form3)
            .then(response => {
              var data = this.utilHelper.handle(this, response);
              if(data !== null) {
                this.$refs[formName].resetFields();
                this.$message.success("修改成功");
                this.setMobile = false;
              };
            }, response => {
              this.$message.error(response.body);
            });
          } else {
            return false;
          }
        });
      },
      logout() {
        this.$http.get("/kc/admin/logout").then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data !== null) {
            router.push('/login');
          }
        }, response => {
          this.$message.error(response.body);
        });
      }
    }
  }
</script>

<style>
  body {
    margin: 0;
  }
  .el-menu-demo {
    position: fixed;
    left: 0;
    right: 0;
    z-index: 1000;
  }
  .el-menu-vertical-demo {
    width: 200px;
    overflow: auto;
  }
  .sidebar {
    height: 100%;
    margin-top: 60px;
    position: fixed;
  }
  .main-container {
    padding: 110px 50px 50px 250px;
  }
  .app-main {
    padding-left: 50px;
  }

</style>
