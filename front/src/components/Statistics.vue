<template>
  <div>
    <el-row :gutter="20">
      <el-col :xs="24" :sm="24" :md="12" :lg="8" :xl="6">
        <el-card class="box-card">
          <div slot="header" class="clearfix">
            <span style="line-height: 36px;">当前量</span>
          </div>
          <div v-for="o in balanceData" :key="o" class="text item">
            {{ o.balance }} {{ o.currency }}
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="24" :md="12" :lg="8" :xl="6">
        <el-card class="box-card">
          <div slot="header" class="clearfix">
            <span style="line-height: 36px;">累计充值</span>
          </div>
          <div v-for="o in rechargeData" :key="o" class="text item">
            {{ o.balance }} {{ o.currency }}
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="24" :md="12" :lg="8" :xl="6">
        <el-card class="box-card">
          <div slot="header" class="clearfix">
            <span style="line-height: 36px;">累计提币</span>
          </div>
          <div v-for="o in withdrawData" :key="o" class="text item">
            {{ o.balance }} {{ o.currency }}
          </div>
        </el-card>
      </el-col>
    </el-row>
    <el-row :gutter="20">
      <el-col :xs="24" :sm="24" :md="12" :lg="8" :xl="6">
        <el-card class="box-card">
          <div slot="header" class="clearfix">
            <span style="line-height: 36px;">注册</span>
          </div>
          <div class="text item">今日注册：{{ signup.signUpToday }}</div>
          <div class="text item">昨日注册：{{ signup.signUpYesterday }}</div>
          <div class="text item">本月注册：{{ signup.signUpMonth }}</div>
          <div class="text item">上月注册：{{ signup.signUpPrevMonth }}</div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="24" :md="12" :lg="8" :xl="6">
        <el-card class="box-card">
          <div slot="header" class="clearfix">
            <span style="line-height: 36px;">FET</span>
          </div>
          <div class="text item">总量：{{ fet.amounts }}</div>
          <div class="text item">冻结：{{ fet.miningAmounts }}</div>
          <div class="text item">流通：{{ fet.activeAmounts }}</div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
  export default {
    data() {
      return {
        balanceData: [],
        rechargeData: [],
        withdrawData: [],
        signup: {
          signUpToday: "-",
          signUpYesterday: "-",
          signUpMonth: "-",
          signUpPrevMonth: "-",
        },
        fet: {
          amounts: "-",
          miningAmounts: "-",
          activeAmounts: "-",
        },
      };
    },
    created () {
      this.getStatistics();
    },
    methods: {
      getStatistics(params) {
        this.$http.get("/kc/admin/statistics")
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data !== null) {
            this.balanceData = data.balance;
            this.rechargeData = data.recharge;
            this.withdrawData = data.withdraw;
            this.signup.signUpToday = data.sign_up_today[0].num;
            this.signup.signUpYesterday = data.sign_up_yesterday[0].num;
            this.signup.signUpMonth = data.sign_up_month[0].num;
            this.signup.signUpPrevMonth = data.sign_up_prev_month[0].num;
            this.fet.amounts = data.fet_stat[0].amounts;
            this.fet.miningAmounts = data.fet_stat[0].mining_amounts;
            this.fet.activeAmounts = this.fet.amounts - this.fet.miningAmounts;
          }
        }, response => {
          this.$message.error(response.body);
        });
      }
    }
  };
</script>
