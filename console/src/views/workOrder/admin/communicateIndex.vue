<template>
  <div v-loading="state.tableData.loading">
    <div
        ref="dialogListRef"
        class="dialog-list"
        :style="dialogListStyle"
        v-infinite-scroll="onHandleCurrentChange"
        :infinite-scroll-disabled="disabled"
        :infinite-scroll-immediate="false"
    >
      <div class="dialog-item" v-for="(item, index) in state.tableData.data" :key="index">
        <div class="dialog-item-header">
          <div class="dialog-item-header-left">
            <div class="dialog-item-header-left-name">
              <div class="dialog-item-header-left-name-text">{{ item.from === 1 ? '用户' : '管理员' }}</div>
              <div class="dialog-item-header-left-name-time"></div>
            </div>
          </div>
          <div class="dialog-item-header-right">
            {{ dayjs.unix(item.createAt).format('YYYY-MM-DD HH:mm:ss') }}
          </div>
        </div>
        <div class="dialog-item-content">
          <div class="dialog-item-content-text" v-html="item.content"></div>
          <template v-if="item?.file&&item?.file.length > 0">
            <div class="dialog-item-content-title">附件（点击附件下载）：</div>
            <div class="dialog-item-content-img">
              <!--附件列表 ：-->
              <div
                  class="dialog-item-content-img-list"
                  v-for="(img, index) in item.file"
                  :key="index"
                  @click="downFile(getFile(img.fid, item.id, true), !isImage(img.fid))"
              >
                <!--                文件+文件名-->
                <div class="content">
                  <el-image class="img" v-if="isImage(img.fid)" :src="getFileType(img.fid)"
                            :preview-src-list="[getFile(img.fid, item.id)]" alt=""/>
                  <!--              文件缩略图，点击下载-->
                  <img v-else class="img" :src="getFileType(img.fid)" alt=""/>
                </div>
                <div class="fileName" :preview-src-list="[getFile(img.fid, item.id)]">{{ img.fid }}</div>
              </div>
            </div>
          </template>
        </div>
      </div>
      <p v-if="state.tableData.loading && !finishStatus" class="loading">加载中...</p>
      <p v-if="state.tableData.noMore && !finishStatus" class="noMore">没有更多了~</p>
    </div>
    <!-- 消息回复   -->
    <div v-if="dialog.type === 'reply' && finishStatus" class="dialog-reply">
      <div class="dialog-reply-content">
        <el-form ref="formRef" :model="ruleForm" label-width="80px" :rules="rules">
          <el-form-item label="回复内容" prop="content">
            <el-col class="mb10">
              <el-input type="textarea" :rows="2" v-model="ruleForm.content" placeholder="请输入回复内容"></el-input>
            </el-col>
          </el-form-item>
        </el-form>
        <el-row justify="space-between" class="w-[100%] mt-5">
          <el-col :span="10">
            <el-form-item label="文件上传">
              当前文件数量：{{ ruleForm.file.length }}
              <el-button class="btn" type="primary" size="small" @click="send">上传</el-button>
            </el-form-item>
          </el-col>
          <el-col :span="1" class="mx-2">
            <el-button class="send" type="primary" @click="submitForm(formRef)">发送</el-button>
          </el-col>
        </el-row>
      </div>
    </div>
    <el-dialog v-model="dialogVisible" title="简短描述" width="30%">
      <div v-html="currentInfo.describe"></div>
    </el-dialog>
    <file-upload ref="fileUploadRef" @send-success="sendSuccess"></file-upload>
  </div>
</template>

<script setup lang="ts">
import {computed, nextTick, onBeforeUnmount, onMounted, reactive, ref, watch} from 'vue';
import dayjs from 'dayjs';
import {getFileType} from '/@/utils/getFileType';
import useDownload from '/@/hooks/useDownload';
import {Session} from '/@/utils/storage';
import commonFunction from '/@/utils/commonFunction';
import {ElMessage, FormInstance, FormRules} from 'element-plus';
import FileUpload from '/@/views/workOrder/user/component/fileUpload.vue';
import {useAdminWorkOrderApi} from '/@/api/workOrder/admin';
import {useRouter} from 'vue-router';
import {useSocketStore} from '/@/stores/webstocket';
import {SendMessageTypeEnum} from '/@/data/enum';
import {storeToRefs} from "pinia";

const currentId = ref();
const dialog = reactive<dialogParams>({
  isShowDialog: false,
  type: 'reply',
  title: '沟通记录',
  submitTxt: '',
});
let ruleForm = ref<any>({
  title: '',
  content: '',
  filename: [] as any[],
  file: [] as any[],
});
const formRef = ref();
const dialogListRef = ref();
const disabled = computed(() => state.tableData.loading || state.tableData.noMore);
const rules = reactive<FormRules>({
  title: [
    {
      required: true,
      message: '请输入工单标题',
      trigger: 'blur',
    },
  ],
  content: [
    {
      required: true,
      message: '请输入工单内容',
      trigger: 'blur',
    },
  ],
});
// 判断当前是手机还是pc
const {isMobile} = commonFunction();
const currentMobile = computed(() => {
  return isMobile();
});
const dialogListStyle = computed(() => {
  let style = {};
  if (finishStatus.value) {
    style = {
      marginBottom: currentMobile.value ? '140px' : '100px',
      maxHeight: currentMobile.value ? 'calc(100vh - 252px) !important' : 'calc(100vh - 211px) !important',
    };
  } else {
    style = {
      marginBottom: '0px',
      maxHeight: 'calc(100vh - 102px) !important',
    };
  }
  return style;
});
// 滚动到底部
const isScrolling = ref(false); //用于判断用户是否在滚动
const scrollToBottom = () => {
  nextTick(() => {
    //注意要使用nexttick以免获取不到dom
    if (!isScrolling.value && dialogListRef.value) {
      dialogListRef.value.scrollTop = dialogListRef.value.scrollHeight;
    }
  });
};
const handleScroll = () => {
  const scrollContainer = dialogListRef.value;
  const scrollTop = scrollContainer.scrollTop;
  const scrollHeight = scrollContainer.scrollHeight;
  const offsetHeight = scrollContainer.offsetHeight;

  if (scrollTop + offsetHeight < scrollHeight) {
    // 用户开始滚动并在最底部之上，取消保持在最底部的效果
    isScrolling.value = true;
  } else {
    // 用户停止滚动并滚动到最底部，开启保持到最底部的效果
    isScrolling.value = false;
  }
};

const state = reactive<any>({
  tableData: {
    data: [],
    total: 0,
    loading: false,
    noMore: false,
    param: {
      page: 1,
      pagesize: 10,
    },
  },
});
const {downFile} = useDownload();
const apiUrl = ref();
const currentInfo = ref();
const dialogVisible = ref<boolean>(false);
const {isImage} = commonFunction();

const reset = () => {
  ruleForm.value = {
    title: '',
    content: '',
    filename: [],
    file: [],
  };
  fileLists.value = [];
};
const router = useRouter();
const fileUploadRef = ref();
const fileLists = ref<any>([]);
const finishStatus = computed(() => {
  return router.currentRoute.value.query.finishStatus;
});
const send = () => {
  fileUploadRef.value.openDialog(ruleForm.value, fileLists.value);
};
const sendSuccess = (data: any, fileList: any) => {
  ruleForm.value = data;
  fileLists.value = fileList;
};
// 发送文件
const submitForm = (formEl: FormInstance | undefined) => {
  if (currentId.value) {
    ruleForm.value.id = currentId.value;
  } else {
    ElMessage({
      type: 'error',
      message: '工单ID不存在',
    });
  }
  if (!formEl) return;
  formEl.validate((valid) => {
    if (valid) {
      useAdminWorkOrderApi()
          .userReplyWorkOrder(ruleForm.value)
          .then((res: any) => {
            if (res.code === "SUCCESS") {
              ElMessage({
                type: 'success',
                message: '发送成功',
              });
              reset();
              // getTableData();
            }
          });
    } else {
      return false;
    }
  });
};
const getFile = (fid: string, id: number, isDownload: boolean = false) => {
  return `${apiUrl.value}/admin/msg/allow-website/order/file?fid=${fid}&communicateID=${id}&xtoken=${Session.get('token')}
										&download=${isDownload}`;
};
const getTableData = async () => {
  await useAdminWorkOrderApi()
      .communicateList({orderID: currentId.value, ...state.tableData.param})
      .then((res: any) => {
        if (res.code === "SUCCESS") {
          infiniteScrollState.value = false;
          if (state.tableData.total == res.data.count) {
            return;
          } else {
            // 判断数据是否需要分页
            let currentPage = Math.ceil(res.data.count / state.tableData.param.pagesize);
            state.tableData.loading = false;
            // 根据分页和页码判断数据的起点和终点
            const start = (state.tableData.param.page - 1) * state.tableData.param.pagesize;
            const end = start + state.tableData.param.pagesize;
            if (res.data.communicate.length > state.tableData.param.pagesize) {
              state.tableData.data = [...state.tableData.data.slice(0, start), ...res.data.communicate, ...state.tableData.data.slice(end)];
            } else {
              state.tableData.data = [...state.tableData.data.slice(0, start), ...res.data.communicate];
            }
            state.tableData.total = state.tableData.data.length;
            state.tableData.count = res.data.count;
            state.tableData.noMore = state.tableData.data.length >= res.data.count;
            if (currentPage <= state.tableData.param.page) {
              onHandleCurrentChange();
            }
            scrollToBottom();
          }
        }
      });
};
// 分页改变
// 无限滚动的状态
const infiniteScrollState = ref(false);
const onHandleCurrentChange = () => {
  if (state.tableData.noMore || state.tableData.total == state.tableData.count) return;
  infiniteScrollState.value = true;
  state.tableData.param.page++;
  getTableData();
};

// pinia 中的orderData监听
const store = useSocketStore()
// pinia 中的orderData监听
const {orderData} = storeToRefs(store);
watch(
    orderData,
    (newVal: any) => {
      if (newVal.orderID == currentId.value) {
        state.tableData.data.push(newVal);
        state.tableData.total = state.tableData.data.length;
        scrollToBottom();
      }
    },
    {
      deep: true,
    }
);

const openDialog = async (row?: any) => {
  currentId.value = row.orderID;
  resetTableData();
  await getTableData();
  dialog.isShowDialog = true;
};
const resetTableData = () => {
  state.tableData.data = [];
  state.tableData.param.page = 1;
  state.tableData.total = 0;
  state.tableData.noMore = false;
};
onMounted(() => {
  apiUrl.value = import.meta.env.VITE_API_URL;
  // 监听滚动事件，判断用户滚动状态
  state.tableData.loading = true;
  router.currentRoute.value.query.id && openDialog({orderID: router.currentRoute.value.query.id});
  window.addEventListener('DOMContentLoaded', () => {
    scrollToBottom();
    dialogListRef.value.addEventListener('scroll', handleScroll);
  });
  useSocketStore().sendMessage(SendMessageTypeEnum.GET_ORDER, {
    data: JSON.stringify({
      orderID: router.currentRoute.value.query.id,
    })
  });
});
onBeforeUnmount(() => {
  window.removeEventListener('DOMContentLoaded', () => {
    dialogListRef.value.removeEventListener('scroll', handleScroll);
  });
});

// 暴露变量
defineExpose({
  openDialog,
});
</script>
<style lang="scss" scoped>
.dialog-list {
  overflow-y: auto;
  padding: 20px;

  .dialog-item {
    padding: 10px 0;
    border-bottom: 2px dashed #eee;

    .dialog-item-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 10px;

      .dialog-item-header-left {
        display: flex;
        align-items: center;

        .dialog-item-header-left-avatar {
          width: 40px;
          height: 40px;
          border-radius: 50%;
          overflow: hidden;

          .img {
            width: 100%;
            height: 100%;
          }
        }

        .dialog-item-header-left-name {
          .dialog-item-header-left-name-text {
            font-size: 14px;
            color: #333;
          }

          .dialog-item-header-left-name-time {
            font-size: 12px;
            color: #999;
          }
        }
      }

      .dialog-item-header-right {
        .el-button {
          padding: 0;
          font-size: 12px;
          color: #999;
        }
      }
    }

    .dialog-item-content {
      .dialog-item-content-text {
        font-size: 14px;
        color: #333;
        margin-bottom: 10px;
      }

      .dialog-item-content-title {
        margin: 5px 0;
      }

      .dialog-item-content-img {
        //display: flex;
        //flex-wrap: wrap;
        .dialog-item-content-img-list {
          cursor: pointer;
          display: flex;
          justify-content: flex-start;
          align-items: center;
          margin-right: 10px;
          margin-bottom: 10px;

          .fileName {
            width: 100%;
            font-size: 14px;
            color: #333;
            margin-bottom: 10px;
          }

          .content {
            .img {
              width: 20px;
              height: 20px;
              margin-right: 10px;
              margin-bottom: 10px;
            }
          }
        }
      }
    }
  }
}

.dialog-reply {
  width: 100%;
  position: absolute;
  bottom: 0;
  left: 0;
  border-top: 1px solid #eee;
  background-color: #fff;
  border-radius: 10px;

  .dialog-reply-content {
    width: 100%;
    height: 100%;
    background: rgba(232, 219, 219, 0.2);
    padding: 10px 20px;

    .el-form-item {
      margin-bottom: 20px;

      .el-input {
        width: 100%;
      }
    }

    .btn {
      margin-left: 20px;
    }

    .send {
      margin-left: 10px;
      width: 100%;
    }
  }
}

.noMore {
  text-align: center;
  font-size: 14px;
  color: #999;
  margin-top: 20px;
}

.loading {
  text-align: center;
  font-size: 14px;
  color: #999;
  margin-top: 20px;
}

//隐藏滚动条
::-webkit-scrollbar {
  display: none;
}
</style>
