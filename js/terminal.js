const go = new Go();
let mod, inst;
WebAssembly.instantiateStreaming(fetch('js/main.wasm'), go.importObject).then(async (result) => {
    mod = result.module;
    inst = result.instance;
    go.run(inst);

    var term = new Terminal();
    var fitAddon = new FitAddon.FitAddon();
    term.loadAddon(fitAddon);
    term.open(document.getElementById('terminal-container'));
    fitAddon.fit();

    var commandHistory = [];
    var commandIndex = -1;
    var currentLine = '';
    var cursorPosition = 0;

    term.write('replee > ');

    term.onKey(e => {
        const printable = !e.domEvent.altKey && !e.domEvent.altGraphKey && !e.domEvent.ctrlKey && !e.domEvent.metaKey;
        if (e.domEvent.keyCode === 13) { // Enter key
            term.writeln('');
            let command = currentLine;
            commandHistory.push(command);
            commandIndex = commandHistory.length;// Assuming replee is the function exposed by your wasm library
            let result = window.replee(command); // Assuming replee is the function exposed by your wasm library
            if (result.startsWith('Error: Unexpected end of input')) {
                // Add an indentation to the next line
                currentLine = '    ';
                cursorPosition = 4;
            } else {
                let lines = result.split('\n');
                for (let line of lines) {
                    term.writeln(line);
                }
                currentLine = '';
                cursorPosition = 0;
            }
        } else if (e.domEvent.keyCode === 38) { // Up arrow key
            if (commandIndex > 0) {
                commandIndex--;
                currentLine = commandHistory[commandIndex];
                cursorPosition = currentLine.length;
                term.write('\r\033[Kreplee > ' + currentLine);
            }
        } else if (e.domEvent.keyCode === 40) { // Down arrow key
            if (commandIndex < commandHistory.length - 1) {
                commandIndex++;
                currentLine = commandHistory[commandIndex];
                cursorPosition = currentLine.length;
                term.write('\r\033[Kreplee > ' + currentLine);
            }
        } else if (e.domEvent.keyCode === 8) { // Backspace key
            if (cursorPosition > 0) {
                currentLine = currentLine.substring(0, cursorPosition - 1) + currentLine.substring(cursorPosition);
                cursorPosition--;
                term.write('\r\033[Kreplee > ' + currentLine);
                term.write('\033[' + (currentLine.length - cursorPosition) + 'D'); // Move the cursor back to the correct position
            }
        } else if (e.domEvent.keyCode === 37) { // Left arrow key
            if (cursorPosition > 0) {
                cursorPosition--;
                term.write('\033[D'); // Move the cursor to the left
            }
        } else if (e.domEvent.keyCode === 39) { // Right arrow key
            if (cursorPosition < currentLine.length) {
                cursorPosition++;
                term.write('\033[C'); // Move the cursor to the right
            }
        } else if (printable) {
            currentLine = currentLine.substring(0, cursorPosition) + e.key + currentLine.substring(cursorPosition);
            cursorPosition++;
            term.write('\r\033[Kreplee > ' + currentLine);
            term.write('\033[' + (currentLine.length - cursorPosition) + 'D'); // Move the cursor back to the correct position
        }
    });
});
