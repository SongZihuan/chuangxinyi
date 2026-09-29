/**
 * 根据类型返回对应的文件图标
 * @param type
 */
import txt from '/@/assets/fileIcon/txt_document_extension_file_format_icon.png';
import word from '/@/assets/fileIcon/doc_document_extension_file_format_icon.png';
import excel from '/@/assets/fileIcon/xls_document_extension_file_format_icon.png';
import ppt from '/@/assets/fileIcon/extension_file_name_ppt_icon.png';
import pdf from '/@/assets/fileIcon/pdf_document_extension_file_format_icon.png';
import video from '/@/assets/fileIcon/mp4_document_extension_file_format_icon.png';
import audio from '/@/assets/fileIcon/mp3_document_documents_file_format_icon.png';
import zip from '/@/assets/fileIcon/zip_file_icon.png';
import gif from '/@/assets/fileIcon/2276047_document_extension_format_gif_paper_icon.png';
import jpg from '/@/assets/fileIcon/2276087_document_extension_format_jpg_paper_icon.png';
import png from '/@/assets/fileIcon/2276091_document_extension_format_paper_png file_icon.png';
import other from '/@/assets/fileIcon/file_document_icon.png';
import dayjs from "dayjs";
export function getFileType(type:string) {
    const fileTypeMap = {
        txt: txt,
        word: word,
        excel: excel,
        ppt: ppt,
        pdf: pdf,
        video: video,
        audio:audio,
        zip: zip,
        gif:gif,
        jpg:jpg,
        png:png,
        jpeg:jpg,
        other: other
    };
    const fileType = type.toLowerCase();
    if (fileType.includes('folder')) {
        return fileTypeMap.folder;
    }
    if (fileType.includes('word')) {
        return fileTypeMap.word;
    }
    if (fileType.includes('excel')) {
        return fileTypeMap.excel;
    }
    if (fileType.includes('ppt')) {
        return fileTypeMap.ppt;
    }
    if (fileType.includes('pdf')) {
        return fileTypeMap.pdf;
    }
    if (fileType.includes('image')) {
        return fileTypeMap.image;
    }
    if (fileType.includes('video')) {
        return fileTypeMap.video;
    }
    if (fileType.includes('audio')) {
        return fileTypeMap.audio;
    }
    if (fileType.includes('zip')) {
        return fileTypeMap.zip;
    }
    if (fileType.includes('gif')) {
        return fileTypeMap.gif;
    }
    if (fileType.includes('jpg') || fileType.includes('jpeg')) {
        return fileTypeMap.jpg;
    }
    if (fileType.includes('png')) {
        return fileTypeMap.png;
    }
    return fileTypeMap.other;
}

/**
 * 加密文件名
 */
export function encryptFileName(fileName:string) {
    let suffixName = fileName.split('.')[1];
    if (!suffixName) {
        return fileName;
    }
    let name =  `${fileName.split('.')[0]}-${dayjs().format('YYYYMMDDHHmmss')}.${suffixName}`;
    return name;
}