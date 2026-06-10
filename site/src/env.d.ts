/// <reference types="astro/client" />

declare module "*.yaml?raw" {
	const content: string;
	export default content;
}

declare module "virtual:hev-ask/digest" {
	const digest: any;
	export default digest;
}
