<template>
  <div>
    <el-form :inline="true" :model="formInline" class="demo-form-inline">
      <el-select v-model="formInline.mining" placeholder="结果" @change="onSubmit">
        <el-option key="0" label="挖矿" value="0"></el-option>
        <el-option key="1" label="推广" value="1"></el-option>
      </el-select>
    </el-form>
    <el-table :data="miningData" border show-summary style="width: 100%; max-width:601px;" v-show="formInline.mining == '0'">
      <el-table-column prop="uid" label="uid" width="200" sortable></el-table-column>
      <el-table-column prop="amount" label="锁仓" width="200" sortable></el-table-column>
      <el-table-column prop="mining" label="收益" width="200" sortable></el-table-column>
    </el-table>
    <el-table :data="shareData" border show-summary style="width: 100%; max-width:801px;" v-show="formInline.mining == '1'">
      <el-table-column prop="uid" label="uid" width="200" sortable></el-table-column>
      <el-table-column prop="extension_rate" label="个人推广算力" width="150" sortable>
        <template slot-scope="scope">
          {{ utilHelper.amountShow(scope.row.extension_rate) }}
        </template>
      </el-table-column>
      <el-table-column prop="extension_reward" label="个人推广收益" width="150" sortable>
        <template slot-scope="scope">
          {{ utilHelper.amountShow(scope.row.extension_reward) }}
        </template>
      </el-table-column>
      <el-table-column prop="competition_rate" label="推广竞赛算力" width="150" sortable>
        <template slot-scope="scope">
          {{ utilHelper.amountShow(scope.row.competition_rate) }}
        </template>
      </el-table-column>
      <el-table-column prop="competition_reward" label="推广竞赛奖励" width="150" sortable>
        <template slot-scope="scope">
          {{ utilHelper.amountShow(scope.row.competition_reward) }}
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
        miningData: [],
        shareData: [],
        datacollection: null,
      };
    },
    created () {
      this.getProfit();
    },
    methods: {
      getProfit() {
        if (this.formInline.mining == "0") {
            this.formInline.desc = "mining";
            this.miningData = [];
        } else {
            this.formInline.desc = "share";
            this.shareData = [];
        }
        this.$http.get("/kc/admin/predistribution", {"params": this.formInline})
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data !== null) {
            if (this.formInline.mining == "0") {
                this.miningData = data;
            } else {
                this.shareData = data;
            }
            
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
