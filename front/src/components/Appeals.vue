<template>
  <div>
    <el-form :inline="true" :model="formInline" class="demo-form-inline">
      <el-form-item label="">
        <el-input v-model="formInline.code" style="width:150px;" placeholder="交易号/广告号"></el-input>
        <el-select v-model="formInline.appeal" placeholder="请选择">
          <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">查询</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="tableData" border style="width: 100%; max-width:1579px;">
      <el-table-column type="expand">
        <template slot-scope="props">
          <el-form label-position="left" inline class="demo-table-expand">
            <el-form-item label="手续费">
              <span>{{ props.row.fee }} {{ props.row.currency }}</span>
            </el-form-item>
            <el-form-item label="所属广告">
              <span>{{ props.row.ad_code }}</span>
            </el-form-item>
            <el-form-item label="广告价格">
              <span>{{ props.row.ad_price }} {{ props.row.currency }}/{{ props.row.unit }}</span>
            </el-form-item>
            <el-form-item label="创建时间">
              <span>{{ utilHelper.timeShow(props.row.create_time) }}</span>
            </el-form-item>
            <el-form-item label="操作时间">
              <span>{{ utilHelper.timeShow(props.row.operate_time) }}</span>
            </el-form-item>
            <el-form-item label="过期时间">
              <span>{{ utilHelper.timeShow(props.row.expire_time) }}</span>
            </el-form-item>
          </el-form>
        </template>
      </el-table-column>
      <el-table-column prop="code" label="交易号" width="150"></el-table-column>
      <el-table-column prop="seller" label="卖方" width="150"></el-table-column>
      <el-table-column prop="buyer" label="买方" width="150"></el-table-column>
      <el-table-column prop="price" label="金额" width="150">
        <template slot-scope="scope">
          {{scope.row.price}} {{scope.row.unit}}
        </template>
      </el-table-column>
      <el-table-column prop="amount" label="数量" width="150">
        <template slot-scope="scope">
          {{scope.row.amount}} {{scope.row.currency}}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="交易状态" width="150">
        <template slot-scope="scope">
          {{ statusShow(scope.row) }}
        </template>
      </el-table-column>
      <el-table-column prop="appeal" label="申诉" width="100">
        <template slot-scope="scope">
          {{ appealShow(scope.row) }}
        </template>
      </el-table-column>
      <el-table-column prop="operate_time" label="申诉时间" width="180">
        <template slot-scope="scope">
          {{ utilHelper.timeShow(scope.row.operate_time) }}
        </template>
      </el-table-column>
      <el-table-column prop="chat" label="操作" width="150">
        <template slot-scope="scope">
          <router-link :to="{ path: '/trade/' + scope.row.code + '/chats' }">聊天</router-link>
          <router-link :to="{ path: '/trade/' + scope.row.code + '/logs' }">历史</router-link>
          <router-link :to="{ path: '/trade/' + scope.row.code + '/datums' }">材料</router-link>
        </template>
      </el-table-column>
      <el-table-column prop="appeal" label="操作" width="250">
        <template slot-scope="scope">
          <div v-show="scope.row.appeal == 1 || scope.row.appeal == 2">
            <el-button type="primary" size="mini" @click="releaseTradePrompt(scope)">放币给买家</el-button>
            <el-button type="danger" size="mini" @click="cancelTradePrompt(scope)">放币给卖家</el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>
    <div style="position:absolute;">
      <el-pagination layout="total, prev, pager, next, jumper" :page-size="formInline.page_size" :total="totalCount" @current-change="currentChange">
      </el-pagination>
    </div>
  </div>
</template>

<style>
  .demo-table-expand {
    font-size: 0;
  }
  .demo-table-expand label {
    width: 90px;
    color: #99a9bf;
  }
  .demo-table-expand .el-form-item {
    margin-right: 0;
    margin-bottom: 0;
    width: 50%;
  }
</style>

<script>
  export default {
    data() {
      return {
        showExtra: false,
        formInline: {
          code: '',
          appeal: true,
          page: 1,
          page_size: 20,
        },
        totalCount: 0,
        tableData: [],
        tableDataExtra: [],
        options: [
          {value: true, label: '申诉中'},
          {value: false, label: '已处理'},
        ]
      };
    },
    created () {
      this.getTrades();
    },
    methods: {
      statusShow(row) {
        switch(row.status) {
          case 0:
            return '待接受';
          case 1:
            return '已注资，待付款';
          case 2:
            return '已付款，待放行';
          case 3:
            return '完成';
          case 4:
            return '取消';
          case 5:
            return '自动取消，锁币';
          case 6:
            return '拒绝';
          case 7:
            return '自动取消，解除锁币';
        }
      },
      appealShow(row) {
        switch(row.appeal) {
          case 0:
            return '';
          case 1:
            return '买家申诉';
          case 2:
            return '卖家申诉';
          case 3:
            return '申诉过';
          case 4:
            return '已处理';
        }
      },
      getTrades() {
        this.$http.get("/kc/admin/appeal-trades", {"params": this.formInline})
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data !== null) {
            this.formInline.page = data.page_no;
            this.formInline.page_size = data.page_size;
            this.totalCount = data.total_count;
            this.tableData = data.list;
          }
        }, response => {
          this.$message.error(response.body);
        });
      },
      onSubmit() {
        this.formInline.page = 1;
        this.getTrades();
      },
      currentChange(pageNo) {
        this.formInline.page = pageNo;
        this.getTrades();
      },
      releaseTradePrompt(scope) {
        this.$prompt('密码', '确定放币给买家吗？', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            inputType: 'password'
          }).then(({ value }) => {
            this.releaseTrade(scope, value);
          })
      },
      releaseTrade(scope, value) {
        this.$http.post("/kc/admin/trade/release", {"code": scope.row.code, "password": value}, {"responseType": "json"})
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data !== null) {
            this.$message.info("操作成功");
            this.tableData[scope.$index].appeal = 4;
          }
        })
      },
      cancelTradePrompt(scope) {
        this.$prompt('密码', '确定放币给卖家吗？', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            inputType: 'password'
          }).then(({ value }) => {
            this.cancelTrade(scope, value);
          })
      },
      cancelTrade(scope, value) {
        this.$http.post("/kc/admin/trade/cancel", {"code": scope.row.code, "password": value}, {"responseType": "json"})
          .then(response => {
            var data = this.utilHelper.handle(this, response);
            if(data !== null) {
              this.$message.info("操作成功");
              this.tableData[scope.$index].appeal = 4;
            }
          })
      },
    }
  };
</script>
