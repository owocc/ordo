import { TableConstructorOptions } from 'npm:cli-table3'

export const blankBorder: TableConstructorOptions['chars'] = {
    'top': '', 'top-mid': '', 'top-left': '', 'top-right': ''
    , 'bottom': '', 'bottom-mid': '', 'bottom-left': '', 'bottom-right': ''
    , 'left': '', 'left-mid': '', 'mid': '', 'mid-mid': ''
    , 'right': '', 'right-mid': '', 'middle': ' '
}