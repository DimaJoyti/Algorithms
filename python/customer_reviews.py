#!/bin/python3

import math
import os
import random
import re
import sys


#
# Complete the 'searchSuggestions' function below.
#
# The function is expected to return a 2D_STRING_ARRAY.
# The function accepts following parameters:
#  1. STRING_ARRAY repository
#  2. STRING customerQuery
#

def searchSuggestions(repository, customerQuery):
    # Write your code here

    def suggestedProducts(self, products: List[str], searchWord: str) -> List[List[str]]:
        _END = '_END'
        trie = {}

        def build(word):
            curr = trie
            for c in word:
                curr = curr.setdefault(c, {})
            curr[_END] = _END
        for word in set(products):
            build(word)
        res = []
        curr = trie

        def dfs(curr, word):
            if len(words) == 3:
                return
            if _END in curr:
                words.append(word)
            if len(words) == 3:
                return
            for key in sorted(curr.keys()):
                if key != _END:
                    dfs(curr[key], word+key)

        w = ""
        for c in searchWord:
            w += c
            words = []
            if c not in curr:
                break
            dfs(curr[c], w)
            if not words:
                break
            res.append(words)
            curr = curr[c]
        return res + [[] for _ in range(len(searchWord)-len(res))]
