export default function useDownload(){
    //通过文件下载url拿到对应的blob对象，然后通过a标签的download属性下载
    const downFile = (fileUrl:string,isDownload:boolean=true)=> {
        if(!isDownload){
            return
        }
        var link = document.createElement('a')
        // const apiUrl = import.meta.env.VITE_ONLINE_URL
        // fileUrl = apiUrl + fileUrl
        // fileUrl = fileUrl.replace("https://", "http://");
        // if (fileUrl.endsWith("/")) {
        //     fileUrl = fileUrl.substring(0, fileUrl.length - 1);
        // }
        link.href = fileUrl
        link.click()
    }
    return {downFile}

}