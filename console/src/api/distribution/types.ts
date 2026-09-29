type valTypes = {
	level: number;
	pre: number;
};

export interface fromType {
	val: valTypes[];
}

export interface editType {
	level?: number;
	pre?: number;
}
