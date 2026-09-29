export type SocketStore = {
    // 连接状态
    isConnected: boolean;
    // 消息内容
    message: string;
    // 重新连接错误
    reconnectError: boolean;
    // 心跳消息发送时间
    heartBeatInterval: number;
    // 心跳定时器
    heartBeatTimer: number;
    // socket实例
    socketInstance: null | any;
    // 工单数据监听
    orderData: any[];
    // token是否链接
    tokenStatus: boolean;
    // 发送token状态
    sendTokenStatus: boolean;
    // 工单完成状态
    orderDataFinishStatus: any[];
};

export type socketType = {
    $connect: () => void;
};