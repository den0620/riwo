/*----------------------------------
    (C) CoffeeLake 2025

    Common Assembler Syntax plugin (CAS)
includes all known machine words for highlighting
translated opcode in Sunflower Monaco window

    It includes all architecture keywords: 
        IA-32: mov call ax bx 
        Aplha: ld r2 r1
        ...
    
    SunFlower.Translator.dll (F# Core) build
uses for translate opcodes to CAS file (raw.asm) and this part
is highlights code in file.
-------------------------------------*/
require(['vs/editor/editor.main'], function() {
    // call Monaco API
    monaco.languages.register({ id: 'cas' });

    monaco.languages.setMonarchTokensProvider('cas', {
        keywords: [
            'mov', 'add', 'sub', 'jmp', 'call', 'callf', 'ret', 'push', 'pop', 'cmp', 'je', 'jne',
            'db', 'dw', 'dd', 'equ', 'times', 'resb', 'resw', 'resd', 'section', 'global'
        ],
        registers: [
            'eax', 'ebx', 'ecx', 'edx', 'esi', 'edi', 'esp', 'ebp',
            'ax', 'bx', 'cx', 'dx', 'al', 'ah', 'bl', 'bh', 'cl', 'ch', 'dl', 'dh'
        ],
        directives: [
            'bits', 'org', '%include', '%define', '%macro', '%endmacro'
        ],
        operators: /[\[\],:+*\-<>]/,
        symbols: /[=!]/,

        tokenizer: {
            root: [
                [/;.*$/, 'comment'],

                [/%[a-z]+/, 'directive'],

                [/\b0x[\da-f]+\b/i, 'number.hex'],
                [/\b0b[01]+\b/, 'number.bin'],
                [/\b\d+\b/, 'number'],

                [/\b(a|b|c|d|e|s|i|p)[a-z]{1,3}\b/, {
                    cases: { '@registers': 'keyword' }
                }],

                [/[a-z_][\w$]*:/, 'type.identifier'],

                [/[a-z_][\w$]*/, {
                    cases: {
                        '@keywords': 'keyword',
                        '@default': 'identifier'
                    }
                }],

                [/"([^"\\]|\\.)*$/, 'string.invalid'],
                [/"/, 'string', '@string']
            ],
            string: [
                [/[^\\"]+/, 'string'],
                [/\\./, 'string.escape'],
                [/"/, 'string', '@pop']
            ]
        }
    });

    const editor = monaco.editor.create(document.getElementById('container'), {
        language: 'cas',
    });
});