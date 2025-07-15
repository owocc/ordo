import { defineBuildConfig } from "unbuild";
export default defineBuildConfig({
	entries:[
		"src/index.ts"
	],
	outDir:"dist",
	declaration:'compatible'
})
