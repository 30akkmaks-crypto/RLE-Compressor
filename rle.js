// rle.js
#!/usr/bin/env node
'use strict';

const fs = require('fs');

// ANSI colors
const COLORS = {
    green: '\x1b[92m',
    red: '\x1b[91m',
    yellow: '\x1b[93m',
    reset: '\x1b[0m'
};

function colorize(text, color) {
    return COLORS[color] + text + COLORS.reset;
}

const ESCAPE = '\\';

function compress(text) {
    if (!text) return '';
    let result = [];
    let i = 0, n = text.length;
    while (i < n) {
        let ch = text[i];
        let j = i + 1;
        while (j < n && text[j] === ch) j++;
        let count = j - i;
        if (count >= 3) {
            result.push(ch, ESCAPE, count.toString());
        } else {
            for (let k = 0; k < count; k++) {
                if (ch === ESCAPE) result.push(ESCAPE);
                result.push(ch);
            }
        }
        i = j;
    }
    return result.join('');
}

function decompress(text) {
    if (!text) return '';
    let result = [];
    let i = 0, n = text.length;
    while (i < n) {
        let ch = text[i];
        if (ch === ESCAPE) {
            if (i + 1 < n && text[i + 1] === ESCAPE) {
                result.push(ESCAPE);
                i += 2;
                continue;
            }
            if (i + 1 >= n) throw new Error('Unexpected end after escape');
            let char = text[i + 1];
            i += 2;
            let numStr = '';
            while (i < n && text[i] >= '0' && text[i] <= '9') {
                numStr += text[i];
                i++;
            }
            if (!numStr) throw new Error('Missing number after escape');
            let count = parseInt(numStr, 10);
            for (let k = 0; k < count; k++) result.push(char);
        } else {
            result.push(ch);
            i++;
        }
    }
    return result.join('');
}

function readInput(filename) {
    if (!filename || filename === '-') {
        // read from stdin
        return fs.readFileSync(0, 'utf-8');
    }
    return fs.readFileSync(filename, 'utf-8');
}

function writeOutput(filename, content) {
    if (!filename || filename === '-') {
        process.stdout.write(content);
    } else {
        fs.writeFileSync(filename, content, 'utf-8');
    }
}

function main() {
    const args = process.argv.slice(2);
    if (args.length < 1) {
        console.log(colorize('Usage: node rle.js compress|decompress [input] [output]', 'yellow'));
        process.exit(1);
    }
    const mode = args[0];
    if (mode !== 'compress' && mode !== 'decompress') {
        console.log(colorize('Invalid mode. Use compress or decompress.', 'red'));
        process.exit(1);
    }
    const inputFile = args.length >= 2 ? args[1] : undefined;
    const outputFile = args.length >= 3 ? args[2] : undefined;

    let data;
    try {
        data = readInput(inputFile);
    } catch (err) {
        console.log(colorize('Error reading input: ' + err.message, 'red'));
        process.exit(1);
    }

    let result;
    try {
        if (mode === 'compress') result = compress(data);
        else result = decompress(data);
    } catch (err) {
        console.log(colorize('Error: ' + err.message, 'red'));
        process.exit(1);
    }

    try {
        writeOutput(outputFile, result);
        if (outputFile && outputFile !== '-') {
            console.log(colorize('Result written to ' + outputFile, 'green'));
        }
    } catch (err) {
        console.log(colorize('Error writing output: ' + err.message, 'red'));
        process.exit(1);
    }
}

main();
