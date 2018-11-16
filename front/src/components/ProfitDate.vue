<template>
  <div>
    <el-form :inline="true" :model="formInline" class="demo-form-inline">
      <el-form-item label="">
        <el-date-picker v-model="formInline.date" type="date" value-format="yyyy-MM-dd" placeholder="选择日期"></el-date-picker>
      </el-form-item>
      <el-select v-model="formInline.mining" placeholder="结果">
        <el-option key="0" label="锁仓" value="0"></el-option>
        <el-option key="1" label="推广" value="1"></el-option>
      </el-select>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">查询</el-button>
      </el-form-item>
      <el-form-item>
        <el-switch v-model="showChart" active-text="图表" inactive-text="表格"></el-switch>
      </el-form-item>
    </el-form>
    <el-table :data="tableData" border show-summary style="width: 100%; max-width:1081px;" v-show="!showChart">
      <el-table-column prop="uid" label="用户(uid)" width="200">
        <template slot-scope="scope">
          {{ scope.row.username }}({{ scope.row.uid }})
        </template>
      </el-table-column>
      <el-table-column prop="currency" label="币种" width="100"></el-table-column>
      <el-table-column prop="reward" label="收益" width="200">
        <template slot-scope="scope">
          {{ utilHelper.amountShow(scope.row.reward) }}
        </template>
      </el-table-column>
      <el-table-column prop="mining" label="锁仓" width="200"></el-table-column>
      <el-table-column label="算力(利率)" width="200">
        <template slot-scope="scope">
          {{ getRate(scope.row.rate, scope.row.interest_rate).toFixed(2) }} ({{ scope.row.interest_rate }})
        </template>
      </el-table-column>
      <el-table-column prop="update_time" label="创建时间" width="180">
        <template slot-scope="scope">
          {{ utilHelper.timeShow(scope.row.create_time) }}
        </template>
      </el-table-column>
    </el-table>
    <div v-show="showChart">
      <bar-chart :height="200" :chart-data="datacollection"></bar-chart>
    </div>
  </div>
</template>

<script>
  import BarChart from './BarChart.js'
  export default {
    components: {
        BarChart
    },
    data() {
      return {
        showChart: false,
        formInline: {
            mining: "1",
            date: '',
        },
        totalCount: 0,
        tableData: [],
        datacollection: null,
      };
    },
    created () {
      this.getProfit();
    },
    methods: {
      getRate(rate, interestRate) {
        return rate / interestRate
      },
      fillData() {
        var datacollection = {
            labels: [],
            datasets: [{
                label: "收益",
                backgroundColor: "#ff6384",
                data: []
            }, {
                label: "锁仓",
                backgroundColor: "#36a2eb",
                data: []
            }],
        };
        datacollection.labels = [];
        datacollection.datasets[0].data = [];
        datacollection.datasets[1].data = [];
        for (var i in this.tableData) {
            datacollection.labels.push(this.tableData[i].username);
            datacollection.datasets[0].data.push(this.tableData[i].reward);
            datacollection.datasets[1].data.push(this.tableData[i].mining);
        }
        this.datacollection = datacollection;
      },
      getProfit() {
        if (this.formInline.mining == "0") {
            this.formInline.desc = "mining";
        } else {
            this.formInline.desc = "share";
        }
        this.$http.get("/kc/admin/profit-date", {"params": this.formInline})
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data !== null) {
            this.tableData = data;
            this.fillData();
          }
        }, response => {
          this.$message.error(response.body);
        });
      },
      onSubmit() {
        this.getProfit();
      },
    }
  };
</script>
