<template>
  <div>
    <el-row>
      <el-col :xs="24" :sm="24" :md="24" :lg="24" style="margin-bottom: 10px">
        <el-form ref="change" :inline="true" label-width="120px">
          <el-form-item label="锁仓利率" prop="miningInterestRate">
            <el-input-number controls-position="right" v-model="miningInterestRate" :min="0" :max="0.005" :debounce="500"></el-input-number>
          </el-form-item>
          <el-form-item label="推广利率" prop="competitionInterestRate">
            <el-input-number controls-position="right" v-model="competitionInterestRate" :min="0" :max="0.005" :debounce="500"></el-input-number>
          </el-form-item>
          <el-button type="primary" @click="onInterestRate">提交</el-button>
        </el-form>
      </el-col>
      <el-col :xs="24" :sm="24" :md="24" :lg="24">
        <router-link :to="{ path: '/simulation/user'}"><el-button type="primary">用户管理</el-button></router-link>
        <router-link :to="{ path: '/simulation/activity'}"><el-button type="primary">活动管理</el-button></router-link>
        <el-button type="primary" @click="onStat" style="float: right">统计</el-button>
        <el-button type="primary" @click="onCalculation" style="float: right;margin-right: 20px">模拟</el-button>
        <el-date-picker v-model="date" type="daterange" start-placeholder="开始日期" end-placeholder="结束日期" value-format="yyyy-MM-dd" style="float: right;margin-right: 20px"></el-date-picker>
        <el-select v-model="selectedUser" filterable clearable placeholder="用户名" style="float: right;margin-right: 20px">
          <el-option v-for="item in users" :key="item" :label="item" :value="item"></el-option>
        </el-select>
      </el-col>
    </el-row>
    <el-row>
      <el-col :xs="24" :sm="24" :md="24" :lg="24">
        <el-table :data="tableData" border show-summary @row-click="rowClick" style="width: 100%;">
          <el-table-column prop="username" label="label" width="100" fixed="left"></el-table-column>
          <el-table-column prop="subscription_all" label="认购总计" sortable width="150"></el-table-column>
          <el-table-column prop="mining_all" label="生息总计" sortable width="150"></el-table-column>
          <el-table-column prop="share_all" label="奖励总计" sortable width="150"></el-table-column>
          <el-table-column prop="unlock_all" label="释放总计" sortable width="150"></el-table-column>
          <el-table-column prop="lock" label="锁定总计" sortable width="150"></el-table-column>
          <el-table-column prop="active" label="流通总计" sortable width="150"></el-table-column>
          <el-table-column label="总资产" sortable width="150">
            <template slot-scope="scope">
              {{ utilHelper.amountShow(scope.row.lock + scope.row.active) }}
            </template>
          </el-table-column>
          <el-table-column prop="subscription_day" label="当日认购" sortable width="150"></el-table-column>
          <el-table-column prop="mining_day" label="当日生息" sortable width="150"></el-table-column>
          <el-table-column prop="share_day" label="当日奖励" sortable width="150"></el-table-column>
          <el-table-column prop="unlock_day" label="当日释放" sortable width="150"></el-table-column>
          <el-table-column prop="lock_day" label="新增锁定" sortable width="150"></el-table-column>
          <el-table-column prop="active_day" label="新增流通" sortable width="150"></el-table-column>
          <el-table-column prop="in_all" label="转入总计" sortable width="150"></el-table-column>
          <el-table-column prop="out_all" label="转出总计" sortable width="150"></el-table-column>
          <el-table-column prop="addlock_all" label="转移总计" sortable width="150"></el-table-column>
          <el-table-column prop="in_day" label="当日转入" sortable width="150"></el-table-column>
          <el-table-column prop="out_day" label=当日转出 sortable width="150"></el-table-column>
          <el-table-column prop="addlock_day" label="当日转移" sortable width="150"></el-table-column>
          <el-table-column prop="date" label="日期" sortable width="150" fixed="right"></el-table-column>
        </el-table>
      </el-col>
    </el-row>
    <el-dialog title="详情" :visible.sync="dialogTableVisible">
      <div style="margin: 20px;font-size: 25px"><span>用户：{{ dialogTable.username }}</span><span style="padding-left: 100px">日期：{{ dialogTable.date }}</span></div>
      <table class="dialog-table">
        <tr>
          <td>认购总计：</td><td class="dialog-table-th">{{ dialogTable.subscription_all }}</td>
          <td>当日认购：</td><td class="dialog-table-th">{{ dialogTable.subscription_day }}</td>
          <td>转入总计：</td><td class="dialog-table-th">{{ dialogTable.in_all }}</td>
          <td>锁定总计：</td><td class="dialog-table-th">{{ dialogTable.lock }}</td>
        </tr>
        <tr>
          <td>生息总计：</td><td class="dialog-table-th">{{ dialogTable.mining_all }}</td>
          <td>当日生息：</td><td class="dialog-table-th">{{ dialogTable.mining_day }}</td>
          <td>转出总计：</td><td class="dialog-table-th">{{ dialogTable.out_all }}</td>
          <td>流通总计：</td><td class="dialog-table-th">{{ dialogTable.active }}</td>
        </tr>
        <tr>
          <td>奖励总计：</td><td class="dialog-table-th">{{ dialogTable.share_all }}</td>
          <td>当日奖励：</td><td class="dialog-table-th">{{ dialogTable.share_day }}</td>
          <td>转移总计：</td><td class="dialog-table-th">{{ dialogTable.addlock_all }}</td>
          <td>新增锁定：</td><td class="dialog-table-th">{{ dialogTable.lock_day }}</td>
        </tr>
        <tr>
          <td>释放总计：</td><td class="dialog-table-th">{{ dialogTable.unlock_all }}</td>
          <td>当日释放：</td><td class="dialog-table-th">{{ dialogTable.unlock_day }}</td>
          <td>当日转入：</td><td class="dialog-table-th">{{ dialogTable.in_day }}</td>
          <td>新增流通：</td><td class="dialog-table-th">{{ dialogTable.active_day }}</td>
        </tr>
        <tr>
          <td>个人总计：</td><td class="dialog-table-th">{{ dialogTable.reward1_all }}</td>
          <td>当日个人：</td><td class="dialog-table-th">{{ dialogTable.reward1_day }}</td>
          <td>当日转出：</td><td class="dialog-table-th">{{ dialogTable.out_day }}</td>
          <td>总资产：</td><td class="dialog-table-th">{{ utilHelper.amountShow(dialogTable.lock + dialogTable.active) }}</td>
        </tr>
        <tr>
          <td>竞赛总计：</td><td class="dialog-table-th">{{ dialogTable.reward2_all }}</td>
          <td>当日竞赛：</td><td class="dialog-table-th">{{ dialogTable.reward2_day }}</td>
          <td>当日转移：</td><td class="dialog-table-th">{{ dialogTable.addlock_day }}</td>
          <td>流通比率：</td><td class="dialog-table-th">{{ utilHelper.amountShow(dialogTable.active / (dialogTable.lock + dialogTable.active)) }}%</td>
        </tr>
      </table>
    </el-dialog>
  </div>
</template>

<style>
  .el-dialog {
    width: 1200px;
  }
  .dialog-table .dialog-table-th {
    width: 200px;
  }
  .el-dialog__body {
    line-height: 30px;
  }
</style>

<script>
  export default {
    data() {
      return {
        dialogTableVisible: false,
        tableData: [],
        tableDataBack: [],
        dialogTable: {},
        selectedUser: "",
        users: [],
        miningInterestRate: 0,
        competitionInterestRate: 0,
        date: ["2018-06-01", "2018-06-30"],
        show: false,
        form: {
          username: "",
          mining_amount: 0,
          inviter_id: "",
        },
        rules: {

        },
        showCreate: false,
        form2: {
          username: "",
          mining_amount: 0,
          inviter_id: "",
        },
        rules2: {

        }
      };
    },
    created () {
      this.getSCurrency();
      this.getSUsers();
    },
    methods: {
      getSCurrency() {
        this.$http.get("/kc/simulation/currency")
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data != null) {
            this.miningInterestRate = data.mining_interest_rate;
            this.competitionInterestRate = data.competition_interest_rate;
          }
        })
      },
      getSUsers() {
        this.$http.get("/kc/simulation/users")
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data != null) {
            var users = [];
            for(var i in data) {
              users.push(data[i].username);
            }
            this.users = users;
          }
        }, response => {
          this.$message.error(response.body);
        });
      },
      onCalculation() {
        this.tableData = [];
        this.$http.get("/kc/simulation/calculation", {params: {username: this.selectedUser, start_date: this.date[0], end_date: this.date[1]}})
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          var tableData = [];
          for (var key in response.data.data) {
            var row = data[key];
            tableData.push(row);
          this.tableData = tableData;
          this.tableDataBack = tableData;
          }
        }, response => {
          this.$message.error(response.body);
        })
      },
      onStat() {
        this.tableData = [];
        this.$http.get("/kc/simulation/stat", {params: {start_date: this.date[0], end_date: this.date[1]}})
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          var tableData = [];
          for (var key in response.data.data) {
            var row = data[key];
            tableData.push(row);
          this.tableData = tableData;
          this.tableDataBack = tableData;
          }
        }, response => {
          this.$message.error(response.body);
        })
      },
      onInterestRate() {
        if(this.reward < 0 || this.reward > 0.005) {
          this.$message.error("锁仓利率必须 >0 <0.005");
          return
        }
        if(this.spread < 0 || this.spread > 0.005) {
          this.$message.error("推广利率必须 >0 <0.005");
          return
        }
        this.$http.post("/kc/simulation/currency", {
          mining_interest_rate: this.miningInterestRate,
          competition_interest_rate: this.competitionInterestRate
        })
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data !== null) {
            this.$message.success("操作成功");
            this.$refs["create"].resetFields();
            this.showCreate = false;
            this.getSCurrency();
          };
        }, response => {
          this.$message.error(response.body);
        });
      },
      rowClick(row, column) {
        this.dialogTableVisible = true;
        this.dialogTable = row;
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
