import {announcementCreateTypes} from "/@/api/announcement/types";

interface AnnouncementTableType extends TableType {
    data: announcementCreateTypes[];
    param: {
        name?: string;
        page: number;
        pagesize: number;
        total?: number;
    };
}

export interface AnnouncementState {
    tableData: AnnouncementTableType;
}

export interface roleMoveTypes {
    id: number;
    isUp: boolean;
}
