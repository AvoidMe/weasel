import ast
from typing import Union
from itertools import zip_longest


def compare_ast(
    node1: Union[ast.expr, list[ast.expr]],
    node2: Union[ast.expr, list[ast.expr]],
) -> bool:
    if type(node1) is not type(node2):
        return False

    if isinstance(node1, ast.AST):
        for k, v in vars(node1).items():
            if k in {
                "lineno",
                "end_lineno",
                "col_offset",
                "end_col_offset",
                "ctx",
            }:
                continue
            if not compare_ast(v, getattr(node2, k)):
                return False
        return True

    elif isinstance(node1, list) and isinstance(node2, list):
        return all(compare_ast(n1, n2) for n1, n2 in zip_longest(node1, node2))
    return node1 == node2
