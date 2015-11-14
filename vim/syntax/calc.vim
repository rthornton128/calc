" Copyright (c) 2014, Rob Thornton
" All rights reserved.
" This source code is governed by a Simplied BSD-License. Please see the
" LICENSE included in this distribution for a copy of the full license
" or, if one is not included, you may also find a copy at
" http://opensource.org/licenses/BSD-2-Clause
"
" Vim Highlighting for Calc 3.0

if exists("b:current_syntax")
	finish
endif

syn case match

" keywords
syn keyword calcStatement define
syn keyword calcConditional if
syn keyword calcExpression func var

hi def link calcStatement Statement
hi def link calcConditional Conditional
hi def link calcExpression Keyword

" Predeclared types
syn keyword calcType bool int

hi def link calcType Type

" Basic literals
syn keyword calcBoolean false true
syn match calcInteger /\d*/

hi def link calcBoolean Boolean
hi def link calcInteger Number

" Comments
syn keyword calcTODO contained TODO FIXME BUG
syn region calcComment start=";" end="$" contains=@calcTODO,@Spell

hi def link calcTODO Todo
hi def link calcComment Comment

" Regions
syn region calcParen start="(" end=")" transparent

" Identifiers
" syn match calcIdentifier /\a[\a\d]*/

" hi def link calcIdentifier Identifier

" Operators
syn match calcOperator "+"
syn match calcOperator "-"
syn match calcOperator "/"
syn match calcOperator "*"
syn match calcOperator "%"
syn match calcOperator "="
syn match calcOperator "=="
syn match calcOperator "!="
syn match calcOperator "<"
syn match calcOperator "<="
syn match calcOperator ">"
syn match calcOperator ">="

syn match calcSpecial ":"

hi def link calcOperator Operator
hi def link calcSpecial Special

let b:current_syntax = "calc"
