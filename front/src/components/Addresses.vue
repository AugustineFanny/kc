<template>
  <div>
    <el-form :inline="true" :model="formInline" class="demo-form-inline">
      <el-form-item>
        <el-input v-model="formInline.uid" style="width:300px;" clearable placeholder="用户ID"></el-input>
      </el-form-item>
      <el-select v-model="formInline.currency" clearable placeholder="全部">
        <el-option key="BTC" label="BTC" value="BTC"></el-option>
        <el-option key="ETH" label="ETH" value="ETH"></el-option>
      </el-select>
      <el-select v-model="formInline.no_zero" clearable placeholder="全部">
        <el-option key="1" label="隐藏0余额" value="1"></el-option>
      </el-select>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">查询</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="tableData" border style="width: 100%; max-width:1131px;">
      <el-table-column prop="uid" label="用户ID" width="100"></el-table-column>
      <el-table-column prop="address" label="地址" width="400"></el-table-column>
      <el-table-column prop="address_index" label="索引" width="150"></el-table-column>
      <el-table-column label="可用/未确认" width="200">
        <template slot-scope="scope">
          {{ scope.row.now_amount }} {{ scope.row.currency }} / {{ scope.row.unconfirmed_amount }} {{ scope.row.currency }}
        </template>
      </el-table-column>
      <el-table-column prop="update_time" label="更新时间" width="180">
        <template slot-scope="scope">
          {{ utilHelper.timeShow(scope.row.update_time) }}
        </template>
      </el-table-column>
      <el-table-column label="累计充值" width="150">
        <template slot-scope="scope">
          {{ scope.row.all_amount }} {{ scope.row.currency }}
        </template>
      </el-table-column>
      <el-table-column fixed="right" label="操作" width="50">
        <template slot-scope="scope">
          <el-button v-show="scope.row.currency === 'ETH'" @click="handleClick(scope.row)" type="text" size="small">转币</el-button>
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
          page: 1,
          page_size: 20,
          uid: null,
          currency: '',
          no_zero: '',
        },
        totalCount: 0,
        tableData: []
      };
    },
    created () {
      this.getAddresses();
    },
    methods: {
      getAddresses(params) {
        this.$http.get("/kc/admin/addresses", {"params": this.formInline})
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
      currentChange(pageNo) {
        this.formInline.page = pageNo;
        this.getAddresses();
      },
      onSubmit() {
        this.formInline.page = 1;
        this.getAddresses();
      },
      handleClick(row) {
        this.$router.push({path:'/admin/ledgerETH',
          query:{
            address: row.address,
            address_index: row.address_index
          }
        })
      }
    }
  };
</script>
