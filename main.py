from functools import partial
import numpy as np
import struct
import os
import math

code = {
    '0':'0',
    '1':'1',
    '2':'2',
    '3':'3',
    '4':'4',
    '5':'5',
    '6':'6',
    '7':'7',
    '8':'8',
    '9':'9',
    '11':'A',
    '12':'B',
    '13':'C',
    '14':'D',
    '15':'E',
    '16':'F',
    '17':'G',
    '18':'H',
    '19':'I',
    '20':'J',
    '21':'K',
    '22':'L',
    '23':'M',
    '24':'N',
    '25':'O',
    '26':'P',
    '27':'Q',
    '28':'R',
    '29':'S',
    '30':'T',
    '31':'U',
    '32':'V',
    '33':'W',
    '34':'X',
    '35':'Y',
    '36':'Z',
    '41':'||',
    '42':'>>',
    '43':'<<',
    '44':'*',
    '45':'#',
    '46':'[]',
    '47':'()',
    '48':'+',
    '49':'-',
    '50':'=',
    '51':'a',
    '52':'b',
    '53':'c',
    '54':'d',
    '55':'e',
    '56':'f',
    '57':'g',
    '58':'h',
    '59':'i',
    '60':'j',
    '61':'k',
    '62':'l',
    '63':'m',
    '64':'n',
    '65':'o',
    '66':'p',
    '67':'q',
    '68':'r',
    '69':'s',
    '70':'t',
    '71':'e',
    '72':'v',
    '73':'w',
    '74':'x',
    '75':'y',
    '76':'z',
    '127':'□'
}

def get_code_key(key):
    if key in code:
        return code[key]
    else:
        return '□'


def get_code_val(val):
    if val in code.values():
        return list(code.keys())[list(code.values()).index(val)]
    else:
        return '127'

def get_bit_val(byte, index):
    """
    得到某个字节中某一位（Bit）的值

    :param byte: 待取值的字节值
    :param index: 待读取位的序号，从右向左0开始，0-7为一个完整字节的8个位
    :returns: 返回读取该位的值，0或1
    """
    if byte & (1 << index):
        return 1
    else:
        return 0

def set_bit_val(byte, index, val):
    """
    更改某个字节中某一位（Bit）的值

    :param byte: 准备更改的字节原值
    :param index: 待更改位的序号，从右向左0开始，0-7为一个完整字节的8个位
    :param val: 目标位预更改的值，0或1
    :returns: 返回更改后字节的值
    """
    if val:
        return byte | (1 << index)
    else:
        return byte & ~(1 << index)

def converting(source_num, source_hex, target_hex):
    # （2， 36）之间的进制转换
    if source_hex > 36 or source_hex < 2:
        return '2 <= source_hex <= 36'
    if target_hex > 36 or target_hex < 2:
        return '2 <= target_hex <= 36'
    str_36 = '0123456789abcdefghijklmnopqrstuvwxyz'
    dict_36 = {}
    for i in range(len(str_36)):
        dict_36[str_36[i]] = i
    str_b = str_36[:target_hex]
    result = ''
    source_str = str(source_num).lower()
    decimal_num = 0
    for i in range(len(source_str)):
        decimal_num += dict_36[source_str[-i-1]] * (source_hex ** i)
    quotient_int = decimal_num
    while quotient_int:
        remainder = quotient_int % target_hex
        quotient_int = quotient_int // target_hex
        result = str_b[remainder] + result
        if quotient_int and quotient_int < target_hex:
            result = str_b[quotient_int] + result
            break
    return result

def ternary_encode(str1):
    ternary = ''
    for s in str1:
        n = get_code_val(s)
        ternary += converting(n,10,2).rjust(7,'0')
    return ternary


def ternary_decode(bin1):
    source = ''
    lenbin1 = len(bin1)
    if lenbin1%7 == 0:
        n = math.ceil(lenbin1/7)
        for i in range(n):
            k = converting(bin1[7*i:7*(i+1)],2,10)
            if k == '':
                k = '0'
            source += get_code_key(k)
    return source

def save_bin(s):
    bin1 = ternary_encode(s)
    n = math.ceil(len(bin1)/8)
    ba = bytearray(n)
    for i in range(n):
        for j in range(7,-1,-1):
            if len(bin1)>int(i*8+7-j):
                ba[i] = set_bit_val(ba[i],j,int(bin1[i*8+7-j]))
            else:
                ba[i] = set_bit_val(ba[i],j,0)
    return ba

def read_bin(bs):
    i=0
    bin2 = ''
    while i<len(bs):
        byte = bs[i]
        b = ''
        for j in range(7,-1,-1):
            b+=str(get_bit_val(byte,j))
        bin2+=b
        i+=1
    n = len(bin2)%7
    return ternary_decode(bin2[:len(bin2)-n])
            

if __name__ == "__main__":
    fd = open('test.bin','rb')
    bs = fd.read()
    str1 = "A23456789"
    ## bin1 = '00000000000001000001000000110110011011010001101010110110'
    bs1 = save_bin(str1)
    print(bs1)
    rb = read_bin(bs1)
    print(rb)