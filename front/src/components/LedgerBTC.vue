<template>
  <div>
    <el-table v-loading="loading" :data="tableData" border style="width: 100%; max-width: 551px; margin-bottom: 100px">
      <el-table-column prop="addr" label="地址" width="350"></el-table-column>
      <el-table-column prop="associatedKeyset" label="path" width="200"></el-table-column>
    </el-table>
    <el-form ref="create" :model="btcForm" :rules="btcRules" label-width="100px" style="width: 500px">
      <el-form-item label="btc_address" prop="address" required>
        <el-input v-model="btcForm.address" auto-complete="off" width="200px"></el-input>
      </el-form-item>
    </el-form>
    <el-row>
      <el-button type="primary" @click="onCreate">BTC转账</el-button>
    </el-row>
  </div>
</template>

<script>
  import btcUtils from '../ledger-btc-transaction'

  export default {
    data() {
      return {
        tableData: [],
        btcForm: {
          address: '17BTcK31CBcwcMcJdiuXjAqiWmBqvMb8DR',
        },
        btcRules: {

        },
        loading: true,
      };
    },
    created () {
      this.getBtcNoZero();
    },
    methods: {
      getBtcNoZero() {
        this.$http.get("/kc/admin/btc-nozero")
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data != null) {
            this.tableData = data;
            this.loading = false;
          }
        }, response => {
          this.$message.error(response.body);
        });
      },
      onCreate() {
        this.$refs["create"].validate((valid) => {
          if (valid) {
            if( this.tableData.length === 0) {
              alert("没有查询到余额");
              return;
            }
            btcUtils.sendTransaction(this.tableData, this.btcForm.address);
          } else {
            return false;
          }
        });
      }
    }
  };
</script>
