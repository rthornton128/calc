.global _main
.global main
.data
fmt: .asciz "%d\12"

.text
_main:
push %ebp
movl %esp, %ebp
subl $16, %esp
call main
movl %eax, 4(%esp)
movl $fmt, (%esp)
call _printf
movl $0, %eax
leave
ret
main:
movl then2, %ebx
movl $true, %eax
then2:
movl $99, %eax
ret
