// rle.cs
using System;
using System.IO;
using System.Text;

class RLE
{
    static string Colorize(string text, string color)
    {
        string col = color switch
        {
            "green" => "\x1b[92m",
            "red" => "\x1b[91m",
            "yellow" => "\x1b[93m",
            _ => "\x1b[0m"
        };
        return col + text + "\x1b[0m";
    }

    const char ESCAPE = '\\';

    static string Compress(string text)
    {
        if (string.IsNullOrEmpty(text)) return "";
        var result = new StringBuilder();
        int i = 0, n = text.Length;
        while (i < n)
        {
            char ch = text[i];
            int j = i + 1;
            while (j < n && text[j] == ch) j++;
            int count = j - i;
            if (count >= 3)
            {
                result.Append(ch);
                result.Append(ESCAPE);
                result.Append(count);
            }
            else
            {
                for (int k = 0; k < count; k++)
                {
                    if (ch == ESCAPE) result.Append(ESCAPE);
                    result.Append(ch);
                }
            }
            i = j;
        }
        return result.ToString();
    }

    static string Decompress(string text)
    {
        if (string.IsNullOrEmpty(text)) return "";
        var result = new StringBuilder();
        int i = 0, n = text.Length;
        while (i < n)
        {
            char ch = text[i];
            if (ch == ESCAPE)
            {
                if (i + 1 < n && text[i + 1] == ESCAPE)
                {
                    result.Append(ESCAPE);
                    i += 2;
                    continue;
                }
                if (i + 1 >= n) throw new Exception("Unexpected end after escape");
                char repeatChar = text[i + 1];
                i += 2;
                string numStr = "";
                while (i < n && char.IsDigit(text[i]))
                {
                    numStr += text[i];
                    i++;
                }
                if (string.IsNullOrEmpty(numStr)) throw new Exception("Missing number after escape");
                int count = int.Parse(numStr);
                result.Append(repeatChar, count);
            }
            else
            {
                result.Append(ch);
                i++;
            }
        }
        return result.ToString();
    }

    static string ReadInput(string filename)
    {
        if (string.IsNullOrEmpty(filename) || filename == "-")
        {
            return Console.In.ReadToEnd();
        }
        return File.ReadAllText(filename, Encoding.UTF8);
    }

    static void WriteOutput(string filename, string content)
    {
        if (string.IsNullOrEmpty(filename) || filename == "-")
        {
            Console.Write(content);
        }
        else
        {
            File.WriteAllText(filename, content, Encoding.UTF8);
        }
    }

    static void Main(string[] args)
    {
        if (args.Length < 1)
        {
            Console.WriteLine(Colorize("Usage: rle compress|decompress [input] [output]", "yellow"));
            return;
        }
        string mode = args[0];
        if (mode != "compress" && mode != "decompress")
        {
            Console.WriteLine(Colorize("Invalid mode. Use compress or decompress.", "red"));
            return;
        }
        string inputFile = args.Length >= 2 ? args[1] : "";
        string outputFile = args.Length >= 3 ? args[2] : "";

        string data;
        try
        {
            data = ReadInput(inputFile);
        }
        catch (Exception e)
        {
            Console.WriteLine(Colorize("Error reading input: " + e.Message, "red"));
            return;
        }

        string result;
        try
        {
            if (mode == "compress") result = Compress(data);
            else result = Decompress(data);
        }
        catch (Exception e)
        {
            Console.WriteLine(Colorize("Error: " + e.Message, "red"));
            return;
        }

        try
        {
            WriteOutput(outputFile, result);
            if (!string.IsNullOrEmpty(outputFile) && outputFile != "-")
            {
                Console.WriteLine(Colorize("Result written to " + outputFile, "green"));
            }
        }
        catch (Exception e)
        {
            Console.WriteLine(Colorize("Error writing output: " + e.Message, "red"));
        }
    }
}
