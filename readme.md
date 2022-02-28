# nvim-screenshot

Neovim Plugin to screenshot lines of code.

# Install

Add the plugin to Neovim and compile the go binary

vim-plug example: 
`Plug 'asamonik/nvim-screenshot', { 'do': 'go build' }`

# basic usage

Call `:Screenshot` without arguments to copy a screenshot of the visually selected 
lines to clipboard.

Call `:Screenshot ./path/` to save them as a png file under the specified path

# TODO
* de-indent code in screenshot
* Bold, Italic and Underline 
