<template>
  <div>
    <el-form :inline="true" :model="formInline" class="demo-form-inline">
      <el-form-item label="">
        <el-input v-model="formInline.ident" style="width:300px;" placeholder="UID"></el-input>
        <el-input v-model="formInline.currency" style="width:200px;" placeholder="currency"></el-input>
        <el-checkbox v-model="formInline.no_zero">只显示已锁仓</el-checkbox>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">查询</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="tableData" border style="width: 100%; max-width:1231px;" @sort-change="sortChange">
      <el-table-column prop="uid" label="UID" sortable="custom" width="100"></el-table-column>
      <el-table-column prop="currency" label="币种" width="150"></el-table-column>
      <el-table-column prop="amount" label="总金额" sortable="custom" width="200">
        <template slot-scope="scope">
          {{ utilHelper.amountShow(scope.row.amount).toFixed(8) }}
        </template>
      </el-table-column>
      <el-table-column label="可用" width="200">
        <template slot-scope="scope">
          {{ utilHelper.amountShow(scope.row.amount - scope.row.lock_amount).toFixed(8) }}
        </template>
      </el-table-column>
      <el-table-column prop="lock_amount" label="冻结" width="200">
        <template slot-scope="scope">
          {{ utilHelper.amountShow(scope.row.lock_amount).toFixed(8) }}
        </template>
      </el-table-column>
      <el-table-column prop="mining_amount" sortable="custom" label="锁仓" width="200">
        <template slot-scope="scope">
          {{ utilHelper.amountShow(scope.row.mining_amount).toFixed(8) }}
        </template>
      </el-table-column>
      <el-table-column prop="update_time" label="更新时间" width="180">
        <template slot-scope="scope">
          {{ utilHelper.timeShow(scope.row.update_time) }}
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
          ident: '',
          currency: '',
          no_zero: null,
          page: 1,
          page_size: 20,
          _sort: "",
          _order: "",
        },
        totalCount: 0,
        tableData: []
      };
    },
    created () {
      if(this.$route.query.uid) {
        this.formInline.ident = this.$route.query.uid;
      }
      this.getWallets();
    },
    methods: {
      getWallets(params, pageNo) {
        this.$http.get("/kc/admin/wallets", {"params": this.formInline})
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
      },
      sortChange(a) {
        this.formInline._sort = a.prop;
        this.formInline._order = a.order;
        this.formInline.page = 1;
        this.getWallets();
      },
    }
  };
</script>
