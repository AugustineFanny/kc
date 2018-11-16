<template>
  <div>
    <el-form :inline="true" :model="formInline" class="demo-form-inline">
      <el-form-item label="">
        <el-input v-model="formInline.ident" style="width:300px;" placeholder="UID/用户名"></el-input>
        <el-input v-model="formInline.code" style="width:300px;" placeholder="广告code"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">查询</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="tableData" border style="width: 100%; max-width:1071px;">
      <el-table-column prop="code" label="广告CODE" sortable="custom" width="150"></el-table-column>
      <el-table-column prop="username" label="所属用户" width="100"></el-table-column>
      <el-table-column prop="unit" label="交易对" width="120">
        <template slot-scope="scope">{{scope.row.unit}}/{{scope.row.currency}}</template>
      </el-table-column>
      <el-table-column prop="direction" label="类型" width="100">
        <template slot-scope="scope">{{scope.row.direction == 'sell' ? '出售' : '购买'}}</template>
      </el-table-column>
      <el-table-column prop="fixed_price" label="定价方式" width="100">
        <template slot-scope="scope">{{scope.row.fixed_price ? '固定' : '浮动'}}</template>
      </el-table-column>
      <el-table-column prop="price" label="价格" width="100">
        <template slot-scope="scope">{{scope.row.fixed_price ? scope.row.price : scope.row.exchange }}</template>
      </el-table-column>
      <el-table-column prop="amount" label="数量" width="100"></el-table-column>
      <el-table-column prop="price" label="限额" width="200">
        <template slot-scope="scope">{{scope.row.min_price}}~{{scope.row.max_price}}</template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template slot-scope="scope">{{scope.row.status ? '禁用' : '正常'}}</template>
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
          code: '',
          page: 1,
          page_size: 20,
        },
        totalCount: 0,
        tableData: []
      };
    },
    created () {
      this.getAds();
    },
    methods: {
      getAds(params, pageNo) {
        this.$http.get("/kc/admin/ads", {"params": this.formInline})
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
        this.getAds();
      },
      currentChange(pageNo) {
        this.formInline.page = pageNo;
        this.getAds();
      }
    }
  };
</script>
