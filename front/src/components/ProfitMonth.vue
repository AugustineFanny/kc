<template>
  <div>
    <el-form :inline="true" :model="formInline" class="demo-form-inline">
      <el-form-item label="">
        <el-date-picker v-model="formInline.date" type="daterange" start-placeholder="开始日期" end-placeholder="结束日期" value-format="yyyy-MM-dd"></el-date-picker>
      </el-form-item>
      <el-select v-model="formInline.mining" placeholder="结果">
        <el-option key="0" label="挖矿" value="0"></el-option>
        <el-option key="1" label="推广" value="1"></el-option>
      </el-select>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">查询</el-button>
      </el-form-item>
      <el-form-item>
        <el-switch v-model="showChart" active-text="图表" inactive-text="表格"></el-switch>
      </el-form-item>
    </el-form>
    <el-table :data="tableData" border show-summary style="width: 100%; max-width:551px;" v-show="!showChart">
      <el-table-column prop="uid" label="用户(uid)" width="200">
        <template slot-scope="scope">
          {{ scope.row.username }}({{ scope.row.uid }})
        </template>
      </el-table-column>
      <el-table-column prop="currency" label="币种" width="150"></el-table-column>
      <el-table-column prop="amount" label="收益" width="200">
        <template slot-scope="scope">
          {{ utilHelper.amountShow(scope.row.amount).toFixed(8) }}
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
  import * as moment from 'moment'
  export default {
    components: {
        BarChart
    },
    data() {
      var defaultStartDay = moment().startOf('month').format("YYYY-MM-DD");
      var defaultEndDay = moment().format("YYYY-MM-DD");
      return {
        showChart: false,
        formInline: {
            mining: "0",
            date: [defaultStartDay, defaultEndDay],
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
      fillData() {
        var datacollection = {
            labels: [],
            datasets: [{
                label: "收益",
                backgroundColor: "#ff6384",
                data: []
            }],
        };
        datacollection.labels = [];
        datacollection.datasets[0].data = [];
        for (var i in this.tableData) {
            datacollection.labels.push(this.tableData[i].username);
            datacollection.datasets[0].data.push(this.tableData[i].amount);
        }
        this.datacollection = datacollection;
      },
      getProfit() {
        if (this.formInline.mining == "0") {
            this.formInline.desc = "mining";
        } else {
            this.formInline.desc = "share";
        }
        this.formInline.start_date = this.formInline.date[0] + " 00:00:00";
        this.formInline.end_date = this.formInline.date[1] + " 23:59:59";
        this.$http.get("/kc/admin/profit-month", {"params": this.formInline})
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
