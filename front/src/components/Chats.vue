<template>
  <div>
    <el-table :data="tableData" border style="width: 100%; max-width:1001px;">
      <el-table-column prop="username" label="用户名" width="150">
        <template slot-scope="scope">
          {{ scope.row.username == firstName ? scope.row.username : '' }}
        </template>
      </el-table-column>
      <el-table-column prop="content" label="内容" width="500"></el-table-column>
      <el-table-column prop="username" label="用户名" width="150">
        <template slot-scope="scope">
          {{ scope.row.username != firstName ? scope.row.username : '' }}
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
        firstName: "",
        tableData: []
      };
    },
    created () {
      this.getChats();
    },
    methods: {
      getChats() {
        this.$http.get("/kc/admin/trade/" + this.$route.params.code + "/chats")
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data !== null) {
            this.tableData = data;
            if(data.length > 0) {
              this.firstName = data[0].username;
            }
          }
        }, response => {
          this.$message.error(response.body);
        });
      }
    }
  };
</script>
