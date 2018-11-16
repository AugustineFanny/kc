<template>
  <div>
    <el-form :inline="true" :model="formInline" class="demo-form-inline">
      <el-form-item label="">
        <el-date-picker v-model="formInline.date" type="daterange" start-placeholder="开始日期" end-placeholder="结束日期" value-format="yyyy-MM-dd"></el-date-picker>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">查询</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="tableData" border style="width: 100%; max-width:751px;">
      <el-table-column prop="username" label="用户(UID)" width="150">
        <template slot-scope="scope">
          {{ scope.row.username }}({{ scope.row.uid }})
        </template>
      </el-table-column>
      <el-table-column prop="m1" label="M1" width="100" sortable></el-table-column>
      <el-table-column prop="m2" label="M2" width="100" sortable></el-table-column>
      <el-table-column prop="m3" label="M3" width="100" sortable></el-table-column>
      <el-table-column prop="m4" label="M4" width="100" sortable></el-table-column>
      <el-table-column prop="m5" label="M5" width="100" sortable></el-table-column>
      <el-table-column prop="m6" label="M6" width="100" sortable></el-table-column>
    </el-table>
  </div>
</template>

<script>
  import * as moment from 'moment';
  export default {
    data() {
      var defaultStartDay = moment().startOf('month').format("YYYY-MM-DD");
      var defaultEndDay = moment().format("YYYY-MM-DD");
      return {
        tableData: [],
        formInline: {
            date: [defaultStartDay, defaultEndDay],
        },
      };
    },
    created () {
      this.getInvites();
    },
    methods: {
      getInvites() {
        if (this.formInline.date) {
          this.formInline.start_date = this.formInline.date[0] + " 00:00:00";
          this.formInline.end_date = this.formInline.date[1] + " 23:59:59";
        } else {
          this.formInline.start_date = "";
          this.formInline.end_date = "";
        }
        this.$http.get("/kc/admin/invites", {"params": this.formInline})
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data !== null) {
            for (var i in data) {
              data[i]["m1"] = parseInt(data[i]["m1"]);
              data[i]["m2"] = parseInt(data[i]["m2"]);
              data[i]["m3"] = parseInt(data[i]["m3"]);
              data[i]["m4"] = parseInt(data[i]["m4"]);
              data[i]["m5"] = parseInt(data[i]["m5"]);
              data[i]["m6"] = parseInt(data[i]["m6"]);
            }
            this.tableData = data;
          }
        }, response => {
          this.$message.error(response.body);
        });
      },
      onSubmit() {
        this.getInvites();
      },
    }
  };
</script>
