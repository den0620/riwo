;
; (C) CoffeeLake 2025
; This is a test file, needed for me to see in Monaco editor
; 

section .text
_start:
    mov rax, 0x10000000
    mov rbx, 0x113345ff
    
    call NextProcedureLabel
    ret

NextProcedureLabel:
    add rax, rbx
    ret

section .data
    some_data db 'Interesting information'