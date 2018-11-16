<template>
  <div>
    <el-row><el-button type="primary" @click="showCreate = true">创建币种</el-button></el-row>
    <el-row>
      <el-col :xs="24" :sm="24" :md="16" :lg="16">
        <el-table :data="tableData" border style="width: 100%; max-width: 951px;">
          <el-table-column prop="id" label="" width="80"></el-table-column>
          <el-table-column prop="currency" label="币" width="100"></el-table-column>
          <el-table-column prop="base_price" label="对rmb价格" width="100"></el-table-column>
          <el-table-column prop="mining_interest_rate" label="挖矿利率" width="100"></el-table-column>
          <el-table-column prop="competition_interest_rate" label="推广竞赛利率" width="100"></el-table-column>
          <el-table-column label="ERC20" width="100">
            <template slot-scope="scope">
              <el-tag :type="scope.row.token ? 'success' : 'primary'" close-transition>{{scope.row.token ? '是' : '否'}}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="可否充值" width="100">
            <template slot-scope="scope">
              <el-tag :type="scope.row.recharge ? 'primary' : 'success'" close-transition>{{scope.row.recharge ? '否' : '是'}}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="可否提币" width="100">
            <template slot-scope="scope">
              <el-tag :type="scope.row.withdraw ? 'primary' : 'success'" close-transition>{{scope.row.withdraw ? '否' : '是'}}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="confirm_num" label="充值确认数(仅用于显示)" width="120"></el-table-column>
          <el-table-column label="可否交易" width="100">
            <template slot-scope="scope">
              <el-tag :type="scope.row.trans_flag ? 'primary' : 'success'" close-transition>{{scope.row.trans_flag ? '否' : '是'}}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="hodl_last" label="锁仓最后奖励时间" width="180">
            <template slot-scope="scope">
              {{ utilHelper.timeShow(scope.row.hodl_last) }}
            </template>
          </el-table-column>
          <el-table-column prop="share_last" label="分享最后奖励时间" width="180">
            <template slot-scope="scope">
              {{ utilHelper.timeShow(scope.row.share_last) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100" fixed="right">
            <template slot-scope="scope">
              <el-button @click="handleClick(scope.row)" type="text" size="small">编辑</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-col>
      <el-col :xs="24" :sm="24" :md="8" :lg="8" v-show="show">
        <el-form ref="change" :model="form" :rules="rules" label-width="120px">
          <el-form-item label="币" prop="currency">
            <el-input readonly v-model="form.currency"></el-input>
          </el-form-item>
          <el-form-item label="对rmb价格" prop="base_price" required>
            <el-input-number controls-position="right" v-model="form.base_price" :min="0" :debounce="500"></el-input-number>
          </el-form-item>
          <el-form-item label="挖矿利率" prop="mining_interest_rate" required>
            <el-input-number controls-position="right" v-model="form.mining_interest_rate" :min="0" :debounce="500"></el-input-number>
          </el-form-item>
          <el-form-item label="推广竞赛利率" prop="competition_interest_rate" required>
            <el-input-number controls-position="right" v-model="form.competition_interest_rate" :min="0" :debounce="500"></el-input-number>
          </el-form-item>
          <el-form-item label="最小锁仓数量" prop="min_lock" required>
            <el-input-number controls-position="right" v-model="form.min_lock" :min="0" :debounce="500"></el-input-number>
          </el-form-item>
          <el-form-item label="ERC20" prop="token" required>
            <el-select v-model="form.token" placeholder="请选择">
              <el-option key="0" label="否" value="0" ></el-option>
              <el-option key="1" label="是" value="1" ></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="合约地址" prop="contract_address">
            <el-input v-model="form.contract_address"></el-input>
          </el-form-item>
          <el-form-item label="小数点位数" prop="decimals" required>
            <el-input-number controls-position="right" v-model="form.decimals" :min="0" :debounce="500"></el-input-number>
          </el-form-item>
          <el-form-item label="充值确认数(仅用于显示)" prop="confirm_num" required>
            <el-input-number controls-position="right" v-model="form.confirm_num" :min="1"></el-input-number>
          </el-form-item>
          <el-form-item label="可否充值" prop="recharge" required>
            <el-select v-model="form.recharge" placeholder="请选择">
              <el-option key="0" label="是" value="0" ></el-option>
              <el-option key="1" label="否" value="1" ></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="可否提币" prop="withdraw" required>
            <el-select v-model="form.withdraw" placeholder="请选择">
              <el-option key="0" label="是" value="0" ></el-option>
              <el-option key="1" label="否" value="1" ></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="可否交易" prop="trans_flag" required>
            <el-select v-model="form.trans_flag" placeholder="请选择">
              <el-option key="0" label="是" value="0" ></el-option>
              <el-option key="1" label="否" value="1" ></el-option>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="onChange">提交</el-button>
            <el-button @click="onCancel">取消</el-button>
          </el-form-item>
        </el-form>
      </el-col>
    </el-row>
    <el-dialog title="创建币种" :visible.sync="showCreate">
      <el-form ref="create" :model="form2" :rules="rules2" label-width="100px" class="demo-ruleForm">
        <el-form-item label="币种" prop="currency">
          <el-input v-model="form2.currency" auto-complete="off"></el-input>
        </el-form-item>
        <el-form-item label="对rmb价格" prop="base_price" required>
          <el-input-number controls-position="right" v-model="form2.base_price" :min="0" :debounce="500"></el-input-number>
        </el-form-item>
        <el-form-item label="挖矿利率" prop="mining_interest_rate" required>
            <el-input-number controls-position="right" v-model="form2.mining_interest_rate" :min="0" :debounce="500"></el-input-number>
          </el-form-item>
          <el-form-item label="推广竞赛利率" prop="competition_interest_rate" required>
            <el-input-number controls-position="right" v-model="form2.competition_interest_rate" :min="0" :debounce="500"></el-input-number>
          </el-form-item>
        <el-form-item label="最小锁仓数量" prop="min_lock" required>
          <el-input-number controls-position="right" v-model="form2.min_lock" :min="0" :debounce="500"></el-input-number>
        </el-form-item>
        <el-form-item label="ERC20" prop="token">
            <el-select v-model="form2.token" placeholder="请选择">
              <el-option key="0" label="否" value="0" ></el-option>
              <el-option key="1" label="是" value="1" ></el-option>
            </el-select>
        </el-form-item>
        <el-form-item label="合约地址" prop="contract_address">
          <el-input v-model="form2.contract_address"></el-input>
        </el-form-item>
        <el-form-item label="小数点位数" prop="decimals" required>
            <el-input-number controls-position="right" v-model="form2.decimals" :min="0" :debounce="500"></el-input-number>
          </el-form-item>
        <el-form-item label="充值确认数(仅用于显示)" prop="confirm_num">
          <el-input-number controls-position="right" v-model="form2.confirm_num" :min="1"></el-input-number>
        </el-form-item>
        <el-form-item label="可否充值" prop="recharge">
          <el-select v-model="form2.recharge" placeholder="请选择">
            <el-option key="0" label="是" value="0" ></el-option>
            <el-option key="1" label="否" value="1" ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="可否提币" prop="withdraw">
          <el-select v-model="form2.withdraw" placeholder="请选择">
            <el-option key="0" label="是" value="0" ></el-option>
            <el-option key="1" label="否" value="1" ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="可否交易" prop="trans_flag">
          <el-select v-model="form2.trans_flag" placeholder="请选择">
            <el-option key="0" label="是" value="0" ></el-option>
            <el-option key="1" label="否" value="1" ></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="showCreate = false">取 消</el-button>
        <el-button type="primary" @click="onCreate">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
  export default {
    data() {
      return {
        tableData: [],
        show: false,
        form: {
          id: '',
          currency: '',
          fee: 0,
          fee_eth: '0',
          min_amount: 0,
          fee_trade: 0,
          token: '0',
          contract_address: '',
          confirm_num: 1,
          recharge: '',
          withdraw: '',
          exchanges: '',
          trans_flag: '',
          base_price: 0,
          decimals: 0,
          min_lock: 0,
          mining_interest_rate: 0,
          competition_interest_rate: 0,
        },
        rules: {

        },
        showCreate: false,
        form2: {
          currency: '',
          fee: 0,
          fee_eth: '0',
          min_amount: 0,
          fee_trade: 0,
          token: '0',
          contract_address: '',
          confirm_num: 1,
          recharge: '1',
          withdraw: '1',
          exchanges: '',
          trans_flag: '1',
          base_price: 0,
          decimals: 0,
          min_lock: 0,
          mining_interest_rate: 0,
          competition_interest_rate: 0,
        },
        rules2: {
          currency: [
            { required: true, message: "请输入币种", trigger: 'blur'}
          ]
        }
      };
    },
    created () {
      this.getCurrencies();
    },
    methods: {
      getCurrencies() {
        this.$http.get("/kc/admin/super/currencies")
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data != null) {
            this.tableData = data;
          }
        }, response => {
          this.$message.error(response.body);
        });
      },
      handleClick(row) {
        this.show = true;
        for (var i in row) {
            this.form[i] = row[i];
        };
        this.form.fee_eth = this.form.fee_eth.toString();
        this.form.token = this.form.token.toString();
        this.form.recharge = this.form.recharge.toString();
        this.form.withdraw = this.form.withdraw.toString();
        this.form.trans_flag = this.form.trans_flag.toString();
      },
      onChange() {
        this.$refs["change"].validate((valid) => {
          if (valid) {
            this.$resource("/kc/admin/super/currency/{id}")
            .save({id: this.form.id}, {
              currency: this.form.currency,
              token: parseInt(this.form.token),
              contract_address: this.form.contract_address,
              confirm_num: this.form.confirm_num,
              recharge: parseInt(this.form.recharge),
              withdraw: parseInt(this.form.withdraw),
              trans_flag: parseInt(this.form.trans_flag),
              base_price: this.form.base_price,
              decimals: this.form.decimals,
              min_lock: this.form.min_lock,
              mining_interest_rate: this.form.mining_interest_rate,
              competition_interest_rate: this.form.competition_interest_rate,
            })
            .then(response => {
              var data = this.utilHelper.handle(this, response);
              if(data !== null) {
                this.$message.success("操作成功");
                this.$refs["change"].resetFields();
                this.show = false;
                this.getCurrencies();
              };
            }, response => {
              this.$message.error(response.body);
            });
          } else {
            return false;
          }
        });
      },
      onCreate() {
        this.$refs["create"].validate((valid) => {
          if (valid) {
            this.$http.post("/kc/admin/super/currencies", {
              currency: this.form2.currency,
              token: parseInt(this.form2.token),
              contract_address: this.form2.contract_address,
              confirm_num: this.form2.confirm_num,
              recharge: parseInt(this.form2.recharge),
              withdraw: parseInt(this.form2.withdraw),
              trans_flag: parseInt(this.form2.trans_flag),
              base_price: this.form2.base_price,
              decimals: this.form2.decimals,
              min_lock: this.form2.min_lock,
              mining_interest_rate: this.form2.mining_interest_rate,
              competition_interest_rate: this.form2.competition_interest_rate,
            })
            .then(response => {
              var data = this.utilHelper.handle(this, response);
              if(data !== null) {
                this.$message.success("操作成功");
                this.$refs["create"].resetFields();
                this.showCreate = false;
                this.getCurrencies();
              };
            }, response => {
              this.$message.error(response.body);
            });
          } else {
            return false;
          }
        });
      },
      onCancel() {
        this.show = false;
      }
    }
  };
</script>

<style>
  .el-row {
    margin-bottom: 20px;
    &:last-child {
      margin-bottom: 0;
    }
  }
</style>
