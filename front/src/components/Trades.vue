<template>
  <div>
    <el-form :inline="true" :model="formInline" class="demo-form-inline">
      <el-form-item label="">
        <el-input v-model="formInline.code" style="width:150px;" placeholder="交易号/广告号"></el-input>
        <el-input v-model="formInline.seller" style="width:150px;" placeholder="卖方"></el-input>
        <el-input v-model="formInline.buyer" style="width:150px;" placeholder="买方"></el-input>
        <el-select v-model="formInline.status" placeholder="请选择" clearable>
          <el-option key="0" label="待接受" value="0"></el-option>
          <el-option key="1" label="已注资，待付款" value="1"></el-option>
          <el-option key="2" label="已付款，待放行" value="2"></el-option>
          <el-option key="3" label="完成" value="3"></el-option>
          <el-option key="4" label="取消" value="4"></el-option>
          <el-option key="5" label="自动取消,锁币中" value="5"></el-option>
          <el-option key="6" label="拒绝" value="6"></el-option>
          <el-option key="7" label="自动取消,解除锁币" value="7"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">查询</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="tableData" border style="width: 100%; max-width:1149px;">
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
      <el-table-column prop="chat" label="操作" width="100">
        <template slot-scope="scope">
          <router-link :to="{ path: '/trade/' + scope.row.code + '/chats' }">聊天</router-link>
          <router-link :to="{ path: '/trade/' + scope.row.code + '/logs' }">历史</router-link>
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
          buyer: '',
          seller: '',
          code: '',
          status: null,
          page: 1,
          page_size: 20,
        },
        totalCount: 0,
        tableData: [],
        tableDataExtra: [],
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
            return '自动取消,锁币中';
          case 6:
            return '拒绝';
          case 7:
            return '自动取消,解除锁币';
          default:
            return row.status;
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
      getTrades(params, pageNo) {
        this.$http.get("/kc/admin/trades", {"params": this.formInline})
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
      }
    }
  };
</script>
