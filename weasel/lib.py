import ast

from weasel.transformers.rule_1 import Rule_1


def optimize(root: ast.AST) -> ast.AST:
    root = Rule_1().visit(root)
    return root
