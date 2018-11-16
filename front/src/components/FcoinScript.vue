<template>
  <div>
    <el-button type="primary" @click="onSubmit">查询</el-button>
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
      this.getLogs();
	  },
    methods: {
      getLogs() {
        this.$http.get("https://api.fcoin.com/v2/public/symbols")
        .then(response => {
          console.log(response)
        }, response => {
          this.$message.error(response.body);
        });
      },
    }
  };
</script>
