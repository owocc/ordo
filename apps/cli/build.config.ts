import { defineBuildConfig } from "unbuild";
export default defineBuildConfig({
	entries:[
		"src/main.tsx"
	],
	outDir:"dist",
	declaration:'compatible'
})
