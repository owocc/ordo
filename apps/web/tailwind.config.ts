import type { Config } from 'tailwindcss'
import { iconsPlugin, getIconCollections } from '@egoist/tailwindcss-icons'

const config: Config = {
    plugins: [
        iconsPlugin({
            collections: getIconCollections(['lucide'])
        })
    ]
}

export default config