if exists('g:loaded_nvim_screenshot')
    finish
endif

let g:loaded_nvim_screenshot = 1

function! s:panic(ch, data, ...) abort
    echom a:data
endfunction

function! s:Start_nvim_screenshot(host) abort
    let binary = nvim_get_runtime_file('main', v:false)[0]
    return jobstart([binary], {
                \ 'rpc': v:true,
                \ 'on_stderr': function('s:panic')
                \})
endfunction

call remote#host#Register('remote', 'x', function('s:Start_nvim_screenshot'))

call remote#host#RegisterPlugin('remote', '0', [
        \ {
            \ 'type': 'command', 
            \ 'name': 'Screenshot', 
            \ 'sync': 1, 
            \ 'opts': {
                \ 'nargs': '?', 
                \ 'range': '', 
                \ 'eval': '[getline(1, "$"),GetSyntax(),synIDattr(hlID("Normal"), "fg"),synIDattr(hlID("Normal"),"bg")]'
                \ }},
        \ ])

function! GetSyntax()

    let s:defaultfg = synIDattr(hlID("Normal"), "fg")
    let s:defaultbg = synIDattr(hlID("Normal"), "bg")

    if s:defaultfg == "" | let s:defaultfg = ( &background == "dark" ? "#ffffff" : "#333333" ) | endif
    if s:defaultbg == "" | let s:defaultbg = ( &background == "dark" ? "#333333" : "#ffffff" ) | endif

    let s:lines = []
    let s:lnum = 1
    let s:linecount = line("$")

    while s:lnum <= s:linecount

        let s:res_line = ""

        let s:line = getline(s:lnum)
        let s:len = strlen(s:line)
        let s:col = 1

        while s:col <= s:len 
            let s:id = synID(s:lnum, s:col, 1)
            let s:fg = synIDattr(synIDtrans(s:id), "fg#")
            let s:bg = synIDattr(synIDtrans(s:id), "bg#")

            if s:fg == "" | let s:fg = s:defaultfg | endif
            if s:bg == "" | let s:bg = s:defaultbg | endif

            let s:res_line .= printf("%s:%s;", s:col-1, s:fg)

            let s:concealinfo = synconcealed(s:lnum, s:col)
            while s:col <= s:len && s:concealinfo == synconcealed(s:lnum, s:col)
                let s:col = s:col + 1 
            endwhile
        endwhile

        let s:lines = add(s:lines, s:res_line[:-2])
        let s:lnum = s:lnum + 1
    endwhile

    return s:lines
endfunction
