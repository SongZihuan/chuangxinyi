export interface applicationTypes {
	name: string;
	webID: number | string;
	url: string;
	icon: string;
	status: number;
	describe?: string;
}

export interface searchTypes {
	name?: string;
	page: number;
	pagesize: number;
	total: number;
}

export interface applicationDeltypes {
	id: number;
}

export interface applicationMoveTypes {
	id: number;
	isUp: boolean;
}

export interface roleActionTypes {
	applicationID: number;
	roleID: number;
}
