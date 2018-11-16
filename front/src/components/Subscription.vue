<template>
  <div>
    <el-form :inline="true" :model="formInline" class="demo-form-inline">
      <el-form-item>
        <el-input v-model="formInline.ident" clearable placeholder="UID/用户名/邮箱"></el-input>
      </el-form-item>
      <el-form-item>
        <el-input v-model="formInline.order" clearable placeholder="订单号"></el-input>
      </el-form-item>
      <el-form-item>
        <el-input v-model="formInline.base" clearable placeholder="付款方式"></el-input>
      </el-form-item>
      <el-select v-model="formInline.status" clearable placeholder="状态">
        <el-option key="0" label="未付款" value="0"></el-option>
        <el-option key="1" label="已付款，待验证" value="1"></el-option>
        <el-option key="2" label="完成" value="2"></el-option>
        <el-option key="3" label="取消" value="3"></el-option>
        <el-option key="4" label="自动取消" value="4"></el-option>
      </el-select>
      <el-form-item label="">
        <el-date-picker v-model="formInline.date" type="date" value-format="yyyy-MM-dd" placeholder="选择日期"></el-date-picker>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">查询</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="tableData" border style="width: 100%; max-width:1479px;" @expand-change="expandChange">
      <el-table-column type="expand">
        <template slot-scope="props">
          <el-form label-position="left" inline class="demo-table-expand">
            <el-form-item label="兑换比例" class="expand-normal">
              <span>1 {{ props.row.base }} : {{ props.row.exchange }} {{ props.row.currency }}</span>
            </el-form-item>
            <el-form-item label="对rmb价格" class="expand-normal">
              <span>{{ props.row.price }}</span>
            </el-form-item>
            <div v-for="item in submissionData[props.row.order]">
              <el-form-item label="HASH" class="demo-table-hash">
                <el-button type="text" @click="aHash(props.row.base, item.txid)">{{ item.txid }}</el-button>
                <span style="color: red;" v-if="item.warn1 || item.warn2">警告：{{ item.warn1 }} {{ item.warn2 }}</span>
              </el-form-item>
              <el-form-item label="状态" class="demo-table-status">
                <span>{{ submissionStatus(item.status) }}</span>
              </el-form-item>
              <el-form-item label="提交时间" class="demo-table-create-time">
                <span>{{ utilHelper.timeShow(item.create_time) }}</span>
              </el-form-item>
              <el-form-item label="备注" class="demo-table-remark">
                <span>{{ item.remark }}</span>
              </el-form-item>
              <el-form-item label="打币截图" class="demo-table-screenshot">
                <img :src='img(item.screenshot)' class="screenshot-img" @click="enlarged(item.screenshot)">
              </el-form-item>
            </div>
          </el-form>
        </template>
      </el-table-column>
      <el-table-column prop="uid" label="UID" width="100"></el-table-column>
      <el-table-column prop="order" label="订单号" width="170"></el-table-column>
      <el-table-column prop="base_amount" label="付款方式" width="150">
        <template slot-scope="scope">
          {{ scope.row.base_amount }} {{ scope.row.base }}
        </template>
      </el-table-column>
      <el-table-column prop="address" label="收款地址" width="400"></el-table-column>
      <el-table-column prop="currency_amount" label="认购数量" width="150">
        <template slot-scope="scope">
          {{ scope.row.currency_amount }} {{ scope.row.currency }}
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template slot-scope="scope">
          {{ orderStatus(scope.row.status) }}
        </template>
      </el-table-column>
      <el-table-column prop="create_time" label="时间" width="180">
        <template slot-scope="scope">
          {{ utilHelper.timeShow(scope.row.create_time) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180">
        <template slot-scope="scope">
          <div v-show="scope.row.status == 1">
            <el-button size="small" type="primary" @click="operation(scope, 2)">通过</el-button>
            <el-button size="small" type="danger" @click="operation(scope, 0)">不通过</el-button>
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
  .demo-table-expand label {
    width: 90px;
    color: #99a9bf;
  }
  .demo-table-expand .el-form-item {
    margin-right: 0;
    margin-bottom: 0;
  }
  .demo-table-expand .expand-normal {
    width: 20%;
  }

  .demo-table-expand .demo-table-hash {
    width: 15%;
  }
  .demo-table-expand .demo-table-status {
    width: 15%;
  }
  .demo-table-expand .demo-table-create-time {
    width: 20%;
  }
  .demo-table-expand .demo-table-remark {
    width: 25%;
  }
  .demo-table-expand .demo-table-screenshot {
    width: 20%;
  }

  .demo-table-expand .screenshot-img {
    width: 100px;
    height: 100px;
  }
  .msgbox {
    width: 800px;
  }
</style>

<script>
  export default {
    data() {
      return {
        formInline: {
          page: 1,
          page_size: 50,
          date: null,
          base: null,
          ident: null,
          status: null,
        },
        totalCount: 0,
        tableData: [],
        submissionData: {},
      };
    },
    created () {
      this.getSubscription();
    },
    methods: {
      orderStatus(status) {
        switch(status) {
          case 0:
            return "未付款";
          case 1:
            return "待验证";
          case 2:
            return "完成";
          case 3:
            return "取消";
          case 4:
            return "自动取消";
        }
      },
      submissionStatus(status) {
        switch(status) {
          case 0:
            return "未审核";
          case 1:
            return "未通过";
          case 2:
            return "完成";
          default:
            return status;
        }
      },
      aHash(base, hash) {
        switch(base) {
          case "BTC":
            var url = "https://btc.com/" + hash; break;
          case "USDT":
            var url = "https://omniexplorer.info/tx/" + hash; break;
          default:
            var url = "https://etherscan.io/tx/" + hash;
        }
        window.open(url,"_blank","toolbar=yes, location=yes, directories=no, status=no, menubar=yes, scrollbars=yes, resizable=no, copyhistory=yes, width=800, height=800");
      },
      img(filename) {
        return '/uphp/gcexserver/submission/' + filename;
      },
      enlarged(filename, deg) {
        const h = this.$createElement;
        if (!deg)
          deg = 0;
        var style = "width: 500px; max-height: 800px; transform: rotate(" + deg + "deg)";
        this.$msgbox({
          title: '消息',
          customClass: 'msgbox',
          message: h('img', {attrs: {src: this.img(filename), style: style}}),
          showCancelButton: false,
          showConfirmButton: true,
          confirmButtonText: '旋转',
        }).then(action => {
          if (action === 'confirm') {
            this.enlarged(filename, deg + 90);
          }
        });
      },
      getSubscription() {
        this.$http.get("/kc/admin/subscriptions", {"params": this.formInline})
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
      currentChange(pageNo) {
        this.formInline.page = pageNo;
        this.getSubscription();
      },
      onSubmit() {
        this.formInline.page = 1;
        this.getSubscription();
      },
      getSubmissions(order) {
        this.$http.get("/kc/admin/order/" + order + "/submissions/")
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if (data !== null) {
            var submissionData = {...this.submissionData};
            submissionData[order] = data;
            this.submissionData = submissionData;
          }
        })
      },
      expandChange(row, expandedRows) {
        this.getSubmissions(row["order"]);
      },
      operation(scope, status) {
        if(status == 2) {
          this.$prompt('请填写“审核通过”', '审核通过', {
            confirmButtonText: '确定',
            cancelButtonText: '取消'
          }).then(({ value }) => {
            if (value != "审核通过") {
              this.$message.error('确定入账请填写"审核通过"');
              return
            }
            this.handle(scope, status, value);
          })
        } else {
          this.$prompt('原因', '审核不通过', {
            confirmButtonText: '确定',
            cancelButtonText: '取消'
          }).then(({ value }) => {
            this.handle(scope, status, value);
          })
        }
      },
      handle(scope, status, remark) {
        this.$http.post("/kc/admin/order", {"id": scope.row.id, "status": status, "remark": remark}, {"responseType": "json"})
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data !== null) {
            this.tableData[scope.$index].status = status;
            this.tableData[scope.$index].status_name = this.orderStatus(status);
            this.$message.info("操作成功");
          }
        })
      }
    }
  };
</script>
