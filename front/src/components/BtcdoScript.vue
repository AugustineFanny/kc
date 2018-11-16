<template>
  <div>
    <el-form ref="create" :model="form" :rules="rules" label-width="120px" class="demo-ruleForm">
      <el-form-item label="交易对" prop="symbol" required>
        <el-select v-model="form.symbol" clearable placeholder="请选择">
          <el-option filterable
            v-for="item in tableData"
            :key="item"
            :label="item"
            :value="item">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="数量" prop="amount" required>
        <el-input-number controls-position="right" v-model="form.amount" :min="0" :debounce="500"></el-input-number>
      </el-form-item>
      <el-form-item label="价格" prop="floatation" required>
        <el-input-number controls-position="right" v-model="form.floatation" :min="0" :max="1" :debounce="500"></el-input-number>
      </el-form-item>
      <el-form-item label="间隔(秒)" prop="rest" required>
        <el-input-number controls-position="right" v-model="form.rest" :min="0.2" :max="20" :debounce="500"></el-input-number>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button type="primary" @click="onCreate">提 交</el-button>
    </div>
    <div>
      <div>交易对：{{ status.symbol }}</div>
      <div>数量：{{ status.amount }}</div>
      <div>价格：{{ status.floatation }}</div>
      <div>间隔(秒)：{{ status.rest }}</div>
      <div>状态：{{ statusShow(status.startFlag) }}</div>
    </div>
    <el-form ref="start" :model="form2" :rules="rules2" label-width="120px" class="demo-ruleForm" style="width: 700px;">
      <el-form-item label="APIKEY" prop="apiKey" required>
        <el-input v-model="form2.apiKey"></el-input>
      </el-form-item>
      <el-form-item label="APISECRET" prop="apiSecret" required>
        <el-input v-model="form2.apiSecret"></el-input>
      </el-form-item>
      <el-form-item label="密码" prop="passwd" required>
        <el-input v-model="form2.passwd"></el-input>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button type="primary" @click="onStart" :disabled="status.startFlag">开 始</el-button>
      <el-button type="primary" @click="onStop" :disabled="!status.startFlag">结 束</el-button>
    </div>
    <el-card id="scroll" class="box-card" style="height: 400px; overflow: auto; font-family: monospace;">
      <div v-for="log in logs" class="text item" style="height: 100%;">
        <span style="display: inline-block; width: 130px;">[{{ log.time }}]:</span>{{ log.log }}
      </div>
    </el-card>
  </div>
</template>

<script>
  export default {
    data() {
      return {
        tableData: ["ETH_USDT", "BDB_USDT", "BDB_ETH", "IOST_BDB"],
        form: {
          symbol: "",
          amount: 0,
          floatation: 0,
          sellMinAmount: 0,
          buyMinAmount: 0,
          rest: 10,
        },
        rules: {

        },
        form2: {

        },
        rules2: {

        },
        status: {

        },
        logs: [],
        interval: null,
      }
    },
    created () {
      this.getConfigure();
      this.interval = setInterval(this.getLogs, 1000);
    },
    updated () {
      var ele = document.getElementById('scroll');
      console.log(ele);
      ele.scrollTop = ele.scrollHeight;
    },
    methods: {
      statusShow(status) {
        if (status) {
          return "进行中";
        } else {
          return "已停止";
        }
      },
      getConfigure() {
        this.$http.get("/kc/public/btcdo-script")
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data != null) {
            this.status = data;
          }
        }, response => {
          this.$message.error(response.body);
        });
      },
      getLogs() {
        this.$http.get("/kc/public/btcdo-script-logs")
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data != null) {
            this.logs = this.logs.concat(data);
          }
        }, response => {
          this.$message.error(response.body);
        });
      },
      onCreate() {
        this.$refs["create"].validate((valid) => {
          if (valid) {

            this.$http.post("/kc/public/btcdo-script", this.form)
            .then(response => {
              var data = this.utilHelper.handle(this, response);
              if(data !== null) {
                this.$message.success("操作成功");
                this.getConfigure();
              };
            }, response => {
              this.$message.error(response.body);
            });
          } else {
            return false;
          }
        });
      },
      onStart() {
        this.$refs["start"].validate((valid) => {
          if (valid) {
            this.$http.post("/kc/public/btcdo-script-start", this.form2)
            .then(response => {
              var data = this.utilHelper.handle(this, response);
              if(data !== null) {
                this.$message.success("操作成功");
                this.getConfigure();
              };
            }, response => {
              this.$message.error(response.body);
            });
          }
        });
      },
      onStop() {
        this.$http.post("/kc/public/btcdo-script-stop")
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data !== null) {
            this.$message.success("操作成功");
            this.getConfigure();
          };
        }, response => {
          this.$message.error(response.body);
        });
      }
    }
  }
</script>
