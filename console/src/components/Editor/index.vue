<template>
  <div class="editor-container">
    <Toolbar :editor="editorRef" :mode="mode"/>
    <Editor
        :mode="mode"
        :defaultConfig="state.editorConfig"
        :style="{ height }"
        v-model="state.editorVal"
        @onCreated="handleCreated"
        @onChange="handleChange"
    />
  </div>
</template>

<script setup lang="ts" name="wngEditor">
import '@wangeditor/editor/dist/css/style.css';
import {reactive, shallowRef, watch, onBeforeUnmount, ref, onMounted} from 'vue';
import {IDomEditor} from '@wangeditor/editor';
import {Toolbar, Editor} from '@wangeditor/editor-for-vue';
import {ElMessage} from "element-plus";
import {Session} from "/@/utils/storage";
import {useBaseApi} from "/@/api/base";
import {uuid} from "@pureadmin/utils";
import dayjs from "dayjs";

// 定义父组件传过来的值
const props = defineProps({
  // 是否禁用
  disable: {
    type: Boolean,
    default: () => false,
  },
  // 内容框默认 placeholder
  placeholder: {
    type: String,
    default: () => '请输入内容...',
  },
  // https://www.wangeditor.com/v5/getting-started.html#mode-%E6%A8%A1%E5%BC%8F
  // 模式，可选 <default|simple>，默认 default
  mode: {
    type: String,
    default: () => 'default',
  },
  // 高度
  height: {
    type: String,
    default: () => '310px',
  },
  // 双向绑定，用于获取 editor.getHtml()
  getHtml: String,
  // 双向绑定，用于获取 editor.getText()
  getText: String,
});

// 定义子组件向父组件传值/事件
const emit = defineEmits(['update:getHtml', 'update:getText']);
// 定义变量内容
const editorRef = shallowRef();
const state = reactive({
  editorConfig: {
    placeholder: props.placeholder,
    // 所有的菜单配置，都要在 MENU_CONF 属性下
    MENU_CONF: {
      // 图片上传
      uploadImage: {
        async customUpload(file: File, insertFn: InsertFnType) {  // TS 语法
          ElMessage({
            message: '图片正在上传中,请耐心等待',
            duration: 0,
            customClass: 'uploadTip',
            iconClass: 'el-icon-loading',
            showClose: true
          });
          let fromData = new FormData();
          fromData.append('file', file);
          //  文件名加密
          let fileName = file.name;
          let suffixName =  fileName.split('.').pop();
          let name =  `${fileName.split('.')[0]}-${uuid()}-${dayjs().format('YYYYMMDDHHmmss')}.${suffixName}`;
          fromData.append('fid', name);
          useBaseApi().uploadFile(fromData).then((res: any) => {
            if(res.code==0){
              ElMessage.closeAll();
              ElMessage.success({
                message: `${name} 上传成功`
              });
              insertFn(getFile(name), name, getFile(name));
            }else{
              ElMessage.error({
                message: `${name} 上传失败，请重新尝试`
              });
            }
          }).catch(() => {
            ElMessage.error({
              message: `${file.name} 上传失败，请重新尝试`
            });
          })
        },
      },
      // 视频上传
      uploadVideo: {
        async customInsert(file: File, insertFn: InsertFnType) {  // TS 语法
          ElMessage({
            message: '图片正在上传中,请耐心等待',
            duration: 0,
            customClass: 'uploadTip',
            iconClass: 'el-icon-loading',
            showClose: true
          });
          let fromData = new FormData();
          fromData.append('file', file);
          //  文件名加密
          let fileName = file.name;
          let suffixName =  fileName.split('.').pop();
          let name =  `${fileName.split('.')[0]}-${uuid()}-${dayjs().format('YYYYMMDDHHmmss')}.${suffixName}`;
          fromData.append('fid', name);
          useBaseApi().uploadFile(fromData).then((res: any) => {
            if(res.code==0){
              ElMessage.closeAll();
              ElMessage.success({
                message: `${name} 上传成功`
              });
              insertFn(getFile(name), name, getFile(name));
            }else{
              ElMessage.error({
                message: `${name} 上传失败，请重新尝试`
              });
            }
          }).catch(() => {
            ElMessage.error({
              message: `${file.name} 上传失败，请重新尝试`
            });
          })
        },
      }
    }
  },
  editorVal: props.getHtml,
});
const apiUrl = ref();
const getFile = (fid: string) => {
  return `${apiUrl.value}/public/ui/file?fid=${fid}&xtoken=${Session.get('token')}&download=false`;
};

// 编辑器回调函数
const handleCreated = (editor: IDomEditor) => {
  editorRef.value = editor;
};
// 编辑器内容改变时
const handleChange = (editor: IDomEditor) => {
  emit('update:getHtml', editor.getHtml());
  emit('update:getText', editor.getText());
};
// 页面销毁时
onBeforeUnmount(() => {
  const editor = editorRef.value;
  if (editor == null) return;
  editor.destroy();
});
// 监听是否禁用改变
// https://gitee.com/lyt-top/vue-next-admin/issues/I4LM7I
watch(
    () => props.disable,
    (bool) => {
      const editor = editorRef.value;
      if (editor == null) return;
      bool ? editor.disable() : editor.enable();
    },
    {
      deep: true,
    }
);
// 监听双向绑定值改变，用于回显
watch(
    () => props.getHtml,
    (val) => {
      state.editorVal = val;
    },
    {
      deep: true,
    }
);
onMounted(() => {
  apiUrl.value = import.meta.env.VITE_API_URL;
});
</script>
