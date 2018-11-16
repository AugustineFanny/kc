<template>
  <div>
    <el-row><el-button type="primary" @click="showCreate = true">创建横幅</el-button></el-row>
    <el-form :inline="true" :model="formInline" class="demo-form-inline">
      <el-select v-model="formInline.status" clearable placeholder="全部">
        <el-option key="-1" label="全部" value=""></el-option>
        <el-option key="0" label="正常" value="0"></el-option>
        <el-option key="1" label="不显示" value="1"></el-option>
      </el-select>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">查询</el-button>
      </el-form-item>
    </el-form>
    <el-row>
      <el-col :xs="24" :sm="24" :md="16" :lg="16">
        <el-table :data="tableData" border style="width: 100%; max-width: 841px;">
          <el-table-column prop="id" label="" width="80"></el-table-column>
          <el-table-column label="大图" width="180">
            <template slot-scope="scope">
              <img :src="img(scope.row.big_banner)" @click="enlarged(scope.row.big_banner)" class="img">
            </template>
          </el-table-column>
          <el-table-column label="中图" width="180">
            <template slot-scope="scope">
              <img :src="img(scope.row.medium_banner)" @click="enlarged(scope.row.medium_banner)" class="img">
            </template>
          </el-table-column>
          <el-table-column label="小图" width="180">
            <template slot-scope="scope">
              <img :src="img(scope.row.small_banner)" @click="enlarged(scope.row.small_banner)" class="img">
            </template>
          </el-table-column>
          <el-table-column prop="src" label="链接" width="120"></el-table-column>
          <el-table-column prop="status" label="状态" width="120">
            <template slot-scope="scope">
              {{ showStatus(scope.row.status) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100" fixed="right">
            <template slot-scope="scope">
              <el-button @click="handleClick(scope.row)" type="text" size="small">编辑</el-button>
            </template>
          </el-table-column>
        </el-table>
        <div style="position:absolute;">
          <el-pagination layout="total, prev, pager, next, jumper" :page-size="formInline.page_size" :total="totalCount" @current-change="currentChange">
          </el-pagination>
        </div>
      </el-col>
      <el-col :xs="24" :sm="24" :md="8" :lg="8" v-show="show">
        <el-form ref="change" :model="form" :rules="rules" label-width="100px">
          <el-form-item label="大图" prop="bigBanner">
            <el-upload
              class="picture-uploader"
              action="https://jsonplaceholder.typicode.com/posts/"
              :show-file-list="false"
              :http-request="updateBigBannerRequest">
              <img v-if="form.bigBannerUrl" :src="form.bigBannerUrl" class="picture">
              <i v-else class="el-icon-plus picture-uploader-icon"></i>
            </el-upload>
          </el-form-item>
          <el-form-item label="中图" prop="mediumBanner">
            <el-upload
              class="picture-uploader"
              action="https://jsonplaceholder.typicode.com/posts/"
              :show-file-list="false"
              :http-request="updateMediumBannerRequest">
              <img v-if="form.mediumBannerUrl" :src="form.mediumBannerUrl" class="picture">
              <i v-else class="el-icon-plus picture-uploader-icon"></i>
            </el-upload>
          </el-form-item>
          <el-form-item label="小图" prop="smallBanner">
            <el-upload
              class="picture-uploader"
              action="https://jsonplaceholder.typicode.com/posts/"
              :show-file-list="false"
              :http-request="updateSmallBannerRequest">
              <img v-if="form.smallBannerUrl" :src="form.smallBannerUrl" class="picture">
              <i v-else class="el-icon-plus picture-uploader-icon"></i>
            </el-upload>
          </el-form-item>
          <el-form-item label="链接" prop="src">
            <el-input v-model="form.src"></el-input>
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="form.status" placeholder="请选择">
              <el-option key="0" label="正常" value="0"></el-option>
              <el-option key="1" label="不显示" value="1"></el-option>
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
        <el-form-item label="大图" prop="bigBanner" required>
          <el-upload
            class="picture-uploader"
            action="https://jsonplaceholder.typicode.com/posts/"
            :show-file-list="false"
            :http-request="bigBannerRequest">
            <img v-if="form2.bigBannerUrl" :src="form2.bigBannerUrl" class="picture">
            <i v-else class="el-icon-plus picture-uploader-icon"></i>
          </el-upload>
        </el-form-item>
        <el-form-item label="中图" prop="mediumBanner" required>
          <el-upload
            class="picture-uploader"
            action="https://jsonplaceholder.typicode.com/posts/"
            :show-file-list="false"
            :http-request="mediumBannerRequest">
            <img v-if="form2.mediumBannerUrl" :src="form2.mediumBannerUrl" class="picture">
            <i v-else class="el-icon-plus picture-uploader-icon"></i>
          </el-upload>
        </el-form-item>
        <el-form-item label="小图" prop="smallBanner" required>
          <el-upload
            class="picture-uploader"
            action="https://jsonplaceholder.typicode.com/posts/"
            :show-file-list="false"
            :http-request="smallBannerRequest">
            <img v-if="form2.smallBannerUrl" :src="form2.smallBannerUrl" class="picture">
            <i v-else class="el-icon-plus picture-uploader-icon"></i>
          </el-upload>
        </el-form-item>
        <el-form-item label="链接" prop="src">
          <el-input v-model="form2.src"></el-input>
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
        formInline: {
          status: "0",
          page: 1,
          page_size: 20,
        },
        tableData: [],
        show: false,
        form: {
          id: null,
          src: "",
          status: "",
          bigBannerUrl: "",
          mediumBannerUrl: "",
          smallBannerUrl: "",
          dialogVisible: false,
          status: "",
        },
        rules: {

        },
        showCreate: false,
        form2: {
          src: "",
          status: "",
          bigBannerUrl: "",
          mediumBannerUrl: "",
          smallBannerUrl: "",
          dialogVisible: false
        },
        rules2: {

        }
      };
    },
    created () {
      this.getBanners();
    },
    methods: {
      onSubmit() {
        this.formInline.page = 1;
        this.getBanners();
      },
      currentChange(pageNo) {
        this.formInline.page = pageNo;
        this.getBanners();
      },
      img(filename) {
        return '/uphp/gcexserver/banner/' + filename;
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
            console.log(deg + 90)
            this.enlarged(filename, deg + 90);
          }
        });
      },
      showStatus(status) {
        switch(status) {
          case 0:
            return "正常";
          case 1:
            return "不显示";
        }
      },
      getBanners() {
        this.$http.get("/kc/admin/banners", {"params": this.formInline})
        .then(response => {
          var data = this.utilHelper.handle(this, response);
          if(data != null) {
            this.formInline.page = data.page_no;
            this.formInline.page_size = data.page_size;
            this.totalCount = data.total_count;
            this.tableData = data.list;
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
        this.form.status = this.form["status"].toString();
        this.form.bigBanner = null;
        this.form.mediumBanner = null;
        this.form.smallBanner = null;
        this.form.bigBannerUrl = null;
        this.form.mediumBannerUrl = null;
        this.form.smallBannerUrl = null;
      },
      bigBannerRequest(obj) {
        this.form2.bigBanner = obj.file;
        this.form2.bigBannerUrl = URL.createObjectURL(obj.file);
      },
      mediumBannerRequest(obj) {
        this.form2.mediumBanner = obj.file;
        this.form2.mediumBannerUrl = URL.createObjectURL(obj.file);
      },
      smallBannerRequest(obj) {
        this.form2.smallBanner = obj.file;
        this.form2.smallBannerUrl = URL.createObjectURL(obj.file);
      },
      updateBigBannerRequest(obj) {
        this.form.bigBanner = obj.file;
        this.form.bigBannerUrl = URL.createObjectURL(obj.file);
      },
      updateMediumBannerRequest(obj) {
        this.form.mediumBanner = obj.file;
        this.form.mediumBannerUrl = URL.createObjectURL(obj.file);
      },
      updateSmallBannerRequest(obj) {
        this.form.smallBanner = obj.file;
        this.form.smallBannerUrl = URL.createObjectURL(obj.file);
      },
      onChange() {
        this.$refs["change"].validate((valid) => {
          if (valid) {
            var formData = new FormData();
            if(this.form.bigBanner != null) {
              formData.append("big_banner", this.form.bigBanner);
            }
            if(this.form.bigBanner != null) {
              formData.append("medium_banner", this.form.mediumBanner);
            }
            if(this.form.bigBanner != null) {
              formData.append("small_banner", this.form.smallBanner);
            }
            formData.append("src", this.form.src);
            formData.append("status", this.form.status);
            this.$resource("/kc/admin/banner/{id}")
            .save({id: this.form.id}, formData)
            .then(response => {
              var data = this.utilHelper.handle(this, response);
              if(data !== null) {
                this.$message.success("操作成功");
                this.$refs["change"].resetFields();
                this.getBanners();
                this.show = false;
                this.form.bigBannerUrl = "";
                this.form.mediumBannerUrl = "";
                this.form.smallbannerUrl = "";
                this.form.dialogVisible = false;
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
            var formData = new FormData();
            formData.append("big_banner", this.form2.bigBanner);
            formData.append("medium_banner", this.form2.mediumBanner);
            formData.append("small_banner", this.form2.smallBanner);
            formData.append("src", this.form2.src);
            this.$http.post("/kc/admin/banners", formData)
            .then(response => {
              var data = this.utilHelper.handle(this, response);
              if(data !== null) {
                this.$message.success("操作成功");
                this.$refs["create"].resetFields();
                this.showCreate = false;
                this.getBanners();
                this.form2.bigBannerUrl = "";
                this.form2.mediumBannerUrl = "";
                this.form2.smallBannerUrl = "";
                this.form2.dialogVisible = false;
              }
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
  .img {
    width: 150px;
    height: 150px;
  }
  .el-row {
    margin-bottom: 20px;
    &:last-child {
      margin-bottom: 0;
    }
  }
  .picture-uploader .el-upload {
    border: 1px dashed #d9d9d9;
    border-radius: 6px;
    cursor: pointer;
    position: relative;
    overflow: hidden;
  }
  .picture-uploader .el-upload:hover {
    border-color: #409EFF;
  }
  .picture-uploader-icon {
    font-size: 28px;
    color: #8c939d;
    width: 178px;
    height: 178px;
    line-height: 178px;
    text-align: center;
  }
  .picture {
    width: 178px;
    height: 178px;
    display: block;
  }
</style>
