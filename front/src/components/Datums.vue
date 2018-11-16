<template>
  <div>
    <el-table :data="tableData" border style="width: 100%; max-width:981px;">
      <el-table-column prop="trade_code" label="交易号" width="150"></el-table-column>
      <el-table-column prop="username" label="用户" width="150"></el-table-column>
      <el-table-column prop="content" label="申诉原因" width="300"></el-table-column>
      <el-table-column prop="datum" label="图片" width="180">
        <template slot-scope="scope">
          <img :src="img(scope.row.datum)" @click="enlarged(scope.row.datum)" class="img">
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
      this.getDatums();
    },
    methods: {
      getDatums(params, pageNo) {
        this.$http.get("/kc/admin/trade/" + this.$route.params.code + "/datums")
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data !== null) {
            this.tableData = data;
          }
        }, response => {
          this.$message.error(response.body);
        });
      },
      img(filename) {
        return '/uphp/appeal/' + filename;
      },
      enlarged(filename) {
        const h = this.$createElement;
        this.$msgbox({
          title: '消息',
          message: h('img', {attrs: {src: this.img(filename)}}),
          showCancelButton: false,
          showConfirmButton: false
        });
      }
    }
  };
</script>

<style type="text/css">
  .img {
    width: 150px;
    height: 150px;
  }
  .el-message-box {
    width: auto;
    min-width: 420px;
  }
</style>
