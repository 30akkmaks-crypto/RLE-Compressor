// rle.cpp
#include <iostream>
#include <fstream>
#include <string>
#include <sstream>
#include <cctype>
#include <stdexcept>

using namespace std;

const string RESET = "\033[0m";
const string GREEN = "\033[92m";
const string RED = "\033[91m";
const string YELLOW = "\033[93m";

string colorize(const string& text, const string& color) {
    return color + text + RESET;
}

const char ESCAPE = '\\';

string compress(const string& text) {
    if (text.empty()) return "";
    ostringstream result;
    size_t i = 0, n = text.size();
    while (i < n) {
        char ch = text[i];
        size_t j = i + 1;
        while (j < n && text[j] == ch) ++j;
        size_t count = j - i;
        if (count >= 3) {
            result << ch << ESCAPE << count;
        } else {
            for (size_t k = 0; k < count; ++k) {
                if (ch == ESCAPE) result << ESCAPE;
                result << ch;
            }
        }
        i = j;
    }
    return result.str();
}

string decompress(const string& text) {
    if (text.empty()) return "";
    ostringstream result;
    size_t i = 0, n = text.size();
    while (i < n) {
        char ch = text[i];
        if (ch == ESCAPE) {
            if (i + 1 < n && text[i + 1] == ESCAPE) {
                result << ESCAPE;
                i += 2;
                continue;
            }
            if (i + 1 >= n) throw runtime_error("Unexpected end after escape");
            char repeat_char = text[i + 1];
            i += 2;
            string num_str;
            while (i < n && isdigit(text[i])) {
                num_str += text[i];
                ++i;
            }
            if (num_str.empty()) throw runtime_error("Missing number after escape");
            int count = stoi(num_str);
            for (int k = 0; k < count; ++k) result << repeat_char;
        } else {
            result << ch;
            ++i;
        }
    }
    return result.str();
}

string readFile(const string& filename) {
    if (filename == "-" || filename.empty()) {
        // read stdin
        stringstream buffer;
        buffer << cin.rdbuf();
        return buffer.str();
    }
    ifstream f(filename);
    if (!f) throw runtime_error("Cannot open file: " + filename);
    stringstream buffer;
    buffer << f.rdbuf();
    return buffer.str();
}

void writeFile(const string& filename, const string& content) {
    if (filename == "-" || filename.empty()) {
        cout << content;
        return;
    }
    ofstream f(filename);
    if (!f) throw runtime_error("Cannot write file: " + filename);
    f << content;
}

int main(int argc, char* argv[]) {
    if (argc < 2) {
        cout << colorize("Usage: rle compress|decompress [input] [output]", YELLOW) << endl;
        return 1;
    }
    string mode = argv[1];
    if (mode != "compress" && mode != "decompress") {
        cout << colorize("Invalid mode. Use compress or decompress.", RED) << endl;
        return 1;
    }
    string inputFile = (argc >= 3) ? argv[2] : "";
    string outputFile = (argc >= 4) ? argv[3] : "";

    string data;
    try {
        data = readFile(inputFile);
    } catch (const exception& e) {
        cout << colorize("Error reading input: " + string(e.what()), RED) << endl;
        return 1;
    }

    string result;
    try {
        if (mode == "compress") result = compress(data);
        else result = decompress(data);
    } catch (const exception& e) {
        cout << colorize("Error: " + string(e.what()), RED) << endl;
        return 1;
    }

    try {
        writeFile(outputFile, result);
        if (!outputFile.empty() && outputFile != "-") {
            cout << colorize("Result written to " + outputFile, GREEN) << endl;
        }
    } catch (const exception& e) {
        cout << colorize("Error writing output: " + string(e.what()), RED) << endl;
        return 1;
    }
    return 0;
}
