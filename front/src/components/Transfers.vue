<template>
  <div>
    <el-form :inline="true" :model="formInline" class="demo-form-inline">
      <el-input v-model="formInline.uid" style="width:200px;" clearable placeholder="uid"></el-input>
      <el-input v-model="formInline.currency" style="width:200px;" clearable placeholder="currency"></el-input>
      <el-select v-model="formInline.direction" clearable placeholder="方向">
        <el-option key="0" label="充" value="0"></el-option>
        <el-option key="1" label="提" value="1"></el-option>
        <el-option key="2" label="站内转" value="2"></el-option>
        <el-option key="3" label="锁仓转" value="3"></el-option>
      </el-select>
      <el-select v-model="formInline.status" clearable placeholder="状态">
        <el-option key="0" label="确认中" value="0"></el-option>
        <el-option key="1" label="异常" value="1"></el-option>
        <el-option key="2" label="成功" value="2"></el-option>
        <el-option key="3" label="待审核" value="3"></el-option>
        <el-option key="4" label="审核未通过" value="4"></el-option>
      </el-select>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">查询</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="tableData" border style="width: 100%;">
      <el-table-column prop="id" label="ID" width="100"></el-table-column>
      <el-table-column prop="uid" label="UID" width="100"></el-table-column>
      <el-table-column prop="currency" label="币" width="100"></el-table-column>
      <el-table-column prop="amount" label="数量" width="120"></el-table-column>
      <el-table-column label="费用" width="120">
        <template slot-scope="scope">
          {{ scope.row.fee }} {{ scope.row.fee_currency }}
        </template>
      </el-table-column>
      <el-table-column prop="status_name" label="状态" width="120"></el-table-column>
      <el-table-column prop="hash" label="TxHash" width="200">
        <template slot-scope="scope"><span class="ellipsis">{{scope.row.hash}}</span></template>
      </el-table-column>
      <el-table-column prop="from" label="源地址" width="200">
        <template slot-scope="scope"><span class="ellipsis">{{scope.row.from}}</span></template>
      </el-table-column>
      <el-table-column prop="direction" label="方向" width="80">
        <template slot-scope="scope">
          <el-tag :type="scope.row.direction ? 'primary' : 'success'" close-transition>{{directionName(scope.row.direction)}}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="to" label="目标地址" width="200">
        <template slot-scope="scope"><span class="ellipsis">{{scope.row.to}}</span></template>
      </el-table-column>
      <el-table-column prop="create_time" label="申请时间" width="180">
        <template slot-scope="scope">{{ utilHelper.timeShow(scope.row.create_time) }}</template>
      </el-table-column>
      <el-table-column prop="create_time" label="到账时间" width="180">
        <template slot-scope="scope">{{ utilHelper.timeShow(scope.row.check_time) }}</template>
      </el-table-column>
      <el-table-column prop="desc" label="描述" width="300">
        <template slot-scope="scope"><span class="ellipsis">{{scope.row.desc}}</span></template>
      </el-table-column>
      <el-table-column prop="remark" label="备注" width="300">
        <template slot-scope="scope"><span class="ellipsis">{{scope.row.remark}}</span></template>
      </el-table-column>
    </el-table>
    <div style="position:absolute;">
      <el-pagination layout="total, prev, pager, next, jumper" :page-size="formInline.page_size" :total="totalCount" @current-change="currentChange">
      </el-pagination>
    </div>
  </div>
</template>

<script>
  export default {
    data() {
      return {
        formInline: {
          uid: '',
          currency: '',
          direction: '',
          status: null,
          page: 1,
          page_size: 20,
        },
        totalCount: 0,
        tableData: []
      };
    },
    created () {
      this.getTransfers();
    },
    methods: {
      directionName(direction) {
        if(direction == 0) {
          return "充";
        } else if (direction == 1) {
          return "提";
        } else if (direction == 2) {
          return "站内转";
        } else if (direction == 3) {
          return "锁仓转";
        } else {
          return "未知状态：" + status;
        }
      },
      statusName(status) {
        if(status == 0) {
          return "确认中";
        } else if (status == 1) {
          return "异常";
        } else if (status == 2) {
          return "成功";
        } else if (status == 3) {
          return "待审核";
        } else if (status == 4) {
          return "审核未通过";
        } else {
          return "未知状态：" + status;
        }
      },
      getTransfers() {
        this.$http.get("/kc/admin/transfers", {"params": this.formInline})
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data !== null) {
            this.formInline.page = data.page_no;
            this.formInline.page_size = data.page_size;
            this.totalCount = data.total_count;
            data.list.forEach((row) => {
              row.status_name = this.statusName(row.status);
            })
            this.tableData = data.list;
          }
        }, response => {
          this.$message.error(response.body);
        });
      },
      currentChange(pageNo) {
        this.formInline.page = pageNo;
        this.getTransfers();
      },
      onSubmit() {
        this.formInline.page = 1;
        this.getTransfers();
      },
      sortChange(a) {
        this.formInline._sort = a.prop;
        this.formInline._order = a.order;
        this.formInline.page = 1;
        this.getTransfers();
      }
    }
  };
</script>

<style type="text/css">
  .ellipsis {
    white-space: nowrap;
  }
</style>
