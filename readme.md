# nvim-screenshot

Neovim Plugin to screenshot lines of code.

## Install

Add the plugin to Neovim and compile the go binary

vim-plug example:<br> 
`Plug 'asamonik/nvim-screenshot', { 'do': 'go build' }`

Packer:<br>
`https://github.com/yolofanhd/nvim-screenshot.git`

## basic usage

Call `:Screenshot` without arguments to copy a screenshot of the visually selected 
lines to clipboard.

Call `:Screenshot ./path/` to save selected lines as a .png file under the specified path

# TODO
* de-indent code in screenshot
* Bold, Italic and Underline 
