# rle.py
#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import sys
import os
import argparse

# ANSI colors
COLORS = {
    'green': '\033[92m',
    'red': '\033[91m',
    'yellow': '\033[93m',
    'reset': '\033[0m'
}

def colorize(text, color):
    return f"{COLORS.get(color, '')}{text}{COLORS['reset']}"

ESCAPE = '\\'

def compress(text):
    """Сжатие RLE с escape-символом."""
    if not text:
        return ''
    result = []
    i = 0
    n = len(text)
    while i < n:
        ch = text[i]
        j = i + 1
        while j < n and text[j] == ch:
            j += 1
        count = j - i
        if count >= 3:
            result.append(ch)
            result.append(ESCAPE)
            result.append(str(count))
        else:
            # Для каждого символа из серии
            for _ in range(count):
                if ch == ESCAPE:
                    result.append(ESCAPE)  # дублируем escape
                result.append(ch)
        i = j
    return ''.join(result)

def decompress(text):
    """Распаковка RLE."""
    if not text:
        return ''
    result = []
    i = 0
    n = len(text)
    while i < n:
        ch = text[i]
        if ch == ESCAPE:
            # Проверяем, что это не экранированный escape
            if i + 1 < n and text[i + 1] == ESCAPE:
                # Это просто символ '\'
                result.append(ESCAPE)
                i += 2
                continue
            # Иначе это управляющая последовательность: символ + число
            if i + 1 >= n:
                raise ValueError("Неожиданный конец строки после escape-символа")
            char = text[i + 1]
            i += 2
            # Считываем число
            num_str = ''
            while i < n and text[i].isdigit():
                num_str += text[i]
                i += 1
            if not num_str:
                raise ValueError("Отсутствует число после escape-последовательности")
            count = int(num_str)
            result.append(char * count)
        else:
            result.append(ch)
            i += 1
    return ''.join(result)

def main():
    parser = argparse.ArgumentParser(description="RLE Compressor/Decompressor")
    parser.add_argument('mode', choices=['compress', 'decompress'], help='Режим работы')
    parser.add_argument('input', nargs='?', help='Входной файл (если не указан, читается stdin)')
    parser.add_argument('output', nargs='?', help='Выходной файл (если не указан, выводится stdout)')
    args = parser.parse_args()

    # Чтение входных данных
    if args.input:
        try:
            with open(args.input, 'r', encoding='utf-8') as f:
                data = f.read()
        except Exception as e:
            sys.exit(colorize(f"Ошибка чтения файла: {e}", 'red'))
    else:
        # Чтение из stdin
        data = sys.stdin.read()
        if not data:
            sys.exit(colorize("Нет входных данных", 'red'))

    # Обработка
    try:
        if args.mode == 'compress':
            result = compress(data)
        else:
            result = decompress(data)
    except Exception as e:
        sys.exit(colorize(f"Ошибка: {e}", 'red'))

    # Вывод результата
    if args.output:
        try:
            with open(args.output, 'w', encoding='utf-8') as f:
                f.write(result)
            print(colorize(f"Результат записан в {args.output}", 'green'))
        except Exception as e:
            sys.exit(colorize(f"Ошибка записи файла: {e}", 'red'))
    else:
        sys.stdout.write(result)

if __name__ == '__main__':
    main()
