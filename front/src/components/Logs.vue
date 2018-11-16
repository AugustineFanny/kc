<template>
  <div>
    <el-form :inline="true" :model="formInline" class="demo-form-inline">
      <el-form-item>
        <el-input v-model="formInline.ident" style="width:300px;" clearable placeholder="操作人员"></el-input>
      </el-form-item>
      <el-form-item>
        <el-input v-model="formInline.api" style="width:200px;" placeholder="调用接口"></el-input>
      </el-form-item>
      <el-select v-model="formInline.status" clearable placeholder="结果">
        <el-option key="1" label="成功" value="1"></el-option>
        <el-option key="0" label="失败" value="0"></el-option>
      </el-select>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">查询</el-button>
      </el-form-item>
    </el-form>
  	<el-table :data="tableData" border style="width: 100%">
      <el-table-column prop="id" label="ID" width="100"></el-table-column>
	    <el-table-column prop="username" label="操作人员(AID)" width="250">
        <template slot-scope="scope">
          {{scope.row.username}}  ({{scope.row.aid}})
        </template>
      </el-table-column>
	    <el-table-column prop="api" label="调用接口" width="160"></el-table-column>
      <el-table-column prop="args" label="传参" width="400"></el-table-column>
      <el-table-column prop="status" label="结果" width="100">
        <template slot-scope="scope">
          <el-tag :type="scope.row.status ? 'primary' : 'success'" close-transition>{{scope.row.status ? '成功' : '失败'}}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="sheet" label="表" width="200"></el-table-column>
      <el-table-column prop="column_id" label="所在行Id" width="100"></el-table-column>
      <el-table-column prop="create_time" label="调用时间" width="180">
        <template slot-scope="scope">
          {{ utilHelper.timeShow(scope.row.create_time) }}
        </template>
      </el-table-column>
      <el-table-column prop="remark" label="备注" width="300"></el-table-column>
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
          api: '',
          status: '',
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
        this.$http.get("/kc/admin/super/logs", {"params": this.formInline})
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
        this.getLogs();
      },
      onSubmit() {
        this.formInline.page = 1;
        this.getLogs();
      }
    }
  };
</script>
