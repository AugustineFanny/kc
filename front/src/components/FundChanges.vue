<template>
  <div>
    <el-form :inline="true" :model="formInline" class="demo-form-inline">
      <el-form-item label="">
        <el-input v-model="formInline.uid" style="width:150px;" placeholder="UID"></el-input>
        <el-input v-model="formInline.currency" style="width:150px;" placeholder="币种"></el-input>
        <el-select v-model="formInline.direction" clearable placeholder="增or减">
          <el-option key="0" label="增" value="0"></el-option>
          <el-option key="1" label="减" value="1"></el-option>
        </el-select>
        <el-select v-model="formInline.desc" clearable placeholder="操作">
          <el-option key="deposit" label="充币" value="deposit"></el-option>
          <el-option key="exchange" label="兑换" value="exchange"></el-option>
          <el-option key="distribute" label="分发" value="distribute"></el-option>
          <el-option key="instation" label="站内转" value="instation"></el-option>
          <el-option key="withdraw" label="提币" value="withdraw"></el-option>
          <el-option key="withdraw" label="挖矿" value="mining"></el-option>
          <el-option key="withdraw" label="推广" value="share"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">查询</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="tableData" border style="width: 100%; max-width:781px;">
      <el-table-column prop="uid" label="UID" width="100"></el-table-column>
      <el-table-column prop="currency" label="币种" width="150"></el-table-column>
      <el-table-column prop="amount" label="总金额" width="200">
        <template slot-scope="scope">
          <span v-if="scope.row.direction == 0" style="color: green">
            + {{ utilHelper.amountShow(scope.row.amount).toFixed(8) }}
          </span>
          <span v-else style="color: red">
            - {{ utilHelper.amountShow(scope.row.amount).toFixed(8) }}
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="desc" label="操作" width="150">
        <template slot-scope="scope">
          {{ showDesc(scope.row.desc) }}
        </template>
      </el-table-column>
      <el-table-column prop="create_time" label="更新时间" width="180">
        <template slot-scope="scope">
          {{ utilHelper.timeShow(scope.row.create_time) }}
        </template>
      </el-table-column>
    </el-table>
    <div style="position:absolute;">
      <el-pagination layout="total, prev, pager, next, jumper" :page-size="formInline.page_size" :total="totalCount" @current-change="currentChange">
      </el-pagination>
    </div>
  </div>
</template>

<script>
  export default {
    data() {
      return {
        formInline: {
          uid: '',
          currency: '',
          direction: null,
          desc: null,
          page: 1,
          page_size: 20,
        },
        totalCount: 0,
        tableData: []
      };
    },
    created () {
      if(this.$route.query.uid) {
        this.formInline.uid = this.$route.query.uid;
      }
      this.getWallets();
    },
    methods: {
      showDesc(desc) {
        switch(desc) {
          case "deposit": return "充币";
          case "exchange": return "兑换";
          case "distribute": return "分发";
          case "instation": return "站内转";
          case "withdraw": return "提币";
          case "trade": return "交易";
          case "mining": return "挖矿";
          case "share": return "推广";
          default: return desc;
        }
      },
      getWallets(params, pageNo) {
        this.$http.get("/kc/admin/fund-changes", {"params": this.formInline})
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
      onSubmit() {
        this.formInline.page = 1;
        this.getWallets();
      },
      currentChange(pageNo) {
        this.formInline.page = pageNo;
        this.getWallets();
      }
    }
  };
</script>
