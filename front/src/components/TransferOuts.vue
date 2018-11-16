<template>
  <div>
    <el-form :inline="true" :model="formInline" class="demo-form-inline">
      <el-form-item label="">
        <el-input v-model="formInline.uid" style="width:200px;" placeholder="uid"></el-input>
        <el-input v-model="formInline.currency" style="width:200px;" placeholder="currency"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">查询</el-button>
      </el-form-item>
    </el-form>
    待审核数量：{{ totalCount }}
    <el-table :data="tableData" border style="width: 100%; max-width: 1321px;">
      <el-table-column prop="uid" label="uid" width="100"></el-table-column>
      <el-table-column prop="currency" label="币" width="100"></el-table-column>
      <el-table-column prop="amount" label="数量" width="120"></el-table-column>
      <el-table-column label="手续费" width="120">
        <template slot-scope="scope">
          {{ scope.row.fee }} {{ scope.row.fee_currency }}
        </template>
      </el-table-column>
      <el-table-column prop="to" label="目标地址" width="400"></el-table-column>
      <el-table-column prop="status_name" label="状态" width="120"></el-table-column>
      <el-table-column prop="create_time" label="申请时间" width="180">
        <template slot-scope="scope">{{ utilHelper.timeShow(scope.row.create_time) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="180">
        <template slot-scope="scope">
          <div v-show="scope.row.status == 3">
            <el-button size="small" @click="operation(scope, 2)">通过</el-button>
            <el-button size="small" type="danger" @click="operation(scope, 4)">不通过</el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
  export default {
    data() {
      return {
        formInline: {
          uid: '',
          currency: '',
          status: 3,
          page: 1,
          page_size: 20,
        },
        totalCount: 0,
        tableData: []
      };
    },
    created () {
      this.getTransfers();
      this.test();
    },
    methods: {
      statusName(status) {
        if(status == 0) {
          return "确认中";
        } else if (status == 1) {
          return "异常";
        } else if (status == 2) {
          return "成功";
        } else if (status == 3) {
          return "待审核";
        } else if (status == 4) {
          return "审核未通过";
        } else {
          return "未知状态：" + status;
        }
      },
      getTransfers() {
        this.$http.get("/kc/admin/transfers", {"params": this.formInline})
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data !== null) {
            this.totalCount = data.total_count;
            data.list.forEach((row) => {
              row.status_name = this.statusName(row.status);
            });
            this.tableData = data.list;
          }
        }, response => {
          this.$message.error(response.body);
        });
      },
      operation(scope, status) {
        if(status == 4) {
          this.$prompt('原因', '审核不通过', {
            confirmButtonText: '确定',
            cancelButtonText: '取消'
          }).then(({ value }) => {
            this.handle(scope, status, value);
          })
        } else if(status == 2) {
          this.$prompt('交易hash', '审核通过', {
            confirmButtonText: '确定',
            cancelButtonText: '取消'
          }).then(({ value }) => {
            this.handle(scope, status, value);
          })
        }
      },
      handle(scope, status, desc) {
        this.$http.post("/kc/admin/transfer-auth", {"id": scope.row.id, "status": status, "desc": desc}, {"responseType": "json"})
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data !== null) {
            this.tableData[scope.$index].status = status;
            this.tableData[scope.$index].status_name = this.statusName(status);
            this.$message.info("操作成功");
          }
        })
      },
      onSubmit() {
        this.formInline.page = 1;
        this.getTransfers();
      }
    }
  };
</script>
