import Conf from "npm:conf";
import type { Config } from "../types/config.ts";

export const config = new Conf<Config>({
    projectName: "ordo",
});
