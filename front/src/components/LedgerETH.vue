<template>
  <div>
    <el-col :xs="24" :sm="24" :md="24" :lg="12" :xl="12">
      <el-form ref="eth" :model="ethForm" :rules="ethRules" label-width="80px" style="width: 550px">
        <el-form-item prop="source" style="margin-bottom: 100px">
          <el-input v-model="ethForm.source" placeholder="ETH 地址,查询余额" auto-complete="off">
            <el-button @click="getBalance" slot="append" icon="el-icon-search"></el-button>
          </el-input>
          <span style="font-size: 12px; color: #778899">余额 {{fromWei(balance)}} ETH</span>
          <el-button v-show="balance !== '0'" @click="sendAllBalance" type="text">发送全部</el-button>
        </el-form-item>
        <el-form-item label="币种" prop="currency" required>
          <el-select v-model="ethForm.currency" placeholder="请选择">
            <el-option v-for="item in currenciesList" :key="item.label" :label="item.label" :value="item.label"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="path" prop="path" required>
          <el-input v-model="ethForm.path" auto-complete="off"></el-input>
          <span style="font-size: 12px; color: #778899">example: 44'/60'/1'/1</span>
        </el-form-item>
        <el-form-item label="转账到地址" prop="address" required>
          <el-input v-model="ethForm.address" auto-complete="off"></el-input>
          <span style="font-size: 12px; color: #778899">example: 0x99B677A1b4b326d2eed60a8DB726b4E1869C5F17</span>
        </el-form-item>
        <el-form-item label="数量" prop="amount" required>
          <el-input v-model="ethForm.amount" auto-complete="off"></el-input>
        </el-form-item>
        <el-form-item label="gasPrice" prop="gasPrice" required>
          <el-input v-model="ethForm.gasPrice" auto-complete="off"><template slot="append">Gwei</template></el-input>
        </el-form-item>
        <el-form-item label="gasLimit" prop="gasLimit" required>
          <el-input v-model="ethForm.gasLimit" auto-complete="off"></el-input>
        </el-form-item>
      </el-form>
      <el-row>
        <el-button type="primary" @click="onCreate" :loading="ethLoading" style="width: 550px">生成交易</el-button>
      </el-row>
    </el-col>
    <el-col :xs="24" :sm="24" :md="24" :lg="12" :xl="12" v-show="showSign">
      <el-form ref="sign" :model="signForm" :rules="signRules" label-width="100px" style="width: 550px">
        <el-form-item label="签名交易" prop="content" required>
          <el-input type="textarea" :rows="6" v-model="signForm.content" readonly></el-input>
        </el-form-item>
      </el-form>
      <el-row>
        <el-button type="success" @click="onSubmit" :loading="signLoading" style="width: 550px">发送交易</el-button>
      </el-row>
    </el-col>
    <el-col :xs="24" :sm="24" :md="24" :lg="12" :xl="12" v-show="showTxHash">
      <el-alert
        title="txHash"
        type="success"
        :description="txHash"
        show-icon
        close-text="查看"
        @close="onCheck"
        style="margin-top: 50px; width: 550px">
      </el-alert>
    </el-col>
  </div>
</template>

<script>
  import ethUtils from '../ledger-eth-transaction'

  export default {
    data() {
      return {
        balance: "0",
        currenciesList: [],
        currencies: {},
        ethForm: {
          source: "",
          path: "",
          address: "",
          amount: "",
          gasPrice: "10",
          gasLimit: "21000",
          currency: "",
          decimal: 0,
        },
        ethRules: {

        },
        signForm: {
          content: "",
        },
        signRules: {

        },
        txHash: "",
        ethLoading: false,
        showSign: false,
        signLoading: false,
        showTxHash: false,
      };
    },
    created () {
      this.getCurrencies();
      this.ethForm.source = this.$route.query.address;
      this.ethForm.path = this.$route.query.address_index;
      if (this.ethForm.source) {
        this.getBalance();
      }
    },
    methods: {
      getCurrencies() {
        this.$http.get("/kc/admin/super/currencies")
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data != null) {
            var currenciesList = [{"label": "ETH", "value": ""}];
            var currencies = {};
            for (var i in data) {
              var cur = data[i];
              if (cur.token === 1) {
                currenciesList.push({"label": cur.currency});
                currencies[cur.currency] = {"decimal": cur.decimals, "contractAddress": cur.contract_address};
              }
            }
            this.currenciesList = currenciesList;
            this.currencies = currencies;
          }
        }, response => {
          this.$message.error(response.body);
        });
      },
      getBalance() {
        if (!this.ethForm.source) {
          this.$message({showClose: true, message: "请填写ETH地址", type: 'error'});
          return
        }
        ethUtils.getBalance(this.ethForm.source).then(res => {
          this.balance = res;
        })
      },
      fromWei(hex) {
        return ethUtils.fromWei(hex);
      },
      sendAllBalance() {
        this.ethForm.amount = ethUtils.allValue(this.balance, this.ethForm.gasPrice, this.ethForm.gasLimit);
      },
      onCreate() {
        this.ethLoading = true;
        this.showTxHash = false;
        this.$refs["eth"].validate((valid) => {
          if (valid) {
            ethUtils
            .sendTransaction(
              this.ethForm.path,
              this.ethForm.address,
              this.ethForm.amount,
              this.ethForm.gasPrice,
              this.ethForm.gasLimit,
              this.currencies[this.ethForm.currency].contractAddress,
              this.currencies[this.ethForm.currency].decimal,
            ).then(res => {
              this.showSign = true;
              this.signForm.content = res;
              this.ethLoading = false;
            }).catch(err => {
              this.$message({showClose: true, duration: 0, message: err, type: 'error'});
              this.ethLoading = false;
            })
          } else {
            this.ethLoading = false;
            return false;
          }
        })
      },
      onCheck() {
        window.open("https://etherscan.io/tx/" + this.txHash);
      },
      onSubmit() {
        this.$refs["sign"].validate((valid) => {
          if (valid) {
            this.signLoading = true;
            this.$http.get("https://api.etherscan.io/api", {
              params: {
                module: "proxy",
                action: "eth_sendRawTransaction",
                hex: this.signForm.content,
                apikey: "YourApiKeyToken"
              }
            }).then(res => {
              this.showSign = false;
              this.signLoading = false;
              this.signForm.content = "";
              if("error" in res.body) {
                var message;
                if(typeof res.body.error === "string" || res.body.error instanceof String) {
                  message = res.body.error;
                } else {
                  message = res.body.error.message;
                }
                this.$message({showClose: true, duration: 0, message: message, type: 'error'});
                return
              }
              this.$message('成功');
              this.showTxHash = true;
              this.txHash = res.body.result;
            }).catch(err => {
              this.$message({showClose: true, duration: 0, message: err, type: 'error'});
              this.showSign = false;
              this.signLoading = false;
              this.signForm.content = "";
            })
          } else {
            return false;
          }
        })
      }
    }
  };
</script>
