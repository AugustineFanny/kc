<template>
  <div>
    <el-table :data="tableData" border style="width: 100%; max-width:1051px;">
      <el-table-column prop="username" label="用户(UID)" width="150">
        <template slot-scope="scope">
          {{ scope.row.username }}({{ scope.row.uid }})
        </template>
      </el-table-column>
      <el-table-column prop="m1_amount" label="M1持仓" width="150" sortable>
        <template slot-scope="scope">
          {{ utilHelper.amountShow(scope.row.m1_amount).toFixed(2) }}
        </template>
      </el-table-column>
      <el-table-column prop="m1_active_amount" label="M1可用" width="150" sortable>
        <template slot-scope="scope">
          {{ utilHelper.amountShow(scope.row.m1_amount - scope.row.m1_lock_amount).toFixed(2) }}
        </template>
      </el-table-column>
      <el-table-column prop="m1_lock_amount" label="M1锁仓" width="150" sortable>
        <template slot-scope="scope">
          {{ utilHelper.amountShow(scope.row.m1_lock_amount).toFixed(2) }}
        </template>
      </el-table-column>
      <el-table-column prop="amount" label="所有下线持仓" width="150" sortable>
        <template slot-scope="scope">
          {{ utilHelper.amountShow(scope.row.amount).toFixed(2) }}
        </template>
      </el-table-column>
      <el-table-column prop="active_amount" label="所有下线可用" width="150" sortable>
        <template slot-scope="scope">
          {{ utilHelper.amountShow(scope.row.amount - scope.row.lock_amount).toFixed(2) }}
        </template>
      </el-table-column>
      <el-table-column prop="lock_amount" label="所有下线锁仓" width="150" sortable>
        <template slot-scope="scope">
          {{ utilHelper.amountShow(scope.row.lock_amount).toFixed(2) }}
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
  import * as moment from 'moment';
  export default {
    data() {
      return {
        tableData: [],
        formInline: {
          ident: ""
        },
      };
    },
    created () {
      this.getChildAmounts();
    },
    methods: {
      getChildAmounts() {
        this.$http.get("/kc/admin/child/amounts")
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data !== null) {
            this.tableData = data;
          }
        }, response => {
          this.$message.error(response.body);
        });
      },
      onSubmit() {
        this.getChildAmounts();
      },
    }
  };
</script>
