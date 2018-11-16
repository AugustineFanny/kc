<template>
  <div>
    <el-table :data="tableData" border style="width: 100%; max-width:651px;">
      <el-table-column prop="trade_code" label="交易号" width="150"></el-table-column>
      <el-table-column prop="status" label="状态" width="150">
        <template slot-scope="scope">
          {{ statusShow(scope.row.status) }}
        </template>
      </el-table-column>
      <el-table-column prop="appeal" label="申诉" width="150">
        <template slot-scope="scope">
          {{ appealShow(scope.row.appeal) }}
        </template>
      </el-table-column>
      <el-table-column prop="create_time" label="时间" width="200">
        <template slot-scope="scope">
          {{ utilHelper.timeShow(scope.row.create_time) }}
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
  export default {
    data() {
      return {
        tableData: []
      };
    },
    created () {
      this.getChats();
    },
    methods: {
      statusShow(status) {
        switch(status) {
          case 1:
            return '已注资，待付款';
          case 2:
            return '已付款，待放行';
          case 3:
            return '完成';
          case 4:
            return '取消';
          case 5:
            return '自动取消,锁币中';
          case 6:
            return '拒绝';
          case 7:
            return '自动取消,已解锁';
          default:
            return '';
        }
      },
      appealShow(appeal) {
        switch(appeal) {
          case 0:
            return '';
          case 1:
            return '买家申诉';
          case 2:
            return '卖家申诉';
          case 3:
            return '已处理,放币给买家';
          case 4:
            return '已处理,放币给卖家';
          default:
            return appeal;
        }
      },
      getChats(params, pageNo) {
        this.$http.get("/kc/admin/trade/" + this.$route.params.code + "/trade_logs")
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data !== null) {
            this.tableData = data;
          }
        }, response => {
          this.$message.error(response.body);
        });
      }
    }
  };
</script>
