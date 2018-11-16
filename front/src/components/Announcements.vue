<template>
  <div>
    <el-row><el-button type="primary" @click="showCreate = true">创建公告</el-button></el-row>
    <el-row>
      <el-col :xs="24" :sm="24" :md="16" :lg="16">
        <el-table :data="tableData" border style="width: 100%; max-width: 881px;">
          <el-table-column prop="id" label="" width="80"></el-table-column>
          <el-table-column prop="title" label="标题" width="500"></el-table-column>
          <el-table-column label="滚动" width="100">
            <template slot-scope="scope">
              <el-tag :type="scope.row.important ? 'primary' : 'success'" close-transition>{{scope.row.important ? '不滚动' : '滚动'}}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="100">
            <template slot-scope="scope">
              <el-tag :type="scope.row.status ? 'primary' : 'success'" close-transition>{{scope.row.status ? '不显示' : '显示'}}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100">
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
        <el-form ref="change" :model="form" :rules="rules" label-width="80px">
          <el-form-item label="标题" prop="title" required>
            <el-input v-model="form.title" auto-complete="off"></el-input>
          </el-form-item>
          <el-form-item label="内容" prop="content" required>
            <el-input type="textarea" :rows="5" v-model="form.content"></el-input>
          </el-form-item>
          <el-form-item label="落款" prop="right_bottom">
            <el-input v-model="form.right_bottom"></el-input>
          </el-form-item>
          <el-form-item label="滚动" prop="important">
            <el-select v-model="form.important" placeholder="请选择">
              <el-option key="0" label="滚动" value="0" ></el-option>
              <el-option key="1" label="不滚动" value="1" ></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="状态" prop="status">
            <el-select v-model="form.status" placeholder="请选择">
              <el-option key="0" label="显示" value="0" ></el-option>
              <el-option key="1" label="不显示" value="1" ></el-option>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="onChange">提交</el-button>
            <el-button @click="onCancel">取消</el-button>
          </el-form-item>
        </el-form>
      </el-col>
    </el-row>
    <el-dialog title="创建公告" :visible.sync="showCreate">
      <el-form ref="create" :model="form2" :rules="rules2" label-width="100px" class="demo-ruleForm">
        <el-form-item label="标题" prop="title" required>
          <el-input v-model="form2.title" auto-complete="off"></el-input>
        </el-form-item>
        <el-form-item label="内容" prop="content" required>
          <el-input type="textarea" :rows="5" v-model="form2.content"></el-input>
        </el-form-item>
        <el-form-item label="落款" prop="right_bottom">
            <el-input v-model="form2.right_bottom"></el-input>
          </el-form-item>
        <el-form-item label="滚动" prop="important">
          <el-select v-model="form2.important" placeholder="请选择">
            <el-option key="0" label="是" value="0" ></el-option>
            <el-option key="1" label="否" value="1" ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="form2.status" placeholder="请选择">
            <el-option key="0" label="显示" value="0" ></el-option>
            <el-option key="1" label="不显示" value="1" ></el-option>
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
        formInline: {
          page: 1,
          page_size: 20,
        },
        totalCount: 0,
        tableData: [],
        show: false,
        form: {
          title: '',
          content: '',
          right_bottom: '',
          important: '0',
          status: '0',
        },
        rules: {

        },
        showCreate: false,
        form2: {
          title: '',
          content: '',
          right_bottom: '',
          important: '0',
          status: '0',
        },
        rules2: {

        }
      };
    },
    created () {
      this.getAnnouncements();
	  },
    methods: {
      getAnnouncements() {
        this.$http.get("/kc/admin/anncs", {"params": this.formInline})
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
        this.form.important = this.form.important.toString();
        this.form.status = this.form.status.toString();
      },
      onChange() {
        this.$refs["change"].validate((valid) => {
          if (valid) {
            this.$resource("/kc/admin/annc/{id}")
            .save({id: this.form.id}, {
              title: this.form.title,
              content: this.form.content,
              right_bottom: this.form.right_bottom,
              important: parseInt(this.form.important),
              status: parseInt(this.form.status),
            })
            .then(response => {
              var data = this.utilHelper.handle(this, response);
              if(data !== null) {
                this.show = false;
                this.$refs["change"].resetFields();
                this.$message.success("操作成功");
                this.getAnnouncements();
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
            this.$http.post("/kc/admin/anncs", {
              title: this.form2.title,
              content: this.form2.content,
              right_bottom: this.form2.right_bottom,
              important: parseInt(this.form2.important),
              status: parseInt(this.form2.status),
            })
            .then(response => {
              var data = this.utilHelper.handle(this, response);
              if(data !== null) {
                this.$message.success("操作成功");
                this.$refs["create"].resetFields();
                this.showCreate = false;
                this.getAnnouncements();
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
      },
      currentChange(pageNo) {
        this.formInline.page = pageNo;
        this.getAnnouncements();
      },
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
