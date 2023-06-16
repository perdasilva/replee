const go = new Go();
let mod, inst;
const modeStartOfInput = "startOfInput";
const modeMidInput = "midInput";
const modeReset = "reset";
WebAssembly.instantiateStreaming(fetch('js/main.wasm'), go.importObject).then(async (result) => {
    mod = result.module;
    inst = result.instance;
    go.run(inst);

    var term = new Terminal({
        theme: {
            foreground: '#00FF00', // Matrix green
        },
        cursorBlink: true,
        cursorStyle: 'bar',
        fontFamily: 'monospace',
        fontSize: 16,
        fontWeight: 'normal',
    });
    var fitAddon = new FitAddon.FitAddon();
    term.loadAddon(fitAddon);
    term.open(document.getElementById('terminal-container'));
    term.focus()
    fitAddon.fit();

    var commandHistory = [];
    var commandIndex = -1;
    var currentLine = '';
    var cursorPosition = 0;
    var mode = modeStartOfInput;

    function printPrompt() {
        term.write('\r\033[Kreplee: ' + currentLine);
        term.write('\033[' + (currentLine.length - (cursorPosition + 1)) + 'D'); // Move the cursor back to the correct position
    }

    function handleResponse(response) {
        if (response.mode === modeStartOfInput) {
            mode = modeStartOfInput;
            term.writeln('');
            let command = currentLine;
            commandHistory.push(command);
            commandIndex = commandHistory.length;
            let request = JSON.stringify({
                mode: modeStartOfInput,
                input: command,
                indent: 0,
            });
            let result = window.replee(request);
            handleResponse(result);
        } else if (response.mode === modeMidInput) {
            mode = modeMidInput;
            term.writeln('');
            let lines = response.output.split('\n');
            for (let line of lines) {
                term.writeln(line);
            }
            currentLine = ' '.repeat(response.indent);
            cursorPosition = currentLine.length;
            printPrompt();
        } else {
            mode = modeStartOfInput;
            // term.writeln('');
            let lines = response.output.split('\n');
            if (response.isErr === true) {
                term.write('\x1b[31m'); // Set the color to red
            } else {
                term.write('\x1b[36m'); // Set the grey to green
            }
            for (let line of lines) {// Write the color escape sequence
                if (line !== "") {
                    term.writeln(line);
                }
            }
            term.write('\x1b[0m'); // Reset the color
            currentLine = '';
            cursorPosition = 0;
            printPrompt();
        }
    }

    printPrompt();

    let isProcessing = false
    term.onKey(e => {
        if (isProcessing) {
            return;
        }
        isProcessing = true;
        const printable = !e.domEvent.altKey && !e.domEvent.altGraphKey && !e.domEvent.ctrlKey && !e.domEvent.metaKey;
        if (e.domEvent.keyCode === 13) { // Enter key
            if (mode === modeStartOfInput) {
                term.writeln('');
                let command = currentLine;
                commandHistory.push(command);
                commandIndex = commandHistory.length;
                let request = JSON.stringify({
                    mode: modeStartOfInput,
                    input: command,
                    indent: 0,
                });
                let result = window.replee(request);
                handleResponse(result);
            } else if (mode === modeMidInput) {
                currentLine += '\n';
                cursorPosition++;
                printPrompt();
            }
        } else if (e.domEvent.keyCode === 38) { // Up arrow key
            if (commandIndex > 0) {
                commandIndex--;
                currentLine = commandHistory[commandIndex];
                cursorPosition = currentLine.length;
                printPrompt();
            }
        } else if (e.domEvent.keyCode === 40) { // Down arrow key
            if (commandIndex < commandHistory.length - 1) {
                commandIndex++;
                currentLine = commandHistory[commandIndex];
                cursorPosition = currentLine.length;
                printPrompt();
            }
        } else if (e.domEvent.keyCode === 8) { // Backspace key
            if (cursorPosition > 0) {
                currentLine = currentLine.substring(0, cursorPosition - 1) + currentLine.substring(cursorPosition);
                cursorPosition--;
                printPrompt();
            }
        } else if (e.domEvent.keyCode === 37) { // Left arrow key
            if (cursorPosition > 0) {
                cursorPosition--;
                printPrompt();
            }
        } else if (e.domEvent.keyCode === 39) { // Right arrow key
            if (cursorPosition < currentLine.length) {
                cursorPosition++;
                printPrompt();
            }
        } else if (printable) {
            currentLine = currentLine.substring(0, cursorPosition) + e.key + currentLine.substring(cursorPosition);
            cursorPosition++;
            printPrompt();
        }
        isProcessing = false;
    });
});