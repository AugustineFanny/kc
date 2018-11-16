<template>
  <div>
    <p v-if="formInline.uid">待认证数量：{{ totalCount }}</p>
    <el-table :data="tableData" border style="width: 100%; max-width:1261px;">
      <el-table-column prop="uid" label="UID" width="100"></el-table-column>
      <el-table-column prop="card" label="信息" width="190">
        <template slot-scope="scope">
          {{ scope.row.country }}<br>
          {{ scope.row.credential_type }}<br>
          {{ scope.row.name }}<br>
          {{ scope.row.card }}
        </template>
      </el-table-column>
      <el-table-column label="有效期" width="130">
        <template slot-scope="scope">
          {{utilHelper.dateShow(scope.row.start_date)}}<br>-<br>{{utilHelper.dateShow(scope.row.end_date)}}
        </template>
      </el-table-column>
      <el-table-column label="正面" width="180">
        <template slot-scope="scope">
          <img :src="img(scope.row.card_front)" @click="enlarged(scope.row.card_front)" class="img">
        </template>
      </el-table-column>
      <el-table-column prop="card_back" label="背面" width="180">
        <template slot-scope="scope">
          <img :src="img(scope.row.card_back)" @click="enlarged(scope.row.card_back)" class="img">
        </template>
      </el-table-column>
      <el-table-column prop="card_hold" label="手持" width="180">
        <template slot-scope="scope">
          <img :src="img(scope.row.card_hold)" @click="enlarged(scope.row.card_hold)" class="img">
        </template>
      </el-table-column>
      <el-table-column prop="status_name" label="状态" width="120"></el-table-column>
      <el-table-column label="操作" width="180">
        <template slot-scope="scope">
          <div v-show="scope.row.status == 1">
            <el-button size="small" @click="operation(scope, 3)">通过</el-button>
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
        tableData: [],
        transform: '0deg',
      };
    },
    created () {
      this.$http.get("/kc/admin/real-names", {"params": this.formInline})
      .then(response => {
        var data = this.utilHelper.handle(this, response);
        if(data !== null) {
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
        return '/uphp/gcexserver/name_auth/' + filename;
      },
      enlarged(filename, deg) {
        const h = this.$createElement;
        if (!deg)
          deg = 0;
        var style = "width: 500px; max-height: 800px; transform: rotate(" + deg + "deg)";
        this.$msgbox({
          title: '消息',
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
        this.$http.post("/kc/admin/name-auth", {"id": scope.row.id, "status": status, "desc": desc}, {"responseType": "json"})
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
