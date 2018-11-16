<template>
  <div>
    待认证数量：{{ totalCount }}
    <el-table :data="tableData" border style="width: 100%;">
      <el-table-column prop="name" label="姓名" width="130"></el-table-column>
      <el-table-column label="用户名/邮箱/手机/生日" width="200">
        <template slot-scope="scope">
          {{ scope.row.username }}<br>
          {{ scope.row.email }}<br>
          {{ scope.row.mobile }}<br>
          {{ scope.row.birthday }}
        </template>
      </el-table-column>
      <el-table-column label="国/省/城市/街道" width="180">
        <template slot-scope="scope">
          {{ scope.row.country }}<br>
          {{ scope.row.province }}<br>
          {{ scope.row.city }}<br>
          {{ scope.row.street }}
        </template>
      </el-table-column>
      <el-table-column label="邮编/身份类型/收入来源" width="200">
        <template slot-scope="scope">
          {{ scope.row.post_code }}<br>
          {{ scope.row.identity_document }}<br>
          {{ scope.row.funds_source }}
        </template>
      </el-table-column>
      <el-table-column label="正面" width="180">
        <template slot-scope="scope">
          <img :src="img(scope.row.photo_id_front)" @click="enlarged(scope.row.photo_id_front)" class="img">
        </template>
      </el-table-column>
      <el-table-column label="背面" width="180">
        <template slot-scope="scope">
          <img :src="img(scope.row.photo_id_back)" @click="enlarged(scope.row.photo_id_back)" class="img">
        </template>
      </el-table-column>
      <el-table-column label="手持" width="180">
        <template slot-scope="scope">
          <img :src="img(scope.row.photo_id_hold)" @click="enlarged(scope.row.photo_id_hold)" class="img">
        </template>
      </el-table-column>
      <el-table-column label="银行账号照片" width="180">
        <template slot-scope="scope">
          <img :src="img(scope.row.bank_account_photo)" @click="enlarged(scope.row.bank_account_photo)" class="img">
        </template>
      </el-table-column>
      <el-table-column prop="status_name" label="状态" width="120"></el-table-column>
      <el-table-column label="操作" width="100" fixed="right">
        <template slot-scope="scope">
          <div v-show="scope.row.status == 1">
            <el-button size="small" @click="operation(scope, 3)">通过</el-button> <br><br>
            <el-button size="small" type="danger" @click="operation(scope, 2)">不通过</el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
  export default {
    data() {
      return {
        formInline: {
          uid: this.$route.params.uid
        },
        totalCount: 0,
        tableData: []
      };
    },
    created () {
      this.$http.get("/kc/admin/kycs", {"params": this.formInline})
      .then(response => {
        var data = this.utilHelper.handle(this, response);
        if(data !== null) {
          console.log(data)
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
    methods: {
      img(filename) {
        return '/uphp/gcexserver/kyc/' + filename;
      },
      enlarged(filename, deg) {
        const h = this.$createElement;
        if (!deg)
          deg = 0;
        var style = "width: 500px; max-height: 800px; transform: rotate(" + deg + "deg)";
        this.$msgbox({
          title: '消息',
          message: h('img', {attrs: {src: this.img(filename), width: "500px", style: style}}),
          showCancelButton: false,
          showConfirmButton: true,
          confirmButtonText: '旋转',
        }).then(action => {
          if (action === 'confirm') {
            console.log(deg + 90)
            this.enlarged(filename, deg + 90);
          }
        });
      },
      statusName(status) {
        if(status == 0) {
          return "未提交认证";
        } else if (status == 1) {
          return "待认证";
        } else if (status == 2) {
          return "认证未通过";
        } else if (status == 3) {
          return "认证通过";
        } else {
          return "未知状态：" + status;
        }
      },
      operation(scope, status) {
        if(status == 2) {
          this.$prompt('原因', '认证不通过', {
            confirmButtonText: '确定',
            cancelButtonText: '取消'
          }).then(({ value }) => {
            this.handle(scope, status, value);
          })
        } else {
          this.handle(scope, status, "");
        }
      },
      handle(scope, status, desc) {
        this.$http.post("/kc/admin/kyc-auth", {"id": scope.row.id, "status": status, "desc": desc}, {"responseType": "json"})
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data !== null) {
            this.tableData[scope.$index].status = status;
            this.tableData[scope.$index].status_name = this.statusName(status);
            this.$message.info("操作成功");
          }
        })
      }
    }
  };
</script>

<style type="text/css">
  .img {
    width: 150px;
    height: 150px;
  }
  .el-message-box {
    width: auto;
    min-width: 420px;
  }
</style>
