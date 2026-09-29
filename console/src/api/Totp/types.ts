export interface totpConfigTypes {
	phone: string;
	secret: string;
}

export interface bindTotpTypes {
	code: string | undefined;
	secret: string;
}
