<template>
  <div>
<!--     <el-form :inline="true" :model="formInline" class="demo-form-inline">
      <el-form-item>
        <el-input v-model="formInline.ident" style="width:300px;" clearable placeholder="操作人员"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">查询</el-button>
      </el-form-item>
    </el-form> -->
  	<el-table :data="tableData" border style="width: 100%; max-width: 251px;">
      <el-table-column prop="username" label="用户名" width="150"></el-table-column>
	    <el-table-column prop="number" label="在线数" width="100"></el-table-column>
	  </el-table>
    <div class="block">
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
          ident: '',
        },
        totalCount: 0,
        tableData: []
      };
    },
    created () {
      this.getLogs();
	  },
    methods: {
      getLogs(params) {
        this.$http.get("/kc/admin/super/online", {"params": this.formInline})
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data !== null) {
            this.formInline.page = data.page_no;
            this.formInline.page_size = data.page_size;
            this.totalCount = data.total_count;
            for (var i in data.list) {
              this.tableData.push({"username": i, "number": data.list[i].length});
            }
          }
        }, response => {
          this.$message.error(response.body);
        });
      },
      currentChange(pageNo) {
        this.formInline.page = pageNo;
        this.getLogs();
      },
      onSubmit() {
        this.formInline.page = 1;
        this.getLogs();
      }
    }
  };
</script>
