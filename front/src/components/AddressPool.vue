<template>
  <div>
    <el-form :inline="true" :model="formInline" class="demo-form-inline">
      <el-select v-model="formInline.currency" clearable placeholder="全部">
        <el-option key="BTC" label="BTC" value="BTC"></el-option>
        <el-option key="ETH" label="ETH" value="ETH"></el-option>
      </el-select>
      <el-select v-model="formInline.flag" clearable placeholder="全部">
        <el-option key="0" label="未使用" value="0"></el-option>
        <el-option key="1" label="已使用" value="1"></el-option>
        <el-option key="2" label="保留" value="2"></el-option>
      </el-select>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">查询</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="tableData" border style="width: 100%; max-width:751px;">
      <el-table-column prop="currency" label="币种" width="100"></el-table-column>
      <el-table-column prop="address" label="地址" width="400"></el-table-column>
      <el-table-column prop="address_index" label="索引" width="150"></el-table-column>
      <el-table-column label="状态" width="100">
        <template slot-scope="scope">
          {{ showFlag(scope.row.flag) }}
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
          currency: '',
          flag: null
        },
        totalCount: 0,
        tableData: []
      };
    },
    created () {
      this.getPool();
    },
    methods: {
      showFlag(flag) {
        switch(flag) {
          case 0:
            return "未使用";
          case 1:
            return "已使用";
          case 2:
            return "保留";
          default:
            return flag;
        }
      },
      getPool(params) {
        this.$http.get("/kc/admin/address/pool", {"params": this.formInline})
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
        this.getPool();
      },
      onSubmit() {
        this.formInline.page = 1;
        this.getPool();
      }
    }
  };
</script>
